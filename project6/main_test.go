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

func Test_parseCInstruction(t *testing.T) {
	type args struct {
		instruction string
	}
	tests := []struct {
		name     string
		args     args
		wantDest string
		wantComp string
		wantJmp  string
	}{
		{
			name:     "1. test_C_instruction_D=A",
			args:     args{instruction: "D=A"},
			wantDest: "D",
			wantComp: "A",
			wantJmp:  "",
		},
		{
			name:     "2. test_C_instruction_D=D+A",
			args:     args{instruction: "D=D+A"},
			wantDest: "D",
			wantComp: "D+A",
			wantJmp:  "",
		},
		{
			name:     "3. test_C_instruction_M=D",
			args:     args{instruction: "M=D"},
			wantDest: "M",
			wantComp: "D",
			wantJmp:  "",
		},
		{
			name:     "3. test_C_instruction_0;JMP",
			args:     args{instruction: "0;JMP"},
			wantDest: "",
			wantComp: "0",
			wantJmp:  "JMP",
		},
		{
			name:     "4. test_C_instruction_D;JEQ",
			args:     args{instruction: "D;JEQ"},
			wantDest: "",
			wantComp: "D",
			wantJmp:  "JEQ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest, comp, jmp, _ := parseCInstruction(tt.args.instruction)
			if dest != tt.wantDest ||
				comp != tt.wantComp ||
				jmp != tt.wantJmp {
				t.Errorf(
					"parseCInstruction() = %v %v %v, want %v %v %v",
					dest, comp, jmp, tt.wantDest, tt.wantComp, tt.wantJmp,
				)
			}
		})
	}
}

func Test_code(t *testing.T) {
	type args struct {
		dest, comp, jump string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1. test_code_D=A",
			args: args{
				dest: "D",
				comp: "A",
				jump: "",
			},
			want: "1110110000010000",
		},
		{
			name: "2. test_code_D=D+A",
			args: args{
				dest: "D",
				comp: "D+A",
				jump: "",
			},
			want: "1110000010010000",
		},
		{
			name: "3. test_code_M=D",
			args: args{
				dest: "M",
				comp: "D",
				jump: "",
			},
			want: "1110001100001000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := code(tt.args.dest, tt.args.comp, tt.args.jump)
			if tt.want != got {
				t.Errorf(
					"code() = %v, want %v",
					got, tt.want,
				)
			}
		})
	}
}
