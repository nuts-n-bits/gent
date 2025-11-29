package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed runtime/lib-typescript.ts
var cgtsLib string

func cgtsProgram(lnkProgram LnkProgram, indent string, newline string) string {
	var b strings.Builder
	write := func(indentCount int, s string) {
		b.WriteString(strings.Repeat(indent, indentCount) + s + newline)
	}
	for topIdent, typeExpr := range lnkProgram.Types {
		topIdent = cgtsBannedWordsMangle(topIdent)
		typeExprLines := cgtsTypeExpr(typeExpr, indent)
		multiWr(write, 0, "export type __"+topIdent+" = ", typeExprLines, "")
		write(0, "")
	}
	for topIdent, typeExpr := range lnkProgram.Types {
		topIdent = cgtsBannedWordsMangle(topIdent)
		write(0, "export class "+topIdent+" {")
		write(1, "static fromJson(a: string): __"+topIdent+" | Error {")
		write(2, "try { ")
		write(3, "const obj = JSON.parse(a);")
		write(3, "return this.fromJsonCore(obj);")
		write(2, "} catch(e) {")
		write(3, `if (!(e instanceof Error)) { return new Error("caught non error"); }`)
		write(3, `return e;`)
		write(2, "}")
		write(2, "")
		write(1, "}")
		write(1, "")
		write(1, "static fromJsonCore(a: $J): __"+topIdent+" | Error {")
		typeParser := cgtsTypeParserJson(typeExpr, indent)
		multiWr(write, 2, "const parser = ", typeParser, ";")
		write(2, "return parser(a);")
		write(1, "}")
		write(1, "")
		write(1, "static toJsonCore(a: __"+topIdent+"): $J {")
		typeWriter := cgtsTypeWriterJson(typeExpr, indent)
		multiWr(write, 2, "const writer: (a: __"+topIdent+") => $J = ", typeWriter, ";")
		write(2, "return writer(a);")
		write(1, "}")
		write(1, "")
		write(1, "static toJson(a: __"+topIdent+"): string {")
		write(2, "return JSON.stringify(this.toJsonCore(a));")
		write(1, "}")
		write(0, "}")
		write(0, "")
	}
	write(0, cgtsLib)
	return b.String()
}

func cgtsTypeParserJson(typeExpr LnkTypeExpr, indent string) []string {
	if typeExpr.OneofBuiltin != nil {
		return cgtsBuiltinParserJson(*typeExpr.OneofBuiltin, indent)
	} else if typeExpr.OneofEnum != nil {
		return cgtsEnumParserJson(*typeExpr.OneofEnum, indent)
	} else if typeExpr.OneofListof != nil {
		return cgtsListParserJson(*typeExpr.OneofListof, indent)
	} else if typeExpr.OneofMapof != nil {
		return cgtsMapParserJson(*typeExpr.OneofMapof, indent)
	} else if typeExpr.OneofStruct != nil {
		return cgtsStructParserJson(*typeExpr.OneofStruct, indent)
	} else if typeExpr.OneofTokenIdent != nil {
		return []string{fmt.Sprintf("(a: $J) => %s.fromJsonCore(a)", typeExpr.OneofTokenIdent.Data)}
	} else if typeExpr.OneofTuple != nil {
		return cgtsTupleParserJson(*typeExpr.OneofTuple, indent)
	} else {
		panic("unreachable")
	}
	//return []string{"UNIMPLEMENTED"}
}

func cgtsTypeWriterJson(typeExpr LnkTypeExpr, indent string) []string {
	if typeExpr.OneofBuiltin != nil {
		return cgtsBuiltinWriterJson(*typeExpr.OneofBuiltin, indent)
	} else if typeExpr.OneofEnum != nil {
		return cgtsEnumWriterJson(*typeExpr.OneofEnum, indent)
	} else if typeExpr.OneofListof != nil {
		return cgtsListWriterJson(*typeExpr.OneofListof, indent)
	} else if typeExpr.OneofMapof != nil {
		return cgtsMapWriterJson(*typeExpr.OneofMapof, indent)
	} else if typeExpr.OneofStruct != nil {
		return cgtsStructWriterJson(*typeExpr.OneofStruct, indent)
	} else if typeExpr.OneofTokenIdent != nil {
		return []string{fmt.Sprintf("a => %s.toJsonCore(a)", typeExpr.OneofTokenIdent.Data)}
	} else if typeExpr.OneofTuple != nil {
		return cgTsTupleWriterJson(*typeExpr.OneofTuple, indent)
	} else {
		panic("unreachable")
	}
	//return []string{"UNIMPLEMENTED"}
}

