import * as StreamList from "./stream-list";
import * as CommandParser from "./command-parser";
import * as Errors from "./errors";
import * as StreamReader from "./stream-reader";

function assert<a>(a: a|undefined): a { 
    if (a === undefined) { throw new Error("a should never be undefined"); }
    return a;
}

class Arsync3 {
        static readonly command = "arsync3";
        //# no args
        #opt_session: string|undefined;  //# required string
        #subs: StreamList.StreamList<Addrev|AddPid|AddPc> = new StreamList.StreamList();
        constructor(stream: StreamReader.StreamStruct) {
                const parsed = CommandParser.CommandCore.parse(stream.text);
                if (parsed instanceof Error) { throw new Errors.KawasError("command malformed", { cause: parsed }); }
                if (parsed.length !== 1) { throw new Errors.KawasError("unexpected multiple commands in one line"); }
                const cmd = parsed[0]!;
                //# no args
                for (let i=0; i<cmd.options.length; i+=2) {
                    if (cmd.options[i] === "session") { this.#opt_session = cmd.options[i+1]; return; }  //# required string
                }
                //# subs: Addrev|AddPid|AddPc
                if (stream.stream === undefined) { 
                    this.#subs.end(); 
                } else {
                    stream.stream.subscribe(
                        (a, l, i) => {
                            for (; i<l.length; i++) {
                                
                            }
                        }, 
                        () => {},
                    );
                }
        }
        get_session(): string { return assert(this.#opt_session); }
}

class Addrev {
        static readonly command = "addrev";
        //# no args
        #opt_revid: string|undefined;  //# required string
        #opt_timestamp: string|undefined;  //# required string
        #opt_userid: string|undefined;  //# required string
        #opt_username: string[] = [];  //# repeated string
        #opt_summary: string|undefined;  //# required string
        #opt_content: string|undefined;  //# required string
}

class AddPid {
        static readonly command = "addpid";
        #args: string|undefined;  //# required string
        //# no opts
        //# no subs
        constructor(parsed: CommandParser.Command) {
                if (parsed.args.length < 1) { throw new Errors.KawasError("missing required argument") }
                this.#args = parsed.args[0];
        }
        getarg_0(): string {
                if (this.#args === undefined) { throw new Errors.NeverError("CU413029561"); }
                return this.#args;
        }
}

class AddPc {

}