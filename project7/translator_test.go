package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslator_WritePushPop(t *testing.T) {

	type args struct {
		cmdType CommandType
		segment string
		idx     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_Translator_WritePushPop_1",
			args: args{
				cmdType: CPUSH,
				segment: "local",
				idx:     1,
			},
			want: `@1
D=A
@LCL
A=M
M=D
@LCL
M=M+1
`,
		},
		{
			name: "test_Translator_WritePushPop_2",
			args: args{
				cmdType: CPOP,
				segment: "local",
				idx:     0,
			},
			want: `@0
D=A
@LCL
M=M-1
@LCL
A=M
M=D
`},
		{
			name: "test_Translator_WritePushPop_3",
			args: args{
				cmdType: CPUSH,
				segment: "argument",
				idx:     1,
			},
			want: `@1
D=A
@ARG
A=M
M=D
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Translator{}
			assert.Equalf(
				t,
				tt.want,
				w.WritePushPop(tt.args.cmdType, tt.args.segment, tt.args.idx),
				"WritePushPop(%v, %v, %v)", tt.args.cmdType, tt.args.segment, tt.args.idx)
		})
	}
}
