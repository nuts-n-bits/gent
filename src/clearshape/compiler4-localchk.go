package main

import (
	"fmt"
	"strings"
)

/*
RULES
1. No toplevel identifier conflict allowed (when compared after cononicalization)
2. No wirename conflict allowed (when compared literally)
   - If progname is used as wirename, prog name is converted to camelCase first
3. No progname conflict allowed (when compared after cononicalization - camelCase)
4. Non minted toplevel identifiers mustn't start with "csres0"
5. Imports cannot collide with top level idents, after canonicalization

*/

type BuiltinType string

const (
	BuiltinTypeString  BuiltinType = "BuiltinTypeString"
	BuiltinTypeBoolean BuiltinType = "BuiltinTypeBoolean"
	BuiltinTypeInt64   BuiltinType = "BuiltinTypeInt64"
	BuiltinTypeUint64  BuiltinType = "BuiltinTypeUint64"
	BuiltinTypeFloat64 BuiltinType = "BuiltinTypeFloat64"
	BuiltinTypeNull    BuiltinType = "BuiltinTypeNull"
	BuiltinTypeBinary  BuiltinType = "BuiltinTypeBinary"
)

func LcIsBuiltIn(s string) (BuiltinType, bool) {
	switch s {
	case "string":
		return BuiltinTypeString, true
	case "u64":
		return BuiltinTypeUint64, true
	case "i64":
		return BuiltinTypeInt64, true
	case "f64":
		return BuiltinTypeFloat64, true
	case "boolean":
		return BuiltinTypeBoolean, true
	case "null":
		return BuiltinTypeNull, true
	case "binary":
		return BuiltinTypeBinary, true
	default:
		return BuiltinTypeString, false
	}
}

type LcProgram struct {
	// Imports is keyed by ident string without normalization
	Imports      map[string]AstImport  `json:"imports"`      // these imports should not be normalized
	TopLevelDefs map[string]LcTypeExpr `json:"topLevelDefs"` // these idents should be normalized to PascalCase
}

type LcTypeExpr struct {
	OneofStruct     *[]LcStructOrEnumLine `json:"struct,omitempty"`
	OneofEnum       *[]LcStructOrEnumLine `json:"enum,omitempty"`
	OneofTuple      *[]LcTypeExpr         `json:"tuple,omitempty"`
	OneofTokenIdent *Token                `json:"tokenIdent,omitempty"`
	OneofBuiltin    *BuiltinType          `json:"builtin,omitempty"`
	OneofListof     *LcTypeExpr           `json:"listOf,omitempty"`
	OneofMapof      *LcTypeExpr           `json:"mapOf,omitempty"`
	OneofImported   *LcImported           `json:"imported,omitempty"`
}

type LcImported struct {
	ImportedIdent Token `json:"importedIdent"`
	ForeignIdent  Token `json:"foreignIdent"`
}

type LcStructOrEnumLine struct {
	WireName   string     `json:"wireName"`
	ProgName   []string   `json:"progName"`
	TypeExpr   LcTypeExpr `json:"typeExpr"`
	Omittable  bool       `json:"omittable"` // this is ignored when its an enum line
	IsReserved bool       `json:"isReserved"`
}

func lcCheckProgram1Of2CheckReservedName(astProgram AstProgram) (*Token, error) {
	for _, importStmt := range astProgram.Imports {
		if err := lcIsReservedIdent(importStmt.ImportedAsIdent.Data); err != nil {
			return &importStmt.ImportedAsIdent, err
		}
	}
	return nil, nil
}

