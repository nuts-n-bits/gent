package main

// import (
// 	_ "embed"
// 	"fmt"
// 	"strings"
// )

// // typescript (.ts)

// //go:embed runtime/lib-typescript.ts
// var cgTsLib string;

// func cgProgramTypescript(chkProgram ChkProgram, indent string) string {
// 	var b strings.Builder;
// 	b.WriteString(cgTsLib);
// 	for _, chkCommand := range chkProgram.commands {
// 		fragment := cgCommandTypescript(chkCommand, indent);
// 		b.WriteString(fragment);
// 	}
// 	for _, chkMix := range chkProgram.mixes {
// 		fragment := cgMixTypescript(chkMix, indent);
// 		b.WriteString(fragment);
// 	}
// 	return b.String();
// }

// func cgCommandTypescript(chkCommand ChkCommandDef, indent string) string {
// 	var b strings.Builder;
// 	write := func (indentCount int, s string) {
// 		b.WriteString(strings.Repeat(indent, indentCount) + s + "\n");
// 	};
// 	commandProgNameTsPascal := hfcgTsBannedWordsMangle(hfNormalizedToPascal(chkCommand.commandProgName));
// 	write(0, "type __" + commandProgNameTsPascal + " = {");
// 	if chkCommand.argExists {
// 		typeExpr := hfcgTsType(chkCommand.argBaseType, chkCommand.argModifier);
// 		write(1, "args: " + typeExpr + ";");
// 	}
// 	for _, lineDef := range chkCommand.optionDefs {
// 		if lineDef.isReserved {
// 			continue;
// 		}
// 		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToCamel(lineDef.progName));
// 		typeExpr := hfcgTsType(lineDef.baseType, lineDef.modifier);
// 		write(1, fieldName + ": " + typeExpr + ";");
// 	}
// 	write(0, "}\n");

