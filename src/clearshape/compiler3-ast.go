package main

import (
	"fmt"
)

type AstProgram struct {
	Imports  []AstImport  `json:"imports"`
	Typedefs []AstTypedef `json:"typedefs"`
}

type AstImport struct {
	ImportSrcLocationString Token `json:"importSrcLocationString"`
	ImportedAsIdent         Token `json:"importedAsIdent"`
}

type AstTypedef struct {
	Ident    Token       `json:"ident"`
	TypeExpr AstTypeExpr `json:"typeExpr"`
}

type AstTypeExpr struct {
	OneofDotAccessExpr *AstDotAccessExpr      `json:"dotAccessExpr,omitempty"`
	OneofEnumDef       *[]AstStructOrEnumLine `json:"enumDef,omitempty"`
	OneofListOf        *AstTypeExpr           `json:"listOf,omitempty"`
	OneofMapOf         *AstTypeExpr           `json:"mapOf,omitempty"`
	OneofStructDef     *[]AstStructOrEnumLine `json:"structDef,omitempty"`
	OneofTypeIdent     *Token                 `json:"typeIdent,omitempty"`
	OneofTupleDef      *[]AstTypeExpr         `json:"tupleDef,omitempty"`
}

type AstDotAccessExpr struct {
	LhsIdent Token `json:"lhsIdent"`
	RhsIdent Token `json:"rhsIdent"`
}

type AstStructOrEnumLine struct {
	WireName   *Token      `json:"wireName"` // can be omitted
	ProgName   Token       `json:"progName"`
	TypeExpr   AstTypeExpr `json:"typeExpr"`
	Omittable  bool        `json:"omittable"`  // ignored if its a enum line
	IsReserved bool        `json:"isReserved"` // if true, then typeExpr and omittable do not matter
}

func rdParseProgram(tokens []Token) (AstProgram, int, error) {
	program, i := AstProgram{}, 0
	for {
		if hfIsKw(tokens[i], "import") {
			importStmt, newI, err := rdParseImport(tokens, i)
			if err != nil {
				return AstProgram{}, newI, err
			}
			i = newI
			program.Imports = append(program.Imports, importStmt)
		} else if hfIsKw(tokens[i], "type") {
			typeDef, newI, err := rdParseTypedef(tokens, i)
			if err != nil {
				return AstProgram{}, newI, err
			}
			i = newI
			program.Typedefs = append(program.Typedefs, typeDef)
		} else if tokens[i].Kind == TokenEof {
			return program, i, nil
		} else {
			return AstProgram{}, i, fmt.Errorf("unexpected token %s(%s)", tokens[i].Kind, tokens[i].Data)
		}
	}
}

func rdParseTypedef(tokens []Token, i int) (AstTypedef, int, error) {
	astTypedef := AstTypedef{}
	// consume keyword type or error
	if hfIsKw(tokens[i], "type") {
		i += 1
	} else {
		return AstTypedef{}, i, fmt.Errorf("expected keyword type")
	}
	// consume type identifier or error
	if tokens[i].Kind == TokenIdent {
		astTypedef.Ident = tokens[i]
		i += 1
	} else {
		return AstTypedef{}, i, fmt.Errorf("expected ident")
	}
	// consume `=` or error
	if tokens[i].Kind == TokenEq {
		i += 1
	} else {
		return AstTypedef{}, i, fmt.Errorf("expected equal sign")
	}
	// consume type expression or error
	if typeExpr, newI, err := rdParseTypeExpression(tokens, i); err != nil {
		return AstTypedef{}, newI, err
	} else {
		astTypedef.TypeExpr = typeExpr
		i = newI
	}
	return astTypedef, i, nil
}

