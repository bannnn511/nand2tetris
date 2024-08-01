package main

import (
	"fmt"
	"strings"
)

type Ast struct {
	tok Token
	lit string
}

type Parser struct {
	elements []Ast  // list of lexical elements
	tok      Token  // current token
	lit      string // current value
	prev     string
	scanner  Scanner // token scanner

	classSB   *SymbolTable // class variable
	routineSB *SymbolTable // subroutine variable
	vmWriter  VmWriter

	// compile state
	className    string
	kind         VariableKind
	variableName string
}

func (p *Parser) Init(src []byte) {
	p.scanner.Init(src)
	p.tok = START
	p.classSB = NewSymbolTable()
	p.routineSB = NewSymbolTable()
	p.vmWriter = *NewVmWriter()
	p.elements = make([]Ast, 0, len(src))
	p.append(Ast{START, ""})
}

func (p *Parser) ParseFile() {
	p.next() //  class
	p.next() // Main
	p.className = p.lit
	p.next() // {
	p.next()

	for p.lit == "static" || p.lit == "field" {
		p.compileClassVarDec()
	}

	for p.lit == "constructor" ||
		p.lit == "function" ||
		p.lit == "method" {
		if p.lit == "function" || p.lit == "method" {
			p.routineSB = NewSymbolTable()
		}
		if p.lit == "method" {
			p.defineVariable(
				Subroutine,
				"this",
				p.className,
				Arg,
			)
		}
		p.compileSubroutine()
	}
}

func (p *Parser) compileSubroutine() {
	// state: subroutine type - construction, function, method
	_, routineLit := p.getState()

	p.next() // state: function type (void)
	p.next() // state: function name

	_, fName := p.getState()
	fName = p.className + "." + fName

	p.next() // state: '('
	p.next()
	paramCount := p.compileParameterList() // parameters
	// state: ')'
	p.next() // '{'

	p.next()
	for p.lit == "var" {
		p.compileVarDec()
	}

	count := p.routineSB.VarCount(Var)
	if routineLit == "constructor" {
		fieldCount := p.classSB.VarCount(Field)
		staticCount := p.classSB.VarCount(Static)
		paramCount = int(fieldCount + staticCount)
	} else if routineLit == "function" {
		varCount := p.routineSB.VarCount(Var)
		paramCount += int(varCount)
	}
	p.vmWriter.WriteFunction(
		routineLit,
		fName,
		count,
		uint(paramCount),
	)

	// if subroutine is method
	// the caller must push the refernce to the opbbject
	if routineLit == "method" {
		p.shouldPushToVariable("this")
		p.shouldPopToVariable("this")
	}

	p.compileStatements()

	p.writeTemplate() // symbol

	p.next()
}

func (p *Parser) compileParameterList() int {
	count := 0
	for p.tok != SYMBOL {
		// p.writeTemplate()
		p.next()

		p.defineVariable(Subroutine, p.lit, p.prev, Arg)
		p.writeTemplate()
		p.next()
		if p.lit == "," {
			p.writeTemplate()
			p.next()
		}
		count++
	}

	return count
}

// <statements>
func (p *Parser) compileStatements() {
	for p.tok == KEYWORD && p.lit != "}" {
		switch p.lit {
		case "let":
			p.compileLet()
		case "while":
			p.compileWhile()
		case "if":
			p.compileIf()
		case "do":
			p.compileDo()
		case "return":
			p.compileReturn()
		}
	}
}

// <returnStatement>
func (p *Parser) compileReturn() {
	// return
	p.next()
	if p.tok != SYMBOL && p.lit != ";" {
		p.compileExpressions()
	} else {
		p.vmWriter.WriteIndentation(4)
		p.vmWriter.out.WriteString("push constant 0\n")
	}
	p.vmWriter.WriteReturn()
	p.next()
}

// <doStatement>
func (p *Parser) compileDo() {
	// state: do

	p.next()
	// state: identifier
	variableName := p.lit
	fName, vType := p.routineSB.TypeOf(variableName)

	// ILLEGAL should be language standard library
	if vType == ILLEGAL {
		fName, vType = p.classSB.TypeOf(variableName)
		if vType == ILLEGAL {
			fName = variableName
		}
	}

	p.next() // '.' or '('

	// method call will use this as argument -> nVars add 1 as offset
	methodArg := 0
	if p.lit == "." {
		fName += p.lit
		p.next() // identifier
		fName += p.lit
		p.next()
	} else if vType == ILLEGAL {
		// method call
		p.vmWriter.WriteFormat("push pointer 0\n")
		fName = p.className + "." + fName
		methodArg++
	}

	p.next() // expressions
	nVars := p.compileExpressionList()
	p.next() // state: ';'

	// if vType is USR -> method function call
	if vType == USR {
		p.shouldPushToVariable(variableName)
		nVars++
	}

	p.vmWriter.WriteDo(fName, nVars+methodArg)
	p.next()
}

