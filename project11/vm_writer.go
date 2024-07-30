package main

import (
	"fmt"
	"strings"
)

var labelCount = 0

type VmWriter struct {
	out         strings.Builder
	indentation int
}

func NewVmWriter() *VmWriter {
	return &VmWriter{}
}

func (w *VmWriter) WriteFunction(
	routine string,
	name string,
	nVars uint,
) {
	w.WriteIndentation()

	tok := Lookup2(routine)
	function := fmt.Sprintf("%v %v %d\n", tok, name, nVars)

	w.out.WriteString(function)
}

func (w *VmWriter) WriteDo(
	fName string,
	nVars int,
) {
	w.WriteIndentation()

	do := fmt.Sprintf("call %v %d\n", fName, nVars)
	w.out.WriteString(do)

	w.WriteIndentation()

	w.out.WriteString("pop temp 0\n")
}

func (w *VmWriter) WriteDoWithReturn(
	fName string,
	nVars int,
	variableKind string,
	index uint,
) {
	w.WriteIndentation()

	do := fmt.Sprintf("call %v %d\n", fName, nVars)
	w.out.WriteString(do)

	w.WriteIndentation()

	popTo := fmt.Sprintf("pop %v %d\n", variableKind, index)
	w.out.WriteString(popTo)
}

func (w *VmWriter) WritePushVariableToStack(segment VariableKind, idx uint) {
	w.WriteIndentation()
	push := fmt.Sprintf("push %v %d\n", segment.String(), idx)
	w.out.WriteString(push)
}

func (w *VmWriter) WritePopVariable(segment VariableKind, idx uint) {
	w.WriteIndentation()
	pop := fmt.Sprintf("pop %v %d\n", segment.String(), idx)
	w.out.WriteString(pop)
}

func (w *VmWriter) WriteIf(label int) {
	w.WriteIndentation()
	w.Write("not\n")
	w.WriteIndentation()
	w.Write(fmt.Sprintf("if-goto L%d\n", label))
}

func (w *VmWriter) WriteGoto(idx int) {
	w.WriteIndentation()
	w.Write(fmt.Sprintf("goto L%d\n", idx))
}

func (w *VmWriter) WriteReturn() {
	w.WriteIndentation()
	w.out.WriteString("push constant 0\n")

	w.WriteIndentation()
	w.out.WriteString("return\n")
}

func (w *VmWriter) WriteOp(op string) {
	w.WriteIndentation()
	switch op {
	case "*":
		w.Write("call Math.multiply 2\n")
	case "+":
		w.Write("add\n")
	case "-":
		w.Write("neg\n")
	case "<":
		w.Write("lt\n")
	case ">":
		w.Write("gt\n")
	case "~":
		w.Write("not\n")
	case "=":
		w.Write("eq\n")
	case "&":
		w.Write("and\n")
	default:
		w.Write(op + "\n")
	}
}

func (w *VmWriter) WriteLabel(label int) {
	w.out.WriteString(fmt.Sprintf("label L%d\n", label))
}

func (w *VmWriter) IncrLabel() {
	labelCount++
}

func (w *VmWriter) GetLabelIdx() int {
	return labelCount
}

func (w *VmWriter) GetLabel() string {
	return fmt.Sprintf("L%d", labelCount)
}

func (w *VmWriter) WriteFalse() {
	w.WriteIndentation()
	w.Write("push constant 1\n")
}

func (w *VmWriter) WriteTrue() {
	w.WriteIndentation()
	w.Write("push constant 1\n")
	w.WriteIndentation()
	w.Write("neg\n")
}

func (w *VmWriter) Write(str string) {
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
	w.WriteIndentation()
	w.Write(str)
}

func (w *VmWriter) WriteIndentation() {
	for i := 0; i < w.indentation; i++ {
		w.Write(" ")
	}
}
