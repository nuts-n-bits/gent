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

type LnkSingleProgram struct {
	FileAbsPath  string                    `json:"fileAbsPath"`  // The absolute path of the file that the LcProgram is generated from
	FileAbsDir   string                    `json:"fileAbsDir"`   // The absolute path of the directory of which the file resides
	Imports      map[string]LnkImport      `json:"imports"`      // Imports is keyed by ident string without normalization
	TopLevelDefs map[string]LcTopLevelType `json:"topLevelDefs"` // these idents should be normalized to PascalCase
}

type LnkImport struct {
	ImportSrcLocationString         Token  `json:"importSrcLocationString"`
	ImportedAsIdent                 Token  `json:"importedAsIdent"`
	ImportSrcLocationAbsoluteString string `json:"importSrcLocationAbsoluteString"`
}

type LnkProcessedBall struct {
	AllPrograms     map[string]LnkSingleProgram `json:"allPrograms"`     // AllPrograms is keyed by the absolute path of the program
	StartingProgram string                      `json:"startingProgram"` // StartingProgram is a string that points to the starting program in the AllPrograms map
}
// start: the starting path, which is a relative path (relative to current wdr) of the file
func lnkGatherSrcFiles(start string) (ret LnkProcessedBall, errFilePath string, lnkErr *LnkErrorUnion) {
	absStart, err := filepath.Abs(start)
	if err != nil {
		return LnkProcessedBall{}, start, &LnkErrorUnion{OneofReadfileErr: err}
	}
	lnkBall := LnkProcessedBall{AllPrograms: make(map[string]LnkSingleProgram, 0)}
	errFilePath, lnkErr = lnkGatherSrcFilesCore(absStart, &lnkBall)
	if lnkErr != nil {
		return LnkProcessedBall{}, errFilePath, lnkErr
	}
	return lnkBall, "", nil
}

// this function expects absolute path
func lnkGatherSrcFilesCore(currentFileAbs string, ball *LnkProcessedBall) (errFilePath string, lnkErr *LnkErrorUnion) {
	lcProg, lnkErr := lnkAbsPathToLcProgram(currentFileAbs)
	if lnkErr != nil {
		return currentFileAbs, lnkErr
	}
	modifiedImports := map[string]LnkImport{}
	fileAbsDir := filepath.Dir(currentFileAbs)
	for lcImportIdent, lcImport := range lcProg.Imports {
		absImportPath := ""
		if filepath.IsAbs(lcImport.ImportSrcLocationString.Data) {
			absImportPath = lcImport.ImportSrcLocationString.Data
		} else {
			absImportPath = filepath.Join(fileAbsDir, lcImport.ImportSrcLocationString.Data)
		}
		modifiedImports[lcImportIdent] = LnkImport{
			ImportSrcLocationString:         lcImport.ImportSrcLocationString,
			ImportedAsIdent:                 lcImport.ImportedAsIdent,
			ImportSrcLocationAbsoluteString: absImportPath,
		}
	}
	lnkProg := LnkSingleProgram{
		FileAbsPath:  currentFileAbs,
		FileAbsDir:   fileAbsDir,
		Imports:      modifiedImports,
		TopLevelDefs: lcProg.TopLevelDefs,
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