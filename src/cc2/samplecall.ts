
import { Arclose, CcCore } from "./samplegen";

const a = Arclose.write({
    caseId: "mape824985",
    caseNumber: "48859275932",
});

console.log("parsing: ", a)
console.log(Arclose.parse(a));