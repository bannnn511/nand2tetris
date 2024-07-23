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
	elements    []Ast // list of lexical elements
	tok         Token // current token
	lit         string
	fileName    string
	scanner     Scanner
	indentation int
	out         strings.Builder
}

func (p *Parser) Init(filename string, src []byte) {
	p.scanner.Init(src)
	p.fileName = filename
	p.tok = START
	p.elements = make([]Ast, 0, len(src))
	p.append(Ast{START, ""})
}

func (p *Parser) ParseFile() {
	p.next()
	p.write("<class>\r\n")
	p.indentation += 2

	p.writeTemplate() //  class
	p.next()
	p.writeTemplate() // Main
	p.next()
	p.writeTemplate() // {
	p.next()

	for p.lit == "static" || p.lit == "field" {
		p.compileClassVarDec()
	}
	for p.lit == "constructor" ||
		p.lit == "function" ||
		p.lit == "method" {
		p.compileSubroutine()
	}

	p.writeTemplate()
	p.indentation -= 2
	p.write("</class>\r\n")
}

func (p *Parser) compileSubroutine() {
	p.writeWithIndentation("<subroutineDec>\r\n")
	p.indentation += 2
	p.writeTemplate() // kw

	p.next()
	p.writeTemplate()

	p.next()
	p.writeTemplate() // identifier

	p.next()
	p.writeTemplate() // (
	p.next()
	p.compileParameterList() // parameter
	p.writeTemplate()        // )
	p.next()

	p.writeWithIndentation("<subroutineBody>\r\n")
	p.indentation += 2
	p.writeTemplate()

	p.next()
	for p.lit == "var" {
		p.compileVarDec()
	}

	p.compileStatements()

	p.writeTemplate() // symbol
	p.indentation -= 2
	p.writeWithIndentation("</subroutineBody>\r\n")
	p.indentation -= 2
	p.writeWithIndentation("</subroutineDec>\r\n")
	p.next()
}

func (p *Parser) compileStatements() {
	p.writeWithIndentation("<statements>\r\n")
	p.indentation += 2

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

	p.indentation -= 2
	p.writeWithIndentation("</statements>\r\n")
}

func (p *Parser) compileReturn() {
	p.writeWithIndentation("<returnStatement>\r\n")
	p.indentation += 2

	p.writeTemplate()

	p.next()
	if p.tok != SYMBOL && p.lit != ";" {
		p.compileExpressions()
	}

	p.writeTemplate()

	p.indentation -= 2
	p.writeWithIndentation("</returnStatement>\r\n")
	p.next()
}

// 'do' subroutineCall ';'
func (p *Parser) compileDo() {
	p.writeWithIndentation("<doStatement>\r\n")
	p.indentation += 2

	p.writeTemplate() // do
	p.next()

	p.writeTemplate() // identifier
	p.next()
	if p.lit == "." {
		p.writeTemplate() // symbol
		p.next()
		p.writeTemplate() // identifier
		p.next()
	}

	p.writeTemplate() // symbol
	p.next()

	p.compileExpressionList()

	p.writeTemplate() // symbol

	p.next()
	p.writeTemplate() // symbol

	p.indentation -= 2
	p.writeWithIndentation("</doStatement>\r\n")
	p.next()
}

// 'if' '(' expression ')' '{' statements '}' ( 'else' '{' statements '}' )?
func (p *Parser) compileIf() {
	p.writeWithIndentation("<ifStatement>\r\n")
	p.indentation += 2

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

	p.indentation -= 2
	p.writeWithIndentation("</ifStatement>\r\n")
}

// 'let' varName ('[' expression ']')? '=' expression ';'
func (p *Parser) compileLet() {
	p.writeWithIndentation("<letStatement>\r\n")
	p.indentation += 2

	p.writeTemplate() // let

	p.next()
	p.writeTemplate() // varName

	p.next()
	if p.lit == "[" {
		p.writeTemplate() // [
		p.next()
		p.compileExpressions()
		p.writeTemplate() // ]
		p.next()
	}

	p.writeTemplate() // =
	p.next()

	p.compileExpressions()
	p.writeTemplate() // ;

	p.indentation -= 2
	p.writeWithIndentation("</letStatement>\r\n")
	p.next()
}

