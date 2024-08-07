package main_test

import (
	"os"
	"strings"
	"testing"

	pkg "project11"

	"github.com/stretchr/testify/assert"
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
			"1. Seven Main.jack",
			fields{
				dest: "./test/Seven/Main.jack",
				want: "./test/Seven/Main.vm",
				out:  "./test/Seven/Main.an",
			},
		},
		{
			"2. 1 if",
			fields{
				dest: "./test/test/if1.jack",
				want: "./test/test/if1.vm",
				out:  "./test/test/if1.an",
			},
		},
		{
			"2. 1 while",
			fields{
				dest: "./test/test/while1.jack",
				want: "./test/test/while1.vm",
				out:  "./test/test/while1.an",
			},
		},
		{
			"3. 2 if",
			fields{
				dest: "./test/test/if2.jack",
				want: "./test/test/if2.vm",
				out:  "./test/test/if2.an",
			},
		},
		{
			"4. ./ConvertToBin",
			fields{
				dest: "./test/ConvertToBin/Main.jack",
				want: "./test/ConvertToBin/Main.vm",
				out:  "./test/ConvertToBin/Main.an",
			},
		},
		{
			"5. ./Square/Main.jack",
			fields{
				dest: "./test/Square/Main.jack",
				want: "./test/Square/Main.vm",
				out:  "./test/Square/Main.an",
			},
		},
		{
			"6. ./Square/Square.jack",
			fields{
				dest: "./test/Square/Square.jack",
				want: "./test/Square/Square.vm",
				out:  "./test/Square/Square.an",
			},
		},
		{
			"7. ./Square/SquareGame.jack",
			fields{
				dest: "./test/Square/SquareGame.jack",
				want: "./test/Square/SquareGame.vm",
				out:  "./test/Square/SquareGame.an",
			},
		},
		{
			"8. ./Average/Average.jack",
			fields{
				dest: "./test/Average/Main.jack",
				want: "./test/Average/Main.vm",
				out:  "./test/Average/Main.an",
			},
		},
		{
			"9. ./Pong/Main.jack",
			fields{
				dest: "./test/Pong/Main.jack",
				want: "./test/Pong/Main.vm",
				out:  "./test/Pong/Main.an",
			},
		},
		{
			"10. ./Pong/Ball.jack",
			fields{
				dest: "./test/Pong/Ball.jack",
				want: "./test/Pong/Ball.vm",
				out:  "./test/Pong/Ball.an",
			},
		},
		{
			"10. ./Pong/Bat.jack",
			fields{
				dest: "./test/Pong/Bat.jack",
				want: "./test/Pong/Bat.vm",
				out:  "./test/Pong/Bat.an",
			},
		},
		{
			"11. ./Pong/PongGame.jack",
			fields{
				dest: "./test/Pong/PongGame.jack",
				want: "./test/Pong/PongGame.vm",
				out:  "./test/Pong/PongGame.an",
			},
		},
		{
			"12. ./ComplexArrays/Main.jack",
			fields{
				dest: "./test/ComplexArrays/Main.jack",
				want: "./test/ComplexArrays/Main.vm",
				out:  "./test/ComplexArrays/Main.an",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, err := os.ReadFile(tt.fields.dest)
			assert.NoError(t, err)
			var p pkg.Parser
			p.Init(src)
			p.ParseFile()

			err = os.WriteFile(tt.fields.out, []byte(p.VmOut()), 0644)
			assert.NoError(t, err)

			want, err := os.ReadFile(tt.fields.want)
			assert.NoError(t, err)
			assert.Equal(
				t,
				strings.TrimSuffix(string(want), "\n"),
				p.VmOut(),
			)
		})
	}
}
