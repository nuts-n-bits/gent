// Kawas = かわす, lit. message exchanges

import { CcCore, Command, parse_one } from "../cc/backing-codec.js";
import { readline, Line } from "./codec.js";


class KawasError extends Error {}
class NeverError extends Error {}

type StreamListSubscriber<T> = {
        push: StreamListPushHandler<T>,
        end: StreamListEndHandler<T>,
        index: number,
};

type StreamListPushHandler<T> = (latest_item: T, list: T[], fresh_idx_begin: number) => void;
type StreamListEndHandler<T> = (list: T[]) => void;

class StreamList<a> implements AsyncIterable<a> {
        backing: a[] = [];
        subscribers: StreamListSubscriber<a>[] = [];
        ended = false
        push(a: a): void {
                if (this.ended) {
                        throw new Error("cannot push to an ended stream");
                }
                this.backing.push(a);
                for (const sub of this.subscribers) {
                        sub.push(a, this.backing, sub.index);
                        sub.index = this.backing.length;
                }
        }
        end(): void {
                if (this.ended) {
                        return;
                }
                this.ended = true;
                for (const sub of this.subscribers) {
                        sub.end(this.backing);
                }
        }
        get_backing(): a[] {
                return this.backing;
        }
        subscribe(push_handler: StreamListPushHandler<a>, end_handler: StreamListEndHandler<a>): void {
                const newest = this.backing[this.backing.length - 1];
                if (this.ended && newest === undefined) {
                        // ended, no content
                        end_handler(this.backing);
                } else if (this.ended && newest !== undefined) {
                        // ended, some contents
                        push_handler(newest, this.backing, 0);
                        end_handler(this.backing);
                } else if (!this.ended && newest === undefined) {  
                        // not ended, no content
                        this.subscribers.push({ push: push_handler, end: end_handler, index: 0 });
                } else if (!this.ended && newest !== undefined){  
                        // not ended, some contents
                        this.subscribers.push({ push: push_handler, end: end_handler, index: this.backing.length });
                        push_handler(newest, this.backing, 0);
                } else {
                        throw new Error("static assert unreachable");
                }
        }
        [Symbol.asyncIterator](): AsyncIterator<a, void, undefined> {
                return (async function*() {
                        while (true) {
                                
                        }
                })();
        }
}

type StreamStruct = {
        text: string,
        stream: StreamList<StreamStruct> | undefined,
}

type Stray = { 
        // if a message bound for a stray channel attempts to open a subchanel,
        // the attempt is ignored.
        strayid: string,
        attempt_listen: string | undefined,
        attempt_close: boolean,
        data: string,
}

class KawasStreamReader {
        mainout: StreamList<StreamStruct> = new StreamList();
        errout: StreamList<Stray|Error> = new StreamList();
        open_channels: Map<string, StreamList<StreamStruct>> = new Map();
        push(line_incoming: string): void {
                const res = readline(line_incoming);
                if (res instanceof Error) { 
                        this.errout.push(new Error("error parsing command", { cause: res })); 
                        return;
                }
                if ("close" in res) {
                        const close_id = res.close
                        const chan_to_close = this.open_channels.get(close_id);
                        if (!chan_to_close) {
                                this.errout.push({ strayid: res.close, attempt_close: true, attempt_listen: undefined, data: "" });
                                return;
                        }
                        chan_to_close.end();
                        this.open_channels.delete(close_id);
                        return;
                }
                // end of all top level guard blocks 
                if (res.channel === "0" && res.listen === undefined) {
                        // message is bound to main channel, it does not have a subchannel that it claims to observe
                        // the message is a standalone main-bound message 
                        this.mainout.push({ text: res.data, stream: undefined });
                } else if (res.channel === "0" && res.listen !== undefined) {
                        // message is bound to main channel but it opens a new channel to observe
                        if (this.open_channels.has(res.listen)) {
                                // illegal behaviour: attempting to listen on an existing channel
                                // the message itself is still emitted via main channel, but a new 
                                // channel is not opened and an error is returned.
                                this.errout.push(new Error("message attempts to listen to an existing channel: " + res.listen)); 
                                this.mainout.push({ text: res.data, stream: undefined });
                                return;
                        }
                        // end of guard
                        const new_channel: StreamList<StreamStruct> = new StreamList();
                        this.open_channels.set(res.listen, new_channel);
                        this.mainout.push({ text: res.data, stream: new_channel });
                } else if (res.channel !== "0" && res.listen === undefined){
                        // message bound for (hopefully) known channel
                        const target_channel = this.open_channels.get(res.channel);
                        if (target_channel === undefined) {
                                // a stray message
                                this.errout.push({ strayid: res.channel, attempt_close: false, attempt_listen: undefined, data: res.data });
                                return;
                        }
                        // end of guard
                        target_channel.push({ text: res.data, stream: undefined });
                } else if (res.channel !== "0" && res.listen !== undefined) {
                        // message bound for (hopefully) known channel, but this time
                        // it also opens a sub channel
                        const target_channel = this.open_channels.get(res.channel);
                        if (target_channel === undefined) {
                                // a stray message. stray messages cannot open new channels, 
                                // but are still emitted via error channel
                                this.errout.push({ strayid: res.channel, attempt_close: false, attempt_listen: res.listen, data: res.data });
                                return;
                        }
                        if (this.open_channels.has(res.listen)) {
                                // illegal behaviour: attempting to listen on an existing channel
                                // the message itself is still emitted via target channel, but a new 
                                // channel is not opened and an error is returned.
                                this.errout.push(new Error("message attempts to listen to an existing channel: " + res.listen)); 
                                target_channel.push({ text: res.data, stream: undefined });
                                return;
                        }
                        // end of guard blocks
                        const new_channel: StreamList<StreamStruct> = new StreamList();
                        target_channel.push({ text: res.data, stream: new_channel });
                        this.open_channels.set(res.listen, new_channel);
                } else {
                        throw new Error("static guarantee unreachable");
                }
        }

}

