package main

import (
	_ "embed"
	"fmt"
	"strings"
)

// golang (.go)

//go:embed runtime/lib.go
var cgGoLib string;

func cgProgramGolang(chkProgram ChkProgram, indent string, packageName string) string {
	var b strings.Builder;
	b.WriteString(strings.Replace(cgGoLib, "RUNTIME_PLACEHOLDER_PEMCJQNAENRICMQR", packageName, 1));
	for _, chkCommand := range chkProgram.commands {
		fragment := cgCommandGolang(chkCommand, indent);
		b.WriteString(fragment);
	}
	for _, chkMix := range chkProgram.mixes {
		fragment := cgMixGolang(chkMix, indent);
		b.WriteString(fragment);
	}
	return b.String();
}

func cgCommandGolang(chkCommand ChkCommandDef, indent string) string {
	var b strings.Builder;
	write := func (indentCount int, s string) {
		b.WriteString(strings.Repeat(indent, indentCount) + s + "\n");
	};
	commandProgNameGoPascal := hfcgGoBannedWordsMangle(hfNormalizedToPascal(chkCommand.commandProgName));
	write(0, "type " + commandProgNameGoPascal + " struct {");
	if chkCommand.argExists {
		typeExpr := hfcgGoType(chkCommand.argBaseType, chkCommand.argModifier);
		write(1, "Args " + typeExpr + ";");
	}
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgGoBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		typeExpr := hfcgGoType(lineDef.baseType, lineDef.modifier);
		write(1, fieldName + " " + typeExpr);
	}
	write(0, "}\n");
	// end of struct definition, begin CoreParse()
	write(0, "func (a *" + commandProgNameGoPascal + ") CoreParse(b Command) (missingFields string, err error) {");
	write(1, "missing := \"\";")
	if chkCommand.argExists {
		write(1, "var bArgs " + hfcgGoType(chkCommand.argBaseType, chkCommand.argModifier) + ";");
		if chkCommand.argModifier == ModifierOptional && true {
			write(1, "if len(b.args) > 0 {");
			write(2, "t := b.args[0];");
			write(2, "bArgs = &t");
			write(1, "}");
		} else if chkCommand.argModifier == ModifierRequired {
			write(1, "if len(b.args) == 0 {");
			write(2, `missing += "args;";`);
			write(1, "}")
			write(1, "bArgs = b.args[0];");
		} else {  // repeated
			write(1, "bArgs = b.args;");

		}
	}
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		typeExpr := hfcgGoType(lineDef.baseType, lineDef.modifier);
		initExpr := hfcgGoInitExpr(lineDef.baseType, lineDef.modifier);
		write(1, "var b" + hfNormalizedToPascal(lineDef.progName) + " " + typeExpr + " = " + initExpr + ";");
		write(1, "var c" + hfNormalizedToPascal(lineDef.progName) + " int = 0;");
	}
	// check command name match
	write(1, `if b.command != "` + chkCommand.commandName + `" {`);
	write(2, `return missing, fmt.Errorf("command name mismatch");`);
	write(1, "}");
	// big for-switch
	write(1, "for i:=0; i<len(b.options); i+=2 {");
	write(2, "switch b.options[i] {");
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgGoBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		if lineDef.longName != "" {
			write(2, fmt.Sprintf(`case "%s", "%s":`, lineDef.shortName, lineDef.longName));
		} else {
			write(2, `case "` + lineDef.shortName + `": `);
		}
		write(3, hfcgGoFieldUpdate(lineDef));
		write(3, "c" + fieldName + " += 1;");
	}
	write(2, "}");
	write(1, "}");
	// end big switch, begin checking required fields
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		if lineDef.modifier == ModifierRequired {
			fieldName := "c" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
			write(1, "if " + fieldName + " < 1 {");
			write(2, `missing += "` + lineDef.shortName + " " + lineDef.longName + `;";`);
			write(1, "}");
		}
	}
	// end checking required. begin repopulating
	if chkCommand.argExists {
		write(1, "a.Args = bArgs; ");
	}
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName1 := hfcgGoBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		fieldName2 := "b" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		write(1, "a." + fieldName1 + " = " + fieldName2 + ";");
	}
	write(1, "return missing, nil;");
	write(0, "}\n");
	// end of CoreParse(), begin CoreEncode()
	write(0, "func (a " + commandProgNameGoPascal + ") CoreEncode() Command {");
	write(1, "args := []string{};");
	write(1, "options := []string{};");
	if chkCommand.argExists {
		if chkCommand.argModifier == ModifierRequired && true {
			write(1, "args = append(args, a.Args);");
		} else if chkCommand.argModifier == ModifierOptional {
			write(1, "if a.Args != nil {");
			write(2, "args[0] = *a.Args;");
			write(1, "}");
		} else {
			write(1, "args = a.Args;\n")
		}
	}

	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgGoBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		if lineDef.baseType == BaseTypeFlag {
			write(1, "if a." + fieldName + " {  // flag");
			write(2, `options = append(options, "` + lineDef.shortName + `", "");`);
			write(1, "}");
		} else if lineDef.modifier == ModifierOptional {
			write(1, "if a." + fieldName + " != nil {");
			write(2, `options = append(options, "` + lineDef.shortName + `", *a.` + fieldName + ");");
			write(1, "}");
		} else if lineDef.modifier == ModifierRepeated {
			write(1, `for _, e := range a.` + fieldName + " {");
			write(2, `options = append(options, "` + lineDef.shortName+ `", e);`);
			write(1, "}");
		} else {		
			write(1, `options = append(options, "` + lineDef.shortName + `", a.` + fieldName + ");");
		}
	}
	write(1, "return Command{");
	write(2, `command: "` + chkCommand.commandName + `",`);
	write(2, "args: args,");
	write(2, "options: options,");
	write(1, "};");
	write(0, "}\n");
	// end of CoreEncode(), begin Parse ()
	write(0, "func (a *" + commandProgNameGoPascal + ") Parse(b string) (missingFields string, err error) {");
	write(1, "cmds, _, err := CcCore.CoreDynParse(b);")
	write(1, "if err != nil {");
	write(2, "return \"\", err;");
	write(1, "}");
	write(1, "if len(cmds) != 1 {");
	write(2, "return \"\", fmt.Errorf(\"expected exactly 1 command line\");");
	write(1, "}");
	write(1, "return a.CoreParse(cmds[0]);");
	write(0, "}\n");
	// end of Parse(), begin Write()
	write(0, "func (a " + commandProgNameGoPascal + ") Write() string {");
	write(1, "return CcCore.CoreDynEncode(a.CoreEncode());")
	write(0, "}\n");
	// end of Write()
	return b.String();
}

