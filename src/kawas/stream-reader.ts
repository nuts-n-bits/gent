// Kawas = かわす, lit. message exchanges

import * as StreamList from "./stream-list";
import * as StreamLine from "./stream-line-parser.js";



export type StreamStruct = {
        text: string,
        stream: StreamList.StreamList<StreamStruct> | undefined,
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
        mainout: StreamList.StreamList<StreamStruct> = new StreamList.StreamList();
        errout: StreamList.StreamList<Stray|Error> = new StreamList.StreamList();
        open_channels: Map<string, StreamList.StreamList<StreamStruct>> = new Map();
        push(line_incoming: string): void {
                const res = StreamLine.readline(line_incoming);
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
                        const new_channel: StreamList.StreamList<StreamStruct> = new StreamList.StreamList();
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
                        const new_channel: StreamList.StreamList<StreamStruct> = new StreamList.StreamList();
                        target_channel.push({ text: res.data, stream: new_channel });
                        this.open_channels.set(res.listen, new_channel);
                } else {
                        throw new Error("static guarantee unreachable");
                }
        }

}