// 	write(0, "export class " + commandProgNameTsPascal + "{");
// 	write(1, "static coreParse(pc: Command): { res: __" + commandProgNameTsPascal + ", missing: string }|Error {");
// 	write(2, "let missing = \"\";");
// 	if chkCommand.argExists {
// 		if chkCommand.argModifier == ModifierOptional && true {
// 			write(2, "let bArgs = pc.args[0];");
// 		} else if chkCommand.argModifier == ModifierRequired {
// 			write(2, "if (pc.args[0] === undefined) {");
// 			write(3, `missing += "args;";`);
// 			write(2, "}")
// 			write(2, "let bArgs = pc.args[0]!;");
// 		} else {
// 			write(2, "let bArgs = pc.args;");
// 		}
// 	}
// 	for _, lineDef := range chkCommand.optionDefs {
// 		if lineDef.isReserved {
// 			continue;
// 		}
// 		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName))
// 		write(2, "let b" + fieldName + " = " + hfcgTsTypeInitExpr(lineDef.baseType, lineDef.modifier) + ";");
// 		write(2, "let c" + fieldName + " = 0;");
// 	}
// 	// check command name match
// 	write(2, `if (pc.command !== "` + chkCommand.commandName + `") {`);
// 	write(3, `return new Error("command name mismatch");`);
// 	write(2, "}")
// 	// big switch statement
// 	write(2, "for (let i=0; i<pc.options.length; i+=2) {");
// 	write(3, "switch(pc.options[i]) {");
// 	for _, lineDef := range chkCommand.optionDefs {
// 		if lineDef.isReserved {
// 			continue;
// 		}
// 		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
// 		write(3, `case "` + lineDef.shortName + `": `);
// 		if lineDef.longName != "" {
// 			write(3, `case "` + lineDef.longName + `": `);
// 		}
// 		write(4, hfcgTsFieldUpdate(lineDef));
// 		write(4, "c" + fieldName + " += 1;");
// 		write(3, "break;");
// 	}
// 	write(3, "}");
// 	write(2, "}");
// 	// end big switch, start checking required fields
// 	for _, lineDef := range chkCommand.optionDefs {
// 		if lineDef.isReserved {
// 			continue;
// 		}
// 		if lineDef.modifier == ModifierRequired {
// 			fieldName := "c" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
// 			write(2, "if (" + fieldName + " < 1) {");
// 			write(3, `missing += "` + lineDef.shortName + " " + lineDef.longName + `;";`);
// 			write(2, "}");
// 		}
// 	}
// 	// end checking required. begin return
// 	write(2, "return {");
// 	write(3, "res: {");
// 	if chkCommand.argExists {
// 		write(4, "args: bArgs, ");
// 	}
// 	for _, lineDef := range chkCommand.optionDefs {
// 		if lineDef.isReserved {
// 			continue;
// 		}
// 		fieldName1 := hfcgTsBannedWordsMangle(hfNormalizedToCamel(lineDef.progName));
// 		fieldName2 := "b" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
// 		write(4, fieldName1 + ": " + fieldName2 + ",");
// 	}
// 	write(3, "},");
// 	write(3, "missing: missing,");
// 	write(2, "}");
// 	write(1, "}");
// 	// end coreParse() function, begin coreEncode() function
// 	write(1, "static coreEncode(a: __" + commandProgNameTsPascal + "): Command {");
// 	if chkCommand.argExists {
// 		if chkCommand.argModifier == ModifierRequired && true {
// 			write(2, "const args = [a.args];");
// 		} else if chkCommand.argModifier == ModifierOptional {
// 			write(2, "const args = a.args !== undefined ? [a.args] : [];")
// 		} else {
// 			write(2, "const args = a.args;");
// 		}
// 	} else {
// 		write(2, "const args = [] as string[];");
// 	}
// 	write(2, "const options = [] as string[];");
// 	for _, lineDef := range chkCommand.optionDefs {
// 		if lineDef.isReserved {
// 			continue;
// 		}
// 		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToCamel(lineDef.progName));
// 		if lineDef.baseType == BaseTypeFlag {
// 			write(2, "if (a." + fieldName + `) { options.push("` + lineDef.shortName + `", ""` + "); }");
// 		} else if lineDef.modifier == ModifierOptional {
// 			write(2, "if (a." + fieldName + ` !== undefined) { options.push("` + lineDef.shortName + `", a.` + fieldName + "); }")
// 		} else if lineDef.modifier == ModifierRepeated {
// 			write(2, `for (const b of a.` + fieldName +`) { options.push("` + lineDef.shortName + `", b); }`);
// 		} else {		
// 			write(2, `options.push("` + lineDef.shortName + `", a.` + fieldName + ");")
// 		}
// 	}
// 	write(2, "return {");
// 	write(3, `command: "` + chkCommand.commandName + `",`);
// 	write(3, "args: args,");
// 	write(3, "options: options,");
// 	write(2, "}");
// 	write(1, "}");
// 	write(1, "static write(a: __" + commandProgNameTsPascal + "): string {");
// 	write(2, "return CcCore.encode(this.coreEncode(a));");
// 	write(1, "}");
// 	write(1, "static parse(s: string): { res: __" + commandProgNameTsPascal + ", missing: string }|Error {");
// 	write(2, "const e = CcCore.parse(s);");
// 	write(2, "if (e instanceof Error) { return e; }");
// 	write(2, "if (e.length !== 1) { return new Error(\"expected exactly 1 command line\"); }");
// 	write(2, "const f = e[0]!;");
// 	write(2, "return this.coreParse(f);");
// 	write(1, "}");
// 	write(1, fmt.Sprintf("static emptyValue(): { res: __%s, missing: string } {", commandProgNameTsPascal));
// 	write(2, fmt.Sprintf(`const ev = this.coreParse({ command: "%s", options: [], args: [] });`, chkCommand.commandName));
// 	write(2, "if (ev instanceof Error) { throw ev; }");
// 	write(2, "return ev;");
// 	write(1, "}")
// 	write(0, "}\n");
// 	return b.String();
// }

