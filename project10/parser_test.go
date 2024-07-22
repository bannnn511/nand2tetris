package main_test

import (
	"fmt"
	pkg "project10"
	"testing"
)

func TestParser_compileWhile(t *testing.T) {

	tests := []struct {
		name string
		src  []byte
	}{
		{
			"1. test main class",
			[]byte(`
			class Main {
				static boolean test;

				function void main() {
     				var SquareGame game;
        			let game = game;
    		}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parser pkg.Parser
			parser.Init("", tt.src)
			parser.ParseFile()
			fmt.Println(parser.Out.String())
		})
	}
}