func lcCheckProgram2Of2CheckCollisionAndUndefined(astProgram AstProgram) (
	lcProg LcProgram, topLevelCollision []LcErrorTokenCollision, undefToks []Token,
) {
	tokenCollisions := []LcErrorTokenCollision{}
	t := lcCheckTopLevelIdentCollision(astProgram.Typedefs, astProgram.Imports)
	tokenCollisions = append(tokenCollisions, t...)
	for _, astTopLevelTypedef := range astProgram.Typedefs {
		lcCheckStructOrEnumFieldNamesRecursive(astTopLevelTypedef.TypeExpr, &tokenCollisions)
	}
	lcImports := map[string]AstImport{}
	lcTlts := map[string]LcTypeExpr{}
	for _, imp := range astProgram.Imports {
		// dont care about duplicate since we have checked top level ident collision above
		lcImports[imp.ImportedAsIdent.Data] = imp
	}
	for _, tld := range astProgram.Typedefs {
		// dont care about duplicate since we have checked top level ident collision above
		pascalNorm := hfNormalizedToPascal(hfNormalizeIdent(tld.Ident.Data))
		lcTlts[pascalNorm] = lcTypeConvertNoCheck(tld.TypeExpr)
	}
	lcProgram := LcProgram{
		Imports:      lcImports,
		TopLevelDefs: lcTlts,
	}
	undefinedToks := lcCheckReferenceExist(lcProgram)
	return lcProgram, tokenCollisions, undefinedToks
}

func lcCheckStructOrEnumFieldNamesRecursive(astType AstTypeExpr, tokenCollisions *[]LcErrorTokenCollision) {
		if astType.OneofStructDef != nil {
			t := lcCheckStructOrEnumFieldNamesCore(*astType.OneofStructDef)
			*tokenCollisions = append(*tokenCollisions, t...)
			for _, subType := range *astType.OneofStructDef {
				lcCheckStructOrEnumFieldNamesRecursive(subType.TypeExpr, tokenCollisions)
			}
		} else if astType.OneofEnumDef != nil {
			t := lcCheckStructOrEnumFieldNamesCore(*astType.OneofEnumDef)
			*tokenCollisions = append(*tokenCollisions, t...)
			for _, subType := range *astType.OneofEnumDef {
				lcCheckStructOrEnumFieldNamesRecursive(subType.TypeExpr, tokenCollisions)
			}
		} else if astType.OneofTupleDef != nil {
			for _, subType := range *astType.OneofTupleDef {
				lcCheckStructOrEnumFieldNamesRecursive(subType, tokenCollisions)
			}
		} else if astType.OneofMapOf != nil {
			lcCheckStructOrEnumFieldNamesRecursive(*astType.OneofMapOf, tokenCollisions)
		} else if astType.OneofListOf != nil {
			lcCheckStructOrEnumFieldNamesRecursive(*astType.OneofListOf, tokenCollisions)
		}
}

func lcCheckStructOrEnumFieldNamesCore(lines []AstStructOrEnumLine) []LcErrorTokenCollision {
	seenProgCamel := map[string][]Token{}
	seenWire := map[string][]Token{}
	for _, line := range lines {
		// compare with camel because camel is more easily collided. want to catch collision even if false pos.
		// consider example: snake(err_400) and snake(err400) are both camel(err400)
		normalizedProg := hfNormalizeIdent(line.ProgName.Data)
		camelProg := hfNormalizedToCamel(normalizedProg)
		seenProgCamel[camelProg] = append(seenProgCamel[camelProg], line.ProgName)
		if line.WireName != nil {
			seenWire[line.WireName.Data] = append(seenWire[line.WireName.Data], *line.WireName)
		} else {
			seenWire[camelProg] = append(seenProgCamel[camelProg], line.ProgName)
		}
	}
	ret := []LcErrorTokenCollision{}
	for k, v := range seenProgCamel {
		if len(v) < 2 {
			continue
		}
		entry := LcErrorTokenCollision{
			ErrT: v,
			Err: fmt.Errorf("Multiple colliding prog names for (%s) in the same struct/enum definition. Remember "+
				"that names mustn't collide after normalization is applied to them", k),
		}
		ret = append(ret, entry)
	}
	for k, v := range seenWire {
		if len(v) < 2 {
			continue
		}
		entry := LcErrorTokenCollision{
			ErrT: v,
			Err: fmt.Errorf("Multiple colliding wire names for (%s) in the same struct/enum definition. Remember "+
				"that if a wire name is not provided, the prog name is coerced into wire name by transforming "+
				"it into camelCase", k),
		}
		ret = append(ret, entry)
	}
	return ret
}

