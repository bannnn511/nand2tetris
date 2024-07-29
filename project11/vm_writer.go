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
	fType string,
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

func (w *VmWriter) WriteReturn() {
	w.writeIndentation()
	w.out.WriteString("push constant 0\n")

	w.writeIndentation()
	w.out.WriteString("return")
}

func (w *VmWriter) WriteOp(op string) {
	switch op {
	case "*":
		w.write("call Math.multiply 2\n")
	case "+":
		w.write("add\n")
	default:
		w.write(op + "\n")
	}
}

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
