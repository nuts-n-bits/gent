
export type Line = {
        channel: string,
        listen: string | undefined,
        data: string,
} | {
        close: string,
}

const CHANID_MAX = 64;
const CHANID_CHARSET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_";
const INLINE_WHITESPACE = " \t";

export function readline(line: string): Line | Error {
        let i = 0;
        let channel = "";
        let listen = "";
        while (true) {  // get header
                const cur = line[i];
                if (cur === undefined) {
                        return new Error("header incomplete");
                } else if (INLINE_WHITESPACE.includes(cur)) {
                        if (channel !== "") {
                                return new Error("encountered whitespace after having encountered channel id")
                        }
                        i += 1;  // skip whitespace as long as header is not done
                } else if (CHANID_CHARSET.includes(cur)) {
                        const [chanid, new_i] = consume_chanid(line, i);
                        channel = chanid;
                        i = new_i;
                } else if (cur === "/") {
                        const next = line[i+1];
                        if (!next || !CHANID_CHARSET.includes(next)) {
                                return new Error("header incomplete");
                        }
                        const [listenid, new_i] = consume_chanid(line, i+1);
                        listen = listenid;
                        i = new_i;
                        if (line[i] !== " " && line[i] !== undefined) {
                                return new Error("header must be separated with body by a space character");
                        }
                        let body = line.substring(i+1);
                        if (channel === "") {  // close-only line
                                if (body !== "") {
                                        return new Error("close-only line must contain no body message");
                                }
                                return { close: listen };
                        } else {
                                return {
                                        channel, 
                                        listen: listen === "0" ? undefined : listen,
                                        data: body,
                                }
                        }
                } else {
                        return new Error("unexpected character");
                }
        }
}

function consume_chanid(s: string, i: number): [string, number] {
        let coll = "";
        while (s[i] && CHANID_CHARSET.includes(s[i]!)) {
                coll += s[i];
                i += 1;
        }
        return [coll, i];
}