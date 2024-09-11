package main

import (
	"fmt"
	"os"
)

const (
	pushLocal = `@LCL
D=M
@%d
D=D+A
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1`

	pushArgument = `@ARG
D=M
@%d
D=D+A
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1`

	pushThis = `@THIS
D=M
@%d
D=D+A
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1`

	pushThat = `@THAT
D=M
@%d
D=D+A
A=D
D=M
@SP
A=M
M=D
@SP
M=M+1`

	pushConstant = `@%d
D=M
@SP
A=M
M=D
@SP
M=M+1`

	pushTemp = `@5
D=M
@%d
D=D+A
A=D
D=M
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

	sub = `@SP
M=M-1
A=M
D=M
A=A-1
M=M-D`

	and = `@SP
M=M-1
A=M
D=M
A=A-1
M=D&M`

	or = `@SP
M=M-1
A=M
D=M
A=A-1
M=D|M`

	not = `@SP
A=M-1
D=M
@NOT_LABEL_%d
D-1;JGE
@SP
A=M-1
M=1
@NOT_END_%d
0;JMP
(NOT_LABEL_%d)
@SP
A=M-1
M=0
(NOT_END_%d)`

	neg = `@SP
A=M-1
M=-M`

	eq = `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M      
D=M-D
@EQ_LABEL_%d
D=D;JEQ
@SP
A=M
M=0
@EQ_END_%d
0;JEQ
(EQ_LABEL_%d) 
@SP
A=M
M=1
(EQ_END_%d)
@SP
M=M+1`

	gt = `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M      
D=M-D
@GT_LABEL_%d
D=D;JGT
@SP
A=M
M=0
@GT_END_%d
0;JEQ
(GT_LABEL_%d) 
@SP
A=M
M=1
(GT_END_%d)
@SP
M=M+1`

	lt = `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M      
D=M-D
@LT_LABEL_%d
D=D;JLT
@SP
A=M
M=0
@LT_END_%d
0;JEQ
(LT_LABEL_%d) 
@SP
A=M
M=1
(LT_END_%d)
@SP
M=M+1`
)

var (
	eqLabelCounter  = 0
	gtLabelCounter  = 0
	ltLabelCounter  = 0
	notLabelCounter = 0
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
	case "sub":
		asm = sub
	case "neg":
		asm = neg
	case "eq":
		asm = fmt.Sprintf(eq, eqLabelCounter, eqLabelCounter, eqLabelCounter, eqLabelCounter)
		eqLabelCounter++
	case "gt":
		asm = fmt.Sprintf(gt, gtLabelCounter, gtLabelCounter, gtLabelCounter, gtLabelCounter)
		gtLabelCounter++
	case "lt":
		asm = fmt.Sprintf(lt, ltLabelCounter, ltLabelCounter, ltLabelCounter, ltLabelCounter)
		ltLabelCounter++
	case "and":
		asm = and
	case "or":
		asm = or
	case "not":
		asm = fmt.Sprintf(not, notLabelCounter, notLabelCounter, notLabelCounter, notLabelCounter)
		notLabelCounter++
	}
	return c.writeAsm(asm)
}

func (c *Code) WritePush(segment string, i uint32) error {
	asm := ""
	switch segment {
	case "local":
		asm = fmt.Sprintf(pushLocal, i)
	case "argument":
		asm = fmt.Sprintf(pushArgument, i)
	case "this":
		asm = fmt.Sprintf(pushThis, i)
	case "that":
		asm = fmt.Sprintf(pushThat, i)
	case "constant":
		asm = fmt.Sprintf(pushConstant, i)
	case "temp":
		asm = fmt.Sprintf(pushTemp, i)

	}
	fmt.Printf("code: write push: asm=%s\n", asm)
	return c.writeAsm(asm)
}

func (c *Code) WritePop(segment string, i uint32) error {
	asm := ""
	switch segment {
	case "local":
		asm = fmt.Sprintf(pushLocal, i)
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
