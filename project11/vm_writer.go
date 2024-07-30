package main

import (
	"fmt"
	"strings"
)

type VmWriter struct {
	out         strings.Builder
	indentation int
}

func NewVmWriter() *VmWriter {
	return &VmWriter{}
}

func (w *VmWriter) writeFunction(
	routine string,
	name string,
	nVars uint,
) {
	w.writeIndentation()

	tok := Lookup2(routine)
	function := fmt.Sprintf("%v %v %d\n", tok, name, nVars)

	w.out.WriteString(function)
}

func (w *VmWriter) WriteDo(
	fName string,
	nVars int,
) {
	w.writeIndentation()

	do := fmt.Sprintf("call %v %d\n", fName, nVars)
	w.out.WriteString(do)

	w.writeIndentation()

	w.out.WriteString("pop temp 0\n")
}

func (w *VmWriter) WriteDoWithReturn(
	fName string,
	nVars int,
	variableKind string,
	index uint32,
) {
	w.writeIndentation()

	do := fmt.Sprintf("call %v %d\n", fName, nVars)
	w.out.WriteString(do)

	w.writeIndentation()

	popTo := fmt.Sprintf("pop %v %d\n", variableKind, index)
	w.out.WriteString(popTo)
}

func (w *VmWriter) WritePushVariableToStack(segment VariableKind, idx uint32) {
	w.writeIndentation()
	w.out.WriteString(fmt.Sprintf("push %v %d\n", segment.String(), idx))
}

func (w *VmWriter) WriteReturn() {
	w.writeIndentation()
	w.out.WriteString("push constant 0\n")

	w.writeIndentation()
	w.out.WriteString("return\n")
}

func (w *VmWriter) WriteOp(op string) {
	w.writeIndentation()
	switch op {
	case "*":
		w.write("call Math.multiply 2\n")
	case "+":
		w.write("add\n")
	case "-":
		w.write("neg\n")
	case "<":
	default:
		w.write(op + "\n")
	}
}

// func (w *VmWriter) WriteOp(op string) {
// 	var symbol string
// 	switch op {
// 	case "<":
// 		symbol = "&lt;"
// 	case ">":
// 		symbol = "&gt;"
// 	case "&":
// 		symbol = "&amp;"
// 	case "-":
// 		symbol = "neg"
// 	default:
// 		symbol = op
// 	}

// 	w.out.WriteString(symbol)
// }

func (w *VmWriter) write(str string) {
	w.out.WriteString(str)
}

func (w *VmWriter) Out() string {
	return w.out.String()
}

func (w *VmWriter) IncrIndent() {
	w.indentation += 4
}

func (w *VmWriter) DecrIndent() {
	w.indentation -= 4
}

func (w *VmWriter) writeWithIndentation(str string) {
	w.writeIndentation()
	w.write(str)
}

func (w *VmWriter) writeIndentation() {
	for i := 0; i < w.indentation; i++ {
		w.write(" ")
	}
}
