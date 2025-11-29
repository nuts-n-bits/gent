package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed runtime/lib-golang.go
var cggoLib string

func cggoProgram(fltProgram FltProgram, indent string, newline string, packageName string) string {
	var b strings.Builder
	write := func(indentCount int, s string) {
		b.WriteString(strings.Repeat(indent, indentCount) + s + newline)
	}
	write(0, "package "+packageName+"\n")
	libNoPkgName := strings.Replace(cggoLib, "package WHATPACKAGENAME\n", "", 1)
	libImport := strings.Join(strings.Split(libNoPkgName, "\n")[0:6], "\n")
	libRest := strings.Join(strings.Split(libNoPkgName, "\n")[6:], "\n")
	write(0, libImport)
	for topIdent, typeExpr := range fltProgram.TopTypes {
		topIdent = cggoBannedWordsMangle(topIdent)
		typeExprLines := cggoTopLevelTypeExpr(typeExpr, indent)
		multiWr(write, 0, "type "+topIdent+" ", typeExprLines, "")
		write(0, "")
	}
	for topIdent, tlt := range fltProgram.TopTypes {
		topIdent = cgtsBannedWordsMangle(topIdent)
		typeParser := cggoTopTypeParserJson(tlt, indent)
		//typeWriter := cggoTopTypeWriterJson(typeExpr, indent)
		multiWr(write, 0, "func (a *"+topIdent+") fromJsonCore (b interface{}) _err ", typeParser, "")
		write(0, "")
	}
	write(0, libRest)
	return b.String()
}

func cggoTopTypeParserJson(topTypeExpr FltTopLevelTypeExpr, indent string) []string {
	if topTypeExpr.OneofBuiltin != nil {
		return cggoBuiltinParserJson(*topTypeExpr.OneofBuiltin, indent)
	} else if topTypeExpr.OneofEnum != nil {
		return cggoTopEnumParserJson(*topTypeExpr.OneofEnum, indent)
	} else if topTypeExpr.OneofListof != nil {
		return cggoTopListParserJson(*topTypeExpr.OneofListof, indent)
	} else if topTypeExpr.OneofMapof != nil {
		return cggoTopMapParserJson(*topTypeExpr.OneofMapof, indent)
	} else if topTypeExpr.OneofStruct != nil {
		return cggoTopStructParserJson(*topTypeExpr.OneofStruct, indent)
	} else if topTypeExpr.OneofTokenIdent != nil {
		return []string{fmt.Sprintf("(a: $J) => %s.fromJsonCore(a)", topTypeExpr.OneofTokenIdent.Data)}
	} else if topTypeExpr.OneofTuple != nil {
		return cggoTopTupleParserJson(*topTypeExpr.OneofTuple, indent)
	} else {
		panic("unreachable")
	}
	//return []string{"UNIMPLEMENTED"}
}

func cggoBuiltinParserJson(builtinType BuiltinType, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if builtinType == BuiltinTypeString && true {
		write(0, `_parseString`)
	} else if builtinType == BuiltinTypeBoolean {
		write(0, `_parseBoolean`)
	} else if builtinType == BuiltinTypeInt64 {
		write(0, `_parseInt64`)
	} else if builtinType == BuiltinTypeUint64 {
		write(0, `_parseUint64`)
	} else if builtinType == BuiltinTypeFloat64 {
		write(0, `_parseFloat64`)
	} else if builtinType == BuiltinTypeNull {
		write(0, `_parseNull`)
	} else if builtinType == BuiltinTypeBinary {
		write(0, `_parseBinary`)
	} else {
		panic("unreachable")
	}
	return b
}

