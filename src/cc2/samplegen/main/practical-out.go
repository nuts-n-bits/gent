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
			if len(cur) > 0 {
				coll = append(coll, cur);
				cur = []Token{};
			}
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
		i += 1;
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

type ccCore struct {}

var CcCore = ccCore{};

func (_ ccCore) CoreDynParse(src string) ([]Command, int, error) {
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

func (_ ccCore) CoreDynEncode(cmd Command) string {
	var b strings.Builder;
	b.WriteString(encodeStr(cmd.command));
	for _, arg := range cmd.args {
		b.WriteString(" ");
		b.WriteString(encodeStr(arg));
	}
	if len(cmd.options) % 2 != 0 {
		cmd.options = append(cmd.options, "");
	}
	for i:=0; i<len(cmd.options); i+=2 {
		b.WriteString(" ");
		b.WriteString(cmd.options[i]);
		if cmd.options[i+1] != "" {
			b.WriteString(" ");
			b.WriteString(encodeStr(cmd.options[i+1]));
		}
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

// BEGIN MACHINE GENERATED CODE
type Arinit struct {
    Namespace string
}

func (a *Arinit) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bNamespace string = "";
    var cNamespace int = 0;
    if b.command != "arinit" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        case "-ns", "--namespace":
            bNamespace = b.options[i+1];
            cNamespace += 1;
        }
    }
    if cNamespace < 1 {
        missing += "-ns --namespace;";
    }
    a.Namespace = bNamespace;
    return missing, nil;
}

func (a Arinit) CoreEncode() Command {
    args := []string{};
    options := []string{};
    options = append(options, "-ns", a.Namespace);
    return Command{
        command: "arinit",
        args: args,
        options: options,
    };
}

func (a *Arinit) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Arinit) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Arsync struct {
    CaseId string
    CaseNumber string
}

func (a *Arsync) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bCaseId string = "";
    var cCaseId int = 0;
    var bCaseNumber string = "";
    var cCaseNumber int = 0;
    if b.command != "arsync" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        case "-i", "--case-id":
            bCaseId = b.options[i+1];
            cCaseId += 1;
        case "-n", "--case-number":
            bCaseNumber = b.options[i+1];
            cCaseNumber += 1;
        }
    }
    if cCaseId < 1 {
        missing += "-i --case-id;";
    }
    if cCaseNumber < 1 {
        missing += "-n --case-number;";
    }
    a.CaseId = bCaseId;
    a.CaseNumber = bCaseNumber;
    return missing, nil;
}

func (a Arsync) CoreEncode() Command {
    args := []string{};
    options := []string{};
    options = append(options, "-i", a.CaseId);
    options = append(options, "-n", a.CaseNumber);
    return Command{
        command: "arsync",
        args: args,
        options: options,
    };
}

func (a *Arsync) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Arsync) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Addrev struct {
    Revid string
    RevTimestamp string
    Uid string
    Uname string
    Summary string
    Content string
}

func (a *Addrev) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bRevid string = "";
    var cRevid int = 0;
    var bRevTimestamp string = "";
    var cRevTimestamp int = 0;
    var bUid string = "";
    var cUid int = 0;
    var bUname string = "";
    var cUname int = 0;
    var bSummary string = "";
    var cSummary int = 0;
    var bContent string = "";
    var cContent int = 0;
    if b.command != "addrev" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        case "-r", "--revid":
            bRevid = b.options[i+1];
            cRevid += 1;
        case "-t", "--rev-timestamp":
            bRevTimestamp = b.options[i+1];
            cRevTimestamp += 1;
        case "-u", "--uid":
            bUid = b.options[i+1];
            cUid += 1;
        case "-un", "--uname":
            bUname = b.options[i+1];
            cUname += 1;
        case "-s", "--summary":
            bSummary = b.options[i+1];
            cSummary += 1;
        case "-c", "--content":
            bContent = b.options[i+1];
            cContent += 1;
        }
    }
    if cRevid < 1 {
        missing += "-r --revid;";
    }
    if cRevTimestamp < 1 {
        missing += "-t --rev-timestamp;";
    }
    if cUid < 1 {
        missing += "-u --uid;";
    }
    if cUname < 1 {
        missing += "-un --uname;";
    }
    if cSummary < 1 {
        missing += "-s --summary;";
    }
    if cContent < 1 {
        missing += "-c --content;";
    }
    a.Revid = bRevid;
    a.RevTimestamp = bRevTimestamp;
    a.Uid = bUid;
    a.Uname = bUname;
    a.Summary = bSummary;
    a.Content = bContent;
    return missing, nil;
}

