// support functions

import { Command } from "./backing-codec";

function noundef<T>(a: T|undefined, desc: string): T { if (a === undefined) { throw new Error(desc + " should not be undefined"); } return a; }
function lastof<T>(as: T[]): T|undefined { return as[as.length-1]; }

class DecimalT {
        constructor(private readonly backing = "0") {}
        check(): boolean { return /^\-?[0-9]+$/.test(this.backing); }
        getString(): string { return this.backing; }
        asBigInt(): bigint { return BigInt(this.backing); }
        asNumber(): number { return Number(this.backing); }
}

class UnsignedDecimalT {
        constructor(private readonly backing = "0") {}
        check(): boolean { return /^[0-9]+$/.test(this.backing); }
        getString(): string { return this.backing; }
        asBigInt(): bigint { return BigInt(this.backing); }
        asNumber(): number { return Number(this.backing); }
}

// arsync: 
//    -rev: string repeated
//    #args: string repeated

class Addrev {
        #arg: string[] = [];
        #option_revid: UnsignedDecimalT;
        #option_revtime: string;
        #option_username: string;
        #option_userid: UnsignedDecimalT;
        #option_editsummary: string;
        #option_contnet: string;

        constructor(fromCommand: Command) {
                const unset = new Set(["opt:revid", "opt:revtime", "opt:username", "opt:userid", "opt:editsummary", "opt:content"]);
                if (fromCommand.options.length % 2 !== 0) {
                        throw new Error("parsed command is malformed"); 
                }
                for (const arg of fromCommand.args) {
                        this.#arg.push(arg);
                }
                for (let i=0; i<fromCommand.options.length; i+=2) {
                        const key: string = fromCommand.options[i]!;
                        const rawvalue: string = fromCommand.options[i+1]!;
                        if (key === "r") { 
                                const value = new UnsignedDecimalT(rawvalue);
                                if (!value.check()) { throw new Error("option -r is malformed"); }
                                this.#option_revid = value; 
                                unset.delete("opt:revid");
                        } else if (key === "t") {
                                this.#option_revtime = rawvalue;
                                unset.delete("opt:revtime");
                        } else if (key === "u") {
                                const value = new UnsignedDecimalT(rawvalue);
                                if (!value.check()) { throw new Error("option -u is malformed"); }
                                this.#option_userid = value;
                        } else if (key === "un") {
                                this.#option_username = rawvalue;
                        } else if (key === "s") {
                                this.#option_editsummary = rawvalue;
                        } else if (key === "c") {
                                this.#option_contnet = rawvalue; 
                        }
                }
        }

        intake_option(key: string, rawvalue: string): Error | undefined {

        }
}

/*



parseStream(readableStream, "receiver-plain-array");
parseStream(readableStream, "receiver-two-stack-queue");
parseStream(readableStream, "receiver-asynciter");
parseStream(readableStream, "receiver-custom", { newcoll: () => { .... }, push: coll => { .... }, end: coll => { .... } });


|0| ["arsync", |1|]
|1| { session: |2|, revs: [...3] }
|2| "XnDf9gHEh"
|3| rev (r  )


|0| arsync -session XnDf9gHEh +1;
|1| addrev -r 38810749 -t 2025-11-02T19:58:22Z -u 116625 -un Hinata -s "test again" -c "test page version 2";
|1| addrev -r 38810250 -t 2025-10-25T

type Line

|0| arsync -session XnDf9gHEh ;
|1| arsync 



*/
