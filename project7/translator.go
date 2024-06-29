package main

import (
	"bufio"
	"fmt"
	"os"
)

type CodeWriter interface {
	WriteArithmetic(command string)
	WritePushPop(cmdType CommandType, segment string, idx int)
}

var _ CodeWriter = (*Translator)(nil)

type Translator struct {
	file bufio.Writer
}

func NewCodeWriter(file *os.File) *Translator {
	return &Translator{
		file: *bufio.NewWriter(file),
	}
}

func (w Translator) WriteArithmetic(command string) {

}

// WritePushPop implements CodeWriter.
func (w Translator) WritePushPop(cmdType CommandType, segment string, idx int) {
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
		"@%d\n%v\n%v\n%v\n%v\n%v\n%v\n",
		idx,
		"D=A",
		segmentCode,
		"A=M",
		"M=D",
		segmentCode,
		"M=M+1",
	)
}

func (w Translator) Close() error {
	if err := w.file.Flush(); err != nil {
		return err
	}

	return nil
}