// returns a list of tokens that is undefined references. len = 0 means no error.
func lcCheckReferenceExist(lcProgram LcProgram) []Token {
	undefinedTokens := []Token{}
	for _, tld := range lcProgram.TopLevelDefs {
		lcCheckReferenceExistCore(tld, &undefinedTokens, lcProgram)
	}
	return undefinedTokens
}

func lcCheckReferenceExistCore(tld LcTypeExpr, undefinedTokens *[]Token, lcProgram LcProgram) {
	if tld.OneofStruct != nil {
		for _, line := range *tld.OneofStruct {
			lcCheckReferenceExistCore(line.TypeExpr, undefinedTokens, lcProgram)
		}
	} else if tld.OneofEnum != nil {
		for _, line := range *tld.OneofEnum {
			lcCheckReferenceExistCore(line.TypeExpr, undefinedTokens, lcProgram)
		}
	} else if tld.OneofTuple != nil {
		for _, inner := range *tld.OneofTuple {
			lcCheckReferenceExistCore(inner, undefinedTokens, lcProgram)
		}
	} else if tld.OneofTokenIdent != nil {
		normPascal := hfNormalizedToPascal(hfNormalizeIdent(tld.OneofTokenIdent.Data))
		if _, has := lcProgram.TopLevelDefs[normPascal]; !has {
			*undefinedTokens = append(*undefinedTokens, *tld.OneofTokenIdent)
		}
	} else if tld.OneofBuiltin != nil {
		// nothing to check against
	} else if tld.OneofListof != nil {
		lcCheckReferenceExistCore(*tld.OneofListof, undefinedTokens, lcProgram)
	} else if tld.OneofMapof != nil {
		lcCheckReferenceExistCore(*tld.OneofMapof, undefinedTokens, lcProgram)
	} else if tld.OneofImported != nil {
		if _, has := lcProgram.Imports[tld.OneofImported.ImportedIdent.Data]; !has {
			*undefinedTokens = append(*undefinedTokens, tld.OneofImported.ImportedIdent)
		}
		if lcIsReservedIdent(tld.OneofImported.ForeignIdent.Data) != nil {
			*undefinedTokens = append(*undefinedTokens, tld.OneofImported.ForeignIdent)
		}
	} else {
		panic("unreachable")
	}
}

type LcErrorTokenCollision struct {
	ErrT []Token `json:"errT"` // always len >= 2
	Err  error   `json:"err"`  // always defined
}

func lcCheckTopLevelIdentCollision(astTypes []AstTypedef, imports []AstImport) []LcErrorTokenCollision {
	seenTopLevelNames := map[string][]Token{}
	for _, tlt := range astTypes {
		pascalName := hfNormalizedToPascal(hfNormalizeIdent(tlt.Ident.Data))
		seenTopLevelNames[pascalName] = append(seenTopLevelNames[pascalName], tlt.Ident)
	}
	for _, imp := range imports {
		pascalName := hfNormalizedToPascal(hfNormalizeIdent(imp.ImportedAsIdent.Data))
		seenTopLevelNames[pascalName] = append(seenTopLevelNames[pascalName], imp.ImportedAsIdent)
	}
	ret := []LcErrorTokenCollision{}
	for k, v := range seenTopLevelNames {
		if len(v) < 2 {
			continue
		}
		entry := LcErrorTokenCollision{
			ErrT: v,
			Err: fmt.Errorf("Multiple colliding names for (%s) defined in top level. Remember "+
				"that names mustn't collide after PascalCase normalization is applied to them", k),
		}
		ret = append(ret, entry)
	}
	return ret
}



