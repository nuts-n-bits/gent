package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// tokenizer ///////////////////////////////////////////////////////////////

type TokenKind string;
const (
	TokenKwCommand TokenKind = "TokenKwCommand"
	TokenKwArgs TokenKind = "TokenKwArgs"
	TokenKwReserved TokenKind= "TokenKwReserved"
	TokenKwAs TokenKind = "TokenKwAs"
	TokenOpenBrace TokenKind = "TokenOpenBrace"
	TokenCloseBrace TokenKind = "TokenCloseBrace"
	TokenColon TokenKind = "TokenColon"
	TokenSemicolon TokenKind = "TokenSemicolon"
	TokenOpenBracket TokenKind = "TokenOpenBracket"
	TokenCloseBracket TokenKind = "TokenCloseBracket"
	TokenQuestion TokenKind = "TokenQuestion"
	TokenIdentLike TokenKind = "TokenIdentLike"
	TokenEof TokenKind = "TokenEof"
)

type Token struct {
	kind TokenKind
	data string
	start int
	end int
}

func lexTokenizer(program string) ([]Token, error, int) {
	tokens := make([]Token, 0);
	i := 0;
	for {
		if i >= len(program) {
			tokens = append(tokens, Token{ kind: TokenEof, start: i, end: i });
			return tokens, nil, 0;
		} else if program[i] == ' ' || program[i] == '\t' || program[i] == '\n' || program[i] == '\r' {
			i += 1;
		} else if program[i] == '/' && program[i+1] == '/' {
			for program[i] != '\n' {
				i += 1;
			}
		} else if program[i] == '{' {
			tokens = append(tokens, Token{ kind: TokenOpenBrace, start: i, end: i+1 });
			i += 1;
		} else if program[i] == '}' {
			tokens = append(tokens, Token{ kind: TokenCloseBrace, start: i, end: i+1 });
			i += 1;
		} else if program[i] == '[' {
			tokens = append(tokens, Token{ kind: TokenOpenBracket, start: i, end: i+1 });
			i += 1;
		} else if program[i] == ']' {
			tokens = append(tokens, Token{ kind: TokenCloseBracket, start: i, end: i+1 });
			i += 1;
		} else if program[i] == '?' {
			tokens = append(tokens, Token{ kind: TokenQuestion, start: i, end: i+1 });
			i += 1;
		} else if program[i] == ':' {
			tokens = append(tokens, Token{ kind: TokenColon, start: i, end: i+1 });
			i += 1;
		} else if program[i] == ';' {
			tokens = append(tokens, Token{ kind: TokenSemicolon, start: i, end: i+1 });
			i += 1;
		} else if hfIsIdentLike(program[i]) {
			ident := lexConsumeIdentLike(program, i);
			i2 := i + len(ident);
			switch ident {
			case "command": 
				tokens = append(tokens, Token{ kind: TokenKwCommand, start: i, end: i2 });
			case "args":
				tokens = append(tokens, Token{ kind: TokenKwArgs, start: i, end: i2 });
			case "reserved":
				tokens = append(tokens, Token{ kind: TokenKwReserved, start: i, end: i2 });
			case "as":
				tokens = append(tokens, Token{ kind: TokenKwAs, start: i, end: i2});
			default:
				tokens = append(tokens, Token{ kind: TokenIdentLike, data: ident, start: i, end: i2 });
			}
			i = i2;
		} else {
			return tokens, fmt.Errorf("unexpected character at %d", i), i;
		}
	}
}

func lexConsumeIdentLike(program string, i int) string {
	identColl := make([]byte, 0);
	for {
		if i >= len(program) {
			return string(identColl);
		} else if hfIsIdentLike(program[i]) {
			identColl = append(identColl, byte(program[i]));
			i += 1;
		} else {
			return string(identColl);
		}
	}
}

//////// tokenizer ends //////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// parser //////////////////////////////////////////////////////////////////

type AstProgram struct {
	commandsDefs []AstCommandBlock
}

type AstCommandBlock struct {
	commandName Token
	lineDefs []AstLineDef
}

type AstLineDef struct {
	argsTok *Token
	shortName *Token  // no short name means the line defines argument type
	longName *Token  // no long name means fallback to short name
	typeBase *Token  // no type base means it's a `resevered` line
	typeMod Modifier  // reserved line makes modifier meaningless
	reservedTok *Token
}

type BaseType string

const (
	BaseTypeString BaseType = "BaseTypeString"
	BaseTypeUdecimal BaseType = "BaseTypeUdecimal"
	BaseTypeDecimal BaseType = "BaseTypeDecimal"
	BaseTypeBase64 BaseType = "BaseTypeBase64"
	BaseTypeFlag BaseType = "BaseTypeFlag"
)

