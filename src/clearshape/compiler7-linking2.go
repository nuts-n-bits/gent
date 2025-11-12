package main

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
}

type LnkTypeExpr struct {
	OneofMintedIdent *string      `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token       `json:"TokenIdent,omitempty"`
	OneofBuiltin     *BuiltinType `json:"builtin,omitempty"`
	OneofListof      *LnkTypeExpr `json:"listOf,omitempty"`
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
	startingProgram := ball.AllPrograms[ball.StartingProgram]
	for tldIdent, tldValue := range startingProgram.TopLevelDefs {
		lnkResolveImportsCore(tldValue, &flatProg)
	}
}

func lnkResolveImportsCore(currentTypeDef LcTopLevelType, flatProg *LnkResolvedFlatProgram) {
	
}

