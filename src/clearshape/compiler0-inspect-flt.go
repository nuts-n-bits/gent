package main

import (
	"encoding/json"
)

func (flt *FltProgram) DebugString() string {
	a, _ := json.Marshal(flt)
	return string(a)
}