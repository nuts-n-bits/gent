/*
The Kiss Encoding: 
	DataFrame(frame_id, body_buffer) => VarInt(frame_id) ^ VarInt(l) ^ [ body buffer ] 
	where body_buffer.length == l
This encoding will be the core serde for the gent data format.
This encoding requires minimum 3 bytes to represent any meaningful data, for example:
	Example data { status: 200 } will probably become b"\x01\x01\xC8"
	Example data { status: 300 } will probably become b"\x01\x02\x01\x2C"
where 
	first byte 0x01 is varint for frame_id, assigned for "status" field in some proto file
	second byte 0x01 or 0x02 is the length of the following body encoded in varint
	the rest is body, in this case simply encoded in binary form (200 = \xC8, 300 = \x01\x2C)
It's a super simple encoding so it is easy to implement correctly.
*/

/** 
 * @member FD_BYTE: How many bytes the decoder will tolerate the FRAME_DESCRIPTOR (i.e. VarInt(frame_id)) before giving 
 * 		up, defaults to 8, which allows any frameId between 0 and 2^56-1
 * @member LD_BYTE: How many bytes the decoder will tolerate the LENGTH_DESCRIPTOR (i.e. the varint that encodes payload 
 * 		length) before giving up, defaults to 8, which allows (0 ~ 2^56-1) to be represented
 * @member BODY_BYTE: How many bytes the decoder will allow the body to build up, i.e. actual length limit for each 
 * 		field. Defaults to 1<<30, which means each field can be 1GiB.
 * This value is a js number so precision is up to 2^53-1 due to IEEE754. *The body in memory will only build up to this 
 * 		length, then decoder will stop working, even if LD_BYTE allows a larger LD to be parsed.*
 */
export type SecT = { FD_BYTE?: number, LD_BYTE?: number, BODY_BYTE?: number };
const SEC_FALLBACK_FD_BYTE = 8;
const SEC_FALLBACK_LD_BYTE = 8;
const SEC_FALLBACK_BODY_BYTE = 1<<30;

