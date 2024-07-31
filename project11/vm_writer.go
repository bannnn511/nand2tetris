package main

import (
	"fmt"
	"strings"
)

type VmWriter struct {
	out         strings.Builder
	indentation int
	labelCount  int
}

func NewVmWriter() *VmWriter {
	return &VmWriter{
		indentation: 0,
		labelCount:  0,
	}
}

func (w *VmWriter) WriteFunction(
	routine string,
	name string,
	nVars uint,
	params uint,
) {
	function := fmt.Sprintf("function %v %d\n", name, nVars)
	w.out.WriteString(function)

	if routine == "constructor" {
		pushArg := fmt.Sprintf("push constant %d\n", params)
		w.WriteWithIndentation(pushArg)

		w.WriteWithIndentation("call Memory.alloc 1\n")
		w.WriteWithIndentation("pop pointer 0\n")
	}
}

func (w *VmWriter) WriteDo(
	fName string,
	nVars int,
) {
	do := fmt.Sprintf("call %v %d\n", fName, nVars)

	w.WriteIndentation(4)
	w.out.WriteString(do)

	w.WriteIndentation(4)
	w.out.WriteString("pop temp 0\n")
}

func (w *VmWriter) WriteDoWithReturn(
	fName string,
	nVars int,
	variableKind string,
	index uint,
) {
	w.WriteIndentation(4)
	do := fmt.Sprintf("call %v %d\n", fName, nVars)
	w.out.WriteString(do)

	w.WriteIndentation(4)
	popTo := fmt.Sprintf("pop %v %d\n", variableKind, index)
	w.out.WriteString(popTo)
}

func (w *VmWriter) WritePushVariableToStack(segment VariableKind, idx uint) {
	w.WriteIndentation(4)
	push := fmt.Sprintf("push %v %d\n", segment.String(), idx)
	w.out.WriteString(push)
}

func (w *VmWriter) WritePopVariable(segment VariableKind, idx uint) {
	w.WriteIndentation(4)
	pop := fmt.Sprintf("pop %v %d\n", segment.String(), idx)
	w.out.WriteString(pop)
}

func (w *VmWriter) WriteIf(label int) {
	w.WriteIndentation(4)
	w.Write("not\n")
	w.WriteIndentation(4)
	w.Write(fmt.Sprintf("if-goto L%d\n", label))
}

func (w *VmWriter) WriteGoto(idx int) {
	w.WriteIndentation(4)
	w.Write(fmt.Sprintf("goto L%d\n", idx))
}

func (w *VmWriter) WriteReturn() {
	w.WriteIndentation(4)
	w.out.WriteString("return\n")
}

func (w *VmWriter) WriteOp(op string) {
	w.WriteIndentation(4)
	switch op {
	case "*":
		w.Write("call Math.multiply 2\n")
	case "+":
		w.Write("add\n")
	case "-":
		w.Write("sub\n")
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
	w.labelCount++
}

func (w *VmWriter) GetLabelIdx() int {
	defer func() {
		w.labelCount += 1
	}()

	return w.labelCount
}

func (w *VmWriter) GetLabel() string {
	return fmt.Sprintf("L%d", w.labelCount)
}

func (w *VmWriter) WriteFalse() {
	w.WriteIndentation(4)
	w.Write("push constant 0\n")
}

func (w *VmWriter) WriteTrue() {
	w.WriteIndentation(4)
	w.Write("push constant 1\n")

	w.WriteIndentation(4)
	w.Write("neg\n")
}

func (w *VmWriter) Write(str string) {
	w.out.WriteString(str)
}

func (w *VmWriter) Out() string {
	return w.out.String()
}

func (w *VmWriter) WriteWithIndentation(str string) {
	w.WriteIndentation(4)
	w.Write(str)
}

func (w *VmWriter) WriteIndentation(indent int) {
	for i := 0; i < indent; i++ {
		w.Write(" ")
	}
}