type Modifier string

const (
	ModifierRequired Modifier = "ModifierRequired"
	ModifierRepeated Modifier = "ModifierRepeated"
	ModifierOptional Modifier = "ModifierOptional"
)

func rdParseProgram(tokens []Token) (AstProgram, int, error) {
	program := AstProgram{};
	i := 0;
	for {
		switch tokens[i].kind {
		case TokenKwCommand:
			commandBlock, newI, err := rdParseCommandBlock(tokens, i);
			if err != nil {
				return AstProgram{}, newI, err;
			}
			program.commandsDefs = append(program.commandsDefs, commandBlock);
			i = newI;
		case TokenEof:
			return program, i, nil;
		default:
			return AstProgram{}, i, fmt.Errorf("unexpected token %s(%s)", tokens[i].kind, tokens[i].data);
		}
	}
}

func rdParseCommandBlock(tokens []Token, i int) (AstCommandBlock, int, error) {
	commandDef := AstCommandBlock{};
	// consume keyword `command`, or error
	if tokens[i].kind == TokenKwCommand {
		i += 1;
	} else {
		return AstCommandBlock{}, i, fmt.Errorf("expected command keyword while parsing command block");
	}
	// consume command name, or error
	if tokens[i].kind == TokenIdentLike {
		commandDef.commandName = tokens[i];
		i += 1;
	} else {
		return AstCommandBlock{}, i, fmt.Errorf("expected command name while parsing command block");
	}
	// consume `{`, or error
	if tokens[i].kind == TokenOpenBrace {
		i += 1;
	} else {		
		return AstCommandBlock{}, i, fmt.Errorf("expected OpenBrace while parsing command block");
	}
	for {
		if tokens[i].kind == TokenCloseBrace {
			i += 1;
			return commandDef, i, nil;
		}
		lineDef := AstLineDef{}
		// consume LHS of line before `:`, could be `-s` or `-s --long-name` or `args`, or error.
		if tokens[i].kind == TokenKwArgs {
			t := tokens[i];
			lineDef.argsTok = &t;
			i += 1;
		} else if tokens[i].kind == TokenIdentLike && tokens[i+1].kind == TokenIdentLike {
			t := tokens[i];
			lineDef.shortName = &t;
			u := tokens[i+1];
			lineDef.longName = &u;
			i += 2;
		} else if tokens[i].kind == TokenIdentLike {
			t := tokens[i];
			lineDef.shortName = &t;
			i += 1;
		} else {
			return AstCommandBlock{}, i, fmt.Errorf("expected keyword `args` or 1-2 option names at " + 
			"the beginning of a line definition in a command block, found %s", tokens[i].kind);
		}
		// consume `:`
		if tokens[i].kind == TokenColon {
			i += 1;
		} else {
			return AstCommandBlock{}, i, fmt.Errorf("expected `:` after option names or `args`, found %s",
			tokens[i].kind);
		}
		// consume RHS of line before `;`, could be `type`, `type[]`, `type?` or `reserved`
		if tokens[i].kind == TokenIdentLike && tokens[i+1].kind == TokenOpenBracket && 
		tokens[i+2].kind == TokenCloseBracket {
			t := tokens[i];
			lineDef.typeBase = &t;
			lineDef.typeMod = ModifierRepeated;
			i += 3;
		} else if tokens[i].kind == TokenIdentLike && tokens[i+1].kind == TokenQuestion {
			t := tokens[i];
			lineDef.typeBase = &t;
			lineDef.typeMod = ModifierOptional;
			i += 2;
		} else if tokens[i].kind == TokenIdentLike {
			t := tokens[i];
			lineDef.typeBase = &t;
			lineDef.typeMod = ModifierRequired;
			i += 1;
		} else if tokens[i].kind == TokenKwReserved {
			t := tokens[i];
			lineDef.reservedTok = &t;
			i += 1;
		} else {
			return AstCommandBlock{}, i, fmt.Errorf(
			"expected type name and modifier or `reserved` keyword after `:`, found %s", tokens[i].kind);
		}
		commandDef.lineDefs = append(commandDef.lineDefs, lineDef);
		// consume `;}` to return, or `}` to return, or `;` to continue, or error.
		if tokens[i].kind == TokenSemicolon && tokens[i+1].kind == TokenCloseBrace {
			i += 2;
			return commandDef, i, nil;
		} else if tokens[i].kind == TokenSemicolon {
			i += 1;
		} else if tokens[i].kind == TokenCloseBrace {
			i += 1;
			return commandDef, i, nil;
		} else {
			return AstCommandBlock{}, i, 
			fmt.Errorf("expected `;` or `}` after a line, found %s", tokens[i].kind);
		}
	}
}

