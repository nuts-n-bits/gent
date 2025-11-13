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

func lnkResolveImportsTltCore(
	currentLcTlt LcTopLevelType, 
	currentLnkProg LnkSingleProgram, 
	lnkBall LnkProcessedBall, 
	flatProg *LnkResolvedFlatProgram,
) (errTok *Token, err error) {
	if currentLcTlt.OneofBuiltin != nil {
		lnkTlt := LnkTopLevelType{OneofBuiltin: currentLcTlt.OneofBuiltin}
		flatProg.Types = append(flatProg.Types, lnkTlt)
	} else if currentLcTlt.OneofImported != nil {
		importAsIdent := currentLcTlt.OneofImported.ImportedIdent.Data
		importedFileLnkImportData, has := currentLnkProg.Imports[importAsIdent]
		if !has {
			panic("shouldnt really happen")
		}
		importedFileAbsPath := importedFileLnkImportData.ImportSrcLocationAbsoluteString
		importedProgram, has := lnkBall.AllPrograms[importedFileAbsPath]
		if !has {
			panic("shouldnt really happen")
		}
		importForeignIdent := currentLcTlt.OneofImported.ForeignIdent.Data
		importedTld, has := importedProgram.TopLevelDefs[importForeignIdent] 
		if !has {
			return &currentLcTlt.OneofImported.ForeignIdent, 
			fmt.Errorf("undefined foreign identifier %s", importForeignIdent)
		}
		nameMint([]string{})
	} else if currentLcTlt.OneofListof != nil {
		// TODO

	} else if currentLcTlt.OneofMapof != nil {
		// TODO

	} else if currentLcTlt.OneofTokenIdent != nil {
		lnkTlt := LnkTopLevelType{OneofTokenIdent: currentLcTlt.OneofTokenIdent}
		flatProg.Types = append(flatProg.Types, lnkTlt)
	} else if currentLcTlt.OneofTopLevelEnum != nil {
		// TODO

	} else if currentLcTlt.OneofTopLevelStruct != nil {
		// TODO

	} else if currentLcTlt.OneofTopLevelTuple != nil {
		// TODO

	} else {
		panic("unreachable")
	}	
}

func lnkResolveImportsCore(lcType LcTypeExpr, flatProg *LnkResolvedFlatProgram) {

}

func lnkCreateLocalTypeForImport(lcTlt LcTopLevelType, flatProg *LnkResolvedFlatProgram) {

}