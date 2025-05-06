/*

Defintion

message Upload {
        #14 session          = text
        #24 revisions        = repeated message Revision
}

message Revision {
        #22 content = optional text
}

message User {
        #1 id     = text 
        #3 handle = text
        #6 groups = repeated message UserGroups
        #8 email   = optional text
        #9 ext_idents = repeated message UserExtIdents
        #10 open_sessions = repeated message {
                #1 device  = text
                #2 t_login = timestamp
        }
        #11 homepage_elements = enum {
                #1 anonymous = nil
                #2 admin = message AdminHomeElements
        }
}

message AdminHomeElements {
        #1 email_reset_token = optional text
        #2 password_token = optional text
        #3 noticeboard_news = repeated message {
                #1 headline = text
                #2 content_excerpt = text
                #3 author_id = text
                #4 author_name = text
        }
        #4 notifications = repeated Notification
}

message UserGroups {
        #1 ns        = text
        #2 name      = text
        #11 t_create = optional timestamp
        #12 expire   = optional timestamp
}

message UserExtIdents {
        #1 ext_host  = text
        #2 ext_ident = text
        #3 ext_email = optional text
        #4 t_create = optional timestamp
        #5 expire = optional timestamp
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

        }
        set_content(a: string): void {

        } 
}