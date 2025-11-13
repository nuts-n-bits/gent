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
	for lcTltIdent, lcTltValue := range startingProgram.TopLevelDefs {
		if lcTltValue.OneofBuiltin != nil {
			lnkTlt := LnkTopLevelType{OneofBuiltin: lcTltValue.OneofBuiltin}
			flatProg.Types = append(flatProg.Types, lnkTlt)
		} else if lcTltValue.OneofImported != nil {
			importAsIdent := lcTltValue.OneofImported.ImportedIdent.Data
			importedFileLnkImportData, has := startingProgram.Imports[importAsIdent]
			if !has {
				panic("shouldnt really happen")
			}
			importedFileAbsPath := importedFileLnkImportData.ImportSrcLocationAbsoluteString
			importedProgram, has := lnkBall.AllPrograms[importedFileAbsPath]
			if !has {
				panic("shouldnt really happen")
			}
			importForeignIdent := lcTltValue.OneofImported.ForeignIdent.Data
			importedTld, has := importedProgram.TopLevelDefs[importForeignIdent] 
			if !has {
				return &lcTltValue.OneofImported.ForeignIdent, 
				fmt.Errorf("undefined foreign identifier %s", importForeignIdent)
			}
			newName := nameMint([]string{"import", importAsIdent, importForeignIdent})
		} else if lcTltValue.OneofListof != nil {
			// TODO

		} else if lcTltValue.OneofMapof != nil {
			// TODO

		} else if lcTltValue.OneofTokenIdent != nil {
			lnkTlt := LnkTopLevelType{OneofTokenIdent: lcTltValue.OneofTokenIdent}
			flatProg.Types = append(flatProg.Types, lnkTlt)
		} else if lcTltValue.OneofTopLevelEnum != nil {
			// TODO

		} else if lcTltValue.OneofTopLevelStruct != nil {
			// TODO

		} else if lcTltValue.OneofTopLevelTuple != nil {
			// TODO

		} else {
			panic("unreachable")
		}	
	}
}

func lnkResolveImportsTltCore(
	currentLcTlt LcTopLevelType, 
	currentLnkProg LnkSingleProgram, 
	lnkBall LnkProcessedBall, 
	flatProg *LnkResolvedFlatProgram,
) (errTok *Token, err error) {

}

func lnkResolveImportsCore(lcType LcTypeExpr, flatProg *LnkResolvedFlatProgram) {

}

func lnkCreateLocalTypeForImport(lcTlt LcTopLevelType, flatProg *LnkResolvedFlatProgram) {

}