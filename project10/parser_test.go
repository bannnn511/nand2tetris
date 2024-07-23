package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParser_ParseFile(t *testing.T) {
	type fields struct {
		dest string
		want string
		out  string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"1. ExpressionLessSquare Main.jack",
			fields{
				dest: "./test/ExpressionLessSquare/Main.jack",
				want: "./test/ExpressionLessSquare/Main.xml",
				out:  "./test/ExpressionLessSquare/Main.test.xml",
			},
		},
		// {
		//	"2. ExpressionLessSquare Square.jack",
		//	fields{
		//		dest: "./test/ExpressionLessSquare/Square.jack",
		//		want: "./test/ExpressionLessSquare/Square.xml",
		//		out:  "./test/ExpressionLessSquare/Square.test.xml",
		//	},
		// },
		// {
		//	"3. ExpressionLessSquare SquareGame.jack",
		//	fields{
		//		dest: "./test/ExpressionLessSquare/SquareGame.jack",
		//		want: "./test/ExpressionLessSquare/SquareGame.xml",
		//		out:  "./test/ExpressionLessSquare/SquareGame.test.xml",
		//	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, err := os.ReadFile(tt.fields.dest)
			assert.NoError(t, err)
			var p Parser
			p.Init("", src)
			p.ParseFile()

			err = os.WriteFile(tt.fields.out, []byte(p.Out()), 0644)
			assert.NoError(t, err)
		})
	}
}
