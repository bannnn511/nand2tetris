package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var ShouldCallSysInit = false

func main() {
	if len(os.Args) < 2 {
		printErr("invalid number of arguments")
	}

	vmFiles := make([]os.File, 0)

	// Open file
	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		printErr(fmt.Sprintf("%s file not exists\n", os.Args[1]))
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		printErr(err.Error())
	}

	outFile := ""
	if fileInfo.IsDir() {
		// is a directory
		vmFiles = getVmFiles(os.Args[1])
		outFile = fmt.Sprintf("%v/%v.asm", os.Args[1], fileInfo.Name())
		defer func(fss []os.File) {
			for _, fs := range fss {
				if err := fs.Close(); err != nil {
					printErr(err.Error())
				}
			}
		}(vmFiles)
	} else {
		// is a file
		if err != nil {
			printErr(err.Error())
		}
		vmFiles = append(vmFiles, *file)
		fileName := os.Args[1][:strings.Index(os.Args[1], ".vm")]
		outFile = fmt.Sprintf("%v.asm", fileName)
	}

	codeFile, err := os.Create(outFile)

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
	translator := NewTranslator(outFile)

	i := 0
	for _, file := range vmFiles {
		parser := NewParser(&file)
		for parser.hasMoreCommand() {
			cmd := ""
			if i == 0 {
				cmd = translator.WriteInit()
				_, _ = codeFile.WriteString(cmd)
				i++
				continue
			}

			parser.advance()

			switch parser.CommandType() {
			case CARITHMETIC:
				cmd = translator.WriteArithmetic(parser.Arg1())
			case CPUSH:
				cmd = translator.WritePushPop(CPUSH, parser.Arg1(), parser.Arg2())
			case CPOP:
				cmd = translator.WritePushPop(CPOP, parser.Arg1(), parser.Arg2())
			case CLABEL:
				cmd = translator.WriteLabel(parser.Arg1())
			case CGOTO:
				cmd = translator.WriteGoto(parser.Arg1())
			case CIF:
				cmd = translator.WriteIfGoto(parser.Arg1())
			case CFUNCTION:
				cmd = translator.WriteFunction(parser.Arg1(), parser.Arg2())
			case CRETURN:
				cmd = translator.WriteReturn()
			case CCALL:
				cmd = translator.WriteCall(parser.Arg1(), parser.Arg2())
			}

			// for debugging command
			_, _ = codeFile.WriteString("// " + parser.arg0 + " " + parser.arg1 + " " + parser.arg2 + "\n")
			_, err := codeFile.WriteString(cmd)
			if err != nil {
				printErr(err.Error())
			}
			i++
		}

	}

	_, err = codeFile.WriteString("(END_PROGRAM)\n" +
		"@END_PROGRAM\n" +
		"0;JMP\n")
	if err != nil {
		printErr(err.Error())
	}
}

func getVmFiles(dir string) []os.File {
	files, err := os.ReadDir(dir)
	vmFiles := make([]os.File, 0, len(files))
	if err != nil {
		printErr(err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".vm") {
			file, err := os.OpenFile(dir+"/"+file.Name(), os.O_RDONLY, 0)
			if err != nil {
				printErr(err.Error())
			}
			vmFiles = append(vmFiles, *file)
		}
	}

	return vmFiles
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
	"gt":     0,
	"lt":     0,
	"eq":     0,
	"return": 0,
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

func (t *Translator) getSegment(segment string, addr string) string {
	var sb strings.Builder
	if segment == "static" {
		return fmt.Sprintf("@%v.%v\n", t.fileName, addr)
	}

	// D = addr
	sb.WriteString(
		"@" + addr + "\n" +
			"D=A\n")

	if segment == "temp" || segment == "pointer" {
		sb.WriteString(
			"@" + SegmentPointer[segment] + "\n" +
				"A=D+A\n",
		)
	} else {
		// A=Base+addr, M=RAM[base + addr]
		sb.WriteString(
			"@" + SegmentPointer[segment] + "\n" +
				"A=D+M\n",
		)
	}

	return sb.String()
}

func (t *Translator) initSegment() string {
	var sb strings.Builder

	sb.WriteString("@256\n" +
		"D=A\n" +
		"@SP\n" +
		"M=D\n")

	sb.WriteString("@261\n" +
		"D=A\n" +
		"@SP\n" +
		"M=D\n" +
		"@Sys.init\n" +
		"0;JMP\n")

	return sb.String()
}

func (t *Translator) WriteInit() string {
	var sb strings.Builder

	if ShouldCallSysInit {
		sb.WriteString(t.initSegment())
	}

	return sb.String()
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

func (t *Translator) WritePush(segment string, addr string) string {
	var sb strings.Builder

	if segment == "constant" {
		sb.WriteString(
			"@" + addr + "\n" +
				"D=A\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n")
	} else {
		sb.WriteString(t.getSegment(segment, addr))
		// assign D=RAM[base + addr]
		sb.WriteString(
			"D=M\n" +
				"@SP\n" +
				"A=M\n" +
				"M=D\n")
	}

	// increment stack pointer
	sb.WriteString(
		"@SP\n" +
			"M=M+1\n")

	return sb.String()
}

func (t *Translator) WritePop(segment string, addr string) string {
	var sb strings.Builder

	sb.WriteString(
		t.getSegment(segment, addr) +
			"D=A\n" +
			"@R13\n" +
			"M=D\n")

	sb.WriteString(gotoTopmostStackVal)

	sb.WriteString(
		"D=M\n" +
			"@R13\n" +
			"A=M\n" +
			"M=D\n")

	sb.WriteString(decrementSp)

	return sb.String()
}

func (t *Translator) WriteArithmetic(op string) string {
	var sb strings.Builder

	switch op {
	case "add":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				"M=D+M\n" +
				decrementSp)
	case "sub":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				"M=M-D\n" +
				decrementSp)
	case "neg":
		sb.WriteString(
			gotoTopmostStackVal +
				"M=-M\n")
	case "and":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				"M=D&M\n" +
				decrementSp)
	case "or":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				"M=D|M\n" +
				decrementSp)
	case "not":
		sb.WriteString(
			gotoTopmostStackVal +
				"M=!M\n")
	case "eq":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				t.writeComparison("eq") +
				decrementSp)
	case "gt":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				t.writeComparison("gt") +
				decrementSp)
	case "lt":
		sb.WriteString(
			gotoTopmostStackVal +
				popIntoD +
				t.writeComparison("lt") +
				decrementSp)
	default:
		panic(fmt.Sprintf("Command %s is not valid", op))
	}
	return sb.String()
}