func cgtsBuiltinWriterJson(builtinType BuiltinType, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if builtinType == BuiltinTypeString && true {
		write(0, `$writeString`)
	} else if builtinType == BuiltinTypeBoolean {
		write(0, `$writeBoolean`)
	} else if builtinType == BuiltinTypeInt64 {
		write(0, `$writeI64`)
	} else if builtinType == BuiltinTypeUint64 {
		write(0, `$writeU64`)
	} else if builtinType == BuiltinTypeFloat64 {
		write(0, `$writeF64`)
	} else if builtinType == BuiltinTypeNull {
		write(0, `$writeNull`)
	} else if builtinType == BuiltinTypeBinary {
		write(0, `$writeBinary`)
	} else {
		panic("unreachable")
	}
	return b
}

func cgtsStructWriterJson(lines []LnkStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "a => {")
	linesWithoutReserved := hfSkipReservedLnkStructOrEnumLines(lines)
	for _, line := range linesWithoutReserved {
		pascalIdent := cgtsBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		innerWriter := cgtsTypeWriterJson(line.TypeExpr, indent)
		innerType := cgtsTypeExpr(line.TypeExpr, indent)
		multiWr2(write, 1, "const wr"+pascalIdent+": (a: ", innerType, ") => $J = ", innerWriter, ";")
	}
	write(1, "const ret: $J = {}")
	for _, line := range linesWithoutReserved {
		pascalIdent := cgtsBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		camleIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		if line.Omittable {
			write(1, fmt.Sprintf("if (a.%s !== undefined) { ret[%s] = wr%s(a.%s); }", camleIdent, cgtsEncodeString(line.WireName), pascalIdent, camleIdent))
		} else {
			write(1, fmt.Sprintf("ret[%s] = wr%s(a.%s);", cgtsEncodeString(line.WireName), pascalIdent, camleIdent))
		}
	}
	write(1, "return ret;")
	write(0, "}")
	return b
}

func cgtsBuiltinParserJson(builtinType BuiltinType, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if builtinType == BuiltinTypeString && true {
		write(0, `$parseString`)
	} else if builtinType == BuiltinTypeBoolean {
		write(0, `$parseBoolean`)
	} else if builtinType == BuiltinTypeInt64 {
		write(0, `$parseI64`)
	} else if builtinType == BuiltinTypeUint64 {
		write(0, `$parseU64`)
	} else if builtinType == BuiltinTypeFloat64 {
		write(0, `$parseF64`)
	} else if builtinType == BuiltinTypeNull {
		write(0, `$parseNull`)
	} else if builtinType == BuiltinTypeBinary {
		write(0, `$parseBinary`)
	} else {
		panic("unreachable")
	}
	return b
}

