package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printErr("invalid number of arguments")
	}

	// Open file
	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		printErr(fmt.Sprintf("%s file not exists\n", os.Args[1]))
	}
	defer file.Close()

	base := filepath.Base(os.Args[1])
	fileNames := strings.Split(base, ".")
	destNames := fileNames[0] + ".asm"

	fileName := os.Args[1][:strings.Index(os.Args[1], ".vm")]
	codeFile, err := os.Create(fmt.Sprintf("%v.asm", fileName))
	defer func(fs *os.File) {
		if err := fs.Close(); err != nil {
			printErr(err.Error())
		}
	}(codeFile)
	if err != nil {
		printErr(err.Error())
	}

	if codeFile == nil {
		printErr("failed to create file")
	}

	// Read file
	parser := NewParser(file)
	translator := NewTranslator(destNames)
	i := 1
	for parser.hasMoreCommand() {
		parser.advance()
		cmd := ""
		if parser.CommandType() == CARITHMETIC {
			cmd = translator.WriteArithmetic(parser.Arg1())
		}
		if parser.CommandType() == CPUSH {
			cmd = translator.WritePushPop(CPUSH, parser.Arg1(), parser.Arg2())
		}
		if parser.CommandType() == CPOP {
			cmd = translator.WritePushPop(CPOP, parser.Arg1(), parser.Arg2())
		}
		_, _ = codeFile.WriteString("// " + parser.arg0 + " " + parser.arg1 + " " + parser.arg2 + "\n")
		_, err := codeFile.WriteString(cmd)
		if err != nil {
			printErr(err.Error())
		}
		i++
	}
	codeFile.WriteString("(END)\n" +
		"@END\n" +
		"0;JMP\n")
}

func hasComment(line string) bool {
	return strings.Contains(line, "//")
}

func removeComment(line string) string {
	offset := strings.Index(line, "//")
	return strings.TrimSpace(line[0:offset])
}

func printErr(err string) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

const gotoTopmostStackVal = "@SP\n" +
	"A=M-1\n"

const popIntoD = "D=M\n" + // store val in D
	"A=A-1\n" // move to second top-most

const decrementSp = "@SP\n" +
	"M=M-1\n"

var SegmentPointer = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"temp":     "R5",
	"pointer":  "R3",
}

var labelCount = map[string]int{
	"gt": 0,
	"lt": 0,
	"eq": 0,
}
var cmpFalse = map[string]string{
	"gt": "JLE", // greater than == not less than or equal to
	"lt": "JGE", // less than == not greater than or equal to
	"eq": "JNE",
} // equal == not not equal

type Translator struct {
	fileName string
}

func NewTranslator(fileName string) *Translator {
	return &Translator{
		fileName: fileName,
	}
}

// WritePushPop implements CodeWriter.
func (t *Translator) WritePushPop(cmdType CommandType, segment string, idx int) string {
	cmd := ""
	addr := fmt.Sprintf("%d", idx)

	// PUSH
	if cmdType == CPUSH {
		cmd = t.WritePush(segment, addr)
	} else if cmdType == CPOP {
		cmd = t.WritePop(segment, addr)
	}

	return cmd
}

func (t *Translator) getSegment(segment string, addr string) string {
	sb := ""
	if segment == "static" {
		return fmt.Sprintf("@%v.%v\n", t.fileName, addr)
	}

	// D = addr
	sb += "@" + addr + "\n" +
		"D=A\n"

	if segment == "temp" || segment == "pointer" {
		sb +=
			"@" + SegmentPointer[segment] + "\n" +
				"A=D+A\n"

	} else {
		// A=Base+addr, M=RAM[base + addr]
		sb += "@" + SegmentPointer[segment] + "\n" +
			"A=D+M\n"
	}

	return sb
}

func (t *Translator) WritePush(segment string, addr string) string {
	sb := ""

	if segment == "constant" {
		sb +=
			"@" + addr + "\n" +
				"D=A\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n"
	} else {
		sb += t.getSegment(segment, addr)
		// assign D=RAM[base + addr]
		sb +=
			"D=M\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n"
	}

	// increment stack pointer
	sb +=
		"@SP\n" +
			"M=M+1\n"

	return sb
}

func (t *Translator) WritePop(segment string, addr string) string {
	sb := ""

	sb +=
		t.getSegment(segment, addr) +
			"D=A\n" +
			"@R13\n" +
			"M=D\n"

	sb += gotoTopmostStackVal

	sb +=
		"D=M\n" +
			"@R13\n" +
			"A=M\n" +
			"M=D\n"

	sb += decrementSp

	return sb
}

func (t *Translator) WriteArithmetic(op string) string {
	sb := ""

	switch op {
	case "add":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				"M=D+M\n" +
				decrementSp
	case "sub":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				"M=M-D\n" +
				decrementSp
	case "neg":
		sb +=
			gotoTopmostStackVal +
				"M=-M\n"
	case "and":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				"M=D&M\n" +
				decrementSp
	case "or":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				"M=D|M\n" +
				decrementSp
	case "not":
		sb +=
			gotoTopmostStackVal +
				"M=!M\n"
	case "eq":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				t.writeComparison("eq") +
				decrementSp
	case "gt":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				t.writeComparison("gt") +
				decrementSp
	case "lt":
		sb +=
			gotoTopmostStackVal +
				popIntoD +
				t.writeComparison("lt") +
				decrementSp
	default:
		panic(fmt.Sprintf("Command %s is not valid", op))
	}
	return sb
}

func (t *Translator) writeComparison(operator string) string {
	sb := ""
	jumpNot := fmt.Sprintf("NOT_%v_%d", operator, labelCount[operator])
	jumpEnd := fmt.Sprintf("END_%v_%d", operator, labelCount[operator])
	labelCount[operator]++
	sb +=
		"D=M-D\n" + // 1st top most - 2nd top most of stack
			"@" + jumpNot + "\n" +
			"D;" + cmpFalse[operator] + "\n" + // if true => write -1
			gotoTopmostStackVal +
			"A=A-1\n" +
			"M=-1\n" +
			"@" + jumpEnd + "\n" +
			"0;JMP\n" +
			"(" + jumpNot + ")\n" +
			gotoTopmostStackVal +
			"A=A-1\n" +
			"M=0\n" +
			"(" + jumpEnd + ")\n"

	return sb
}

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
	if hasComment(line) {
		line = removeComment(line)
	}
	if len(line) == 0 {
		return
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
