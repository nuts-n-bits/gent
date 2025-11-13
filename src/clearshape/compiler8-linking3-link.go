package main

import "fmt"

type LnkResolvedFlatProgram struct {
	Types []LnkTopLevelType `json:"types"`
}

type LnkTopLevelType struct {
	OneofTopLevelStruct *[]LnkStructOrEnumLine `json:"topLevelStruct,omitempty"`
	OneofTopLevelEnum   *[]LnkStructOrEnumLine `json:"topLevelEnum,omitempty"`
	OneofTopLevelTuple  *[]LnkTypeExpr         `json:"topLevelTuple,omitempty"`
	OneofTokenIdent     *Token                 `json:"tokenIdent,omitempty"`
	OneofBuiltin        *BuiltinType           `json:"builtin,omitempty"`
	OneofListof         *LnkTypeExpr           `json:"listOf,omitempty"`
	OneofMapof          *LnkTypeExpr           `json:"mapOf,omitempty"`
}

type LnkTypeExpr struct {
	OneofMintedIdent *string      `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token       `json:"TokenIdent,omitempty"`
	OneofBuiltin     *BuiltinType `json:"builtin,omitempty"`
	OneofListof      *LnkTypeExpr `json:"listOf,omitempty"`
	OneofMapof       *LnkTypeExpr `json:"mapOf,omitempty"`
}

type LnkStructOrEnumLine struct {
	WireName   string      `json:"wireName"`
	ProgName   []string    `json:"progName"`
	TypeExpr   LnkTypeExpr `json:"typeExpr"`
	Omittable  bool        `json:"omittable"` // this is ignored when its an enum line
	IsReserved bool        `json:"isReserved"`
}

func lnkResolveImports(ball LnkProcessedBall) (LnkResolvedFlatProgram, error) {
	flatProg := LnkResolvedFlatProgram{}
	startingProgram, exists := ball.AllPrograms[ball.StartingProgram]
	if !exists {
		panic(fmt.Sprintf("shouldn't really happen - starting program not defined: %s", ball.StartingProgram))
	}
	for tldIdent, tldValue := range startingProgram.TopLevelDefs {
		lnkResolveImportsTltCore(tldValue, &flatProg)
	}
}

func lnkResolveImportsTltCore(currentLcTlt LcTopLevelType, flatProg *LnkResolvedFlatProgram) {
	var lnkTlt LnkTopLevelType
	currentLcTlt.match(func(structDef *[]LcStructOrEnumLine) {
		// TODO
	}, func(enumDef *[]LcStructOrEnumLine) {
		// TODO
	}, func(tupleDef *[]LcTypeExpr) {
		// TODO
	}, func(tokenIdent *Token) {
		lnkTlt = LnkTopLevelType{OneofTokenIdent: tokenIdent}
	}, func(builtin *BuiltinType) {
		lnkTlt = LnkTopLevelType{OneofBuiltin: builtin}
	}, func(listOf *LcTypeExpr) {
		// TODO
	}, func(mapOf *LcTypeExpr) {
		// TODO
	}, func(imported *LcImported) {

	})
	flatProg.Types = append(flatProg.Types, )
}

func lnkResolveImportsCore(lcType ) {

}