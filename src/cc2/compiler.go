package main

import (
	_ "embed"
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
	TokenKwReserved TokenKind= "TokenKwReserved"
	TokenKwAs TokenKind = "TokenKwAs"
	TokenKwMix TokenKind = "TokenKwMix"
	TokenOpenBrace TokenKind = "TokenOpenBrace"
	TokenCloseBrace TokenKind = "TokenCloseBrace"
	TokenOpenParen TokenKind = "TokenOpenParen"
	TokenCloseParen TokenKind = "TokenCloseParen"
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
		} else if program[i] == '(' {
			tokens = append(tokens, Token{ kind: TokenOpenParen, start: i, end: i+1 });
			i += 1;
		} else if program[i] == ')' {
			tokens = append(tokens, Token{ kind: TokenCloseParen, start: i, end: i+1 });
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
			case "reserved":
				tokens = append(tokens, Token{ kind: TokenKwReserved, start: i, end: i2 });
			case "as":
				tokens = append(tokens, Token{ kind: TokenKwAs, start: i, end: i2 });
			case "mix":
				tokens = append(tokens, Token{ kind: TokenKwMix, start: i, end: i2 });
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
	mixDefs []AstMixBlock
}

type AstCommandBlock struct {
	commandName Token
	argTypeBase *Token
	argModifier Modifier
	lineDefs []AstLineDef
}

type AstLineDef struct {
	shortName *Token  // no short name means the line defines argument type
	longName *Token  // no long name means fallback to short name
	typeBase *Token  // no type base means it's a `resevered` line
	typeMod Modifier  // reserved line makes modifier meaningless
	reservedTok *Token
}

type AstMixBlock struct {
	mixName *Token
	lineDefs []AstMixLineDef
}