export class DecodeWorktable {
	public readonly SEC_FD_BYTE: number;
	public readonly SEC_LD_BYTE: number;
	public readonly SEC_BODY_BYTE: number;
	constructor(sec = {} as SecT) {
		this.SEC_FD_BYTE = sec.FD_BYTE ?? SEC_FALLBACK_FD_BYTE;
		this.SEC_LD_BYTE = sec.LD_BYTE ?? SEC_FALLBACK_LD_BYTE;
		this.SEC_BODY_BYTE = sec.BODY_BYTE ?? SEC_FALLBACK_BODY_BYTE;
	}
	private inputQ: Uint8Array[] = [];
	private inputSumLen = 0;
	private inputPeek(i: number): number | undefined {
		if (this.inputSumLen <= i) { return undefined; }
		let current_chunk_index = 0;
		let input_partial_sum_len = this.inputQ[0]!.length;
		while (i > input_partial_sum_len) {
			current_chunk_index += 1;
			if (this.inputQ[current_chunk_index] === undefined) { return undefined; }
			input_partial_sum_len += this.inputQ[current_chunk_index]!.length;
		}
		const i2 = i - input_partial_sum_len + this.inputQ[current_chunk_index]!.length;
		return this.inputQ[current_chunk_index]![i2];
	}
	public load(buf: Uint8Array) { 
		if (buf.length === 0) { return; }
		this.inputQ.push(buf);
		this.inputSumLen += buf.length;
	}
	public drop(): Uint8Array[] {
		const droppedQ = this.inputQ;
		this.inputQ = [];
		this.inputSumLen = 0;
		return droppedQ;
	}
	public len(): number {
		return this.inputSumLen;
	}
	/** @returns
	 * - If returns DataFrame, this is one complete DataFrame and the corresponding data has been removed from internal 
	 * 		state. The remaining input might contain more DataFrames and you are encouraged to call step() again to get 
	 * 		these frames.
	 * - If returns number, the data has not been touched and the only reason would be that the current input data does 
	 * 		not yet form a complete frame. If the number is -1, then the header is incomplete. If the number is > 0, 
	 * 		then that means how many bytes are missing from the body. If you call step() again without loading more 
	 * 		input, it will return the same number again and there is no point to do this. Load more chunks first, then 
	 * 		call again.
	 * - If returns Error, then one of the sec Errors have been triggered. Either drop() and put data into another 
	 * 		worktable with a bigger sec threshold, or maybe data is corrupt. Either way, data is preserved and not 
	 * 		touched since the last time step() returned a DataFrame. */
	public step(): { frameid: bigint, data: Uint8Array[] } | number | Error {
		if (this.inputSumLen === 0) { return -1; }
		const fdvi_coll = [] as number[]; // frame descriptor varint collection
		for (let i=0; true; i++) {
			if (i >= this.SEC_FD_BYTE) { return new Error("frame descriptor overlength limit " + this.SEC_FD_BYTE); }
			const varint_byte = this.inputPeek(i);
			// step() does nothing since the inputQ does not have enough data to form a varint header
			if (varint_byte === undefined) { return -1; }  
			fdvi_coll.push(varint_byte);
			if ((varint_byte & 0x80) === 0) { break; }  // continue out of the loop
		}
		const ldvi_coll = [] as number[]; // length descriptor varint collection
		for (let i=0; true; i++) {
			if (i >= this.SEC_LD_BYTE) { return new Error("length descriptor overlength limit " + this.SEC_LD_BYTE); }
			const varint_byte = this.inputPeek(i + fdvi_coll.length);
			// step() does nothing since the inputQ does not have enough data to form a varint header
			if (varint_byte === undefined) { return -1; }  
			ldvi_coll.push(varint_byte);
			if ((varint_byte & 0x80) === 0) { break; }  // continue out of the loop
		}
		const frame_id = decode_varint(fdvi_coll);
		const body_len = Number(decode_varint(ldvi_coll));
		if (body_len > this.SEC_BODY_BYTE) { return new Error("body overlength limit " + this.SEC_BODY_BYTE); }
		if (this.inputSumLen < (body_len + fdvi_coll.length + ldvi_coll.length)) { 
			// step() does nothing again, inputQ not enough data to form a DataFrame
			return body_len + fdvi_coll.length + ldvi_coll.length - this.inputSumLen; 
		}  
		/* Discard return value = */ this.unload(fdvi_coll.length + ldvi_coll.length);
		const payload = this.unload(body_len);
		return { frameid: frame_id, data: payload };
	}
	// this is similar to step(), except it expects its body to be LD(bytes) ++ Ld(bytes) ++ ...
	// it should only be used the contents of a repeated frame that is using compact representation.
	public stepCompact(): Uint8Array[] | number | Error {
		if (this.inputSumLen === 0) { return -1; }
		const ldvi_coll = [] as number[]; // length descriptor varint collection
		for (let i=0; true; i++) {
			if (i >= this.SEC_LD_BYTE) { return new Error("length descriptor overlength limit " + this.SEC_LD_BYTE); }
			const varint_byte = this.inputPeek(i);
			// step() does nothing since the inputQ does not have enough data to form a varint header
			if (varint_byte === undefined) { return -1; }  
			ldvi_coll.push(varint_byte);
			if ((varint_byte & 0x80) === 0) { break; }  // continue out of the loop
		}
		const body_len = Number(decode_varint(ldvi_coll));
		if (body_len > this.SEC_BODY_BYTE) { return new Error("body overlength limit " + this.SEC_BODY_BYTE); }
		if (this.inputSumLen < (body_len + ldvi_coll.length)) { 
			// step() does nothing again, inputQ not enough data to form a DataFrame
			return body_len + ldvi_coll.length - this.inputSumLen; 
		}  
		/* Discard return value = */ this.unload(ldvi_coll.length);
		const payload = this.unload(body_len);
		return payload;
	}
	// assume have enough!
	private unload(bytes: number): Uint8Array[] {
		if (bytes === 0) { return [new Uint8Array(0)]; }
		const coll: Uint8Array[] = [];
		let deficit = bytes;
		while (deficit > 0) {
			const first_buf_len = this.inputQ[0]!.length;
			if (deficit >= first_buf_len) {
				deficit -= first_buf_len;
				coll.push(this.inputQ.shift()!); 
			} else {
				coll.push(this.inputQ[0]!.subarray(0, deficit));
				this.inputQ[0] = this.inputQ[0]!.subarray(deficit);
				deficit = 0;
			}
		}
		this.inputSumLen -= bytes;
		return coll;
	}
}

function encode_varint(varint: bigint): number[] {
	if (varint < 0n) { throw new Error("negative varint is not accepted"); }
	if (varint === 0n) { return [0]; }
	const buf: number[] = [];
	while (varint != 0n) {
		buf.push(Number((varint & 0x7Fn) | 0x80n))
		varint = varint >> 7n
	}
	buf[buf.length-1]! &= 0x7F;
	return buf;
}

/** @param buf assumes bounds is correct, will not check continuation bit. assume at least 1 byte. */
function decode_varint(buf: Iterable<number>): bigint {
	let count = 0, ret = 0n;
	for (const byte of buf) {
		ret |= (BigInt(byte & 0x7F) << BigInt(count*7));
		count += 1;
	}
	return ret;
}

export function total_len(ui8ai: Iterable<Uint8Array>): number {
        let sum_length = 0;
        for (const ui8a of ui8ai) { sum_length += ui8a.length; }
        return sum_length;
}

export function header_of(frameid: bigint, data: Uint8Array): Uint8Array { 
	const frame_header = encode_varint(frameid).concat(encode_varint(BigInt(data.length)));
	return new Uint8Array(frame_header);
}

export function header_of_frag(frameid: bigint, data: Uint8Array[]): Uint8Array { 
	const frame_header = encode_varint(frameid).concat(encode_varint(BigInt(total_len(data))));
	return new Uint8Array(frame_header);
}

export function delimit_compact_no_rem(buffers: Uint8Array[]): Uint8Array[][] | Error {
	const dwt = new DecodeWorktable();
	for (const buf of buffers) {
		dwt.load(buf)
	}
	const coll = [] as Uint8Array[][];
	while (true) {
		if (dwt.len() == 0) {
			break;
		}
		const stepResult = dwt.stepCompact();
		if (stepResult instanceof Error) {
			return stepResult;
		}
		if (typeof stepResult === "number") {
			return new Error("incomplete message detected, signal = " + stepResult);
		}
		coll.push(stepResult);
	}
	return coll;
}