func cgtsStructParserJson(lines []LnkStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "(a: $J) => {")
	write(1, `if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing struct"); }`)
	write(1, `const copycat = { ...a };`)
	linesWithoutReserved := hfSkipReservedLnkStructOrEnumLines(lines)
	write(1, "// for each field: create parsers")
	for _, line := range linesWithoutReserved {
		typeParser := cgtsTypeParserJson(line.TypeExpr, indent)
		pascalIdent := cgtsBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		multiWr(write, 1, "const parser"+pascalIdent+" = ", typeParser, ";")
	}
	write(1, "// for required fields only: check presence")
	for _, line := range linesWithoutReserved {
		camelIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		fieldAccessor := "copycat[" + cgtsEncodeString(line.WireName) + "]"
		if !line.Omittable {
			write(1, fmt.Sprintf(
				`if (%s === undefined) { return new Error("required field '%s' (wire name '" + %s + "') is undefined") }`,
				fieldAccessor, camelIdent, cgtsEncodeString(line.WireName)),
			)
		}
	}
	write(1, "// for each field: parse, respecting requiredness, early return on error")
	for _, line := range linesWithoutReserved {
		pascalIdent := cgtsBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		camelIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		fieldAccessor := "copycat[" + cgtsEncodeString(line.WireName) + "]"
		if line.Omittable {
			write(1, fmt.Sprintf(
				"const parsed%s = %s === undefined ? undefined : parser%s(%s);",
				pascalIdent, fieldAccessor, pascalIdent, fieldAccessor),
			)
		} else {
			write(1, fmt.Sprintf("const parsed%s = parser%s(%s);", pascalIdent, pascalIdent, fieldAccessor))
		}
		write(1, fmt.Sprintf(
			`if (parsed%s instanceof Error) { return new Error("error when parsing field %s (wire name '" + %s + "')", { cause: parsed%s }); }`,
			pascalIdent, camelIdent, cgtsEncodeString(line.WireName), pascalIdent,
		))
	}
	write(1, "// for each field: delete field from copycat object")
	for _, line := range linesWithoutReserved {
		fieldAccessor := "copycat[" + cgtsEncodeString(line.WireName) + "]"
		write(1, "delete "+fieldAccessor+";")
	}
	write(1, `if (Object.keys(copycat).length > 0) { return new Error("unknown fields present: " + Object.keys(copycat).join(", ")); }`)
	write(1, "return {")
	for _, line := range linesWithoutReserved {
		camelIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		pascalIdent := cgtsBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		write(2, camelIdent+`: parsed`+pascalIdent+", ")
	}
	write(1, "}")
	write(0, "}")
	return b
}

func cgtsEnumParserJson(lines []LnkStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	typeLines := cgtsTypeEnum(lines, indent)
	write(0, "(a: $J) => {")
	multiWr(write, 1, "type retType = ", typeLines, ";")
	write(1, `if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing enum"); }`)
	write(1, `const entries = Object.entries(a);`)
	write(1, `if (entries.length !== 1) { return new Error("enum values must contain exactly 1 field"); } `)
	write(1, `const [k, v] = entries[0]!;`)
	linesWithoutReserved := hfSkipReservedLnkStructOrEnumLines(lines)
	write(1, "switch (k) {")
	for _, line := range linesWithoutReserved {
		typeParser := cgtsTypeParserJson(line.TypeExpr, indent)
		pascalIdent := cgtsBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		camelIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		write(1, "case "+cgtsEncodeString(line.WireName)+": ")
		multiWr(write, 2, "const parser"+pascalIdent+" = ", typeParser, ";")
		write(2, fmt.Sprintf("const parsed%s = parser%s(v);", pascalIdent, pascalIdent))
		write(2, fmt.Sprintf("return { %s: parsed%s } as retType; ", camelIdent, pascalIdent))
		write(1, "break;")
	}
	write(1, "default: ")
	write(2, `return new Error("unknown variant name while parsing enum, expected one of " + [`)
	for _, line := range linesWithoutReserved {
		camelIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		write(3, fmt.Sprintf(`"%s ('" + %s + "')",`, camelIdent, cgtsEncodeString(line.WireName)))
	}
	write(2, `].join(" / "));`)
	write(1, "}")
	write(0, "}")
	return b
}

func cgtsEnumWriterJson(lines []LnkStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	typeLines := cgtsTypeEnum(lines, indent)
	write(0, "a => {")
	multiWr(write, 1, "type retType = ", typeLines, ";")
	linesWithoutReserved := hfSkipReservedLnkStructOrEnumLines(lines)
	for _, line := range linesWithoutReserved {
		camelIdent := cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		write(1, fmt.Sprintf(`if ("%s" in a) {`, camelIdent))
		multiWr(write, 2, "const writerInner = ", cgtsTypeWriterJson(line.TypeExpr, indent), ";")
		write(2, fmt.Sprintf(`return { %s: writerInner(a.%s) };`, cgtsEncodeString(line.WireName), camelIdent))
		write(1, "}")
	}
	write(1, "return $never(a);")
	write(0, "}")
	return b
}

func cgtsListParserJson(innerType LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "(a: $J) => {")
	write(1, `if (!(a instanceof Array)) { return new Error("expected array while parsing list"); }`)
	innerTypeDecl := cgtsTypeExpr(innerType, indent)
	multiWr(write, 1, "const coll = [] as ", innerTypeDecl, "[];")
	innerParser := cgtsTypeParserJson(innerType, indent)
	multiWr(write, 1, "const parser = ", innerParser, ";")
	write(1, `for (const elem of a) {`)
	write(2, "const parsed = parser(elem);")
	write(2, `if (parsed instanceof Error) { return new Error("failed to parse list inner type", { cause: parsed }); } `)
	write(2, "coll.push(parsed);")
	write(1, `}`)
	write(1, `return coll;`)
	write(0, "}")
	return b
}

