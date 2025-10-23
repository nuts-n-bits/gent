// This file contains generated code and should not be modified by hand other than when debugging
// Change the source protocol file (usually with extension .cdef) and re-run codegen to update.
// BEGIN RUNTIME LIBRARY

export type Command = {
	command: string, 
	args: string[],
	options: string[],  // always even number of elements arranged as (k, v, k, v, ...)
}

const _LineBreak = 4, _QuotationIndicator = 5;
type _Token = string | typeof _QuotationIndicator | typeof _LineBreak;
type _Token2 = string | typeof _QuotationIndicator;
// All printable characters on ANSI keyboard, less backtick (`), apos ('), quote ("), and backslash (\).
const _NONQUOTE_CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDFEGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()-_=+[{]}|;:,<.>/?";

function _tokenizer(src: string): [Error|_Token[], number] {
	const coll: _Token[] = [];
	let i = 0;
	let str_cannot_begin_at = -1;
	while (true) {
		const cur_ch = src[i];
		if (cur_ch === undefined) { 
			return [coll, 0]; 
		} else if (cur_ch === "\n") {
			coll.push(_LineBreak);
			i += 1;
		} else if (cur_ch === "\r" || cur_ch === " " || cur_ch === "\t") { 
			i += 1; 
		} else if (cur_ch === "\"") {
			if (i === str_cannot_begin_at) { return [new Error("quoted string term cannot appear back-to-back with a previous term"), i]; }
			const [res, new_i] = _consume_quoted(src, "\"", i+1);
			if (res instanceof Error) { return [res, new_i]; }
			coll.push(_QuotationIndicator, res);
			str_cannot_begin_at = new_i;
			i = new_i;
		} else if (_NONQUOTE_CHARSET.includes(cur_ch)) {
			if (i === str_cannot_begin_at) { return [new Error("non-quoted string term cannot appear back-to-back with a previous term"), i]; }
			const [res, new_i] = _consume_nonquoted(src, i);
			coll.push(res);
			str_cannot_begin_at = new_i;
			i = new_i;
		} else {
			return [new Error("unexpected character"), i];
		}
	}
}

