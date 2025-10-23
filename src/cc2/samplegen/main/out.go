// BEGIN RUNTIME LIBRARY

package main

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

type _TokenKind int

const (
	_TokenNonQuotedString _TokenKind = iota
	_TokenQuotedString
	_TokenLineBreak
)

type _Token struct {
	data string
	kind _TokenKind
}

// All printable characters on ANSI keyboard, less backtick (`), apos ('), quote ("), and backslash (\).
var _NONQUOTE_CHARSET = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDFEGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()-_=+[{]}|;:,<.>/?");

func _tokenizer(src string) ([][]_Token, int, error) {
	coll := [][]_Token{};
	cur := []_Token{};
	i := 0;
	strCannotBeginAt := -1;
	for {
		if i >= len(src) {
			if len(cur) > 0 {
				coll = append(coll, cur);
				cur = []_Token{};
			}
			return coll, i, nil;
		} else if src[i] == '\n' {
			if len(cur) > 0 {
				coll = append(coll, cur);
				cur = []_Token{};
			}
			i += 1;
		} else if src[i] == '\r' || src[i] == '\t' || src[i] == ' ' {
			i += 1;
		} else if src[i] == '"' {
			if i == strCannotBeginAt { 
				return [][]_Token{}, i, fmt.Errorf("quoted string term cannot appear back-to-back with a previous term"); 
			}
			str, newI, err := _consumeQuoted(src, '"', i+1);
			if err != nil { 
				return [][]_Token{}, newI, err; 
			}
			cur = append(cur, _Token{ data: str, kind: _TokenQuotedString });
			strCannotBeginAt = newI;
			i = newI;
		} else if slices.Contains(_NONQUOTE_CHARSET, src[i]) {
			if i == strCannotBeginAt { 
				return [][]_Token{}, i, fmt.Errorf("non-quoted string term cannot appear back-to-back with a previous term"); 
			}

			str, newI := _consumeNonquoted(src, i);
			cur = append(cur, _Token{ data: str, kind: _TokenNonQuotedString} );
			strCannotBeginAt = newI;
			i = newI;
		} else {
			return [][]_Token{}, i, fmt.Errorf("unexpected character");
		}
	}
}

func _consumeQuoted(src string, delim byte, i int) (string, int, error) {
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

func _consumeNonquoted(src string, i int) (string, int) {
	var b strings.Builder;
	for {

		if i < len(src) && slices.Contains(_NONQUOTE_CHARSET, src[i]) { 
			b.WriteByte(src[i]);
			i = i+1; 
		} else { 
			return b.String(), i; 
		}
	}
}

func _parseOne(toks []_Token) (Command, int, error) {
	command, i, positionalMode := Command{}, 0, false;
	if len(toks) == 0 {
		return Command{}, i, fmt.Errorf("empty token stream");
	} else {
		command.command = toks[i].data;
		i += 1;
	}
	for {
		if i >= len(toks) {
			return command, i, nil;
		} else if toks[i].kind == _TokenQuotedString {
			command.args = append(command.args, toks[i].data);
			i += 1;
		} else if toks[i].data[0] != '-' || positionalMode {
			command.args = append(command.args, toks[i].data);
			i += 1;
		} else if toks[i].data == "--" {
			positionalMode = true;
			i += 1;
		} else if i+1 >= len(toks) || (toks[i+1].kind == _TokenNonQuotedString && toks[i+1].data[0] == '-') {
			command.options = append(command.options, toks[i].data, "");
			i += 1;
		} else {
			command.options = append(command.options, toks[i].data, toks[i+1].data);
			i += 2;
		}
	}
}

type _ccCore struct {}

var CcCore = _ccCore{};

func (_ _ccCore) CoreDynParse(src string) ([]Command, int, error) {
	tokenss, i, err := _tokenizer(src);
	coll := []Command{};
	if err != nil {
		return []Command{}, i, err;
	}
	for _, tokens := range tokenss {
		command, _, err := _parseOne(tokens);
		if err != nil {
			return []Command{}, 0, err;
		}
		coll = append(coll, command);
	}
	return coll, 0, nil;
}

func (_ _ccCore) CoreDynEncode(cmd Command) string {
	var b strings.Builder;
	b.WriteString(_encodeStr(cmd.command));
	for _, arg := range cmd.args {
		b.WriteString(" ");
		b.WriteString(_encodeStr(arg));
	}
	if len(cmd.options) % 2 != 0 {
		cmd.options = append(cmd.options, "");
	}
	for i:=0; i<len(cmd.options); i+=2 {
		b.WriteString(" ");
		b.WriteString(cmd.options[i]);
		if cmd.options[i+1] != "" {
			b.WriteString(" ");
			b.WriteString(_encodeStr(cmd.options[i+1]));
		}
	}
	return b.String();
}

func _nqtest(tested string) bool {
	for _, byte := range []byte(tested) {
		if !slices.Contains(_NONQUOTE_CHARSET, byte) { 
			return false; 
		}
	}
	return true;
}

func _encodeStr(s string) string {
	if len(s) > 0 && len(s) < 50 && s[0]!= '-' && _nqtest(s) {
		return s;
	}
	t := strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\"", "\\\"");
	return "\"" + strings.ReplaceAll(strings.ReplaceAll(t, "\r", "\\r"), "\n", "\\n") + "\"";
}

// BEGIN MACHINE GENERATED CODE
type Token struct {
    Args string;
}

func (a *Token) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "token" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a Token) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "token",
        args: args,
        options: options,
    };
}

func (a *Token) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Token) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Proxy struct {
    Args string;
    Enable bool
}

func (a *Proxy) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    var bEnable bool = false;
    var cEnable int = 0;
    if b.command != "proxy" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        case "-enable": 
            bEnable = true;
            cEnable += 1;
        }
    }
    a.Args = bArgs; 
    a.Enable = bEnable;
    return missing, nil;
}

func (a Proxy) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    if a.Enable {  // flag
        options = append(options, "-enable", "");
    }
    return Command{
        command: "proxy",
        args: args,
        options: options,
    };
}

func (a *Proxy) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Proxy) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Record struct {
    Args string;
}

func (a *Record) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "record" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a Record) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "record",
        args: args,
        options: options,
    };
}

func (a *Record) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Record) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Groups struct {
    Args []string;
}

func (a *Groups) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs []string;
    bArgs = b.args;
    if b.command != "groups" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a Groups) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = a.Args;

    return Command{
        command: "groups",
        args: args,
        options: options,
    };
}

func (a *Groups) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Groups) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type LogChannel struct {
    Args string;
}

func (a *LogChannel) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "log_channel" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a LogChannel) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "log_channel",
        args: args,
        options: options,
    };
}

func (a *LogChannel) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a LogChannel) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MainSite struct {
    Args string;
}

func (a *MainSite) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "main_site" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MainSite) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "main_site",
        args: args,
        options: options,
    };
}

func (a *MainSite) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MainSite) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type OauthAuthUrl struct {
    Args string;
}

func (a *OauthAuthUrl) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "oauth_auth_url" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a OauthAuthUrl) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "oauth_auth_url",
        args: args,
        options: options,
    };
}

func (a *OauthAuthUrl) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a OauthAuthUrl) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type OauthQueryUrl struct {
    Args string;
}

func (a *OauthQueryUrl) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "oauth_query_url" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a OauthQueryUrl) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "oauth_query_url",
        args: args,
        options: options,
    };
}

func (a *OauthQueryUrl) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a OauthQueryUrl) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type OauthQueryKey struct {
    Args string;
}

func (a *OauthQueryKey) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "oauth_query_key" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a OauthQueryKey) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "oauth_query_key",
        args: args,
        options: options,
    };
}

func (a *OauthQueryKey) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a OauthQueryKey) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type WikiList struct {
    Args []string;
}

func (a *WikiList) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs []string;
    bArgs = b.args;
    if b.command != "wiki_list" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a WikiList) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = a.Args;

    return Command{
        command: "wiki_list",
        args: args,
        options: options,
    };
}

func (a *WikiList) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a WikiList) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Blacklist struct {
    Args []string;
}

func (a *Blacklist) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs []string;
    bArgs = b.args;
    if b.command != "blacklist" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a Blacklist) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = a.Args;

    return Command{
        command: "blacklist",
        args: args,
        options: options,
    };
}