// 'while' '(' expression ')' '{' statements '}'
func (p *Parser) compileWhile() {
	p.writeWithIndentation("<whileStatement>\r\n")
	p.indentation += 2

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

	p.indentation -= 2
	p.writeWithIndentation("</whileStatement>\r\n")
	p.next()
}

func (p *Parser) compileExpressions() {
	p.writeWithIndentation("<expression>\r\n")
	p.indentation += 2

	p.CompileTerm()
	for p.tok == SYMBOL && IsOp(p.lit) {
		p.writeTemplate()
		p.next()
		p.CompileTerm()
	}

	p.indentation -= 2
	p.writeWithIndentation("</expression>\r\n")
}

func (p *Parser) CompileTerm() {
	p.writeWithIndentation("<term>\r\n")
	p.indentation += 2

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
			//p.next()
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

	p.indentation -= 2
	p.writeWithIndentation("</term>\r\n")
}

func (p *Parser) compileExpressionList() {
	p.writeWithIndentation("<expressionList>\r\n")
	p.indentation += 2

	if p.tok != SYMBOL && p.lit != ")" {
		p.compileExpressions()
		for p.tok == SYMBOL && p.lit == "," {
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressions()
		}
	}

	// if after '(' is a '(' -> new expression
	if p.lit == "(" {
		p.next()
		if p.lit == ")" {
			goto exit
		}
		p.compileExpressions()
		for p.tok == SYMBOL && p.lit == "," {
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressions()
		}
	}

exit:
	p.indentation -= 2
	p.writeWithIndentation("</expressionList>\r\n")
}

func (p *Parser) compileVarDec() {
	p.writeWithIndentation("<varDec>\r\n")
	p.indentation += 2

	p.writeTemplate()
	p.next()
	p.compileTypeAndVarName()

	p.indentation -= 2
	p.writeWithIndentation("</varDec>\r\n")
}

func (p *Parser) compileTypeAndVarName() {
	p.writeTemplate()
	p.next()

	p.writeTemplate()
	p.next()

	for p.lit == "," {
		p.writeTemplate() // symbol
		p.next()

		p.writeTemplate() // identifier
		p.next()
	}

	p.writeTemplate()
	p.next()
}

func (p *Parser) compileParameterList() {
	p.writeWithIndentation("<parameterList>\r\n")
	p.indentation += 2

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

	p.indentation -= 2
	p.writeWithIndentation("</parameterList>\r\n")
}

func (p *Parser) compileClassVarDec() {
	p.writeWithIndentation("<classVarDec>\r\n")
	p.indentation += 2
	p.writeTemplate()
	p.next()

	p.compileTypeAndVarName()

	p.indentation -= 2
	p.writeWithIndentation("</classVarDec>\r\n")
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

func (p *Parser) writeWithIndentation(str string) {
	p.writeIndentation()
	p.write(str)
}

func (p *Parser) writeIndentation() {
	for i := 0; i < p.indentation; i++ {
		p.write(" ")
	}
}

const template = "<%v> %v </%v>\r\n"

var i = 0

// writeTemplate writes KEYWORD, IDENT and SYMBOL token
func (p *Parser) writeTemplate() {
	if p.tok == EOF {
		return
	}

	if p.lit == "moveLeft" {
		println("here")
	}

	if p.lit == "drawRectangle" {
		i++
		if i == 7 {
			println("her")
		}
	}

	if p.tok == SYMBOL {
		p.writeSymbol()
		return
	}

	if p.tok == INT {
		p.writeIndentation()
		p.write(fmt.Sprintf(template, "integerConstant", p.lit, "integerConstant"))
		return
	}

	if p.tok == CHAR {
		p.writeIndentation()
		p.write(fmt.Sprintf(template, "stringConstant", p.lit, "stringConstant"))
		return
	}

	p.writeIndentation()
	p.write(fmt.Sprintf(template, p.tok, p.lit, p.tok))
}

func (p *Parser) writeSymbol() {
	var symbol string
	switch p.lit {
	case "<":
		symbol = "&lt;"
	case ">":
		symbol = "&gt;"
	case "&":
		symbol = "&amp;"
	default:
		symbol = p.lit
	}
	p.writeIndentation()
	p.write("<symbol> " + symbol + " </symbol>\r\n")
}

func (p *Parser) write(str string) {
	p.out.WriteString(str)
}

func (p *Parser) Out() string {
	return p.out.String()
}