func cggoTopStructParserJson(lines []FltStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "{")
	write(1, "bMap, err := _parseJsonObject(b)")
	write(1, "if err != nil {")
	write(2, "return err")
	write(1, "}")
	linesWithoutReserved := hfSkipReservedFltStructOrEnumLines(lines)
	for _, line := range linesWithoutReserved {
		camelIdent := cggoBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		innerLines := cggoInnerTypeExpr(line.TypeExpr, indent)
		if line.Omittable {
			multiWr(write, 1, "var "+camelIdent+" *", innerLines, "")
		} else {
			multiWr(write, 1, "var "+camelIdent+" ", innerLines, "")
		}
	}
	for _, line := range linesWithoutReserved {
		camelIdent := cggoBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		pascalIdent := cggoBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		innerTypeParser := cggoInnerTypeParser(line.TypeExpr, indent, "parsed", "err", "t0", "t1", "t2")
		write(1, "if v, ok := bMap["+cggoEncodeString(line.WireName)+"]; ok {")
		write(2, "// parse v into "+camelIdent+" here")
		multiWr(write, 2, "", innerTypeParser, "")
		write(2, "if err != nil {")
		write(3, fmt.Sprintf(`return _newerr("error when parsing field %s (wire name '"+%s+"')", err)`, pascalIdent, cggoEncodeString(line.WireName)))
		write(2, "}")
		if line.Omittable {
			write(2, camelIdent+" = &parsed")
		} else {
			write(2, camelIdent+" = parsed")
		}
		if !line.Omittable {
			write(1, "} else {")
			write(2, fmt.Sprintf(`return _newerr("missing required field %s (wire name '"+%s+"')", nil)`, pascalIdent, cggoEncodeString(line.WireName)))
		}
		write(1, "}")
	}
	for _, line := range linesWithoutReserved {
		camelIdent := cggoBannedWordsMangle(hfNormalizedToCamel(line.ProgName))
		pascalIdent := cggoBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		write(1, "a."+pascalIdent+" = "+camelIdent)
	}
	write(1, "return nil")
	write(0, "}")
	return b
}

func cggoTopEnumParserJson(lines []FltStructOrEnumLine, indent string) []string {
	return []string{"{", indent + "UNIMPLEMENTED", "}"}
}

func cggoTopListParserJson(innerType FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "{")
	listParserLines := cggoListParserJson(innerType, indent)
	multiWr(write, 1, "listParser := ", listParserLines, "")
	write(1, "parsed, err := listParser(b)")
	write(1, "if err != nil {")
	write(2, `return _newerr("error while parsing list", err)`)
	write(1, "}")
	write(1, "*a = parsed")
	write(1, "return nil")
	write(0, "}")
	return b
}

func cggoListParserJson(innerType FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	innerTypeLines := cggoInnerTypeExpr(innerType, indent)
	multiWr(write, 0, "func (a interface{}) ([]", innerTypeLines, ", _err) {")
	write(1, "list, err := _parseJsonList(a)")
	write(1, "if err != nil {")
	multiWr(write, 2, "return []", innerTypeLines, `{}, _newerr("error parsing list", err)`)
	write(1, "}")
	multiWr(write, 1, "ret := []", innerTypeLines, "{}")
	write(1, "for i, v := range list {")
	innerParserLines := cggoInnerTypeParser(innerType, indent, "parsed", "err", "t0", "t1", "t2")
	multiWr(write, 2, "", innerParserLines, "")
	write(2, "if err != nil {")
	write(3, "iStr := strconv.Itoa(i)")
	multiWr(write, 3, "return []", innerTypeLines, `{}, _newerr("error parsing item " + iStr + " inside list", err)`)
	write(2, "}")
	write(2, "ret = append(ret, parsed)")
	write(1, "}")
	write(1, "return ret, nil")
	write(0, "}")
	return b
}

func cggoTopMapParserJson(innerType FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "{")
	mapParserLines := cggoMapParserJson(innerType, indent)
	multiWr(write, 1, "mapParser := ", mapParserLines, "")
	write(1, "parsed, err := mapParser(b)")
	write(1, "if err != nil {")
	write(2, `return _newerr("error while parsing map", err)`)
	write(1, "}")
	write(1, "*a = parsed")
	write(1, "return nil")
	write(0, "}")
	return b
}

func cggoMapParserJson(innerType FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	innerTypeLines := cggoInnerTypeExpr(innerType, indent)
	multiWr(write, 0, "func (a interface{}) (map[string]", innerTypeLines, ", _err) {")
	write(1, "map_, err := _parseJsonObject(a)")
	write(1, "if err != nil {")
	multiWr(write, 2, "return map[string]", innerTypeLines, `{}, _newerr("error parsing map", err)`)
	write(1, "}")
	multiWr(write, 1, "ret := map[string]", innerTypeLines, "{}")
	write(1, "for k, v := range map_ {")
	innerParserLines := cggoInnerTypeParser(innerType, indent, "parsed", "err", "t0", "t1", "t2")
	multiWr(write, 2, "", innerParserLines, "")
	write(2, "if err != nil {")
	multiWr(write, 3, "return map[string]", innerTypeLines, `{}, _newerr("error parsing key '" + k + "' inside map", err)`)
	write(2, "}")
	write(2, "ret[k] = parsed")
	write(1, "}")
	write(1, "return ret, nil")
	write(0, "}")
	return b
}

