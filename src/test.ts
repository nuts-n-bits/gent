import { KisseCore, DataFrame } from "./core-serde.js";

function frameOf(id: bigint, content: string): DataFrame {
        return {
                frame_id: id,
                payload: new TextEncoder().encode(content),
        }
}

function fromFrame(f: DataFrame): [bigint, string] {
        return [f.frame_id, new TextDecoder().decode(f.payload)];
}

const encoded = [
        ...KisseCore.encode(frameOf(218n, "hello")),
        ...KisseCore.encode(frameOf(58375837n, "gtejuisog5rhe4oiwghertisughirtuskghursghjiurt")),
        ...KisseCore.encode(frameOf(567534892048354348230437546478349234375643n, "")),
];

console.log("encoded:", encoded);

const decoded = [...KisseCore.decode(encoded, { FD_BYTE: 100 })];

console.log(decoded);

console.log(decoded.map(fromFrame));