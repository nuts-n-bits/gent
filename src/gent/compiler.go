package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type TokenEnum string;
const (
	TokenKwStruct TokenEnum = "TokenKwStruct";
	TokenKwEnum TokenEnum = "TokenKwEnum";
	TokenKwImport TokenEnum = "TokenKwImport";
	TokenKwReserved TokenEnum = "TokenKwReserved";
	TokenKwAs TokenEnum = "TokenKwAs";
	TokenOpenBrace TokenEnum = "TokenOpenBrace";
	TokenCloseBrace TokenEnum = "TokenCloseBrace";
	TokenEq TokenEnum = "TokenEq";
	TokenOpenParen TokenEnum = "TokenOpenParen";
	TokenCloseParen TokenEnum = "TokenCloseParen";
	TokenColon TokenEnum = "TokenColon";
	TokenColonColon TokenEnum = "TokenColonColon";
	TokenComma TokenEnum = "TokenComma";
	TokenEof TokenEnum = "TokenEof";
	// the following 3 kinds of tokens should have an associated data
	TokenIdent TokenEnum = "TokenIdent";
	TokenStringLit TokenEnum = "TokenStringLit";
	TokenNumericLit TokenEnum = "TokenNumericLit";
);

type Token struct {
	kind TokenEnum
	data string
	posStart int
	posEnd int
	err error
}

type FieldTypeBuiltin string

const (
	FieldTypeBuiltinString FieldTypeBuiltin = "FieldTypeBuiltinString";
	FieldTypeBuiltinBytes FieldTypeBuiltin = "FieldTypeBuiltinBytes";
	FieldTypeBuiltinBool FieldTypeBuiltin = "FieldTypeBuiltinBool";
	FieldTypeBuiltinU64 FieldTypeBuiltin = "FieldTypeBuiltinU64";
	FieldTypeBuiltinI64 FieldTypeBuiltin = "FieldTypeBuiltinI64";
	FieldTypeBuiltinF64 FieldTypeBuiltin = "FieldTypeBuiltinF64";
)

type IdentKind string;
const (
	IdentKindSnake IdentKind = "IdentKindSnake";
	IdentKindPascal IdentKind = "IdentKindPascal";
	IdentKindCamel IdentKind = "IdentKindCamel";
	IdentKindConst IdentKind = "IdentKindConst";
);

// AstXxx types are raw asts before type checking. They correspond 1-1 to the raw program text

type AstProgram struct {
	imports []AstImport
	topLevelStructs []AstTopLevelStruct
	topLevelEnums []AstTopLevelEnum
}

type AstImport struct {
	fileSpecifierStringLitToken Token
	symbols []AstImportSymbolTuple
}

type AstTopLevelStruct struct {
	identifierToken Token
	definitions []AstLineDefinition
}

type AstAnonymousStruct struct {
	definitions []AstLineDefinition
}

type AstTopLevelEnum struct {
	identifierToken Token
	definitions []AstLineDefinition
}

type AstAnonymousEnum struct {
	definitions []AstLineDefinition
}

type AstLineDefinition struct {
	fieldNumberLiteralToken Token
	fieldIdentifierToken Token

	// could be nil
	fieldModifierIdentifierToken *Token

	// only one of the following 3 fields is supposed to be present: 
	// oneofFieldTypeIdentifierToken, fieldTypeAnonymousEnum, fieldTypeAnonymousStruct
	oneofFieldTypeIdentifierToken *Token
	oneofFieldTypeAnonymousEnum *AstAnonymousEnum
	oneofFieldTypeAnonymousStruct *AstAnonymousStruct	
}

type FieldModifier string
const (
	FieldModifierRequired FieldModifier = "Required"
	FieldModifierRepeated FieldModifier = "Repeated"
	FieldModifierOptional FieldModifier = "Optional"
)

type AstImportSymbolTuple struct {
	originalNameIdentifierToken Token

	// coule be nil
	nameInCurrentScopeIdentifierToken *Token
}

// CheckedXxx types have passed type checking and is a semantic representation of the program

