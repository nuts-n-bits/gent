// in linking 1, we gather all the imported files recursively and check if each imported symbol exists

package main

import (
	"fmt"
	"path/filepath"
)

type LnkErrorUnion struct {
	OneofReadfileErr     error                    `json:"oneofReadfileErr,omitempty"`
	OneofLexErr          *LnkLexErr               `json:"oneofLexErr,omitempty"`
	OneofParseErr        *LnkTokenErr1            `json:"oneofParseErr,omitempty"`
	OneofReservedNameErr *LnkTokenErr1            `json:"oneofReservedNameErr,omitempty"`
	OneofLcCollisionErr  *[]LcErrorTokenCollision `json:"oneofLcCollisionErr,omitempty"`
	OneofUndefErr        *[]Token                 `json:"oneofUndefErr,omitempty"`
}

type LnkLexErr struct {
	Pos int   `json:"pos"`
	Err error `json:"err"`
}

type LnkTokenErr1 struct {
	ErrToken Token `json:"errToken"`
	Err      error `json:"err"`
}

type Lnk1SingleProgram struct {
	FileAbsPath  string                      `json:"fileAbsPath"`  // The absolute path of the file that the LcProgram is generated from
	FileAbsDir   string                      `json:"fileAbsDir"`   // The absolute path of the directory of which the file resides
	Imports      map[string]Lnk1ImportStmt   `json:"imports"`      // Imports is keyed by ident string without normalization
	TopLevelDefs map[string]Lnk1TopLevelType `json:"topLevelDefs"` // these idents should be normalized to PascalCase
}

type Lnk1ImportStmt struct {
	ImportSrcLocationString         Token  `json:"importSrcLocationString"`
	ImportedAsIdent                 Token  `json:"importedAsIdent"`
	ImportSrcLocationAbsoluteString string `json:"importSrcLocationAbsoluteString"`
}

type Lnk1ProcessedBall struct {
	AllPrograms     map[string]Lnk1SingleProgram `json:"allPrograms"`     // AllPrograms is keyed by the absolute path of the program
	StartingProgram string                       `json:"startingProgram"` // StartingProgram is a string that points to the starting program in the AllPrograms map
}

type Lnk1TopLevelType struct {
	OneofTopLevelStruct *[]Lnk1StructOrEnumLine `json:"topLevelStruct,omitempty"`
	OneofTopLevelEnum   *[]Lnk1StructOrEnumLine `json:"topLevelEnum,omitempty"`
	OneofTopLevelTuple  *[]Lnk1TypeExpr         `json:"topLevelTuple,omitempty"`
	OneofTokenIdent     *Token                  `json:"tokenIdent,omitempty"`
	OneofBuiltin        *BuiltinType            `json:"builtin,omitempty"`
	OneofListof         *Lnk1TypeExpr           `json:"listOf,omitempty"`
	OneofMapof          *Lnk1TypeExpr           `json:"mapOf,omitempty"`
	OneofImported       *Lnk1Imported           `json:"imported,omitempty"`
}

type Lnk1TypeExpr struct {
	OneofMintedIdent *string       `json:"mintedIdent,omitempty"`
	OneofTokenIdent  *Token        `json:"TokenIdent,omitempty"`
	OneofBuiltin     *BuiltinType  `json:"builtin,omitempty"`
	OneofImported    *Lnk1Imported `json:"imported,omitempty"`
	OneofListof      *Lnk1TypeExpr `json:"listOf,omitempty"`
	OneofMapof       *Lnk1TypeExpr `json:"mapOf,omitempty"`
}

type Lnk1Imported struct {
	ImportedIdent                   Token  `json:"importedIdent"`
	ForeignIdent                    Token  `json:"foreignIdent"`
	ImportSrcLocationAbsoluteString string `json:"ImportSrcLocationAbsoluteString"`
}