func (a *Blacklist) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Blacklist) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageStart struct {
    Args string;
}

func (a *MessageStart) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-start" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageStart) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-start",
        args: args,
        options: options,
    };
}

func (a *MessageStart) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageStart) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessagePolicy struct {
    Args string;
}

func (a *MessagePolicy) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-policy" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessagePolicy) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-policy",
        args: args,
        options: options,
    };
}

func (a *MessagePolicy) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessagePolicy) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageInsufficientRight struct {
    Args string;
}

func (a *MessageInsufficientRight) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-insufficient_right" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageInsufficientRight) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-insufficient_right",
        args: args,
        options: options,
    };
}

func (a *MessageInsufficientRight) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageInsufficientRight) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageGeneralPrompt struct {
    Args string;
}

func (a *MessageGeneralPrompt) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-general_prompt" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageGeneralPrompt) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-general_prompt",
        args: args,
        options: options,
    };
}

func (a *MessageGeneralPrompt) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageGeneralPrompt) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageTelegramIdError struct {
    Args string;
}

func (a *MessageTelegramIdError) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-telegram_id_error" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageTelegramIdError) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-telegram_id_error",
        args: args,
        options: options,
    };
}

func (a *MessageTelegramIdError) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageTelegramIdError) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageRestoreSilence struct {
    Args string;
}

func (a *MessageRestoreSilence) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-restore_silence" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageRestoreSilence) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-restore_silence",
        args: args,
        options: options,
    };
}

func (a *MessageRestoreSilence) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageRestoreSilence) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmAlready struct {
    Args string;
}

func (a *MessageConfirmAlready) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_already" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmAlready) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_already",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmAlready) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmAlready) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmOtherTg struct {
    Args string;
}

func (a *MessageConfirmOtherTg) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_other_tg" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmOtherTg) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_other_tg",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmOtherTg) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmOtherTg) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmConflict struct {
    Args string;
}

func (a *MessageConfirmConflict) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_conflict" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmConflict) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_conflict",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmConflict) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmConflict) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmChecking struct {
    Args string;
}

func (a *MessageConfirmChecking) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_checking" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmChecking) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_checking",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmChecking) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmChecking) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmUserNotFound struct {
    Args string;
}

func (a *MessageConfirmUserNotFound) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_user_not_found" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmUserNotFound) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_user_not_found",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmUserNotFound) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmUserNotFound) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmButton struct {
    Args string;
}

func (a *MessageConfirmButton) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_button" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmButton) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_button",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmButton) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmButton) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmWait struct {
    Args string;
}

func (a *MessageConfirmWait) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_wait" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmWait) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_wait",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmWait) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmWait) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmConfirming struct {
    Args string;
}

func (a *MessageConfirmConfirming) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_confirming" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmConfirming) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_confirming",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmConfirming) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmConfirming) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmIneligible struct {
    Args string;
}

func (a *MessageConfirmIneligible) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_ineligible" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmIneligible) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_ineligible",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmIneligible) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmIneligible) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmSessionLost struct {
    Args string;
}

func (a *MessageConfirmSessionLost) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_session_lost" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmSessionLost) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_session_lost",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmSessionLost) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmSessionLost) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmComplete struct {
    Args string;
}

func (a *MessageConfirmComplete) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_complete" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmComplete) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_complete",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmComplete) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmComplete) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmFailed struct {
    Args string;
}

func (a *MessageConfirmFailed) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_failed" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmFailed) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_failed",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmFailed) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmFailed) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageConfirmLog struct {
    Args string;
}

func (a *MessageConfirmLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-confirm_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageConfirmLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-confirm_log",
        args: args,
        options: options,
    };
}

func (a *MessageConfirmLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageConfirmLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDeconfirmPrompt struct {
    Args string;
}

func (a *MessageDeconfirmPrompt) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-deconfirm_prompt" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDeconfirmPrompt) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-deconfirm_prompt",
        args: args,
        options: options,
    };
}

func (a *MessageDeconfirmPrompt) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDeconfirmPrompt) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDeconfirmButton struct {
    Args string;
}

func (a *MessageDeconfirmButton) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-deconfirm_button" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDeconfirmButton) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-deconfirm_button",
        args: args,
        options: options,
    };
}

func (a *MessageDeconfirmButton) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDeconfirmButton) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDeconfirmSucc struct {
    Args string;
}

func (a *MessageDeconfirmSucc) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-deconfirm_succ" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDeconfirmSucc) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-deconfirm_succ",
        args: args,
        options: options,
    };
}

func (a *MessageDeconfirmSucc) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDeconfirmSucc) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDeconfirmNotConfirmed struct {
    Args string;
}

func (a *MessageDeconfirmNotConfirmed) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-deconfirm_not_confirmed" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDeconfirmNotConfirmed) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-deconfirm_not_confirmed",
        args: args,
        options: options,
    };
}

func (a *MessageDeconfirmNotConfirmed) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDeconfirmNotConfirmed) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDeconfirmLog struct {
    Args string;
}

func (a *MessageDeconfirmLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-deconfirm_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDeconfirmLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-deconfirm_log",
        args: args,
        options: options,
    };
}

func (a *MessageDeconfirmLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDeconfirmLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageNewMemberHint struct {
    Args string;
}

func (a *MessageNewMemberHint) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-new_member_hint" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageNewMemberHint) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-new_member_hint",
        args: args,
        options: options,
    };
}

func (a *MessageNewMemberHint) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageNewMemberHint) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageAddWhitelistPrompt struct {
    Args string;
}

func (a *MessageAddWhitelistPrompt) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-add_whitelist_prompt" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageAddWhitelistPrompt) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-add_whitelist_prompt",
        args: args,
        options: options,
    };
}

func (a *MessageAddWhitelistPrompt) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageAddWhitelistPrompt) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageAddWhitelistSucc struct {
    Args string;
}

func (a *MessageAddWhitelistSucc) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-add_whitelist_succ" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageAddWhitelistSucc) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-add_whitelist_succ",
        args: args,
        options: options,
    };
}

func (a *MessageAddWhitelistSucc) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageAddWhitelistSucc) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageAddWhitelistLog struct {
    Args string;
}

func (a *MessageAddWhitelistLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-add_whitelist_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageAddWhitelistLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-add_whitelist_log",
        args: args,
        options: options,
    };
}

func (a *MessageAddWhitelistLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageAddWhitelistLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageRemoveWhitelistPrompt struct {
    Args string;
}

func (a *MessageRemoveWhitelistPrompt) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-remove_whitelist_prompt" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageRemoveWhitelistPrompt) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-remove_whitelist_prompt",
        args: args,
        options: options,
    };
}

func (a *MessageRemoveWhitelistPrompt) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageRemoveWhitelistPrompt) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageRemoveWhitelistNotFound struct {
    Args string;
}

func (a *MessageRemoveWhitelistNotFound) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-remove_whitelist_not_found" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageRemoveWhitelistNotFound) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-remove_whitelist_not_found",
        args: args,
        options: options,
    };
}

func (a *MessageRemoveWhitelistNotFound) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageRemoveWhitelistNotFound) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageRemoveWhitelistLog struct {
    Args string;
}

func (a *MessageRemoveWhitelistLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-remove_whitelist_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageRemoveWhitelistLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-remove_whitelist_log",
        args: args,
        options: options,
    };
}

func (a *MessageRemoveWhitelistLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageRemoveWhitelistLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageRemoveWhitelistSucc struct {
    Args string;
}

func (a *MessageRemoveWhitelistSucc) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-remove_whitelist_succ" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageRemoveWhitelistSucc) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-remove_whitelist_succ",
        args: args,
        options: options,
    };
}

func (a *MessageRemoveWhitelistSucc) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageRemoveWhitelistSucc) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisHead struct {
    Args string;
}

func (a *MessageWhoisHead) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_head" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisHead) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_head",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisHead) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisHead) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisPrompt struct {
    Args string;
}

func (a *MessageWhoisPrompt) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_prompt" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisPrompt) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_prompt",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisPrompt) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisPrompt) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisNotFound struct {
    Args string;
}