type CheckedProgram struct {
	globalScope map[string]CheckedItem
	globalIdentToken map[string]Token
}

type CheckedItem struct {
	isImported bool
	oneofEnum *[]CheckedLineItem
	oneofStruct *[]CheckedLineItem
}

type CheckedLineItem struct {
	fieldNumber uint64
	fieldIdentifierKindOriginal IdentKind
	fieldIdentifierParsed []string
	
	oneofFieldTypeAnonymousDefined *CheckedItem
	oneofFieldTypeReference string
}

// Examples:
// PascalCase -> 
// 		(IdentTypePascal, ["pascal", "case"], nil)
// snake_case -> 
// 		(IdentTypeSnake, ["snake", "case"], nil)
// camelCase -> 
// 		(IdentTypeCamel, ["camel", "case"], nil)
// ambiguous -> 
// 		(IdentTypeCamel, ["ambiguous"], nil)
// AMBIGUOUS -> 
// 		(IdentTypePascal, ["a", "m", "b", "i", "g", "u", "o", "u", "s"], nil)
// Ada_Case -> 
// 		(IdentTypeConst, ["ada", "case"], nil)
// CONST_CASE -> 
//		(IdentTypeConst, ["const", "case"], nil)
// _preceding_underscore_not_allowed -> 
// 		(IdentTypeSnake, [], fmt.Errorf("preceding underscore not allowed!"))
func parseIdent(ident string) (IdentKind, []string, error) {
	match, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]*$", ident);
	if err != nil { 
		return IdentKindSnake, []string{}, err; 
	} else if !match {
		return IdentKindSnake, []string{}, fmt.Errorf(
			"Identifier must be strictly snake_case, camelCase, PascalCase, or CONST_CASE. " + 
			"Preceding underscore and empty identifier not allowed.",
		);
	}
	match, err = regexp.MatchString("__", ident);
	if err != nil {
		return IdentKindSnake, []string{}, err;
	} else if match {
		return IdentKindSnake, []string{}, fmt.Errorf("Consecutive underscores not allowed."); 
	}
	identRunes := []rune(ident);
	if strings.Contains(ident, "_") {
		coll := strings.Split(ident, "_");
		for i, e := range coll {
			coll[i] = strings.ToLower(e);
		}
		if unicode.IsUpper(identRunes[0]) {
			return IdentKindConst, coll, nil;
		} else {
			return IdentKindSnake, coll, nil;
		}
	} else {
		originalUpper := unicode.IsUpper(identRunes[0]);
		identRunes[0] = unicode.ToUpper(identRunes[0]);
		coll := [][]rune{};
		for _, e := range identRunes {
			if unicode.IsUpper(e) {
				coll = append(coll, []rune{});
			}
			coll[len(coll) - 1] = append(coll[len(coll) - 1], e);
		}
		coll2 := []string{};
		for _, e := range coll {
			coll2 = append(coll2, strings.ToLower(string(e)));
		}
		if originalUpper {
			return IdentKindPascal, coll2, nil;
		} else {
			return IdentKindCamel, coll2, nil;
		}
	}
}

func toCamel(li []string) string {
	for i, e := range li {
		if i != 0 {
			runes := []rune(e);
			runes[0] = unicode.ToUpper(runes[0]);
			li[i] = string(runes);
		}
	}
	return strings.Join(li, "");
}

func toPascal(li []string) string {
	for i, e := range li {
		runes := []rune(e);
		runes[0] = unicode.ToUpper(runes[0]);
		li[i] = string(runes);
	}
	return strings.Join(li, "");
}

func toSnake(li []string) string {
	return strings.Join(li, "_");
}

