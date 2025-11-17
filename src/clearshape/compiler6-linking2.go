package main

import "fmt"

type LnkProgram struct {
	Types map[string]LnkTypeExpr `json:"types"`
}

type LnkTypeExpr struct {
	OneofStruct     *[]LnkStructOrEnumLine `json:"topLevelStruct,omitempty"`
	OneofEnum       *[]LnkStructOrEnumLine `json:"topLevelEnum,omitempty"`
	OneofTuple      *[]LnkTypeExpr         `json:"topLevelTuple,omitempty"`
	OneofTokenIdent *Token                 `json:"tokenIdent,omitempty"`
	OneofBuiltin    *BuiltinType           `json:"builtin,omitempty"`
	OneofListof     *LnkTypeExpr           `json:"listOf,omitempty"`
	OneofMapof      *LnkTypeExpr           `json:"mapOf,omitempty"`
}

type LnkStructOrEnumLine struct {
	WireName   string      `json:"wireName"`
	ProgName   []string    `json:"progName"`
	TypeExpr   LnkTypeExpr `json:"typeExpr"`
	Omittable  bool        `json:"omittable"` // this is ignored when its an enum line
	IsReserved bool        `json:"isReserved"`
}

func lnkResolveImports(ball LnkProcessedBall) (LnkProgram, *Token, error) {
	lnkProg := LnkProgram{}
	startingProgram, exists := ball.AllPrograms[ball.StartingProgram]
	if !exists {
		panic(fmt.Sprintf("shouldn't really happen - starting program not defined: %s", ball.StartingProgram))
	}
	for lcTltIdent, lcType := range startingProgram.TopLevelDefs {
		lnkType, errT, err := lnkResolveTypeCore(lcType, startingProgram, &ball)
		if err != nil {
			return LnkProgram{}, errT, err
		}
		lnkProg.Types[lcTltIdent] = lnkType
	}
	return lnkProg, nil, nil
}

func lnkResolveTypeCore(lcType LcTypeExpr, currentProgram LnkSingleProgram, lnkBall *LnkProcessedBall) (LnkTypeExpr, *Token, error) {
	if lcType.OneofBuiltin != nil {
		return LnkTypeExpr{OneofBuiltin: lcType.OneofBuiltin}, nil, nil
	} else if lcType.OneofEnum != nil {
		enumLines, errT, err := lnkResolveLcStructOrEnum(*lcType.OneofEnum, currentProgram, lnkBall)
		if err != nil {
			return LnkTypeExpr{}, errT, err
		}
		return LnkTypeExpr{OneofEnum: &enumLines}, nil, nil
	} else if lcType.OneofImported != nil {
		importAsIdent := lcType.OneofImported.ImportedIdent.Data
		importedFileLnkImportData, has := currentProgram.Imports[importAsIdent]
		if !has {
			panic("shouldnt really happen")
		}
		importedFileAbsPath := importedFileLnkImportData.ImportSrcLocationAbsoluteString
		importedProgram, has := lnkBall.AllPrograms[importedFileAbsPath]
		if !has {
			panic("shouldnt really happen")
		}
		importForeignIdent := lcType.OneofImported.ForeignIdent.Data
		importedLcType, has := importedProgram.TopLevelDefs[importForeignIdent]
		if !has {
			return LnkTypeExpr{}, &lcType.OneofImported.ForeignIdent,
				fmt.Errorf("undefined foreign identifier %s", importForeignIdent)
		}
		return lnkResolveTypeCore(importedLcType, importedProgram, lnkBall)
	} else if lcType.OneofListof != nil {
		lnkInnerType, errT, err := lnkResolveTypeCore(*lcType.OneofListof, currentProgram, lnkBall)
		if err != nil {
			return LnkTypeExpr{}, errT, err
		}
		return LnkTypeExpr{OneofListof: &lnkInnerType}, nil, nil
	} else if lcType.OneofMapof != nil {
		lnkInnerType, errT, err := lnkResolveTypeCore(*lcType.OneofMapof, currentProgram, lnkBall)
		if err != nil {
			return LnkTypeExpr{}, errT, err
		}
		return LnkTypeExpr{OneofMapof: &lnkInnerType}, nil, nil
	} else if lcType.OneofStruct != nil {
		structLines, errT, err := lnkResolveLcStructOrEnum(*lcType.OneofStruct, currentProgram, lnkBall)
		if err != nil {
			return LnkTypeExpr{}, errT, err
		}
		return LnkTypeExpr{OneofStruct: &structLines}, nil, nil
	} else if lcType.OneofTokenIdent != nil {
		return LnkTypeExpr{OneofTokenIdent: lcType.OneofTokenIdent}, nil, nil
	} else if lcType.OneofTuple != nil {
		lnkTupleTypes, errT, err := lnkResolveLcTuple(*lcType.OneofTuple, currentProgram, lnkBall)
		if err != nil {
			return LnkTypeExpr{}, errT, err
		}
		return LnkTypeExpr{OneofTuple: &lnkTupleTypes}, nil, nil
	} else {
		panic("unreachable")
	}
}

func lnkResolveLcStructOrEnum(
	lines []LcStructOrEnumLine, 
	currentProgram LnkSingleProgram, 
	lnkBall *LnkProcessedBall,
) ([]LnkStructOrEnumLine, *Token, error) {
	coll := []LnkStructOrEnumLine{}
	for _, lcLine := range lines {
		lnkType, errT, err := lnkResolveTypeCore(lcLine.TypeExpr, currentProgram, lnkBall)
		if err != nil {
			return []LnkStructOrEnumLine{}, errT, err
		}
		lnkLine := LnkStructOrEnumLine{
			WireName: lcLine.WireName,
			ProgName: lcLine.ProgName,
			TypeExpr: lnkType,
			Omittable: lcLine.Omittable,
			IsReserved: lcLine.IsReserved,
		}
		coll = append(coll, lnkLine)
	}
	return coll, nil, nil
}

func lnkResolveLcTuple(
	lcTypes []LcTypeExpr, 
	currentProgram LnkSingleProgram, 
	lnkBall *LnkProcessedBall,
) ([]LnkTypeExpr, *Token, error) {
	coll := []LnkTypeExpr{}
	for _, lctype := range lcTypes {
		lnkType, errT, err := lnkResolveTypeCore(lctype, currentProgram, lnkBall)
		if err != nil {
			return []LnkTypeExpr{}, errT, err
		}
		coll = append(coll, lnkType)
	}
	return coll, nil, nil
}