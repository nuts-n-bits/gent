/*
        Spec:


        |0| arinit -ns zhwp;      // --> parsed as { chans: ["0"],       leading: "arinit" , args: [ ] ,     rec: { ns: "zhwp"               } , sub: undefined }
        |0| arsync -s 1284 +1;    // --> parsed as { chans: ["0", "+1"], leading: "arsync" , args: [ ] ,     rec: { session: "XnDf9Ghep"     } , sub: "1"       }
        |1| addrev -r 3981 +2;    // --> parsed as { chans: ["1", "+2"], leading: "addrev" , args: [ ] ,     rec: { r: "398172", u: "199272" } , sub: "2"       }
        |1| addrev -r 3992 +2;    // --> parsed as { chans: ["1", "+2"], leading: "addrev" , args: [ ] ,     rec: { r: "399478"              } , sub: "2"       }  <-- this represents an error! resubscription to channel "2".
        |0| arsync -s 4991 +3;    // --> parsed as { chans: ["0", "+3"], leading: "arsync" , args: [ ] ,     rec: { session: "XnDf9Hgpe"     } , sub: undefined }
        |2| rrd -rs RD1 -rs RD2;  // --> parsed as { chans: ["2"],       leading: "rrd"    , args: [ ] ,     rec: { reason: "RRD#2"          } , sub: undefined }  <-- this represents an error! the record key "reason" appears twice.

        This message: |0| ..... +1;
        Is parsed into: { chans: ["0"], listen: "1", ..... }
        It means: This message goes to channel 0 (i.e. it targets channel 0), it subscribes to channel 1.
        
        This message: /2;
        Is parsed into: { chanends: "2" } 
        It means: This message closes channel 2. subsequent messages to channel 2 will not be sent to listeners. 
                  They will be sent to the error channel instead, like any messages to an undeclared channel (except channel "0").
        
        This message: +1;
        Cannot be parsed. Any messages must start with channel declarations like |1 2 3| to indicate target channels,
        or be of form /n; which signals the closure of channel n.


        
*/
function never(_n?: never): never { throw new Error("never"); }

export type StreamLine = {
        chans: string[],
        listen: string|undefined;
        leading: string,
        args: string[],
        options: Record<string, string>,
} | {
        close: string,
}

const Semicolon = 2, ChanDeclIndicator = 5, ChanSubscribeIndicator = 6, ChanCloseIndicator = 7, PairIndicator = 8;
type Token = string | typeof Semicolon | typeof ChanDeclIndicator | typeof ChanSubscribeIndicator | typeof ChanCloseIndicator | typeof PairIndicator;
type Token2 = Exclude<Token, typeof Semicolon>;
const ALPHANUMERIC_CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDFEGHIJKLMNOPQRSTUVWXYZ";
// All printable characters on ANSI keyboard, less backtick (`), apos ('), quote ("), brackets ([]), semicolon (;), and backslash (\).
const NONQUOTE_CHARSET = ALPHANUMERIC_CHARSET + "~!@#$%^&*()-_=+{}|:,<.>/?";
const CHANID_CHARSET = ALPHANUMERIC_CHARSET + "_";
const WHITESPACE_CHARSET = " \n\r\t";

export function tokenizer(src: string): [Error|Token[], number] {
        const coll: Token[] = [];
        let i = 0;
        let str_cannot_begin_at = -1;
        while (true) {
                const cur_ch = src[i];
                if (cur_ch === undefined) { return [coll, 0]; }
                else if (cur_ch === ";") { coll.push(Semicolon); i += 1; }
                else if (cur_ch === "[") { 
                        const [res, new_i] = consume_channel_decl_block(src, i+1);
                        if (res instanceof Error) { return [res, new_i]; }
                        if (res.length === 0) { return [new Error("empty channel control block"), new_i]; }
                        coll.push(ChanDeclIndicator, res);
                        i = new_i;
                }
                else if (WHITESPACE_CHARSET.includes(cur_ch)) { i += 1; /* skip whitespace */ }
                else if (cur_ch === "-") {
                        const next_ch = src[i+1];
                        if (next_ch === "\"") {
                                coll.push(PairIndicator);
                                i += 1;
                        } else if (next_ch && NONQUOTE_CHARSET.includes(next_ch)) {
                                const [res, new_i] = consume_nonquoted(src, i);
                                coll.push(PairIndicator, res.substring(1));
                                str_cannot_begin_at = new_i;
                                i = new_i;
                        } else if (next_ch !== undefined && " \n\r\t;".includes(next_ch)) {
                                coll.push(PairIndicator, "");
                                i += 1;
                        } else {
                                return [new Error("expected quoted or nonquoted string after pair indicator `-`"), i];
                        }
                }
                else if (cur_ch === "/" || cur_ch === "+") {
                        const [res, new_i] = consume_nonquoted(src, i+1);
                        if (badchanid(res)) { return [new Error("malformed channel id"), i+1]; }
                        coll.push(cur_ch === "/" ? ChanCloseIndicator : ChanSubscribeIndicator, res);
                        str_cannot_begin_at = new_i;
                        i = new_i;
                }
                else if (cur_ch === "\"") {
                        if (i === str_cannot_begin_at) { return [new Error("quoted string term cannot appear back-to-back with a previous term"), i]; }
                        const [res, new_i] = consume_quoted(src, i+1);
                        if (res instanceof Error) { return [res, new_i]; }
                        else { coll.push(res); }
                        str_cannot_begin_at = new_i;
                        i = new_i;
                }
                else if (NONQUOTE_CHARSET.includes(cur_ch)) {
                        if (i === str_cannot_begin_at) { return [new Error("non-quoted string term cannot appear back-to-back with a previous term"), i]; }
                        const [res, new_i] = consume_nonquoted(src, i);
                        coll.push(res);
                        str_cannot_begin_at = new_i;
                        i = new_i;
                }
        }
}

