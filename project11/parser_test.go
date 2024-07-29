package main_test

import (
	"os"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, err := os.ReadFile(tt.fields.dest)
			assert.NoError(t, err)
			var p pkg.Parser
			p.Init("", src)
			p.ParseFile()

			err = os.WriteFile(tt.fields.out, []byte(p.VmOut()), 0644)
			assert.NoError(t, err)

			want, err := os.ReadFile(tt.fields.want)
			assert.NoError(t, err)
			assert.Equal(t, string(want), p.VmOut())
		})
	}
}
