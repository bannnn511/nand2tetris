package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_writePushFormat(t *testing.T) {
	type args struct {
		segment string
		idx     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_writePushFormat_1",
			args: args{
				segment: "LCL",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, writePushFormat(tt.args.segment, tt.args.idx), "writePushFormat(%v, %v)", tt.args.segment, tt.args.idx)
		})
	}
}