func cgMixGolang(chkMixDef ChkMixDef, indent string) string {
	var b strings.Builder;
	write := func (indentCount int, s string) {
		b.WriteString(strings.Repeat(indent, indentCount) + s + "\n");
	};
	mixNameGoPascal := hfcgGoBannedWordsMangle(hfNormalizedToPascal(chkMixDef.mixProgName));
	// begin type definition
	write(0, "type " + mixNameGoPascal + " struct {");
	for _, command := range chkMixDef.mixCommands {
		typeIdent := hfcgGoBannedWordsMangle(hfNormalizedToPascal(command.commandProgName));
		if command.modifier == ModifierOptional && true {
			typeIdent = "*" + typeIdent;
		} else if command.modifier == ModifierRepeated {
			typeIdent = "[]" + typeIdent;
		}
		write(1, hfcgGoBannedWordsMangle(hfNormalizedToPascal(command.commandProgName)) + " " + typeIdent);
	}
	write(0, "}\n");
	// end type definition
	



	// begin class definition
	
	write(0, fmt.Sprintf(`func (a *%s) Parse(wiredata string) (string, error) {`, mixNameGoPascal));
	write(1, "res, _, err := CcCore.CoreDynParse(wiredata);");
	write(1, `if err != nil { return "", err; }`);
	write(1, "missing := \"\";");
	write(1, "mixMissing := \"\"");
	for _, mixCommand := range chkMixDef.mixCommands {
		commandBackingVariableName := "b" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		commandMissingVariableName := "m" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		commandCountingVariableName := "c" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		backingVariableTypeBase := hfNormalizedToPascal(mixCommand.commandProgName);
		backingVariableTypeModified := "";
		if mixCommand.modifier == ModifierOptional && true {
			backingVariableTypeModified = "*" + backingVariableTypeBase + " = nil";
		} else if mixCommand.modifier == ModifierRepeated {
			backingVariableTypeModified = "[]" + backingVariableTypeBase + " = []" + backingVariableTypeBase + "{}";
		} else {
			backingVariableTypeModified = backingVariableTypeBase;
		}
		write(1, "var " + commandBackingVariableName + " " + backingVariableTypeModified + ";");
		write(1, "var " + commandMissingVariableName + " = \"\";");
		write(1, "var " + commandCountingVariableName + " = 0;");
	}
	// big for switch
	write(1, "for _, cmd := range res {");
	write(2, "switch cmd.command {");
	for _, mixCommand := range chkMixDef.mixCommands {
		commandBackingVariableName := "b" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		commandMissingVariableName := "m" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		commandCountingVariableName := "c" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		write(2, `case "` + mixCommand.commandName + `":`);
		write(3, "var parsed " + hfNormalizedToPascal(mixCommand.commandProgName) + ";");
		write(3, "mis, err := parsed.CoreParse(cmd);");
		write(3, `if err != nil { return "", err; }`);
		if mixCommand.modifier == ModifierRepeated && true {
			write(3, fmt.Sprintf("%s = append(%s, parsed);", commandBackingVariableName, commandBackingVariableName));
		} else if mixCommand.modifier == ModifierOptional {
			write(3, commandBackingVariableName + " = &parsed;");
		} else {
			write(3, commandBackingVariableName + " = parsed;");
		}
		write(3, commandMissingVariableName + " += mis;");
		write(3, commandCountingVariableName + " += 1;");
	}
	write(2, "}");
	write(1, "}");
	// end of big for switch
	// check missing
	for _, mixCommand := range chkMixDef.mixCommands {
		if mixCommand.modifier == ModifierRequired {
			commandCountingVariableName := "c" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
			write(1, "if " + commandCountingVariableName + " < 1 {");
			write(2, "mixMissing += \"" + mixCommand.commandName + ";\";");
			write(1, "}")
		}
	}
	write(1, `if len(mixMissing) > 0 { missing += "mix:" + mixMissing + "\n"; }`);
	for _, mixCommand := range chkMixDef.mixCommands {
		commandMissingVariableName := "m" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		write(1, fmt.Sprintf(`if len(%s) > 0 { missing += "%s:" + %s + "\n"; }`, commandMissingVariableName, mixCommand.commandName, commandMissingVariableName));
	}
	for _, mixCommand := range chkMixDef.mixCommands {
		commandGoFieldName := hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		commandBackingVariableName := "b" + hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		write(1, "a." + commandGoFieldName + " = " + commandBackingVariableName + ";");
	} 
	write(1, "return missing, nil;");
	write(0, "}\n");
	// end of Parse(), begin of Write()
	write(0, fmt.Sprintf("func (a *%s) Write() string {", mixNameGoPascal));
	write(1, `coll := "";`);
	for _, mixCommand := range chkMixDef.mixCommands {
		commandGoFieldName := hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		commandGoClassName := hfcgGoBannedWordsMangle(hfNormalizedToPascal(mixCommand.commandProgName));
		if mixCommand.modifier == ModifierOptional && true {
			write(1, fmt.Sprintf("if a.%s != nil {", commandGoFieldName));
			write(2, fmt.Sprintf(`coll += a.%s.Write() + "\n";`, commandGoClassName));
			write(1, "}");
		} else if mixCommand.modifier == ModifierRepeated {
			write(1, fmt.Sprintf("for _, b := range a.%s {", commandGoFieldName));
			write(2, `coll += b.Write() + "\n";`);
			write(1, "}");
		} else {
			write(1, fmt.Sprintf(`coll += a.%s.Write() + "\n";`, commandGoClassName));
		}
	}
	write(1, "return coll;")
	write(0, "}");
	return b.String();
}