func (a *MessageWhoisNotFound) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_not_found" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisNotFound) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_not_found",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisNotFound) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisNotFound) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisSelf struct {
    Args string;
}

func (a *MessageWhoisSelf) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_self" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisSelf) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_self",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisSelf) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisSelf) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisBot struct {
    Args string;
}

func (a *MessageWhoisBot) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_bot" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisBot) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_bot",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisBot) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisBot) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisHasMw struct {
    Args string;
}

func (a *MessageWhoisHasMw) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_has_mw" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisHasMw) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_has_mw",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisHasMw) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisHasMw) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisNoMw struct {
    Args string;
}

func (a *MessageWhoisNoMw) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_no_mw" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisNoMw) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_no_mw",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisNoMw) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisNoMw) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisWhitelisted struct {
    Args string;
}

func (a *MessageWhoisWhitelisted) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_whitelisted" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisWhitelisted) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_whitelisted",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisWhitelisted) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisWhitelisted) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageWhoisTgNameUnavailable struct {
    Args string;
}

func (a *MessageWhoisTgNameUnavailable) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-whois_tg_name_unavailable" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageWhoisTgNameUnavailable) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-whois_tg_name_unavailable",
        args: args,
        options: options,
    };
}

func (a *MessageWhoisTgNameUnavailable) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageWhoisTgNameUnavailable) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageRefuseLog struct {
    Args string;
}

func (a *MessageRefuseLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-refuse_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageRefuseLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-refuse_log",
        args: args,
        options: options,
    };
}

func (a *MessageRefuseLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageRefuseLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageAcceptLog struct {
    Args string;
}

func (a *MessageAcceptLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-accept_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageAcceptLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-accept_log",
        args: args,
        options: options,
    };
}

func (a *MessageAcceptLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageAcceptLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageLiftRestrictionAlert struct {
    Args string;
}

func (a *MessageLiftRestrictionAlert) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-lift_restriction_alert" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageLiftRestrictionAlert) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-lift_restriction_alert",
        args: args,
        options: options,
    };
}

func (a *MessageLiftRestrictionAlert) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageLiftRestrictionAlert) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageSilenceAlert struct {
    Args string;
}

func (a *MessageSilenceAlert) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-silence_alert" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageSilenceAlert) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-silence_alert",
        args: args,
        options: options,
    };
}

func (a *MessageSilenceAlert) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageSilenceAlert) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageEnable struct {
    Args string;
}

func (a *MessageEnable) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-enable" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageEnable) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-enable",
        args: args,
        options: options,
    };
}

func (a *MessageEnable) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageEnable) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDisable struct {
    Args string;
}

func (a *MessageDisable) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-disable" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDisable) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-disable",
        args: args,
        options: options,
    };
}

func (a *MessageDisable) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDisable) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageEnableLog struct {
    Args string;
}

func (a *MessageEnableLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-enable_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageEnableLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-enable_log",
        args: args,
        options: options,
    };
}

func (a *MessageEnableLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageEnableLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type MessageDisableLog struct {
    Args string;
}

func (a *MessageDisableLog) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "message-disable_log" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a MessageDisableLog) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "message-disable_log",
        args: args,
        options: options,
    };
}

func (a *MessageDisableLog) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a MessageDisableLog) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type ConfigMix struct {
    Token Token
    Proxy Proxy
    Record Record
    Groups Groups
    LogChannel LogChannel
    MainSite MainSite
    OauthAuthUrl OauthAuthUrl
    OauthQueryUrl OauthQueryUrl
    OauthQueryKey OauthQueryKey
    WikiList WikiList
    Blacklist Blacklist
    MessageStart MessageStart
    MessagePolicy MessagePolicy
    MessageInsufficientRight MessageInsufficientRight
    MessageGeneralPrompt MessageGeneralPrompt
    MessageTelegramIdError MessageTelegramIdError
    MessageRestoreSilence MessageRestoreSilence
    MessageConfirmAlready MessageConfirmAlready
    MessageConfirmOtherTg MessageConfirmOtherTg
    MessageConfirmConflict MessageConfirmConflict
    MessageConfirmChecking MessageConfirmChecking
    MessageConfirmUserNotFound MessageConfirmUserNotFound
    MessageConfirmButton MessageConfirmButton
    MessageConfirmWait MessageConfirmWait
    MessageConfirmConfirming MessageConfirmConfirming
    MessageConfirmIneligible MessageConfirmIneligible
    MessageConfirmSessionLost MessageConfirmSessionLost
    MessageConfirmComplete MessageConfirmComplete
    MessageConfirmFailed MessageConfirmFailed
    MessageConfirmLog MessageConfirmLog
    MessageDeconfirmPrompt MessageDeconfirmPrompt
    MessageDeconfirmButton MessageDeconfirmButton
    MessageDeconfirmSucc MessageDeconfirmSucc
    MessageDeconfirmNotConfirmed MessageDeconfirmNotConfirmed
    MessageDeconfirmLog MessageDeconfirmLog
    MessageNewMemberHint MessageNewMemberHint
    MessageAddWhitelistPrompt MessageAddWhitelistPrompt
    MessageAddWhitelistSucc MessageAddWhitelistSucc
    MessageAddWhitelistLog MessageAddWhitelistLog
    MessageRemoveWhitelistPrompt MessageRemoveWhitelistPrompt
    MessageRemoveWhitelistNotFound MessageRemoveWhitelistNotFound
    MessageRemoveWhitelistLog MessageRemoveWhitelistLog
    MessageRemoveWhitelistSucc MessageRemoveWhitelistSucc
    MessageWhoisHead MessageWhoisHead
    MessageWhoisPrompt MessageWhoisPrompt
    MessageWhoisNotFound MessageWhoisNotFound
    MessageWhoisSelf MessageWhoisSelf
    MessageWhoisBot MessageWhoisBot
    MessageWhoisHasMw MessageWhoisHasMw
    MessageWhoisNoMw MessageWhoisNoMw
    MessageWhoisWhitelisted MessageWhoisWhitelisted
    MessageWhoisTgNameUnavailable MessageWhoisTgNameUnavailable
    MessageRefuseLog MessageRefuseLog
    MessageAcceptLog MessageAcceptLog
    MessageLiftRestrictionAlert MessageLiftRestrictionAlert
    MessageSilenceAlert MessageSilenceAlert
    MessageEnable MessageEnable
    MessageDisable MessageDisable
    MessageEnableLog MessageEnableLog
    MessageDisableLog MessageDisableLog
}