func tokenizer(sourceStr string) ([]Token, error) {
	i, coll, source := 0, []Token{}, []rune(sourceStr);
	for i < len(source) {
		if source[i] == ' ' || source[i] == '\n' || source[i] == '\r' || source[i] == '\t' {
			i += 1;
		} else if source[i] == '{' {
			coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenOpenBrace });
			i += 1;
		} else if source[i] == '}' { 
			coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenCloseBrace });
			i += 1;
		} else if source[i] == '(' { 
			coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenOpenParen });
			i += 1;
		} else if source[i] == ')' { 
			coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenCloseParen });
			i += 1;
		} else if source[i] == '=' { 
			coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenEq });
			i += 1;
		} else if source[i] == ',' { 
			coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenComma });
			i += 1;
		} else if source[i] == ':' { 
			if len(source) >= i && source[i+1] == ':' {
				coll = append(coll, Token{ posStart: i, posEnd: i+2, kind: TokenColonColon });
				i += 2;
			} else {
				coll = append(coll, Token{ posStart: i, posEnd: i+1, kind: TokenColon });
				i += 1;
			}
		} else if source[i] == '"' {
			strlitContent, newI, err := tokenizeConsumeStrlit(source, i+1);
			if err != nil {
				return []Token{}, err;
			}
			coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenStringLit, data: strlitContent });
			i = newI;
		} else if isAscii0To9(source[i]) {
			numlit, newI := tokenizeConsumeIdent(source, i);
			match, err := regexp.MatchString("^[0-9]+$", numlit);
			if err != nil {
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenNumericLit, data: numlit, err: err });
			} else if !match {
				coll = append(coll, Token{ 
					posStart: i, posEnd: newI, kind: TokenNumericLit, data: numlit, err: fmt.Errorf(
						"Numeric literal must only consist ascii characters 0-9",
					),
				});
			} else {
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenNumericLit, data: numlit });
			}
			i = newI;
		} else if isIdentStart(source[i]) {
			ident, newI := tokenizeConsumeIdent(source, i);
			switch ident {
			case "struct":
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenKwStruct });
			case "enum":
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenKwEnum });
			case "import":
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenKwImport });
			case "reserved":
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenKwReserved });
			case "as": 
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenKwAs });
			default:
				coll = append(coll, Token{ posStart: i, posEnd: newI, kind: TokenIdent, data: ident });
			}
			i = newI;
		} else {
			return []Token{}, fmt.Errorf("Unexpected character %c", source[i]);
		}
	}
	return coll, nil;
}

// returns the identifier (possibly a keyword) as a string, and the position of the cursor position immediately 
// following the identifier
func tokenizeConsumeIdent(source []rune, identStart int) (string, int) {
	i, coll := identStart, []rune{};
	for isIdentCont(source[i]) {
		coll = append(coll, source[i]);
		i += 1;
	}
	return string(coll), i;
}

func tokenizeConsumeStrlit(source []rune, strlitStart int) (string, int, error) {
	i, coll := strlitStart, []rune{};
	for {
		switch source[i] {
		case '"':
			return string(coll), i+1, nil;
		case '\\':
			if i+1 >= len(source) {
				return string(coll), i+1, fmt.Errorf("Unterminated string literal");
			} else if source[i+1] == '\\' {
				coll = append(coll, '\\');
				i += 2;
			} else if source[i+1] == '"' {
				coll = append(coll, '"');
				i += 2;
			} else if source[i+1] == '\t' {
				coll = append(coll, '\t');
				i += 2;
			} else if source[i+1] == '\n' {
				coll = append(coll, '\n');
				i += 2;
			} else if source[i+1] == '\r' {
				coll = append(coll, '\r');
				i += 2;
			} else {
				return string(coll), i+1, fmt.Errorf("Unexepcted escape sequence");
			}
		default:
			coll = append(coll, source[i]);
			i += 1;
		}
	}
}

func isUpperAsciiLetter(r rune) bool {
	return r >= 'A' && r <= 'Z';
}

func isLowerAsciiLetter(r rune) bool {
	return r >= 'a' && r <= 'z';
}

func isAscii0To9(r rune) bool {
	return r >= '0' && r <= '9';
}

func isIdentStart(r rune) bool {
	return isUpperAsciiLetter(r) || isLowerAsciiLetter(r);
}

func isIdentCont(r rune) bool {
	return isUpperAsciiLetter(r) || isLowerAsciiLetter(r) || isAscii0To9(r) || r == '_';
}

