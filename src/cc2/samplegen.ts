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
		}
		else if (cur_ch === "\r" && src[i+1] === "\n") {
			coll.push(LineBreak);
			i += 2;
		}
		else if (cur_ch === " " || cur_ch === "\t") { 
			i += 1; 
		}
		else if (cur_ch === "\"") {
			if (i === str_cannot_begin_at) { return [new Error("quoted string term cannot appear back-to-back with a previous term"), i]; }
			const [res, new_i] = consume_quoted(src, "\"", i+1);
			if (res instanceof Error) { return [res, new_i]; }
			coll.push(QuotationIndicator, res);
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

function parse_one(toks: Token2[]): [Command | Error, number] {
	const command: Command = {
		command: "",
		args: [], 
		options: [],
	};
	let i = 0; 
	let positionalMode = false;
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
		for (let i = 0; i<dataline.options.length; i+=2) { coll += " " + dataline.options[i]! + " " + encode_str(dataline.options[i+1] ?? ""); }
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

type __Arinit = {
	namespace: string;
}

export class Arinit{
	static coreParse(pc: Command, checkReq = true): __Arinit|Error {
		let bNamespace = "";
		let cNamespace = 0;
		if (pc.command !== "arinit") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			case "-ns": 
			case "--namespace": 
				bNamespace = pc.options[i+1]!;
				cNamespace += 1;
			break;
			}
		}
		if (cNamespace < 1 && checkReq) {
			return new Error("missing field -ns --namespace");
		}
		return {
			namespace: bNamespace,
		}
	}
	static coreEncode(a: __Arinit): Command {
		const args = [] as string[];
		const options = [] as string[];
		options.push("-ns", a.namespace);
		return {
			command: "arinit",
			args: args,
			options: options,
		}
	}
	static write(a: __Arinit): string {
		return CcCore.encode(this.coreEncode(a));
	}
	static parse(s: string, checkReq = true): __Arinit|Error {
		const e = CcCore.parse(s);
		if (e instanceof Error) { return e; }
		if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
		const f = e[0]!;
		return this.coreParse(f, checkReq);
	}
}

type __Arsync = {
	caseId: string;
	caseNumber: string;
}

export class Arsync{
	static coreParse(pc: Command, checkReq = true): __Arsync|Error {
		let bCaseId = "";
		let cCaseId = 0;
		let bCaseNumber = "";
		let cCaseNumber = 0;
		if (pc.command !== "arsync") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			case "-i": 
			case "--case-id": 
				bCaseId = pc.options[i+1]!;
				cCaseId += 1;
			break;
			case "-n": 
			case "--case-number": 
				bCaseNumber = pc.options[i+1]!;
				cCaseNumber += 1;
			break;
			}
		}
		if (cCaseId < 1 && checkReq) {
			return new Error("missing field -i --case-id");
		}
		if (cCaseNumber < 1 && checkReq) {
			return new Error("missing field -n --case-number");
		}
		return {
			caseId: bCaseId,
			caseNumber: bCaseNumber,
		}
	}
	static coreEncode(a: __Arsync): Command {
		const args = [] as string[];
		const options = [] as string[];
		options.push("-i", a.caseId);
		options.push("-n", a.caseNumber);
		return {
			command: "arsync",
			args: args,
			options: options,
		}
	}
	static write(a: __Arsync): string {
		return CcCore.encode(this.coreEncode(a));
	}
	static parse(s: string, checkReq = true): __Arsync|Error {
		const e = CcCore.parse(s);
		if (e instanceof Error) { return e; }
		if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
		const f = e[0]!;
		return this.coreParse(f, checkReq);
	}
}

type __Addrev = {
	revid: string;
	revTimestamp: string;
	uid: string|undefined;
	uname: boolean;
	summary: string[];
	content: string;
}