func (a *ConfigMix) Parse(wiredata string) (string, error) {
    res, _, err := CcCore.CoreDynParse(wiredata);
    if err != nil { return "", err; }
    missing := "";
    mixMissing := ""
    var bToken Token;
    var mToken = "";
    var cToken = 0;
    var bProxy Proxy;
    var mProxy = "";
    var cProxy = 0;
    var bRecord Record;
    var mRecord = "";
    var cRecord = 0;
    var bGroups Groups;
    var mGroups = "";
    var cGroups = 0;
    var bLogChannel LogChannel;
    var mLogChannel = "";
    var cLogChannel = 0;
    var bMainSite MainSite;
    var mMainSite = "";
    var cMainSite = 0;
    var bOauthAuthUrl OauthAuthUrl;
    var mOauthAuthUrl = "";
    var cOauthAuthUrl = 0;
    var bOauthQueryUrl OauthQueryUrl;
    var mOauthQueryUrl = "";
    var cOauthQueryUrl = 0;
    var bOauthQueryKey OauthQueryKey;
    var mOauthQueryKey = "";
    var cOauthQueryKey = 0;
    var bWikiList WikiList;
    var mWikiList = "";
    var cWikiList = 0;
    var bBlacklist Blacklist;
    var mBlacklist = "";
    var cBlacklist = 0;
    var bMessageStart MessageStart;
    var mMessageStart = "";
    var cMessageStart = 0;
    var bMessagePolicy MessagePolicy;
    var mMessagePolicy = "";
    var cMessagePolicy = 0;
    var bMessageInsufficientRight MessageInsufficientRight;
    var mMessageInsufficientRight = "";
    var cMessageInsufficientRight = 0;
    var bMessageGeneralPrompt MessageGeneralPrompt;
    var mMessageGeneralPrompt = "";
    var cMessageGeneralPrompt = 0;
    var bMessageTelegramIdError MessageTelegramIdError;
    var mMessageTelegramIdError = "";
    var cMessageTelegramIdError = 0;
    var bMessageRestoreSilence MessageRestoreSilence;
    var mMessageRestoreSilence = "";
    var cMessageRestoreSilence = 0;
    var bMessageConfirmAlready MessageConfirmAlready;
    var mMessageConfirmAlready = "";
    var cMessageConfirmAlready = 0;
    var bMessageConfirmOtherTg MessageConfirmOtherTg;
    var mMessageConfirmOtherTg = "";
    var cMessageConfirmOtherTg = 0;
    var bMessageConfirmConflict MessageConfirmConflict;
    var mMessageConfirmConflict = "";
    var cMessageConfirmConflict = 0;
    var bMessageConfirmChecking MessageConfirmChecking;
    var mMessageConfirmChecking = "";
    var cMessageConfirmChecking = 0;
    var bMessageConfirmUserNotFound MessageConfirmUserNotFound;
    var mMessageConfirmUserNotFound = "";
    var cMessageConfirmUserNotFound = 0;
    var bMessageConfirmButton MessageConfirmButton;
    var mMessageConfirmButton = "";
    var cMessageConfirmButton = 0;
    var bMessageConfirmWait MessageConfirmWait;
    var mMessageConfirmWait = "";
    var cMessageConfirmWait = 0;
    var bMessageConfirmConfirming MessageConfirmConfirming;
    var mMessageConfirmConfirming = "";
    var cMessageConfirmConfirming = 0;
    var bMessageConfirmIneligible MessageConfirmIneligible;
    var mMessageConfirmIneligible = "";
    var cMessageConfirmIneligible = 0;
    var bMessageConfirmSessionLost MessageConfirmSessionLost;
    var mMessageConfirmSessionLost = "";
    var cMessageConfirmSessionLost = 0;
    var bMessageConfirmComplete MessageConfirmComplete;
    var mMessageConfirmComplete = "";
    var cMessageConfirmComplete = 0;
    var bMessageConfirmFailed MessageConfirmFailed;
    var mMessageConfirmFailed = "";
    var cMessageConfirmFailed = 0;
    var bMessageConfirmLog MessageConfirmLog;
    var mMessageConfirmLog = "";
    var cMessageConfirmLog = 0;
    var bMessageDeconfirmPrompt MessageDeconfirmPrompt;
    var mMessageDeconfirmPrompt = "";
    var cMessageDeconfirmPrompt = 0;
    var bMessageDeconfirmButton MessageDeconfirmButton;
    var mMessageDeconfirmButton = "";
    var cMessageDeconfirmButton = 0;
    var bMessageDeconfirmSucc MessageDeconfirmSucc;
    var mMessageDeconfirmSucc = "";
    var cMessageDeconfirmSucc = 0;
    var bMessageDeconfirmNotConfirmed MessageDeconfirmNotConfirmed;
    var mMessageDeconfirmNotConfirmed = "";
    var cMessageDeconfirmNotConfirmed = 0;
    var bMessageDeconfirmLog MessageDeconfirmLog;
    var mMessageDeconfirmLog = "";
    var cMessageDeconfirmLog = 0;
    var bMessageNewMemberHint MessageNewMemberHint;
    var mMessageNewMemberHint = "";
    var cMessageNewMemberHint = 0;
    var bMessageAddWhitelistPrompt MessageAddWhitelistPrompt;
    var mMessageAddWhitelistPrompt = "";
    var cMessageAddWhitelistPrompt = 0;
    var bMessageAddWhitelistSucc MessageAddWhitelistSucc;
    var mMessageAddWhitelistSucc = "";
    var cMessageAddWhitelistSucc = 0;
    var bMessageAddWhitelistLog MessageAddWhitelistLog;
    var mMessageAddWhitelistLog = "";
    var cMessageAddWhitelistLog = 0;
    var bMessageRemoveWhitelistPrompt MessageRemoveWhitelistPrompt;
    var mMessageRemoveWhitelistPrompt = "";
    var cMessageRemoveWhitelistPrompt = 0;
    var bMessageRemoveWhitelistNotFound MessageRemoveWhitelistNotFound;
    var mMessageRemoveWhitelistNotFound = "";
    var cMessageRemoveWhitelistNotFound = 0;
    var bMessageRemoveWhitelistLog MessageRemoveWhitelistLog;
    var mMessageRemoveWhitelistLog = "";
    var cMessageRemoveWhitelistLog = 0;
    var bMessageRemoveWhitelistSucc MessageRemoveWhitelistSucc;
    var mMessageRemoveWhitelistSucc = "";
    var cMessageRemoveWhitelistSucc = 0;
    var bMessageWhoisHead MessageWhoisHead;
    var mMessageWhoisHead = "";
    var cMessageWhoisHead = 0;
    var bMessageWhoisPrompt MessageWhoisPrompt;
    var mMessageWhoisPrompt = "";
    var cMessageWhoisPrompt = 0;
    var bMessageWhoisNotFound MessageWhoisNotFound;
    var mMessageWhoisNotFound = "";
    var cMessageWhoisNotFound = 0;
    var bMessageWhoisSelf MessageWhoisSelf;
    var mMessageWhoisSelf = "";
    var cMessageWhoisSelf = 0;
    var bMessageWhoisBot MessageWhoisBot;
    var mMessageWhoisBot = "";
    var cMessageWhoisBot = 0;
    var bMessageWhoisHasMw MessageWhoisHasMw;
    var mMessageWhoisHasMw = "";
    var cMessageWhoisHasMw = 0;
    var bMessageWhoisNoMw MessageWhoisNoMw;
    var mMessageWhoisNoMw = "";
    var cMessageWhoisNoMw = 0;
    var bMessageWhoisWhitelisted MessageWhoisWhitelisted;
    var mMessageWhoisWhitelisted = "";
    var cMessageWhoisWhitelisted = 0;
    var bMessageWhoisTgNameUnavailable MessageWhoisTgNameUnavailable;
    var mMessageWhoisTgNameUnavailable = "";
    var cMessageWhoisTgNameUnavailable = 0;
    var bMessageRefuseLog MessageRefuseLog;
    var mMessageRefuseLog = "";
    var cMessageRefuseLog = 0;
    var bMessageAcceptLog MessageAcceptLog;
    var mMessageAcceptLog = "";
    var cMessageAcceptLog = 0;
    var bMessageLiftRestrictionAlert MessageLiftRestrictionAlert;
    var mMessageLiftRestrictionAlert = "";
    var cMessageLiftRestrictionAlert = 0;
    var bMessageSilenceAlert MessageSilenceAlert;
    var mMessageSilenceAlert = "";
    var cMessageSilenceAlert = 0;
    var bMessageEnable MessageEnable;
    var mMessageEnable = "";
    var cMessageEnable = 0;
    var bMessageDisable MessageDisable;
    var mMessageDisable = "";
    var cMessageDisable = 0;
    var bMessageEnableLog MessageEnableLog;
    var mMessageEnableLog = "";
    var cMessageEnableLog = 0;
    var bMessageDisableLog MessageDisableLog;
    var mMessageDisableLog = "";
    var cMessageDisableLog = 0;
    for _, cmd := range res {
        switch cmd.command {
        case "token":
            var parsed Token;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bToken = parsed;
            mToken += mis;
            cToken += 1;
        case "proxy":
            var parsed Proxy;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bProxy = parsed;
            mProxy += mis;
            cProxy += 1;
        case "record":
            var parsed Record;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bRecord = parsed;
            mRecord += mis;
            cRecord += 1;
        case "groups":
            var parsed Groups;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bGroups = parsed;
            mGroups += mis;
            cGroups += 1;
        case "log_channel":
            var parsed LogChannel;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bLogChannel = parsed;
            mLogChannel += mis;
            cLogChannel += 1;
        case "main_site":
            var parsed MainSite;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMainSite = parsed;
            mMainSite += mis;
            cMainSite += 1;
        case "oauth_auth_url":
            var parsed OauthAuthUrl;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bOauthAuthUrl = parsed;
            mOauthAuthUrl += mis;
            cOauthAuthUrl += 1;
        case "oauth_query_url":
            var parsed OauthQueryUrl;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bOauthQueryUrl = parsed;
            mOauthQueryUrl += mis;
            cOauthQueryUrl += 1;
        case "oauth_query_key":
            var parsed OauthQueryKey;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bOauthQueryKey = parsed;
            mOauthQueryKey += mis;
            cOauthQueryKey += 1;
        case "wiki_list":
            var parsed WikiList;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bWikiList = parsed;
            mWikiList += mis;
            cWikiList += 1;
        case "blacklist":
            var parsed Blacklist;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bBlacklist = parsed;
            mBlacklist += mis;
            cBlacklist += 1;
        case "message-start":
            var parsed MessageStart;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageStart = parsed;
            mMessageStart += mis;
            cMessageStart += 1;
        case "message-policy":
            var parsed MessagePolicy;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessagePolicy = parsed;
            mMessagePolicy += mis;
            cMessagePolicy += 1;
        case "message-insufficient_right":
            var parsed MessageInsufficientRight;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageInsufficientRight = parsed;
            mMessageInsufficientRight += mis;
            cMessageInsufficientRight += 1;
        case "message-general_prompt":
            var parsed MessageGeneralPrompt;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageGeneralPrompt = parsed;
            mMessageGeneralPrompt += mis;
            cMessageGeneralPrompt += 1;
        case "message-telegram_id_error":
            var parsed MessageTelegramIdError;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageTelegramIdError = parsed;
            mMessageTelegramIdError += mis;
            cMessageTelegramIdError += 1;
        case "message-restore_silence":
            var parsed MessageRestoreSilence;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageRestoreSilence = parsed;
            mMessageRestoreSilence += mis;
            cMessageRestoreSilence += 1;
        case "message-confirm_already":
            var parsed MessageConfirmAlready;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmAlready = parsed;
            mMessageConfirmAlready += mis;
            cMessageConfirmAlready += 1;
        case "message-confirm_other_tg":
            var parsed MessageConfirmOtherTg;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmOtherTg = parsed;
            mMessageConfirmOtherTg += mis;
            cMessageConfirmOtherTg += 1;
        case "message-confirm_conflict":
            var parsed MessageConfirmConflict;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmConflict = parsed;
            mMessageConfirmConflict += mis;
            cMessageConfirmConflict += 1;
        case "message-confirm_checking":
            var parsed MessageConfirmChecking;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmChecking = parsed;
            mMessageConfirmChecking += mis;
            cMessageConfirmChecking += 1;
        case "message-confirm_user_not_found":
            var parsed MessageConfirmUserNotFound;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmUserNotFound = parsed;
            mMessageConfirmUserNotFound += mis;
            cMessageConfirmUserNotFound += 1;
        case "message-confirm_button":
            var parsed MessageConfirmButton;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmButton = parsed;
            mMessageConfirmButton += mis;
            cMessageConfirmButton += 1;
        case "message-confirm_wait":
            var parsed MessageConfirmWait;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmWait = parsed;
            mMessageConfirmWait += mis;
            cMessageConfirmWait += 1;
        case "message-confirm_confirming":
            var parsed MessageConfirmConfirming;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmConfirming = parsed;
            mMessageConfirmConfirming += mis;
            cMessageConfirmConfirming += 1;
        case "message-confirm_ineligible":
            var parsed MessageConfirmIneligible;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmIneligible = parsed;
            mMessageConfirmIneligible += mis;
            cMessageConfirmIneligible += 1;
        case "message-confirm_session_lost":
            var parsed MessageConfirmSessionLost;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmSessionLost = parsed;
            mMessageConfirmSessionLost += mis;
            cMessageConfirmSessionLost += 1;
        case "message-confirm_complete":
            var parsed MessageConfirmComplete;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmComplete = parsed;
            mMessageConfirmComplete += mis;
            cMessageConfirmComplete += 1;
        case "message-confirm_failed":
            var parsed MessageConfirmFailed;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmFailed = parsed;
            mMessageConfirmFailed += mis;
            cMessageConfirmFailed += 1;
        case "message-confirm_log":
            var parsed MessageConfirmLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageConfirmLog = parsed;
            mMessageConfirmLog += mis;
            cMessageConfirmLog += 1;
        case "message-deconfirm_prompt":
            var parsed MessageDeconfirmPrompt;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDeconfirmPrompt = parsed;
            mMessageDeconfirmPrompt += mis;
            cMessageDeconfirmPrompt += 1;
        case "message-deconfirm_button":
            var parsed MessageDeconfirmButton;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDeconfirmButton = parsed;
            mMessageDeconfirmButton += mis;
            cMessageDeconfirmButton += 1;
        case "message-deconfirm_succ":
            var parsed MessageDeconfirmSucc;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDeconfirmSucc = parsed;
            mMessageDeconfirmSucc += mis;
            cMessageDeconfirmSucc += 1;
        case "message-deconfirm_not_confirmed":
            var parsed MessageDeconfirmNotConfirmed;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDeconfirmNotConfirmed = parsed;
            mMessageDeconfirmNotConfirmed += mis;
            cMessageDeconfirmNotConfirmed += 1;
        case "message-deconfirm_log":
            var parsed MessageDeconfirmLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDeconfirmLog = parsed;
            mMessageDeconfirmLog += mis;
            cMessageDeconfirmLog += 1;
        case "message-new_member_hint":
            var parsed MessageNewMemberHint;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageNewMemberHint = parsed;
            mMessageNewMemberHint += mis;
            cMessageNewMemberHint += 1;
        case "message-add_whitelist_prompt":
            var parsed MessageAddWhitelistPrompt;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageAddWhitelistPrompt = parsed;
            mMessageAddWhitelistPrompt += mis;
            cMessageAddWhitelistPrompt += 1;
        case "message-add_whitelist_succ":
            var parsed MessageAddWhitelistSucc;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageAddWhitelistSucc = parsed;
            mMessageAddWhitelistSucc += mis;
            cMessageAddWhitelistSucc += 1;
        case "message-add_whitelist_log":
            var parsed MessageAddWhitelistLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageAddWhitelistLog = parsed;
            mMessageAddWhitelistLog += mis;
            cMessageAddWhitelistLog += 1;
        case "message-remove_whitelist_prompt":
            var parsed MessageRemoveWhitelistPrompt;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageRemoveWhitelistPrompt = parsed;
            mMessageRemoveWhitelistPrompt += mis;
            cMessageRemoveWhitelistPrompt += 1;
        case "message-remove_whitelist_not_found":
            var parsed MessageRemoveWhitelistNotFound;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageRemoveWhitelistNotFound = parsed;
            mMessageRemoveWhitelistNotFound += mis;
            cMessageRemoveWhitelistNotFound += 1;
        case "message-remove_whitelist_log":
            var parsed MessageRemoveWhitelistLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageRemoveWhitelistLog = parsed;
            mMessageRemoveWhitelistLog += mis;
            cMessageRemoveWhitelistLog += 1;
        case "message-remove_whitelist_succ":
            var parsed MessageRemoveWhitelistSucc;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageRemoveWhitelistSucc = parsed;
            mMessageRemoveWhitelistSucc += mis;
            cMessageRemoveWhitelistSucc += 1;
        case "message-whois_head":
            var parsed MessageWhoisHead;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisHead = parsed;
            mMessageWhoisHead += mis;
            cMessageWhoisHead += 1;
        case "message-whois_prompt":
            var parsed MessageWhoisPrompt;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisPrompt = parsed;
            mMessageWhoisPrompt += mis;
            cMessageWhoisPrompt += 1;
        case "message-whois_not_found":
            var parsed MessageWhoisNotFound;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisNotFound = parsed;
            mMessageWhoisNotFound += mis;
            cMessageWhoisNotFound += 1;
        case "message-whois_self":
            var parsed MessageWhoisSelf;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisSelf = parsed;
            mMessageWhoisSelf += mis;
            cMessageWhoisSelf += 1;
        case "message-whois_bot":
            var parsed MessageWhoisBot;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisBot = parsed;
            mMessageWhoisBot += mis;
            cMessageWhoisBot += 1;
        case "message-whois_has_mw":
            var parsed MessageWhoisHasMw;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisHasMw = parsed;
            mMessageWhoisHasMw += mis;
            cMessageWhoisHasMw += 1;
        case "message-whois_no_mw":
            var parsed MessageWhoisNoMw;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisNoMw = parsed;
            mMessageWhoisNoMw += mis;
            cMessageWhoisNoMw += 1;
        case "message-whois_whitelisted":
            var parsed MessageWhoisWhitelisted;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisWhitelisted = parsed;
            mMessageWhoisWhitelisted += mis;
            cMessageWhoisWhitelisted += 1;
        case "message-whois_tg_name_unavailable":
            var parsed MessageWhoisTgNameUnavailable;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageWhoisTgNameUnavailable = parsed;
            mMessageWhoisTgNameUnavailable += mis;
            cMessageWhoisTgNameUnavailable += 1;
        case "message-refuse_log":
            var parsed MessageRefuseLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageRefuseLog = parsed;
            mMessageRefuseLog += mis;
            cMessageRefuseLog += 1;
        case "message-accept_log":
            var parsed MessageAcceptLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageAcceptLog = parsed;
            mMessageAcceptLog += mis;
            cMessageAcceptLog += 1;
        case "message-lift_restriction_alert":
            var parsed MessageLiftRestrictionAlert;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageLiftRestrictionAlert = parsed;
            mMessageLiftRestrictionAlert += mis;
            cMessageLiftRestrictionAlert += 1;
        case "message-silence_alert":
            var parsed MessageSilenceAlert;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageSilenceAlert = parsed;
            mMessageSilenceAlert += mis;
            cMessageSilenceAlert += 1;
        case "message-enable":
            var parsed MessageEnable;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageEnable = parsed;
            mMessageEnable += mis;
            cMessageEnable += 1;
        case "message-disable":
            var parsed MessageDisable;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDisable = parsed;
            mMessageDisable += mis;
            cMessageDisable += 1;
        case "message-enable_log":
            var parsed MessageEnableLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageEnableLog = parsed;
            mMessageEnableLog += mis;
            cMessageEnableLog += 1;
        case "message-disable_log":
            var parsed MessageDisableLog;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bMessageDisableLog = parsed;
            mMessageDisableLog += mis;
            cMessageDisableLog += 1;
        }
    }
    if cToken < 1 {
        mixMissing += "token;";
    }
    if cProxy < 1 {
        mixMissing += "proxy;";
    }
    if cRecord < 1 {
        mixMissing += "record;";
    }
    if cGroups < 1 {
        mixMissing += "groups;";
    }
    if cLogChannel < 1 {
        mixMissing += "log_channel;";
    }
    if cMainSite < 1 {
        mixMissing += "main_site;";
    }
    if cOauthAuthUrl < 1 {
        mixMissing += "oauth_auth_url;";
    }
    if cOauthQueryUrl < 1 {
        mixMissing += "oauth_query_url;";
    }
    if cOauthQueryKey < 1 {
        mixMissing += "oauth_query_key;";
    }
    if cWikiList < 1 {
        mixMissing += "wiki_list;";
    }
    if cBlacklist < 1 {
        mixMissing += "blacklist;";
    }
    if cMessageStart < 1 {
        mixMissing += "message-start;";
    }
    if cMessagePolicy < 1 {
        mixMissing += "message-policy;";
    }
    if cMessageInsufficientRight < 1 {
        mixMissing += "message-insufficient_right;";
    }
    if cMessageGeneralPrompt < 1 {
        mixMissing += "message-general_prompt;";
    }
    if cMessageTelegramIdError < 1 {
        mixMissing += "message-telegram_id_error;";
    }
    if cMessageRestoreSilence < 1 {
        mixMissing += "message-restore_silence;";
    }
    if cMessageConfirmAlready < 1 {
        mixMissing += "message-confirm_already;";
    }
    if cMessageConfirmOtherTg < 1 {
        mixMissing += "message-confirm_other_tg;";
    }
    if cMessageConfirmConflict < 1 {
        mixMissing += "message-confirm_conflict;";
    }
    if cMessageConfirmChecking < 1 {
        mixMissing += "message-confirm_checking;";
    }
    if cMessageConfirmUserNotFound < 1 {
        mixMissing += "message-confirm_user_not_found;";
    }
    if cMessageConfirmButton < 1 {
        mixMissing += "message-confirm_button;";
    }
    if cMessageConfirmWait < 1 {
        mixMissing += "message-confirm_wait;";
    }
    if cMessageConfirmConfirming < 1 {
        mixMissing += "message-confirm_confirming;";
    }
    if cMessageConfirmIneligible < 1 {
        mixMissing += "message-confirm_ineligible;";
    }
    if cMessageConfirmSessionLost < 1 {
        mixMissing += "message-confirm_session_lost;";
    }
    if cMessageConfirmComplete < 1 {
        mixMissing += "message-confirm_complete;";
    }
    if cMessageConfirmFailed < 1 {
        mixMissing += "message-confirm_failed;";
    }
    if cMessageConfirmLog < 1 {
        mixMissing += "message-confirm_log;";
    }
    if cMessageDeconfirmPrompt < 1 {
        mixMissing += "message-deconfirm_prompt;";
    }
    if cMessageDeconfirmButton < 1 {
        mixMissing += "message-deconfirm_button;";
    }
    if cMessageDeconfirmSucc < 1 {
        mixMissing += "message-deconfirm_succ;";
    }
    if cMessageDeconfirmNotConfirmed < 1 {
        mixMissing += "message-deconfirm_not_confirmed;";
    }
    if cMessageDeconfirmLog < 1 {
        mixMissing += "message-deconfirm_log;";
    }
    if cMessageNewMemberHint < 1 {
        mixMissing += "message-new_member_hint;";
    }
    if cMessageAddWhitelistPrompt < 1 {
        mixMissing += "message-add_whitelist_prompt;";
    }
    if cMessageAddWhitelistSucc < 1 {
        mixMissing += "message-add_whitelist_succ;";
    }
    if cMessageAddWhitelistLog < 1 {
        mixMissing += "message-add_whitelist_log;";
    }
    if cMessageRemoveWhitelistPrompt < 1 {
        mixMissing += "message-remove_whitelist_prompt;";
    }
    if cMessageRemoveWhitelistNotFound < 1 {
        mixMissing += "message-remove_whitelist_not_found;";
    }
    if cMessageRemoveWhitelistLog < 1 {
        mixMissing += "message-remove_whitelist_log;";
    }
    if cMessageRemoveWhitelistSucc < 1 {
        mixMissing += "message-remove_whitelist_succ;";
    }
    if cMessageWhoisHead < 1 {
        mixMissing += "message-whois_head;";
    }
    if cMessageWhoisPrompt < 1 {
        mixMissing += "message-whois_prompt;";
    }
    if cMessageWhoisNotFound < 1 {
        mixMissing += "message-whois_not_found;";
    }
    if cMessageWhoisSelf < 1 {
        mixMissing += "message-whois_self;";
    }
    if cMessageWhoisBot < 1 {
        mixMissing += "message-whois_bot;";
    }
    if cMessageWhoisHasMw < 1 {
        mixMissing += "message-whois_has_mw;";
    }
    if cMessageWhoisNoMw < 1 {
        mixMissing += "message-whois_no_mw;";
    }
    if cMessageWhoisWhitelisted < 1 {
        mixMissing += "message-whois_whitelisted;";
    }
    if cMessageWhoisTgNameUnavailable < 1 {
        mixMissing += "message-whois_tg_name_unavailable;";
    }
    if cMessageRefuseLog < 1 {
        mixMissing += "message-refuse_log;";
    }
    if cMessageAcceptLog < 1 {
        mixMissing += "message-accept_log;";
    }
    if cMessageLiftRestrictionAlert < 1 {
        mixMissing += "message-lift_restriction_alert;";
    }
    if cMessageSilenceAlert < 1 {
        mixMissing += "message-silence_alert;";
    }
    if cMessageEnable < 1 {
        mixMissing += "message-enable;";
    }
    if cMessageDisable < 1 {
        mixMissing += "message-disable;";
    }
    if cMessageEnableLog < 1 {
        mixMissing += "message-enable_log;";
    }
    if cMessageDisableLog < 1 {
        mixMissing += "message-disable_log;";
    }
    if len(mixMissing) > 0 { missing += "mix:" + mixMissing + "\n"; }
    if len(mToken) > 0 { missing += "token:" + mToken + "\n"; }
    if len(mProxy) > 0 { missing += "proxy:" + mProxy + "\n"; }
    if len(mRecord) > 0 { missing += "record:" + mRecord + "\n"; }
    if len(mGroups) > 0 { missing += "groups:" + mGroups + "\n"; }
    if len(mLogChannel) > 0 { missing += "log_channel:" + mLogChannel + "\n"; }
    if len(mMainSite) > 0 { missing += "main_site:" + mMainSite + "\n"; }
    if len(mOauthAuthUrl) > 0 { missing += "oauth_auth_url:" + mOauthAuthUrl + "\n"; }
    if len(mOauthQueryUrl) > 0 { missing += "oauth_query_url:" + mOauthQueryUrl + "\n"; }
    if len(mOauthQueryKey) > 0 { missing += "oauth_query_key:" + mOauthQueryKey + "\n"; }
    if len(mWikiList) > 0 { missing += "wiki_list:" + mWikiList + "\n"; }
    if len(mBlacklist) > 0 { missing += "blacklist:" + mBlacklist + "\n"; }
    if len(mMessageStart) > 0 { missing += "message-start:" + mMessageStart + "\n"; }
    if len(mMessagePolicy) > 0 { missing += "message-policy:" + mMessagePolicy + "\n"; }
    if len(mMessageInsufficientRight) > 0 { missing += "message-insufficient_right:" + mMessageInsufficientRight + "\n"; }
    if len(mMessageGeneralPrompt) > 0 { missing += "message-general_prompt:" + mMessageGeneralPrompt + "\n"; }
    if len(mMessageTelegramIdError) > 0 { missing += "message-telegram_id_error:" + mMessageTelegramIdError + "\n"; }
    if len(mMessageRestoreSilence) > 0 { missing += "message-restore_silence:" + mMessageRestoreSilence + "\n"; }
    if len(mMessageConfirmAlready) > 0 { missing += "message-confirm_already:" + mMessageConfirmAlready + "\n"; }
    if len(mMessageConfirmOtherTg) > 0 { missing += "message-confirm_other_tg:" + mMessageConfirmOtherTg + "\n"; }
    if len(mMessageConfirmConflict) > 0 { missing += "message-confirm_conflict:" + mMessageConfirmConflict + "\n"; }
    if len(mMessageConfirmChecking) > 0 { missing += "message-confirm_checking:" + mMessageConfirmChecking + "\n"; }
    if len(mMessageConfirmUserNotFound) > 0 { missing += "message-confirm_user_not_found:" + mMessageConfirmUserNotFound + "\n"; }
    if len(mMessageConfirmButton) > 0 { missing += "message-confirm_button:" + mMessageConfirmButton + "\n"; }
    if len(mMessageConfirmWait) > 0 { missing += "message-confirm_wait:" + mMessageConfirmWait + "\n"; }
    if len(mMessageConfirmConfirming) > 0 { missing += "message-confirm_confirming:" + mMessageConfirmConfirming + "\n"; }
    if len(mMessageConfirmIneligible) > 0 { missing += "message-confirm_ineligible:" + mMessageConfirmIneligible + "\n"; }
    if len(mMessageConfirmSessionLost) > 0 { missing += "message-confirm_session_lost:" + mMessageConfirmSessionLost + "\n"; }
    if len(mMessageConfirmComplete) > 0 { missing += "message-confirm_complete:" + mMessageConfirmComplete + "\n"; }
    if len(mMessageConfirmFailed) > 0 { missing += "message-confirm_failed:" + mMessageConfirmFailed + "\n"; }
    if len(mMessageConfirmLog) > 0 { missing += "message-confirm_log:" + mMessageConfirmLog + "\n"; }
    if len(mMessageDeconfirmPrompt) > 0 { missing += "message-deconfirm_prompt:" + mMessageDeconfirmPrompt + "\n"; }
    if len(mMessageDeconfirmButton) > 0 { missing += "message-deconfirm_button:" + mMessageDeconfirmButton + "\n"; }
    if len(mMessageDeconfirmSucc) > 0 { missing += "message-deconfirm_succ:" + mMessageDeconfirmSucc + "\n"; }
    if len(mMessageDeconfirmNotConfirmed) > 0 { missing += "message-deconfirm_not_confirmed:" + mMessageDeconfirmNotConfirmed + "\n"; }
    if len(mMessageDeconfirmLog) > 0 { missing += "message-deconfirm_log:" + mMessageDeconfirmLog + "\n"; }
    if len(mMessageNewMemberHint) > 0 { missing += "message-new_member_hint:" + mMessageNewMemberHint + "\n"; }
    if len(mMessageAddWhitelistPrompt) > 0 { missing += "message-add_whitelist_prompt:" + mMessageAddWhitelistPrompt + "\n"; }
    if len(mMessageAddWhitelistSucc) > 0 { missing += "message-add_whitelist_succ:" + mMessageAddWhitelistSucc + "\n"; }
    if len(mMessageAddWhitelistLog) > 0 { missing += "message-add_whitelist_log:" + mMessageAddWhitelistLog + "\n"; }
    if len(mMessageRemoveWhitelistPrompt) > 0 { missing += "message-remove_whitelist_prompt:" + mMessageRemoveWhitelistPrompt + "\n"; }
    if len(mMessageRemoveWhitelistNotFound) > 0 { missing += "message-remove_whitelist_not_found:" + mMessageRemoveWhitelistNotFound + "\n"; }
    if len(mMessageRemoveWhitelistLog) > 0 { missing += "message-remove_whitelist_log:" + mMessageRemoveWhitelistLog + "\n"; }
    if len(mMessageRemoveWhitelistSucc) > 0 { missing += "message-remove_whitelist_succ:" + mMessageRemoveWhitelistSucc + "\n"; }
    if len(mMessageWhoisHead) > 0 { missing += "message-whois_head:" + mMessageWhoisHead + "\n"; }
    if len(mMessageWhoisPrompt) > 0 { missing += "message-whois_prompt:" + mMessageWhoisPrompt + "\n"; }
    if len(mMessageWhoisNotFound) > 0 { missing += "message-whois_not_found:" + mMessageWhoisNotFound + "\n"; }
    if len(mMessageWhoisSelf) > 0 { missing += "message-whois_self:" + mMessageWhoisSelf + "\n"; }
    if len(mMessageWhoisBot) > 0 { missing += "message-whois_bot:" + mMessageWhoisBot + "\n"; }
    if len(mMessageWhoisHasMw) > 0 { missing += "message-whois_has_mw:" + mMessageWhoisHasMw + "\n"; }
    if len(mMessageWhoisNoMw) > 0 { missing += "message-whois_no_mw:" + mMessageWhoisNoMw + "\n"; }
    if len(mMessageWhoisWhitelisted) > 0 { missing += "message-whois_whitelisted:" + mMessageWhoisWhitelisted + "\n"; }
    if len(mMessageWhoisTgNameUnavailable) > 0 { missing += "message-whois_tg_name_unavailable:" + mMessageWhoisTgNameUnavailable + "\n"; }
    if len(mMessageRefuseLog) > 0 { missing += "message-refuse_log:" + mMessageRefuseLog + "\n"; }
    if len(mMessageAcceptLog) > 0 { missing += "message-accept_log:" + mMessageAcceptLog + "\n"; }
    if len(mMessageLiftRestrictionAlert) > 0 { missing += "message-lift_restriction_alert:" + mMessageLiftRestrictionAlert + "\n"; }
    if len(mMessageSilenceAlert) > 0 { missing += "message-silence_alert:" + mMessageSilenceAlert + "\n"; }
    if len(mMessageEnable) > 0 { missing += "message-enable:" + mMessageEnable + "\n"; }
    if len(mMessageDisable) > 0 { missing += "message-disable:" + mMessageDisable + "\n"; }
    if len(mMessageEnableLog) > 0 { missing += "message-enable_log:" + mMessageEnableLog + "\n"; }
    if len(mMessageDisableLog) > 0 { missing += "message-disable_log:" + mMessageDisableLog + "\n"; }
    a.Token = bToken;
    a.Proxy = bProxy;
    a.Record = bRecord;
    a.Groups = bGroups;
    a.LogChannel = bLogChannel;
    a.MainSite = bMainSite;
    a.OauthAuthUrl = bOauthAuthUrl;
    a.OauthQueryUrl = bOauthQueryUrl;
    a.OauthQueryKey = bOauthQueryKey;
    a.WikiList = bWikiList;
    a.Blacklist = bBlacklist;
    a.MessageStart = bMessageStart;
    a.MessagePolicy = bMessagePolicy;
    a.MessageInsufficientRight = bMessageInsufficientRight;
    a.MessageGeneralPrompt = bMessageGeneralPrompt;
    a.MessageTelegramIdError = bMessageTelegramIdError;
    a.MessageRestoreSilence = bMessageRestoreSilence;
    a.MessageConfirmAlready = bMessageConfirmAlready;
    a.MessageConfirmOtherTg = bMessageConfirmOtherTg;
    a.MessageConfirmConflict = bMessageConfirmConflict;
    a.MessageConfirmChecking = bMessageConfirmChecking;
    a.MessageConfirmUserNotFound = bMessageConfirmUserNotFound;
    a.MessageConfirmButton = bMessageConfirmButton;
    a.MessageConfirmWait = bMessageConfirmWait;
    a.MessageConfirmConfirming = bMessageConfirmConfirming;
    a.MessageConfirmIneligible = bMessageConfirmIneligible;
    a.MessageConfirmSessionLost = bMessageConfirmSessionLost;
    a.MessageConfirmComplete = bMessageConfirmComplete;
    a.MessageConfirmFailed = bMessageConfirmFailed;
    a.MessageConfirmLog = bMessageConfirmLog;
    a.MessageDeconfirmPrompt = bMessageDeconfirmPrompt;
    a.MessageDeconfirmButton = bMessageDeconfirmButton;
    a.MessageDeconfirmSucc = bMessageDeconfirmSucc;
    a.MessageDeconfirmNotConfirmed = bMessageDeconfirmNotConfirmed;
    a.MessageDeconfirmLog = bMessageDeconfirmLog;
    a.MessageNewMemberHint = bMessageNewMemberHint;
    a.MessageAddWhitelistPrompt = bMessageAddWhitelistPrompt;
    a.MessageAddWhitelistSucc = bMessageAddWhitelistSucc;
    a.MessageAddWhitelistLog = bMessageAddWhitelistLog;
    a.MessageRemoveWhitelistPrompt = bMessageRemoveWhitelistPrompt;
    a.MessageRemoveWhitelistNotFound = bMessageRemoveWhitelistNotFound;
    a.MessageRemoveWhitelistLog = bMessageRemoveWhitelistLog;
    a.MessageRemoveWhitelistSucc = bMessageRemoveWhitelistSucc;
    a.MessageWhoisHead = bMessageWhoisHead;
    a.MessageWhoisPrompt = bMessageWhoisPrompt;
    a.MessageWhoisNotFound = bMessageWhoisNotFound;
    a.MessageWhoisSelf = bMessageWhoisSelf;
    a.MessageWhoisBot = bMessageWhoisBot;
    a.MessageWhoisHasMw = bMessageWhoisHasMw;
    a.MessageWhoisNoMw = bMessageWhoisNoMw;
    a.MessageWhoisWhitelisted = bMessageWhoisWhitelisted;
    a.MessageWhoisTgNameUnavailable = bMessageWhoisTgNameUnavailable;
    a.MessageRefuseLog = bMessageRefuseLog;
    a.MessageAcceptLog = bMessageAcceptLog;
    a.MessageLiftRestrictionAlert = bMessageLiftRestrictionAlert;
    a.MessageSilenceAlert = bMessageSilenceAlert;
    a.MessageEnable = bMessageEnable;
    a.MessageDisable = bMessageDisable;
    a.MessageEnableLog = bMessageEnableLog;
    a.MessageDisableLog = bMessageDisableLog;
    return missing, nil;
}