//////// parser ends /////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// checker /////////////////////////////////////////////////////////////////

type ChkProgram struct {
	commands []ChkCommandDef
}

type ChkCommandDef struct {
	commandName string
	commandProgName string
	commandProgNameSrcTok *Token
	argExists bool
	argBaseType BaseType  // can be empty in case of !argExists
	argModifier Modifier  // can be empty in case of !argExists
	optionDefs []ChkOptionDef
}

type ChkOptionDef struct {
	shortName string
	longName string  // can be empty in case longName is not specified
	progName []string
	isReserved bool
	baseType BaseType  // can be empty in case of reserved
	modifier Modifier  // can be empty in case of reserved
}

func chkProgram(astProgram AstProgram) (ChkProgram, *Token, error) {
	checkedProgram := ChkProgram{};
	// check each astCommandBlock
	for _, astCommandBlock := range astProgram.commandsDefs {
		checkedCommandDef, errTok, err := chkCommandDef(astCommandBlock);
		if err != nil {
			return ChkProgram{}, errTok, err;
		}
		checkedProgram.commands = append(checkedProgram.commands, checkedCommandDef);
	}
	// check duplicate commandProgName
	seen := map[string]bool{};
	for _, checkedCommandDef := range checkedProgram.commands {
		if seen[checkedCommandDef.commandProgName] {
			return ChkProgram{}, checkedCommandDef.commandProgNameSrcTok, 
			fmt.Errorf("duplicate command program name %s (from command: %s)", 
			checkedCommandDef.commandProgName, checkedCommandDef.commandName);
		}
		seen[checkedCommandDef.commandProgName] = true;
	}
	return checkedProgram, nil, nil;
}

func chkCommandDef(astCommandBlock AstCommandBlock) (ChkCommandDef, *Token, error) {
	checkedCommandDef := ChkCommandDef{};
	// populate commandName
	err := hfCommandNameCheck(astCommandBlock.commandName.data);
	if err != nil {
		return ChkCommandDef{}, &astCommandBlock.commandName, err;
	}
	checkedCommandDef.commandName = astCommandBlock.commandName.data;
	// populate progName by transforming commandName
	checkedCommandDef.commandProgName = hfNormalizedToPascal(hfNormalizeIdentLike(astCommandBlock.commandName.data));
	checkedCommandDef.commandProgNameSrcTok = &astCommandBlock.commandName;
	// convert each line inside the block
	for _, astLineDef := range astCommandBlock.lineDefs {
		if astLineDef.argsTok != nil {  // this means this line is `args: ...;` ignore shortName, ignore longName
			if checkedCommandDef.argExists {
				return ChkCommandDef{}, astLineDef.argsTok, fmt.Errorf("duplicate argumetns definition");
			}
			if astLineDef.reservedTok != nil {
				return ChkCommandDef{}, astLineDef.reservedTok, fmt.Errorf("args cannot be reserved");
			}
			// populate argBaseType and argModifier
			baseType, modifier, err := hfIsLegalType(astLineDef.typeBase.data, astLineDef.typeMod, true);
			if err != nil {
				return ChkCommandDef{}, astLineDef.typeBase, err;
			}
			checkedCommandDef.argExists = true;
			checkedCommandDef.argBaseType = baseType;
			checkedCommandDef.argModifier = modifier;
			continue;
		} 
		// no argsTok, this line defines shortName and maybe longName
		// populate the corresponding fields in checked struct
		checkedOptionDef := ChkOptionDef{};
		err := hfOptionNameCheck(astLineDef.shortName.data, false);
		if err != nil {
			return ChkCommandDef{}, astLineDef.shortName, err;
		}
		checkedOptionDef.shortName = astLineDef.shortName.data;
		if astLineDef.longName != nil {
			err := hfOptionNameCheck(astLineDef.longName.data, true);
			if err != nil {
				return ChkCommandDef{}, astLineDef.longName, err;
			}
			checkedOptionDef.longName = astLineDef.longName.data;
			checkedOptionDef.progName = hfNormalizeIdentLike(astLineDef.longName.data);
		} else {
			err := hfOptionNameCheck(astLineDef.shortName.data, true);
			if err != nil {
				return ChkCommandDef{}, astLineDef.shortName, err;
			}
			checkedOptionDef.progName = hfNormalizeIdentLike(astLineDef.shortName.data);
		}
		// populate argBaseType and argModifier if not reserved
		if astLineDef.reservedTok != nil {
			checkedOptionDef.isReserved = true;
		} else {
			baseType, modifier, err := hfIsLegalType(astLineDef.typeBase.data, astLineDef.typeMod, false);
			if err != nil {
				return ChkCommandDef{}, astLineDef.typeBase, err;
			}
			checkedOptionDef.baseType = baseType;
			checkedOptionDef.modifier = modifier;
		}
		checkedCommandDef.optionDefs = append(checkedCommandDef.optionDefs, checkedOptionDef);
	}
	// check duplicate shortName, longName, and progName
	seen := map[string]bool{};
	for _, astLine := range astCommandBlock.lineDefs {
		// check short name
		if astLine.shortName != nil {
			if seen[astLine.shortName.data] {
				return ChkCommandDef{}, astLine.shortName, 
				fmt.Errorf("duplicate option name %s", astLine.shortName.data);
			} else {
				seen[astLine.shortName.data] = true;
			}
		}
		// check long name
		if astLine.longName != nil {
			if seen[astLine.longName.data] {
				return ChkCommandDef{}, astLine.longName, 
				fmt.Errorf("duplicate option name %s", astLine.longName.data);
			} else {
				seen[astLine.longName.data] = true;
			}
		}
		// check prog name
		if astLine.shortName == nil {
			// no need to check, is an `args: type;` line
		} else if astLine.shortName != nil && astLine.longName == nil {
			// is a shorthand line `-s: type`, normalize shortName
			pascalName := hfNormalizedToPascal(hfNormalizeIdentLike(astLine.shortName.data));
			if seen[pascalName] {
				return ChkCommandDef{}, astLine.shortName, 
				fmt.Errorf("duplicate normalized name %s (%s)", pascalName, astLine.shortName.data);
			}
			seen[pascalName] = true;
		} else {
			// is a full line `-s --shorthand: type`, normalize longName
			pascalName := hfNormalizedToPascal(hfNormalizeIdentLike(astLine.longName.data));
			if seen[pascalName] {
				return ChkCommandDef{}, astLine.longName,
				fmt.Errorf("duplicated normalized name %s (%s)", pascalName, astLine.longName.data);
			}
			seen[pascalName] = true;
		}
	}
	return checkedCommandDef, nil, nil;
}

