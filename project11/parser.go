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
	elements []Ast   // list of lexical elements
	tok      Token   // current token
	lit      string  // current value
	scanner  Scanner // token scanner

	classSB   *SymbolTable // class variable
	routineSB *SymbolTable // subroutine variable
	vmWriter  VmWriter

	// compile state
	className    string
	kind         VariableKind
	variableName string
}

func (p *Parser) Init(filename string, src []byte) {
	p.scanner.Init(src)
	p.tok = START
	p.classSB = NewSymbolTable()
	p.routineSB = NewSymbolTable()
	p.vmWriter = *NewVmWriter()
	p.elements = make([]Ast, 0, len(src))
	p.append(Ast{START, ""})
}

func (p *Parser) ParseFile() {
	p.next()

	//  class
	p.next()

	// Main
	p.className = p.lit
	p.next()

	// {
	p.next()

	for p.lit == "static" || p.lit == "field" {
		p.compileClassVarDec()
	}

	for p.lit == "constructor" ||
		p.lit == "function" ||
		p.lit == "method" {
		p.compileSubroutine()
	}

	// }
	// p.writeTemplate()

}

func (p *Parser) compileSubroutine() {
	// construction, function, method
	_, routineLit := p.getState()

	p.next()
	_, typeLit := p.getState()

	p.next()
	_, fName := p.getState()
	fName = p.className + "." + fName

	p.vmWriter.writeFunction(routineLit, fName, 0, typeLit)

	p.next()
	p.writeTemplate() // (
	p.next()
	p.compileParameterList() // parameter
	p.writeTemplate()        // )
	p.next()

	p.writeTemplate()

	p.next()
	for p.lit == "var" {
		p.compileVarDec()
	}

	p.compileStatements()

	p.writeTemplate() // symbol

	p.next()
}

func (p *Parser) compileParameterList() {

	for p.tok != SYMBOL {
		p.writeTemplate()
		p.next()

		p.writeTemplate()
		p.next()

		if p.lit == "," {
			p.writeTemplate()
			p.next()
		}
	}

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
	p.vmWriter.IncrIndent()
	// return
	p.vmWriter.WriteReturn()

	p.next()
	if p.tok != SYMBOL && p.lit != ";" {
		p.compileExpressions()
	}

	p.writeTemplate()
	p.vmWriter.DecrIndent()

	p.next()
}

// <doStatement>
func (p *Parser) compileDo() {
	p.vmWriter.IncrIndent()

	// state: do

	p.next()
	// state: identifier
	fName := p.lit

	p.next()
	if p.lit == "." {
		// state: '.'
		fName += p.lit
		p.next()

		// state: identifier
		fName += p.lit
		p.next()
	}

	// state: (
	p.next()
	nVars := p.compileExpressionList()
	// state: )

	p.next()

	// state: ';'
	p.vmWriter.WriteDo(fName, nVars)

	p.vmWriter.DecrIndent()
	p.next()
}

// 'if' '(' expression ')' '{' statements '}' ( 'else' '{' statements '}' )?
func (p *Parser) compileIf() {

	p.writeTemplate() // if

	p.next()
	p.writeTemplate() // (

	p.next()
	p.compileExpressions()

	p.writeTemplate() // )

	p.next()
	p.writeTemplate() // {

	p.next()
	p.compileStatements()
	p.writeTemplate() // }

	p.next()
	if p.tok == KEYWORD && p.lit == "else" {
		p.writeTemplate() // else

		p.next()
		p.writeTemplate() // {
		p.next()

		// if else body is empty -> dont call next
		prev := p.lit
		p.compileStatements()
		if prev != "}" && p.lit != "}" {
			p.next()
		}

		p.writeTemplate() // }
		p.next()
	}

}

// <letStatement>
func (p *Parser) compileLet() {
	p.vmWriter.IncrIndent()
	// state: let

	p.next()
	// state: varName
	p.variableName = p.lit

	p.next()
	if p.lit == "[" {
		p.writeTemplate() // [
		p.next()
		p.compileExpressions()
		p.writeTemplate() // ]
		p.next()
	}

	// state: =
	p.next()
	p.compileExpressions2()
	// state: ;

	p.vmWriter.DecrIndent()
	p.next()
}

// 'while' '(' expression ')' '{' statements '}'
func (p *Parser) compileWhile() {

	p.writeTemplate() // while

	p.next()
	p.writeTemplate() // (

	p.next()
	p.compileExpressions()

	p.writeTemplate() // )

	p.next()
	p.writeTemplate() // {

	p.next()
	p.compileStatements()

	p.writeTemplate() // }

	p.next()
}

func (p *Parser) compileExpressions() {

	p.CompileTerm()
	for p.tok == SYMBOL && IsOp(p.lit) {
		p.writeTemplate()
		p.next()
		p.CompileTerm()
	}

}

