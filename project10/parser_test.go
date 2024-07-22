package main_test

import (
	"fmt"
	pkg "project10"
	"testing"
)

func TestParser_ParseFile(t *testing.T) {

	tests := []struct {
		name string
		src  []byte
	}{
		{
			"1. test main class",
			[]byte(`
class Square {

   field int x, y; // screen location of the square's top-left corner
   field int size; // length of this square, in pixels

   /** Constructs a new square with a given location and size. */
   constructor Square new(int Ax, int Ay, int Asize) {
      let x = Ax;
      let y = Ay;
      let size = Asize;
      do draw();
      return this;
   }
		`),
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
