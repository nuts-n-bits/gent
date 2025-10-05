// This file contains generated code and should not be modified by hand other than when debugging
// Change the source protocol file (usually with extension .cdef) and re-run codegen to update.
// BEGIN RUNTIME LIBRARY

export type Command = {
	command: string, 
	args: string[],
	options: string[],  // always even number of elements arranged as (k, v, k, v, ...)
}

const LineBreak = 4, QuotationIndicator = 5;
type Token = string | typeof QuotationIndicator | typeof LineBreak;
type Token2 = string | typeof QuotationIndicator;
// All printable characters on ANSI keyboard, less backtick (`), apos ('), quote ("), and backslash (\).
const NONQUOTE_CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDFEGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()-_=+[{]}|;:,<.>/?";

function tokenizer(src: string): [Error|Token[], number] {
	const coll: Token[] = [];
	let i = 0;
	let str_cannot_begin_at = -1;
	while (true) {
		const cur_ch = src[i];
		if (cur_ch === undefined) { 
			return [coll, 0]; 
		} else if (cur_ch === "\n") {
			coll.push(LineBreak);
			i += 1;
		} else if (cur_ch === "\r" || cur_ch === " " || cur_ch === "\t") { 
			i += 1; 
		} else if (cur_ch === "\"") {
			if (i === str_cannot_begin_at) { return [new Error("quoted string term cannot appear back-to-back with a previous term"), i]; }
			const [res, new_i] = consume_quoted(src, "\"", i+1);
			if (res instanceof Error) { return [res, new_i]; }
			coll.push(QuotationIndicator, res);
			str_cannot_begin_at = new_i;
			i = new_i;
		} else if (NONQUOTE_CHARSET.includes(cur_ch)) {
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
		if (cur === undefined) { 
			return [new Error("unexpected eof while consuming quoted"), i]; 
		} else if (cur === "\\") {
			if (src[i+1] === "n") { coll += "\n"; i += 2; }
			else if (src[i+1] === "r") { coll += "\r"; i += 2; }
			else if (src[i+1] === "\\") { coll += "\\"; i += 2; }
			else if (src[i+1] === "t") { coll += "\t"; i += 2; }
			else if (src[i+1] === "\"") { coll += "\""; i += 2; }
			else { return [new Error("unexpected escape sequence"), i+1]; }
		} else if (cur === delim) { 
			return [coll, i+1]; 
		} else { 
			coll += cur; i = i+1; 
		}
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

function parse_one(toks: Token2[]): [Command | Error, number] {
	const command: Command = { command: "", args: [],  options: [] };
	let i = 0, positionalMode = false;
	const first = toks[i], second = toks[i+1];
	if (first === QuotationIndicator) {
		if (typeof second !== "string") { return [new Error("malformed token stream 2745"), i]; }
		command.command = second;
		i += 2;
	} else if (first === undefined) {
		return [new Error("empty token stream 3991"), i];
	} else {
		command.command = first;
		i += 1;
	}
	while (true) {
		const cur = toks[i];
		const next1 = toks[i+1];
		if (cur === undefined) {
			return [command, i+1];
		} else if (cur === QuotationIndicator) {
			if (typeof next1 !== "string") { return [new Error("malformed tokenstream 2524"), i]; }
			command.args.push(next1);
			i += 2;
		} else if (cur[0] !== "-" || positionalMode) {
			command.args.push(cur);
			i += 1;
		} else if (cur === "--") {
			positionalMode = true;
			i += 1;
		} else if (next1 === QuotationIndicator) {
			const next2 = toks[i+2];
			if (typeof next2 !== "string") { return [new Error("malformed token stream 4991"), i]; }
			command.options.push(cur, next2);
			i += 3;
		} else if (next1 === undefined || next1[0] === "-") {
			command.options.push(cur, "");
			i += 1;
		} else {
			command.options.push(cur, next1);
			i += 2;
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
		for (let i = 0; i<dataline.options.length; i+=2) { coll += " " + dataline.options[i]! + " " + encode_str(dataline.options[i+1]!); }
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

// BEGIN MACHINE GENERATED CODE

