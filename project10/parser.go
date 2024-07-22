package main

import (
	"fmt"
	"strings"
)

const debug = false

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
	Out         strings.Builder
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
	p.write("<class>\n")
	p.indentation++

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
	p.indentation--
	p.write("</class>\n")
}

func (p *Parser) write(str string) {
	p.Out.WriteString(str)
}

func (p *Parser) compileSubroutine() {
	p.writeWithIndentation("<subroutineDec>\n")
	p.indentation++
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

	p.writeWithIndentation("<subroutineBody>\n")
	p.indentation++
	p.writeTemplate()

	p.next()
	for p.lit == "var" {
		p.compileVarDec()
	}

	p.compileStatements()

	p.writeTemplate() // symbol
	p.indentation--
	p.writeWithIndentation("</subroutineBody>\n")
	p.indentation--
	p.writeWithIndentation("</subroutineDec>\n")
}

func (p *Parser) compileStatements() {
	p.writeWithIndentation("<statements>\n")
	p.indentation++

	for p.tok == KEYWORD {
		switch p.lit {
		case "let":
			p.compileLet()
		case "while":
			p.compileWhile()
		case "if":
		case "do":
		case "return":
		}
	}

	p.indentation--
	p.writeWithIndentation("</statements>\n")
}

// 'let' varName ('[' expression ']')? '=' expression ';'
func (p *Parser) compileLet() {
	p.writeWithIndentation("<letStatement>\n")
	p.indentation++

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

	p.indentation--
	p.writeWithIndentation("</letStatement>\n")
	p.next()
}

// 'while' '(' expression ')' '{' statements '}'
func (p *Parser) compileWhile() {
	p.writeWithIndentation("<whileStatement>\n")
	p.indentation++

	p.next()
	p.writeTemplate() // (

	p.next()
	p.compileExpressions()

	p.next()
	p.writeTemplate() // )

	p.next()
	p.writeTemplate() // {

	p.next()
	p.compileStatements()

	p.next()
	p.writeTemplate() // }

	p.indentation--
	p.writeWithIndentation("</whileStatement>\n")
	p.next()
}

func (p *Parser) compileExpressions() {
	p.writeWithIndentation("<expression>\n")
	p.indentation++

	p.CompileTerm()
	for p.tok == SYMBOL && IsOp(p.lit) {
		p.writeTemplate()
		p.next()
		p.CompileTerm()
	}

	p.indentation--
	p.writeWithIndentation("</expression>\n")
}

func (p *Parser) CompileTerm() {
	p.writeWithIndentation("<term>\n")
	p.indentation++

	if p.tok == KEYWORD {
		p.writeTemplate()
	} else if p.tok == IDENT {
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
		} else if p.lit == "~" || p.lit == "-" {
			p.writeTemplate() // symbol
			p.next()
			p.CompileTerm()
		}
	}

	p.indentation--
	p.writeWithIndentation("</term>\n")
}

func (p *Parser) compileExpressionList() {
	p.writeWithIndentation("<expressionList>\n")
	p.indentation++

	if p.tok != SYMBOL && p.lit != ")" ||
		p.lit == "(" {
		p.compileExpressions()
		for p.tok == SYMBOL && p.lit == "," {
			p.writeTemplate() // symbol
			p.next()
			p.compileExpressions()
		}
	}

	p.indentation--
	p.writeWithIndentation("</expressionList>\n")
}

func (p *Parser) compileVarDec() {
	p.writeWithIndentation("<varDec>\n")
	p.indentation++

	p.writeTemplate()
	p.next()

	p.indentation--
	p.writeWithIndentation("</varDec>\n")
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
	p.writeWithIndentation("<parameterList>\n")
	p.indentation++

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

	p.indentation--
	p.writeWithIndentation("</parameterList>\n")
}

func (p *Parser) compileClassVarDec() {
	p.writeWithIndentation("<classVarDec>\n")
	p.indentation++
	p.writeTemplate()
	p.next()

	p.compileTypeAndVarName()

	p.indentation--
	p.writeWithIndentation("</classVarDec>\n")
}

func (p *Parser) append(ele Ast) {
	p.elements = append(p.elements, ele)
}

func (p *Parser) GetXML() string {
	var sb strings.Builder
	for _, node := range p.elements {
		switch node.tok {
		case START:
			sb.WriteString("<tokens>\n")
		case EOF:
			sb.WriteString("</tokens>")
		case COMMENT:
			continue
		default:
			v := fmt.Sprintf("<%v> %v </%v>\n", node.tok, node.lit, node.tok)
			sb.WriteString(v)
		}
	}

	return sb.String()
}

func (p *Parser) printTree() {
	for _, node := range p.elements {
		switch node.tok {
		case START:
			fmt.Println("<tokens>")
		case EOF:
			fmt.Println("</tokens>")
		case COMMENT:
			continue
		default:
			fmt.Printf("<%v> %v </%v>\n", node.tok, node.lit, node.tok)
		}
	}
}

// func (p *Parser) tokenized() {
// 	for p.tok != EOF {
// 		p.next()
// 	}
// }

func (p *Parser) next() {
	tok, lit := p.scanner.Scan()
	if tok == COMMENT {
		p.next()
		return
	}
	p.tok = tok
	p.lit = lit
	if debug {
		fmt.Println("current state", p.tok, p.lit)
	}
	p.elements = append(p.elements, Ast{tok, lit})
}

const template = "<%v> %v </%v>\n"

func (p *Parser) writeWithIndentation(str string) {
	p.writeIndentation()
	p.write(str)
}

func (p *Parser) writeIndentation() {
	for i := 0; i < p.indentation; i++ {
		p.write(" ")
	}
}

// writeTemplate writes KEYWORD, IDENT and SYMBOL token
func (p *Parser) writeTemplate() {
	if p.tok == EOF {
		return
	}
	p.writeIndentation()
	p.write(fmt.Sprintf(template, p.tok, p.lit, p.tok))
}