func rdParseImport(tokens []Token, i int) (AstImport, int, error) {
	astImport := AstImport{}
	// consume keyword import
	if hfIsKw(tokens[i], "import") {
		i += 1
	} else {
		return AstImport{}, i, fmt.Errorf("expected `import` keyword, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	// consume path string
	if tokens[i].Kind == TokenString {
		astImport.ImportSrcLocationString = tokens[i]
		i += 1
	} else {
		return AstImport{}, i, fmt.Errorf("expected source location string, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	// consume keyword as
	if hfIsKw(tokens[i], "as") {
		i += 1
	} else {
		return AstImport{}, i, fmt.Errorf("expected `as` keyword, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	// consume local identifier
	if tokens[i].Kind == TokenIdent {
		astImport.ImportedAsIdent = tokens[i]
		i += 1
	} else {
		return AstImport{}, i, fmt.Errorf("expected identifier, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	return astImport, i, nil
}

func rdParseTypeExpression(tokens []Token, i int) (AstTypeExpr, int, error) {
	astTypeExpr := AstTypeExpr{}
	if tokens[i].Kind == TokenOpenParen { // (  /* TypeExpr */  )
		astTypeExprInner, newI, err := rdParseTypeExpression(tokens, i+1)
		if err != nil {
			return AstTypeExpr{}, newI, err
		}
		i = newI
		// `)` -> return, other -> error
		if tokens[i].Kind == TokenCloseParen {
			astTypeExpr = astTypeExprInner
			i += 1
		} else {
			return AstTypeExpr{}, i, fmt.Errorf("expected CloseParen, got %s(%s)", tokens[i].Kind, tokens[i].Data)
		}
	} else if tokens[i].Kind == TokenOpenBrace { // { struct }
		structLines, newI, err := rdParseTypeExprStruct(tokens, i, true)
		if err != nil {
			return AstTypeExpr{}, newI, err
		}
		i = newI
		astTypeExpr.OneofStructDef = &structLines
	} else if hfIsKw(tokens[i], "enum") { // enum { enum }
		structLines, newI, err := rdParseTypeExprEnum(tokens, i)
		if err != nil {
			return AstTypeExpr{}, newI, err
		}
		i = newI
		astTypeExpr.OneofEnumDef = &structLines
	} else if tokens[i].Kind == TokenIdent && tokens[i+1].Kind == TokenDot && tokens[i+2].Kind == TokenIdent { // field.access
		astTypeExpr.OneofDotAccessExpr = &AstDotAccessExpr{
			LhsIdent: tokens[i],
			RhsIdent: tokens[i+2],
		}
		i += 3
	} else if tokens[i].Kind == TokenOpenBracket { // [tuple, tuple]
		tupleTypes, newI, err := rdParseTypeExprTuple(tokens, i)
		if err != nil {
			return AstTypeExpr{}, newI, err
		}
		astTypeExpr.OneofTupleDef = &tupleTypes
		i = newI
	} else if tokens[i].Kind == TokenIdent { // typeident
		astTypeExpr.OneofTypeIdent = &tokens[i]
		i += 1
	} else {
		return AstTypeExpr{}, i, fmt.Errorf(
			"expected one of struct/tuple definition or type identifier, found %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	for tokens[i].Kind == TokenOpenBracket && tokens[i+1].Kind == TokenCloseBracket { // XXX[]
		t := astTypeExpr
		astTypeExpr = AstTypeExpr{OneofListOf: &t}
		i += 2
	}
	return astTypeExpr, i, nil
}

func rdParseTypeExprEnum(tokens []Token, i int) ([]AstStructOrEnumLine, int, error) {
	// consume `enum`
	if hfIsKw(tokens[i], "enum") {
		i += 1
	} else {
		return []AstStructOrEnumLine{}, i, fmt.Errorf("expected keyword enum, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	// consume as if its a regular struct
	lines, newI, err := rdParseTypeExprStruct(tokens, i, false)
	return lines, newI, err
}

// set allowQuestionMark to true when parsing struct { ... } and false when parsing enum { ... }
func rdParseTypeExprStruct(tokens []Token, i int, allowQuestionMark bool) ([]AstStructOrEnumLine, int, error) {
	// consume `{`
	if tokens[i].Kind == TokenOpenBrace {
		i += 1
	} else {
		return []AstStructOrEnumLine{}, i, fmt.Errorf("expected OpenBrace, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	// end early if `}`
	if tokens[i].Kind == TokenCloseBrace {
		return []AstStructOrEnumLine{}, i + 1, nil
	}
	structLines := []AstStructOrEnumLine{}
	for {
		line := AstStructOrEnumLine{}
		// consume field names
		if tokens[i].Kind == TokenIdent && tokens[i+1].Kind == TokenIdent { // wirename progname both idents
			line.WireName = &tokens[i]
			line.ProgName = tokens[i+1]
			i += 2
		} else if tokens[i].Kind == TokenString && tokens[i+1].Kind == TokenIdent { // "wirename" progname
			line.WireName = &tokens[i]
			line.ProgName = tokens[i+1]
			i += 2
		} else if tokens[i].Kind == TokenIdent { // no wire name, only progname
			line.WireName = &tokens[i]
			line.ProgName = tokens[i]
			i += 1
		} else {
			return []AstStructOrEnumLine{}, i, fmt.Errorf("expected field name, got %s(%s)", tokens[i].Kind, tokens[i].Data)
		}
		// consume `?` optionally
		if tokens[i].Kind == TokenQuestion {
			if allowQuestionMark {
				line.Omittable = true
				i += 1
			} else {
				return []AstStructOrEnumLine{}, i, fmt.Errorf("Question mark not allowed inside enum")
			}
		}
		// consume `:`
		if tokens[i].Kind == TokenColon {
			i += 1
		} else {
			return []AstStructOrEnumLine{}, i, fmt.Errorf("expected Colon, got %s(%s)", tokens[i].Kind, tokens[i].Data)
		}
		// consume typeExpr
		if typeExpr, newI, err := rdParseTypeExpression(tokens, i); err != nil {
			return []AstStructOrEnumLine{}, newI, err
		} else {
			line.TypeExpr = typeExpr
			i = newI
		}
		// append
		structLines = append(structLines, line)
		// `}` -> end, `,}` -> end, `,` -> nextline, other -> error
		if tokens[i].Kind == TokenCloseBrace {
			return structLines, i + 1, nil
		} else if tokens[i].Kind == TokenComma && tokens[i+1].Kind == TokenCloseBrace {
			return structLines, i + 2, nil
		} else if tokens[i].Kind == TokenComma {
			i += 1
			continue
		} else {
			return []AstStructOrEnumLine{}, i, fmt.Errorf(
				"expected `,` or `,}` or `}` at the end of a struct line, got %s(%s)", tokens[i].Kind, tokens[i].Data)
		}
	}
}

func rdParseTypeExprTuple(tokens []Token, i int) ([]AstTypeExpr, int, error) {
	// consume `[`
	if tokens[i].Kind == TokenOpenBracket {
		i += 1
	} else {
		return []AstTypeExpr{}, i, fmt.Errorf("expected OpenBracket, got %s(%s)", tokens[i].Kind, tokens[i].Data)
	}
	// end early if `]`
	if tokens[i].Kind == TokenCloseBrace {
		return []AstTypeExpr{}, i + 1, nil
	}
	types := []AstTypeExpr{}
	for {
		// consume typeExpr
		if typeExpr, newI, err := rdParseTypeExpression(tokens, i); err != nil {
			return []AstTypeExpr{}, newI, err
		} else {
			types = append(types, typeExpr)
			i = newI
		}
		// `]` -> end, `,]` => end, `,` -> nextitem, other -> error
		if tokens[i].Kind == TokenComma && tokens[i+1].Kind == TokenCloseBracket {
			return types, i + 2, nil
		} else if tokens[i].Kind == TokenCloseBracket {
			return types, i + 1, nil
		} else if tokens[i].Kind == TokenComma {
			i += 1
			continue
		} else {
			return []AstTypeExpr{}, i, fmt.Errorf(
				"expected `,` or `,]` or `]` after a tuple item, got %s(%s)", tokens[i].Kind, tokens[i].Data)
		}
	}
}