func (a Addrev) CoreEncode() Command {
    args := []string{};
    options := []string{};
    options = append(options, "-r", a.Revid);
    options = append(options, "-t", a.RevTimestamp);
    options = append(options, "-u", a.Uid);
    options = append(options, "-un", a.Uname);
    options = append(options, "-s", a.Summary);
    options = append(options, "-c", a.Content);
    return Command{
        command: "addrev",
        args: args,
        options: options,
    };
}

func (a *Addrev) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Addrev) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Addpid struct {
    Args string;
}

func (a *Addpid) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "addpid" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a Addpid) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "addpid",
        args: args,
        options: options,
    };
}

func (a *Addpid) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Addpid) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Addpc struct {
    Args string;
}

func (a *Addpc) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    if b.command != "addpc" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        }
    }
    a.Args = bArgs; 
    return missing, nil;
}

func (a Addpc) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    return Command{
        command: "addpc",
        args: args,
        options: options,
    };
}

func (a *Addpc) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Addpc) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Arclose struct {
    Args string;
    CaseId string
    CaseNumber string
}

func (a *Arclose) CoreParse(b Command) (missingFields string, err error) {
    missing := "";
    var bArgs string;
    if len(b.args) == 0 {
        missing += "args;";
    }
    bArgs = b.args[0];
    var bCaseId string = "";
    var cCaseId int = 0;
    var bCaseNumber string = "";
    var cCaseNumber int = 0;
    if b.command != "arclose" {
        return missing, fmt.Errorf("command name mismatch");
    }
    for i:=0; i<len(b.options); i+=2 {
        switch b.options[i] {
        case "-i", "--case-id":
            bCaseId = b.options[i+1];
            cCaseId += 1;
        case "-n", "--case-number":
            bCaseNumber = b.options[i+1];
            cCaseNumber += 1;
        }
    }
    if cCaseId < 1 {
        missing += "-i --case-id;";
    }
    if cCaseNumber < 1 {
        missing += "-n --case-number;";
    }
    a.Args = bArgs; 
    a.CaseId = bCaseId;
    a.CaseNumber = bCaseNumber;
    return missing, nil;
}

func (a Arclose) CoreEncode() Command {
    args := []string{};
    options := []string{};
    args = append(args, a.Args);
    options = append(options, "-i", a.CaseId);
    options = append(options, "-n", a.CaseNumber);
    return Command{
        command: "arclose",
        args: args,
        options: options,
    };
}

func (a *Arclose) Parse(b string) (missingFields string, err error) {
    cmds, _, err := CcCore.CoreDynParse(b);
    if err != nil {
        return "", err;
    }
    if len(cmds) != 1 {
        return "", fmt.Errorf("expected exactly 1 command line");
    }
    return a.CoreParse(cmds[0]);
}

func (a Arclose) Write() string {
    return CcCore.CoreDynEncode(a.CoreEncode());
}

type Upload1 struct {
    Arinit Arinit
}

func (a *Upload1) Parse(wiredata string) (string, error) {
    res, _, err := CcCore.CoreDynParse(wiredata);
    if err != nil { return "", err; }
    missing := "";
    mixMissing := ""
    var bArinit Arinit;
    var mArinit = "";
    var cArinit = 0;
    for _, cmd := range res {
        switch cmd.command {
        case "arinit":
            var parsed Arinit;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bArinit = parsed;
            mArinit += mis;
            cArinit += 1;
        }
    }
    if cArinit < 1 {
        mixMissing += "arinit;";
    }
    if len(mixMissing) > 0 { missing += "mix:" + mixMissing + "\n"; }
    if len(mArinit) > 0 { missing += "arinit:" + mArinit + "\n"; }
    a.Arinit = bArinit;
    return missing, nil;
}

