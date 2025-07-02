const MAX_CHANID = (2n ** 64n) - 1n;
const MIN_CHANID = 0n;
class InvalidError extends Error {}
class InvalidJsReprError extends InvalidError {}
class InvalidWireError extends InvalidError {}
class InternalError extends InvalidError {}

type HeaderA = {
        target_channels: bigint[], 
        listen_channel: bigint | undefined,
}

type CloseFrame = {
        close_channel: bigint,
}

function encode_header(h: HeaderA): Uint8Array {
        if (h.target_channels.length === 0) { throw new InvalidJsReprError("a message must target at least one channel"); }
        if (h.listen_channel === 0n) { throw new InvalidJsReprError("cannot subscribe to channel 0"); }
        if (h.target_channels.length === 1 && h.listen_channel === undefined) { // compact form
                return encode_header_core(h.target_channels);
        } else {
                return encode_header_core([...h.target_channels, h.listen_channel ?? 0n]);
        }

}

function decode_header(message: Uint8Array, maxbyte: number, maxchans: number): { header: HeaderA, offset: number } | "EEOF" | "ETOOLONG" | "ETOOMANY" {
        const res = decode_header_core(message, maxbyte, maxchans);
        if (typeof res === "string") { return res; }
        if (res.chans.length === 0) {
                throw new InvalidWireError("assert unreachable");
        } else if (res.chans.length === 1) { 
                return {
                        header: {
                                target_channels: res.chans,
                                listen_channel: undefined,
                        },
                        offset: res.offset,
                }
        } else {
                const listen = res.chans.pop()!;
                return {
                        header: {
                                target_channels: res.chans,
                                listen_channel: listen === 0n ? undefined : listen,
                        },
                        offset: res.offset,
                }
        }
}

function encode_header_core(bigints: bigint[]): Uint8Array {
        const coll = [] as number[];
        for (let bigint of bigints) {
                if (bigint > MAX_CHANID || bigint < MIN_CHANID) { 
                        throw new Error("channel id must be inside range u64"); 
                }
                const subcoll = [] as number[];
                do {
                        subcoll.push( Number( 0b0011_1111n & bigint ) );
                        bigint = bigint >> 6n;
                } while (bigint > 0n)
                subcoll[subcoll.length-1] = subcoll[subcoll.length-1]! | 0b0100_0000;
                coll.push(...subcoll);
        }
        coll[coll.length-1] = coll[coll.length-1]! | 0b1000_0000;
        return new Uint8Array(coll);
}

function decode_header_core(message: Uint8Array, maxbyte: number, maxchans: number): { chans: bigint[], offset: number } | "EEOF" | "ETOOLONG" | "ETOOMANY" {
        const coll = [] as bigint[];
        let i = 0;
        let temp = 0n;
        let tempi = 0n;
        while (true) {
                if (i >= maxbyte) { return "ETOOLONG"; }
                if (coll.length >= maxchans) { return "ETOOMANY"; }
                const cur = message[i];
                i += 1;
                if (cur === undefined) {
                        return "EEOF";
                }
                const chanid_ends = (cur & 0b0100_0000) !== 0;
                const header_ends = (cur & 0b1000_0000) !== 0;
                const value = BigInt(cur & 0b0011_1111);
                temp = temp | (value << 6n * tempi);
                tempi += 1n;
                if (chanid_ends) {
                        coll.push(temp);
                        temp = 0n;
                        tempi = 0n;
                }
                if (header_ends && !chanid_ends) {
                        throw new InvalidWireError("invalid header bit pattern");
                }
                if (header_ends) {
                        break;
                }
        }
        return { chans: coll, offset: i };
}


const enc = encode_header({
        target_channels: [655372876n, 8888888888888n, 999999999999n, 456765434567n, 65434567654345454n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n, 6n],
        listen_channel: 2827868268682n,
});

console.log(enc);

console.log(decode_header(enc, 100, 27));
