package main

import (
	"fmt"
	"strings"
)

type TokenKind string

const (
	TokenOpenBrace    TokenKind = "TokenOpenBrace"
	TokenCloseBrace   TokenKind = "TokenCloseBrace"
	TokenOpenParen    TokenKind = "TokenOpenParen"
	TokenCloseParen   TokenKind = "TokenCloseParen"
	TokenColon        TokenKind = "TokenColon"
	TokenPipe         TokenKind = "TokenPipe"
	TokenSemicolon    TokenKind = "TokenSemicolon"
	TokenOpenBracket  TokenKind = "TokenOpenBracket"
	TokenCloseBracket TokenKind = "TokenCloseBracket"
	TokenQuestion     TokenKind = "TokenQuestion"
	TokenComma        TokenKind = "TokenComma"
	TokenEq           TokenKind = "TokenEq"
	TokenIdent        TokenKind = "TokenIdent"
	TokenDot          TokenKind = "TokenDot"
	TokenString       TokenKind = "TokenString"
	TokenEof          TokenKind = "TokenEof"
)

type Token struct {
	Kind  TokenKind `json:"kind"`
	Data  string    `json:"data"`
	Start int       `json:"start"`
	End   int       `json:"end"`
}

func lexTokenizer(program string) ([]Token, int, error) {
	tokens := []Token{}
	i := 0
	for {
		if i >= len(program) {
			tokens = append(tokens, Token{Kind: TokenEof, Start: i, End: i})
			return tokens, 0, nil
		} else if program[i] == ' ' || program[i] == '\t' || program[i] == '\n' || program[i] == '\r' {
			i += 1
		} else if program[i] == '/' && program[i+1] == '/' {
			for program[i] != '\n' {
				i += 1
			}
		} else if program[i] == '{' {
			tokens = append(tokens, Token{Kind: TokenOpenBrace, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '}' {
			tokens = append(tokens, Token{Kind: TokenCloseBrace, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '(' {
			tokens = append(tokens, Token{Kind: TokenOpenParen, Start: i, End: i + 1})
			i += 1
		} else if program[i] == ')' {
			tokens = append(tokens, Token{Kind: TokenCloseParen, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '[' {
			tokens = append(tokens, Token{Kind: TokenOpenBracket, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '|' {
			tokens = append(tokens, Token{Kind: TokenPipe, Start: i, End: i + 1})
			i += 1
		} else if program[i] == ']' {
			tokens = append(tokens, Token{Kind: TokenCloseBracket, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '?' {
			tokens = append(tokens, Token{Kind: TokenQuestion, Start: i, End: i + 1})
			i += 1
		} else if program[i] == ':' {
			tokens = append(tokens, Token{Kind: TokenColon, Start: i, End: i + 1})
			i += 1
		} else if program[i] == ';' {
			tokens = append(tokens, Token{Kind: TokenSemicolon, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '.' {
			tokens = append(tokens, Token{Kind: TokenDot, Start: i, End: i + 1})
			i += 1
		} else if program[i] == ',' {
			tokens = append(tokens, Token{Kind: TokenComma, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '=' {
			tokens = append(tokens, Token{Kind: TokenEq, Start: i, End: i + 1})
			i += 1
		} else if program[i] == '"' {
			strContent, newI, err := lexConsumeString(program, i+1)
			if err != nil {
				return tokens, i, err
			}
			tokens = append(tokens, Token{Kind: TokenString, Start: i, End: newI, Data: strContent})
			i = newI
		} else if lexIsIdentStart(program[i]) {
			ident := lexConsumeIdent(program, i)
			i2 := i + len(ident)
			tokens = append(tokens, Token{Kind: TokenIdent, Data: ident, Start: i, End: i2})
			i = i2
		} else {
			return tokens, i, fmt.Errorf("unexpected character at %d", i)
		}
	}
}

func lexIsIdentStart(r byte) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_'
}

func lexIsIdentCont(r byte) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_'
}

func lexConsumeIdent(program string, i int) string {
	identColl := []byte{}
	for {
		if i >= len(program) {
			return string(identColl)
		} else if lexIsIdentCont(program[i]) {
			identColl = append(identColl, byte(program[i]))
			i += 1
		} else {
			return string(identColl)
		}
	}
}

func lexConsumeString(program string, i int) (string, int, error) {
	var b strings.Builder
	for {
		if i >= len(program) {
			return "", i, fmt.Errorf("unclosed string literal")
		} else if program[i] == '"' {
			return b.String(), i + 1, nil
		} else if program[i] == '\n' || program[i] == '\r' {
			return "", i, fmt.Errorf("illegal CR LF characters in string literal")
		} else if program[i] == '\\' {
			if i+1 >= len(program) {
				return "", i, fmt.Errorf("unfinished escape sequence")
			} else if program[i+1] == 'n' {
				b.WriteString("\n")
				i += 2
			} else if program[i+1] == 't' {
				b.WriteString("\t")
				i += 2
			} else if program[i+1] == 'r' {
				b.WriteString("\r")
				i += 2
			} else if program[i+1] == '\\' {
				b.WriteString("\\")
				i += 2
			} else if program[i+1] == '"' {
				b.WriteString("\"")
				i += 2
			} else {
				return "", i + 1, fmt.Errorf("illegal escape sequence at %d", i+1)
			}
		} else {
			b.WriteByte(program[i])
			i += 1
		}
	}
}
