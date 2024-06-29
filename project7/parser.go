package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type CommandType string

const (
	CPUSH       = "C_PUSH"
	CPOP        = "C_POP"
	CARITHMETIC = "C_ARITHMETIC"
	CLABEL      = "C_LABEL"
	CGOTO       = "C_GOTO"
	CIF         = "C_IF"
	CFUNCTION   = "C_FUNCTION"
	CRETURN     = "C_RETURN"
	CCALL       = "C_CALL"
)

type Parser struct {
	file     *bufio.Scanner
	curLine  string
	arg0     string
	arg1     string
	arg2     string
	mCmdType CommandType
}

func NewParser(file io.Reader) *Parser {
	return &Parser{
		file: bufio.NewScanner(file),
	}
}

func (p *Parser) CommandType() CommandType {
	return p.mCmdType
}

func (p *Parser) hasMoreCommand() bool {
	return p.file.Scan()
}

func (p *Parser) advance() {
	line := strings.TrimSpace(p.file.Text())
	if line == "" || hasComment(line) {
		line = removeComment(line)
	}
	p.curLine = line
	cmds := strings.Split(line, " ")

	if len(cmds) > 1 {
		p.arg1 = cmds[1]
	}
	if len(cmds) > 2 {
		p.arg2 = cmds[2]
	}

	p.arg0 = cmds[0]
	if cmds[0] == "push" {
		p.mCmdType = CPUSH
	} else if cmds[0] == "pop" {
		p.mCmdType = CPOP
	} else {
		p.mCmdType = CARITHMETIC
	}
}

func (p *Parser) Arg1() string {
	if p.CommandType() == CARITHMETIC {
		return p.arg0
	}

	return p.arg1
}

func (p *Parser) Arg2() int {
	if p.CommandType() == CPUSH ||
		p.CommandType() == CPOP ||
		p.CommandType() == CFUNCTION ||
		p.CommandType() == CCALL {
		value, _ := strconv.Atoi(p.arg2)
		return value
	}
	return 0
}
