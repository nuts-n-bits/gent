type Definition = {
        type: "struct" | "enum",
        ident: string,
        fields: DefinitionField[],
}

type DefinitionField = { 
        fd: number, 
        ident: string, 
        type: "str" | "bytes" | "num", 
        kind: "optional" | "required" | "repeated",
}

const def_example: Definition = {
        type: "struct",
        ident: "Arinit",
        fields: [
                { fd: 1, ident: "ns", type: "str", kind: "required" },
                { fd: 2, ident: "session", type: "str", kind: "repeated" },
                { fd: 3, ident: "pageCommonName", type: "str", kind: "optional" },
        ],
}

function type_code_of(f: DefinitionField): string {
        let s = f.type;
        let k = f.kind;
        let built = "";
        
        if (s === "str") { built += "string"; }
        else if (s === "num") { built += "string"; }
        else if (s === "bytes") { built += "Uint8Array"; }
        else { built += "unknown"; }

        if (k === "required") { }
        else if (k === "optional") { built += " | undefined"; }
        else if (k === "repeated") { built += "[]"; }

        return built;
}

function generate_class_ts(def: Definition): string {
        
        const type_ident = `___R___t_${def.ident}`;
        const class_ident = `${def.ident}`;

        const type_lines = def.fields.map(f => "\n\t" + f.ident + ": " + type_code_of(f) + ",");
        const type_string = `type ${type_ident} = {${type_lines.join("")}\n}`;

        const intake_frame_lines = def.fields.map(f => `case ${f.fd}: this.#_B_${f.ident} = ${f.type}()`);

        
        const class_constructor_lines = def.fields.map(f => `\n\t\tthis.#_B_${f.ident} = a.${f.ident};`)
        const class_constructor = `\tsetFieldsFromObject(a: ${type_ident}): void {${class_constructor_lines.join("")}\n\t}`
        
        const setter_lines = def.fields.map(f => `\n\tset_${f.ident}(a: ${type_code_of(f)}): void { this.#_B_${f.ident} = a; }`);
        const getter_lines = def.fields.map(f => `\n\tget_${f.ident}(): ${type_code_of(f)} { const a = this.#_B_${f.ident}; ${f.kind === "optional" ? "" : `if (a === undefined) { throw new Error(\"Field ${f.ident} is not supposed to be undefined\"); } ` }return a; }`)

        const class_field_lines = def.fields.map(f => "\n\t#_B_" + f.ident + ": " + type_code_of(f) + (f.kind === "optional" ? " = undefined;" : f.kind === "repeated" ? " = [];" : " | undefined = undefined;"));
        
        const class_string = `class ${class_ident} {${class_field_lines.join("")}\n\n${class_constructor}\n${setter_lines.join("")}\n${getter_lines.join("")}\n}`;
        
        return type_string + "\n\n" + class_string;
}

console.log(generate_class_ts(def_example));
/*

type ___R___t_Arinit = {
        ns: string,
        session: string[],
        pageCommonName: string | undefined,
}

class Arinit {
        #_B_ns: string | undefined = undefined;
        #_B_session: string[] = [];
        #_B_pageCommonName: string | undefined = undefined;

        setFieldsFromObject(a: ___R___t_Arinit) {
                this.#_B_ns = a.ns;
                this.#_B_session = a.session;
                this.#_B_pageCommonName = a.pageCommonName;
        }

        set_ns(a: string): void { this.#_B_ns = a; },
        set_session(a: string[]): void { this.#_B_session = a; },
        set_pageCommonName(a: string | undefined): void { this.#_B_pageCommonName = a; }
}

*/