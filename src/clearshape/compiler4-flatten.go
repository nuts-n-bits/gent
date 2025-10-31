package main

import (
	"strconv"
)

type FltProgram struct {
	// Imports is keyed by ident string without normalization
	Imports      []AstImport       `json:"imports"`      // these imports should be normalized to snake_case
	TopLevelDefs []FltTopLevelType `json:"topLevelDefs"` // these idents should be normalized to PascalCase
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

// func chkProgram(astProgram AstProgram) (ChkProgram, *Token, error) {
// 	chkProgram := ChkProgram{};
// 	for _, typedef := range astProgram.Typedefs {
// 		topLevelIdent := typedef.Ident.Data;
// 		typedef.TypeExpr
// 	}
// }

//func chkResolveTypeExpr

var ctr int = 0

func counter() string {
	ret := strconv.Itoa(ctr)
	ctr += 1
	return "csres0" + ret
}

func fltFlattenProgram(astProgram AstProgram) (ret FltProgram) {
	fltProgram := FltProgram{}
	fltProgram.Imports = astProgram.Imports
	for _, astTypedef := range astProgram.Typedefs {
		newTLT, addedTLT := fltSerializeTop(astTypedef.TypeExpr)
		newTLT.Oneof01TopLevelName = &astTypedef.Ident
		fltProgram.TopLevelDefs = append(fltProgram.TopLevelDefs, newTLT)
		fltProgram.TopLevelDefs = append(fltProgram.TopLevelDefs, addedTLT...)
	}
	return fltProgram
}

func fltSerializeInner(te AstTypeExpr) (selfT FltTypeExpr, addedTLT []FltTopLevelType) {
	if te.OneofDotAccessExpr != nil {
		return FltTypeExpr{OneofImported: &FltImported{
				ImportedIdent: te.OneofDotAccessExpr.LhsIdent,
				ForeignIdent:  te.OneofDotAccessExpr.RhsIdent}},
			[]FltTopLevelType{}
	} else if te.OneofListOf != nil {
		resSelfT, resAddedT := fltSerializeInner(*te.OneofListOf)
		return FltTypeExpr{OneofListof: &resSelfT}, resAddedT
	} else if te.OneofEnumDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofEnumDef)
		mintedTopIdent := counter()
		addedTlt = append(addedTlt, FltTopLevelType{
			Oneof01TopLevelMintedName: &mintedTopIdent,
			OneofTopLevelEnum:         &linesFlat,
		})
		return FltTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else if te.OneofTupleDef != nil {
		tuple, addedTLT := fltProcessTuple(*te.OneofTupleDef)
		mintedTopIdent := counter()
		addedTLT = append(addedTLT, FltTopLevelType{
			Oneof01TopLevelMintedName: &mintedTopIdent,
			OneofTopLevelTuple:        &tuple,
		})
		return FltTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTLT
	} else if te.OneofTypeIdent != nil {
		return FltTypeExpr{OneofTokenIdent: te.OneofTypeIdent},
			[]FltTopLevelType{}
	} else if te.OneofStructDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofStructDef)
		mintedTopIdent := counter()
		addedTlt = append(addedTlt, FltTopLevelType{
			Oneof01TopLevelMintedName: &mintedTopIdent,
			OneofTopLevelStruct:       &linesFlat,
		})
		return FltTypeExpr{OneofMintedIdent: &mintedTopIdent}, addedTlt
	} else {
		panic("unreachable")
	}
}

func fltSerializeTop(te AstTypeExpr) (selfTLT FltTopLevelType, addedTLT []FltTopLevelType) {
	if te.OneofDotAccessExpr != nil {
		return FltTopLevelType{OneofImported: &FltImported{
				ImportedIdent: te.OneofDotAccessExpr.LhsIdent,
				ForeignIdent:  te.OneofDotAccessExpr.RhsIdent}},
			[]FltTopLevelType{}
	} else if te.OneofListOf != nil {
		resSelfT, resAddedT := fltSerializeInner(*te.OneofListOf)
		return FltTopLevelType{OneofListof: &resSelfT}, resAddedT
	} else if te.OneofEnumDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofEnumDef)
		return FltTopLevelType{OneofTopLevelEnum: &linesFlat}, addedTlt
	} else if te.OneofTupleDef != nil {
		tuple, addedTLT := fltProcessTuple(*te.OneofTupleDef)
		return FltTopLevelType{OneofTopLevelTuple: &tuple}, addedTLT
	} else if te.OneofTypeIdent != nil {
		return FltTopLevelType{OneofTokenIdent: te.OneofTypeIdent}, []FltTopLevelType{}
	} else if te.OneofStructDef != nil {
		linesFlat, addedTlt := fltProcessStructOrEnum(*te.OneofStructDef)
		return FltTopLevelType{OneofTopLevelStruct: &linesFlat}, addedTlt
	} else {
		panic("unreachable")
	}
}

func fltProcessTuple(astTypes []AstTypeExpr) (resTuple []FltTypeExpr, addedTlt []FltTopLevelType) {
	tltColl := []FltTopLevelType{}
	teColl := []FltTypeExpr{}
	for _, astType := range astTypes {
		// work on ast type expr recursively
		resS, resT := fltSerializeInner(astType)
		tltColl = append(tltColl, resT...)
		teColl = append(teColl, resS)
	}
	return teColl, tltColl
}

func fltProcessStructOrEnum(astLines []AstStructOrEnumLine) (
	linesFlat []FltStructOrEnumLine, addedTlt []FltTopLevelType,
) {
	tltColl := []FltTopLevelType{}
	linesColl := []FltStructOrEnumLine{}
	for _, astLine := range astLines {
		// work on ast type expr recursively
		resS, resT := fltSerializeInner(astLine.TypeExpr)
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