//////// checker ends ////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// codegen /////////////////////////////////////////////////////////////////

// typescript (.ts)

//go:embed runtime/lib.ts
var cgTsLib string;

func cgProgramTypescript(chkProgram ChkProgram, indent string) string {
	var b strings.Builder;
	b.WriteString(cgTsLib);
	for _, chkCommand := range chkProgram.commands {
		fragment := cgCommandTypescript(chkCommand, indent);
		b.WriteString(fragment);
	}
	return b.String();
}

func cgCommandTypescript(chkCommand ChkCommandDef, indent string) string {
	var b strings.Builder;
	write := func (indentCount int, s string) {
		b.WriteString(strings.Repeat(indent, indentCount) + s);
	}
	write(0, "type __" + chkCommand.commandProgName + " = {\n");
	if chkCommand.argExists {
		typeExpr := hfcgTsType(chkCommand.argBaseType, chkCommand.argModifier);
		write(1, "args: " + typeExpr + ";\n");
	}
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToCamel(lineDef.progName));
		typeExpr := hfcgTsType(lineDef.baseType, lineDef.modifier);
		write(1, fieldName + ": " + typeExpr + ";\n");
	}
	write(0, "}\n\n");

	write(0, "export class " + chkCommand.commandProgName + "{\n");
	write(1, "static coreParse(pc: Command, checkReq = true): __" + chkCommand.commandProgName + "|Error {\n");
	if chkCommand.argExists {
		if chkCommand.argModifier == ModifierOptional && true {
			write(2, "let bArgs = pc.args[0];\n");
		} else if chkCommand.argModifier == ModifierRequired {
			write(2, "if (pc.args[0] === undefined) {\n");
			write(3, `return new Error("missing required arguments");` + "\n");
			write(2, "}\n")
			write(2, "let bArgs = pc.args[0]!;\n");
		} else {
			write(2, "let bArgs = pc.args;\n")
		}
	}
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName))
		write(2, "let b" + fieldName + " = " + hfcgTsTypeInitExpr(lineDef.baseType, lineDef.modifier) + ";\n");
		write(2, "let c" + fieldName + " = 0;\n");
	}
	write(2, `if (pc.command !== "` + chkCommand.commandName + `") {` + "\n");
	write(3, `return new Error("command name mismatch");` + "\n");
	write(2, "}\n")
	// big switch statement
	write(2, "for (let i=0; i<pc.options.length; i+=2) {\n");
	write(3, "switch(pc.options[i]) {\n");
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		write(3, `case "` + lineDef.shortName + `": ` + "\n");
		if lineDef.longName != "" {
			write(3, `case "` + lineDef.longName + `": ` + "\n");
		}
		write(4, hfcgTsFieldUpdate(lineDef));
		write(4, "c" + fieldName + " += 1;\n");
		write(3, "break;\n");
	}
	write(3, "}\n")
	write(2, "}\n")
	// end big switch, start checking required fields
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		if lineDef.modifier == ModifierRequired {
			fieldName := "c" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
			write(2, "if (" + fieldName + " < 1 && checkReq) {\n");
			write(3, `return new Error("missing field ` + lineDef.shortName + " " + lineDef.longName + `");` + "\n");
			write(2, "}\n");
		}
	}
	// end checking required. begin return
	write(2, "return {\n");
	if chkCommand.argExists {
		write(3, "args: bArgs, \n");
	}
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName1 := hfcgTsBannedWordsMangle(hfNormalizedToCamel(lineDef.progName));
		fieldName2 := "b" + hfcgTsBannedWordsMangle(hfNormalizedToPascal(lineDef.progName));
		write(3, fieldName1 + ": " + fieldName2 + ",\n");
	}
	write(2, "}\n");
	write(1, "}\n");
	// end coreParse() function, begin coreEncode() function
	write(1, "static coreEncode(a: __" + chkCommand.commandProgName + "): Command {\n");
	if chkCommand.argExists {
		if chkCommand.argModifier == ModifierRequired && true {
			write(2, "const args = [a.args];\n");
		} else if chkCommand.argModifier == ModifierOptional {
			write(2, "const args = a.args !== undefined ? [a.args] : [];\n")
		} else {
			write(2, "const args = a.args;\n")
		}
	} else {
		write(2, "const args = [] as string[];\n");
	}
	write(2, "const options = [] as string[];\n");
	for _, lineDef := range chkCommand.optionDefs {
		if lineDef.isReserved {
			continue;
		}
		fieldName := hfcgTsBannedWordsMangle(hfNormalizedToCamel(lineDef.progName));
		if lineDef.baseType == BaseTypeFlag {
			write(2, "if (a." + fieldName + `) { options.push("` + lineDef.shortName + `", ""` + "); }\n")
		} else if lineDef.modifier == ModifierOptional {
			write(2, "if (a." + fieldName + ` !== undefined) { options.push("` + lineDef.shortName + `", a.` + fieldName + "); }\n")
		} else if lineDef.modifier == ModifierRepeated {
			write(2, `for (const b of a.` + fieldName +`) { options.push("` + lineDef.shortName + `", b); }` + "\n")
		} else {		
			write(2, `options.push("` + lineDef.shortName + `", a.` + fieldName + ");\n")
		}
	}
	write(2, "return {\n");
	write(3, `command: "` + chkCommand.commandName + `",` + "\n");
	write(3, "args: args,\n");
	write(3, "options: options,\n");
	write(2, "}\n");
	write(1, "}\n");
	write(1, "static write(a: __" + chkCommand.commandProgName + "): string {\n");
	write(2, "return CcCore.encode(this.coreEncode(a));\n");
	write(1, "}\n");
	write(1, "static parse(s: string, checkReq = true): __" + chkCommand.commandProgName + "|Error {\n");
	write(2, "const e = CcCore.parse(s);\n");
	write(2, "if (e instanceof Error) { return e; }\n");
	write(2, "if (e.length !== 1) { return new Error(\"expected exactly 1 command line\"); }\n");
	write(2, "const f = e[0]!;\n");
	write(2, "return this.coreParse(f, checkReq);\n");
	write(1, "}\n");
	write(0, "}\n\n");
	return b.String();
}

