package goPackageNamePlaceholder

import (
	"fmt"
	"slices"
	"strings"
)

type Command struct {
	command string
	args []string
	options []string // always even number of elements arranged as (k, v, k, v, ...)
}

type TokenKind int

const (
	TokenNonQuotedString TokenKind = iota
	TokenQuotedString
	TokenLineBreak
)

type Token struct {
	data string
	kind TokenKind
}

// All printable characters on ANSI keyboard, less backtick (`), apos ('), quote ("), and backslash (\).
var NONQUOTE_CHARSET = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDFEGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()-_=+[{]}|;:,<.>/?");

func tokenizer(src string) ([][]Token, int, error) {
	coll := [][]Token{};
	cur := []Token{};
	i := 0;
	strCannotBeginAt := -1;
	for {
		if i >= len(src) {
			return coll, i, nil;
		} else if src[i] == '\n' {
			if len(cur) > 0 {
				coll = append(coll, cur);
				cur = []Token{};
			}
			i += 1;
		} else if src[i] == '\r' || src[i] == '\t' || src[i] == ' ' {
			i += 1;
		} else if src[i] == '"' {
			if i == strCannotBeginAt { 
				return [][]Token{}, i, fmt.Errorf("quoted string term cannot appear back-to-back with a previous term"); 
			}
			str, newI, err := consumeQuoted(src, '"', i+1);
			if err != nil { 
				return [][]Token{}, newI, err; 
			}
			cur = append(cur, Token{ data: str, kind: TokenQuotedString });
			strCannotBeginAt = newI;
			i = newI;
		} else if slices.Contains(NONQUOTE_CHARSET, src[i]) {
			if i == strCannotBeginAt { 
				return [][]Token{}, i, fmt.Errorf("non-quoted string term cannot appear back-to-back with a previous term"); 
			}

			str, newI := consumeNonquoted(src, i);
			cur = append(cur, Token{ data: str, kind: TokenNonQuotedString} );
			strCannotBeginAt = newI;
			i = newI;
		} else {
			return [][]Token{}, i, fmt.Errorf("unexpected character");
		}
	}
}

func consumeQuoted(src string, delim byte, i int) (string, int, error) {
	var b strings.Builder;
	for {
		if i >= len(src) { 
			return "", i, fmt.Errorf("unexpected eof while consuming quoted");
		} else if (src[i] == '\\') {
			if i+1 >= len(src) {
				return "", i, fmt.Errorf("unexpected eof while consuming escape sequence");
			}
			switch src[i+1] {
			case 'n' : 
				b.WriteString("\n");
				i += 2;
			case 'r' : 
				b.WriteString("\r");
				i += 2;
			case '\\': 
				b.WriteString("\\");
				i += 2;
			case 't' : 
				b.WriteString("\t");
				i += 2;
			case '"': 
				b.WriteString("\"");
				i += 2;
			default: 
				return "", i, fmt.Errorf("unexpected escape sequence while consuming quoted string");
			}
		} else if src[i] == delim { 
			return b.String(), i+1, nil; 
		} else { 
			b.WriteByte(src[i]);
			i = i+1; 
		}
	}
}

func consumeNonquoted(src string, i int) (string, int) {
	var b strings.Builder;
	for {

		if i < len(src) && slices.Contains(NONQUOTE_CHARSET, src[i]) { 
			b.WriteByte(src[i]);
			i = i+1; 
		} else { 
			return b.String(), i; 
		}
	}
}

func parseOne(toks []Token) (Command, int, error) {
	command, i, positionalMode := Command{}, 0, false;
	if len(toks) == 0 {
		return Command{}, i, fmt.Errorf("empty token stream");
	} else {
		command.command = toks[i].data;
	}
	for {
		if i >= len(toks) {
			return command, i, nil;
		} else if toks[i].kind == TokenQuotedString {
			command.args = append(command.args, toks[i].data);
			i += 1;
		} else if toks[i].data[0] != '-' || positionalMode {
			command.args = append(command.args, toks[i].data);
			i += 1;
		} else if toks[i].data == "--" {
			positionalMode = true;
			i += 1;
		} else if i+1 >= len(toks) || (toks[i+1].kind == TokenNonQuotedString && toks[i+1].data[0] == '-') {
			command.options = append(command.options, toks[i].data, "");
			i += 1;
		} else {
			command.options = append(command.options, toks[i].data, toks[i+1].data);
			i += 2;
		}
	}
}

func CoreDynParse(src string) ([]Command, int, error) {
	tokenss, i, err := tokenizer(src);
	coll := []Command{};
	if err != nil {
		return []Command{}, i, err;
	}
	for _, tokens := range tokenss {
		command, _, err := parseOne(tokens);
		if err != nil {
			return []Command{}, 0, err;
		}
		coll = append(coll, command);
	}
	return coll, 0, nil;
}

func CoreDynEncode(cmd Command) string {
	var b strings.Builder;
	b.WriteString(encodeStr(cmd.command));
	for _, arg := range cmd.args {
		b.WriteString(" ");
		b.WriteString(encodeStr(arg));
	}
	for i:=0; i<len(cmd.options); i+=2 {
		b.WriteString(" ");
		b.WriteString(cmd.options[i]);
		b.WriteString(" ");
		b.WriteString(cmd.options[i+1]);
	}
	return b.String();
}

func nqtest(tested string) bool {
	for _, byte := range []byte(tested) {
		if !slices.Contains(NONQUOTE_CHARSET, byte) { 
			return false; 
		}
	}
	return true;
}

func encodeStr(s string) string {
	if len(s) > 0 && len(s) < 50 && s[0]!= '-' && nqtest(s) {
		return s;
	}
	t := strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"");
	return "\"" + strings.ReplaceAll(strings.ReplaceAll(t, "\r", "\\r"), "\n", "\\n") + "\"";
}