func main() {
	source := `
	
		import "some_text_file.gent"::{ A, B as Compiler }

		struct Arsync4 {
			1 = revision_id: string (optional)
			2 = caseid: struct {
				1 = field1: struct {
					2 = whatever: enum {
					
					}
				}
			}
		}

	`;

	tokens, err := tokenizer(source);

	fmt.Printf("ERR: (%#v), TOKENS: (%#v)\n\n\n", err, tokens);

	astprogram, err := rdConsumeProgram(tokens, 0);

	fmt.Printf("ERR: %#v, AST: %+v", err, astprogram);

}

func rdConsumeProgram(tokens []Token, cursor int) (AstProgram, error) {
	astProgram := AstProgram{};
	for {
		if cursor >= len(tokens) {
			break;
		}
		parsedTopImport, newCursor, err := rdConsumeTopLevelImport(tokens, cursor);
		if err != nil {
			break;
		}
		cursor = newCursor;
		astProgram.imports = append(astProgram.imports, parsedTopImport);
	}
	for {
		if cursor >= len(tokens) {
			break;
		}
		parsedTopStruct, newCursor, err1 := rdConsumeTopLevelStruct(tokens, cursor);
		if err1 == nil {
			astProgram.topLevelStructs = append(astProgram.topLevelStructs, parsedTopStruct);
			cursor = newCursor;
			continue;
		}
		parsedTopEnum, newCursor, err2 := rdConsumeTopLevelEnum(tokens, cursor);
		if err2 == nil {
			astProgram.topLevelEnums = append(astProgram.topLevelEnums, parsedTopEnum);
			cursor = newCursor;
			continue;
		}
		return AstProgram{}, fmt.Errorf("expected struct or enum definition (%s) (%s)", err1, err2);
	}
	return astProgram, nil;
}

func rdConsumeTopLevelStruct(tokens []Token, cursor int) (topLevelStruct AstTopLevelStruct, newCursor int, err error) {
	if tokens[cursor].kind != TokenKwStruct {
		return AstTopLevelStruct{}, 0, fmt.Errorf("expected keyword struct while consuming top level struct");
	}
	cursor += 1;
	if (tokens[cursor].kind != TokenIdent) {
		return AstTopLevelStruct{}, 0, fmt.Errorf("expected identifier while consuming top level struct");
	}
	identToken := tokens[cursor];
	cursor += 1;
	structBodyLines, newCursor, err := rdConsumeStructOrEnumBody(tokens, cursor);
	if err != nil {
		return AstTopLevelStruct{}, 0, err;
	}
	cursor = newCursor;
	return AstTopLevelStruct{
		identifierToken: identToken,
		definitions: structBodyLines,
	}, cursor, nil;	
}

func rdConsumeTopLevelEnum(tokens []Token, cursor int) (topLevelEnum AstTopLevelEnum, newCursor int, err error) {
	if tokens[cursor].kind != TokenKwEnum {
		return AstTopLevelEnum{}, 0, fmt.Errorf("expected keyword enum while consuming top level enum");
	}
	cursor += 1;
	if (tokens[cursor].kind != TokenIdent) {
		return AstTopLevelEnum{}, 0, fmt.Errorf("expected identifier while consuming top level enum");
	}
	identToken := tokens[cursor];
	cursor += 1;
	enumBodyLines, newCursor, err := rdConsumeStructOrEnumBody(tokens, cursor);
	if err != nil {
		return AstTopLevelEnum{}, 0, err;
	}
	cursor = newCursor;
	return AstTopLevelEnum{
		identifierToken: identToken,
		definitions: enumBodyLines,
	}, cursor, nil;	
}

