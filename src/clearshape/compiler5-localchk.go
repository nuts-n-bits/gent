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
	default:
		return BuiltinTypeString, false
	}
}

type LcProgram struct {
	// Imports is keyed by ident string without normalization
	Imports      map[string]AstImport      `json:"imports"`      // these imports should not be normalized
	TopLevelDefs map[string]LcTopLevelType `json:"topLevelDefs"` // these idents should be normalized to PascalCase
}

type LcTopLevelType struct {
	OneofTopLevelStruct *[]LcStructOrEnumLine `json:"topLevelStruct,omitempty"`
	OneofTopLevelEnum   *[]LcStructOrEnumLine `json:"topLevelEnum,omitempty"`
	OneofTopLevelTuple  *[]LcTypeExpr         `json:"topLevelTuple,omitempty"`
	OneofTokenIdent     *Token                `json:"tokenIdent,omitempty"`
	OneofBuiltin        *BuiltinType          `json:"builtin,omitempty"`
	OneofListof         *LcTypeExpr           `json:"listOf,omitempty"`
	OneofImported       *LcImported           `json:"imported,omitempty"`
}

type LcTypeExpr struct {
	OneofMintedIdent *string      `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token       `json:"TokenIdent,omitempty"`
	OneofBuiltin     *BuiltinType `json:"builtin,omitempty"`
	OneofImported    *LcImported  `json:"imported,omitempty"`
	OneofListof      *LcTypeExpr  `json:"listOf,omitempty"`
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

func (lctlt *LcTopLevelType) match(
	handlerForTopLevelStruct func(*[]LcStructOrEnumLine),
	handlerForTopLevelEnum func(*[]LcStructOrEnumLine),
	handlerForTopLevelTuple func(*[]LcTypeExpr),
	handlerForTokenIdent func(*Token),
	handlerForBuiltin func(*BuiltinType),
	handlerForListof func(*LcTypeExpr),
	handlerForImported func(*LcImported),
) {
	if lctlt.OneofTopLevelStruct != nil {
		handlerForTopLevelStruct(lctlt.OneofTopLevelStruct)
	} else if lctlt.OneofTopLevelEnum != nil {
		handlerForTopLevelEnum(lctlt.OneofTopLevelEnum)
	} else if lctlt.OneofTopLevelTuple != nil {
		handlerForTopLevelTuple(lctlt.OneofTopLevelTuple)
	} else if lctlt.OneofTokenIdent != nil {
		handlerForTokenIdent(lctlt.OneofTokenIdent)
	} else if lctlt.OneofBuiltin != nil {
		handlerForBuiltin(lctlt.OneofBuiltin)
	} else if lctlt.OneofListof != nil {
		handlerForListof(lctlt.OneofListof)
	} else if lctlt.OneofImported != nil {
		handlerForImported(lctlt.OneofImported)
	} else {
		panic("unreachable")
	}
}

func lcCheckProgram1Of2CheckReservedName(fltProgram FltProgram) (*Token, error) {
	return lcCheckTopLevelReservedName(fltProgram.TopLevelTypedefs)
}

func lcCheckTopLevelReservedName(a []FltTopLevelType) (*Token, error) {
	for _, e := range a {
		if e.Oneof01TopLevelName != nil {
			if lcIsReservedIdent(e.Oneof01TopLevelName.Data) {
				return e.Oneof01TopLevelName,
					fmt.Errorf("Identifiers cannot begin with `csres` as it is reserved. This rule is case insensitive.")
			}
		}
	}
	return nil, nil
}

func lcCheckProgram2Of2CheckCollisionAndUndefined(fltProgram FltProgram) (
	lcProg LcProgram, topLevelCollision []LcErrorTokenCollision, undefToks []Token,
) {
	tokenCollisions := []LcErrorTokenCollision{}
	t := lcCheckTopLevelIdentCollision(fltProgram.TopLevelTypedefs, fltProgram.Imports)
	tokenCollisions = append(tokenCollisions, t...)
	for _, tld := range fltProgram.TopLevelTypedefs {
		if tld.OneofTopLevelStruct != nil {
			t := lcCheckStructOrEnumFieldNames(*tld.OneofTopLevelStruct)
			tokenCollisions = append(tokenCollisions, t...)
		} else if tld.OneofTopLevelEnum != nil {
			t := lcCheckStructOrEnumFieldNames(*tld.OneofTopLevelEnum)
			tokenCollisions = append(tokenCollisions, t...)
		}
	}
	lcImports := map[string]AstImport{}
	lcTlts := map[string]LcTopLevelType{}
	for _, imp := range fltProgram.Imports {
		// dont care about duplicate since we have checked top level ident collision above
		lcImports[imp.ImportedAsIdent.Data] = imp
	}
	for _, tlt := range fltProgram.TopLevelTypedefs {
		if tlt.Oneof01TopLevelMintedName != nil {
			pascalNorm := hfNormalizedToPascal(hfNormalizeIdent(*tlt.Oneof01TopLevelMintedName))
			lcTlts[pascalNorm] = lcTltConvertNoCheckStripName(tlt)
		} else if tlt.Oneof01TopLevelName != nil {
			pascalNorm := hfNormalizedToPascal(hfNormalizeIdent(tlt.Oneof01TopLevelName.Data))
			lcTlts[pascalNorm] = lcTltConvertNoCheckStripName(tlt)
		} else {
			panic("unreachable")
		}
	}
	lcProgram := LcProgram{
		Imports:      lcImports,
		TopLevelDefs: lcTlts,
	}
	undefinedToks := lcCheckReferenceExist(lcProgram)
	return lcProgram, tokenCollisions, undefinedToks
}

// returns a list of tokens that is undefined references. len = 0 means no error.
func lcCheckReferenceExist(lcProgram LcProgram) []Token {
	undefinedTokens := []Token{}
	for _, tld := range lcProgram.TopLevelDefs {
		if tld.OneofTopLevelStruct != nil {
			for _, line := range *tld.OneofTopLevelStruct {
				lcCheckReferenceExistInner(line.TypeExpr, &undefinedTokens, lcProgram)
			}
		} else if tld.OneofTopLevelEnum != nil {
			for _, line := range *tld.OneofTopLevelEnum {
				lcCheckReferenceExistInner(line.TypeExpr, &undefinedTokens, lcProgram)
			}
		} else if tld.OneofTopLevelTuple != nil {
			for _, inner := range *tld.OneofTopLevelTuple {
				lcCheckReferenceExistInner(inner, &undefinedTokens, lcProgram)
			}
		} else if tld.OneofTokenIdent != nil {
			normPascal := hfNormalizedToPascal(hfNormalizeIdent(tld.OneofTokenIdent.Data))
			if _, has := lcProgram.TopLevelDefs[normPascal]; !has {
				undefinedTokens = append(undefinedTokens, *tld.OneofTokenIdent)
			}
		} else if tld.OneofBuiltin != nil {
			// nothing to check against
		} else if tld.OneofListof != nil {
			lcCheckReferenceExistInner(*tld.OneofListof, &undefinedTokens, lcProgram)
		} else if tld.OneofImported != nil {
			if _, has := lcProgram.Imports[tld.OneofImported.ImportedIdent.Data]; !has {
				undefinedTokens = append(undefinedTokens, tld.OneofImported.ImportedIdent)
			}
			if lcIsReservedIdent(tld.OneofImported.ForeignIdent.Data) {
				undefinedTokens = append(undefinedTokens, tld.OneofImported.ForeignIdent)
			}
		} else {
			panic("unreachable")
		}
	}
	return undefinedTokens
}

func lcCheckReferenceExistInner(fltType LcTypeExpr, undefToks *[]Token, lcProgram LcProgram) {
	if fltType.OneofMintedIdent != nil {
		// do nothing for now. minted ident should not be undefined
	} else if fltType.OneofTokenIdent != nil {
		normPascal := hfNormalizedToPascal(hfNormalizeIdent(fltType.OneofTokenIdent.Data))
		if _, has := lcProgram.TopLevelDefs[normPascal]; !has {
			*undefToks = append(*undefToks, *fltType.OneofTokenIdent)
		}
	} else if fltType.OneofBuiltin != nil {
		// nothing to check
	} else if fltType.OneofImported != nil {
		if _, has := lcProgram.Imports[fltType.OneofImported.ImportedIdent.Data]; !has {
			*undefToks = append(*undefToks, fltType.OneofImported.ImportedIdent)
		}
	} else if fltType.OneofListof != nil {
		lcCheckReferenceExistInner(*fltType.OneofListof, undefToks, lcProgram)
	} else {
		panic("unreachable")
	}
}

type LcErrorTokenCollision struct {
	ErrT []Token `json:"errT"` // always len >= 2
	Err  error   `json:"err"`  // always defined
}

func lcCheckTopLevelIdentCollision(tlts []FltTopLevelType, imports []AstImport) []LcErrorTokenCollision {
	seenTli := map[string][]Token{}
	for _, tlt := range tlts {
		if tlt.Oneof01TopLevelName != nil {
			pascalName := hfNormalizedToPascal(hfNormalizeIdent(tlt.Oneof01TopLevelName.Data))
			seenTli[pascalName] = append(seenTli[pascalName], *tlt.Oneof01TopLevelName)
		} else if tlt.Oneof01TopLevelMintedName != nil {
			// do nothing for now since minted name is not supposed to collide anyways
		} else {
			panic("unreachable")
		}
	}
	for _, imp := range imports {
		pascalName := hfNormalizedToPascal(hfNormalizeIdent(imp.ImportedAsIdent.Data))
		seenTli[pascalName] = append(seenTli[pascalName], imp.ImportedAsIdent)
	}
	ret := []LcErrorTokenCollision{}
	for k, v := range seenTli {
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

func lcCheckStructOrEnumFieldNames(lines []FltStructOrEnumLine) []LcErrorTokenCollision {
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

func lcTltConvertNoCheckStripName(fltTlt FltTopLevelType) LcTopLevelType {
	if fltTlt.OneofTopLevelStruct != nil {
		t := lcStructOrEnumLineConvertNoCheck(*fltTlt.OneofTopLevelStruct)
		return LcTopLevelType{OneofTopLevelStruct: &t}
	} else if fltTlt.OneofTopLevelEnum != nil {
		t := lcStructOrEnumLineConvertNoCheck(*fltTlt.OneofTopLevelEnum)
		return LcTopLevelType{OneofTopLevelStruct: &t}
	} else if fltTlt.OneofTopLevelTuple != nil {
		t := lcTupleConvertNoCheck(*fltTlt.OneofTopLevelTuple)
		return LcTopLevelType{OneofTopLevelTuple: &t}
	} else if fltTlt.OneofTokenIdent != nil {
		if builtin, is := LcIsBuiltIn(fltTlt.OneofTokenIdent.Data); is {
			return LcTopLevelType{OneofBuiltin: &builtin}
		} else {
			return LcTopLevelType{OneofTokenIdent: fltTlt.OneofTokenIdent}
		}
	} else if fltTlt.OneofListof != nil {
		t := lcTypeExprConvertNoCheck(*fltTlt.OneofListof)
		return LcTopLevelType{OneofListof: &t}
	} else if fltTlt.OneofImported != nil {
		return LcTopLevelType{OneofImported: &LcImported{
			ImportedIdent: fltTlt.OneofImported.ImportedIdent,
			ForeignIdent:  fltTlt.OneofImported.ForeignIdent,
		}}
	} else {
		panic("unreachable")
	}
}

func lcStructOrEnumLineConvertNoCheck(fltLine []FltStructOrEnumLine) []LcStructOrEnumLine {
	coll := []LcStructOrEnumLine{}
	for _, a := range fltLine {
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
			TypeExpr:   lcTypeExprConvertNoCheck(a.TypeExpr),
			Omittable:  a.Omittable,
			IsReserved: a.IsReserved,
		}
		coll = append(coll, acc)
	}
	return coll
}

func lcTupleConvertNoCheck(fltLine []FltTypeExpr) []LcTypeExpr {
	coll := []LcTypeExpr{}
	for _, e := range fltLine {
		coll = append(coll, lcTypeExprConvertNoCheck(e))
	}
	return coll
}

func lcTypeExprConvertNoCheck(fltTypeExpr FltTypeExpr) LcTypeExpr {
	if fltTypeExpr.OneofMintedIdent != nil {
		return LcTypeExpr{OneofMintedIdent: fltTypeExpr.OneofMintedIdent}
	} else if fltTypeExpr.OneofTokenIdent != nil {
		if builtin, is := LcIsBuiltIn(fltTypeExpr.OneofTokenIdent.Data); is {
			return LcTypeExpr{OneofBuiltin: &builtin}
		} else {
			return LcTypeExpr{OneofTokenIdent: fltTypeExpr.OneofTokenIdent}
		}
	} else if fltTypeExpr.OneofImported != nil {
		return LcTypeExpr{OneofImported: &LcImported{
			ImportedIdent: fltTypeExpr.OneofImported.ImportedIdent,
			ForeignIdent:  fltTypeExpr.OneofImported.ForeignIdent,
		}}
	} else if fltTypeExpr.OneofListof != nil {
		t := lcTypeExprConvertNoCheck(*fltTypeExpr.OneofListof)
		return LcTypeExpr{OneofListof: &t}
	} else {
		panic("unreachable")
	}
}

func lcIsReservedIdent(ident string) bool {
	allLowerCamel := strings.ToLower(hfNormalizedToCamel(hfNormalizeIdent(ident)))
	return strings.HasPrefix(allLowerCamel, "csres")
}