func (p *Parser) CompileTerm() {

	switch p.tok {
	case INT:
		p.writeTemplate()
		p.next()
	case CHAR:
		p.writeTemplate()
		p.next()
	case KEYWORD:
		p.writeTemplate()
		p.next()
	case IDENT:
		p.writeTemplate()
		p.next()
		if p.lit == "[" {
			p.writeTemplate()
			p.next()
			p.compileExpressions()
			p.writeTemplate()
			p.next()
		} else if p.lit == "." {
			p.writeTemplate() // symbol
			p.next()
			p.writeTemplate() // identifier
			p.next()
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressionList()
			p.writeTemplate() // symbol
			p.next()
		} else if p.lit == "(" {
			p.writeTemplate() // symbol
			p.compileExpressionList()
			p.next()
			p.writeTemplate() // symbol
			p.next()
		}
		// else if p.lit == "~" || p.lit == "-" {
		//	p.writeTemplate() // symbol
		//	p.next()
		//	p.CompileTerm()
		// }
	case SYMBOL:
		if p.lit == "(" {
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressions()
			p.writeTemplate() // symbol
			p.next()
		} else if p.lit == "~" || p.lit == "-" {
			p.writeTemplate() // symbol
			p.next()
			p.CompileTerm()
		}
	default:
	}

}

// term(op term)*
func (p *Parser) compileExpressions2() {
	p.compileTerm2() // term
	for p.tok == SYMBOL && IsOp(p.lit) {
		// op
		op := p.lit
		p.next()

		// term
		p.compileTerm2()

		// op
		p.vmWriter.writeIndentation()
		p.vmWriter.WriteOp(op)
	}

}

func (p *Parser) compileTerm2() {
	switch p.tok {
	case INT:
		p.vmWriter.writeIndentation()
		p.vmWriter.write(fmt.Sprintf("push constant %v\n", p.lit))
		p.next()
	case CHAR:
		p.writeTemplate()
		p.next()
	case KEYWORD:
		p.writeTemplate()
		p.next()
	case IDENT:
		// state: identifier
		identifier := p.lit
		p.next()
		if p.lit == "[" {
			p.writeTemplate()
			p.next()
			p.compileExpressions2()
			p.writeTemplate()
			p.next()
		} else if p.lit == "." {
			// state: '.'
			fName := identifier + p.lit
			p.next()

			// state: identifier
			fName += p.lit
			p.next()

			// state: symbol
			p.next()

			nVars := p.compileExpressionList()
			if nVars == 0 {
				p.vmWriter.WriteDo(fName, nVars)
			} else {
				vKind := p.routineSB.KindOf(p.variableName)
				index := p.routineSB.IndexOf(p.variableName)
				p.vmWriter.WriteDoWithReturn(fName, nVars, vKind.String(), index)
			}

			// state: symbol
			p.next()
		} else if p.lit == "(" {
			p.writeTemplate() // symbol
			p.compileExpressionList()
			p.next()
			p.writeTemplate() // symbol
			p.next()
		}
		// else if p.lit == "~" || p.lit == "-" {
		//	p.writeTemplate() // symbol
		//	p.next()
		//	p.CompileTerm()
		// }
	case SYMBOL:
		if p.lit == "(" {
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressions2()
			p.writeTemplate() // symbol
			p.next()
		} else if p.lit == "~" || p.lit == "-" {
			p.writeTemplate() // symbol
			p.next()
			p.compileTerm2()
		}
	default:
	}
}

// <expressionList>
func (p *Parser) compileExpressionList() int {
	count := 0

	if p.tok != SYMBOL && p.lit != ")" {
		count++
		p.compileExpressions2()
		for p.tok == SYMBOL && p.lit == "," {
			// symbol ','
			p.next()
			p.compileExpressions2()
			count++
		}
	}

	// if after '(' is a '(' -> new expression
	if p.lit == "(" {
		p.compileExpressions2()
		for p.tok == SYMBOL && p.lit == "," {
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressions2()
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

	// subroutine symbol table must be reconstruct for each subroutine
	p.routineSB = NewSymbolTable()

	p.kind = Var
	p.writeTemplate()
	p.next()
	p.compileTypeAndVarName(Subroutine)

}

func (p *Parser) compileTypeAndVarName(scope VariableScope) {
	vType := p.lit
	// variable type
	// p.writeTemplate()
	p.next()

	name := p.lit
	// variable identifier
	// p.writeTemplate()
	p.next()

	p.defineVariable(scope, name, vType, p.kind)

	for p.lit == "," {
		// p.writeTemplate() // ',''
		p.next()

		name := p.lit
		p.defineVariable(scope, name, vType, p.kind)
		// p.writeTemplate() // variable identifier
		p.next()
	}

	// p.writeTemplate()
	p.next()
}

func (p *Parser) compileClassVarDec() {
	// p.writeWithIndentation("<classVarDec>\r\n")

	p.kind = WhichKind(p.lit)
	// var
	// p.writeTemplate()
	p.next()

	p.compileTypeAndVarName(Class)

	//
	// p.writeWithIndentation("</classVarDec>\r\n")
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
	p.lit = lit
	p.elements = append(p.elements, Ast{tok, lit})
}

const template = "<%v> %v </%v>\r\n"

// writeTemplate writes KEYWORD, IDENT and SYMBOL token
func (p *Parser) writeTemplate() {
	if p.tok == EOF {
		return
	}

	if p.tok == SYMBOL {
		p.writeSymbol()
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

func (p *Parser) writeSymbol() {
	// var symbol string
	// switch p.lit {
	// case "<":
	// 	symbol = "&lt;"
	// case ">":
	// 	symbol = "&gt;"
	// case "&":
	// 	symbol = "&amp;"
	// default:
	// 	symbol = p.lit
	// }

	// p.write("<symbol> " + symbol + " </symbol>\r\n")
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