// golang (.go)

func cgCommandGolang(chkCommand ChkCommandDef, indent string) string {
	return "???";
}

// ocaml (.ml)

// rust (.rs)

// java (.java)

// python (.py)

// csharp (.cs)

// fsharp (.fs)

// 

//////// codegen ends ////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// helper functions ////////////////////////////////////////////////////////

func hfIsIdentLike(r byte) bool {
	return r == '-' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_';
}

func hfIsLegalType(baseTypeName string, modifier Modifier, isArg bool) (BaseType, Modifier, error) {
	switch baseTypeName {
	case "string":
		return BaseTypeString, modifier, nil;
	case "udecimal":
		return BaseTypeUdecimal, modifier, nil;
	case "decimal":
		return BaseTypeDecimal, modifier, nil;
	case "flag":
		if isArg {
			return BaseTypeString, modifier, fmt.Errorf("base type flag is not allowed as an args type");
		}
		if modifier != ModifierRequired {
			return BaseTypeString, modifier, fmt.Errorf("type modifier is not allowed on base type flag");
		}
		return BaseTypeFlag, ModifierOptional, nil;
	case "base64":
		return BaseTypeBase64, modifier, nil;
	default:
		return BaseTypeString, modifier, fmt.Errorf("unrecognized base type %s", baseTypeName);
	}
}

