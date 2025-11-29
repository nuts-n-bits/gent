export type __C = { [_: string]: string }

export type __D = string[]

export type __A = {
    session: string,
    name: {
        s: string,
    }[],
    name2: string[],
    time2?: undefined | string,
    time4?: undefined | {
        y: bigint,
        m: bigint,
        d: bigint,
    },
    technicalIdentifier: { [_: string]: bigint },
    another?: undefined | __B,
    str: {
        struct: string,
        struct2: string,
    },
    enu: {
        struct: string
    } | {
        struct2: string
    },
}

export type __B = {
    s: __A,
    v: [
        string,
        string,
    ],
    c: [
        bigint,
        string,
    ],
    d: [string],
    map: { [_: string]: __A },
    bin: { [_: string]: Uint8Array[] },
}

export class A {
    static fromJson(a: string): __A | Error {
        try { 
            const obj = JSON.parse(a);
            return this.fromJsonCore(obj);
        } catch(e) {
            if (!(e instanceof Error)) { return new Error("caught non error"); }
            return e;
        }
        
    }
    
    static fromJsonCore(a: $J): __A | Error {
        const parser = (a: $J) => {
            if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing struct"); }
            const copycat = { ...a };
            // for each field: create parsers
            const parserSession = $parseString;
            const parserName = (a: $J) => {
                if (!(a instanceof Array)) { return new Error("expected array while parsing list"); }
                const coll = [] as {
                    s: string,
                }[];
                const parser = (a: $J) => {
                    if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing struct"); }
                    const copycat = { ...a };
                    // for each field: create parsers
                    const parserS = $parseString;
                    // for required fields only: check presence
                    if (copycat["s"] === undefined) { return new Error("required field 's' (wire name '" + "s" + "') is undefined") }
                    // for each field: parse, respecting requiredness, early return on error
                    const parsedS = parserS(copycat["s"]);
                    if (parsedS instanceof Error) { return new Error("error when parsing field s (wire name '" + "s" + "')", { cause: parsedS }); }
                    // for each field: delete field from copycat object
                    delete copycat["s"];
                    if (Object.keys(copycat).length > 0) { return new Error("unknown fields present: " + Object.keys(copycat).join(", ")); }
                    return {
                        s: parsedS, 
                    }
                };
                for (const elem of a) {
                    const parsed = parser(elem);
                    if (parsed instanceof Error) { return new Error("failed to parse list inner type", { cause: parsed }); } 
                    coll.push(parsed);
                }
                return coll;
            };
            const parserName2 = (a: $J) => {
                if (!(a instanceof Array)) { return new Error("expected array while parsing list"); }
                const coll = [] as string[];
                const parser = $parseString;
                for (const elem of a) {
                    const parsed = parser(elem);
                    if (parsed instanceof Error) { return new Error("failed to parse list inner type", { cause: parsed }); } 
                    coll.push(parsed);
                }
                return coll;
            };
            const parserTime2 = $parseString;
            const parserTime4 = (a: $J) => {
                if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing struct"); }
                const copycat = { ...a };
                // for each field: create parsers
                const parserY = $parseI64;
                const parserM = $parseI64;
                const parserD = $parseI64;
                // for required fields only: check presence
                if (copycat["y"] === undefined) { return new Error("required field 'y' (wire name '" + "y" + "') is undefined") }
                if (copycat["m"] === undefined) { return new Error("required field 'm' (wire name '" + "m" + "') is undefined") }
                if (copycat["d"] === undefined) { return new Error("required field 'd' (wire name '" + "d" + "') is undefined") }
                // for each field: parse, respecting requiredness, early return on error
                const parsedY = parserY(copycat["y"]);
                if (parsedY instanceof Error) { return new Error("error when parsing field y (wire name '" + "y" + "')", { cause: parsedY }); }
                const parsedM = parserM(copycat["m"]);
                if (parsedM instanceof Error) { return new Error("error when parsing field m (wire name '" + "m" + "')", { cause: parsedM }); }
                const parsedD = parserD(copycat["d"]);
                if (parsedD instanceof Error) { return new Error("error when parsing field d (wire name '" + "d" + "')", { cause: parsedD }); }
                // for each field: delete field from copycat object
                delete copycat["y"];
                delete copycat["m"];
                delete copycat["d"];
                if (Object.keys(copycat).length > 0) { return new Error("unknown fields present: " + Object.keys(copycat).join(", ")); }
                return {
                    y: parsedY, 
                    m: parsedM, 
                    d: parsedD, 
                }
            };
            const parserTechnicalIdentifier = (a: $J) => {
                if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing map"); }
                const coll = {} as { [i: string]: bigint };
                const parser = $parseI64;
                for (const k in a) {
                    const parsed = parser(a[k]!);
                    if (parsed instanceof Error) { return new Error("failed to parse map's inner type", { cause: parsed }); } 
                    coll[k] = parsed;
                }
                return coll;
            };
            const parserAnother = (a: $J) => B.fromJsonCore(a);
            const parserStr = (a: $J) => {
                if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing struct"); }
                const copycat = { ...a };
                // for each field: create parsers
                const parserStruct = $parseString;
                const parserStruct2 = $parseString;
                // for required fields only: check presence
                if (copycat["struct"] === undefined) { return new Error("required field 'struct' (wire name '" + "struct" + "') is undefined") }
                if (copycat["struct2"] === undefined) { return new Error("required field 'struct2' (wire name '" + "struct2" + "') is undefined") }
                // for each field: parse, respecting requiredness, early return on error
                const parsedStruct = parserStruct(copycat["struct"]);
                if (parsedStruct instanceof Error) { return new Error("error when parsing field struct (wire name '" + "struct" + "')", { cause: parsedStruct }); }
                const parsedStruct2 = parserStruct2(copycat["struct2"]);
                if (parsedStruct2 instanceof Error) { return new Error("error when parsing field struct2 (wire name '" + "struct2" + "')", { cause: parsedStruct2 }); }
                // for each field: delete field from copycat object
                delete copycat["struct"];
                delete copycat["struct2"];
                if (Object.keys(copycat).length > 0) { return new Error("unknown fields present: " + Object.keys(copycat).join(", ")); }
                return {
                    struct: parsedStruct, 
                    struct2: parsedStruct2, 
                }
            };
            const parserEnu = (a: $J) => {
                type retType = {
                    struct: string
                } | {
                    struct2: string
                };
                if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing enum"); }
                const entries = Object.entries(a);
                if (entries.length !== 1) { return new Error("enum values must contain exactly 1 field"); } 
                const [k, v] = entries[0]!;
                switch (k) {
                case "struct": 
                    const parserStruct = $parseString;
                    const parsedStruct = parserStruct(v);
                    return { struct: parsedStruct } as retType; 
                break;
                case "struct2": 
                    const parserStruct2 = $parseString;
                    const parsedStruct2 = parserStruct2(v);
                    return { struct2: parsedStruct2 } as retType; 
                break;
                default: 
                    return new Error("unknown variant name while parsing enum, expected one of " + [
                        "struct ('" + "struct" + "')",
                        "struct2 ('" + "struct2" + "')",
                    ].join(" / "));
                }
            };
            // for required fields only: check presence
            if (copycat["s"] === undefined) { return new Error("required field 'session' (wire name '" + "s" + "') is undefined") }
            if (copycat["n"] === undefined) { return new Error("required field 'name' (wire name '" + "n" + "') is undefined") }
            if (copycat["n2"] === undefined) { return new Error("required field 'name2' (wire name '" + "n2" + "') is undefined") }
            if (copycat["ti"] === undefined) { return new Error("required field 'technicalIdentifier' (wire name '" + "ti" + "') is undefined") }
            if (copycat["str"] === undefined) { return new Error("required field 'str' (wire name '" + "str" + "') is undefined") }
            if (copycat["enu"] === undefined) { return new Error("required field 'enu' (wire name '" + "enu" + "') is undefined") }
            // for each field: parse, respecting requiredness, early return on error
            const parsedSession = parserSession(copycat["s"]);
            if (parsedSession instanceof Error) { return new Error("error when parsing field session (wire name '" + "s" + "')", { cause: parsedSession }); }
            const parsedName = parserName(copycat["n"]);
            if (parsedName instanceof Error) { return new Error("error when parsing field name (wire name '" + "n" + "')", { cause: parsedName }); }
            const parsedName2 = parserName2(copycat["n2"]);
            if (parsedName2 instanceof Error) { return new Error("error when parsing field name2 (wire name '" + "n2" + "')", { cause: parsedName2 }); }
            const parsedTime2 = copycat["t2"] === undefined ? undefined : parserTime2(copycat["t2"]);
            if (parsedTime2 instanceof Error) { return new Error("error when parsing field time2 (wire name '" + "t2" + "')", { cause: parsedTime2 }); }
            const parsedTime4 = copycat["t4"] === undefined ? undefined : parserTime4(copycat["t4"]);
            if (parsedTime4 instanceof Error) { return new Error("error when parsing field time4 (wire name '" + "t4" + "')", { cause: parsedTime4 }); }
            const parsedTechnicalIdentifier = parserTechnicalIdentifier(copycat["ti"]);
            if (parsedTechnicalIdentifier instanceof Error) { return new Error("error when parsing field technicalIdentifier (wire name '" + "ti" + "')", { cause: parsedTechnicalIdentifier }); }
            const parsedAnother = copycat["a"] === undefined ? undefined : parserAnother(copycat["a"]);
            if (parsedAnother instanceof Error) { return new Error("error when parsing field another (wire name '" + "a" + "')", { cause: parsedAnother }); }
            const parsedStr = parserStr(copycat["str"]);
            if (parsedStr instanceof Error) { return new Error("error when parsing field str (wire name '" + "str" + "')", { cause: parsedStr }); }
            const parsedEnu = parserEnu(copycat["enu"]);
            if (parsedEnu instanceof Error) { return new Error("error when parsing field enu (wire name '" + "enu" + "')", { cause: parsedEnu }); }
            // for each field: delete field from copycat object
            delete copycat["s"];
            delete copycat["n"];
            delete copycat["n2"];
            delete copycat["t2"];
            delete copycat["t4"];
            delete copycat["ti"];
            delete copycat["a"];
            delete copycat["str"];
            delete copycat["enu"];
            if (Object.keys(copycat).length > 0) { return new Error("unknown fields present: " + Object.keys(copycat).join(", ")); }
            return {
                session: parsedSession, 
                name: parsedName, 
                name2: parsedName2, 
                time2: parsedTime2, 
                time4: parsedTime4, 
                technicalIdentifier: parsedTechnicalIdentifier, 
                another: parsedAnother, 
                str: parsedStr, 
                enu: parsedEnu, 
            }
        };
        return parser(a);
    }
    