function _consume_quoted(src: string, delim: string, i: number): [string|Error, number] {
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

function _consume_nonquoted(src: string, i: number): [string, number] {
	let coll = "";
	while (true) {
		const cur = src[i];
		if (cur !== undefined && _NONQUOTE_CHARSET.includes(cur)) { coll += cur; i = i+1; }
		else { return [coll, i]; }
	}
}

function _parse_one(toks: _Token2[]): [Command | Error, number] {
	const command: Command = { command: "", args: [],  options: [] };
	let i = 0, positionalMode = false;
	const first = toks[i], second = toks[i+1];
	if (first === _QuotationIndicator) {
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
		} else if (cur === _QuotationIndicator) {
			if (typeof next1 !== "string") { return [new Error("malformed tokenstream 2524"), i]; }
			command.args.push(next1);
			i += 2;
		} else if (cur[0] !== "-" || positionalMode) {
			command.args.push(cur);
			i += 1;
		} else if (cur === "--") {
			positionalMode = true;
			i += 1;
		} else if (next1 === _QuotationIndicator) {
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

function _group_tok(toks: _Token[]): _Token2[][] {
	const coll: _Token2[][] = [[]];
	for (const tok of toks) {
		if (tok === _LineBreak) { coll.push([]); }
		else { coll[coll.length-1]!.push(tok); }
	}
	return coll.filter(a => a.length > 0);
}

export class CcCore {
	static parse(src: string): Command[] | Error {
		const [tokens, i] = _tokenizer(src);
		if (tokens instanceof Error) { return tokens; }
		const grouped_toks = _group_tok(tokens);
		const coll = [] as Command[];
		for (const grouped of grouped_toks) {
			const [parsed, i] = _parse_one(grouped);
			if (parsed instanceof Error) { return parsed; }
			coll.push(parsed);
		}
		return coll;
	}
	static encode(dataline: Command): string {
		let coll = _encode_str(dataline.command);
		for (const arg of dataline.args) { coll += " " + _encode_str(arg); }
		if (dataline.options.length % 2 !== 0) { dataline.options.push(""); }
		for (let i = 0; i<dataline.options.length; i+=2) { 
			coll += " " + dataline.options[i]!
			if (dataline.options[i+1] !== "") {
				coll += " " + _encode_str(dataline.options[i+1]!);
			}
		}
		return coll;
	}
}

function _nqtest(tested: string): boolean {
	for (const ch of tested) { if (!_NONQUOTE_CHARSET.includes(ch)) { return false; } }
	return true;
}

function _encode_str(s: string): string {
	const noneed_quote = s.length > 0 && s.length < 50 && s[0] !== "-" && _nqtest(s);
	if (noneed_quote) { return s; }
	s.replaceAll("\\", "\\\\").replaceAll("\"", "\\\"").replaceAll("\r", "\\r").replaceAll("\n", "\\n");
	return "\"" + s + "\"";
}

// BEGIN MACHINE GENERATED CODE

type __Token = {
    args: string;
}

export class Token{
    static coreParse(pc: Command): { res: __Token, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "token") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __Token): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "token",
            args: args,
            options: options,
        }
    }
    static write(a: __Token): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __Token, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __Token, missing: string } {
        const ev = this.coreParse({ command: "token", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __Proxy = {
    args: string;
    enable: boolean;
}

export class Proxy{
    static coreParse(pc: Command): { res: __Proxy, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        let bEnable = false;
        let cEnable = 0;
        if (pc.command !== "proxy") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            case "-enable": 
                bEnable = true;
                cEnable += 1;
            break;
            }
        }
        return {
            res: {
                args: bArgs, 
                enable: bEnable,
            },
            missing: missing,
        }
    }
    static coreEncode(a: __Proxy): Command {
        const args = [a.args];
        const options = [] as string[];
        if (a.enable) { options.push("-enable", ""); }
        return {
            command: "proxy",
            args: args,
            options: options,
        }
    }
    static write(a: __Proxy): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __Proxy, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __Proxy, missing: string } {
        const ev = this.coreParse({ command: "proxy", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __Record = {
    args: string;
}

export class Record{
    static coreParse(pc: Command): { res: __Record, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "record") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __Record): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "record",
            args: args,
            options: options,
        }
    }
    static write(a: __Record): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __Record, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __Record, missing: string } {
        const ev = this.coreParse({ command: "record", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __Groups = {
    args: string[];
}

export class Groups{
    static coreParse(pc: Command): { res: __Groups, missing: string }|Error {
        let missing = "";
        let bArgs = pc.args;
        if (pc.command !== "groups") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __Groups): Command {
        const args = a.args;
        const options = [] as string[];
        return {
            command: "groups",
            args: args,
            options: options,
        }
    }
    static write(a: __Groups): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __Groups, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __Groups, missing: string } {
        const ev = this.coreParse({ command: "groups", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __LogChannel = {
    args: string;
}

export class LogChannel{
    static coreParse(pc: Command): { res: __LogChannel, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "log_channel") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __LogChannel): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "log_channel",
            args: args,
            options: options,
        }
    }
    static write(a: __LogChannel): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __LogChannel, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __LogChannel, missing: string } {
        const ev = this.coreParse({ command: "log_channel", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MainSite = {
    args: string;
}

export class MainSite{
    static coreParse(pc: Command): { res: __MainSite, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "main_site") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MainSite): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "main_site",
            args: args,
            options: options,
        }
    }
    static write(a: __MainSite): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MainSite, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MainSite, missing: string } {
        const ev = this.coreParse({ command: "main_site", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __OauthAuthUrl = {
    args: string;
}

export class OauthAuthUrl{
    static coreParse(pc: Command): { res: __OauthAuthUrl, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "oauth_auth_url") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __OauthAuthUrl): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "oauth_auth_url",
            args: args,
            options: options,
        }
    }
    static write(a: __OauthAuthUrl): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __OauthAuthUrl, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __OauthAuthUrl, missing: string } {
        const ev = this.coreParse({ command: "oauth_auth_url", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __OauthQueryUrl = {
    args: string;
}

export class OauthQueryUrl{
    static coreParse(pc: Command): { res: __OauthQueryUrl, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "oauth_query_url") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __OauthQueryUrl): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "oauth_query_url",
            args: args,
            options: options,
        }
    }
    static write(a: __OauthQueryUrl): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __OauthQueryUrl, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __OauthQueryUrl, missing: string } {
        const ev = this.coreParse({ command: "oauth_query_url", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __OauthQueryKey = {
    args: string;
}

export class OauthQueryKey{
    static coreParse(pc: Command): { res: __OauthQueryKey, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "oauth_query_key") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __OauthQueryKey): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "oauth_query_key",
            args: args,
            options: options,
        }
    }
    static write(a: __OauthQueryKey): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __OauthQueryKey, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __OauthQueryKey, missing: string } {
        const ev = this.coreParse({ command: "oauth_query_key", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __WikiList = {
    args: string[];
}

export class WikiList{
    static coreParse(pc: Command): { res: __WikiList, missing: string }|Error {
        let missing = "";
        let bArgs = pc.args;
        if (pc.command !== "wiki_list") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __WikiList): Command {
        const args = a.args;
        const options = [] as string[];
        return {
            command: "wiki_list",
            args: args,
            options: options,
        }
    }
    static write(a: __WikiList): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __WikiList, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __WikiList, missing: string } {
        const ev = this.coreParse({ command: "wiki_list", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __Blacklist = {
    args: string[];
}

export class Blacklist{
    static coreParse(pc: Command): { res: __Blacklist, missing: string }|Error {
        let missing = "";
        let bArgs = pc.args;
        if (pc.command !== "blacklist") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __Blacklist): Command {
        const args = a.args;
        const options = [] as string[];
        return {
            command: "blacklist",
            args: args,
            options: options,
        }
    }
    static write(a: __Blacklist): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __Blacklist, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __Blacklist, missing: string } {
        const ev = this.coreParse({ command: "blacklist", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageStart = {
    args: string;
}

export class MessageStart{
    static coreParse(pc: Command): { res: __MessageStart, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-start") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageStart): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-start",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageStart): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageStart, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageStart, missing: string } {
        const ev = this.coreParse({ command: "message-start", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessagePolicy = {
    args: string;
}

export class MessagePolicy{
    static coreParse(pc: Command): { res: __MessagePolicy, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-policy") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessagePolicy): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-policy",
            args: args,
            options: options,
        }
    }
    static write(a: __MessagePolicy): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessagePolicy, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessagePolicy, missing: string } {
        const ev = this.coreParse({ command: "message-policy", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageInsufficientRight = {
    args: string;
}

export class MessageInsufficientRight{
    static coreParse(pc: Command): { res: __MessageInsufficientRight, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-insufficient_right") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageInsufficientRight): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-insufficient_right",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageInsufficientRight): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageInsufficientRight, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageInsufficientRight, missing: string } {
        const ev = this.coreParse({ command: "message-insufficient_right", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageGeneralPrompt = {
    args: string;
}

export class MessageGeneralPrompt{
    static coreParse(pc: Command): { res: __MessageGeneralPrompt, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-general_prompt") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageGeneralPrompt): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-general_prompt",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageGeneralPrompt): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageGeneralPrompt, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageGeneralPrompt, missing: string } {
        const ev = this.coreParse({ command: "message-general_prompt", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageTelegramIdError = {
    args: string;
}

export class MessageTelegramIdError{
    static coreParse(pc: Command): { res: __MessageTelegramIdError, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-telegram_id_error") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageTelegramIdError): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-telegram_id_error",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageTelegramIdError): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageTelegramIdError, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageTelegramIdError, missing: string } {
        const ev = this.coreParse({ command: "message-telegram_id_error", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageRestoreSilence = {
    args: string;
}

export class MessageRestoreSilence{
    static coreParse(pc: Command): { res: __MessageRestoreSilence, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-restore_silence") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageRestoreSilence): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-restore_silence",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageRestoreSilence): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageRestoreSilence, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageRestoreSilence, missing: string } {
        const ev = this.coreParse({ command: "message-restore_silence", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmAlready = {
    args: string;
}

export class MessageConfirmAlready{
    static coreParse(pc: Command): { res: __MessageConfirmAlready, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_already") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmAlready): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_already",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmAlready): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmAlready, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmAlready, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_already", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmOtherTg = {
    args: string;
}

export class MessageConfirmOtherTg{
    static coreParse(pc: Command): { res: __MessageConfirmOtherTg, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_other_tg") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmOtherTg): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_other_tg",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmOtherTg): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmOtherTg, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmOtherTg, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_other_tg", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmConflict = {
    args: string;
}

export class MessageConfirmConflict{
    static coreParse(pc: Command): { res: __MessageConfirmConflict, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_conflict") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmConflict): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_conflict",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmConflict): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmConflict, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmConflict, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_conflict", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmChecking = {
    args: string;
}

export class MessageConfirmChecking{
    static coreParse(pc: Command): { res: __MessageConfirmChecking, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_checking") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmChecking): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_checking",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmChecking): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmChecking, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmChecking, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_checking", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmUserNotFound = {
    args: string;
}

export class MessageConfirmUserNotFound{
    static coreParse(pc: Command): { res: __MessageConfirmUserNotFound, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_user_not_found") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmUserNotFound): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_user_not_found",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmUserNotFound): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmUserNotFound, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmUserNotFound, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_user_not_found", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmButton = {
    args: string;
}

export class MessageConfirmButton{
    static coreParse(pc: Command): { res: __MessageConfirmButton, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_button") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmButton): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_button",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmButton): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmButton, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmButton, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_button", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmWait = {
    args: string;
}

export class MessageConfirmWait{
    static coreParse(pc: Command): { res: __MessageConfirmWait, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_wait") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmWait): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_wait",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmWait): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmWait, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmWait, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_wait", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmConfirming = {
    args: string;
}

export class MessageConfirmConfirming{
    static coreParse(pc: Command): { res: __MessageConfirmConfirming, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_confirming") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmConfirming): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_confirming",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmConfirming): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmConfirming, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmConfirming, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_confirming", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmIneligible = {
    args: string;
}

export class MessageConfirmIneligible{
    static coreParse(pc: Command): { res: __MessageConfirmIneligible, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_ineligible") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmIneligible): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_ineligible",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmIneligible): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmIneligible, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmIneligible, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_ineligible", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmSessionLost = {
    args: string;
}

export class MessageConfirmSessionLost{
    static coreParse(pc: Command): { res: __MessageConfirmSessionLost, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_session_lost") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmSessionLost): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_session_lost",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmSessionLost): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmSessionLost, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmSessionLost, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_session_lost", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmComplete = {
    args: string;
}

export class MessageConfirmComplete{
    static coreParse(pc: Command): { res: __MessageConfirmComplete, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_complete") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmComplete): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_complete",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmComplete): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmComplete, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmComplete, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_complete", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmFailed = {
    args: string;
}

export class MessageConfirmFailed{
    static coreParse(pc: Command): { res: __MessageConfirmFailed, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_failed") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmFailed): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_failed",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmFailed): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmFailed, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmFailed, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_failed", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageConfirmLog = {
    args: string;
}

export class MessageConfirmLog{
    static coreParse(pc: Command): { res: __MessageConfirmLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-confirm_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageConfirmLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-confirm_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageConfirmLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageConfirmLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageConfirmLog, missing: string } {
        const ev = this.coreParse({ command: "message-confirm_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDeconfirmPrompt = {
    args: string;
}

export class MessageDeconfirmPrompt{
    static coreParse(pc: Command): { res: __MessageDeconfirmPrompt, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-deconfirm_prompt") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDeconfirmPrompt): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-deconfirm_prompt",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDeconfirmPrompt): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDeconfirmPrompt, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDeconfirmPrompt, missing: string } {
        const ev = this.coreParse({ command: "message-deconfirm_prompt", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDeconfirmButton = {
    args: string;
}

export class MessageDeconfirmButton{
    static coreParse(pc: Command): { res: __MessageDeconfirmButton, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-deconfirm_button") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDeconfirmButton): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-deconfirm_button",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDeconfirmButton): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDeconfirmButton, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDeconfirmButton, missing: string } {
        const ev = this.coreParse({ command: "message-deconfirm_button", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDeconfirmSucc = {
    args: string;
}

export class MessageDeconfirmSucc{
    static coreParse(pc: Command): { res: __MessageDeconfirmSucc, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-deconfirm_succ") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDeconfirmSucc): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-deconfirm_succ",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDeconfirmSucc): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDeconfirmSucc, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDeconfirmSucc, missing: string } {
        const ev = this.coreParse({ command: "message-deconfirm_succ", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDeconfirmNotConfirmed = {
    args: string;
}

export class MessageDeconfirmNotConfirmed{
    static coreParse(pc: Command): { res: __MessageDeconfirmNotConfirmed, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-deconfirm_not_confirmed") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDeconfirmNotConfirmed): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-deconfirm_not_confirmed",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDeconfirmNotConfirmed): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDeconfirmNotConfirmed, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDeconfirmNotConfirmed, missing: string } {
        const ev = this.coreParse({ command: "message-deconfirm_not_confirmed", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDeconfirmLog = {
    args: string;
}

export class MessageDeconfirmLog{
    static coreParse(pc: Command): { res: __MessageDeconfirmLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-deconfirm_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDeconfirmLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-deconfirm_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDeconfirmLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDeconfirmLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDeconfirmLog, missing: string } {
        const ev = this.coreParse({ command: "message-deconfirm_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageNewMemberHint = {
    args: string;
}

export class MessageNewMemberHint{
    static coreParse(pc: Command): { res: __MessageNewMemberHint, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-new_member_hint") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageNewMemberHint): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-new_member_hint",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageNewMemberHint): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageNewMemberHint, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageNewMemberHint, missing: string } {
        const ev = this.coreParse({ command: "message-new_member_hint", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageAddWhitelistPrompt = {
    args: string;
}

export class MessageAddWhitelistPrompt{
    static coreParse(pc: Command): { res: __MessageAddWhitelistPrompt, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-add_whitelist_prompt") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageAddWhitelistPrompt): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-add_whitelist_prompt",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageAddWhitelistPrompt): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageAddWhitelistPrompt, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageAddWhitelistPrompt, missing: string } {
        const ev = this.coreParse({ command: "message-add_whitelist_prompt", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageAddWhitelistSucc = {
    args: string;
}

export class MessageAddWhitelistSucc{
    static coreParse(pc: Command): { res: __MessageAddWhitelistSucc, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-add_whitelist_succ") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageAddWhitelistSucc): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-add_whitelist_succ",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageAddWhitelistSucc): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageAddWhitelistSucc, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageAddWhitelistSucc, missing: string } {
        const ev = this.coreParse({ command: "message-add_whitelist_succ", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageAddWhitelistLog = {
    args: string;
}

export class MessageAddWhitelistLog{
    static coreParse(pc: Command): { res: __MessageAddWhitelistLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-add_whitelist_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageAddWhitelistLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-add_whitelist_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageAddWhitelistLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageAddWhitelistLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageAddWhitelistLog, missing: string } {
        const ev = this.coreParse({ command: "message-add_whitelist_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageRemoveWhitelistPrompt = {
    args: string;
}

export class MessageRemoveWhitelistPrompt{
    static coreParse(pc: Command): { res: __MessageRemoveWhitelistPrompt, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-remove_whitelist_prompt") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageRemoveWhitelistPrompt): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-remove_whitelist_prompt",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageRemoveWhitelistPrompt): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageRemoveWhitelistPrompt, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageRemoveWhitelistPrompt, missing: string } {
        const ev = this.coreParse({ command: "message-remove_whitelist_prompt", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageRemoveWhitelistNotFound = {
    args: string;
}

export class MessageRemoveWhitelistNotFound{
    static coreParse(pc: Command): { res: __MessageRemoveWhitelistNotFound, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-remove_whitelist_not_found") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageRemoveWhitelistNotFound): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-remove_whitelist_not_found",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageRemoveWhitelistNotFound): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageRemoveWhitelistNotFound, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageRemoveWhitelistNotFound, missing: string } {
        const ev = this.coreParse({ command: "message-remove_whitelist_not_found", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageRemoveWhitelistLog = {
    args: string;
}

export class MessageRemoveWhitelistLog{
    static coreParse(pc: Command): { res: __MessageRemoveWhitelistLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-remove_whitelist_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageRemoveWhitelistLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-remove_whitelist_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageRemoveWhitelistLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageRemoveWhitelistLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageRemoveWhitelistLog, missing: string } {
        const ev = this.coreParse({ command: "message-remove_whitelist_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageRemoveWhitelistSucc = {
    args: string;
}

export class MessageRemoveWhitelistSucc{
    static coreParse(pc: Command): { res: __MessageRemoveWhitelistSucc, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-remove_whitelist_succ") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageRemoveWhitelistSucc): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-remove_whitelist_succ",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageRemoveWhitelistSucc): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageRemoveWhitelistSucc, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageRemoveWhitelistSucc, missing: string } {
        const ev = this.coreParse({ command: "message-remove_whitelist_succ", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisHead = {
    args: string;
}

export class MessageWhoisHead{
    static coreParse(pc: Command): { res: __MessageWhoisHead, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_head") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisHead): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_head",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisHead): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisHead, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisHead, missing: string } {
        const ev = this.coreParse({ command: "message-whois_head", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisPrompt = {
    args: string;
}

export class MessageWhoisPrompt{
    static coreParse(pc: Command): { res: __MessageWhoisPrompt, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_prompt") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisPrompt): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_prompt",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisPrompt): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisPrompt, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisPrompt, missing: string } {
        const ev = this.coreParse({ command: "message-whois_prompt", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisNotFound = {
    args: string;
}

export class MessageWhoisNotFound{
    static coreParse(pc: Command): { res: __MessageWhoisNotFound, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_not_found") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisNotFound): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_not_found",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisNotFound): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisNotFound, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisNotFound, missing: string } {
        const ev = this.coreParse({ command: "message-whois_not_found", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisSelf = {
    args: string;
}

export class MessageWhoisSelf{
    static coreParse(pc: Command): { res: __MessageWhoisSelf, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_self") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisSelf): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_self",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisSelf): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisSelf, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisSelf, missing: string } {
        const ev = this.coreParse({ command: "message-whois_self", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisBot = {
    args: string;
}

export class MessageWhoisBot{
    static coreParse(pc: Command): { res: __MessageWhoisBot, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_bot") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisBot): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_bot",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisBot): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisBot, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisBot, missing: string } {
        const ev = this.coreParse({ command: "message-whois_bot", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisHasMw = {
    args: string;
}

export class MessageWhoisHasMw{
    static coreParse(pc: Command): { res: __MessageWhoisHasMw, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_has_mw") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisHasMw): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_has_mw",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisHasMw): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisHasMw, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisHasMw, missing: string } {
        const ev = this.coreParse({ command: "message-whois_has_mw", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisNoMw = {
    args: string;
}

export class MessageWhoisNoMw{
    static coreParse(pc: Command): { res: __MessageWhoisNoMw, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_no_mw") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisNoMw): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_no_mw",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisNoMw): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisNoMw, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisNoMw, missing: string } {
        const ev = this.coreParse({ command: "message-whois_no_mw", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisWhitelisted = {
    args: string;
}

export class MessageWhoisWhitelisted{
    static coreParse(pc: Command): { res: __MessageWhoisWhitelisted, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_whitelisted") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisWhitelisted): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_whitelisted",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisWhitelisted): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisWhitelisted, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisWhitelisted, missing: string } {
        const ev = this.coreParse({ command: "message-whois_whitelisted", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageWhoisTgNameUnavailable = {
    args: string;
}

export class MessageWhoisTgNameUnavailable{
    static coreParse(pc: Command): { res: __MessageWhoisTgNameUnavailable, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-whois_tg_name_unavailable") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageWhoisTgNameUnavailable): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-whois_tg_name_unavailable",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageWhoisTgNameUnavailable): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageWhoisTgNameUnavailable, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageWhoisTgNameUnavailable, missing: string } {
        const ev = this.coreParse({ command: "message-whois_tg_name_unavailable", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageRefuseLog = {
    args: string;
}

export class MessageRefuseLog{
    static coreParse(pc: Command): { res: __MessageRefuseLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-refuse_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageRefuseLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-refuse_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageRefuseLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageRefuseLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageRefuseLog, missing: string } {
        const ev = this.coreParse({ command: "message-refuse_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageAcceptLog = {
    args: string;
}

export class MessageAcceptLog{
    static coreParse(pc: Command): { res: __MessageAcceptLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-accept_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageAcceptLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-accept_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageAcceptLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageAcceptLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageAcceptLog, missing: string } {
        const ev = this.coreParse({ command: "message-accept_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageLiftRestrictionAlert = {
    args: string;
}

export class MessageLiftRestrictionAlert{
    static coreParse(pc: Command): { res: __MessageLiftRestrictionAlert, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-lift_restriction_alert") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageLiftRestrictionAlert): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-lift_restriction_alert",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageLiftRestrictionAlert): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageLiftRestrictionAlert, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageLiftRestrictionAlert, missing: string } {
        const ev = this.coreParse({ command: "message-lift_restriction_alert", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageSilenceAlert = {
    args: string;
}

export class MessageSilenceAlert{
    static coreParse(pc: Command): { res: __MessageSilenceAlert, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-silence_alert") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageSilenceAlert): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-silence_alert",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageSilenceAlert): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageSilenceAlert, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageSilenceAlert, missing: string } {
        const ev = this.coreParse({ command: "message-silence_alert", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageEnable = {
    args: string;
}

export class MessageEnable{
    static coreParse(pc: Command): { res: __MessageEnable, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-enable") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageEnable): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-enable",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageEnable): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageEnable, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageEnable, missing: string } {
        const ev = this.coreParse({ command: "message-enable", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDisable = {
    args: string;
}

export class MessageDisable{
    static coreParse(pc: Command): { res: __MessageDisable, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-disable") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDisable): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-disable",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDisable): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDisable, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDisable, missing: string } {
        const ev = this.coreParse({ command: "message-disable", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageEnableLog = {
    args: string;
}

export class MessageEnableLog{
    static coreParse(pc: Command): { res: __MessageEnableLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-enable_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageEnableLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-enable_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageEnableLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageEnableLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageEnableLog, missing: string } {
        const ev = this.coreParse({ command: "message-enable_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __MessageDisableLog = {
    args: string;
}

export class MessageDisableLog{
    static coreParse(pc: Command): { res: __MessageDisableLog, missing: string }|Error {
        let missing = "";
        if (pc.args[0] === undefined) {
            missing += "args;";
        }
        let bArgs = pc.args[0]!;
        if (pc.command !== "message-disable_log") {
            return new Error("command name mismatch");
        }
        for (let i=0; i<pc.options.length; i+=2) {
            switch(pc.options[i]) {
            }
        }
        return {
            res: {
                args: bArgs, 
            },
            missing: missing,
        }
    }
    static coreEncode(a: __MessageDisableLog): Command {
        const args = [a.args];
        const options = [] as string[];
        return {
            command: "message-disable_log",
            args: args,
            options: options,
        }
    }
    static write(a: __MessageDisableLog): string {
        return CcCore.encode(this.coreEncode(a));
    }
    static parse(s: string): { res: __MessageDisableLog, missing: string }|Error {
        const e = CcCore.parse(s);
        if (e instanceof Error) { return e; }
        if (e.length !== 1) { return new Error("expected exactly 1 command line"); }
        const f = e[0]!;
        return this.coreParse(f);
    }
    static emptyValue(): { res: __MessageDisableLog, missing: string } {
        const ev = this.coreParse({ command: "message-disable_log", options: [], args: [] });
        if (ev instanceof Error) { throw ev; }
        return ev;
    }
}

type __ConfigMix = {
    token: __Token,
    proxy: __Proxy,
    record: __Record,
    groups: __Groups,
    logChannel: __LogChannel,
    mainSite: __MainSite,
    oauthAuthUrl: __OauthAuthUrl,
    oauthQueryUrl: __OauthQueryUrl,
    oauthQueryKey: __OauthQueryKey,
    wikiList: __WikiList,
    blacklist: __Blacklist,
    messageStart: __MessageStart,
    messagePolicy: __MessagePolicy,
    messageInsufficientRight: __MessageInsufficientRight,
    messageGeneralPrompt: __MessageGeneralPrompt,
    messageTelegramIdError: __MessageTelegramIdError,
    messageRestoreSilence: __MessageRestoreSilence,
    messageConfirmAlready: __MessageConfirmAlready,
    messageConfirmOtherTg: __MessageConfirmOtherTg,
    messageConfirmConflict: __MessageConfirmConflict,
    messageConfirmChecking: __MessageConfirmChecking,
    messageConfirmUserNotFound: __MessageConfirmUserNotFound,
    messageConfirmButton: __MessageConfirmButton,
    messageConfirmWait: __MessageConfirmWait,
    messageConfirmConfirming: __MessageConfirmConfirming,
    messageConfirmIneligible: __MessageConfirmIneligible,
    messageConfirmSessionLost: __MessageConfirmSessionLost,
    messageConfirmComplete: __MessageConfirmComplete,
    messageConfirmFailed: __MessageConfirmFailed,
    messageConfirmLog: __MessageConfirmLog,
    messageDeconfirmPrompt: __MessageDeconfirmPrompt,
    messageDeconfirmButton: __MessageDeconfirmButton,
    messageDeconfirmSucc: __MessageDeconfirmSucc,
    messageDeconfirmNotConfirmed: __MessageDeconfirmNotConfirmed,
    messageDeconfirmLog: __MessageDeconfirmLog,
    messageNewMemberHint: __MessageNewMemberHint,
    messageAddWhitelistPrompt: __MessageAddWhitelistPrompt,
    messageAddWhitelistSucc: __MessageAddWhitelistSucc,
    messageAddWhitelistLog: __MessageAddWhitelistLog,
    messageRemoveWhitelistPrompt: __MessageRemoveWhitelistPrompt,
    messageRemoveWhitelistNotFound: __MessageRemoveWhitelistNotFound,
    messageRemoveWhitelistLog: __MessageRemoveWhitelistLog,
    messageRemoveWhitelistSucc: __MessageRemoveWhitelistSucc,
    messageWhoisHead: __MessageWhoisHead,
    messageWhoisPrompt: __MessageWhoisPrompt,
    messageWhoisNotFound: __MessageWhoisNotFound,
    messageWhoisSelf: __MessageWhoisSelf,
    messageWhoisBot: __MessageWhoisBot,
    messageWhoisHasMw: __MessageWhoisHasMw,
    messageWhoisNoMw: __MessageWhoisNoMw,
    messageWhoisWhitelisted: __MessageWhoisWhitelisted,
    messageWhoisTgNameUnavailable: __MessageWhoisTgNameUnavailable,
    messageRefuseLog: __MessageRefuseLog,
    messageAcceptLog: __MessageAcceptLog,
    messageLiftRestrictionAlert: __MessageLiftRestrictionAlert,
    messageSilenceAlert: __MessageSilenceAlert,
    messageEnable: __MessageEnable,
    messageDisable: __MessageDisable,
    messageEnableLog: __MessageEnableLog,
    messageDisableLog: __MessageDisableLog,
}

export class ConfigMix {
    static parse(wiredata: string): { mix: __ConfigMix, missing: string }|Error {
        const res = CcCore.parse(wiredata);
        if (res instanceof Error) { return res; }
        let missing = "";
        let mixMissing = ""
        let bToken: __Token;
        let mToken = "";
        let cToken = 0;
        let bProxy: __Proxy;
        let mProxy = "";
        let cProxy = 0;
        let bRecord: __Record;
        let mRecord = "";
        let cRecord = 0;
        let bGroups: __Groups;
        let mGroups = "";
        let cGroups = 0;
        let bLogChannel: __LogChannel;
        let mLogChannel = "";
        let cLogChannel = 0;
        let bMainSite: __MainSite;
        let mMainSite = "";
        let cMainSite = 0;
        let bOauthAuthUrl: __OauthAuthUrl;
        let mOauthAuthUrl = "";
        let cOauthAuthUrl = 0;
        let bOauthQueryUrl: __OauthQueryUrl;
        let mOauthQueryUrl = "";
        let cOauthQueryUrl = 0;
        let bOauthQueryKey: __OauthQueryKey;
        let mOauthQueryKey = "";
        let cOauthQueryKey = 0;
        let bWikiList: __WikiList;
        let mWikiList = "";
        let cWikiList = 0;
        let bBlacklist: __Blacklist;
        let mBlacklist = "";
        let cBlacklist = 0;
        let bMessageStart: __MessageStart;
        let mMessageStart = "";
        let cMessageStart = 0;
        let bMessagePolicy: __MessagePolicy;
        let mMessagePolicy = "";
        let cMessagePolicy = 0;
        let bMessageInsufficientRight: __MessageInsufficientRight;
        let mMessageInsufficientRight = "";
        let cMessageInsufficientRight = 0;
        let bMessageGeneralPrompt: __MessageGeneralPrompt;
        let mMessageGeneralPrompt = "";
        let cMessageGeneralPrompt = 0;
        let bMessageTelegramIdError: __MessageTelegramIdError;
        let mMessageTelegramIdError = "";
        let cMessageTelegramIdError = 0;
        let bMessageRestoreSilence: __MessageRestoreSilence;
        let mMessageRestoreSilence = "";
        let cMessageRestoreSilence = 0;
        let bMessageConfirmAlready: __MessageConfirmAlready;
        let mMessageConfirmAlready = "";
        let cMessageConfirmAlready = 0;
        let bMessageConfirmOtherTg: __MessageConfirmOtherTg;
        let mMessageConfirmOtherTg = "";
        let cMessageConfirmOtherTg = 0;
        let bMessageConfirmConflict: __MessageConfirmConflict;
        let mMessageConfirmConflict = "";
        let cMessageConfirmConflict = 0;
        let bMessageConfirmChecking: __MessageConfirmChecking;
        let mMessageConfirmChecking = "";
        let cMessageConfirmChecking = 0;
        let bMessageConfirmUserNotFound: __MessageConfirmUserNotFound;
        let mMessageConfirmUserNotFound = "";
        let cMessageConfirmUserNotFound = 0;
        let bMessageConfirmButton: __MessageConfirmButton;
        let mMessageConfirmButton = "";
        let cMessageConfirmButton = 0;
        let bMessageConfirmWait: __MessageConfirmWait;
        let mMessageConfirmWait = "";
        let cMessageConfirmWait = 0;
        let bMessageConfirmConfirming: __MessageConfirmConfirming;
        let mMessageConfirmConfirming = "";
        let cMessageConfirmConfirming = 0;
        let bMessageConfirmIneligible: __MessageConfirmIneligible;
        let mMessageConfirmIneligible = "";
        let cMessageConfirmIneligible = 0;
        let bMessageConfirmSessionLost: __MessageConfirmSessionLost;
        let mMessageConfirmSessionLost = "";
        let cMessageConfirmSessionLost = 0;
        let bMessageConfirmComplete: __MessageConfirmComplete;
        let mMessageConfirmComplete = "";
        let cMessageConfirmComplete = 0;
        let bMessageConfirmFailed: __MessageConfirmFailed;
        let mMessageConfirmFailed = "";
        let cMessageConfirmFailed = 0;
        let bMessageConfirmLog: __MessageConfirmLog;
        let mMessageConfirmLog = "";
        let cMessageConfirmLog = 0;
        let bMessageDeconfirmPrompt: __MessageDeconfirmPrompt;
        let mMessageDeconfirmPrompt = "";
        let cMessageDeconfirmPrompt = 0;
        let bMessageDeconfirmButton: __MessageDeconfirmButton;
        let mMessageDeconfirmButton = "";
        let cMessageDeconfirmButton = 0;
        let bMessageDeconfirmSucc: __MessageDeconfirmSucc;
        let mMessageDeconfirmSucc = "";
        let cMessageDeconfirmSucc = 0;
        let bMessageDeconfirmNotConfirmed: __MessageDeconfirmNotConfirmed;
        let mMessageDeconfirmNotConfirmed = "";
        let cMessageDeconfirmNotConfirmed = 0;
        let bMessageDeconfirmLog: __MessageDeconfirmLog;
        let mMessageDeconfirmLog = "";
        let cMessageDeconfirmLog = 0;
        let bMessageNewMemberHint: __MessageNewMemberHint;
        let mMessageNewMemberHint = "";
        let cMessageNewMemberHint = 0;
        let bMessageAddWhitelistPrompt: __MessageAddWhitelistPrompt;
        let mMessageAddWhitelistPrompt = "";
        let cMessageAddWhitelistPrompt = 0;
        let bMessageAddWhitelistSucc: __MessageAddWhitelistSucc;
        let mMessageAddWhitelistSucc = "";
        let cMessageAddWhitelistSucc = 0;
        let bMessageAddWhitelistLog: __MessageAddWhitelistLog;
        let mMessageAddWhitelistLog = "";
        let cMessageAddWhitelistLog = 0;
        let bMessageRemoveWhitelistPrompt: __MessageRemoveWhitelistPrompt;
        let mMessageRemoveWhitelistPrompt = "";
        let cMessageRemoveWhitelistPrompt = 0;
        let bMessageRemoveWhitelistNotFound: __MessageRemoveWhitelistNotFound;
        let mMessageRemoveWhitelistNotFound = "";
        let cMessageRemoveWhitelistNotFound = 0;
        let bMessageRemoveWhitelistLog: __MessageRemoveWhitelistLog;
        let mMessageRemoveWhitelistLog = "";
        let cMessageRemoveWhitelistLog = 0;
        let bMessageRemoveWhitelistSucc: __MessageRemoveWhitelistSucc;
        let mMessageRemoveWhitelistSucc = "";
        let cMessageRemoveWhitelistSucc = 0;
        let bMessageWhoisHead: __MessageWhoisHead;
        let mMessageWhoisHead = "";
        let cMessageWhoisHead = 0;
        let bMessageWhoisPrompt: __MessageWhoisPrompt;
        let mMessageWhoisPrompt = "";
        let cMessageWhoisPrompt = 0;
        let bMessageWhoisNotFound: __MessageWhoisNotFound;
        let mMessageWhoisNotFound = "";
        let cMessageWhoisNotFound = 0;
        let bMessageWhoisSelf: __MessageWhoisSelf;
        let mMessageWhoisSelf = "";
        let cMessageWhoisSelf = 0;
        let bMessageWhoisBot: __MessageWhoisBot;
        let mMessageWhoisBot = "";
        let cMessageWhoisBot = 0;
        let bMessageWhoisHasMw: __MessageWhoisHasMw;
        let mMessageWhoisHasMw = "";
        let cMessageWhoisHasMw = 0;
        let bMessageWhoisNoMw: __MessageWhoisNoMw;
        let mMessageWhoisNoMw = "";
        let cMessageWhoisNoMw = 0;
        let bMessageWhoisWhitelisted: __MessageWhoisWhitelisted;
        let mMessageWhoisWhitelisted = "";
        let cMessageWhoisWhitelisted = 0;
        let bMessageWhoisTgNameUnavailable: __MessageWhoisTgNameUnavailable;
        let mMessageWhoisTgNameUnavailable = "";
        let cMessageWhoisTgNameUnavailable = 0;
        let bMessageRefuseLog: __MessageRefuseLog;
        let mMessageRefuseLog = "";
        let cMessageRefuseLog = 0;
        let bMessageAcceptLog: __MessageAcceptLog;
        let mMessageAcceptLog = "";
        let cMessageAcceptLog = 0;
        let bMessageLiftRestrictionAlert: __MessageLiftRestrictionAlert;
        let mMessageLiftRestrictionAlert = "";
        let cMessageLiftRestrictionAlert = 0;
        let bMessageSilenceAlert: __MessageSilenceAlert;
        let mMessageSilenceAlert = "";
        let cMessageSilenceAlert = 0;
        let bMessageEnable: __MessageEnable;
        let mMessageEnable = "";
        let cMessageEnable = 0;
        let bMessageDisable: __MessageDisable;
        let mMessageDisable = "";
        let cMessageDisable = 0;
        let bMessageEnableLog: __MessageEnableLog;
        let mMessageEnableLog = "";
        let cMessageEnableLog = 0;
        let bMessageDisableLog: __MessageDisableLog;
        let mMessageDisableLog = "";
        let cMessageDisableLog = 0;
        for (const cmd of res) {
            switch (cmd.command) {
            case "token":
                const parsed0 = Token.coreParse(cmd);
                if (parsed0 instanceof Error) { return parsed0; }
                bToken = parsed0.res;
                mToken += parsed0.missing;
                cToken += 1;
            break;
            case "proxy":
                const parsed1 = Proxy.coreParse(cmd);
                if (parsed1 instanceof Error) { return parsed1; }
                bProxy = parsed1.res;
                mProxy += parsed1.missing;
                cProxy += 1;
            break;
            case "record":
                const parsed2 = Record.coreParse(cmd);
                if (parsed2 instanceof Error) { return parsed2; }
                bRecord = parsed2.res;
                mRecord += parsed2.missing;
                cRecord += 1;
            break;
            case "groups":
                const parsed3 = Groups.coreParse(cmd);
                if (parsed3 instanceof Error) { return parsed3; }
                bGroups = parsed3.res;
                mGroups += parsed3.missing;
                cGroups += 1;
            break;
            case "log_channel":
                const parsed4 = LogChannel.coreParse(cmd);
                if (parsed4 instanceof Error) { return parsed4; }
                bLogChannel = parsed4.res;
                mLogChannel += parsed4.missing;
                cLogChannel += 1;
            break;
            case "main_site":
                const parsed5 = MainSite.coreParse(cmd);
                if (parsed5 instanceof Error) { return parsed5; }
                bMainSite = parsed5.res;
                mMainSite += parsed5.missing;
                cMainSite += 1;
            break;
            case "oauth_auth_url":
                const parsed6 = OauthAuthUrl.coreParse(cmd);
                if (parsed6 instanceof Error) { return parsed6; }
                bOauthAuthUrl = parsed6.res;
                mOauthAuthUrl += parsed6.missing;
                cOauthAuthUrl += 1;
            break;
            case "oauth_query_url":
                const parsed7 = OauthQueryUrl.coreParse(cmd);
                if (parsed7 instanceof Error) { return parsed7; }
                bOauthQueryUrl = parsed7.res;
                mOauthQueryUrl += parsed7.missing;
                cOauthQueryUrl += 1;
            break;
            case "oauth_query_key":
                const parsed8 = OauthQueryKey.coreParse(cmd);
                if (parsed8 instanceof Error) { return parsed8; }
                bOauthQueryKey = parsed8.res;
                mOauthQueryKey += parsed8.missing;
                cOauthQueryKey += 1;
            break;
            case "wiki_list":
                const parsed9 = WikiList.coreParse(cmd);
                if (parsed9 instanceof Error) { return parsed9; }
                bWikiList = parsed9.res;
                mWikiList += parsed9.missing;
                cWikiList += 1;
            break;
            case "blacklist":
                const parsed10 = Blacklist.coreParse(cmd);
                if (parsed10 instanceof Error) { return parsed10; }
                bBlacklist = parsed10.res;
                mBlacklist += parsed10.missing;
                cBlacklist += 1;
            break;
            case "message-start":
                const parsed11 = MessageStart.coreParse(cmd);
                if (parsed11 instanceof Error) { return parsed11; }
                bMessageStart = parsed11.res;
                mMessageStart += parsed11.missing;
                cMessageStart += 1;
            break;
            case "message-policy":
                const parsed12 = MessagePolicy.coreParse(cmd);
                if (parsed12 instanceof Error) { return parsed12; }
                bMessagePolicy = parsed12.res;
                mMessagePolicy += parsed12.missing;
                cMessagePolicy += 1;
            break;
            case "message-insufficient_right":
                const parsed13 = MessageInsufficientRight.coreParse(cmd);
                if (parsed13 instanceof Error) { return parsed13; }
                bMessageInsufficientRight = parsed13.res;
                mMessageInsufficientRight += parsed13.missing;
                cMessageInsufficientRight += 1;
            break;
            case "message-general_prompt":
                const parsed14 = MessageGeneralPrompt.coreParse(cmd);
                if (parsed14 instanceof Error) { return parsed14; }
                bMessageGeneralPrompt = parsed14.res;
                mMessageGeneralPrompt += parsed14.missing;
                cMessageGeneralPrompt += 1;
            break;
            case "message-telegram_id_error":
                const parsed15 = MessageTelegramIdError.coreParse(cmd);
                if (parsed15 instanceof Error) { return parsed15; }
                bMessageTelegramIdError = parsed15.res;
                mMessageTelegramIdError += parsed15.missing;
                cMessageTelegramIdError += 1;
            break;
            case "message-restore_silence":
                const parsed16 = MessageRestoreSilence.coreParse(cmd);
                if (parsed16 instanceof Error) { return parsed16; }
                bMessageRestoreSilence = parsed16.res;
                mMessageRestoreSilence += parsed16.missing;
                cMessageRestoreSilence += 1;
            break;
            case "message-confirm_already":
                const parsed17 = MessageConfirmAlready.coreParse(cmd);
                if (parsed17 instanceof Error) { return parsed17; }
                bMessageConfirmAlready = parsed17.res;
                mMessageConfirmAlready += parsed17.missing;
                cMessageConfirmAlready += 1;
            break;
            case "message-confirm_other_tg":
                const parsed18 = MessageConfirmOtherTg.coreParse(cmd);
                if (parsed18 instanceof Error) { return parsed18; }
                bMessageConfirmOtherTg = parsed18.res;
                mMessageConfirmOtherTg += parsed18.missing;
                cMessageConfirmOtherTg += 1;
            break;
            case "message-confirm_conflict":
                const parsed19 = MessageConfirmConflict.coreParse(cmd);
                if (parsed19 instanceof Error) { return parsed19; }
                bMessageConfirmConflict = parsed19.res;
                mMessageConfirmConflict += parsed19.missing;
                cMessageConfirmConflict += 1;
            break;
            case "message-confirm_checking":
                const parsed20 = MessageConfirmChecking.coreParse(cmd);
                if (parsed20 instanceof Error) { return parsed20; }
                bMessageConfirmChecking = parsed20.res;
                mMessageConfirmChecking += parsed20.missing;
                cMessageConfirmChecking += 1;
            break;
            case "message-confirm_user_not_found":
                const parsed21 = MessageConfirmUserNotFound.coreParse(cmd);
                if (parsed21 instanceof Error) { return parsed21; }
                bMessageConfirmUserNotFound = parsed21.res;
                mMessageConfirmUserNotFound += parsed21.missing;
                cMessageConfirmUserNotFound += 1;
            break;
            case "message-confirm_button":
                const parsed22 = MessageConfirmButton.coreParse(cmd);
                if (parsed22 instanceof Error) { return parsed22; }
                bMessageConfirmButton = parsed22.res;
                mMessageConfirmButton += parsed22.missing;
                cMessageConfirmButton += 1;
            break;
            case "message-confirm_wait":
                const parsed23 = MessageConfirmWait.coreParse(cmd);
                if (parsed23 instanceof Error) { return parsed23; }
                bMessageConfirmWait = parsed23.res;
                mMessageConfirmWait += parsed23.missing;
                cMessageConfirmWait += 1;
            break;
            case "message-confirm_confirming":
                const parsed24 = MessageConfirmConfirming.coreParse(cmd);
                if (parsed24 instanceof Error) { return parsed24; }
                bMessageConfirmConfirming = parsed24.res;
                mMessageConfirmConfirming += parsed24.missing;
                cMessageConfirmConfirming += 1;
            break;
            case "message-confirm_ineligible":
                const parsed25 = MessageConfirmIneligible.coreParse(cmd);
                if (parsed25 instanceof Error) { return parsed25; }
                bMessageConfirmIneligible = parsed25.res;
                mMessageConfirmIneligible += parsed25.missing;
                cMessageConfirmIneligible += 1;
            break;
            case "message-confirm_session_lost":
                const parsed26 = MessageConfirmSessionLost.coreParse(cmd);
                if (parsed26 instanceof Error) { return parsed26; }
                bMessageConfirmSessionLost = parsed26.res;
                mMessageConfirmSessionLost += parsed26.missing;
                cMessageConfirmSessionLost += 1;
            break;
            case "message-confirm_complete":
                const parsed27 = MessageConfirmComplete.coreParse(cmd);
                if (parsed27 instanceof Error) { return parsed27; }
                bMessageConfirmComplete = parsed27.res;
                mMessageConfirmComplete += parsed27.missing;
                cMessageConfirmComplete += 1;
            break;
            case "message-confirm_failed":
                const parsed28 = MessageConfirmFailed.coreParse(cmd);
                if (parsed28 instanceof Error) { return parsed28; }
                bMessageConfirmFailed = parsed28.res;
                mMessageConfirmFailed += parsed28.missing;
                cMessageConfirmFailed += 1;
            break;
            case "message-confirm_log":
                const parsed29 = MessageConfirmLog.coreParse(cmd);
                if (parsed29 instanceof Error) { return parsed29; }
                bMessageConfirmLog = parsed29.res;
                mMessageConfirmLog += parsed29.missing;
                cMessageConfirmLog += 1;
            break;
            case "message-deconfirm_prompt":
                const parsed30 = MessageDeconfirmPrompt.coreParse(cmd);
                if (parsed30 instanceof Error) { return parsed30; }
                bMessageDeconfirmPrompt = parsed30.res;
                mMessageDeconfirmPrompt += parsed30.missing;
                cMessageDeconfirmPrompt += 1;
            break;
            case "message-deconfirm_button":
                const parsed31 = MessageDeconfirmButton.coreParse(cmd);
                if (parsed31 instanceof Error) { return parsed31; }
                bMessageDeconfirmButton = parsed31.res;
                mMessageDeconfirmButton += parsed31.missing;
                cMessageDeconfirmButton += 1;
            break;
            case "message-deconfirm_succ":
                const parsed32 = MessageDeconfirmSucc.coreParse(cmd);
                if (parsed32 instanceof Error) { return parsed32; }
                bMessageDeconfirmSucc = parsed32.res;
                mMessageDeconfirmSucc += parsed32.missing;
                cMessageDeconfirmSucc += 1;
            break;
            case "message-deconfirm_not_confirmed":
                const parsed33 = MessageDeconfirmNotConfirmed.coreParse(cmd);
                if (parsed33 instanceof Error) { return parsed33; }
                bMessageDeconfirmNotConfirmed = parsed33.res;
                mMessageDeconfirmNotConfirmed += parsed33.missing;
                cMessageDeconfirmNotConfirmed += 1;
            break;
            case "message-deconfirm_log":
                const parsed34 = MessageDeconfirmLog.coreParse(cmd);
                if (parsed34 instanceof Error) { return parsed34; }
                bMessageDeconfirmLog = parsed34.res;
                mMessageDeconfirmLog += parsed34.missing;
                cMessageDeconfirmLog += 1;
            break;
            case "message-new_member_hint":
                const parsed35 = MessageNewMemberHint.coreParse(cmd);
                if (parsed35 instanceof Error) { return parsed35; }
                bMessageNewMemberHint = parsed35.res;
                mMessageNewMemberHint += parsed35.missing;
                cMessageNewMemberHint += 1;
            break;
            case "message-add_whitelist_prompt":
                const parsed36 = MessageAddWhitelistPrompt.coreParse(cmd);
                if (parsed36 instanceof Error) { return parsed36; }
                bMessageAddWhitelistPrompt = parsed36.res;
                mMessageAddWhitelistPrompt += parsed36.missing;
                cMessageAddWhitelistPrompt += 1;
            break;
            case "message-add_whitelist_succ":
                const parsed37 = MessageAddWhitelistSucc.coreParse(cmd);
                if (parsed37 instanceof Error) { return parsed37; }
                bMessageAddWhitelistSucc = parsed37.res;
                mMessageAddWhitelistSucc += parsed37.missing;
                cMessageAddWhitelistSucc += 1;
            break;
            case "message-add_whitelist_log":
                const parsed38 = MessageAddWhitelistLog.coreParse(cmd);
                if (parsed38 instanceof Error) { return parsed38; }
                bMessageAddWhitelistLog = parsed38.res;
                mMessageAddWhitelistLog += parsed38.missing;
                cMessageAddWhitelistLog += 1;
            break;
            case "message-remove_whitelist_prompt":
                const parsed39 = MessageRemoveWhitelistPrompt.coreParse(cmd);
                if (parsed39 instanceof Error) { return parsed39; }
                bMessageRemoveWhitelistPrompt = parsed39.res;
                mMessageRemoveWhitelistPrompt += parsed39.missing;
                cMessageRemoveWhitelistPrompt += 1;
            break;
            case "message-remove_whitelist_not_found":
                const parsed40 = MessageRemoveWhitelistNotFound.coreParse(cmd);
                if (parsed40 instanceof Error) { return parsed40; }
                bMessageRemoveWhitelistNotFound = parsed40.res;
                mMessageRemoveWhitelistNotFound += parsed40.missing;
                cMessageRemoveWhitelistNotFound += 1;
            break;
            case "message-remove_whitelist_log":
                const parsed41 = MessageRemoveWhitelistLog.coreParse(cmd);
                if (parsed41 instanceof Error) { return parsed41; }
                bMessageRemoveWhitelistLog = parsed41.res;
                mMessageRemoveWhitelistLog += parsed41.missing;
                cMessageRemoveWhitelistLog += 1;
            break;
            case "message-remove_whitelist_succ":
                const parsed42 = MessageRemoveWhitelistSucc.coreParse(cmd);
                if (parsed42 instanceof Error) { return parsed42; }
                bMessageRemoveWhitelistSucc = parsed42.res;
                mMessageRemoveWhitelistSucc += parsed42.missing;
                cMessageRemoveWhitelistSucc += 1;
            break;
            case "message-whois_head":
                const parsed43 = MessageWhoisHead.coreParse(cmd);
                if (parsed43 instanceof Error) { return parsed43; }
                bMessageWhoisHead = parsed43.res;
                mMessageWhoisHead += parsed43.missing;
                cMessageWhoisHead += 1;
            break;
            case "message-whois_prompt":
                const parsed44 = MessageWhoisPrompt.coreParse(cmd);
                if (parsed44 instanceof Error) { return parsed44; }
                bMessageWhoisPrompt = parsed44.res;
                mMessageWhoisPrompt += parsed44.missing;
                cMessageWhoisPrompt += 1;
            break;
            case "message-whois_not_found":
                const parsed45 = MessageWhoisNotFound.coreParse(cmd);
                if (parsed45 instanceof Error) { return parsed45; }
                bMessageWhoisNotFound = parsed45.res;
                mMessageWhoisNotFound += parsed45.missing;
                cMessageWhoisNotFound += 1;
            break;
            case "message-whois_self":
                const parsed46 = MessageWhoisSelf.coreParse(cmd);
                if (parsed46 instanceof Error) { return parsed46; }
                bMessageWhoisSelf = parsed46.res;
                mMessageWhoisSelf += parsed46.missing;
                cMessageWhoisSelf += 1;
            break;
            case "message-whois_bot":
                const parsed47 = MessageWhoisBot.coreParse(cmd);
                if (parsed47 instanceof Error) { return parsed47; }
                bMessageWhoisBot = parsed47.res;
                mMessageWhoisBot += parsed47.missing;
                cMessageWhoisBot += 1;
            break;
            case "message-whois_has_mw":
                const parsed48 = MessageWhoisHasMw.coreParse(cmd);
                if (parsed48 instanceof Error) { return parsed48; }
                bMessageWhoisHasMw = parsed48.res;
                mMessageWhoisHasMw += parsed48.missing;
                cMessageWhoisHasMw += 1;
            break;
            case "message-whois_no_mw":
                const parsed49 = MessageWhoisNoMw.coreParse(cmd);
                if (parsed49 instanceof Error) { return parsed49; }
                bMessageWhoisNoMw = parsed49.res;
                mMessageWhoisNoMw += parsed49.missing;
                cMessageWhoisNoMw += 1;
            break;
            case "message-whois_whitelisted":
                const parsed50 = MessageWhoisWhitelisted.coreParse(cmd);
                if (parsed50 instanceof Error) { return parsed50; }
                bMessageWhoisWhitelisted = parsed50.res;
                mMessageWhoisWhitelisted += parsed50.missing;
                cMessageWhoisWhitelisted += 1;
            break;
            case "message-whois_tg_name_unavailable":
                const parsed51 = MessageWhoisTgNameUnavailable.coreParse(cmd);
                if (parsed51 instanceof Error) { return parsed51; }
                bMessageWhoisTgNameUnavailable = parsed51.res;
                mMessageWhoisTgNameUnavailable += parsed51.missing;
                cMessageWhoisTgNameUnavailable += 1;
            break;
            case "message-refuse_log":
                const parsed52 = MessageRefuseLog.coreParse(cmd);
                if (parsed52 instanceof Error) { return parsed52; }
                bMessageRefuseLog = parsed52.res;
                mMessageRefuseLog += parsed52.missing;
                cMessageRefuseLog += 1;
            break;
            case "message-accept_log":
                const parsed53 = MessageAcceptLog.coreParse(cmd);
                if (parsed53 instanceof Error) { return parsed53; }
                bMessageAcceptLog = parsed53.res;
                mMessageAcceptLog += parsed53.missing;
                cMessageAcceptLog += 1;
            break;
            case "message-lift_restriction_alert":
                const parsed54 = MessageLiftRestrictionAlert.coreParse(cmd);
                if (parsed54 instanceof Error) { return parsed54; }
                bMessageLiftRestrictionAlert = parsed54.res;
                mMessageLiftRestrictionAlert += parsed54.missing;
                cMessageLiftRestrictionAlert += 1;
            break;
            case "message-silence_alert":
                const parsed55 = MessageSilenceAlert.coreParse(cmd);
                if (parsed55 instanceof Error) { return parsed55; }
                bMessageSilenceAlert = parsed55.res;
                mMessageSilenceAlert += parsed55.missing;
                cMessageSilenceAlert += 1;
            break;
            case "message-enable":
                const parsed56 = MessageEnable.coreParse(cmd);
                if (parsed56 instanceof Error) { return parsed56; }
                bMessageEnable = parsed56.res;
                mMessageEnable += parsed56.missing;
                cMessageEnable += 1;
            break;
            case "message-disable":
                const parsed57 = MessageDisable.coreParse(cmd);
                if (parsed57 instanceof Error) { return parsed57; }
                bMessageDisable = parsed57.res;
                mMessageDisable += parsed57.missing;
                cMessageDisable += 1;
            break;
            case "message-enable_log":
                const parsed58 = MessageEnableLog.coreParse(cmd);
                if (parsed58 instanceof Error) { return parsed58; }
                bMessageEnableLog = parsed58.res;
                mMessageEnableLog += parsed58.missing;
                cMessageEnableLog += 1;
            break;
            case "message-disable_log":
                const parsed59 = MessageDisableLog.coreParse(cmd);
                if (parsed59 instanceof Error) { return parsed59; }
                bMessageDisableLog = parsed59.res;
                mMessageDisableLog += parsed59.missing;
                cMessageDisableLog += 1;
            break;
            }
        }
        if (cToken < 1) {
            mixMissing += "token;";
            const { res, missing } = Token.emptyValue();
            bToken = res;
        }
        if (cProxy < 1) {
            mixMissing += "proxy;";
            const { res, missing } = Proxy.emptyValue();
            bProxy = res;
        }
        if (cRecord < 1) {
            mixMissing += "record;";
            const { res, missing } = Record.emptyValue();
            bRecord = res;
        }
        if (cGroups < 1) {
            mixMissing += "groups;";
            const { res, missing } = Groups.emptyValue();
            bGroups = res;
        }
        if (cLogChannel < 1) {
            mixMissing += "log_channel;";
            const { res, missing } = LogChannel.emptyValue();
            bLogChannel = res;
        }
        if (cMainSite < 1) {
            mixMissing += "main_site;";
            const { res, missing } = MainSite.emptyValue();
            bMainSite = res;
        }
        if (cOauthAuthUrl < 1) {
            mixMissing += "oauth_auth_url;";
            const { res, missing } = OauthAuthUrl.emptyValue();
            bOauthAuthUrl = res;
        }
        if (cOauthQueryUrl < 1) {
            mixMissing += "oauth_query_url;";
            const { res, missing } = OauthQueryUrl.emptyValue();
            bOauthQueryUrl = res;
        }
        if (cOauthQueryKey < 1) {
            mixMissing += "oauth_query_key;";
            const { res, missing } = OauthQueryKey.emptyValue();
            bOauthQueryKey = res;
        }
        if (cWikiList < 1) {
            mixMissing += "wiki_list;";
            const { res, missing } = WikiList.emptyValue();
            bWikiList = res;
        }
        if (cBlacklist < 1) {
            mixMissing += "blacklist;";
            const { res, missing } = Blacklist.emptyValue();
            bBlacklist = res;
        }
        if (cMessageStart < 1) {
            mixMissing += "message-start;";
            const { res, missing } = MessageStart.emptyValue();
            bMessageStart = res;
        }
        if (cMessagePolicy < 1) {
            mixMissing += "message-policy;";
            const { res, missing } = MessagePolicy.emptyValue();
            bMessagePolicy = res;
        }
        if (cMessageInsufficientRight < 1) {
            mixMissing += "message-insufficient_right;";
            const { res, missing } = MessageInsufficientRight.emptyValue();
            bMessageInsufficientRight = res;
        }
        if (cMessageGeneralPrompt < 1) {
            mixMissing += "message-general_prompt;";
            const { res, missing } = MessageGeneralPrompt.emptyValue();
            bMessageGeneralPrompt = res;
        }
        if (cMessageTelegramIdError < 1) {
            mixMissing += "message-telegram_id_error;";
            const { res, missing } = MessageTelegramIdError.emptyValue();
            bMessageTelegramIdError = res;
        }
        if (cMessageRestoreSilence < 1) {
            mixMissing += "message-restore_silence;";
            const { res, missing } = MessageRestoreSilence.emptyValue();
            bMessageRestoreSilence = res;
        }
        if (cMessageConfirmAlready < 1) {
            mixMissing += "message-confirm_already;";
            const { res, missing } = MessageConfirmAlready.emptyValue();
            bMessageConfirmAlready = res;
        }
        if (cMessageConfirmOtherTg < 1) {
            mixMissing += "message-confirm_other_tg;";
            const { res, missing } = MessageConfirmOtherTg.emptyValue();
            bMessageConfirmOtherTg = res;
        }
        if (cMessageConfirmConflict < 1) {
            mixMissing += "message-confirm_conflict;";
            const { res, missing } = MessageConfirmConflict.emptyValue();
            bMessageConfirmConflict = res;
        }
        if (cMessageConfirmChecking < 1) {
            mixMissing += "message-confirm_checking;";
            const { res, missing } = MessageConfirmChecking.emptyValue();
            bMessageConfirmChecking = res;
        }
        if (cMessageConfirmUserNotFound < 1) {
            mixMissing += "message-confirm_user_not_found;";
            const { res, missing } = MessageConfirmUserNotFound.emptyValue();
            bMessageConfirmUserNotFound = res;
        }
        if (cMessageConfirmButton < 1) {
            mixMissing += "message-confirm_button;";
            const { res, missing } = MessageConfirmButton.emptyValue();
            bMessageConfirmButton = res;
        }
        if (cMessageConfirmWait < 1) {
            mixMissing += "message-confirm_wait;";
            const { res, missing } = MessageConfirmWait.emptyValue();
            bMessageConfirmWait = res;
        }
        if (cMessageConfirmConfirming < 1) {
            mixMissing += "message-confirm_confirming;";
            const { res, missing } = MessageConfirmConfirming.emptyValue();
            bMessageConfirmConfirming = res;
        }
        if (cMessageConfirmIneligible < 1) {
            mixMissing += "message-confirm_ineligible;";
            const { res, missing } = MessageConfirmIneligible.emptyValue();
            bMessageConfirmIneligible = res;
        }
        if (cMessageConfirmSessionLost < 1) {
            mixMissing += "message-confirm_session_lost;";
            const { res, missing } = MessageConfirmSessionLost.emptyValue();
            bMessageConfirmSessionLost = res;
        }
        if (cMessageConfirmComplete < 1) {
            mixMissing += "message-confirm_complete;";
            const { res, missing } = MessageConfirmComplete.emptyValue();
            bMessageConfirmComplete = res;
        }
        if (cMessageConfirmFailed < 1) {
            mixMissing += "message-confirm_failed;";
            const { res, missing } = MessageConfirmFailed.emptyValue();
            bMessageConfirmFailed = res;
        }
        if (cMessageConfirmLog < 1) {
            mixMissing += "message-confirm_log;";
            const { res, missing } = MessageConfirmLog.emptyValue();
            bMessageConfirmLog = res;
        }
        if (cMessageDeconfirmPrompt < 1) {
            mixMissing += "message-deconfirm_prompt;";
            const { res, missing } = MessageDeconfirmPrompt.emptyValue();
            bMessageDeconfirmPrompt = res;
        }
        if (cMessageDeconfirmButton < 1) {
            mixMissing += "message-deconfirm_button;";
            const { res, missing } = MessageDeconfirmButton.emptyValue();
            bMessageDeconfirmButton = res;
        }
        if (cMessageDeconfirmSucc < 1) {
            mixMissing += "message-deconfirm_succ;";
            const { res, missing } = MessageDeconfirmSucc.emptyValue();
            bMessageDeconfirmSucc = res;
        }
        if (cMessageDeconfirmNotConfirmed < 1) {
            mixMissing += "message-deconfirm_not_confirmed;";
            const { res, missing } = MessageDeconfirmNotConfirmed.emptyValue();
            bMessageDeconfirmNotConfirmed = res;
        }
        if (cMessageDeconfirmLog < 1) {
            mixMissing += "message-deconfirm_log;";
            const { res, missing } = MessageDeconfirmLog.emptyValue();
            bMessageDeconfirmLog = res;
        }
        if (cMessageNewMemberHint < 1) {
            mixMissing += "message-new_member_hint;";
            const { res, missing } = MessageNewMemberHint.emptyValue();
            bMessageNewMemberHint = res;
        }
        if (cMessageAddWhitelistPrompt < 1) {
            mixMissing += "message-add_whitelist_prompt;";
            const { res, missing } = MessageAddWhitelistPrompt.emptyValue();
            bMessageAddWhitelistPrompt = res;
        }
        if (cMessageAddWhitelistSucc < 1) {
            mixMissing += "message-add_whitelist_succ;";
            const { res, missing } = MessageAddWhitelistSucc.emptyValue();
            bMessageAddWhitelistSucc = res;
        }
        if (cMessageAddWhitelistLog < 1) {
            mixMissing += "message-add_whitelist_log;";
            const { res, missing } = MessageAddWhitelistLog.emptyValue();
            bMessageAddWhitelistLog = res;
        }
        if (cMessageRemoveWhitelistPrompt < 1) {
            mixMissing += "message-remove_whitelist_prompt;";
            const { res, missing } = MessageRemoveWhitelistPrompt.emptyValue();
            bMessageRemoveWhitelistPrompt = res;
        }
        if (cMessageRemoveWhitelistNotFound < 1) {
            mixMissing += "message-remove_whitelist_not_found;";
            const { res, missing } = MessageRemoveWhitelistNotFound.emptyValue();
            bMessageRemoveWhitelistNotFound = res;
        }
        if (cMessageRemoveWhitelistLog < 1) {
            mixMissing += "message-remove_whitelist_log;";
            const { res, missing } = MessageRemoveWhitelistLog.emptyValue();
            bMessageRemoveWhitelistLog = res;
        }
        if (cMessageRemoveWhitelistSucc < 1) {
            mixMissing += "message-remove_whitelist_succ;";
            const { res, missing } = MessageRemoveWhitelistSucc.emptyValue();
            bMessageRemoveWhitelistSucc = res;
        }
        if (cMessageWhoisHead < 1) {
            mixMissing += "message-whois_head;";
            const { res, missing } = MessageWhoisHead.emptyValue();
            bMessageWhoisHead = res;
        }
        if (cMessageWhoisPrompt < 1) {
            mixMissing += "message-whois_prompt;";
            const { res, missing } = MessageWhoisPrompt.emptyValue();
            bMessageWhoisPrompt = res;
        }
        if (cMessageWhoisNotFound < 1) {
            mixMissing += "message-whois_not_found;";
            const { res, missing } = MessageWhoisNotFound.emptyValue();
            bMessageWhoisNotFound = res;
        }
        if (cMessageWhoisSelf < 1) {
            mixMissing += "message-whois_self;";
            const { res, missing } = MessageWhoisSelf.emptyValue();
            bMessageWhoisSelf = res;
        }
        if (cMessageWhoisBot < 1) {
            mixMissing += "message-whois_bot;";
            const { res, missing } = MessageWhoisBot.emptyValue();
            bMessageWhoisBot = res;
        }
        if (cMessageWhoisHasMw < 1) {
            mixMissing += "message-whois_has_mw;";
            const { res, missing } = MessageWhoisHasMw.emptyValue();
            bMessageWhoisHasMw = res;
        }
        if (cMessageWhoisNoMw < 1) {
            mixMissing += "message-whois_no_mw;";
            const { res, missing } = MessageWhoisNoMw.emptyValue();
            bMessageWhoisNoMw = res;
        }
        if (cMessageWhoisWhitelisted < 1) {
            mixMissing += "message-whois_whitelisted;";
            const { res, missing } = MessageWhoisWhitelisted.emptyValue();
            bMessageWhoisWhitelisted = res;
        }
        if (cMessageWhoisTgNameUnavailable < 1) {
            mixMissing += "message-whois_tg_name_unavailable;";
            const { res, missing } = MessageWhoisTgNameUnavailable.emptyValue();
            bMessageWhoisTgNameUnavailable = res;
        }
        if (cMessageRefuseLog < 1) {
            mixMissing += "message-refuse_log;";
            const { res, missing } = MessageRefuseLog.emptyValue();
            bMessageRefuseLog = res;
        }
        if (cMessageAcceptLog < 1) {
            mixMissing += "message-accept_log;";
            const { res, missing } = MessageAcceptLog.emptyValue();
            bMessageAcceptLog = res;
        }
        if (cMessageLiftRestrictionAlert < 1) {
            mixMissing += "message-lift_restriction_alert;";
            const { res, missing } = MessageLiftRestrictionAlert.emptyValue();
            bMessageLiftRestrictionAlert = res;
        }
        if (cMessageSilenceAlert < 1) {
            mixMissing += "message-silence_alert;";
            const { res, missing } = MessageSilenceAlert.emptyValue();
            bMessageSilenceAlert = res;
        }
        if (cMessageEnable < 1) {
            mixMissing += "message-enable;";
            const { res, missing } = MessageEnable.emptyValue();
            bMessageEnable = res;
        }
        if (cMessageDisable < 1) {
            mixMissing += "message-disable;";
            const { res, missing } = MessageDisable.emptyValue();
            bMessageDisable = res;
        }
        if (cMessageEnableLog < 1) {
            mixMissing += "message-enable_log;";
            const { res, missing } = MessageEnableLog.emptyValue();
            bMessageEnableLog = res;
        }
        if (cMessageDisableLog < 1) {
            mixMissing += "message-disable_log;";
            const { res, missing } = MessageDisableLog.emptyValue();
            bMessageDisableLog = res;
        }
        if (mixMissing.length > 0) { missing += "mix:" + mixMissing + "\n"; }
        if (mToken.length > 0) { missing += "token:" + mToken + "\n"; }
        if (mProxy.length > 0) { missing += "proxy:" + mProxy + "\n"; }
        if (mRecord.length > 0) { missing += "record:" + mRecord + "\n"; }
        if (mGroups.length > 0) { missing += "groups:" + mGroups + "\n"; }
        if (mLogChannel.length > 0) { missing += "log_channel:" + mLogChannel + "\n"; }
        if (mMainSite.length > 0) { missing += "main_site:" + mMainSite + "\n"; }
        if (mOauthAuthUrl.length > 0) { missing += "oauth_auth_url:" + mOauthAuthUrl + "\n"; }
        if (mOauthQueryUrl.length > 0) { missing += "oauth_query_url:" + mOauthQueryUrl + "\n"; }
        if (mOauthQueryKey.length > 0) { missing += "oauth_query_key:" + mOauthQueryKey + "\n"; }
        if (mWikiList.length > 0) { missing += "wiki_list:" + mWikiList + "\n"; }
        if (mBlacklist.length > 0) { missing += "blacklist:" + mBlacklist + "\n"; }
        if (mMessageStart.length > 0) { missing += "message-start:" + mMessageStart + "\n"; }
        if (mMessagePolicy.length > 0) { missing += "message-policy:" + mMessagePolicy + "\n"; }
        if (mMessageInsufficientRight.length > 0) { missing += "message-insufficient_right:" + mMessageInsufficientRight + "\n"; }
        if (mMessageGeneralPrompt.length > 0) { missing += "message-general_prompt:" + mMessageGeneralPrompt + "\n"; }
        if (mMessageTelegramIdError.length > 0) { missing += "message-telegram_id_error:" + mMessageTelegramIdError + "\n"; }
        if (mMessageRestoreSilence.length > 0) { missing += "message-restore_silence:" + mMessageRestoreSilence + "\n"; }
        if (mMessageConfirmAlready.length > 0) { missing += "message-confirm_already:" + mMessageConfirmAlready + "\n"; }
        if (mMessageConfirmOtherTg.length > 0) { missing += "message-confirm_other_tg:" + mMessageConfirmOtherTg + "\n"; }
        if (mMessageConfirmConflict.length > 0) { missing += "message-confirm_conflict:" + mMessageConfirmConflict + "\n"; }
        if (mMessageConfirmChecking.length > 0) { missing += "message-confirm_checking:" + mMessageConfirmChecking + "\n"; }
        if (mMessageConfirmUserNotFound.length > 0) { missing += "message-confirm_user_not_found:" + mMessageConfirmUserNotFound + "\n"; }
        if (mMessageConfirmButton.length > 0) { missing += "message-confirm_button:" + mMessageConfirmButton + "\n"; }
        if (mMessageConfirmWait.length > 0) { missing += "message-confirm_wait:" + mMessageConfirmWait + "\n"; }
        if (mMessageConfirmConfirming.length > 0) { missing += "message-confirm_confirming:" + mMessageConfirmConfirming + "\n"; }
        if (mMessageConfirmIneligible.length > 0) { missing += "message-confirm_ineligible:" + mMessageConfirmIneligible + "\n"; }
        if (mMessageConfirmSessionLost.length > 0) { missing += "message-confirm_session_lost:" + mMessageConfirmSessionLost + "\n"; }
        if (mMessageConfirmComplete.length > 0) { missing += "message-confirm_complete:" + mMessageConfirmComplete + "\n"; }
        if (mMessageConfirmFailed.length > 0) { missing += "message-confirm_failed:" + mMessageConfirmFailed + "\n"; }
        if (mMessageConfirmLog.length > 0) { missing += "message-confirm_log:" + mMessageConfirmLog + "\n"; }
        if (mMessageDeconfirmPrompt.length > 0) { missing += "message-deconfirm_prompt:" + mMessageDeconfirmPrompt + "\n"; }
        if (mMessageDeconfirmButton.length > 0) { missing += "message-deconfirm_button:" + mMessageDeconfirmButton + "\n"; }
        if (mMessageDeconfirmSucc.length > 0) { missing += "message-deconfirm_succ:" + mMessageDeconfirmSucc + "\n"; }
        if (mMessageDeconfirmNotConfirmed.length > 0) { missing += "message-deconfirm_not_confirmed:" + mMessageDeconfirmNotConfirmed + "\n"; }
        if (mMessageDeconfirmLog.length > 0) { missing += "message-deconfirm_log:" + mMessageDeconfirmLog + "\n"; }
        if (mMessageNewMemberHint.length > 0) { missing += "message-new_member_hint:" + mMessageNewMemberHint + "\n"; }
        if (mMessageAddWhitelistPrompt.length > 0) { missing += "message-add_whitelist_prompt:" + mMessageAddWhitelistPrompt + "\n"; }
        if (mMessageAddWhitelistSucc.length > 0) { missing += "message-add_whitelist_succ:" + mMessageAddWhitelistSucc + "\n"; }
        if (mMessageAddWhitelistLog.length > 0) { missing += "message-add_whitelist_log:" + mMessageAddWhitelistLog + "\n"; }
        if (mMessageRemoveWhitelistPrompt.length > 0) { missing += "message-remove_whitelist_prompt:" + mMessageRemoveWhitelistPrompt + "\n"; }
        if (mMessageRemoveWhitelistNotFound.length > 0) { missing += "message-remove_whitelist_not_found:" + mMessageRemoveWhitelistNotFound + "\n"; }
        if (mMessageRemoveWhitelistLog.length > 0) { missing += "message-remove_whitelist_log:" + mMessageRemoveWhitelistLog + "\n"; }
        if (mMessageRemoveWhitelistSucc.length > 0) { missing += "message-remove_whitelist_succ:" + mMessageRemoveWhitelistSucc + "\n"; }
        if (mMessageWhoisHead.length > 0) { missing += "message-whois_head:" + mMessageWhoisHead + "\n"; }
        if (mMessageWhoisPrompt.length > 0) { missing += "message-whois_prompt:" + mMessageWhoisPrompt + "\n"; }
        if (mMessageWhoisNotFound.length > 0) { missing += "message-whois_not_found:" + mMessageWhoisNotFound + "\n"; }
        if (mMessageWhoisSelf.length > 0) { missing += "message-whois_self:" + mMessageWhoisSelf + "\n"; }
        if (mMessageWhoisBot.length > 0) { missing += "message-whois_bot:" + mMessageWhoisBot + "\n"; }
        if (mMessageWhoisHasMw.length > 0) { missing += "message-whois_has_mw:" + mMessageWhoisHasMw + "\n"; }
        if (mMessageWhoisNoMw.length > 0) { missing += "message-whois_no_mw:" + mMessageWhoisNoMw + "\n"; }
        if (mMessageWhoisWhitelisted.length > 0) { missing += "message-whois_whitelisted:" + mMessageWhoisWhitelisted + "\n"; }
        if (mMessageWhoisTgNameUnavailable.length > 0) { missing += "message-whois_tg_name_unavailable:" + mMessageWhoisTgNameUnavailable + "\n"; }
        if (mMessageRefuseLog.length > 0) { missing += "message-refuse_log:" + mMessageRefuseLog + "\n"; }
        if (mMessageAcceptLog.length > 0) { missing += "message-accept_log:" + mMessageAcceptLog + "\n"; }
        if (mMessageLiftRestrictionAlert.length > 0) { missing += "message-lift_restriction_alert:" + mMessageLiftRestrictionAlert + "\n"; }
        if (mMessageSilenceAlert.length > 0) { missing += "message-silence_alert:" + mMessageSilenceAlert + "\n"; }
        if (mMessageEnable.length > 0) { missing += "message-enable:" + mMessageEnable + "\n"; }
        if (mMessageDisable.length > 0) { missing += "message-disable:" + mMessageDisable + "\n"; }
        if (mMessageEnableLog.length > 0) { missing += "message-enable_log:" + mMessageEnableLog + "\n"; }
        if (mMessageDisableLog.length > 0) { missing += "message-disable_log:" + mMessageDisableLog + "\n"; }
        return { missing: missing, mix: {
            token: bToken!,
            proxy: bProxy!,
            record: bRecord!,
            groups: bGroups!,
            logChannel: bLogChannel!,
            mainSite: bMainSite!,
            oauthAuthUrl: bOauthAuthUrl!,
            oauthQueryUrl: bOauthQueryUrl!,
            oauthQueryKey: bOauthQueryKey!,
            wikiList: bWikiList!,
            blacklist: bBlacklist!,
            messageStart: bMessageStart!,
            messagePolicy: bMessagePolicy!,
            messageInsufficientRight: bMessageInsufficientRight!,
            messageGeneralPrompt: bMessageGeneralPrompt!,
            messageTelegramIdError: bMessageTelegramIdError!,
            messageRestoreSilence: bMessageRestoreSilence!,
            messageConfirmAlready: bMessageConfirmAlready!,
            messageConfirmOtherTg: bMessageConfirmOtherTg!,
            messageConfirmConflict: bMessageConfirmConflict!,
            messageConfirmChecking: bMessageConfirmChecking!,
            messageConfirmUserNotFound: bMessageConfirmUserNotFound!,
            messageConfirmButton: bMessageConfirmButton!,
            messageConfirmWait: bMessageConfirmWait!,
            messageConfirmConfirming: bMessageConfirmConfirming!,
            messageConfirmIneligible: bMessageConfirmIneligible!,
            messageConfirmSessionLost: bMessageConfirmSessionLost!,
            messageConfirmComplete: bMessageConfirmComplete!,
            messageConfirmFailed: bMessageConfirmFailed!,
            messageConfirmLog: bMessageConfirmLog!,
            messageDeconfirmPrompt: bMessageDeconfirmPrompt!,
            messageDeconfirmButton: bMessageDeconfirmButton!,
            messageDeconfirmSucc: bMessageDeconfirmSucc!,
            messageDeconfirmNotConfirmed: bMessageDeconfirmNotConfirmed!,
            messageDeconfirmLog: bMessageDeconfirmLog!,
            messageNewMemberHint: bMessageNewMemberHint!,
            messageAddWhitelistPrompt: bMessageAddWhitelistPrompt!,
            messageAddWhitelistSucc: bMessageAddWhitelistSucc!,
            messageAddWhitelistLog: bMessageAddWhitelistLog!,
            messageRemoveWhitelistPrompt: bMessageRemoveWhitelistPrompt!,
            messageRemoveWhitelistNotFound: bMessageRemoveWhitelistNotFound!,
            messageRemoveWhitelistLog: bMessageRemoveWhitelistLog!,
            messageRemoveWhitelistSucc: bMessageRemoveWhitelistSucc!,
            messageWhoisHead: bMessageWhoisHead!,
            messageWhoisPrompt: bMessageWhoisPrompt!,
            messageWhoisNotFound: bMessageWhoisNotFound!,
            messageWhoisSelf: bMessageWhoisSelf!,
            messageWhoisBot: bMessageWhoisBot!,
            messageWhoisHasMw: bMessageWhoisHasMw!,
            messageWhoisNoMw: bMessageWhoisNoMw!,
            messageWhoisWhitelisted: bMessageWhoisWhitelisted!,
            messageWhoisTgNameUnavailable: bMessageWhoisTgNameUnavailable!,
            messageRefuseLog: bMessageRefuseLog!,
            messageAcceptLog: bMessageAcceptLog!,
            messageLiftRestrictionAlert: bMessageLiftRestrictionAlert!,
            messageSilenceAlert: bMessageSilenceAlert!,
            messageEnable: bMessageEnable!,
            messageDisable: bMessageDisable!,
            messageEnableLog: bMessageEnableLog!,
            messageDisableLog: bMessageDisableLog!,
        } }; 
    }
    static write(a: __ConfigMix): string {
        let coll = "";
        coll += Token.write(a.token) + "\n";
        coll += Proxy.write(a.proxy) + "\n";
        coll += Record.write(a.record) + "\n";
        coll += Groups.write(a.groups) + "\n";
        coll += LogChannel.write(a.logChannel) + "\n";
        coll += MainSite.write(a.mainSite) + "\n";
        coll += OauthAuthUrl.write(a.oauthAuthUrl) + "\n";
        coll += OauthQueryUrl.write(a.oauthQueryUrl) + "\n";
        coll += OauthQueryKey.write(a.oauthQueryKey) + "\n";
        coll += WikiList.write(a.wikiList) + "\n";
        coll += Blacklist.write(a.blacklist) + "\n";
        coll += MessageStart.write(a.messageStart) + "\n";
        coll += MessagePolicy.write(a.messagePolicy) + "\n";
        coll += MessageInsufficientRight.write(a.messageInsufficientRight) + "\n";
        coll += MessageGeneralPrompt.write(a.messageGeneralPrompt) + "\n";
        coll += MessageTelegramIdError.write(a.messageTelegramIdError) + "\n";
        coll += MessageRestoreSilence.write(a.messageRestoreSilence) + "\n";
        coll += MessageConfirmAlready.write(a.messageConfirmAlready) + "\n";
        coll += MessageConfirmOtherTg.write(a.messageConfirmOtherTg) + "\n";
        coll += MessageConfirmConflict.write(a.messageConfirmConflict) + "\n";
        coll += MessageConfirmChecking.write(a.messageConfirmChecking) + "\n";
        coll += MessageConfirmUserNotFound.write(a.messageConfirmUserNotFound) + "\n";
        coll += MessageConfirmButton.write(a.messageConfirmButton) + "\n";
        coll += MessageConfirmWait.write(a.messageConfirmWait) + "\n";
        coll += MessageConfirmConfirming.write(a.messageConfirmConfirming) + "\n";
        coll += MessageConfirmIneligible.write(a.messageConfirmIneligible) + "\n";
        coll += MessageConfirmSessionLost.write(a.messageConfirmSessionLost) + "\n";
        coll += MessageConfirmComplete.write(a.messageConfirmComplete) + "\n";
        coll += MessageConfirmFailed.write(a.messageConfirmFailed) + "\n";
        coll += MessageConfirmLog.write(a.messageConfirmLog) + "\n";
        coll += MessageDeconfirmPrompt.write(a.messageDeconfirmPrompt) + "\n";
        coll += MessageDeconfirmButton.write(a.messageDeconfirmButton) + "\n";
        coll += MessageDeconfirmSucc.write(a.messageDeconfirmSucc) + "\n";
        coll += MessageDeconfirmNotConfirmed.write(a.messageDeconfirmNotConfirmed) + "\n";
        coll += MessageDeconfirmLog.write(a.messageDeconfirmLog) + "\n";
        coll += MessageNewMemberHint.write(a.messageNewMemberHint) + "\n";
        coll += MessageAddWhitelistPrompt.write(a.messageAddWhitelistPrompt) + "\n";
        coll += MessageAddWhitelistSucc.write(a.messageAddWhitelistSucc) + "\n";
        coll += MessageAddWhitelistLog.write(a.messageAddWhitelistLog) + "\n";
        coll += MessageRemoveWhitelistPrompt.write(a.messageRemoveWhitelistPrompt) + "\n";
        coll += MessageRemoveWhitelistNotFound.write(a.messageRemoveWhitelistNotFound) + "\n";
        coll += MessageRemoveWhitelistLog.write(a.messageRemoveWhitelistLog) + "\n";
        coll += MessageRemoveWhitelistSucc.write(a.messageRemoveWhitelistSucc) + "\n";
        coll += MessageWhoisHead.write(a.messageWhoisHead) + "\n";
        coll += MessageWhoisPrompt.write(a.messageWhoisPrompt) + "\n";
        coll += MessageWhoisNotFound.write(a.messageWhoisNotFound) + "\n";
        coll += MessageWhoisSelf.write(a.messageWhoisSelf) + "\n";
        coll += MessageWhoisBot.write(a.messageWhoisBot) + "\n";
        coll += MessageWhoisHasMw.write(a.messageWhoisHasMw) + "\n";
        coll += MessageWhoisNoMw.write(a.messageWhoisNoMw) + "\n";
        coll += MessageWhoisWhitelisted.write(a.messageWhoisWhitelisted) + "\n";
        coll += MessageWhoisTgNameUnavailable.write(a.messageWhoisTgNameUnavailable) + "\n";
        coll += MessageRefuseLog.write(a.messageRefuseLog) + "\n";
        coll += MessageAcceptLog.write(a.messageAcceptLog) + "\n";
        coll += MessageLiftRestrictionAlert.write(a.messageLiftRestrictionAlert) + "\n";
        coll += MessageSilenceAlert.write(a.messageSilenceAlert) + "\n";
        coll += MessageEnable.write(a.messageEnable) + "\n";
        coll += MessageDisable.write(a.messageDisable) + "\n";
        coll += MessageEnableLog.write(a.messageEnableLog) + "\n";
        coll += MessageDisableLog.write(a.messageDisableLog) + "\n";
        return coll;
    }
}
