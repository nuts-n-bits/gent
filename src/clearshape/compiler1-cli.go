package main

import (
	"encoding/json"
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
	tokens, errI, err := lexTokenizer(programStr)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("TOK\n%#v", tokens)
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
	tokens, errI, err := lexTokenizer(programStr)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("TOK\n%#v", tokens)
	program, errI, err := rdParseProgram(tokens)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("\n\n\nAST\n%s", program.DebugString())
	programFlt := fltFlattenProgram(program)
	fmt.Printf("\n\n\nFLT\n%s", programFlt.DebugString())
}

func show_lc(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file")
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file")
	}
	programStr, err := readFile(args.rest[0])
	if err != nil {
		log.Fatalf("ERR: %s", err.Error())
	}
	tokens, errI, err := lexTokenizer(programStr)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("TOK\n%#v", tokens)
	astProgram, errI, err := rdParseProgram(tokens)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), tokens[errI])
	}
	fmt.Printf("\n\n\nAST\n%s", astProgram.DebugString())
	errT, err := lcCheckProgram1Of2CheckReservedName(astProgram)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), errT)
	}
	programLc, topLevelCollisions, undefinedRefs := lcCheckProgram2Of2CheckCollisionAndUndefined(astProgram)
	if len(topLevelCollisions) > 0 || len(undefinedRefs) > 0 {
		if len(topLevelCollisions) > 0 {
			log.Printf("\n\nDetected %d duplicate identifiers: %#v\n", len(topLevelCollisions), topLevelCollisions)
		}
		if len(undefinedRefs) > 0 {
			log.Printf("\n\nDetected %d undefined references: %#v\n", len(undefinedRefs), undefinedRefs)
		}
		log.Fatalf("exiting due to previous error(s)")
	}
	fmt.Printf("\n\n\nLC\n%s", programLc.DebugString())
}

func show_lnk_ball(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file")
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file")
	}
	linkedBall, errPath, err := lnkGatherSrcFiles(args.rest[0])
	if err != nil {
		errDesc, mErr := json.Marshal(err)
		if mErr != nil {
			panic("shouldn't really happen")
		}
		log.Fatalf("In file %s, encountered error: %s (%s)", errPath, err.ErrToStr(), errDesc)
	}
	str, err1 := json.Marshal(linkedBall)
	if err1 != nil {
		log.Fatalf("Cannot json marshal linked ball")
	}
	log.Printf("\n\n\nLNK-BALL\n%s", str)
}

func show_lnk(args ProgramCliParameters) {
	if len(args.rest) == 0 {
		log.Fatal("ERR: No input file")
	} else if len(args.rest) > 1 {
		log.Fatal("ERR: Multiple input file")
	}
	linkedBall, errPath, errU := lnkGatherSrcFiles(args.rest[0])
	if errU != nil {
		errDesc, mErr := json.Marshal(errU)
		if mErr != nil {
			panic("shouldn't really happen")
		}
		log.Fatalf("In file %s, encountered error: %s (%s)", errPath, errU.ErrToStr(), errDesc)
	}
	str, err := json.Marshal(linkedBall)
	if err != nil {
		log.Fatalf("Cannot json marshal linked ball")
	}
	log.Printf("\n\n\nLNK-BALL\n%s", str)
	lnkProgram, errT, err := lnkResolveImports(linkedBall)
	if err != nil {
		log.Fatalf("ERR: %s (at %#v)", err.Error(), errT)
	}
	fmt.Printf("\n\n\nLNK\n%s", lnkProgram.DebugString())
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
	} else if args.verb == "show-lc" {
		show_lc(args)
	} else if args.verb == "show-lnk-ball" {
		show_lnk_ball(args)
	} else if args.verb == "show-lnk" {
		show_lnk(args)
	}

	fmt.Printf("\n\n")
}