async function main() {   
        //@ts-ignore 
        //const { inspect } = await import("util");   
        const ksr = new KawasStreamReader();
        async function sleep(ms: number) { return new Promise(res => setTimeout(res, ms)); }
        ksr.push("0/debug test");
        ksr.push("debug/0 test again 1");
        //@ts-ignore
        console.log("\n\n\n==========", Deno.inspect(ksr, { depth: 100 }));
        await sleep(500);
        ksr.push("debug/0 test again 2");
        //@ts-ignore
        console.log("\n\n\n==========", Deno.inspect(ksr, { depth: 100 }));
        await sleep(500);
        ksr.push("0/1 arsync3 -session XnDf9gHEh");
        ksr.push("1/0 addrev -r 485727 -u 487528 -t 20250304T125959Z -un Hinata -s \"\" -c \".....\"");
        //@ts-ignore
        console.log("\n\n\n==========", Deno.inspect(ksr, { depth: 100 }));
        ksr.push("1/0 addrev -r 485727 -u 487528 -t 20250304T125959Z -un Hinata -s \"\" -c \".....\"");
        await sleep(500);
        //@ts-ignore
        console.log("\n\n\n==========", Deno.inspect(ksr, { depth: 100 }));
        ksr.push("1/0 addrev -r 485727 -u 487528 -t 20250304T125959Z -un Hinata -s \"\" -c \".....\"");
        await sleep(500);
        //@ts-ignore
        console.log("\n\n\n==========", Deno.inspect(ksr, { depth: 100 }));
        ksr.push("1/0 addrev -r 485727 -u 487528 -t 20250304T125959Z -un Hinata -s \"\" -c \".....\"");
        await sleep(500);
        //@ts-ignore
        console.log("\n\n\n==========", Deno.inspect(ksr, { depth: 100 }));
}

main();



// class Arsync3 {
//         static readonly command = "arsync3";
//         //# no args
//         #opt_session: string|undefined;  //# required string
//         #subs: StreamList<Addrev|AddPid|AddPc> = new StreamList();
//         constructor(stream: StreamStruct<string>) {
//                 const parsed = CcCore.parse(stream.struct);
//                 if (parsed instanceof Error) { throw new KawasError(); }
//         }
//         get_revid(): string {}
//         get_timestamp(): string {}
//         get_userid(): string {}
//         get_username(): string[] {}
//         get_summary(): string {}
//         get_content(): string {}
// }

// class Addrev {
//         static readonly command = "addrev";
//         //# no args
//         #opt_revid: string|undefined;  //# required string
//         #opt_timestamp: string|undefined;  //# required string
//         #opt_userid: string|undefined;  //# required string
//         #opt_username: string[] = [];  //# repeated string
//         #opt_summary: string|undefined;  //# required string
//         #opt_content: string|undefined;  //# required string
// }

// class AddPid {
//         static readonly command = "addpid";
//         #args: string|undefined;  //# required string
//         //# no opts
//         //# no subs
//         constructor(parsed: Command) {
//                 if (parsed.args.length < 1) { throw new KawasError("missing required argument") }
//                 this.#args = parsed.args[0];
//         }
//         getarg_0(): string {
//                 if (this.#args === undefined) { throw new NeverError("CU413029561"); }
//                 return this.#args;
//         }
// }

// class AddPc {

// }