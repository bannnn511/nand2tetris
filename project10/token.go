package main

import (
	"strconv"
)

type Token int

// The list of tokens.
const (
	// Speicial tokens
	ILLEGAL Token = iota
	EOF
	COMMENT
	IDENT

	KEYWORD
	keyword_beg
	CLASS
	CONSTRUCTOR
	FUNCTION
	METHOD
	FIELD
	STATIC
	VAR
	INT
	CHAR
	BOOLEAN
	VOID
	TRUE
	FALSE
	NULL
	THIS
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	keyword_end

	SYMBOL
	symbol_beg
	LBRACE // {
	LPAREN // (
	LBRACK // [

	RBRACE // }
	RPAREN // )
	RBRACK // ]
	COMMA  // ,
	symbol_end

	operator_beg
	ADD // +
	SUB // -
	QUO // /
	AND // &
	OR  // |
	GTR // >
	LSS // <
	EQL // =
	NOT // ~
	operator_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",
	IDENT:   "identifier",

	KEYWORD:     "keyword",
	CLASS:       "class",
	CONSTRUCTOR: "constructor",
	FUNCTION:    "function",
	METHOD:      "method",
	FIELD:       "field",
	STATIC:      "static",
	VAR:         "var",
	INT:         "int",
	CHAR:        "char",
	BOOLEAN:     "boolean",
	VOID:        "void",
	TRUE:        "true",
	FALSE:       "false",
	NULL:        "null",
	THIS:        "this",
	LET:         "let",
	DO:          "do",
	IF:          "if",
	ELSE:        "else",
	WHILE:       "while",
	RETURN:      "return",

	LBRACE: "LBRACE",
	LPAREN: "LPAREN",
	LBRACK: "LBRACK",

	RBRACE: "RBRACE",
	RPAREN: "RPAREN",
	RBRACK: "RBRACK",
	COMMA:  "COMMA",

	ADD: "ADD",
	SUB: "SUB",
	QUO: "QUO",
	AND: "AND",
	OR:  "OR",
	GTR: "GTR",
	LSS: "LSS",
	EQL: "EQL",
	NOT: "NOT",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}

	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if _, ok := keywords[ident]; ok {
		return KEYWORD
	}

	return IDENT
}
