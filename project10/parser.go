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
	elements []Ast // list of lexical elements
	tok      Token // current token
	fileName string
	scanner  Scanner
	out      string
}

func (p *Parser) init(filename string, src []byte) {
	p.scanner.Init(src)
	p.fileName = filename
	p.tok = START
	p.elements = make([]Ast, 0, len(src))
	p.append(Ast{START, ""})
}

func (p *Parser) parseFile() {
	for p.tok != EOF {
		p.next()
	}
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

func (p *Parser) next() {
	tok, lit := p.scanner.Scan()
	p.tok = tok
	p.elements = append(p.elements, Ast{tok, lit})
}

func (p *Parser) expect(tok Token) {
	if p.tok != tok {
		p.errorExpected("'" + tok.String() + "'")
	}
	p.next()
}

func (p *Parser) errorExpected(msg string) {
	msg = "expected" + msg
	printErr(msg)
}