func lcTypeConvertNoCheck(astType AstTypeExpr) LcTypeExpr {
	if astType.OneofStructDef != nil {
		t := lcStructOrEnumLineConvertNoCheck(*astType.OneofStructDef)
		return LcTypeExpr{OneofStruct: &t}
	} else if astType.OneofEnumDef != nil {
		t := lcStructOrEnumLineConvertNoCheck(*astType.OneofEnumDef)
		return LcTypeExpr{OneofEnum: &t}
	} else if astType.OneofTupleDef != nil {
		t := lcTupleConvertNoCheck(*astType.OneofTupleDef)
		return LcTypeExpr{OneofTuple: &t}
	} else if astType.OneofTypeIdent != nil {
		if builtin, is := LcIsBuiltIn(astType.OneofTypeIdent.Data); is {
			return LcTypeExpr{OneofBuiltin: &builtin}
		} else {
			return LcTypeExpr{OneofTokenIdent: astType.OneofTypeIdent}
		}
	} else if astType.OneofListOf != nil {
		t := lcTypeConvertNoCheck(*astType.OneofListOf)
		return LcTypeExpr{OneofListof: &t}
	} else if astType.OneofMapOf != nil {
		t := lcTypeConvertNoCheck(*astType.OneofMapOf)
		return LcTypeExpr{OneofMapof: &t}
	} else if astType.OneofDotAccessExpr != nil {
		return LcTypeExpr{OneofImported: &LcImported{
			ImportedIdent: astType.OneofDotAccessExpr.LhsIdent,
			ForeignIdent:  astType.OneofDotAccessExpr.RhsIdent,
		}}
	} else {
		panic("unreachable")
	}
}

func lcStructOrEnumLineConvertNoCheck(astLine []AstStructOrEnumLine) []LcStructOrEnumLine {
	coll := []LcStructOrEnumLine{}
	for _, a := range astLine {
		wireName := ""
		progNameNorm := hfNormalizeIdent(a.ProgName.Data)
		if a.WireName != nil {
			wireName = a.WireName.Data
		} else {
			wireName = hfNormalizedToCamel(progNameNorm)
		}
		acc := LcStructOrEnumLine{
			WireName:   wireName,
			ProgName:   progNameNorm,
			TypeExpr:   lcTypeConvertNoCheck(a.TypeExpr),
			Omittable:  a.Omittable,
			IsReserved: a.IsReserved,
		}
		coll = append(coll, acc)
	}
	return coll
}

func lcTupleConvertNoCheck(astTuple []AstTypeExpr) []LcTypeExpr {
	coll := []LcTypeExpr{}
	for _, e := range astTuple {
		coll = append(coll, lcTypeConvertNoCheck(e))
	}
	return coll
}

func lcIsReservedIdent(ident string) error {
	allLowerCamel := strings.ToLower(hfNormalizedToCamel(hfNormalizeIdent(ident)))
	if strings.HasPrefix(allLowerCamel, "csres") {
		return fmt.Errorf("Identifiers cannot begin with `csres` as it is reserved. This rule is case insensitive.")
	} else if allLowerCamel == "enum" {
		return fmt.Errorf("Cannot use `enum` as identifier as it is reserved.")
	} else if allLowerCamel == "import" {
		return fmt.Errorf("Cannot use `import` as identifier as it is reserved.")
	} else if allLowerCamel == "map" {
		return fmt.Errorf("Cannot use `map` as identifier as it is reserved.")
	} else if allLowerCamel == "as" {
		return fmt.Errorf("Cannot use `as` as identifier as it is reserved.")
	} else if allLowerCamel == "type" {
		return fmt.Errorf("Cannot use `type` as identifier as it is reserved.")
	} else {
		return nil
	}
}