function consume_quoted(src: string, i: number): [string|Error, number] {
        let coll = "";
        while(true) {
                const cur = src[i];
                if (cur === undefined) { return [new Error("unexpected eof while consuming quoted"), i]; }
                else if (cur === "\\") {
                        if (src[i+1] === "n") { coll += "\n"; i += 2; }
                        else if (src[i+1] === "r") { coll += "\r"; i += 2; }
                        else if (src[i+1] === "\\") { coll += "\\"; i += 2; }
                        else if (src[i+1] === "\t") { coll += "\t"; i += 2; }
                        else if (src[i+1] === "\"") { coll += "\""; i += 2; }
                        else { return [new Error("unexpected escape sequence"), i+1]; }
                }
                else if (cur === "\"") { return [coll, i+1]; }
                else { coll += cur; i = i+1; }
        }
}

function consume_nonquoted(src: string, i: number): [string, number] {
        let coll = "";
        while (true) {
                const cur = src[i];
                if (cur !== undefined && NONQUOTE_CHARSET.includes(cur)) { coll += cur; i = i+1; }
                else { return [coll, i]; }
        }
}

function consume_channel_decl_block(src: string, i: number): [string|Error, number] {
        let coll = "";
        while (true) {
                const cur = src[i];
                if (cur === undefined) { return [new Error("unexpected eof while consuming channel declaration block"), i]; }
                else if (cur === "]") { return [coll.trim(), i+1]; }
                else if (CHANID_CHARSET.includes(cur)) { coll += cur; i += 1; }
                else if (WHITESPACE_CHARSET.includes(cur)) { if (coll[coll.length-1] !== " ") { coll += " "; } i += 1; }
                else { return [new Error("unexpected character while parsing channel declaration block"), i]; }
        }
}

export function parse_one(toks: Token2[], i = 0): [StreamLine | Error, number] {
        
        // consume a line with shape "/n;"
        if (toks[i] === ChanCloseIndicator) {
                const chanid = toks[1];
                if (typeof chanid !== "string") { return [new Error("expecting string term following channel close indicator"), i+1]; }
                if (toks[i+2] !== undefined) { return [new Error("expecting eof after channel closure declaration"), i+2]; }
                if (badchanid(chanid)) { return [new Error("malformed chanid"), i+1]; }
                return [{ close: chanid }, i];
        }

        // if the line has shape "|n| .... ;"
        // consume channel declatation block
        const cur = toks[i], next = toks[i+1];
        if (cur !== ChanDeclIndicator) { return [new Error("a dataline must begin with a well-formed channel declaration block"), i]; }
        if (typeof next !== "string") { return [new Error("lexer internal error - channel declaration indicator is not followed by a string term"), i]; }
        const ret_chans = next.split(" ");
        if (ret_chans.length === 0) { return [new Error("empty channel declaration block"), i]; }
        i += 2;

        const parsed: StreamLine = {
                args: [],
                chans: ret_chans,
                leading: "",
                listen: undefined,
                options: {},
        }

        // consume leading term
        const first = toks[i], second = toks[i+1];
        if (first === PairIndicator) {
                if (typeof second !== "string") { return [new Error("malformed token stream 2741"), i]; }
                parsed.leading = "-" + second;
                i += 2;
        } else if (typeof first === "string") {
                parsed.leading = first!;
                i += 1;
        } else {
                return [new Error("expected leading string term, found channel controls"), i];
        }

        while (true) {
                const cur = toks[i];
                const next = toks[i+1];
                if (cur === ChanDeclIndicator) {
                        return [new Error("channel declaration is only expected at the beginning of a dataline"), i];
                } else if (cur === ChanCloseIndicator) {
                        return [new Error("channel closure declaration is only expected as an independent dataline"), i];
                } else if (cur === PairIndicator) {
                        const nextnext = toks[i+2];
                        if (typeof next !== "string") { return [new Error("lexer internal error: at least one string term must come after a pair indicator"), i+1]; }
                        if (parsed.options[next] !== undefined) { return [new Error("repeated record entry"), i+1]; }
                        if (typeof nextnext === "string") { parsed.options[next] = nextnext; i += 3; }
                        else { parsed.options[next] = ""; i += 2; }
                } else if (typeof cur === "string") {
                        parsed.args.push(cur);
                        i += 1;
                } else if (cur === ChanSubscribeIndicator) {
                        if (typeof next !== "string") { return [new Error("lexer internal error: ChanSubscribeIndicator must be followed by string term"), i+1]; }
                        if (parsed.listen !== undefined) { return [new Error("cannot listen to multiple channels"), i]; }
                        parsed.listen = next;
                        i += 2;
                } else if (cur === undefined) {
                        return [parsed, i+1];
                } else {
                        console.log(cur);
                        never(cur);
                }
        }
}

