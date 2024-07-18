package main

type Parser struct {
	fileName string
	tok      Token
	scanner  Scanner
}

func (p *Parser) init(filename string, src []byte) {
	p.scanner.Init(src)
	p.fileName = filename

	p.next()
}

func (p *Parser) parseFile() {
	for p.tok != EOF {
		tok, _ := p.scanner.Scan()
		p.tok = tok
	}
}

func (p *Parser) expect(tok Token) {
	if p.tok != tok {
		p.errorExpected("'" + tok.String() + "'")
	}
	p.next()
}

func (p *Parser) next0() {
	for {

	}
}

func (p *Parser) next() {
	p.next0()
}

func (p *Parser) errorExpected(msg string) {
	msg = "expected" + msg
	printErr(msg)
}
