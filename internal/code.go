package main

import (
	"fmt"
	"os"
)

const (
	pushConstant = `@%d
D=A
@SP
A=M
M=D
@SP
M=M+1`

	add = `@SP
M=M-1
A=M
D=M
A=A-1
M=D+M`
)

type Code struct {
	outputFile *os.File
}

func NewCode(file *os.File) *Code {
	return &Code{outputFile: file}
}

func (c *Code) WriteArithmetic(command string) error {
	asm := ""
	switch command {
	case "add":
		asm = add
	}
	return c.writeAsm(asm)
}

func (c *Code) WritePush(segment string, i uint32) error {
	asm := ""
	switch segment {
	case "constant":
		asm = fmt.Sprintf(pushConstant, i)
	}
	fmt.Printf("code: write push: asm=%s\n", asm)
	return c.writeAsm(asm)
}

func (c *Code) writeAsm(asm string) error {
	_, err := c.outputFile.Write([]byte(asm + "\r\n"))
	if err != nil {
		return fmt.Errorf("write asm: %s", err.Error())
	}
	return nil
}

func (c *Code) Close() error {
	return c.outputFile.Close()
}
