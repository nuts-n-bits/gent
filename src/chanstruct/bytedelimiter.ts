function* encode(a: Uint8Array): Iterable<Uint8Array> {
        for (let i=0; i<a.length;) {
                const rest_length = a.length - i;
                if (rest_length > 127) {
                        const chunk = new Uint8Array(128);
                        chunk[0] = 0xFF;
                        chunk.set(a.slice(i, i+127), 1);
                        i += 127;
                        yield chunk;
                } else {
                        const chunk = new Uint8Array(rest_length + 1);
                        chunk[0] = rest_length;
                        chunk.set(a.slice(i, i+rest_length), 1);
                        yield chunk;
                        return;
                }
        }
}

function txen(a: string) { return new TextEncoder().encode(a); }
function txde(a: Uint8Array) { return new TextDecoder().decode(a); }

console.log([...encode(txen("hello world, my biggest dream and fear and et cetera, et cetera. gtesogntmjsintjs gterijsgotjresigt ershjbt rsigotsje igborstjeghit rsjigotrsh jrtioht rjhiotrjhio btge"))])

; 

`
xyaaaaaa bbbbbbbb



/admin-note/mapsjeiviwv 

XYAAAAAA BBBBBBBB

XYAAAAAA BBBBBBBB
^
 \
  +------
`