func rdConsumeTopLevelImport(tokens []Token, cursor int) (importStmt AstImport, newCursor int, err error) {
	if tokens[cursor].kind != TokenKwImport {
		return AstImport{}, 0, fmt.Errorf("expected keyword import while consuming import statement");
	}
	cursor += 1;
	if tokens[cursor].kind != TokenStringLit {
		return AstImport{}, 0, fmt.Errorf("expected string literal while consuming import statement");
	}
	importFileSpecifierStringLiteralToken := tokens[cursor];
	cursor += 1;
	if tokens[cursor].kind != TokenColonColon {
		return AstImport{}, 0, fmt.Errorf("expected token ColonColon while consuming import statement");
	}
	cursor += 1;
	if tokens[cursor].kind != TokenOpenBrace {
		return AstImport{}, 0, fmt.Errorf("expected token OpenBrace while consuming import statement");
	}
	cursor += 1;
	coll := []AstImportSymbolTuple{}; 
	for {
		if tokens[cursor].kind != TokenIdent {
			return AstImport{}, 0, fmt.Errorf("expected token identifier while consuming import statement");
		}
		importedSymbolIdentifierToken := tokens[cursor];
		var importedSymbolInCurrentNamespaceIdentifierToken *Token;
		cursor += 1;
		if tokens[cursor].kind == TokenKwAs {
			cursor += 1;
			if tokens[cursor].kind != TokenIdent {
				return AstImport{}, 0, fmt.Errorf("expected token identifier following the as keyword");
			}
			importedSymbolInCurrentNamespaceIdentifierToken = &tokens[cursor];
			cursor += 1;
		}
		importTuple := AstImportSymbolTuple{
			originalNameIdentifierToken: importedSymbolIdentifierToken,
			nameInCurrentScopeIdentifierToken: importedSymbolInCurrentNamespaceIdentifierToken,
		}
		coll = append(coll, importTuple);
		if tokens[cursor].kind == TokenComma {
			cursor += 1;
		}
		if tokens[cursor].kind == TokenCloseBrace {
			return AstImport{
				fileSpecifierStringLitToken: importFileSpecifierStringLiteralToken,
				symbols: coll,
			}, cursor+1, nil;
		}
	}
}

func rdConsumeAnonymousStruct(tokens []Token, cursor int) (parsed AstAnonymousStruct, newCursor int, err error) {
	if tokens[cursor].kind != TokenKwStruct {
		return AstAnonymousStruct{}, 0, fmt.Errorf("expected keyword struct while parsing anonymous struct");
	}
	cursor += 1;
	bodyLines, newCursor, err := rdConsumeStructOrEnumBody(tokens, cursor);
	if err != nil {
		return AstAnonymousStruct{}, 0, err;
	}
	cursor = newCursor;
	return AstAnonymousStruct{ definitions: bodyLines }, cursor, nil;
}

func rdConsumeAnonymousEnum(tokens []Token, cursor int) (parsed AstAnonymousEnum, newCursor int, err error) {
	if tokens[cursor].kind != TokenKwEnum {
		return AstAnonymousEnum{}, 0, fmt.Errorf("expected keyword enum while parsing anonymous enum");
	}
	cursor += 1;
	bodyLines, newCursor, err := rdConsumeStructOrEnumBody(tokens, cursor);
	if err != nil {
		return AstAnonymousEnum{}, 0, err;
	}
	cursor = newCursor;
	return AstAnonymousEnum{ definitions: bodyLines }, cursor, nil;
}