func cgtsListWriterJson(innerType LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, `a => {`)
	write(1, `const coll = [] as $J[];`)
	write(1, `for (const elem of a) {`)
	multiWr2(write, 2, `const innerWriter: (a: `, cgtsTypeExpr(innerType, indent), `) => $J = `, cgtsTypeWriterJson(innerType, indent), `;`)
	write(2, `coll.push(innerWriter(elem));`)
	write(1, `}`)
	write(1, `return coll;`)
	write(0, `}`)
	return b
}

func cgtsMapParserJson(innerType LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "(a: $J) => {")
	write(1, `if (typeof a !== "object" || a === null || a instanceof Array) { return new Error("expected object when parsing map"); }`)
	innerTypeDecl := cgtsTypeExpr(innerType, indent)
	multiWr(write, 1, "const coll = {} as { [i: string]: ", innerTypeDecl, " };")
	innerParser := cgtsTypeParserJson(innerType, indent)
	multiWr(write, 1, "const parser = ", innerParser, ";")
	write(1, `for (const k in a) {`)
	write(2, "const parsed = parser(a[k]!);")
	write(2, `if (parsed instanceof Error) { return new Error("failed to parse map's inner type", { cause: parsed }); } `)
	write(2, "coll[k] = parsed;")
	write(1, `}`)
	write(1, `return coll;`)
	write(0, "}")
	return b
}

func cgtsMapWriterJson(innerType LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, `a => {`)
	write(1, `const coll = {} as { [_: string]: $J };`)
	write(1, `for (const k in a) {`)
	multiWr2(write, 2, `const innerWriter: (a: `, cgtsTypeExpr(innerType, indent), `) => $J = `, cgtsTypeWriterJson(innerType, indent), `;`)
	write(2, `coll[k] = innerWriter(a[k]!);`)
	write(1, `}`)
	write(1, `return coll;`)
	write(0, `}`)
	return b
}

func cgtsTupleParserJson(innerTypes []LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	retTypeLines := cgtsTypeTuple(innerTypes, indent)
	write(0, "(a: $J) => {")
	write(1, `if (!(a instanceof Array)) { return new Error("expected array when parsing tuple"); }`)
	write(1, `if (a.length !== `+strconv.Itoa(len(innerTypes))+`) { return new Error("wrong tuple length"); }`)
	for i, innerType := range innerTypes {
		iStr := strconv.Itoa(i)
		typeParser := cgtsTypeParserJson(innerType, indent)
		multiWr(write, 1, "const parser"+iStr+" = ", typeParser, ";")
	}
	for i := range innerTypes {
		iStr := strconv.Itoa(i)
		write(1, fmt.Sprintf("const parsed%s = parser%s(a[%s]!);", iStr, iStr, iStr))
		write(1, fmt.Sprintf(
			`if (parsed%s instanceof Error) { return new Error("failed to parse item #%s in tuple", { cause: parsed%s }); }`,
			iStr, iStr, iStr,
		))
	}
	write(1, "return [")
	for i := range innerTypes {
		iStr := strconv.Itoa(i)
		write(2, "parsed"+iStr+",")
	}
	multiWr(write, 1, "] as ", retTypeLines, ";")
	write(0, "}")
	return b
}

func cgTsTupleWriterJson(innerTypes []LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, `a => {`)
	for i, innerType := range innerTypes {
		iStr := strconv.Itoa(i)
		typeWriter := cgtsTypeWriterJson(innerType, indent)
		multiWr(write, 1, "const writer"+iStr+" = ", typeWriter, ";")
	}
	for i := range innerTypes {
		iStr := strconv.Itoa(i)
		write(1, fmt.Sprintf("const written%s = writer%s(a[%s]!);", iStr, iStr, iStr))
	}
	write(1, "return [")
	for i := range innerTypes {
		iStr := strconv.Itoa(i)
		write(2, "written"+iStr+",")
	}
	write(1, "];")
	write(0, `}`)
	return b
}