// func cgMixTypescript(chkMixDef ChkMixDef, indent string) string {
// 	var b strings.Builder;
// 	write := func (indentCount int, s string) {
// 		b.WriteString(strings.Repeat(indent, indentCount) + s + "\n");
// 	};
// 	mixNameTsPascal := hfcgTsBannedWordsMangle(hfNormalizedToPascal(chkMixDef.mixProgName));
// 	// begin type definition
// 	write(0, "type __" + mixNameTsPascal + " = {");
// 	for _, command := range chkMixDef.mixCommands {
// 		typeIdent := "__" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(command.commandProgName));
// 		if command.modifier == ModifierOptional && true {
// 			typeIdent += "|undefined";
// 		} else if command.modifier == ModifierRepeated {
// 			typeIdent += "[]";
// 		}
// 		write(1, hfcgTsBannedWordsMangle(hfNormalizedToCamel(command.commandProgName)) + ": " + typeIdent + ",");
// 	}
// 	write(0, "}\n");
// 	// end type definition
// 	// begin class definition
// 	write(0, "export class " + mixNameTsPascal + " {");
// 	write(1, "static parse(wiredata: string): { mix: __" + mixNameTsPascal + ", missing: string }|Error {");
// 	write(2, "const res = CcCore.parse(wiredata);");
// 	write(2, "if (res instanceof Error) { return res; }");
// 	write(2, "let missing = \"\";");
// 	write(2, "let mixMissing = \"\"");
// 	for _, mixCommand := range chkMixDef.mixCommands {
// 		commandBackingVariableName := "b" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		commandMissingVariableName := "m" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		commandCountingVariableName := "c" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		backingVariableTypeBase := "__" + hfNormalizedToPascal(mixCommand.commandProgName);
// 		backingVariableTrailer := "";
// 		if mixCommand.modifier == ModifierOptional && true {
// 			backingVariableTrailer = "|undefined = undefined";
// 		} else if mixCommand.modifier == ModifierRepeated {
// 			backingVariableTrailer = "[] = []";
// 		}
// 		write(2, "let " + commandBackingVariableName + ": " + backingVariableTypeBase + backingVariableTrailer + ";");
// 		write(2, "let " + commandMissingVariableName + " = \"\";");
// 		write(2, "let " + commandCountingVariableName + " = 0;");
// 	}
// 	// big for switch
// 	write(2, "for (const cmd of res) {");
// 	write(3, "switch (cmd.command) {");
// 	for i, mixCommand := range chkMixDef.mixCommands {
// 		tempVarName := "parsed" + fmt.Sprint(i);
// 		commandBackingVariableName := "b" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		commandMissingVariableName := "m" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		commandCountingVariableName := "c" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		write(3, `case "` + mixCommand.commandName + `":`);
// 			write(4, "const " + tempVarName + " = " + hfNormalizedToPascal(mixCommand.commandProgName) + ".coreParse(cmd);");
// 			write(4, "if (" + tempVarName + " instanceof Error) { return " + tempVarName + "; }");
// 			if mixCommand.modifier == ModifierRepeated {
// 				write(4, commandBackingVariableName + ".push(" + tempVarName + ".res);");
// 			} else {
// 				write(4, commandBackingVariableName + " = " + tempVarName + ".res;");
// 			}
// 			write(4, commandMissingVariableName + " += " + tempVarName + ".missing;");
// 			write(4, commandCountingVariableName + " += 1;");
// 		write(3, "break;");
// 	}
// 	write(3, "}");
// 	write(2, "}");
// 	// end of big for switch
// 	// check missing
// 	for _, mixCommand := range chkMixDef.mixCommands {
// 		if mixCommand.modifier == ModifierRequired {
// 			commandClassName := hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 			commandBackingVariableName := "b" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 			commandCountingVariableName := "c" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 			write(2, "if (" + commandCountingVariableName + " < 1) {");
// 			write(3, "mixMissing += \"" + mixCommand.commandName + ";\";");
// 			write(3, fmt.Sprintf("const { res, missing } = %s.emptyValue();", commandClassName));
// 			write(3, commandBackingVariableName + " = res;");
// 			write(2, "}")
// 		}
// 	}
// 	write(2, `if (mixMissing.length > 0) { missing += "mix:" + mixMissing + "\n"; }`);
// 	for _, mixCommand := range chkMixDef.mixCommands {
// 		commandMissingVariableName := "m" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		write(2, fmt.Sprintf(`if (%s.length > 0) { missing += "%s:" + %s + "\n"; }`, commandMissingVariableName, mixCommand.commandName, commandMissingVariableName));
// 	}
// 	write(2, "return { missing: missing, mix: {");
// 	for _, mixCommand := range chkMixDef.mixCommands {
// 		commandBackingVariableName := "b" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		commandTsFieldName := hfcgTsBannedWordsMangle(hfNormalizedToCamel(mixCommand.commandProgName));
// 		write(3, commandTsFieldName + ": " + commandBackingVariableName + "!,");
// 	}
// 	write(2, "} }; ");
// 	write(1, "}");
// 	// end of parse(), begin of write()
// 	write(1, fmt.Sprintf("static write(a: __%s): string {", mixNameTsPascal));
// 	write(2, "let coll = \"\";");
// 	for _, mixCommand := range chkMixDef.mixCommands {
// 		commandTsFieldName := hfcgTsBannedWordsMangle(hfNormalizedToCamel(mixCommand.commandProgName));
// 		commandTsClassName := hfcgTsBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
// 		if mixCommand.modifier == ModifierOptional && true {
// 			write(2, fmt.Sprintf("if (a.%s !== undefined) {", commandTsFieldName));
// 			write(3, fmt.Sprintf(`coll += %s.write(a.%s) + "\n";`, commandTsClassName, commandTsFieldName));
// 			write(2, "}");
// 		} else if mixCommand.modifier == ModifierRepeated {
// 			write(2, fmt.Sprintf("for (const b of a.%s) {", commandTsFieldName));
// 			write(3, fmt.Sprintf(`coll += %s.write(b) + "\n";`, commandTsClassName));
// 			write(2, "}");
// 		} else {
// 			write(2, fmt.Sprintf(`coll += %s.write(a.%s) + "\n";`, commandTsClassName, commandTsFieldName));
// 		}
// 	}
// 	write(2, "return coll;")
// 	write(1, "}");
// 	write(0, "}");
// 	return b.String();
// }


