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
					lit: `negative`,
				},
			},
		},
		{
			name: "8. test comment //",
			src:  []byte("// aaa"),
			wants: []lexical{
				{
					tok: pkg.COMMENT,
					lit: "// aaa",
				},
			},
		},
		{
			name: "9. test comment /** comment */",
			src:  []byte("/** aaa */"),
			wants: []lexical{
				{
					tok: pkg.COMMENT,
					lit: "/** aaa */",
				},
			},
		},
		{
			name: `10.let length = Keyboard.readInt("HOW MANY NUMBERS? ");`,
			src:  []byte(`let length = Keyboard.readInt("HOW MANY NUMBERS? ");`),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "let",
				},
				{
					tok: pkg.IDENT,
					lit: "length",
				},
				{
					tok: pkg.SYMBOL,
					lit: "=",
				},
				{
					tok: pkg.IDENT,
					lit: "Keyboard",
				},
				{
					tok: pkg.SYMBOL,
					lit: ".",
				},
				{
					tok: pkg.IDENT,
					lit: "readInt",
				},
				{
					tok: pkg.SYMBOL,
					lit: "(",
				},
				{
					tok: pkg.CHAR,
					lit: `HOW MANY NUMBERS? `,
				},
				{
					tok: pkg.SYMBOL,
					lit: ")",
				},
				{
					tok: pkg.SYMBOL,
					lit: ";",
				},
			},
		},
		{
			name: "11. let a[1] = a[2];",
			src:  []byte(`let a[1] = a[2];`),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "let",
				},
				{
					tok: pkg.IDENT,
					lit: "a",
				},
				{
					tok: pkg.SYMBOL,
					lit: "[",
				},
				{
					tok: pkg.INT,
					lit: "1",
				},
				{
					tok: pkg.SYMBOL,
					lit: "]",
				},
				{
					tok: pkg.SYMBOL,
					lit: "=",
				},
				{
					tok: pkg.IDENT,
					lit: "a",
				},
				{
					tok: pkg.SYMBOL,
					lit: "[",
				},
				{
					tok: pkg.INT,
					lit: "2",
				},
				{
					tok: pkg.SYMBOL,
					lit: "]",
				},
				{
					tok: pkg.SYMBOL,
					lit: ";",
				},
			},
		},
		{
			name: "12. var int length;",
			src:  []byte("var int length;"),
			wants: []lexical{
				{
					tok: pkg.KEYWORD,
					lit: "var",
				},
				{
					tok: pkg.KEYWORD,
					lit: "int",
				},
				{
					tok: pkg.KEYWORD,
					lit: "length",
				},
				{
					tok: pkg.SYMBOL,
					lit: ";",
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

			t.Run("test EOF", func(t *testing.T) {
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
		})
	}
}
