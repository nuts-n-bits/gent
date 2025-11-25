type a = enum {

}

type b = enum {
    a: string,
    b: i64,
    c: enum {
        a: string,
        b: i64,
    }
}