function group_tok(toks: Token[]): Token2[][] | Error {
        const coll: Token2[][] = [[]];
        for (const tok of toks) {
                if (tok === Semicolon) { coll.push([]); }
                else { coll[coll.length-1]!.push(tok); }
        }
        if (coll[coll.length-1]!.length > 0) { return new Error("missing semicolon at the end"); }
        return coll.filter(a => a.length > 0);
}

export class StreamlineCore {
        static parse_all(src: string): StreamLine[] | Error {
                const [tokens, i] = tokenizer(src);
                if (tokens instanceof Error) { return tokens; }
                const grouped_toks = group_tok(tokens);
                if (grouped_toks instanceof Error) { return grouped_toks; }
                const coll = [] as StreamLine[];
                for (const grouped of grouped_toks) {
                        const [parsed, i] = parse_one(grouped);
                        if (parsed instanceof Error) { return parsed; }
                        coll.push(parsed);
                }
                return coll;
        }
        static parse_one(src: string): StreamLine | Error {
                const [tokens, i] = tokenizer(src);
                if (tokens instanceof Error) { return tokens; }
                const grouped_toks = group_tok(tokens);
                if (grouped_toks instanceof Error) { return grouped_toks; }
                if (grouped_toks.length !== 1) { return new Error("expected exactly one dataline, got " + grouped_toks.length); }
                const [parsed, i2] = parse_one(grouped_toks[0]!);
                return parsed;
        }
        static encode(parsed: StreamLine): string|Error {
                if ("close" in parsed) {
                        if (badchanid(parsed.close)) { return new Error("malformed channel id"); }
                        return "/" + parsed.close + ";";
                }
                let coll = "";
                if (parsed.chans.length < 0) { return new Error("no channel declaration"); }
                for (const target of parsed.chans) { if (badchanid(target)) { return new Error("no channel declaration"); } }
                coll += "[" + parsed.chans.join(" ") + "]";
                coll += " " + encode_str(parsed.leading);
                for (const arg of parsed.args) { coll += " " + encode_str(arg); }
                if (parsed.listen !== undefined) {
                        if (badchanid(parsed.listen)) { return new Error("malformed channel id"); }
                        coll += " +" + parsed.listen;
                }
                for (const [k, v] of Object.entries(parsed.options)) { coll += " -" + encode_str(k) + " " + encode_str(v); }
                return coll + ";";
        }
}

function charsettest(charset: string, tested: string): boolean {
        for (const ch of tested) { if (!NONQUOTE_CHARSET.includes(ch)) { return false; } }
        return true;
}

function encode_str(s: string): string {
        const need_quote = s.length === 0 || s.length > 50 || s[0] === "-" || s[0] === "/" || s[0] === "+" || !charsettest(NONQUOTE_CHARSET, s);
        if (!need_quote) { return s; }
        s.replaceAll("\\", "\\\\").replaceAll("\"", "\\\"").replaceAll("\r", "\\r").replaceAll("\n", "\\n");
        return "\"" + s + "\"";
}

function badchanid(chanid: string): boolean {
        if (chanid.length === 0) { return true; }
        if (!charsettest(CHANID_CHARSET, chanid)) { return true; }
        return false;
}