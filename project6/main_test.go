package main

import (
	"testing"
)

var sb = &SymbolTable{}

func TestMain(t *testing.M) {
	sb = NewSymbolTable()
	t.Run()
}

func Test_parse(t *testing.T) {
	type args struct {
		instruction string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_A_instruction_@value",
			args: args{
				instruction: "@2",
			},
			want: "0000000000000010",
		},
		{
			name: "test_A_instruction_@value",
			args: args{
				instruction: "@3",
			},
			want: "0000000000000011",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parse(sb, tt.args.instruction); got != tt.want {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
