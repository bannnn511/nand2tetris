package main_test

import (
	"testing"

	pkg "project11"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTable(t *testing.T) {
	type fields struct {
		name      string
		tok       string
		kind      pkg.VariableKind
		wantCount uint
		wantTok   pkg.Token
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			"1. var int value",
			fields{
				name:      "value",
				tok:       "int",
				kind:      pkg.Var,
				wantCount: 0,
				wantTok:   pkg.INT,
			},
		},
		{
			"2. var int value2",
			fields{
				name:      "value2",
				tok:       "int",
				kind:      pkg.Var,
				wantCount: 1,
				wantTok:   pkg.INT,
			},
		},
	}

	sb := pkg.NewSymbolTable()

	count := uint32(0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb.Define(tt.fields.tok, tt.fields.name, tt.fields.kind)
			want := tt.fields.wantCount
			got := sb.IndexOf(tt.fields.name)
			assert.Equal(
				t,
				want,
				got,
				"IndexOf",
			)

			assert.Equal(
				t,
				sb.KindOf(tt.fields.name),
				pkg.Var,
				"KindOf",
			)

			_, gotType := sb.TypeOf(tt.fields.name)

			assert.Equal(
				t,
				tt.fields.wantTok,
				gotType,
				"TypeOf",
			)

			count++
		})
	}
}
