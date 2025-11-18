package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

type FltProgram struct {
	TopTypes map[string]FltTopLevelTypeExpr `json:"topLevelDefs"`
}

type FltTopLevelTypeExpr struct {
	OneofBuiltin    *BuiltinType           `json:"builtin,omitempty"`
	OneofEnum       *[]FltStructOrEnumLine `json:"enum,omitempty"`
	OneofListof     *FltInnerTypeExpr      `json:"listOf,omitempty"`
	OneofMapof      *FltInnerTypeExpr      `json:"mapOf,omitempty"`
	OneofStruct     *[]FltStructOrEnumLine `json:"struct,omitempty"`
	OneofTokenIdent *Token                 `json:"tokenIdent,omitempty"`
	OneofTuple      *[]FltInnerTypeExpr    `json:"tuple,omitempty"`
}

type FltInnerTypeExpr struct {
	OneofBuiltin     *BuiltinType      `json:"builtin,omitempty"`
	OneofMintedIdent *string           `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token            `json:"TokenIdent,omitempty"`
	OneofListof      *FltInnerTypeExpr `json:"listOf,omitempty"`
	OneofMapof       *FltInnerTypeExpr `json:"mapOf,omitempty"`
}

type FltStructOrEnumLine struct {
	WireName   string           `json:"wireName"`
	ProgName   []string         `json:"progName"`
	TypeExpr   FltInnerTypeExpr `json:"typeExpr"`
	Omittable  bool             `json:"omittable"` // this is ignored when its an enum line
	IsReserved bool             `json:"isReserved"`
}

// expect a list of identifiers, critically these identifiers must not start with ascii 0-9
// TODO: if the resulting minted identifier is longer than 64 characters, the entire thing is truncated to the first 32
// characters, plus another 32 characters of hex hash that is obtained by running the entire string through sha2-256
// (sha2-256 outputs 64 hex asciis, the first 32 are used)
func nameMint(ss []string) string {
	fmt.Printf("minting: %#v \n", ss)
	coll := ""
	for _, e := range ss {
		coll += strconv.Itoa(len(e))
		coll += e
	}
	coll = "Csres0" + coll
	if len(coll) <= 64 {
		return coll
	}
	sum := hex.EncodeToString(sha256.New().Sum([]byte(coll)))
	coll = "Csres1" + coll[6:32] + (sum)[0:32]
	return coll
}

func copyAppend(ss []string, s string) []string {
	ret := make([]string, len(ss)+1)
	copy(ret, ss)
	ret[len(ss)] = s
	return ret
}

func fltFlattenProgram(lnkProgram LnkProgram) FltProgram {
	ret := FltProgram{TopTypes: map[string]FltTopLevelTypeExpr{}}
	for topLevelIdent, typeExpr := range lnkProgram.Types {
		selfTLT, addedTLT := fltSerializeTop(topLevelIdent, typeExpr)
		fltMergeMaps(&ret.TopTypes, map[string]FltTopLevelTypeExpr{topLevelIdent: selfTLT})
		fltMergeMaps(&ret.TopTypes, addedTLT)
	}
	return ret
}

// the flt pass should never produce name collisions in the top level, regardless of user input -- because names
// LIKE "Csres%" is reserved so its safe to always create new names that start with these characters, a name collision
// always means something is wrong with the name minting process. therefore it panics instead of returning an error
func fltMergeMaps(sink *map[string]FltTopLevelTypeExpr, source map[string]FltTopLevelTypeExpr) {
	for k, v := range source {
		if _, has := (*sink)[k]; has {
			panic(fmt.Sprintf("duplicate identifier string %s while merging maps!", k))
		}
		(*sink)[k] = v
	}
}

// a helper function for fltMergeMaps
func fltAddToMap(sink *map[string]FltTopLevelTypeExpr, key string, value FltTopLevelTypeExpr) {
	fltMergeMaps(sink, map[string]FltTopLevelTypeExpr{key: value})
}

func fltSerializeInner(te LnkTypeExpr, nestedIdentsSoFar []string) (selfInnerT FltInnerTypeExpr, addedTlt map[string]FltTopLevelTypeExpr) {
	if te.OneofBuiltin != nil {
		return FltInnerTypeExpr{OneofBuiltin: te.OneofBuiltin}, map[string]FltTopLevelTypeExpr{}
	} else if te.OneofEnum != nil {
		selfInnerT, addedTlt := fltProcessStructOrEnum(*te.OneofEnum, nestedIdentsSoFar)
		mintedTopIdent := nameMint(nestedIdentsSoFar)
		fltAddToMap(&addedTlt, mintedTopIdent, FltTopLevelTypeExpr{OneofEnum: &selfInnerT})
		return FltInnerTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else if te.OneofListof != nil {
		resSelfT, addedTlt := fltSerializeInner(*te.OneofListof, nestedIdentsSoFar)
		return FltInnerTypeExpr{OneofListof: &resSelfT}, addedTlt
	} else if te.OneofMapof != nil {
		resSelfT, addedTLT := fltSerializeInner(*te.OneofMapof, nestedIdentsSoFar)
		return FltInnerTypeExpr{OneofMapof: &resSelfT}, addedTLT
	} else if te.OneofStruct != nil {
		selfInnerT, addedTlt := fltProcessStructOrEnum(*te.OneofStruct, nestedIdentsSoFar)
		mintedTopIdent := nameMint(nestedIdentsSoFar)
		fltAddToMap(&addedTlt, mintedTopIdent, FltTopLevelTypeExpr{OneofStruct: &selfInnerT})
		return FltInnerTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else if te.OneofTokenIdent != nil {
		return FltInnerTypeExpr{OneofTokenIdent: te.OneofTokenIdent}, map[string]FltTopLevelTypeExpr{}
	} else if te.OneofTuple != nil {
		selfInnerT, addedTlt := fltProcessTuple(*te.OneofTuple, nestedIdentsSoFar)
		mintedTopIdent := nameMint(nestedIdentsSoFar)
		fltAddToMap(&addedTlt, mintedTopIdent, FltTopLevelTypeExpr{OneofTuple: &selfInnerT})
		return FltInnerTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else {
		panic("unreachable")
	}
}

func fltSerializeTop(topIdent string, typeExpr LnkTypeExpr) (selfTlt FltTopLevelTypeExpr, addedTlt map[string]FltTopLevelTypeExpr) {
	if typeExpr.OneofBuiltin != nil {
		return FltTopLevelTypeExpr{OneofBuiltin: typeExpr.OneofBuiltin}, map[string]FltTopLevelTypeExpr{}
	} else if typeExpr.OneofEnum != nil {
		selfTlt, addedTlt := fltProcessStructOrEnum(*typeExpr.OneofEnum, []string{topIdent})
		return FltTopLevelTypeExpr{OneofEnum: &selfTlt}, addedTlt
	} else if typeExpr.OneofListof != nil {
		selfTlt, addedTlt := fltSerializeInner(*typeExpr.OneofListof, []string{topIdent})
		return FltTopLevelTypeExpr{OneofListof: &selfTlt}, addedTlt
	} else if typeExpr.OneofMapof != nil {
		selfTlt, addedTlt := fltSerializeInner(*typeExpr.OneofMapof, []string{topIdent})
		return FltTopLevelTypeExpr{OneofMapof: &selfTlt}, addedTlt
	} else if typeExpr.OneofStruct != nil {
		selfTlt, addedTlt := fltProcessStructOrEnum(*typeExpr.OneofStruct, []string{topIdent})
		return FltTopLevelTypeExpr{OneofStruct: &selfTlt}, addedTlt
	} else if typeExpr.OneofTokenIdent != nil {
		return FltTopLevelTypeExpr{OneofTokenIdent: typeExpr.OneofTokenIdent}, map[string]FltTopLevelTypeExpr{}
	} else if typeExpr.OneofTuple != nil {
		selfTlt, addedTlt := fltProcessTuple(*typeExpr.OneofTuple, []string{topIdent})
		return FltTopLevelTypeExpr{OneofTuple: &selfTlt}, addedTlt
	} else {
		panic("unreachable")
	}
}

func fltProcessTuple(lnkTypes []LnkTypeExpr, nestedIdentsSoFar []string) (
	resTuple []FltInnerTypeExpr, addedTlt map[string]FltTopLevelTypeExpr,
) {
	addedTlt = map[string]FltTopLevelTypeExpr{}
	tupleTypeColl := []FltInnerTypeExpr{}
	for i, lnkType := range lnkTypes {
		// work on ast type expr recursively
		selfInnerT, addedTltInner := fltSerializeInner(lnkType, copyAppend(nestedIdentsSoFar, "field"+strconv.Itoa(i)))
		fltMergeMaps(&addedTlt, addedTltInner)
		tupleTypeColl = append(tupleTypeColl, selfInnerT)
	}
	return tupleTypeColl, addedTlt
}

func fltProcessStructOrEnum(lnkLines []LnkStructOrEnumLine, nestedIdentsSoFar []string) (
	linesFlat []FltStructOrEnumLine, addedTlt2 map[string]FltTopLevelTypeExpr,
) {
	addedTlt := map[string]FltTopLevelTypeExpr{}
	linesColl := []FltStructOrEnumLine{}
	for _, lnkLine := range lnkLines {
		// work on ast type expr recursively
		selfInnerT, addedTltInner :=
			fltSerializeInner(lnkLine.TypeExpr, copyAppend(nestedIdentsSoFar, hfNormalizedToCamel(lnkLine.ProgName)))
		fltMergeMaps(&addedTlt, addedTltInner)
		line := FltStructOrEnumLine{
			WireName:   lnkLine.WireName,
			ProgName:   lnkLine.ProgName,
			IsReserved: lnkLine.IsReserved,
			TypeExpr:   selfInnerT,
			Omittable:  lnkLine.Omittable,
		}
		linesColl = append(linesColl, line)
	}
	return linesColl, addedTlt
}