func cggoTopTupleParserJson(innerTypes []FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "{")
	write(1, "bList, err := _parseJsonList(b)")
	write(1, "if err != nil {")
	write(2, "return err")
	write(1, "}")
	write(1, "if len(bList) != " + strconv.Itoa(len(innerTypes)) + " {")
	write(2, `return _newerr("tuple does not have the correct length", nil)`)
	write(1, "}")
	for i, innerType := range innerTypes {
		iStr := strconv.Itoa(i)
		camelIdent := cggoBannedWordsMangle("field" + iStr)
		innerLines := cggoInnerTypeExpr(innerType, indent)
		multiWr(write, 1, "var "+camelIdent+" ", innerLines, "")
	}
	for i, innerType := range innerTypes {
		iStr := strconv.Itoa(i)
		camelIdent := cggoBannedWordsMangle("field" + iStr)
		pascalIdent := cggoBannedWordsMangle("Field" + iStr)
		innerTypeParser := cggoInnerTypeParser(innerType, indent, "parsed", "err", "t0", "t1", "t2")

		multiWr(write, 1, "", innerTypeParser, "")
		write(1, "if err != nil {")
		write(2, fmt.Sprintf(`return _newerr("error when parsing %s", err)`, pascalIdent))
		write(1, "}")
		write(1, camelIdent+" = parsed")

	}
	for i, _ := range innerTypes {
		iStr := strconv.Itoa(i)
		camelIdent := cggoBannedWordsMangle("field" + iStr)
		pascalIdent := cggoBannedWordsMangle("Field" + iStr)
		write(1, "a."+pascalIdent+" = "+camelIdent)
	}
	write(1, "return nil")
	write(0, "}")
	return b
}

// this function assumes the symbols "t0-t9" "parsed" "err" are available in local scope, and will assume "v" is the input!!
func cggoInnerTypeParser(innerType FltInnerTypeExpr, indent string, freeNameRes, freeNameErr, freeNameT0, freeNameT1, freeNameT2 string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if innerType.OneofBuiltin != nil {
		multiWr(write, 0, freeNameRes+", "+freeNameErr+" := ", cggoBuiltinParserJson(*innerType.OneofBuiltin, indent), "(v)")
	} else if innerType.OneofMintedIdent != nil {
		write(0, freeNameRes+" := "+*innerType.OneofMintedIdent+"{}")
		write(0, freeNameErr+" := parsed.fromJsonCore(v)")
	} else if innerType.OneofTokenIdent != nil {
		write(0, freeNameRes+" := "+innerType.OneofTokenIdent.Data+"{}")
		write(0, freeNameErr+" := parsed.fromJsonCore(v)")
	} else if innerType.OneofListof != nil {
		listParserLines := cggoListParserJson(*innerType.OneofListof, indent)
		multiWr(write, 0, freeNameT0+" := ", listParserLines, "")
		write(0, freeNameRes+", "+freeNameErr+" := "+freeNameT0+"(v)")
	} else if innerType.OneofMapof != nil {
		mapParserLines := cggoMapParserJson(*innerType.OneofMapof, indent)
		multiWr(write, 0, freeNameT0+" := ", mapParserLines, "")
		write(0, freeNameRes+", "+freeNameErr+" := "+freeNameT0+"(v)")
	} else {
		panic("unreachable")
	}
	return b
}

func cggoTopLevelTypeExpr(fltTopLevelTypeExpr FltTopLevelTypeExpr, indent string) []string {
	if fltTopLevelTypeExpr.OneofBuiltin != nil {
		return []string{cggoTypeBuiltIn(*fltTopLevelTypeExpr.OneofBuiltin)}
	} else if fltTopLevelTypeExpr.OneofEnum != nil {
		return cggoTypeEnum(*fltTopLevelTypeExpr.OneofEnum, indent)
	} else if fltTopLevelTypeExpr.OneofListof != nil {
		b := []string{}
		write := func(indentCount int, s string) {
			b = append(b, strings.Repeat(indent, indentCount)+s)
		}
		multiWr(write, 0, "[]", cggoInnerTypeExpr(*fltTopLevelTypeExpr.OneofListof, indent), "")
		return b
	} else if fltTopLevelTypeExpr.OneofMapof != nil {
		b := []string{}
		write := func(indentCount int, s string) {
			b = append(b, strings.Repeat(indent, indentCount)+s)
		}
		multiWr(write, 0, "map[string]", cggoInnerTypeExpr(*fltTopLevelTypeExpr.OneofMapof, indent), "")
		return b
	} else if fltTopLevelTypeExpr.OneofStruct != nil {
		return cggoTypeStruct(*fltTopLevelTypeExpr.OneofStruct, indent)
	} else if fltTopLevelTypeExpr.OneofTokenIdent != nil {
		return []string{cggoBannedWordsMangle(hfNormalizedToPascal(hfNormalizeIdent(fltTopLevelTypeExpr.OneofTokenIdent.Data)))}
	} else if fltTopLevelTypeExpr.OneofTuple != nil {
		return cggoTypeTuple(*fltTopLevelTypeExpr.OneofTuple, indent)
	} else {
		panic("unreachable")
	}
	//return []string{"UNIMPLEMENTED"}
}

