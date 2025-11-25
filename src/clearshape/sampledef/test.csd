type A = {
    s session: string,
    n name: {s: string}[],
    n2 name2: string[],
    t2 time2?: string,
    t4 time4?: { y: i64, m: i64, d: i64 },
    ti technicalIdentifier: map(i64),
    //a another: B,
    str: {
        struct: string,
        struct2: string,
    },
    enu: enum {
        struct: string, 
        struct2: string,
    }
}

type B = {
    s: A,
    v: [string, string],
    c: [i64, string],
    d: [string],
    map: map(A),
}

