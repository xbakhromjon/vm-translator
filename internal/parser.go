package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	C_ARITHMETIC = "c_arithmetic"
	C_PUSH       = "c_push"
	C_POP        = "c_pop"
)

type Parser struct {
	current   string
	inputFile *os.File
	scanner   *bufio.Scanner
}

func NewParser(file *os.File) *Parser {
	return &Parser{inputFile: file, scanner: bufio.NewScanner(file)}
}

func (p *Parser) Advance() bool {
	for {
		if !p.scanner.Scan() {
			log.Printf("parser: advance input file: %s \n", p.scanner.Err())
			return false
		}
		text := p.scanner.Text()
		text = strings.TrimSpace(text)
		log.Printf("parser: advance text: %s \n", text)
		time.Sleep(1 * time.Second)
		if text == "" || strings.HasPrefix(text, "//") {
			continue
		}
		p.current = text
		log.Printf("parser: advance: current line: %s \n", p.current)
		time.Sleep(1 * time.Second)
		return true
	}
}

func (p *Parser) CommandType() (string, error) {
	cmd, err := p.Command()
	if err != nil {
		return "", err
	}
	switch cmd {
	case "push":
		return C_PUSH, nil
	case "pop":
		return C_POP, nil
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return C_ARITHMETIC, nil
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}

func (p *Parser) Command() (string, error) {
	cmd := strings.Split(p.current, " ")[0]
	log.Printf("parser: command: %s", cmd)
	return cmd, nil
}

func (p *Parser) Arg1() (string, error) {
	arg1 := strings.Split(p.current, " ")[1]
	log.Printf("parser: arg1: %s", arg1)
	return arg1, nil
}

func (p *Parser) Arg2() (uint32, error) {
	arg2 := strings.Split(p.current, " ")[2]
	log.Printf("parser: arg2: %s", arg2)
	index, _ := strconv.ParseUint(arg2, 10, 32)
	return uint32(index), nil
}

func (p *Parser) Close() error {
	return p.inputFile.Close()
}