// // helper functions

// func hfcgTsTypeInitExpr(baseType BaseType, modifier Modifier) string {
// 	var zeroVal, typeName string;
// 	switch baseType {
// 	case BaseTypeUdecimal, BaseTypeDecimal, BaseTypeString, BaseTypeBase64:
// 		zeroVal = "\"\"";
// 		typeName = "string";
// 	case BaseTypeFlag:
// 		return "false";
// 	default: 
// 		zeroVal = "\"\"";
// 		typeName = "string";
// 	}
// 	switch modifier {
// 	case ModifierOptional: 
// 		return "undefined as " + typeName + "|undefined";
// 	case ModifierRequired:
// 		return zeroVal;
// 	case ModifierRepeated:
// 		return "[] as " + typeName + "[]";
// 	default: 
// 		return zeroVal;
// 	}
// }

// func hfcgTsType(baseType BaseType, modifier Modifier) string {
// 	var t string;
// 	var m string;
// 	switch modifier {
// 	case ModifierOptional: 
// 		m = "|undefined";
// 	case ModifierRequired:
// 		m = "";
// 	case ModifierRepeated:
// 		m = "[]";
// 	}
// 	switch baseType {
// 	case BaseTypeUdecimal, BaseTypeDecimal, BaseTypeString, BaseTypeBase64:
// 		t = "string" + m;
// 	case BaseTypeFlag:
// 		t = "boolean";
// 	default: 
// 		t = "string" + m;
// 	}
// 	return t;
// }

// func hfcgTsFieldUpdate(lineDef ChkOptionDef) string {
// 	progName := "b" + hfNormalizedToPascal(lineDef.progName);
// 	if lineDef.baseType == BaseTypeFlag {
// 		return progName + " = true;"
// 	} else if lineDef.modifier == ModifierRepeated {
// 		return progName + ".push(pc.options[i+1]!);";
// 	} else {
// 		return progName + " = pc.options[i+1]!;";
// 	}
// }

// func hfcgTsBannedWordsMangle(ident string) string {
// 	switch ident {
// 	case "break", "case", "catch", "class", "const", "continue", "debugger", "default", "delete", "do", "else", "enum", "export", "extends",  "false", "finally", "for", "function",  "if", "import", "in", "instanceof",  "new", "null",  "return",  "super", "switch",  "this", "throw", "true", "try", "typeof", "var", "void", "while", "with", "yield":
// 		return ident + "_";
// 	}
// 	return hfcgOtherBannedWordsMangle(ident);
// }