    static toJsonCore(a: __A): $J {
        const writer: (a: __A) => $J = a => {
            const wrSession: (a: string) => $J = $writeString;
            const wrName: (a: {
                s: string,
            }[]) => $J = a => {
                const coll = [] as $J[];
                for (const elem of a) {
                    const innerWriter: (a: {
                        s: string,
                    }) => $J = a => {
                        const wrS: (a: string) => $J = $writeString;
                        const ret: $J = {}
                        ret["s"] = wrS(a.s);
                        return ret;
                    }
                    coll.push(innerWriter(elem));
                }
                return coll;
            }
            const wrName2: (a: string[]) => $J = a => {
                const coll = [] as $J[];
                for (const elem of a) {
                    const innerWriter: (a: string) => $J = $writeString;
                    coll.push(innerWriter(elem));
                }
                return coll;
            };
            const wrTime2: (a: string) => $J = $writeString;
            const wrTime4: (a: {
                y: bigint,
                m: bigint,
                d: bigint,
            }) => $J = a => {
                const wrY: (a: bigint) => $J = $writeI64;
                const wrM: (a: bigint) => $J = $writeI64;
                const wrD: (a: bigint) => $J = $writeI64;
                const ret: $J = {}
                ret["y"] = wrY(a.y);
                ret["m"] = wrM(a.m);
                ret["d"] = wrD(a.d);
                return ret;
            }
            const wrTechnicalIdentifier: (a: { [_: string]: bigint }) => $J = a => {
                const coll = {} as { [_: string]: $J };
                for (const k in a) {
                    const innerWriter: (a: bigint) => $J = $writeI64;
                    coll[k] = innerWriter(a[k]!);
                }
                return coll;
            };
            const wrAnother: (a: __B) => $J = a => B.toJsonCore(a);
            const wrStr: (a: {
                struct: string,
                struct2: string,
            }) => $J = a => {
                const wrStruct: (a: string) => $J = $writeString;
                const wrStruct2: (a: string) => $J = $writeString;
                const ret: $J = {}
                ret["struct"] = wrStruct(a.struct);
                ret["struct2"] = wrStruct2(a.struct2);
                return ret;
            }
            const wrEnu: (a: {
                struct: string
            } | {
                struct2: string
            }) => $J = a => {
                type retType = {
                    struct: string
                } | {
                    struct2: string
                };
                if ("struct" in a) {
                    const writerInner = $writeString;
                    return { "struct": writerInner(a.struct) };
                }
                if ("struct2" in a) {
                    const writerInner = $writeString;
                    return { "struct2": writerInner(a.struct2) };
                }
                return $never(a);
            }
            const ret: $J = {}
            ret["s"] = wrSession(a.session);
            ret["n"] = wrName(a.name);
            ret["n2"] = wrName2(a.name2);
            if (a.time2 !== undefined) { ret["t2"] = wrTime2(a.time2); }
            if (a.time4 !== undefined) { ret["t4"] = wrTime4(a.time4); }
            ret["ti"] = wrTechnicalIdentifier(a.technicalIdentifier);
            if (a.another !== undefined) { ret["a"] = wrAnother(a.another); }
            ret["str"] = wrStr(a.str);
            ret["enu"] = wrEnu(a.enu);
            return ret;
        };
        return writer(a);
    }
    
