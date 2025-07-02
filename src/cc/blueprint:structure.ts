type StructDef = { [tag: string]: StructDef };
function struct_def_check(def: StructDef, depth = 0, type_tags_by_depth: Set<string>[] = []): "OK" | Error {
        if (type_tags_by_depth[depth] === undefined) { type_tags_by_depth[depth] = new Set(); }
        for (const [tag, _] of Object.entries(def)) { 
                if (is_tag_exists_at_depth_less_than_n(type_tags_by_depth, tag, depth)) { return new Error(`Type tag "${tag}" appears at a parent level. The same type tag can only appear at the same level to avoid creating ambiguous sequences. Consider aliasing the command name if you need the same command at different levels.`); }
                type_tags_by_depth[depth]!.add(tag);
        }
        for (const [_, children] of Object.entries(def)) {
                const res = struct_def_check(children, depth+1, type_tags_by_depth);
                if (res instanceof Error) { return res; }
        }
        return "OK";
}
function is_tag_exists_at_depth_less_than_n(type_tags_by_depth: Set<string>[], tag: string, depth_n: number): boolean {
        for (let i=0; i<depth_n; i++) {
                if (type_tags_by_depth[i]?.has(tag)) { return true; }
        }
        return false;
}

type ParsedStruct<T> = { value: T, children: ParsedStruct<T>[] }
function struct_parser<T>(def: StructDef, sequence: T[], tagof: (t: T) => string, i = 0): [ParsedStruct<T>[], number] {
        const coll: ParsedStruct<T>[] = [];
        while (true) {
                const cur = sequence[i];
                if (cur === undefined) { break; }  // reached the end of the sequence -> return collected parsed structs
                const tag = tagof(cur);
                const curdef = def[tag];
                if (curdef === undefined) { break; }  // not in our definition layer -> return to parent caller
                const value = cur;
                const [children, new_i] = struct_parser(curdef, sequence, tagof, i+1);
                i = new_i;
                coll.push({ value, children });
        }
        return [coll, i];
}

const test: StructDef = {
        "arinit": {},
        "arsync": {
                "addrev": {
                        "rrdrecord": {},
                },
                "addpc": {},
                "addpid": {},
        },
        "arclose": {},
}

const testseq = [
        { command: "arinit"   , args: [], options: [], },
        { command: "arsync"   , args: [], options: [], },
        { command: "addrev"   , args: [], options: [], },
        { command: "addpc"    , args: [], options: [], },
        { command: "addrev"   , args: [], options: [], },
        { command: "rrdrecord", args: [], options: [], },
        { command: "addrev"   , args: [], options: [], },
        { command: "arclose"  , args: [], options: [], },
]

console.log(struct_parser(test, testseq, i => i.command)[0][1]?.children);
/*

(arsync 
        (session XnDf93Glw94W1)
        (addrev (r 1488291)(t 2025-03-24T12:55:59Z)(u 100294)(un Hinata)(s "created a new page")(c "....."))
        (addpid 1973742)
        (addpc Wikipedia:Test WP:Test)
)

const session = cc.type(
        cc.ident("session"),
        cc.value.unique(cc.string()),
        cc.field.unique(cc.type(""))
)

type FieldType = { 
        list: number,
        type: "string" | "decimal" | "udecimal" | "boolean" | "base64",
        kind: "repeated" | "required" | "optional" | "none",
}

(arinit (zhwp) ((user Hinata) (uid 198473)) (
        (arsync (zhwp) ((session FEgnwiceifjmaGRGrwg))
                (addrev () ())
        )
))

type SeType = {
        ident: string,
        posval: FieldType,

}

(arsync 
        (session )
)

cc.type(
        cc.ident("arsync"),
        cc.field.unique(cc.type(
                cc.ident("session"),
                cc.value.string.unique(),
        )),
        cc.field.repeated(cc.type(
                cc.ident("addrev"),
                cc.field.
        ))
)


arsync [session XnDfefGEH] {
        addpid 947925
        addpc "大大小小的事情"
        addrev [r 37782; t 2024-10-15T16:23:07Z; un Hinata; u 117137; s ""; c "......"]
}

special ones: "->" "<-" "-X" "[" "]"

arsync [session XnDfefGEH] <- 0
-> 0 addpc 398583
-> 0 addrev [r 325325; t 2025-03-11T06:19:24Z; s ""; un Hinata; u 117118; c "......"] <- 1
-> 1 rrd RRD1 RRD2 RDTEXT RDSUMMARY [un Hinata; s "理由：RRD#2，侵犯版权"]
-> 1 revalttext "......"
-| 1
-| 0


- arsync -1 [session XnDfZGhGt6]
1 addpc WP:Test
1 addrev [r 4366543; t 2025-03-25T06:28:44Z; un hinata; u 87769; s ""; c "....."] -2
1 addrev [......]
1 addpid 4436376

channel


arsync -session XndFefFGh <- 0

arsync <- 0
-> 0 session XnDfeE3H2m2
-> 0 addpid 498812
-> 0 addpc "大大小小的事情"
-> 

(arsync (session (<- 3)))
(-> 3 XnDf3D3)

(arsync (session XnDf3D3))


*/