func isAsciiLetterUpper(r byte) bool {
	return r >= 'A' && r <= 'Z';
}

func isAsciiLetterLower(r byte) bool {
	return r >= 'a' && r <= 'z';
}

func isAsciiLetter(r byte) bool {
	return isAsciiLetterLower(r) || isAsciiLetterUpper(r);
}

func isAscii0To9(r byte) bool {
	return r >= '0' && r <= '9';
}

func hfCommandNameCheck(commandName string) error {
	if len(commandName) < 1 {
		return fmt.Errorf("command name cannot be empty");
	}
	if pass, b := hfCheckOptionAndCommandNameCharset(commandName); !pass {
		return fmt.Errorf("command name %s contains disallowed character %c", commandName, b);
	}
	if !isAsciiLetter(commandName[0]) {
		return fmt.Errorf("command name %s must start with a letter when no explicit identifier is provided", 
		commandName);
	}
	return nil;
}

func hfCommandProgramNameCheck(commandProgramName string) error {
	if len(commandProgramName) < 1 {
		return fmt.Errorf("command program name cannot be empty");
	}
	if pass, b := hfCheckCommandProgramNameCharset(commandProgramName); !pass {
		return fmt.Errorf("command program name %s contains disallowed character %c", commandProgramName, b);
	}
	if !isAsciiLetterUpper(commandProgramName[0]) {
		return fmt.Errorf("command program name %s did not start with an upper case letter", commandProgramName);
	}
	return nil;
}

func hfOptionNameCheck(optionName string, isLong bool) error {
	if len(optionName) < 1 {
		return fmt.Errorf("option name cannot be empty");
	}
	if pass, b := hfCheckOptionAndCommandNameCharset(optionName); !pass {
		return fmt.Errorf("option name %s contains disallowed character %c", optionName, b);
	}
	if optionName[0] != '-' {
		return fmt.Errorf("option name %s did not start with a dash", optionName);
	}
	if optionName == "--" {
		return fmt.Errorf("option name cannot be `--`");
	}
	if isLong {
		i := 0;
		for i < len(optionName) && optionName[i] == '-' {
			i += 1;
		}
		if i >= len(optionName) {
			return fmt.Errorf("long option name must contain a letter");
		} 
		if !isAsciiLetter(optionName[i]) {
			return fmt.Errorf("long option name's first non-dash character must be a letter");
		}
	}
	return nil;
}

func hfNormalizeIdentLike(identLike string) []string {
	coll := [][]byte{};
	cur := []byte{};
	for _, c := range []byte(identLike) {
		if (c == '_' || c == '-') && len(cur) > 0 {
			coll = append(coll, cur);
			cur = []byte{};
		} else if (c == '_' || c == '-') {
			// do nothing
		} else if isAsciiLetterUpper(c) && len(cur) > 0 {
			coll = append(coll, cur);
			cur = []byte{ c };
		} else if isAsciiLetterUpper(c) {
			cur = append(cur, c);
		} else {
			cur = append(cur, c);
		}
	}
	if len(cur) > 0 {
		coll = append(coll, cur);
	}
	ret := []string{};
	for _, cs := range coll {
		ret = append(ret, strings.ToLower(string(cs)));
	}
	return ret;
}

func hfNormalizedToPascal(li []string) string {
	li2 := []string{};
	for _, e := range li {
		runes := []rune(e);
		runes[0] = unicode.ToUpper(runes[0]);
		li2 = append(li2, string(runes));
	}
	return strings.Join(li2, "");
}

func hfNormalizedToCamel(li []string) string {
	li2 := []string{};
	for i, e := range li {
		runes := []rune(e);
		if i > 0 {
			runes[0] = unicode.ToUpper(runes[0]);
		}
		li2 = append(li2, string(runes));
	}
	return strings.Join(li2, "");
}

