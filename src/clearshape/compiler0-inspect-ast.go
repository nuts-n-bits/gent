package main

import (
	"encoding/json"
)

func (ast *AstProgram) DebugString() string {
	a, _ := json.Marshal(ast)
	return string(a)
}