func rdConsumeStructOrEnumBody(tokens []Token, cursor int) (bodyLines []AstLineDefinition, newCursor int, err error) {
	if tokens[cursor].kind != TokenOpenBrace {
		return []AstLineDefinition{}, 0, fmt.Errorf("expected token OpenBrace while parsing struct or enum body");
	}
	cursor += 1;
	coll := []AstLineDefinition{};
	for tokens[cursor].kind != TokenCloseBrace {
		if tokens[cursor].kind != TokenNumericLit {
			return []AstLineDefinition{}, 0, fmt.Errorf("expected numeric literal while parsing struct or enum body");
		}
		fieldNumberLiteralToken := tokens[cursor];
		cursor += 1;
		if tokens[cursor].kind != TokenEq {
			return []AstLineDefinition{}, 0, fmt.Errorf("expected token Eq while parsing struct or enum body");
		}
		cursor += 1;
		if tokens[cursor].kind != TokenIdent {
			return []AstLineDefinition{}, 0, fmt.Errorf("expected identifier while parsing struct or enum body");
		}
		fieldIdentifierToken := tokens[cursor];
		cursor += 1;
		if tokens[cursor].kind != TokenColon {
			return []AstLineDefinition{}, 0, fmt.Errorf("expected token Colon while parsing struct or enum body");
		}
		cursor += 1;
		// start consuming type
		var fieldTypeIdentifierToken *Token;
		var fieldTypeAnonymousEnum *AstAnonymousEnum;
		var fieldTypeAnonymousStruct *AstAnonymousStruct;
		switch tokens[cursor].kind {
		case TokenIdent:
			fieldTypeIdentifierToken = &tokens[cursor];
			cursor += 1;
		case TokenKwEnum:
			anonEnum, newCursor, err := rdConsumeAnonymousEnum(tokens, cursor);
			if err != nil {
				return []AstLineDefinition{}, 0, err;
			}
			fieldTypeAnonymousEnum = &anonEnum;
			cursor = newCursor;
		case TokenKwStruct:
			anonStruct, newCursor, err := rdConsumeAnonymousStruct(tokens, cursor);
			if err != nil {
				return []AstLineDefinition{}, 0, err;
			}
			fieldTypeAnonymousStruct = &anonStruct;
			cursor = newCursor;
		default:
			return []AstLineDefinition{}, 0, 
				fmt.Errorf("expected identifier (at the type position) while parsing struct or enum body");
		}
		var fieldModifierIdentifierToken *Token;
		if tokens[cursor].kind == TokenOpenParen {
			cursor += 1;
			if tokens[cursor].kind != TokenIdent {
				return []AstLineDefinition{}, 0, fmt.Errorf("expected identifier in a field modifier");
			}
			fieldModifierIdentifierToken = &tokens[cursor];
			cursor += 1;
			if tokens[cursor].kind != TokenCloseParen {
				return []AstLineDefinition{}, 0, fmt.Errorf("expected token CloseParen while parsing field modifier");
			}
			cursor += 1;
		}
		lineDef := AstLineDefinition{
			fieldNumberLiteralToken: fieldNumberLiteralToken,
			fieldIdentifierToken: fieldIdentifierToken,
			oneofFieldTypeIdentifierToken: fieldTypeIdentifierToken,
			oneofFieldTypeAnonymousEnum: fieldTypeAnonymousEnum,
			oneofFieldTypeAnonymousStruct: fieldTypeAnonymousStruct,
			fieldModifierIdentifierToken: fieldModifierIdentifierToken,
		}
		coll = append(coll, lineDef);
	}
	return coll, cursor+1, nil;
}

func checkProgram(astProgram AstProgram) (checked CheckedProgram, errToken *Token, err error) {
	checkedProgram := CheckedProgram{};
	// ignore imports for now ......
	for _, e := range astProgram.topLevelEnums {
		_, parsedIdent, err := parseIdent(e.identifierToken.data);
		if err != nil {
			return CheckedProgram{}, &e.identifierToken, err;
		}
		pascalIdent := toPascal(parsedIdent);
		checked, errToken, err := checkStructOrEnumBody(e.definitions);
		if err != nil {
			return CheckedProgram{}, errToken, err;
		}
		if _, identAlreadyExists := checkedProgram.globalScope[pascalIdent]; identAlreadyExists {
			t := checkedProgram.globalIdentToken[pascalIdent];
			return CheckedProgram{}, &t, fmt.Errorf("duplicate identifier");
		}
		checkedItem := CheckedItem{
			isImported: false,
			oneofEnum: &checked,
		}
		checkedProgram.globalScope[pascalIdent] = checkedItem;
		checkedProgram.globalIdentToken[pascalIdent] = e.identifierToken;
	}
	for _, e := range astProgram.topLevelStructs {
		_, parsedIdent, err := parseIdent(e.identifierToken.data);
		if err != nil {
			return CheckedProgram{}, &e.identifierToken, err;
		}
		pascalIdent := toPascal(parsedIdent);
		checked, errToken, err := checkStructOrEnumBody(e.definitions);
		if err != nil {
			return CheckedProgram{}, errToken, err;
		}
		if _, identAlreadyExists := checkedProgram.globalScope[pascalIdent]; identAlreadyExists {
			t := checkedProgram.globalIdentToken[pascalIdent];
			return CheckedProgram{}, &t, fmt.Errorf("duplicate identifier");
		}
		checkedItem := CheckedItem{
			isImported: false,
			oneofStruct: &checked,
		}
		checkedProgram.globalScope[pascalIdent] = checkedItem;
		checkedProgram.globalIdentToken[pascalIdent] = e.identifierToken;
	}
	return checkedProgram, nil, nil;
}

