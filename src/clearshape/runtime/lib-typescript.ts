// BEGIN RUNTIME LIBRARY

type $J = 
    string | number | boolean | null | { [i: string]: $J } | $J[]

const $$i64Max = (2n << 63n) - 1n
const $$i64Min = -1n * (2n << 63n)
const $$u64Max = (2n << 64n) - 1n
const $$u64Min = 0n

function $never(_?: never): never { throw new Error("unreachable") }

function $parseI64(a: $J): bigint | Error {	
    try {
        if (typeof a !== "string") { return new Error("expecting string"); }
        if (a.length > 99) { return new Error("i64 bigint decimal length overflow"); }
        const i64 = BigInt(a);
        if (i64 > $$i64Max) { return new Error("i64 range overflow"); }
        if (i64 < $$i64Min) { return new Error("i64 range underflow"); }
        return i64;
    } catch (e) {
        if (e instanceof Error) { return e; }
        else { return new Error("caught non error", { cause: e }); }
    }
}

function $parseU64(a: $J): bigint | Error {	
    try {
        if (typeof a !== "string") { return new Error("expecting string"); }
        if (a.length > 99) { return new Error("u64 bigint decimal length overflow"); }
        const u64 = BigInt(a);
        if (u64 > $$u64Max) { return new Error("u64 range overflow"); }
        if (u64 < $$u64Min) { return new Error("u64 range underflow"); }
        return u64;
    } catch (e) {
        if (e instanceof Error) { return e; }
        else { return new Error("caught non error", { cause: e }); }
    }
}

function $parseF64(a: $J): number | Error {	
    try {
        if (typeof a !== "string") { return new Error("expecting string"); }
        if (a.length > 99) { return new Error("f64 decimal length overflow"); }
        const num = Number(a);
        if (isNaN(num)) { return new Error("f64 is NaN"); }
        return num;
    } catch (e) {
        if (e instanceof Error) { return e; }
        else { return new Error("caught non error", { cause: e }); }
    }
}

function $parseBoolean(a: $J): boolean|Error {
    return typeof a === "boolean" ? a : new Error("expected boolean");
}

function $parseString(a: $J): string|Error {
    return typeof a === "string" ? a : new Error("expected string");
}

function $parseNull(a: $J): null|Error {
    return a === null ? a : new Error("expected null");
}

function $parseBinary(a: $J): Uint8Array|Error {
    if (typeof a !== "string") { return new Error("expecting string"); }
    if (a.length % 4 !== 0) { return new Error("Invalid input") }
    return $fromb64(a)
}

const $b64dict = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
const $b64idict = Object.fromEntries($b64dict.split("").map((char, idx) => [char, idx]));
$b64idict["="] = 0;

function $tob64(buf: Uint8Array): string {
	let str = "";
	for (let i=0; i<buf.length; i+=3) {
		const word24 = ((buf[i] ?? 0) << 16) | ((buf[i+1] ?? 0) << 8) | (buf[i+2] ?? 0)
		const b64char4 = $b64dict[(word24      ) & 0b00_11_11_11];
		const b64char3 = $b64dict[(word24 >> 6 ) & 0b00_11_11_11];
		const b64char2 = $b64dict[(word24 >> 12) & 0b00_11_11_11];
		const b64char1 = $b64dict[(word24 >> 18) & 0b00_11_11_11];
		str += (b64char1! + b64char2! + b64char3! + b64char4!);
	}
	if (buf.length % 3 === 0) { return str; }
	if (buf.length % 3 === 1) { return str.slice(0, str.length-2) + "=="; }
	if (buf.length % 3 === 2) { return str.slice(0, str.length-1) + "="; }
	throw new Error("lol impossible");
}

function $fromb64(b64: string): Uint8Array|Error {
	if (b64.length % 4 === 1) { return new Error("Invalid input"); }
	while(b64.length % 4 !== 0) { b64 += "="; }
	if (b64.length === 0) { return new Uint8Array(0); }
	const missing = b64[b64.length-2] === "=" ? 2 : b64[b64.length-1] === "=" ? 1 : 0;
    if (b64[b64.length-2] === "=" && b64[b64.length-1] !== "=") { return new Error("malformed base64 input"); }
	const buf = new Uint8Array(b64.length / 4 * 3 - missing);
	for (let i=0; i<b64.length; i+=4) {
		const word24 = ($b64idict[b64[i]!]! << 18) | ($b64idict[b64[i+1]!]! << 12) | ($b64idict[b64[i+2]!]! << 6) | $b64idict[b64[i+3]!]!;
		const byte3 = (word24       ) & 0xFF;
		const byte2 = (word24 >> 8  ) & 0xFF;
		const byte1 = (word24 >> 16 ) & 0xFF;
		if (i !== b64.length - 4 || missing === 0) {
			buf.set([byte1, byte2, byte3], i/4*3);
		}
		else {
			if (missing === 1) { buf.set([byte1, byte2], i/4*3); }
			else if (missing === 2) { buf.set([byte1], i/4*3); }
			else { throw new Error("lol impossible"); }
		}
	}
	return buf;
}

function $writeI64(a: bigint): $J { return a.toString(); }
function $writeU64(a: bigint): $J { return a.toString(); }
function $writeF64(a: number): $J { return a.toString(); }
function $writeBoolean(a: boolean): $J { return a; }
function $writeNull(a: null): $J { return null; }
function $writeBinray(a: Uint8Array): $J { return $tob64(a) }
function $writeString(a: string): $J { return a; }
