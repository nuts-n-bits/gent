/*

Defintion

message Upload {
        #14 session          = text
        #24 revisions        = repeated message Revision
}

message Revision {
        #22 content = optional text
}

*/

abstract class GeneratedClass {
        abstract parseFromBuffer(buffer: Uint8Array): this
}

import { DataFrame } from "./core-serde";

function $intakeFrame(f: DataFrame, b: any, fd: string, fn: string, t: any, k: string): void {
        let value: any;
        if (fn === "text") {
                const text = new TextDecoder().decode(f.payload);
                value = text;
        }
        else if (t instanceof GeneratedClass) {
                const class_inst = t.parseFromBuffer(f.payload);
                value = class_inst;
        }

        if (k === "repeated") {
                if (b[fn] === undefined) { b[fn] = []; }
                b[fn].push(value);
        } else if (k === "optional") {
                if (b[fn] === undefined) { b[fn] = value; }
                else { throw new Error(`Optional field ${fn} (${fd}) encountered repeated value`); }
        } else { // k === "unique"
                if (b[fn] === undefined) { b[fn] = value; }
                else { throw new Error(`Unique field ${fn} (${fd}) encountered repeated value`); }
        }
}

type __RESERVEDst_6Upload = {
        session: string,
        revisions: Revision[],
}

class Upload {
        b = {} as any
        intake_frame(f: DataFrame) {
                switch(f.frame_id) {
                        case 14n: $intakeFrame(f, this.b, "14", "session", "text", "unique"); break;
                        case 24n: $intakeFrame(f, this.b, "24", "revisions", Revision, "repeated"); break; 
                }
        }
        check() {
                $checkField(this.b, "14", "session", "text", "unique");
                $checkField(this.b, "24", "revisions", Revision, "repeated");
        }
        as_obj(): __RESERVEDst_6Upload {
                return {
                        session: $accessField(this.b, "14", "session", "text", "unique"),
                        revisions: $accessField(this.b, "24", "revisions", Revision, "repeated"),
                }
        }
}

type __RESERVEDst_8Revision = {
        content: string | undefined, 
}



declare function $fromObj(o: Object, b: Object, fd: string, fn: string, t: any, k: string): void
declare function $checkField(b: any, fd: string, fn: string, t: any, k: string): void;
declare function $accessField(b: any, fd: string, fn: string, t: any, k: string): any;

class Revision {
        b = {} as any
        intake_frame(f: DataFrame) {
                switch(f.frame_id) {
                        case 5n: $intakeFrame(f, this.b, "5", "content", "text", "unique"); break;
                }
        }
        get_content(): string {
                if 
        }
        set_content(a: string): void {

        } 
}