func cgtsTypeExpr(lnkType LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if lnkType.OneofBuiltin != nil {
		write(0, cgtsTypeBuiltIn(*lnkType.OneofBuiltin))
	} else if lnkType.OneofEnum != nil {
		innerLines := cgtsTypeEnum(*lnkType.OneofEnum, indent)
		for _, innerLine := range innerLines {
			write(0, innerLine)
		}
	} else if lnkType.OneofListof != nil {
		innerLines := cgtsTypeExpr(*lnkType.OneofListof, indent)
		for i, innerLine := range innerLines {
			if i == len(innerLines)-1 {
				write(0, innerLine+"[]")
			} else {
				write(0, innerLine)
			}
		}
	} else if lnkType.OneofMapof != nil {
		innerLines := cgtsTypeExpr(*lnkType.OneofMapof, indent)
		multiWr(write, 0, "{ [_: string]: ", innerLines, " }")
	} else if lnkType.OneofStruct != nil {
		innerLines := cgtsTypeStruct(*lnkType.OneofStruct, indent)
		for _, innerLine := range innerLines {
			write(0, innerLine)
		}
	} else if lnkType.OneofTokenIdent != nil {
		write(0, "__"+lnkType.OneofTokenIdent.Data)
	} else if lnkType.OneofTuple != nil {
		innerLines := cgtsTypeTuple(*lnkType.OneofTuple, indent)
		for _, innerLine := range innerLines {
			write(0, innerLine)
		}
	} else {
		panic("unreachable")
	}
	return b
}

func cgtsTypeBuiltIn(builtinType BuiltinType) string {
	if builtinType == BuiltinTypeString && true {
		return "string"
	} else if builtinType == BuiltinTypeBoolean {
		return "boolean"
	} else if builtinType == BuiltinTypeInt64 {
		return "bigint"
	} else if builtinType == BuiltinTypeUint64 {
		return "bigint"
	} else if builtinType == BuiltinTypeFloat64 {
		return "number"
	} else if builtinType == BuiltinTypeNull {
		return "null"
	} else if builtinType == BuiltinTypeBinary {
		return "Uint8Array"
	} else {
		panic("unreachable")
	}
}

func cgtsTypeEnum(lines []LnkStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "{")
	for i, line := range lines {
		if i != 0 {
			write(0, "} | {")
		}
		innerLines := cgtsTypeExpr(line.TypeExpr, indent)
		write(1, cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))+": "+innerLines[0])
		for _, innerLine := range innerLines[1:] {
			write(1, innerLine)
		}
	}
	write(0, "}")
	return b
}

func cgtsTypeStruct(lines []LnkStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "{")
	for _, line := range lines {
		innerLines := cgtsTypeExpr(line.TypeExpr, indent)
		if line.Omittable {
			multiWr(write, 1, cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))+"?: undefined | ", innerLines, ",")
		} else {
			multiWr(write, 1, cgtsBannedWordsMangle(hfNormalizedToCamel(line.ProgName))+": ", innerLines, ",")
		}
	}
	write(0, "}")
	return b
}

func cgtsTypeTuple(types []LnkTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if len(types) == 0 {
		write(0, "[]")
	} else if len(types) == 1 {
		innerLines := cgtsTypeExpr(types[0], indent)
		multiWr(write, 0, "[", innerLines, "]")
	} else {
		write(0, "[")
		for _, tupleType := range types {
			innerLines := cgtsTypeExpr(tupleType, indent)
			multiWr(write, 1, "", innerLines, ",")
		}
		write(0, "]")
	}
	return b
}

func cgtsEncodeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, `"`, "\\\"")
	s = strings.ReplaceAll(s, `'`, "\\'")
	return `"` + s + `"`
}

func cgtsBannedWordsMangle(ident string) string {
	switch ident {
	case "break", "case", "catch", "class", "const", "continue", "debugger", "default", "delete", "do", "else", "enum", "export", "extends", "false", "finally", "for", "function", "if", "import", "in", "instanceof", "new", "null", "return", "super", "switch", "this", "throw", "true", "try", "typeof", "var", "void", "while", "with", "yield":
		return ident + "_"
	case "JSON", "Json", "json", "Error", "Object", "Array", "BigInt", "Number", "Uint8Array":
		return ident + "_"
	default:
		return ident
	}
}