    static toJson(a: __A): string {
        return JSON.stringify(this.toJsonCore(a));
    }
}

export class B {
    static fromJson(a: string): __B | Error {
        try { 
            const obj = JSON.parse(a);
            return this.fromJsonCore(obj);
        } catch(e) {
            if (!(e instanceof Error)) { return new Error("caught non error"); }
            return e;
        }
        
    }
    
    static fromJsonCore(a: $J): __B | Error {
        const parser = (a: $J) => {
            if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing struct"); }
            const copycat = { ...a };
            // for each field: create parsers
            const parserS = (a: $J) => A.fromJsonCore(a);
            const parserV = (a: $J) => {
                if (!(a instanceof Array)) { return new Error("expected array when parsing tuple"); }
                if (a.length !== 2) { return new Error("wrong tuple length"); }
                const parser0 = $parseString;
                const parser1 = $parseString;
                const parsed0 = parser0(a[0]!);
                if (parsed0 instanceof Error) { return new Error("failed to parse item #0 in tuple", { cause: parsed0 }); }
                const parsed1 = parser1(a[1]!);
                if (parsed1 instanceof Error) { return new Error("failed to parse item #1 in tuple", { cause: parsed1 }); }
                return [
                    parsed0,
                    parsed1,
                ] as [
                    string,
                    string,
                ];
            };
            const parserC = (a: $J) => {
                if (!(a instanceof Array)) { return new Error("expected array when parsing tuple"); }
                if (a.length !== 2) { return new Error("wrong tuple length"); }
                const parser0 = $parseI64;
                const parser1 = $parseString;
                const parsed0 = parser0(a[0]!);
                if (parsed0 instanceof Error) { return new Error("failed to parse item #0 in tuple", { cause: parsed0 }); }
                const parsed1 = parser1(a[1]!);
                if (parsed1 instanceof Error) { return new Error("failed to parse item #1 in tuple", { cause: parsed1 }); }
                return [
                    parsed0,
                    parsed1,
                ] as [
                    bigint,
                    string,
                ];
            };
            const parserD = (a: $J) => {
                if (!(a instanceof Array)) { return new Error("expected array when parsing tuple"); }
                if (a.length !== 1) { return new Error("wrong tuple length"); }
                const parser0 = $parseString;
                const parsed0 = parser0(a[0]!);
                if (parsed0 instanceof Error) { return new Error("failed to parse item #0 in tuple", { cause: parsed0 }); }
                return [
                    parsed0,
                ] as [string];
            };
            const parserMap = (a: $J) => {
                if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing map"); }
                const coll = {} as { [i: string]: __A };
                const parser = (a: $J) => A.fromJsonCore(a);
                for (const k in a) {
                    const parsed = parser(a[k]!);
                    if (parsed instanceof Error) { return new Error("failed to parse map's inner type", { cause: parsed }); } 
                    coll[k] = parsed;
                }
                return coll;
            };
            const parserBin = (a: $J) => {
                if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing map"); }
                const coll = {} as { [i: string]: Uint8Array[] };
                const parser = (a: $J) => {
                    if (!(a instanceof Array)) { return new Error("expected array while parsing list"); }
                    const coll = [] as Uint8Array[];
                    const parser = $parseBinary;
                    for (const elem of a) {
                        const parsed = parser(elem);
                        if (parsed instanceof Error) { return new Error("failed to parse list inner type", { cause: parsed }); } 
                        coll.push(parsed);
                    }
                    return coll;
                };
                for (const k in a) {
                    const parsed = parser(a[k]!);
                    if (parsed instanceof Error) { return new Error("failed to parse map's inner type", { cause: parsed }); } 
                    coll[k] = parsed;
                }
                return coll;
            };
            // for required fields only: check presence
            if (copycat["s"] === undefined) { return new Error("required field 's' (wire name '" + "s" + "') is undefined") }
            if (copycat["v"] === undefined) { return new Error("required field 'v' (wire name '" + "v" + "') is undefined") }
            if (copycat["c"] === undefined) { return new Error("required field 'c' (wire name '" + "c" + "') is undefined") }
            if (copycat["d"] === undefined) { return new Error("required field 'd' (wire name '" + "d" + "') is undefined") }
            if (copycat["map"] === undefined) { return new Error("required field 'map' (wire name '" + "map" + "') is undefined") }
            if (copycat["bin"] === undefined) { return new Error("required field 'bin' (wire name '" + "bin" + "') is undefined") }
            // for each field: parse, respecting requiredness, early return on error
            const parsedS = parserS(copycat["s"]);
            if (parsedS instanceof Error) { return new Error("error when parsing field s (wire name '" + "s" + "')", { cause: parsedS }); }
            const parsedV = parserV(copycat["v"]);
            if (parsedV instanceof Error) { return new Error("error when parsing field v (wire name '" + "v" + "')", { cause: parsedV }); }
            const parsedC = parserC(copycat["c"]);
            if (parsedC instanceof Error) { return new Error("error when parsing field c (wire name '" + "c" + "')", { cause: parsedC }); }
            const parsedD = parserD(copycat["d"]);
            if (parsedD instanceof Error) { return new Error("error when parsing field d (wire name '" + "d" + "')", { cause: parsedD }); }
            const parsedMap = parserMap(copycat["map"]);
            if (parsedMap instanceof Error) { return new Error("error when parsing field map (wire name '" + "map" + "')", { cause: parsedMap }); }
            const parsedBin = parserBin(copycat["bin"]);
            if (parsedBin instanceof Error) { return new Error("error when parsing field bin (wire name '" + "bin" + "')", { cause: parsedBin }); }
            // for each field: delete field from copycat object
            delete copycat["s"];
            delete copycat["v"];
            delete copycat["c"];
            delete copycat["d"];
            delete copycat["map"];
            delete copycat["bin"];
            if (Object.keys(copycat).length > 0) { return new Error("unknown fields present: " + Object.keys(copycat).join(", ")); }
            return {
                s: parsedS, 
                v: parsedV, 
                c: parsedC, 
                d: parsedD, 
                map: parsedMap, 
                bin: parsedBin, 
            }
        };
        return parser(a);
    }
    
