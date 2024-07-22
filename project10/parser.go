package main

import (
	"fmt"
	"strings"
)

const debug = true

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

	}

	p.writeTemplate()
	p.indentation--
	p.write("</class>\n")
}

func (p *Parser) write(str string) {
	p.Out.WriteString(str)
}

func (p *Parser) compileSubroutine() {
	p.writeIndentation()
	p.write("<subroutineDec>\n")

	p.write("</subroutineDec>\n")
}

func (p *Parser) compileClassVarDec() {
	p.writeIndentation()
	p.write("<classVarDec>\n")
	p.indentation++
	p.writeTemplate()
	p.next()

	p.compileTypeAndVarName()

	p.indentation--
	p.write("</classVarDec>\n")
}

func (p *Parser) compileTypeAndVarName() {
	p.writeTemplate()
	p.next()
	p.writeTemplate()
	p.next()
	for p.lit == "," {
		p.writeTemplate()
		p.next()
		p.writeTemplate()
		p.next()
	}
	p.writeTemplate()
	p.next()
}

func (p *Parser) compileWhile() {
	p.expect(KEYWORD, "while")
	p.next()
	p.expect(SYMBOL, "(")
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
	p.tok = tok
	p.lit = lit
	if debug {
		fmt.Println("current state", p.tok, p.lit)
	}
	p.elements = append(p.elements, Ast{tok, lit})
}

func (p *Parser) expect(tok Token, lit string) {
	if p.tok != tok {
		p.errorExpected("'" + tok.String() + "'")
	}

	if p.lit != lit {
		p.errorExpected("'" + lit + "'")
	}
}

func (p *Parser) errorExpected(msg string) {
	msg = "expected" + msg
	printErr(msg)
}

const template = "<%v> %v </%v>\n"

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
