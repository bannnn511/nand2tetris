package main

import (
	"bufio"
	"fmt"
	"os"
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
}

var arithCmds = [9]string{
	"add", "sub", "neq",
	"eq", "gt", "lt",
	"and", "or", "not",
}

type CodeWriter interface {
	WriteArithmetic(command string)
	WritePushPop(cmdType CommandType, segment string, idx int)
}

var _ CodeWriter = (*CodeWrite)(nil)

type CodeWrite struct {
	file bufio.Writer
}

func NewCodeWriter(file *os.File) *CodeWrite {
	return &CodeWrite{
		file: *bufio.NewWriter(file),
	}
}

func (w CodeWrite) WriteArithmetic(command string) {

}

// WritePushPop implements CodeWriter.
func (w CodeWrite) WritePushPop(cmdType CommandType, segment string, idx int) {
	cmd := ""
	switch cmdType {
	case CPUSH:
		segmentPointer := ""
		switch segment {
		case "local":
			segmentPointer = "LCL"
		case "argument":
			segmentPointer = "ARG"
		case "this":
			segmentPointer = "THIS"
		case "that":
			segmentPointer = "THAT"
		}
		cmd = writePushFormat(segmentPointer, idx)
	}

	_, err := w.file.WriteString(cmd)
	if err != nil {
		printErr(err.Error())
	}
}

func writePushFormat(segment string, idx int) string {
	segmentCode := fmt.Sprintf("@%v", segment)
	return fmt.Sprintf(
		"%d\n%v\n%v\n%v\n%v\n%v\n%v\n",
		idx,
		"D=A",
		segmentCode,
		"A=M",
		"M=D",
		segmentCode,
		"M=M+1",
	)
}

func (w CodeWrite) Close() error {
	if err := w.file.Flush(); err != nil {
		return err
	}

	return nil
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
