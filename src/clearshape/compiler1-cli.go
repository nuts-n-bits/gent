package main

import (
	"fmt"
	"log"
	"os"
)

type ProgramCliParameters struct {
	tsOut         string
	mlOut         string
	goOut         string
	goPackageName string
	verb          string
	name          string
	rest          []string
	indent        string
}

func hfcliParseArgs(args []string) ProgramCliParameters {
	programCliParameters := ProgramCliParameters{}
	nameFilled, verbFilled, positionalMode := false, false, false
	isOption := func(s string) bool { return len(s) != 0 && s[0] == '-' }
	for i := 0; i < len(args); {
		if i >= len(args) {
			return programCliParameters
		} else if isOption(args[i]) && !positionalMode {
			if args[i] == "--" {
				positionalMode = true
				i += 1
				continue
			}
			this := args[i]
			next := ""
			if i+1 < len(args) && isOption(args[i+1]) {
				i += 1
			} else if i+1 < len(args) {
				next = args[i+1]
				i += 2
			} else {
				i += 1
			}
			switch this {
			case "--ts-out":
				programCliParameters.tsOut = next
			case "--go-out":
				programCliParameters.goOut = next
			case "--go-package-name":
				programCliParameters.goPackageName = next
			case "--ml-out":
				programCliParameters.mlOut = next
			case "--indent":
				programCliParameters.indent = next
			}
		} else if verbFilled {
			programCliParameters.rest = append(programCliParameters.rest, args[i])
			i += 1
		} else if nameFilled {
			programCliParameters.verb = args[i]
			verbFilled = true
			i += 1
		} else {
			programCliParameters.name = args[i]
			nameFilled = true
			i += 1
		}
	}
	return programCliParameters
}

func readFile(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(fileName string, fileContent string) error {
	return os.WriteFile(fileName, []byte(fileContent), 0666)
}

// func build(args ProgramCliParameters) {
// 	if len(args.rest) == 0 {
// 		log.Fatal("ERR: No input file");
// 	} else if len(args.rest) > 1 {
// 		log.Fatal("ERR: Multiple input file");
// 	}
// 	programStr, err := readFile(args.rest[0]);
// 	if err != nil {
// 		log.Fatalf("ERR: %s", err.Error());
// 	}
// 	tokens, err, errI := lexTokenizer(programStr);
// 	if err != nil {
// 		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), tokens[errI].start, tokens[errI].end);
// 	}
// 	program, errI, err := rdParseProgram(tokens);
// 	if err != nil {
// 		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), tokens[errI].start, tokens[errI].end);
// 	}
// 	checked, errT, err := chkProgram(program);
// 	if err != nil {
// 		log.Fatalf("ERR: %s (at %d-%d)", err.Error(), errT.start, errT.end);
// 	}
// 	// parsing complete
// 	indent := "";
// 	switch args.indent {
// 	case "", "4":
// 		indent = "    ";
// 	case "tab":
// 		indent = "\t";
// 	case "2":
// 		indent = "  ";
// 	default:
// 		log.Fatal("--indent must be `4`, `2`, `tab` or unspecified");
// 	}
// 	//
// 	if args.tsOut != "" {
// 		program := cgProgramTypescript(checked, indent);
// 		writeFile(args.tsOut, program);
// 	}
// 	if args.goOut != "" {
// 		if args.goPackageName == "" {
// 			log.Fatal("--go-package-name must be present when --go-out is specified");
// 		}
// 		program := cgProgramGolang(checked, indent, args.goPackageName);
// 		writeFile(args.goOut, program);
// 	}
// }

func show_ast(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file")
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file")
	}
	programStr, err := readFile(args.rest[0])
	if err != nil {
		log.Fatalf("ERR: %s", err.Error())
	}
	tokens, err, errI := lexTokenizer(programStr)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("TOK\n%#v", tokens);
	program, errI, err := rdParseProgram(tokens)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("\n\n\nAST\n%s", program.DebugString())
}

func show_flt(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file")
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file")
	}
	programStr, err := readFile(args.rest[0])
	if err != nil {
		log.Fatalf("ERR: %s", err.Error())
	}
	tokens, err, errI := lexTokenizer(programStr)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("TOK\n%#v", tokens);
	program, errI, err := rdParseProgram(tokens)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("\n\n\nAST\n%s", program.DebugString())
	programFlt := fltFlattenProgram(program);
	fmt.Printf("\n\n\nFLT\n%s", programFlt.DebugString())
}

func main() {

	args := hfcliParseArgs(os.Args)

	//fmt.Printf("//args: %#v", args);

	if args.verb == "build" && true {
		//build(args);
		fmt.Print("Cannot build yet")
	} else if args.verb == "show-ast" {
		show_ast(args)
	} else if args.verb == "show-flt" {
		show_flt(args)
	}
}