// helper functions

func hfcgGoType(baseType BaseType, modifier Modifier) string {
	var t string;
	var m string;
	switch modifier {
	case ModifierOptional:
		m = "*";
	case ModifierRequired: 
		m = "";
	case ModifierRepeated:
		m = "[]";
	}
	switch baseType {
	case BaseTypeUdecimal, BaseTypeDecimal, BaseTypeString, BaseTypeBase64:
		t = m + "string";
	case BaseTypeFlag:
		t = "bool";
	default:
		t = m + "string";
	}
	return t;
}

func hfcgGoInitExpr(baseType BaseType, modifier Modifier) string {
	var initExpr string;
	switch baseType {
	case BaseTypeUdecimal, BaseTypeDecimal, BaseTypeString, BaseTypeBase64:
		initExpr = "\"\"";
	case BaseTypeFlag:
		return "false";
	default: 
		initExpr = "\"\"";
	}
	switch modifier {
	case ModifierOptional: 
		return "nil";
	case ModifierRequired:
		return initExpr;
	case ModifierRepeated:
		return "make(" + hfcgGoType(baseType, modifier) + ", 0)";
	default: 
		return initExpr;
	}
}

func hfcgGoFieldUpdate(lineDef ChkOptionDef) string {
	progName := "b" + hfNormalizedToPascal(lineDef.progName);
	if lineDef.baseType == BaseTypeFlag {
		return progName + " = true;";
	} else if lineDef.modifier == ModifierRepeated {
		return progName + " = append(" + progName + ", b.options[i+1]);";
	} else if lineDef.modifier == ModifierOptional {
		return progName + " = &b.options[i+1];";
	} else {
		return progName + " = b.options[i+1];";
	}
}

func hfcgGoBannedWordsMangle(ident string) string {
	switch ident {
	case "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "type", "var":
		return ident + "_";
	}
	return hfcgOtherBannedWordsMangle(ident);
}