func (a *ConfigMix) Write() string {
    coll := "";
    coll += a.Token.Write() + "\n";
    coll += a.Proxy.Write() + "\n";
    coll += a.Record.Write() + "\n";
    coll += a.Groups.Write() + "\n";
    coll += a.LogChannel.Write() + "\n";
    coll += a.MainSite.Write() + "\n";
    coll += a.OauthAuthUrl.Write() + "\n";
    coll += a.OauthQueryUrl.Write() + "\n";
    coll += a.OauthQueryKey.Write() + "\n";
    coll += a.WikiList.Write() + "\n";
    coll += a.Blacklist.Write() + "\n";
    coll += a.MessageStart.Write() + "\n";
    coll += a.MessagePolicy.Write() + "\n";
    coll += a.MessageInsufficientRight.Write() + "\n";
    coll += a.MessageGeneralPrompt.Write() + "\n";
    coll += a.MessageTelegramIdError.Write() + "\n";
    coll += a.MessageRestoreSilence.Write() + "\n";
    coll += a.MessageConfirmAlready.Write() + "\n";
    coll += a.MessageConfirmOtherTg.Write() + "\n";
    coll += a.MessageConfirmConflict.Write() + "\n";
    coll += a.MessageConfirmChecking.Write() + "\n";
    coll += a.MessageConfirmUserNotFound.Write() + "\n";
    coll += a.MessageConfirmButton.Write() + "\n";
    coll += a.MessageConfirmWait.Write() + "\n";
    coll += a.MessageConfirmConfirming.Write() + "\n";
    coll += a.MessageConfirmIneligible.Write() + "\n";
    coll += a.MessageConfirmSessionLost.Write() + "\n";
    coll += a.MessageConfirmComplete.Write() + "\n";
    coll += a.MessageConfirmFailed.Write() + "\n";
    coll += a.MessageConfirmLog.Write() + "\n";
    coll += a.MessageDeconfirmPrompt.Write() + "\n";
    coll += a.MessageDeconfirmButton.Write() + "\n";
    coll += a.MessageDeconfirmSucc.Write() + "\n";
    coll += a.MessageDeconfirmNotConfirmed.Write() + "\n";
    coll += a.MessageDeconfirmLog.Write() + "\n";
    coll += a.MessageNewMemberHint.Write() + "\n";
    coll += a.MessageAddWhitelistPrompt.Write() + "\n";
    coll += a.MessageAddWhitelistSucc.Write() + "\n";
    coll += a.MessageAddWhitelistLog.Write() + "\n";
    coll += a.MessageRemoveWhitelistPrompt.Write() + "\n";
    coll += a.MessageRemoveWhitelistNotFound.Write() + "\n";
    coll += a.MessageRemoveWhitelistLog.Write() + "\n";
    coll += a.MessageRemoveWhitelistSucc.Write() + "\n";
    coll += a.MessageWhoisHead.Write() + "\n";
    coll += a.MessageWhoisPrompt.Write() + "\n";
    coll += a.MessageWhoisNotFound.Write() + "\n";
    coll += a.MessageWhoisSelf.Write() + "\n";
    coll += a.MessageWhoisBot.Write() + "\n";
    coll += a.MessageWhoisHasMw.Write() + "\n";
    coll += a.MessageWhoisNoMw.Write() + "\n";
    coll += a.MessageWhoisWhitelisted.Write() + "\n";
    coll += a.MessageWhoisTgNameUnavailable.Write() + "\n";
    coll += a.MessageRefuseLog.Write() + "\n";
    coll += a.MessageAcceptLog.Write() + "\n";
    coll += a.MessageLiftRestrictionAlert.Write() + "\n";
    coll += a.MessageSilenceAlert.Write() + "\n";
    coll += a.MessageEnable.Write() + "\n";
    coll += a.MessageDisable.Write() + "\n";
    coll += a.MessageEnableLog.Write() + "\n";
    coll += a.MessageDisableLog.Write() + "\n";
    return coll;
}