func hfNormalizedToSnake(li []string) string {
	return strings.Join(li, "_");
}

func hfCheckOptionAndCommandNameCharset(optionName string) (bool, byte) {
	for _, c := range []byte(optionName) {
		if !(isAsciiLetter(c) || isAscii0To9(c) || c == '-' || c == '_') {
			return false, c;
		}
	}
	return true, '0';  // '0' does not matter
}

func hfCheckCommandProgramNameCharset(optionName string) (bool, byte) {
	for _, c := range []byte(optionName) {
		if !(isAsciiLetter(c) || isAscii0To9(c)) {
			return false, c;
		}
	}
	return true, '0';  // '0' does not matter
}

func hfcgTsType(baseType BaseType, modifier Modifier) string {
	var t string;
	var m string;
	switch modifier {
	case ModifierOptional: 
		m = "|undefined";
	case ModifierRequired:
		m = "";
	case ModifierRepeated:
		m = "[]";
	}
	switch baseType {
	case BaseTypeUdecimal, BaseTypeDecimal, BaseTypeString, BaseTypeBase64:
		t = "string" + m;
	case BaseTypeFlag:
		t = "boolean";
	default: 
		t = "string" + m;
	}
	return t;
}

func hfcgTsTypeInitExpr(baseType BaseType, modifier Modifier) string {
	var zeroVal, typeName string;
	switch baseType {
	case BaseTypeUdecimal, BaseTypeDecimal, BaseTypeString, BaseTypeBase64:
		zeroVal = "\"\"";
		typeName = "string";
	case BaseTypeFlag:
		return "false";
	default: 
		zeroVal = "\"\"";
		typeName = "string";
	}
	switch modifier {
	case ModifierOptional: 
		return "undefined as " + typeName + "|undefined";
	case ModifierRequired:
		return zeroVal;
	case ModifierRepeated:
		return "[] as " + typeName + "[]";
	default: 
		return zeroVal;
	}
}

func hfcgTsFieldUpdate(lineDef ChkOptionDef) string {
	progName := "b" + hfNormalizedToPascal(lineDef.progName);
	if lineDef.baseType == BaseTypeFlag {
		return progName + " = true;\n"
	} else if lineDef.modifier == ModifierRepeated {
		return progName + ".push(pc.options[i+1]!);\n";
	} else {
		return progName + " = pc.options[i+1]!;\n";
	}
}

func hfcgTsBannedWordsMangle(ident string) string {
	switch ident {
	case "break", "case", "catch", "class", "const", "continue", "debugger", "default", "delete", "do", "else", "enum", "export", "extends",  "false", "finally", "for", "function",  "if", "import", "in", "instanceof",  "new", "null",  "return",  "super", "switch",  "this", "throw", "true", "try", "typeof", "var", "void", "while", "with", "yield":
		return ident + "_";
	}
	return ident;
}

func hfcgGoBannedWordsMangle(ident string) string {
	switch ident {
	case "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "type", "var":
		return ident + "_";
	}
	return ident;
}

func hfcgMlBannedWordsMangle(ident string) string {
	switch ident {
	case "and", "as", "assert", "asr", "begin",  "class", "constraint", "do", "done", "downto", "else", "end", "exception", "external", "false", "for", "fun", "function", "functor", "if", "in", "include", "inherit", "initializer", "land", "lazy", "let", "lor", "lsl", "lsr", "lxor", "match", "method", "mod", "module", "mutable", "new", "nonrec", "object", "of", "open", "or", "private", "rec", "sig", "struct", "then", "to", "true", "try", "type", "val", "virtual", "when", "while", "with":
		return ident + "_";
	}
	return ident;
}

func hfcgRsBannedWordsMangle(ident string) string {
	switch ident {
	case "abstract", "as", "async", "await", "become", "box", "break", "const", "continue", "crate", "do", "dyn", "else", "enum", "extern", "false", "final", "fn", "for", "gen", "if", "impl", "in", "let", "loop", "macro", "match", "mod", "move", "mut", "override", "priv", "pub", "ref", "return",  "Self", "self", "static", "struct", "super", "trait", "true", "try", "type", "typeof", "unsafe", "unsized", "use", "virtual", "where", "while", "yield":
		return "r#" + ident;
	}
	return ident;
}

