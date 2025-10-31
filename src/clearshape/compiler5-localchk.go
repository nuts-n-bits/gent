package main

import "fmt"

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
	OneofListof         *LcTypeExpr           `json:"listOf,omitempty"`
	OneofImported       *LcImported           `json:"imported,omitempty"`
}

type LcTypeExpr struct {
	OneofMintedIdent *string     `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token      `json:"TokenIdent,omitempty"`
	OneofImported    *LcImported `json:"imported,omitempty"`
	OneofListof      *LcTypeExpr `json:"listOf,omitempty"`
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

func lcCheckProgram(fltProgram FltProgram) (topLevelCollision []LcErrorTokenCollision) {
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
	if len(tokenCollisions) > 0 {
		return tokenCollisions
	}
	lcImports := map[string]AstImport{}
	lcTlts := map[string]LcTopLevelType{}
	for _, imp := range fltProgram.Imports {
		// wont duplicate since we have checked top level ident collision above
		lcImports[imp.ImportedAsIdent.Data] = imp
	}
	for _, tlt := range fltProgram.TopLevelTypedefs {
		if tlt.Oneof01TopLevelMintedName != nil {
			pascalNorm := hfNormalizedToPascal(hfNormalizeIdent(*tlt.Oneof01TopLevelMintedName))
			lcTlts[pascalNorm] = LcTltConvertNoCheckStripName(tlt)
		} else if tlt.Oneof01TopLevelName != nil {

		} else {
			panic("unreachable")
		}
	}
}

func lcCheckReferenceExist(lcProgram LcProgram) ([]Token, error) {
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
		} else if tld.OneofListof != nil {
			lcCheckReferenceExistInner(*tld.OneofListof, &undefinedTokens, lcProgram)
		} else if tld.OneofImported != nil {
			if _, has := lcProgram.Imports[tld.OneofImported.ImportedIdent.Data]; !has {
				undefinedTokens = append(undefinedTokens, tld.OneofImported.ImportedIdent)
			}
		} else {
			panic("unreachable")
		}
	}
	if len(undefinedTokens) > 0 {
		return undefinedTokens, fmt.Errorf("Undefined identifier(s)")
	} else {
		return undefinedTokens, nil
	}
}

func lcCheckReferenceExistInner(fltType LcTypeExpr, undefToks *[]Token, lcProgram LcProgram) {
	if fltType.OneofMintedIdent != nil {
		// do nothing for now. minted ident should not be undefined
	} else if fltType.OneofTokenIdent != nil {
		normPascal := hfNormalizedToPascal(hfNormalizeIdent(fltType.OneofTokenIdent.Data))
		if _, has := lcProgram.TopLevelDefs[normPascal]; !has {
			*undefToks = append(*undefToks, *fltType.OneofTokenIdent)
		}
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
	// always len >= 2
	errT []Token
	// always defined
	err error
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
			errT: v,
			err: fmt.Errorf("Multiple colliding names for (%s) defined in top level. Remember "+
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
			errT: v,
			err: fmt.Errorf("Multiple colliding prog names for (%s) in the same struct/enum definition. Remember "+
				"that names mustn't collide after normalization is applied to them", k),
		}
		ret = append(ret, entry)
	}
	for k, v := range seenWire {
		if len(v) < 2 {
			continue
		}
		entry := LcErrorTokenCollision{
			errT: v,
			err: fmt.Errorf("Multiple colliding wire names for (%s) in the same struct/enum definition. Remember "+
				"that if a wire name is not provided, the prog name is coerced into wire name by transforming "+
				"it into camelCase", k),
		}
		ret = append(ret, entry)
	}
	return ret
}

func LcTltConvertNoCheckStripName(fltTlt FltTopLevelType) LcTopLevelType {
	if fltTlt.OneofTopLevelStruct != nil {
		t := LcStructOrEnumLineConvertNoCheck(*fltTlt.OneofTopLevelStruct)
		return LcTopLevelType{OneofTopLevelStruct: &t}
	} else if fltTlt.OneofTopLevelEnum != nil {
		t := LcStructOrEnumLineConvertNoCheck(*fltTlt.OneofTopLevelEnum)
		return LcTopLevelType{OneofTopLevelStruct: &t}
	} else if fltTlt.OneofTopLevelTuple != nil {
		t := LcTupleConvertNoCheck(*fltTlt.OneofTopLevelTuple)
		return LcTopLevelType{OneofTopLevelTuple: &t}
	} else if fltTlt.OneofTokenIdent != nil {
		return LcTopLevelType{OneofTokenIdent: fltTlt.OneofTokenIdent}
	} else if fltTlt.OneofListof != nil {
		t := LcTypeExprConvertNoCheck(*fltTlt.OneofListof)
		return LcTopLevelType{OneofListof: &t}
	} else if fltTlt.OneofImported != nil {
		return LcTopLevelType{OneofImported: &LcImported{
			ImportedIdent: fltTlt.OneofImported.ImportedIdent,
			ForeignIdent: fltTlt.OneofImported.ForeignIdent,
		}}
	} else {
		panic("unreachable")
	}
}

func LcStructOrEnumLineConvertNoCheck(fltLine []FltStructOrEnumLine) []LcStructOrEnumLine {
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
			WireName: wireName,
			ProgName: progNameNorm,
			TypeExpr: LcTypeExprConvertNoCheck(a.TypeExpr),
			Omittable: a.Omittable,
			IsReserved: a.IsReserved,
		}
	}
	return coll 
}

func LcTupleConvertNoCheck(fltLine []FltTypeExpr) []LcTypeExpr {
	coll := []LcTypeExpr{}
	for _, e := range fltLine {
		coll = append(coll, LcTypeExprConvertNoCheck(e))
	}
	return coll
}

func LcTypeExprConvertNoCheck(fltTypeExpr FltTypeExpr) LcTypeExpr {
 	if fltTypeExpr.OneofMintedIdent != nil {
		return LcTypeExpr{OneofMintedIdent: fltTypeExpr.OneofMintedIdent}
	} else if fltTypeExpr.OneofTokenIdent != nil {
		return LcTypeExpr{OneofTokenIdent: fltTypeExpr.OneofTokenIdent}
	} else if fltTypeExpr.OneofImported != nil {
		return LcTypeExpr{OneofImported: &LcImported{
			ImportedIdent: fltTypeExpr.OneofImported.ImportedIdent,
			ForeignIdent: fltTypeExpr.OneofImported.ForeignIdent,
		}}
	} else if fltTypeExpr.OneofListof != nil {
		t := LcTypeExprConvertNoCheck(*fltTypeExpr.OneofListof)
		return LcTypeExpr{OneofListof: &t}
	} else {
		panic("unreachable")
	}
}
