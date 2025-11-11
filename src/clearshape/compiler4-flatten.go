package main

import (
	"fmt"
	"strconv"
)

type FltProgram struct {
	Imports          []AstImport       `json:"imports"`
	TopLevelTypedefs []FltTopLevelType `json:"topLevelDefs"`
}

type FltTopLevelType struct {
	Oneof01TopLevelName       *Token                 `json:"topLevelName,omitempty"`
	Oneof01TopLevelMintedName *string                `json:"topLevelMintedName,omitempty"`
	OneofTopLevelStruct       *[]FltStructOrEnumLine `json:"topLevelStruct,omitempty"`
	OneofTopLevelEnum         *[]FltStructOrEnumLine `json:"topLevelEnum,omitempty"`
	OneofTopLevelTuple        *[]FltTypeExpr         `json:"topLevelTuple,omitempty"`
	OneofTokenIdent           *Token                 `json:"tokenIdent,omitempty"`
	OneofListof               *FltTypeExpr           `json:"listOf,omitempty"`
	OneofImported             *FltImported           `json:"imported,omitempty"`
}

type FltTypeExpr struct {
	OneofMintedIdent *string      `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token       `json:"TokenIdent,omitempty"`
	OneofImported    *FltImported `json:"imported,omitempty"`
	OneofListof      *FltTypeExpr `json:"listOf,omitempty"`
}

type FltImported struct {
	ImportedIdent Token `json:"importedIdent"`
	ForeignIdent  Token `json:"foreignIdent"`
}

type FltStructOrEnumLine struct {
	WireName   *Token      `json:"wireName"`
	ProgName   Token       `json:"progName"`
	TypeExpr   FltTypeExpr `json:"typeExpr"`
	Omittable  bool        `json:"omittable"` // this is ignored when its an enum line
	IsReserved bool        `json:"isReserved"`
}

// expect a list of identifiers, critically these identifiers must not start with ascii 0-9 
func nameMint(ss []string) string {
	fmt.Printf("minting: %#v \n", ss)
	coll := ""
	for _, e := range ss {
		coll += strconv.Itoa(len(e))
		coll += e
	}
	fmt.Printf("minted: %s \n\n\n", coll)
	return "Csres0" + coll
}

func copyAppend(ss []string, s string) []string {
	ret := make([]string, len(ss) + 1);
	copy(ret, ss)
	ret[len(ss)] = s
	return ret
}

func fltFlattenProgram(astProgram AstProgram) (ret FltProgram) {
	fltProgram := FltProgram{}
	fltProgram.Imports = astProgram.Imports
	for _, astTypedef := range astProgram.Typedefs {
		newTLT, addedTLT := fltSerializeTop(astTypedef.Ident.Data, astTypedef.TypeExpr)
		newTLT.Oneof01TopLevelName = &astTypedef.Ident
		fltProgram.TopLevelTypedefs = append(fltProgram.TopLevelTypedefs, newTLT)
		fltProgram.TopLevelTypedefs = append(fltProgram.TopLevelTypedefs, addedTLT...)
	}
	return fltProgram
}

func fltSerializeInner(te AstTypeExpr, nestedIdentsSoFar []string) (selfT FltTypeExpr, addedTLT []FltTopLevelType) {
	if te.OneofDotAccessExpr != nil {
		return FltTypeExpr{OneofImported: &FltImported{
				ImportedIdent: te.OneofDotAccessExpr.LhsIdent,
				ForeignIdent:  te.OneofDotAccessExpr.RhsIdent}},
			[]FltTopLevelType{}
	} else if te.OneofListOf != nil {
		resSelfT, resAddedT := fltSerializeInner(*te.OneofListOf, nestedIdentsSoFar)
		return FltTypeExpr{OneofListof: &resSelfT}, resAddedT
	} else if te.OneofEnumDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofEnumDef, nestedIdentsSoFar)
		mintedTopIdent := nameMint(nestedIdentsSoFar)
		addedTlt = append(addedTlt, FltTopLevelType{
			Oneof01TopLevelMintedName: &mintedTopIdent,
			OneofTopLevelEnum:         &linesFlat,
		})
		return FltTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else if te.OneofTupleDef != nil {
		tuple, addedTLT := fltProcessTuple(*te.OneofTupleDef, nestedIdentsSoFar)
		mintedTopIdent := nameMint(nestedIdentsSoFar)
		addedTLT = append(addedTLT, FltTopLevelType{
			Oneof01TopLevelMintedName: &mintedTopIdent,
			OneofTopLevelTuple:        &tuple,
		})
		return FltTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTLT
	} else if te.OneofTypeIdent != nil {
		return FltTypeExpr{OneofTokenIdent: te.OneofTypeIdent},
			[]FltTopLevelType{}
	} else if te.OneofStructDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofStructDef, nestedIdentsSoFar)
		mintedTopIdent := nameMint(nestedIdentsSoFar)
		addedTlt = append(addedTlt, FltTopLevelType{
			Oneof01TopLevelMintedName: &mintedTopIdent,
			OneofTopLevelStruct:       &linesFlat,
		})
		return FltTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else {
		panic("unreachable")
	}
}

func fltSerializeTop(topIdent string, te AstTypeExpr) (selfTLT FltTopLevelType, addedTLT []FltTopLevelType) {
	if te.OneofDotAccessExpr != nil {
		return FltTopLevelType{OneofImported: &FltImported{
				ImportedIdent: te.OneofDotAccessExpr.LhsIdent,
				ForeignIdent:  te.OneofDotAccessExpr.RhsIdent}},
			[]FltTopLevelType{}
	} else if te.OneofListOf != nil {
		resSelfT, resAddedT := fltSerializeInner(*te.OneofListOf, []string{topIdent})
		return FltTopLevelType{OneofListof: &resSelfT}, resAddedT
	} else if te.OneofEnumDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofEnumDef, []string{topIdent})
		return FltTopLevelType{OneofTopLevelEnum: &linesFlat}, addedTlt
	} else if te.OneofTupleDef != nil {
		tuple, addedTLT := fltProcessTuple(*te.OneofTupleDef, []string{topIdent})
		return FltTopLevelType{OneofTopLevelTuple: &tuple}, addedTLT
	} else if te.OneofTypeIdent != nil {
		return FltTopLevelType{OneofTokenIdent: te.OneofTypeIdent}, []FltTopLevelType{}
	} else if te.OneofStructDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofStructDef, []string{topIdent})
		return FltTopLevelType{OneofTopLevelStruct: &linesFlat}, addedTlt
	} else {
		panic("unreachable")
	}
}

func fltProcessTuple(astTypes []AstTypeExpr, nestedIdentsSoFar []string) (resTuple []FltTypeExpr, addedTlt []FltTopLevelType) {
	tltColl := []FltTopLevelType{}
	teColl := []FltTypeExpr{}
	for i, astType := range astTypes {
		// work on ast type expr recursively
		resS, resT := fltSerializeInner(astType, copyAppend(nestedIdentsSoFar, "field" + strconv.Itoa(i)))
		tltColl = append(tltColl, resT...)
		teColl = append(teColl, resS)
	}
	return teColl, tltColl
}

func fltProcessStructOrEnum(astLines []AstStructOrEnumLine, nestedIdentsSoFar []string) (
	linesFlat []FltStructOrEnumLine, addedTlt []FltTopLevelType,
) {
	tltColl := []FltTopLevelType{}
	linesColl := []FltStructOrEnumLine{}
	for _, astLine := range astLines {
		// work on ast type expr recursively
		resS, resT := fltSerializeInner(astLine.TypeExpr, copyAppend(nestedIdentsSoFar, astLine.ProgName.Data))
		tltColl = append(tltColl, resT...)
		line := FltStructOrEnumLine{
			WireName:   astLine.WireName,
			ProgName:   astLine.ProgName,
			IsReserved: astLine.IsReserved,
			TypeExpr:   resS,
			Omittable:  astLine.Omittable,
		}
		linesColl = append(linesColl, line)
	}
	return linesColl, tltColl
}