export class Addrev{
	static coreParse(pc: Command, checkReq = true): __Addrev|Error {
		let bRevid = "";
		let cRevid = 0;
		let bRevTimestamp = "";
		let cRevTimestamp = 0;
		let bUid = undefined as string|undefined;
		let cUid = 0;
		let bUname = false;
		let cUname = 0;
		let bSummary = [] as string[];
		let cSummary = 0;
		let bContent = "";
		let cContent = 0;
		if (pc.command !== "addrev") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			case "-r": 
			case "--revid": 
				bRevid = pc.options[i+1]!;
				cRevid += 1;
			break;
			case "-t": 
			case "--rev-timestamp": 
				bRevTimestamp = pc.options[i+1]!;
				cRevTimestamp += 1;
			break;
			case "-u": 
			case "--uid": 
				bUid = pc.options[i+1]!;
				cUid += 1;
			break;
			case "-un": 
			case "--uname": 
				bUname = true;
				cUname += 1;
			break;
			case "-s": 
			case "--summary": 
				bSummary.push(pc.options[i+1]!);
				cSummary += 1;
			break;
			case "-c": 
			case "--content": 
				bContent = pc.options[i+1]!;
				cContent += 1;
			break;
			}
		}
		if (cRevid < 1 && checkReq) {
			return new Error("missing field -r --revid");
		}
		if (cRevTimestamp < 1 && checkReq) {
			return new Error("missing field -t --rev-timestamp");
		}
		if (cContent < 1 && checkReq) {
			return new Error("missing field -c --content");
		}
		return {
			revid: bRevid,
			revTimestamp: bRevTimestamp,
			uid: bUid,
			uname: bUname,
			summary: bSummary,
			content: bContent,
		}
	}
	static coreEncode(a: __Addrev): Command {
		const args = [] as string[];
		const options = [] as string[];
		options.push("-r", a.revid);
		options.push("-t", a.revTimestamp);
		if (a.uid !== undefined) { options.push("-u", a.uid); }
		if (a.uname) { options.push("-un", ""); }
		for (const b of a.summary) { options.push("-s", b); }
		options.push("-c", a.content);
		return {
			command: "addrev",
			args: args,
			options: options,
		}
	}
	static write(a: __Addrev): string {
		return CcCore.encode(this.coreEncode(a));
	}
	static parse(s: string, checkReq = true): __Addrev|Error {
		const e = CcCore.parse(s);
		if (e instanceof Error) { return e; }
		if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
		const f = e[0]!;
		return this.coreParse(f, checkReq);
	}
}

type __Addpid = {
	args: string;
}

export class Addpid{
	static coreParse(pc: Command, checkReq = true): __Addpid|Error {
		if (pc.args[0] === undefined) {
			return new Error("missing required arguments");
		}
		let bArgs = pc.args[0]!;
		if (pc.command !== "addpid") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			}
		}
		return {
			args: bArgs, 
		}
	}
	static coreEncode(a: __Addpid): Command {
		const args = [a.args];
		const options = [] as string[];
		return {
			command: "addpid",
			args: args,
			options: options,
		}
	}
	static write(a: __Addpid): string {
		return CcCore.encode(this.coreEncode(a));
	}
	static parse(s: string, checkReq = true): __Addpid|Error {
		const e = CcCore.parse(s);
		if (e instanceof Error) { return e; }
		if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
		const f = e[0]!;
		return this.coreParse(f, checkReq);
	}
}

type __Addpc = {
	args: string[];
}

export class Addpc{
	static coreParse(pc: Command, checkReq = true): __Addpc|Error {
		let bArgs = pc.args;
		if (pc.command !== "addpc") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			}
		}
		return {
			args: bArgs, 
		}
	}
	static coreEncode(a: __Addpc): Command {
		const args = a.args;
		const options = [] as string[];
		return {
			command: "addpc",
			args: args,
			options: options,
		}
	}
	static write(a: __Addpc): string {
		return CcCore.encode(this.coreEncode(a));
	}
	static parse(s: string, checkReq = true): __Addpc|Error {
		const e = CcCore.parse(s);
		if (e instanceof Error) { return e; }
		if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
		const f = e[0]!;
		return this.coreParse(f, checkReq);
	}
}

type __Arclose = {
	caseId: string;
	caseNumber: string;
}

export class Arclose{
	static coreParse(pc: Command, checkReq = true): __Arclose|Error {
		let bCaseId = "";
		let cCaseId = 0;
		let bCaseNumber = "";
		let cCaseNumber = 0;
		if (pc.command !== "arclose") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			case "-i": 
			case "--case-id": 
				bCaseId = pc.options[i+1]!;
				cCaseId += 1;
			break;
			case "-n": 
			case "--case-number": 
				bCaseNumber = pc.options[i+1]!;
				cCaseNumber += 1;
			break;
			}
		}
		if (cCaseId < 1 && checkReq) {
			return new Error("missing field -i --case-id");
		}
		if (cCaseNumber < 1 && checkReq) {
			return new Error("missing field -n --case-number");
		}
		return {
			caseId: bCaseId,
			caseNumber: bCaseNumber,
		}
	}
	static coreEncode(a: __Arclose): Command {
		const args = [] as string[];
		const options = [] as string[];
		options.push("-i", a.caseId);
		options.push("-n", a.caseNumber);
		return {
			command: "arclose",
			args: args,
			options: options,
		}
	}
	static write(a: __Arclose): string {
		return CcCore.encode(this.coreEncode(a));
	}
	static parse(s: string, checkReq = true): __Arclose|Error {
		const e = CcCore.parse(s);
		if (e instanceof Error) { return e; }
		if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
		const f = e[0]!;
		return this.coreParse(f, checkReq);
	}
}