// 'if' '(' expression ')' '{' statements '}' ( 'else' '{' statements '}' )?
func (p *Parser) compileIf() {
	// state:if
	p.next() // '('
	p.next() // exp
	p.compileExpressions()

	l1 := p.vmWriter.GetLabelIdx()
	l2 := p.vmWriter.GetLabelIdx()

	p.vmWriter.WriteIf(l2)

	p.next() // '{'
	p.next()
	p.compileStatements()

	p.vmWriter.WriteGoto(l1)
	p.vmWriter.WriteLabel(l2)

	p.next() // else
	if p.tok == KEYWORD && p.lit == "else" {
		// state: else
		p.next() // '{'
		p.next()

		// if else body is empty -> dont call next
		prev := p.lit
		p.compileStatements()
		if prev != "}" && p.lit != "}" {
			p.next()
		}

		// '}'
		p.next()
	}
	p.vmWriter.WriteLabel(l1)
}

// <letStatement>
// handle array for case arr[expression1] = expression2
func (p *Parser) compileLet() {
	// state: let

	p.next() // state: varName
	p.variableName = p.lit

	p.next() // '=' or '['
	isArray := false
	if p.lit == "[" {
		isArray = true
		p.next()                               // array index
		p.compileExpressions()                 // compute expression 1
		p.shouldPushToVariable(p.variableName) // push arr
		p.variableName = ""                    // set empty so that pop dont need for this case
		p.vmWriter.WriteFormat("add\n")        // top stack value = arr[expression1]
		// state: ']'
		p.next()
	}

	p.next() // state:
	p.compileExpressions()
	p.shouldPopToVariable(p.variableName)
	if isArray {
		p.vmWriter.WritePopArrayExpression()
	}
	p.variableName = ""

	p.next()
}

// 'while' '(' expression ')' '{' statements '}'
func (p *Parser) compileWhile() {
	// state: while

	l1 := p.vmWriter.GetLabelIdx()
	p.vmWriter.WriteLabel(l1)
	l2 := p.vmWriter.GetLabelIdx()

	p.next() // state: '('
	p.next() // state: ')'
	p.compileExpressions()

	p.vmWriter.WriteIf(l2)
	p.next() // state: {
	p.next()
	p.compileStatements()

	p.vmWriter.WriteGoto(l1)
	p.vmWriter.WriteLabel(l2)

	// state: '}'
	p.next()
}

// term(op term)*
func (p *Parser) compileExpressions() {
	p.compileTerm2() // term
	for p.tok == SYMBOL && IsOp(p.lit) {
		// op
		op := p.lit
		p.next() // term
		p.compileTerm2()
		p.vmWriter.WriteOp(op)
	}
}

func (p *Parser) shouldPushToVariable(name string) {
	if name == "" {
		return
	}
	if p.classSB.IsExists(name) {
		kind, idx := p.classSB.GetSegment(name)
		p.vmWriter.WritePushVariableToStack(kind, idx)
	} else if p.routineSB.IsExists(name) {
		kind, idx := p.routineSB.GetSegment(name)
		p.vmWriter.WritePushVariableToStack(kind, idx)
	}
}

func (p *Parser) shouldPopToVariable(name string) {
	if name == "" {
		return
	}

	if p.classSB.IsExists(name) {
		kind, idx := p.classSB.GetSegment(name)
		p.vmWriter.WritePopVariable(kind, idx)
	} else if p.routineSB.IsExists(name) {
		if name == "this" {
			p.vmWriter.WriteFormat("pop pointer 0\n")
		} else {
			kind, idx := p.routineSB.GetSegment(name)
			p.vmWriter.WritePopVariable(kind, idx)
		}
	}
}

