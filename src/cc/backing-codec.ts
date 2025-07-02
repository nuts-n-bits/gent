/*

arinit -ns zhwp
arsync -session XnDf9Ghep
addrev -r 398172 -t 2025-03-08T12:29:05Z -u 199272 -un Hinata -s "......" -c "......" 
addrev -r 399478
arsync -session XnDf9Hgpe
rrd -reason RRD#1 -reason RRD#2


*/
function never(_n?: never): never { throw new Error("never"); }

export type Command = {
        command: string, 
        args: string[],
        options: string[],  // always even number of elements arranged as k-v-k-v-...
}

const PairIndicator = 3, LineBreak = 4;
type Token = string | typeof PairIndicator | typeof LineBreak;
type Token2 = string | typeof PairIndicator;
// All printable characters on ANSI keyboard, less backtick (`), apos ('), quote ("), and backslash (\).
const NONQUOTE_CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDFEGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()-_=+[{]}|;:,<.>/?";

export function tokenizer(src: string): [Error|Token[], number] {
        const coll: Token[] = [];
        let i = 0;
        let str_cannot_begin_at = -1;
        while (true) {
                const cur_ch = src[i];
                if (cur_ch === undefined) { return [coll, 0]; }
                else if (cur_ch === "\n") {
                        coll.push(LineBreak);
                        i += 1;
                }
                else if (cur_ch === "\r" && src[i+1] === "\n") {
                        coll.push(LineBreak);
                        i += 2;
                }
                else if (" \t".includes(cur_ch)) { i += 1; }
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
                        } else if (next_ch === undefined || " \n\r\t".includes(next_ch)) {
                                coll.push(PairIndicator, "");
                                i += 1;
                        } else {
                                return [new Error("unexpected character after pair indicator `-`"), i];
                        }
                }
                else if (cur_ch === "\"") {
                        if (i === str_cannot_begin_at) { return [new Error("quoted string term cannot appear back-to-back with a previous term"), i]; }
                        const [res, new_i] = consume_quoted(src, "\"", i+1);
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
                } else {
                        return [new Error("unexpected character"), i];
                }
        }
}

function consume_quoted(src: string, delim: string, i: number): [string|Error, number] {
        let coll = "";
        while(true) {
                const cur = src[i];
                if (cur === undefined) { return [new Error("unexpected eof while consuming quoted"), i]; }
                else if (cur === "\\") {
                        if (src[i+1] === "n") { coll += "\n"; i += 2; }
                        else if (src[i+1] === "r") { coll += "\r"; i += 2; }
                        else if (src[i+1] === "\\") { coll += "\\"; i += 2; }
                        else if (src[i+1] === "t") { coll += "\t"; i += 2; }
                        else if (src[i+1] === "\"") { coll += "\""; i += 2; }
                        else if (src[i+1] === "\'") { coll += "\'"; i += 2; }
                        else if (src[i+1] === "\`") { coll += "\`"; i += 2; }
                        else { return [new Error("unexpected escape sequence"), i+1]; }
                }
                else if (cur === delim) { return [coll, i+1]; }
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

export function parse_one(toks: Token2[]): [Command | Error, number] {
        const command: Command = {
                command: "",
                args: [], 
                options: [],
        };
        let i = 0; 
        const first = toks[i], second = toks[i+1];
        if (first === PairIndicator) {
                if (typeof second !== "string") { return [new Error("malformed token stream 2741"), i]; }
                command.command = "-" + second;
                i += 2;
        } else {
                command.command = first!;
                i += 1;
        }
        while (true) {
                const cur = toks[i];
                const next = toks[i+1];
                if (cur === PairIndicator) {
                        const nextnext = toks[i+2];
                        if (typeof next !== "string") {
                                return [new Error("malformed tokenstream 2519"), i];
                        }
                        if (typeof nextnext === "string") {
                                command.options.push(next, nextnext);
                                i += 3;
                        } else {
                                command.options.push(next, "");
                                i += 2;
                        }
                } else if (typeof cur === "string") {
                        command.args.push(cur);
                        i += 1;
                } else if (cur === undefined) {
                        return [command, i+1];
                } else {
                        never(cur);
                }
        }
}

function group_tok(toks: Token[]): Token2[][] {
        const coll: Token2[][] = [[]];
        for (const tok of toks) {
                if (tok === LineBreak) { coll.push([]); }
                else { coll[coll.length-1]!.push(tok); }
        }
        return coll.filter(a => a.length > 0);
}

export class CcCore {
        static parse(src: string): Command[] | Error {
                const [tokens, i] = tokenizer(src);
                if (tokens instanceof Error) { return tokens; }
                const grouped_toks = group_tok(tokens);
                if (grouped_toks instanceof Error) { return grouped_toks; }
                const coll = [] as Command[];
                for (const grouped of grouped_toks) {
                        const [parsed, i] = parse_one(grouped);
                        if (parsed instanceof Error) { return parsed; }
                        coll.push(parsed);
                }
                return coll;
        }
        static encode(dataline: Command): string {
                let coll = encode_str(dataline.command);
                for (const arg of dataline.args) { coll += " " + encode_str(arg); }
                for (let i = 0; i<dataline.options.length; ) { coll += " -" + encode_str(dataline.options[i]!) + " " + encode_str(dataline.options[i+1] ?? ""); }
                return coll;
        }
}

function nqtest(tested: string): boolean {
        for (const ch of tested) { if (!NONQUOTE_CHARSET.includes(ch)) { return false; } }
        return true;
}

function encode_str(s: string): string {
        const noneed_quote = s.length > 0 && s.length < 50 && s[0] !== "-" && nqtest(s);
        if (noneed_quote) { return s; }
        s.replaceAll("\\", "\\\\").replaceAll("\"", "\\\"").replaceAll("\r", "\\r").replaceAll("\n", "\\n");
        return "\"" + s + "\"";
}