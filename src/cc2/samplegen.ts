
// command arinit {
// 	-ns --namespace: string;
//	-c --comment: string[];
// }

import { CcCore } from "./lib";

type ParsedCommand = { 
	command: string; 
	args: string[]; 
	options: string[]; 
}

type __Arinit = {
	namespace: string;
	comment: string[];
}

const cccoreencode = CcCore.encode
const cccoreparse = CcCore.parse

class Arinit {
	static coreParse(pc: ParsedCommand, checkReq = true): __Arinit | Error {
		let bNamespace = "";
		let cNamespace = 0;
		let bComment = [] as string[];
		let cComment = 0;
		if (pc.command !== "arinit") {
			return new Error("command name mismatch");
		}
		for (let i=0; i<pc.options.length; i+=2) {
			switch(pc.options[i]) {
			case "-ns": case "--namespace": 
				bNamespace = pc.options[i+1]!;
				cNamespace += 1;
			break;
			case "-c": case "--comment":
				bComment.push(pc.options[i+1]!);
				cComment += 1;
			break;
			}
		}
		if (cNamespace === 0 && checkReq) {
			return new Error("missing field -ns --namespace");
		}
		return {
			namespace: bNamespace,
			comment: bComment,
		}
	}

	static coreEncode(a: __Arinit): ParsedCommand {
		return {
			command: "arinit",
			args: [],
			options: [
				"-ns", a.namespace
			],
		}
	}
}

















type __Addrev = {
        revid: string;
        revTimestamp: string;
        uid: string|undefined;
        uname: boolean;
        summary: string[];
        content: string;
}

class Addrev{
        static coreParse(pc: ParsedCommand, checkReq = true): __Addrev|Error {
                let bRevid = "";
                let cRevid = 0;
                let bRevTimestamp = "";
                let cRevTimestamp = 0;
                let bUid = undefined as string|undefined;
                let cUid = 0;
                let bUname = false;
                let cUname = 0;
                let bSummary = [] as string[];
                let cSummary = 0;
                let bContent = "";
                let cContent = 0;
                if (pc.command !== "addrev") {
                        return new Error("command name mismatch");
                }
                for (let i=0; i<pc.options.length; i+=2) {
                        switch(pc.options[i]) {
                        case "-r": 
                        case "--revid": 
                                bRevid = pc.options[i+1]!;
                                cRevid += 1;
                        break;
                        case "-t": 
                        case "--rev-timestamp": 
                                bRevTimestamp = pc.options[i+1]!;
                                cRevTimestamp += 1;
                        break;
                        case "-u": 
                        case "--uid": 
                                bUid = pc.options[i+1]!;
                                cUid += 1;
                        break;
                        case "-un": 
                        case "--uname": 
                                bUname = true;
                                cUname += 1;
                        break;
                        case "-s": 
                        case "--summary": 
                                bSummary.push(pc.options[i+1]!);
                                cSummary += 1;
                        break;
                        case "-c": 
                        case "--content": 
                                bContent = pc.options[i+1]!;
                                cContent += 1;
                        break;
                        }
                }
                if (cRevid < 1 && checkReq) {
                        return new Error("missing field -r --revid");
                }
                if (cRevTimestamp < 1 && checkReq) {
                        return new Error("missing field -t --rev-timestamp");
                }
                if (cContent < 1 && checkReq) {
                        return new Error("missing field -c --content");
                }
                return {
                        revid: bRevid,
                        revTimestamp: bRevTimestamp,
                        uid: bUid,
                        uname: bUname,
                        summary: bSummary,
                        content: bContent,
                }
        }
	static coreEncode(a: __Addrev): ParsedCommand {
                const args = [] as string[];
                const options = [] as string[];
                options.push("-r", a.revid);
                options.push("-t", a.revTimestamp);
                if (a.uid !== undefined) { options.push("-u", a.uid); }
                if (a.uname) { options.push("-un", ""); }
                for (const b of a.summary) { options.push("-s", b); }
                options.push("-c", a.content);
                return {
                        command: "addrev",
                        args: args,
                        options: options,
                }
        }
	
}