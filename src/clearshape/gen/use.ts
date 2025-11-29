import { A } from "./out"

const a = A.fromJson(`
{"s":"ioiu887","n":[],"n2":[],"ti":{},"str":{"struct":"","struct2":""},"enu":{"struct2":""},"t2":"","t4":{"y":"1799987654345689","m":"","d":""}}
`)

console.log(a)

if (a instanceof Error) { throw a }

console.log(A.toJson(a))