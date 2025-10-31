package main

/*
RULES
1. No toplevel identifier conflict allowed (when compared after cononicalization)
2. No wirename conflict allowed (when compared literally)
   - If progname is used as wirename, prog name is converted to snake_case first
3. No progname conflict allowed (when compared after cononicalization - camelCase)
4. Non minted toplevel identifiers mustn't start with "csres0"

*/

type BuiltinType string

const (
	BuiltinTypeString  BuiltinType = "BuiltinTypeString"
	BuiltinTypeBoolean BuiltinType = "BuiltinTypeBoolean"
	BuiltinTypeInt64   BuiltinType = "BuiltinTypeInt64"
	BuiltinTypeUint64  BuiltinType = "BuiltinTypeUint64"
	BuiltinTypeFloat64 BuiltinType = "BuiltinTypeFloat64"
	BuiltinTypeNull    BuiltinType = "BuiltinTypeNull"
)

// func lckCheckProgram(FltProgram) (errT1, errT2 *Token, err error) {

// }