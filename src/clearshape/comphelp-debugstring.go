package main

import (
	"encoding/json"
)

func (ast *AstProgram) DebugString() string {
	a, _ := json.Marshal(ast)
	return string(a)
}

func (flt *FltProgram) DebugString() string {
	a, _ := json.Marshal(flt)
	return string(a)
}

func (lc *LcProgram) DebugString() string{
	a, _ := json.Marshal(lc)
	return string(a)
}

func (lnk *LnkProgram) DebugString() string{
	a, _ := json.Marshal(lnk)
	return string(a)
}