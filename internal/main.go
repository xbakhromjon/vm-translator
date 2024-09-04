package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filepath := GetArg("-f")

	log.Printf("Reading file %s", filepath)

	vmFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("open input vm file: %s", err.Error())
	}

	parser := NewParser(vmFile)

	defer func() {
		err := parser.Close()
		if err != nil {
			log.Fatalf("close input vm file: %s", err.Error())
		}
	}()

	outputFilePath := strings.Replace(filepath, ".vm", ".asm", 1)

	outputFile, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("open outputFile: %s", err.Error())
		}
		outputFile, err = os.Create(outputFilePath)
		if err != nil {
			log.Fatalf("create outputFile: %s", err.Error())
		}
	} else {
		err := outputFile.Truncate(0)
		if err != nil {
			log.Fatalf("truncate outputFile: %s", err.Error())
		}
	}

	codeWriter := NewCode(outputFile)

	defer func() {
		err := codeWriter.Close()
		if err != nil {
			log.Fatalf("close asm outputFile: %s", err.Error())
		}
	}()

	// starting translating process
	for ok := parser.Advance(); ok; ok = parser.Advance() {
		cmdType, err := parser.CommandType()
		if err != nil {
			log.Fatalf("getting command type: %s", err.Error())
		}
		fmt.Printf("cmdType: %s \n", cmdType)

		switch cmdType {
		case C_PUSH:
			arg1, err := parser.Arg1()
			if err != nil {
				log.Fatalf("get arg1 value: %s", err.Error())
			}
			log.Printf("push arg1 %v \n", arg1)

			arg2, err := parser.Arg2()
			if err != nil {
				log.Fatalf("get arg2 value: %s", err.Error())
			}
			log.Printf("push arg2 %v \n", arg2)
			err = codeWriter.WritePush(arg1, arg2)
			if err != nil {
				log.Fatalf("write push value: %s", err.Error())
			}
		case C_ARITHMETIC:
			cmd, err := parser.Command()
			if err != nil {
				log.Fatalf("get command: %s", err.Error())
			}
			err = codeWriter.WriteArithmetic(cmd)
			if err != nil {
				log.Fatalf("write arithmetic value: %s", err.Error())
			}
		}

		log.Printf("parsed line: %s \n", parser.current)
	}

}