    static toJsonCore(a: __B): $J {
        const writer: (a: __B) => $J = a => {
            const wrS: (a: __A) => $J = a => A.toJsonCore(a);
            const wrV: (a: [
                string,
                string,
            ]) => $J = a => {
                const writer0 = $writeString;
                const writer1 = $writeString;
                const written0 = writer0(a[0]!);
                const written1 = writer1(a[1]!);
                return [
                    written0,
                    written1,
                ];
            }
            const wrC: (a: [
                bigint,
                string,
            ]) => $J = a => {
                const writer0 = $writeI64;
                const writer1 = $writeString;
                const written0 = writer0(a[0]!);
                const written1 = writer1(a[1]!);
                return [
                    written0,
                    written1,
                ];
            }
            const wrD: (a: [string]) => $J = a => {
                const writer0 = $writeString;
                const written0 = writer0(a[0]!);
                return [
                    written0,
                ];
            };
            const wrMap: (a: { [_: string]: __A }) => $J = a => {
                const coll = {} as { [_: string]: $J };
                for (const k in a) {
                    const innerWriter: (a: __A) => $J = a => A.toJsonCore(a);
                    coll[k] = innerWriter(a[k]!);
                }
                return coll;
            };
            const wrBin: (a: { [_: string]: Uint8Array[] }) => $J = a => {
                const coll = {} as { [_: string]: $J };
                for (const k in a) {
                    const innerWriter: (a: Uint8Array[]) => $J = a => {
                        const coll = [] as $J[];
                        for (const elem of a) {
                            const innerWriter: (a: Uint8Array) => $J = $writeBinary;
                            coll.push(innerWriter(elem));
                        }
                        return coll;
                    };
                    coll[k] = innerWriter(a[k]!);
                }
                return coll;
            };
            const ret: $J = {}
            ret["s"] = wrS(a.s);
            ret["v"] = wrV(a.v);
            ret["c"] = wrC(a.c);
            ret["d"] = wrD(a.d);
            ret["map"] = wrMap(a.map);
            ret["bin"] = wrBin(a.bin);
            return ret;
        };
        return writer(a);
    }
    