func hfcgJavaBannedWordsMangle(ident string) string {
	switch ident {
	case "abstract", "assert", "boolean", "break", "byte", "case", "catch", "char", "class", "const", "continue", "default", "do", "double", "else", "enum", "extends", "final", "finally", "float", "for", "goto", "if", "implements", "import", "instanceof", "int", "interface", "long", "native", "new", "package", "private", "protected", "public", "return", "short", "static", "strictfp", "super", "switch", "synchronized", "this", "throw", "throws", "transient", "try", "void", "volatile", "while", "true", "false", "null":
		return ident + "_";
	}
	return ident;
}

type ProgramCliParameters struct {
	tsOut string
	mlOut string
	goOut string
	goPackageName string
	verb string
	name string
	rest []string
}

func hfcliParseArgs(args []string) ProgramCliParameters {
	programCliParameters := ProgramCliParameters{};
	nameFilled, verbFilled := false, false;
	isOption := func (s string) bool { return len(s) != 0 && s[0] == '-'; }
	for i:=0; i<len(args); {
		if i>=len(args) {
			return programCliParameters;
		} else if isOption(args[i]) {
			this := args[i];
			next := "";
			if i+1<len(args) && isOption(args[i+1]) {
				i += 1;
			} else if i+1<len(args) {
				next = args[i+1];
				i += 2;
			} else {
				i += 1;
			}
			switch this {
			case "--ts-out": 
				programCliParameters.tsOut = next;
			case "--go-out":
				programCliParameters.goOut = next;
			case "--go-package-name":
				programCliParameters.goPackageName = next;
			case "--ml-out":
				programCliParameters.mlOut = next;
			}
		} else if verbFilled {
			programCliParameters.rest = append(programCliParameters.rest, args[i]);
			i += 1;
		} else if nameFilled {
			programCliParameters.verb = args[i];
			verbFilled = true;
			i += 1;
		} else {
			programCliParameters.name = args[i];
			nameFilled = true;
			i += 1;
		}
	}
	return programCliParameters;
}

//////// helper functions end ////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// cli begins //////////////////////////////////////////////////////////////

func readFile() string {
	data, err := os.ReadFile("example.txt");
	if err != nil {
		log.Fatal(err);
	}
	return string(data);
}

func main() {

	args := hfcliParseArgs(os.Args)

	tsOut := flag.String("ts-out", "", "target filename for ts artefact");
	goOut := flag.String("go-out", "", "target filename for go artefact");
	mlOut := flag.String("file", "", "target filename for ocaml artefact");
	//javaOut := flag.String("file", "", "target filename for java artefact");
	//rsOut := flag.String("file", "", "target filename for rust artefact");
	//pyOut := flag.String("file", "", "target filename for python artefact");
	flag.Parse()

	fmt.Println("--ts-out:", *tsOut);
	fmt.Println("--go-out:", *goOut);
	fmt.Println("--ml-out:", *mlOut);

	fmt.Println("positionals:", flag.Args());

	fmt.Printf("args: %#v", args);

	// norm := hfOptionNameCheck("requireAuth", true)
	// norm2 := hfNormalizeIdentLike("--requireAuth");
	// norm3 := hfNormalizedToPascal(hfNormalizeIdentLike("--requireAuth"));
	// norm4 := hfNormalizedToCamel(hfNormalizeIdentLike("--requireAuth"));
	// norm5 := hfNormalizedToSnake(hfNormalizeIdentLike("--requireAuth"));

	// fmt.Printf("%#v, %#v, %#v, %#v, %#v\n\n\n", norm, norm2, norm3, norm4, norm5);

// 	tokens, _, _ := lexTokenizer(`
	

// command arinit {
// 	-ns --namespace: string;  // comment is now supported!
// }

// command arsync {
// 	-s --session: reserved;
// 	-i --case-id: string;
// 	-n --case-number: udecimal;
// }

// command addrev {
// 	-r --revid: udecimal;
// 	-t --rev-timestamp: string;
// 	-u --uid: udecimal?;
// 	-un --uname: flag;
// 	-s --summary: string[];
// 	-c --content: string;
// }

// command addpid {
// 	args: udecimal;
// }

// command addpc {
// 	args: string[];
// }

// command arclose {
// 	-s --session: reserved;
// 	-i --case-id: string;
// 	-n --case-number: udecimal;
// }


// 	`);

	// fmt.Printf("%#v, %#v, %#v\n\n\n", errI, err, tokens);

	//program, _, _ := rdParseProgram(tokens)
	
	// fmt.Printf("%#v, %#v, %#v\n\n\n", errI, err, program);

	//checked, _, _ := chkProgram(program);

	// fmt.Printf("%#v, %#v, %#v\n\n\n", errT, err, checked);

	//tscode := cgProgramTypescript(checked, "\t");

	//fmt.Printf("%s\n\n\n", tscode);


}