func checkStructOrEnumBody(lineDefinitions []AstLineDefinition) (
	checked []CheckedLineItem, 
	errToken *Token,  // might be nil!
	err error,
) {
	checkedBody := make([]CheckedLineItem, 0);
	seenFieldNumbers := map[uint64]bool{};
	seenFieldIdentSnakeNormalized := map[string]bool{};
	for _, e := range lineDefinitions {
		checkedLineItem := CheckedLineItem{};
		// check field id
		fieldId, err := strconv.ParseUint(e.fieldNumberLiteralToken.data, 10, 64);
		if err != nil {
			e.fieldNumberLiteralToken.err = err;
			return make([]CheckedLineItem, 0), &e.fieldNumberLiteralToken, err;
		}
		if seenFieldNumbers[fieldId] {
			return make([]CheckedLineItem, 0), &e.fieldNumberLiteralToken, 
				fmt.Errorf("Duplicate field number %d", fieldId);
		} 
		seenFieldNumbers[fieldId] = true;
		checkedLineItem.fieldNumber = fieldId;
		// check field identifier
		identKind, identAsList, err := parseIdent(e.fieldIdentifierToken.data);
		if err != nil {
			e.fieldIdentifierToken.err = err;
			return make([]CheckedLineItem, 0), &e.fieldIdentifierToken, err;
		}
		identAsSnake := toSnake(identAsList);
		if seenFieldIdentSnakeNormalized[identAsSnake] {
			return make([]CheckedLineItem, 0), &e.fieldIdentifierToken, 
				fmt.Errorf("Duplicate field name: %s", identAsSnake);
		}
		seenFieldIdentSnakeNormalized[identAsSnake] = true;
		checkedLineItem.fieldIdentifierKindOriginal = identKind;
		checkedLineItem.fieldIdentifierParsed = identAsList;
		// check field type
		if e.oneofFieldTypeAnonymousEnum != nil {
			checkedEnumBody, errToken, err := checkStructOrEnumBody(e.oneofFieldTypeAnonymousEnum.definitions);
			if err != nil {
				return make([]CheckedLineItem, 0), errToken, err;
			}
			checkedItem := CheckedItem{
				isImported: false,
				oneofEnum: &checkedEnumBody,
			}
			checkedLineItem.oneofFieldTypeAnonymousDefined = &checkedItem;
		} else if e.oneofFieldTypeAnonymousStruct != nil {
			checkedStructBody, errToken, err := checkStructOrEnumBody(e.oneofFieldTypeAnonymousEnum.definitions);
			if err != nil {
				return make([]CheckedLineItem, 0), errToken, err;
			}
			checkedItem := CheckedItem{
				isImported: false,
				oneofStruct: &checkedStructBody,
			}
			checkedLineItem.oneofFieldTypeAnonymousDefined = &checkedItem;
		} else if e.oneofFieldTypeIdentifierToken != nil {
			checkedLineItem.oneofFieldTypeReference = e.oneofFieldTypeIdentifierToken.data;
		} else {
			return make([]CheckedLineItem, 0), nil, 
				fmt.Errorf("internal error: all oneofXXX fields are nil");
		}
		// field number ready, field identifier ready, field type ready
		checkedBody = append(checkedBody, checkedLineItem);
	}
	return checkedBody, nil, nil;
}