type AstMixLineDef struct {
	commandName Token
	modifier Modifier
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
		case TokenKwMix:
			mixBlock, newI, err := rdParseMixBlock(tokens, i);
			if err != nil {
				return AstProgram{}, newI, err;
			}
			program.mixDefs = append(program.mixDefs, mixBlock);
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
	// consume `(`, or error
	if tokens[i].kind == TokenOpenParen {
		i += 1;
	} else {		
		return AstCommandBlock{}, i, fmt.Errorf("expected TokenOpenParen while parsing command block");
	}
	// consume type expression or error
	if tokens[i].kind != TokenCloseParen {	
		typeBaseTok, mod, newI, err := rdConsumeTypeExpression(tokens, i);
		if err != nil {
			return AstCommandBlock{}, newI, err;
		}
		commandDef.argTypeBase = typeBaseTok;
		commandDef.argModifier = mod;
		i = newI;
	}
	// consume `)`, or error
	if tokens[i].kind == TokenCloseParen {
		i += 1;
	} else {		
		return AstCommandBlock{}, i, fmt.Errorf("expected TokenCloseParen while parsing command block");
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
		// consume LHS of line before `:`, could be `-s`, or `-s --long-name`, or error.
		if tokens[i].kind == TokenIdentLike && tokens[i+1].kind == TokenIdentLike {
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
			return AstCommandBlock{}, i, fmt.Errorf("expected 1-2 option names at " + 
			"the beginning of a line definition in a command block, found %s", tokens[i].kind);
		}
		// consume `:`
		if tokens[i].kind == TokenColon {
			i += 1;
		} else {
			return AstCommandBlock{}, i, fmt.Errorf("expected `:` after option names or `args`, found %s",
			tokens[i].kind);
		}
		// consume RHS of line before `;`, could be `reserved` or TypeExpression (`type`, `type[]`, `type?`)
		if tokens[i].kind == TokenKwReserved {
			t := tokens[i];
			lineDef.reservedTok = &t;
			i += 1;
		} else {
			typeBaseTok, mod, newI, err := rdConsumeTypeExpression(tokens, i);
			if err != nil {
				return AstCommandBlock{}, newI, err;
			}
			lineDef.typeBase = typeBaseTok;
			lineDef.typeMod = mod;
			i = newI;
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

func rdParseMixBlock(toks []Token, i int) (AstMixBlock, int, error) {
	astMixBlock := AstMixBlock{};
	// consume mix keyword
	if toks[i].kind != TokenKwMix {
		return AstMixBlock{}, i, fmt.Errorf("expected keyword `mix` while consuming mix block");
	}
	i += 1;
	if toks[i].kind != TokenIdentLike {
		return AstMixBlock{}, i, fmt.Errorf("expected ident like while consuming mix block"); 
	}
	astMixBlock.mixName = &toks[i];
	i += 1;
	if toks[i].kind != TokenOpenBrace {
		return AstMixBlock{}, i, fmt.Errorf("expexted OpenBrace while consuming mix block");
	}
	i += 1;
	for {
		if toks[i].kind == TokenCloseBrace && true {
			i += 1;
			return astMixBlock, i, nil;
		} else if toks[i].kind == TokenIdentLike {
			astMixLineDef := AstMixLineDef{};
			astMixLineDef.commandName = toks[i];
			i += 1;
			if toks[i].kind == TokenQuestion {
				astMixLineDef.modifier = ModifierOptional;
				i += 1;
			} else if toks[i].kind == TokenOpenBracket && toks[i+1].kind == TokenCloseBracket {
				astMixLineDef.modifier = ModifierRepeated;
				i += 2;
			} else {
				astMixLineDef.modifier = ModifierRequired;
			}
			astMixBlock.lineDefs = append(astMixBlock.lineDefs, astMixLineDef);
			if toks[i].kind == TokenSemicolon {
				i += 1;
			}
			continue;
		} else {
			return AstMixBlock{}, i, fmt.Errorf("expected IdentLike or CloseBrace while parsing inside a mix block");
		}
	}
}

func rdConsumeTypeExpression(toks []Token, i int) (typeBase *Token, mod Modifier, newI int, err error) {
	if toks[i].kind == TokenIdentLike && toks[i+1].kind == TokenOpenBracket && toks[i+2].kind == TokenCloseBracket {
		t := toks[i];
		return &t, ModifierRepeated, i+3, nil;
	} else if toks[i].kind == TokenIdentLike && toks[i+1].kind == TokenQuestion {
		t := toks[i];
		return &t, ModifierOptional, i+2, nil;
	} else if toks[i].kind == TokenIdentLike {
		t := toks[i];
		return &t, ModifierRequired, i+1, nil;
	} else {
		return nil, ModifierRequired, i, fmt.Errorf("expected type expression, found %s", toks[i].kind);
	}
}

//////// parser ends /////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// checker /////////////////////////////////////////////////////////////////

type ChkProgram struct {
	commands []ChkCommandDef
	mixes []ChkMixDef
}

type ChkCommandDef struct {
	commandName string
	commandProgName []string
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

type ChkMixDef struct {
	mixName string
	mixProgName []string
	mixProgNameSrcTok Token
	mixCommands []ChkMixCommand
}

type ChkMixCommand struct {
	commandName string
	commandNameSrcTok Token
	commandProgName []string
	modifier Modifier
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
	for _, astMixBlock := range astProgram.mixDefs {
		checkedMixBlock, errTok, err := chkMixDef(astMixBlock);
		if err != nil {
			return ChkProgram{}, errTok, err;
		}
		checkedProgram.mixes = append(checkedProgram.mixes, checkedMixBlock);
	}
	// check duplicate commandProgName
	seenProgNamesPascal := map[string]bool{};
	seenCommandNames := map[string]bool{};
	for _, checkedCommandDef := range checkedProgram.commands {
		normalized := hfNormalizedToPascal(checkedCommandDef.commandProgName);
		if seenProgNamesPascal[normalized] {
			return ChkProgram{}, checkedCommandDef.commandProgNameSrcTok, 
			fmt.Errorf("duplicate command program name %s (from command: %s)", 
			normalized, checkedCommandDef.commandName);
		}
		seenProgNamesPascal[normalized] = true;
		seenCommandNames[checkedCommandDef.commandName] = true;  // commands in mix blocks must refer to known commands
	}
	for _, checkedMixDef := range checkedProgram.mixes {
		normalized := hfNormalizedToPascal(checkedMixDef.mixProgName);
		if seenProgNamesPascal[normalized] {
			return ChkProgram{}, &checkedMixDef.mixProgNameSrcTok, 
			fmt.Errorf("duplicate mix program name %s (from mix: %s)",
			normalized, checkedMixDef.mixName);
		}
		seenProgNamesPascal[normalized] = true;
		for _, mixCommand := range checkedMixDef.mixCommands {
			if !seenCommandNames[mixCommand.commandName] {
				return ChkProgram{}, &mixCommand.commandNameSrcTok, 
				fmt.Errorf("undeclared command name (%s)", mixCommand.commandName);
			}
		}
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
	checkedCommandDef.commandProgName = hfNormalizeIdentLike(astCommandBlock.commandName.data);
	checkedCommandDef.commandProgNameSrcTok = &astCommandBlock.commandName;
	// transfer arguments type and modifier
	if astCommandBlock.argTypeBase != nil {
		baseType, modifier, err := hfIsLegalType(astCommandBlock.argTypeBase.data, astCommandBlock.argModifier, true);
		if err != nil {
			return ChkCommandDef{}, astCommandBlock.argTypeBase, err;
		}
		checkedCommandDef.argExists = true;
		checkedCommandDef.argBaseType = baseType;
		checkedCommandDef.argModifier = modifier;
	}
	// convert each line inside the block
	for _, astLineDef := range astCommandBlock.lineDefs {
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
		// populate option BaseType and option Modifier if not reserved
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

func chkMixDef(astMixBlock AstMixBlock) (checked ChkMixDef, errTok *Token, err0 error) {
	checkedMixDef := ChkMixDef{};
	// populate mix name
	err := hfCommandNameCheck(astMixBlock.mixName.data);
	if err != nil {
		return ChkMixDef{}, astMixBlock.mixName, err;
	}
	checkedMixDef.mixName = astMixBlock.mixName.data;
	// pupulate prog name by transforming mix name
	checkedMixDef.mixProgName = hfNormalizeIdentLike(astMixBlock.mixName.data);
	checkedMixDef.mixProgNameSrcTok = *astMixBlock.mixName;
	// convert each line inside the block, scanning for duplicates
	seenCommandNames := map[string]bool{};
	for _, lineDef := range astMixBlock.lineDefs {
		checkedMixCommand := ChkMixCommand{};
		checkedMixCommand.commandName = lineDef.commandName.data;
		checkedMixCommand.commandNameSrcTok = lineDef.commandName;
		checkedMixCommand.commandProgName = hfNormalizeIdentLike(lineDef.commandName.data);
		checkedMixCommand.modifier = lineDef.modifier;
		if seenCommandNames[lineDef.commandName.data] {
			return ChkMixDef{}, &lineDef.commandName, 
			fmt.Errorf("duplicate command name %s", lineDef.commandName.data);
		} else {
			seenCommandNames[lineDef.commandName.data] = true;
		}
		checkedMixDef.mixCommands = append(checkedMixDef.mixCommands, checkedMixCommand);
	}
	return checkedMixDef, nil, nil;
}

//////// checker ends ////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////// codegen /////////////////////////////////////////////////////////////////

// ocaml (.ml)

// rust (.rs)

// java (.java)

// python (.py)

// csharp (.cs)

// fsharp (.fs)

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
	if pass, badByte := hfCheckOptionAndCommandNameCharset(commandName); !pass {
		return fmt.Errorf("command name %s contains disallowed character %c", commandName, badByte);
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

// This applies to all output langauges. These names are reserved by the runtime library.
func hfcgOtherBannedWordsMangle(ident string) string {
	switch ident {
	case "args", "CcCore":
		return ident + "_";
	}
	return ident;
}

func hfcgMlBannedWordsMangle(ident string) string {
	switch ident {
	case "and", "as", "assert", "asr", "begin",  "class", "constraint", "do", "done", "downto", "else", "end", "exception", "external", "false", "for", "fun", "function", "functor", "if", "in", "include", "inherit", "initializer", "land", "lazy", "let", "lor", "lsl", "lsr", "lxor", "match", "method", "mod", "module", "mutable", "new", "nonrec", "object", "of", "open", "or", "private", "rec", "sig", "struct", "then", "to", "true", "try", "type", "val", "virtual", "when", "while", "with":
		return ident + "_";
	}
	return hfcgOtherBannedWordsMangle(ident);
}

func hfcgRsBannedWordsMangle(ident string) string {
	switch ident {
	case "abstract", "as", "async", "await", "become", "box", "break", "const", "continue", "crate", "do", "dyn", "else", "enum", "extern", "false", "final", "fn", "for", "gen", "if", "impl", "in", "let", "loop", "macro", "match", "mod", "move", "mut", "override", "priv", "pub", "ref", "return",  "Self", "self", "static", "struct", "super", "trait", "true", "try", "type", "typeof", "unsafe", "unsized", "use", "virtual", "where", "while", "yield":
		return "r#" + ident;
	}
	return hfcgOtherBannedWordsMangle(ident);
}

func hfcgJavaBannedWordsMangle(ident string) string {
	switch ident {
	case "abstract", "assert", "boolean", "break", "byte", "case", "catch", "char", "class", "const", "continue", "default", "do", "double", "else", "enum", "extends", "final", "finally", "float", "for", "goto", "if", "implements", "import", "instanceof", "int", "interface", "long", "native", "new", "package", "private", "protected", "public", "return", "short", "static", "strictfp", "super", "switch", "synchronized", "this", "throw", "throws", "transient", "try", "void", "volatile", "while", "true", "false", "null":
		return ident + "_";
	}
	return hfcgOtherBannedWordsMangle(ident);
}

type ProgramCliParameters struct {
	tsOut string
	mlOut string
	goOut string
	goPackageName string
	verb string
	name string
	rest []string
	indent string
}

func hfcliParseArgs(args []string) ProgramCliParameters {
	programCliParameters := ProgramCliParameters{};
	nameFilled, verbFilled, positionalMode := false, false, false;
	isOption := func (s string) bool { return len(s) != 0 && s[0] == '-'; }
	for i:=0; i<len(args); {
		if i>=len(args) {
			return programCliParameters;
		} else if isOption(args[i]) && !positionalMode {
			if args[i] == "--" {
				positionalMode = true;
				i += 1;
				continue;
			}
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
			case "--indent":
				programCliParameters.indent = next;
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

func readFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName);
	if err != nil {
		return "", err;
	}
	return string(data), nil;
}

func writeFile(fileName string, fileContent string) error {
	return os.WriteFile(fileName, []byte(fileContent), 0666);
}

func build(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file");
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file");
	}
	programStr, err := readFile(args.rest[0]);
	if err != nil {
		log.Fatalf("ERR: %s", err.Error());
	}
	tokens, err, errI := lexTokenizer(programStr);
	if err != nil {
		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), tokens[errI].start, tokens[errI].end);
	}
	program, errI, err := rdParseProgram(tokens);
	if err != nil {
		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), tokens[errI].start, tokens[errI].end);
	}
	checked, errT, err := chkProgram(program);
	if err != nil {
		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), errT.start, errT.end);
	}
	// parsing complete
	indent := "";
	switch args.indent {
	case "", "4": 
		indent = "    ";
	case "tab":
		indent = "\t";
	case "2":
		indent = "  ";
	default:
		log.Fatal("--indent must be `4`, `2`, `tab` or unspecified");
	}
	//
	if args.tsOut != "" {
		program := cgProgramTypescript(checked, indent);
		writeFile(args.tsOut, program);
	}
	if args.goOut != "" {
		if args.goPackageName == "" {
			log.Fatal("--go-package-name must be present when --go-out is specified");
		}
		program := cgProgramGolang(checked, indent, args.goPackageName);
		writeFile(args.goOut, program);
	}
}

func show_ast(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file");
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file");
	}
	programStr, err := readFile(args.rest[0]);
	if err != nil {
		log.Fatalf("ERR: %s", err.Error());
	}
	tokens, err, errI := lexTokenizer(programStr);
	if err != nil {
		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), tokens[errI].start, tokens[errI].end);
	}
	program, errI, err := rdParseProgram(tokens);
	if err != nil {
		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), tokens[errI].start, tokens[errI].end);
	}
	fmt.Printf("====== AST ====== \n\n%#v", program);
}

func main() {

	args := hfcliParseArgs(os.Args)

	//fmt.Printf("//args: %#v", args);

	if args.verb == "build" && true {
		build(args);
	} else if args.verb == "show-ast" {
		show_ast(args);
	}
}