func cggoInnerTypeExpr(fltInnerTypeExpr FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	if fltInnerTypeExpr.OneofBuiltin != nil {
		write(0, cggoTypeBuiltIn(*fltInnerTypeExpr.OneofBuiltin))
	} else if fltInnerTypeExpr.OneofMintedIdent != nil {
		write(0, cggoBannedWordsMangle(hfNormalizedToPascal(hfNormalizeIdent(*fltInnerTypeExpr.OneofMintedIdent))))
	} else if fltInnerTypeExpr.OneofTokenIdent != nil {
		write(0, cggoBannedWordsMangle(hfNormalizedToPascal(hfNormalizeIdent(fltInnerTypeExpr.OneofTokenIdent.Data))))
	} else if fltInnerTypeExpr.OneofListof != nil {
		innerLines := cggoInnerTypeExpr(*fltInnerTypeExpr.OneofListof, indent)
		multiWr(write, 0, "[]", innerLines, "")
	} else if fltInnerTypeExpr.OneofMapof != nil {
		innerLines := cggoInnerTypeExpr(*fltInnerTypeExpr.OneofMapof, indent)
		multiWr(write, 0, "map[string]", innerLines, "")
	} else {
		panic("unreachable")
	}
	return b
}

func cggoTypeStruct(lines []FltStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "struct {")
	for _, line := range lines {
		pascalIdent := cggoBannedWordsMangle(hfNormalizedToPascal(line.ProgName))
		innerLines := cggoInnerTypeExpr(line.TypeExpr, indent)
		if line.Omittable {
			multiWr(write, 1, pascalIdent+" *", innerLines, "")
		} else {
			multiWr(write, 1, pascalIdent+" ", innerLines, "")
		}
	}
	write(0, "}")
	return b
}

func cggoTypeEnum(lines []FltStructOrEnumLine, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "struct {")
	for _, line := range lines {
		pascalIdent := cggoBannedWordsMangle("Oneof" + hfNormalizedToPascal(line.ProgName))
		innerLines := cggoInnerTypeExpr(line.TypeExpr, indent)
		multiWr(write, 1, pascalIdent+" *", innerLines, "")
	}
	write(0, "}")
	return b
}

func cggoTypeTuple(innerTypes []FltInnerTypeExpr, indent string) []string {
	b := []string{}
	write := func(indentCount int, s string) {
		b = append(b, strings.Repeat(indent, indentCount)+s)
	}
	write(0, "struct {")
	for i, innerType := range innerTypes {
		iStr := strconv.Itoa(i)
		pascalIdent := "Field" + iStr
		innerLines := cggoInnerTypeExpr(innerType, indent)
		multiWr(write, 1, pascalIdent+" ", innerLines, "")
	}
	write(0, "}")
	return b
}

func cggoTypeBuiltIn(builtinType BuiltinType) string {
	if builtinType == BuiltinTypeString && true {
		return "string"
	} else if builtinType == BuiltinTypeBoolean {
		return "boolean"
	} else if builtinType == BuiltinTypeInt64 {
		return "int64"
	} else if builtinType == BuiltinTypeUint64 {
		return "uint64"
	} else if builtinType == BuiltinTypeFloat64 {
		return "float64"
	} else if builtinType == BuiltinTypeNull {
		return "null"
	} else if builtinType == BuiltinTypeBinary {
		return "[]byte"
	} else {
		panic("unreachable")
	}
}

func cggoBannedWordsMangle(ident string) string {
	switch ident {
	case "break", "case", "chan", "const", "continue", "default", "defer", "else", "fallthrough", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "type", "var":
		return ident + "_"
	case "JSON", "Json", "json", "error", "null":
		return ident + "_"
	default:
		return ident
	}
}

func cggoEncodeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, `"`, "\\\"")
	s = strings.ReplaceAll(s, `'`, "\\'")
	return `"` + s + `"`
}
