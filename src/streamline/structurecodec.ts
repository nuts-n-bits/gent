import { StreamLine } from "./linecodex";


type value<colltype> = {
        args: [],
        options: [string, string][],
        subs: undefined | colltype,
}

type Valuecontainer = value<Valuecontainer>[]

function new_coll(): Valuecontainer {
        return [];
}

function pusher(v: Valuecontainer, b: value<Valuecontainer>) {
        v.push(b);
}

function stopper(v: Valuecontainer) {
        return;
}

struct_parser([], new_coll, pusher, stopper);


async function* struct_parser<colltype>(
        line_stream: Iterable<StreamLine> | AsyncIterable<StreamLine>, 
        new_coll: () => colltype, 
        pusher: (a: colltype, b: value<colltype>) => void, 
        stopper: (a: colltype) => void
): AsyncIterable<value<colltype>> {
        const open_channels: Map<string, colltype> = new Map();
        for await (const line of line_stream) {
                let coll: colltype | undefined;
                if (line.subchannels.length > 0) { 
                        coll = new_coll();
                        for (const ch of line.channels) {
                                open_channels.set()
                        }
                } else { 
                        coll = undefined; 
                }
                const value: value<colltype> = {
                        args: [], 
                        options: [],
                        subs: coll,
                }
        }
}