type Lnk1StructOrEnumLine struct {
	WireName   string       `json:"wireName"`
	ProgName   []string     `json:"progName"`
	TypeExpr   Lnk1TypeExpr `json:"typeExpr"`
	Omittable  bool         `json:"omittable"` // this is ignored when its an enum line
	IsReserved bool         `json:"isReserved"`
}

// start: the starting path, which is a relative path (relative to current wdr) of the file
func lnkGatherSrcFiles(start string) (ret Lnk1ProcessedBall, errFilePath string, lnkErr *LnkErrorUnion) {
	absStart, err := filepath.Abs(start)
	if err != nil {
		return Lnk1ProcessedBall{}, start, &LnkErrorUnion{OneofReadfileErr: err}
	}
	lnkBall := Lnk1ProcessedBall{AllPrograms: make(map[string]Lnk1SingleProgram, 0), StartingProgram: absStart}
	errFilePath, lnkErr = lnkGatherSrcFilesCore(absStart, &lnkBall)
	if lnkErr != nil {
		return Lnk1ProcessedBall{}, errFilePath, lnkErr
	}
	return lnkBall, "", nil
}

// this function expects absolute path
func lnkGatherSrcFilesCore(currentFileAbs string, ball *Lnk1ProcessedBall) (errFilePath string, lnkErr *LnkErrorUnion) {
	lcProg, lnkErr := lnkAbsPathToLcProgram(currentFileAbs)
	if lnkErr != nil {
		return currentFileAbs, lnkErr
	}
	modifiedImports := map[string]Lnk1ImportStmt{}
	fileAbsDir := filepath.Dir(currentFileAbs)
	for lcImportIdent, lcImport := range lcProg.Imports {
		absImportPath := ""
		if filepath.IsAbs(lcImport.ImportSrcLocationString.Data) {
			absImportPath = lcImport.ImportSrcLocationString.Data
		} else {
			absImportPath = filepath.Join(fileAbsDir, lcImport.ImportSrcLocationString.Data)
		}
		modifiedImports[lcImportIdent] = Lnk1ImportStmt{
			ImportSrcLocationString:         lcImport.ImportSrcLocationString,
			ImportedAsIdent:                 lcImport.ImportedAsIdent,
			ImportSrcLocationAbsoluteString: absImportPath,
		}
	}
	modifierdLcTlt := map[string]Lnk1TopLevelType{}
	for lcTltIdent, lcTlt := range lcProg.TopLevelDefs {
		modifierdLcTlt[lcTltIdent] = lnkAddAbsPathToTlt(lcTlt)
	}
	lnkProg := Lnk1SingleProgram{
		FileAbsPath:  currentFileAbs,
		FileAbsDir:   fileAbsDir,
		Imports:      modifiedImports,
		TopLevelDefs: modifierdLcTlt,
	}
	ball.AllPrograms[currentFileAbs] = lnkProg
	for _, lnkImport := range lnkProg.Imports {
		if _, has := ball.AllPrograms[lnkImport.ImportSrcLocationAbsoluteString]; has {
			continue
		}
		errFilePath, lnkErr := lnkGatherSrcFilesCore(lnkImport.ImportSrcLocationAbsoluteString, ball)
		if lnkErr != nil {
			return errFilePath, lnkErr
		}
	}
	return "", nil
}

func lnkAddAbsPathToTlt(lcTlt LcTopLevelType) Lnk1TopLevelType {
	if lcTlt.OneofTopLevelStruct != nil {

	} else if lcTlt.OneofTopLevelEnum != nil {

	} else if lcTlt.OneofTopLevelTuple != nil {

	} else if lcTlt.OneofTokenIdent != nil {

	} else if lcTlt.OneofBuiltin != nil {

	} else if lcTlt.OneofListof != nil {

	} else if lcTlt.OneofMapof != nil {

	} else if lcTlt.OneofImported != nil {

	} else {
		panic("unreachable")
	}
}