func (t *Translator) WriteLabel(label string) string {
	return "(" + label + ")\n"
}

func (t *Translator) WriteGoto(label string) string {
	var sb strings.Builder
	sb.WriteString(
		"@" + label + "\n" +
			"0;JMP\n")

	return sb.String()
}

func (t *Translator) WriteIfGoto(label string) string {
	var sb strings.Builder
	sb.WriteString(
		gotoTopmostStackVal +
			popIntoD +
			decrementSp +
			"@" + label + "\n" +
			"D;JNE\n")

	return sb.String()
}

func (t *Translator) WriteFunction(label string, nVars int) string {
	var sb strings.Builder
	sb.WriteString(
		"(" + label + ")\n")

	// init variable to )
	for i := 0; i < nVars; i++ {
		sb.WriteString("@SP\n" +
			"A=M\n" +
			"M=0\n" +
			"@SP\n" +
			"M=M+1\n")
	}

	return sb.String()
}

func (t *Translator) WriteReturn() string {
	var sb strings.Builder
	sb.WriteString("@LCL\n" + "D=M\n" + "@R13\n" + "M=D\n")
	sb.WriteString("@5\n" + "A=D-A\n" + "M=D\n" + "@R14\n" + "M=D\n")
	sb.WriteString(gotoTopmostStackVal +
		popIntoD +
		"@ARG\n" +
		"A=M\n" +
		"M=D\n")
	sb.WriteString("@ARG\n" +
		"D=M+1\n" +
		"@SP\n" +
		"M=D\n")
	sb.WriteString("@R13\n" + "AM=M-1\n" + "D=M\n" + "@THAT\n" + "M=D\n")
	sb.WriteString("@R13\n" + "AM=M-1\n" + "D=M\n" + "@THIS\n" + "M=D\n")
	sb.WriteString("@R13\n" + "AM=M-1\n" + "D=M\n" + "@ARG\n" + "M=D\n")
	sb.WriteString("@R13\n" + "AM=M-1\n" + "D=M\n" + "@LCL\n" + "M=D\n")
	sb.WriteString("@R14\n" + "A=M\n" + "0;JMP\n")

	return sb.String()
}

func (t *Translator) WriteCall(functionName string, nVars int) string {
	var sb strings.Builder
	counter := labelCount[functionName]
	returnAddr := fmt.Sprintf("RETURN_ADDR_%s_%d", functionName, counter)
	labelCount[functionName]++
	sb.WriteString("@" + returnAddr + "\n" + "D=A\n" + "@SP\n" + "A=M\n" + "M=D\n" + "@SP\n" + "M=M+1\n")
	sb.WriteString("@LCL\n" + "D=M\n" + "@SP\n" + "A=M\n" + "M=D\n" + "@SP\n" + "M=M+1\n")
	sb.WriteString("@ARG\n" + "D=M\n" + "@SP\n" + "A=M\n" + "M=D\n" + "@SP\n" + "M=M+1\n")
	sb.WriteString("@THIS\n" + "D=M\n" + "@SP\n" + "A=M\n" + "M=D\n" + "@SP\n" + "M=M+1\n")
	sb.WriteString("@THAT\n" + "D=M\n" + "@SP\n" + "A=M\n" + "M=D\n" + "@SP\n" + "M=M+1\n")
	sb.WriteString("@SP\n" + "D=M\n" + "@5\n" + "D=D-A\n" + "@" + returnAddr + "\n" + "D=D-A\n" + "@ARG\n" + "M=D\n")
	sb.WriteString("@SP\n" + "D=M\n" + "@LCL\n" + "M=D\n")
	sb.WriteString("@" + functionName + "\n" + "0;JMP\n")
	sb.WriteString("(" + returnAddr + ")\n")

	return sb.String()
}

func (t *Translator) writeComparison(operator string) string {
	var sb strings.Builder
	jumpNot := fmt.Sprintf("NOT_%v_%d", operator, labelCount[operator])
	jumpEnd := fmt.Sprintf("END_%v_%d", operator, labelCount[operator])
	labelCount[operator]++
	sb.WriteString(
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
			"(" + jumpEnd + ")\n")

	return sb.String()
}

type CommandType string

const (
	CCMT        = "C_CMT" // comment
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
		p.mCmdType = CCMT
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
	switch cmds[0] {
	case "push":
		p.mCmdType = CPUSH
	case "pop":
		p.mCmdType = CPOP
	case "label":
		p.mCmdType = CLABEL
	case "goto":
		p.mCmdType = CGOTO
	case "if-goto":
		p.mCmdType = CIF
	case "function":
		p.mCmdType = CFUNCTION
	case "return":
		p.mCmdType = CRETURN
	case "call":
		p.mCmdType = CCALL
	default:
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
