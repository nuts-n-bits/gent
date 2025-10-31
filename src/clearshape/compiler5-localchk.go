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
	Imports      map[string]AstImport      `json:"imports"`      // these imports should be normalized to snake_case
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

// func lcCheckProgram(FltProgram) (errT1, errT2 *Token, err error) {

// }

func lcCheckTopLevelIdentCollision(tlt []FltTopLevelType, imports []AstImport) {

}

type LcErrorEntryA struct {
	errT []Token
	err  error
}

func lcCheckStructOrEnumFieldNames(lines []FltStructOrEnumLine) []LcErrorEntryA {
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
	ret := []LcErrorEntryA{}
	for k, v := range seenProgCamel {
		if len(v) < 2 {
			continue
		}
		entry := LcErrorEntryA{
			errT: v,
			err: fmt.Errorf("Multiple colliding prog names for (%s) in the same struct/enum definition. Remember " +
				"that names mustn't collide after normalization is applied to them", k),
		}
		ret = append(ret, entry)
	}
	for k, v := range seenWire {
		if len(v) < 2 {
			continue
		}
		entry := LcErrorEntryA{
			errT: v,
			err: fmt.Errorf("Multiple colliding wire names for (%s) in the same struct/enum definition. Remember " +
				"that if a wire name is not provided, the prog name is coerced into wire name by transforming " + 
				"it into camelCase", k),
		}
		ret = append(ret, entry)
	}
	return ret
}