func (a *Upload1) Write() string {
    coll := "";
    coll += a.Arinit.Write() + "\n";
    return coll;
}
type Upload2 struct {
    Arsync Arsync
    Addrev []Addrev
    Addpid []Addpid
    Addpc []Addpc
}

func (a *Upload2) Parse(wiredata string) (string, error) {
    res, _, err := CcCore.CoreDynParse(wiredata);
    if err != nil { return "", err; }
    missing := "";
    mixMissing := ""
    var bArsync Arsync;
    var mArsync = "";
    var cArsync = 0;
    var bAddrev []Addrev = []Addrev{};
    var mAddrev = "";
    var cAddrev = 0;
    var bAddpid []Addpid = []Addpid{};
    var mAddpid = "";
    var cAddpid = 0;
    var bAddpc []Addpc = []Addpc{};
    var mAddpc = "";
    var cAddpc = 0;
    for _, cmd := range res {
        switch cmd.command {
        case "arsync":
            var parsed Arsync;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bArsync = parsed;
            mArsync += mis;
            cArsync += 1;
        case "addrev":
            var parsed Addrev;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bAddrev = append(bAddrev, parsed);
            mAddrev += mis;
            cAddrev += 1;
        case "addpid":
            var parsed Addpid;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bAddpid = append(bAddpid, parsed);
            mAddpid += mis;
            cAddpid += 1;
        case "addpc":
            var parsed Addpc;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bAddpc = append(bAddpc, parsed);
            mAddpc += mis;
            cAddpc += 1;
        }
    }
    if cArsync < 1 {
        mixMissing += "arsync;";
    }
    if len(mixMissing) > 0 { missing += "mix:" + mixMissing + "\n"; }
    if len(mArsync) > 0 { missing += "arsync:" + mArsync + "\n"; }
    if len(mAddrev) > 0 { missing += "addrev:" + mAddrev + "\n"; }
    if len(mAddpid) > 0 { missing += "addpid:" + mAddpid + "\n"; }
    if len(mAddpc) > 0 { missing += "addpc:" + mAddpc + "\n"; }
    a.Arsync = bArsync;
    a.Addrev = bAddrev;
    a.Addpid = bAddpid;
    a.Addpc = bAddpc;
    return missing, nil;
}

func (a *Upload2) Write() string {
    coll := "";
    coll += a.Arsync.Write() + "\n";
    for _, b := range a.Addrev {
        coll += b.Write() + "\n";
    }
    for _, b := range a.Addpid {
        coll += b.Write() + "\n";
    }
    for _, b := range a.Addpc {
        coll += b.Write() + "\n";
    }
    return coll;
}
type Upload3 struct {
    Arclose Arclose
}

func (a *Upload3) Parse(wiredata string) (string, error) {
    res, _, err := CcCore.CoreDynParse(wiredata);
    if err != nil { return "", err; }
    missing := "";
    mixMissing := ""
    var bArclose Arclose;
    var mArclose = "";
    var cArclose = 0;
    for _, cmd := range res {
        switch cmd.command {
        case "arclose":
            var parsed Arclose;
            mis, err := parsed.CoreParse(cmd);
            if err != nil { return "", err; }
            bArclose = parsed;
            mArclose += mis;
            cArclose += 1;
        }
    }
    if cArclose < 1 {
        mixMissing += "arclose;";
    }
    if len(mixMissing) > 0 { missing += "mix:" + mixMissing + "\n"; }
    if len(mArclose) > 0 { missing += "arclose:" + mArclose + "\n"; }
    a.Arclose = bArclose;
    return missing, nil;
}

func (a *Upload3) Write() string {
    coll := "";
    coll += a.Arclose.Write() + "\n";
    return coll;
}