func (p *Parser) compileTerm2() {
	switch p.tok {
	case INT:
		p.vmWriter.WriteFormat(fmt.Sprintf("push constant %v\n", p.lit))
		p.next()
	case CHAR:
		p.vmWriter.WriteString(p.lit)
		p.next()
	case KEYWORD:
		switch p.lit {
		case "true":
			p.vmWriter.WriteTrue()
			p.shouldPopToVariable(p.variableName)
			p.variableName = ""
		case "false":
			p.vmWriter.WriteFalse()
			p.shouldPopToVariable(p.variableName)
			p.variableName = ""
		case "this":
			// return this
			if p.prev == "return" {
				p.vmWriter.WriteFormat("push pointer 0\n")
			}

			// method call with 'this'
			if p.prev == "(" {
				p.vmWriter.WriteFormat("push pointer 0\n")
			} else {
				p.shouldPopToVariable(p.variableName)
			}
		default:
			println("implement other keyword case", p.lit)
		}
		p.next()
	case IDENT:
		// state: identifier
		identifier := p.lit
		p.shouldPushToVariable(identifier)
		p.lit = ""

		p.next()
		if p.lit == "[" {
			p.next() // array expression
			p.compileExpressions()
			p.vmWriter.WriteFormat("add\n")
			p.vmWriter.WriteFormat("pop pointer 1\n")
			p.vmWriter.WriteFormat("push that 0\n")
			// state: ']'
			p.next()
		} else if p.lit == "." {
			// state: '.'
			fName := identifier + p.lit
			p.next() // state: identifier
			fName += p.lit
			p.next() // state: symbol
			p.next()

			nVars := p.compileExpressionList()
			vKind := p.routineSB.KindOf(p.variableName)
			if vKind == Undefined {
				vKind = p.classSB.KindOf(p.variableName)
			}
			index := p.routineSB.IndexOf(p.variableName)

			p.vmWriter.WriteDoWithReturn(fName, nVars, vKind.String(), index)
			p.variableName = ""

			// state: symbol
			p.next()
		} else if p.lit == "(" {
			p.writeTemplate() // symbol
			p.compileExpressionList()
			p.next() // symbol
			p.next()
		} else if p.lit == "~" || p.lit == "-" {
			op := p.lit
			p.next()
			p.compileTerm2()
			p.vmWriter.WriteOp(op)
		}
	case SYMBOL:
		if p.lit == "(" {
			p.writeTemplate() // symbol
			p.next()          // symbol
			p.compileExpressions()
			p.next()
		} else if p.lit == "~" || p.lit == "-" {
			op := p.lit
			p.next()
			p.compileTerm2()
			if op == "~" {
				p.vmWriter.WriteFormat("not\n")
			} else {
				p.vmWriter.WriteFormat("neg\n")
			}
		}
	default:
	}
}

// <expressionList>
func (p *Parser) compileExpressionList() int {
	count := 0
	if p.tok != SYMBOL && p.lit != ")" {
		count++
		p.compileExpressions()
		for p.tok == SYMBOL && p.lit == "," {
			// symbol ','
			p.next()
			p.compileExpressions()
			count++
		}
	}

	// if after '(' is a '(' -> new expression
	if p.lit == "(" {
		p.compileExpressions()
		count++
		for p.tok == SYMBOL && p.lit == "," {
			// symbol
			p.next()
			p.compileExpressions()
			count++
		}
	}

	return count
}

type VariableScope int

const (
	Class VariableScope = iota
	Subroutine
)

func (p *Parser) compileVarDec() {
	p.kind = Var
	// p.writeTemplate()
	p.next()
	p.compileTypeAndVarName(Subroutine)
}

func (p *Parser) compileTypeAndVarName(scope VariableScope) {
	vType := p.lit
	// state: variable type
	p.next()

	name := p.lit
	// state: variable identifier
	p.next()
	p.defineVariable(scope, name, vType, p.kind)
	for p.lit == "," {
		// state ','
		p.next()

		name := p.lit
		p.defineVariable(scope, name, vType, p.kind)

		// state: variable identifier
		p.next()
	}

	// p.writeTemplate()
	p.next()
}

func (p *Parser) compileClassVarDec() {
	p.kind = WhichKind(p.lit)
	// var
	// p.writeTemplate()
	p.next()
	p.compileTypeAndVarName(Class)
}

func (p *Parser) append(ele Ast) {
	p.elements = append(p.elements, ele)
}

func (p *Parser) next() {
	tok, lit := p.scanner.Scan()
	if tok == COMMENT {
		p.next()
		return
	}
	p.tok = tok
	p.prev = p.lit
	p.lit = lit
	p.elements = append(p.elements, Ast{tok, lit})
}

// writeTemplate writes KEYWORD, IDENT and SYMBOL token
func (p *Parser) writeTemplate() {
	if p.tok == EOF {
		return
	}

	if p.tok == SYMBOL {
		// p.writeSymbol()
		return
	}

	if p.tok == INT {
		// p.write(fmt.Sprintf(template, "integerConstant", p.lit, "integerConstant"))
		return
	}

	if p.tok == CHAR {
		// p.write(fmt.Sprintf(template, "stringConstant", p.lit, "stringConstant"))
		return
	}

	// p.write(fmt.Sprintf(template, p.tok, p.lit, p.tok))
}

func (p *Parser) VmOut() string {
	return strings.TrimSuffix(p.vmWriter.Out(), "\n")
}

func (p *Parser) defineVariable(
	scope VariableScope,
	name string,
	vType string,
	kind VariableKind,
) {
	switch scope {
	case Class:
		p.classSB.Define(vType, name, kind)
	case Subroutine:
		p.routineSB.Define(vType, name, kind)
	}
}

func (p *Parser) getState() (Token, string) {
	return p.tok, p.lit
}