func lnkAddAbsPathToTypeExpr(lcTe LcTypeExpr) Lnk1TypeExpr {
	if lcTe.OneofMintedIdent != nil {
		return Lnk1TypeExpr{OneofMintedIdent: lcTe.OneofMintedIdent}
	} else if lcTe.OneofTokenIdent != nil {
		return Lnk1TypeExpr{OneofTokenIdent: lcTe.OneofTokenIdent}
	} else if lcTe.OneofBuiltin != nil {
		return Lnk1TypeExpr{OneofBuiltin: lcTe.OneofBuiltin}
	} else if lcTe.OneofImported != nil {
		
	} else if lcTe.OneofListof != nil {
		t := lnkAddAbsPathToTypeExpr(*lcTe.OneofListof)
		return Lnk1TypeExpr{OneofListof: &t}
	} else if lcTe.OneofMapof != nil {
		t := lnkAddAbsPathToTypeExpr(*lcTe.OneofMapof)
		return Lnk1TypeExpr{OneofMapof: &t}
	} else {
		panic("unreachable")
	}
}

func lnkAbsPathToLcProgram(pathStr string) (LcProgram, *LnkErrorUnion) {
	file, err := readFile(pathStr)
	if err != nil {
		return LcProgram{}, &LnkErrorUnion{OneofReadfileErr: err}
	}
	tokens, errI, err := lexTokenizer(file)
	if err != nil {
		return LcProgram{}, &LnkErrorUnion{OneofLexErr: &LnkLexErr{Pos: errI, Err: err}}
	}
	astProgram, errI, err := rdParseProgram(tokens)
	if err != nil {
		return LcProgram{}, &LnkErrorUnion{OneofParseErr: &LnkTokenErr1{ErrToken: tokens[errI], Err: err}}
	}
	fltProgram := fltFlattenProgram(astProgram)
	errT, err := lcCheckProgram1Of2CheckReservedName(fltProgram)
	if err != nil {
		return LcProgram{}, &LnkErrorUnion{OneofReservedNameErr: &LnkTokenErr1{ErrToken: *errT, Err: err}}
	}
	lcprog, collision, undef := lcCheckProgram2Of2CheckCollisionAndUndefined(fltProgram)
	if len(collision) > 0 {
		return LcProgram{}, &LnkErrorUnion{OneofLcCollisionErr: &collision}
	}
	if len(undef) > 0 {
		return LcProgram{}, &LnkErrorUnion{OneofUndefErr: &undef}
	}
	return lcprog, nil
}

func (lnkErr LnkErrorUnion) ErrToStr() string {
	if lnkErr.OneofReadfileErr != nil {
		return lnkErr.OneofReadfileErr.Error()
	} else if lnkErr.OneofLexErr != nil {
		pos, err := lnkErr.OneofLexErr.Pos, lnkErr.OneofLexErr.Err
		return fmt.Sprintf("Lex error: at %d encountered error: %s", pos, err.Error())
	} else if lnkErr.OneofParseErr != nil {
		errTok, err := lnkErr.OneofParseErr.ErrToken, lnkErr.OneofParseErr.Err
		return fmt.Sprintf("Parse error: at token %s(%s)(%d-%d) encountered error: %s",
			errTok.Kind, errTok.Data, errTok.Start, errTok.End, err.Error())
	} else if lnkErr.OneofReservedNameErr != nil {
		errTok, err := lnkErr.OneofReservedNameErr.ErrToken, lnkErr.OneofReservedNameErr.Err
		return fmt.Sprintf("Reserved name: at token %s(%s)(%d-%d) encountered error: %s",
			errTok.Kind, errTok.Data, errTok.Start, errTok.End, err.Error())
	} else if lnkErr.OneofLcCollisionErr != nil {
		errs := lnkErr.OneofLcCollisionErr
		return fmt.Sprintf("%d collisions found", len(*errs))
	} else if lnkErr.OneofUndefErr != nil {
		errs := lnkErr.OneofUndefErr
		return fmt.Sprintf("%d undefs found", len(*errs))
	} else {
		panic("unreachable")
	}
}
