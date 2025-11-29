package main

import (
	"strings"
	"unicode"
)

func isAsciiLetterUpper(r byte) bool {
	return r >= 'A' && r <= 'Z'
}

func hfNormalizeIdent(identLike string) []string {
	coll := [][]byte{}
	cur := []byte{}
	for _, c := range []byte(identLike) {
		if (c == '_' || c == '-') && len(cur) > 0 {
			coll = append(coll, cur)
			cur = []byte{}
		} else if c == '_' || c == '-' {
			// do nothing
		} else if isAsciiLetterUpper(c) && len(cur) > 0 {
			coll = append(coll, cur)
			cur = []byte{c}
		} else if isAsciiLetterUpper(c) {
			cur = append(cur, c)
		} else {
			cur = append(cur, c)
		}
	}
	if len(cur) > 0 {
		coll = append(coll, cur)
	}
	ret := []string{}
	for _, cs := range coll {
		ret = append(ret, strings.ToLower(string(cs)))
	}
	return ret
}

func hfNormalizedToPascal(li []string) string {
	li2 := []string{}
	for _, e := range li {
		runes := []rune(e)
		runes[0] = unicode.ToUpper(runes[0])
		li2 = append(li2, string(runes))
	}
	return strings.Join(li2, "")
}

func hfNormalizedToCamel(li []string) string {
	li2 := []string{}
	for i, e := range li {
		runes := []rune(e)
		if i > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		li2 = append(li2, string(runes))
	}
	return strings.Join(li2, "")
}

func hfNormalizedToSnake(li []string) string {
	return strings.Join(li, "_")
}

func hfIsKw(tok Token, ident string) bool {
	return tok.Kind == TokenIdent && tok.Data == ident
}

func hfSkipReservedLnkStructOrEnumLines(lines []LnkStructOrEnumLine) []LnkStructOrEnumLine {
	coll := []LnkStructOrEnumLine{}
	for _, line := range lines {
		if line.IsReserved {
			continue
		}
		coll = append(coll, line)
	}
	return coll
}

func hfSkipReservedFltStructOrEnumLines(lines []FltStructOrEnumLine) []FltStructOrEnumLine {
	coll := []FltStructOrEnumLine{}
	for _, line := range lines {
		if line.IsReserved {
			continue
		}
		coll = append(coll, line)
	}
	return coll
}

func multiWr(write func(int, string), indentCount int, before string, betweenLines []string, after string) {
	if len(betweenLines) == 0 {
		write(indentCount, before+after)
		return
	}
	for i, line := range betweenLines {
		if len(betweenLines) == 1 {
			write(indentCount, before+line+after)
		} else if i == 0 {
			write(indentCount, before+line)
		} else if i == len(betweenLines)-1 {
			write(indentCount, line+after)
		} else {
			write(indentCount, line)
		}
	}
}

func multiWr2(
	write func(int, string),
	indentCount int,
	before string,
	firstMultiLine []string,
	between string,
	secondMultiLine []string,
	after string,
) {
	if len(firstMultiLine) == 1 && len(secondMultiLine) == 1 {
		write(indentCount, before+firstMultiLine[0]+between+secondMultiLine[0]+after)
	} else if len(firstMultiLine) == 1 {
		multiWr(write, indentCount, before+firstMultiLine[0]+between, secondMultiLine, after)
	} else if len(secondMultiLine) == 1 {
		multiWr(write, indentCount, before, firstMultiLine, between+secondMultiLine[0]+after)
	} else if len(firstMultiLine) == 0 && len(secondMultiLine) == 0 {
		write(indentCount, before+between+after)
	} else if len(firstMultiLine) == 0 {
		multiWr(write, indentCount, before+between, secondMultiLine, after)
	} else if len(secondMultiLine) == 0 {
		multiWr(write, indentCount, before, firstMultiLine, between+after)
	} else {
		for i, line := range firstMultiLine {
			if i == 0 {
				write(indentCount, before+line)
			} else if i != len(firstMultiLine)-1 {
				write(indentCount, line)
			}
		}
		write(indentCount, firstMultiLine[len(firstMultiLine)-1]+between+secondMultiLine[0])
		for i, line := range secondMultiLine {
			if i != 0 {
				write(indentCount, line)
			} else if i == len(secondMultiLine)-1 {
				write(indentCount, line+after)
			}
		}
	}
}
