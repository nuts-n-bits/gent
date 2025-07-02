import { StreamlineCore, parse_one, tokenizer } from "./linecodex"
function doThrow<a>(a: a|Error): a { if (a instanceof Error) { throw a; } return a; }

const teststr = `



[0 1 2] arinit -ns zhwp;
[0] arsync -session "Xn;Df9Ghep" +1;
[1] addrev -r 398172 -t 2025-03-08T12:29:05Z -u 199272 -un Hinata -s "......"-c "......" +2;
[1] --addrev -r 399478 +2;
/1;

                        

Frame Type A: ---------
  BAEFRAME: [byte 0x00] (chanid c_1) (chanid c_2) ... (chanid c_n)
  // close channel c_i for i = 1 .. n        
-----------------------

Frame Type B: ---------
  BAEFRAME: [byte 0x01] (varint b1) (varint b2) (varint b3) ... (varint bn) (varint s)
  BAEFRAME: [data]
-----------------------

Frame Type C: ---------
  BAEFRAME: [byte 0x02] (varint b1) (varint b2) (varint b3) ... (varint bn) (varint s)
  BAEFRAME: [data]
-----------------------


`

const par = doThrow(StreamlineCore.parse_all(teststr));

console.log( par );

console.log( par.map(a => doThrow(StreamlineCore.encode(a)) ) )