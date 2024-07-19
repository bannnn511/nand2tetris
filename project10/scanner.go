package main

import (
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	// immutable state
	src []byte

	// scanning state
	ch       rune // current character
	offset   int  // character offset
	rdOffset int  // reading offset (position after current character)
}

const eof = -1

func (s *Scanner) Init(src []byte) {
	s.src = src
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0

	s.next()
}

func (s *Scanner) Scan() (tok Token, lit string) {
	s.skipWhiteSpace()

	if s.isEOF() {
		return EOF, ""
	}

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = Lookup(lit)
	case isDecimal(ch):
		tok, lit = s.scanNumber()
	default:
		s.next()
		switch ch {
		case '"':
			tok = CHAR
			lit = s.scanString()
		case eof:
			tok = EOF
		case '/':
			if s.ch == '/' || s.ch == '*' {
				// comment
				tok = COMMENT
				lit = s.scanComment()
			} else {
				// division
				panic("implement division case")
			}
		default:
			lit = string(ch)
			tok = SYMBOL
		}
	}

	return
}

func (s *Scanner) next() {
	if s.rdOffset >= len(s.src) {
		s.ch = eof
		return
	}

	s.offset = s.rdOffset
	if s.ch == '\n' {
		panic("handle new line")
	}

	r, w := rune(s.src[s.rdOffset]), 1
	switch {
	case r == 0:
		panic("illegal character NULl")
	case r > utf8.RuneSelf:
		r, w = utf8.DecodeRune(s.src[s.rdOffset:])
		if r == utf8.RuneError && w == 1 {
			printErr("illegal UTF-8 encoding")
		}
	}
	s.ch = r
	s.rdOffset += w
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for rdOffset, b := range s.src[s.rdOffset:] {
		if 'a' <= b && b <= 'z' ||
			'A' <= b && b <= 'Z' ||
			'0' <= b && b <= '9' {
			continue
		}
		s.rdOffset += rdOffset
		s.offset = s.rdOffset
		s.ch = rune(b)
		s.rdOffset++
		break
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanNumber() (tok Token, lit string) {
	offs := s.offset
	tok = INT
	for rdOffset, b := range s.src[s.rdOffset:] {
		if '0' <= b && b <= '9' {
			continue
		}
		s.rdOffset += rdOffset
		s.offset = s.rdOffset
		s.ch = rune(b)
		s.rdOffset++
		break
	}
	lit = string(s.src[offs:s.offset])

	return
}

func (s *Scanner) scanString() string {
	offs := s.offset
	for {
		ch := s.ch
		s.next()
		if ch == '"' {
			break
		}
	}

	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanComment() string {
	offs := s.offset - 1
	//-style comment
	if s.ch == '/' {
		s.next()
		for s.ch != '\n' && s.ch > 0 {
			s.next()
		}
		goto exit
	}

	//**-style comment
	s.next()
	s.next()
	for s.ch != '/' && s.ch > 0 {
		s.next()
	}

exit:
	lit := string(s.src[offs:s.rdOffset])
	s.next()
	return lit
}

func (s *Scanner) skipWhiteSpace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *Scanner) isEOF() bool {
	return s.ch == eof
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' ||
		ch == '_' ||
		ch >= utf8.RuneSelf && unicode.IsLetter(ch)

}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