    static toJson(a: __B): string {
        return JSON.stringify(this.toJsonCore(a));
    }
}

export class C {
    static fromJson(a: string): __C | Error {
        try { 
            const obj = JSON.parse(a);
            return this.fromJsonCore(obj);
        } catch(e) {
            if (!(e instanceof Error)) { return new Error("caught non error"); }
            return e;
        }
        
    }
    
    static fromJsonCore(a: $J): __C | Error {
        const parser = (a: $J) => {
            if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing map"); }
            const coll = {} as { [i: string]: string };
            const parser = $parseString;
            for (const k in a) {
                const parsed = parser(a[k]!);
                if (parsed instanceof Error) { return new Error("failed to parse map's inner type", { cause: parsed }); } 
                coll[k] = parsed;
            }
            return coll;
        };
        return parser(a);
    }
    
    static toJsonCore(a: __C): $J {
        const writer: (a: __C) => $J = a => {
            const coll = {} as { [_: string]: $J };
            for (const k in a) {
                const innerWriter: (a: string) => $J = $writeString;
                coll[k] = innerWriter(a[k]!);
            }
            return coll;
        };
        return writer(a);
    }
    
    static toJson(a: __C): string {
        return JSON.stringify(this.toJsonCore(a));
    }
}

export class D {
    static fromJson(a: string): __D | Error {
        try { 
            const obj = JSON.parse(a);
            return this.fromJsonCore(obj);
        } catch(e) {
            if (!(e instanceof Error)) { return new Error("caught non error"); }
            return e;
        }
        
    }
    
    static fromJsonCore(a: $J): __D | Error {
        const parser = (a: $J) => {
            if (!(a instanceof Array)) { return new Error("expected array while parsing list"); }
            const coll = [] as string[];
            const parser = $parseString;
            for (const elem of a) {
                const parsed = parser(elem);
                if (parsed instanceof Error) { return new Error("failed to parse list inner type", { cause: parsed }); } 
                coll.push(parsed);
            }
            return coll;
        };
        return parser(a);
    }
    
    static toJsonCore(a: __D): $J {
        const writer: (a: __D) => $J = a => {
            const coll = [] as $J[];
            for (const elem of a) {
                const innerWriter: (a: string) => $J = $writeString;
                coll.push(innerWriter(elem));
            }
            return coll;
        };
        return writer(a);
    }
    
    static toJson(a: __D): string {
        return JSON.stringify(this.toJsonCore(a));
    }
}

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
function $writeBinary(a: Uint8Array): $J { return $tob64(a) }
function $writeString(a: string): $J { return a; }

