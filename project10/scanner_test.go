package main_test

import (
	pkg "project10"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner_Scan(t *testing.T) {
	type lexical struct {
		tok pkg.Token
		lit string
	}

	tests := []struct {
		name  string
		src   []byte
		wants []lexical
	}{
		{
			name: "1. class Main {",
			src:  []byte("class Main {"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "class",
				},
				{
					tok: pkg.IDENT,
					lit: "Main",
				},
				{
					tok: pkg.SYMBOL,
					lit: "{",
				},
			},
		},
		{
			name: "2. static boolean test;",
			src:  []byte("static boolean test;"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "static",
				},
				{
					tok: pkg.KEYWORD,
					lit: "boolean",
				},
				{
					tok: pkg.IDENT,
					lit: "test",
				},
				{
					tok: pkg.SYMBOL,
					lit: ";",
				},
			},
		},
		{
			name: "3. function void main() {",
			src:  []byte("function void main() {"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "function",
				},
				{
					tok: pkg.KEYWORD,
					lit: "void",
				},
				{
					tok: pkg.IDENT,
					lit: "main",
				},
				{
					tok: pkg.SYMBOL,
					lit: "(",
				},
				{
					tok: pkg.SYMBOL,
					lit: ")",
				},
				{
					tok: pkg.SYMBOL,
					lit: "{",
				},
			},
		},
		{
			name: "4. var SquareGame game;",
			src:  []byte("var SquareGame game;"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "var",
				},
				{
					tok: pkg.IDENT,
					lit: "SquareGame",
				},
				{
					tok: pkg.IDENT,
					lit: "game",
				},
				{
					tok: pkg.SYMBOL,
					lit: ";",
				},
			},
		},
		{
			name: "5. let game = game;",
			src:  []byte("let game = game;"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "let",
				},
				{
					tok: pkg.IDENT,
					lit: "game",
				},
				{
					tok: pkg.SYMBOL,
					lit: "=",
				},
				{
					tok: pkg.IDENT,
					lit: "game",
				},
				{
					tok: pkg.SYMBOL,
					lit: ";",
				},
			},
		},
		{
			name: "6. if (x<0) {",
			src:  []byte("if (x<0) {"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "if",
				},
				{
					tok: pkg.SYMBOL,
					lit: "(",
				},
				{
					tok: pkg.IDENT,
					lit: "x",
				},
				{
					tok: pkg.SYMBOL,
					lit: "<",
				},
				{
					tok: pkg.INT,
					lit: "0",
				},
				{
					tok: pkg.SYMBOL,
					lit: ")",
				},
				{
					tok: pkg.SYMBOL,
					lit: "{",
				},
			},
		},
		{
			name: "7. let state = \"negative\"",
			src:  []byte(`let state = "negative"`),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "let",
				},
				{
					tok: pkg.IDENT,
					lit: "state",
				},
				{
					tok: pkg.SYMBOL,
					lit: "=",
				},
				{
					tok: pkg.CHAR,
					lit: "negative",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var scanner pkg.Scanner
			scanner.Init(tt.src)
			for _, want := range tt.wants {
				tok, lit := scanner.Scan()
				got := lexical{
					tok,
					lit,
				}
				assert.Equalf(
					t,
					want,
					got,
					"got %v, want %v",
					got,
					want,
				)

			}

			tok, lit := scanner.Scan()
			assert.Equal(
				t,
				pkg.EOF,
				tok,
				"expexted EOF, got %v and %v",
				tok,
				lit,
			)
		})
	}
}
