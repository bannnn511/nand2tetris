package main_test

import (
	"testing"

	pkg "project11"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTable_Define(t *testing.T) {
	type fields struct {
		name      string
		tok       pkg.Token
		kind      pkg.VariableKind
		wantCount uint32
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			"1. var int value",
			fields{
				name:      "value",
				tok:       pkg.INT,
				kind:      pkg.Var,
				wantCount: 0,
			},
		},
		{
			"2. var int value2",
			fields{
				name:      "value2",
				tok:       pkg.INT,
				kind:      pkg.Var,
				wantCount: 1,
			},
		},
	}

	sb := pkg.NewSymbolTable()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb.Define(tt.fields.tok, tt.fields.name, tt.fields.kind)
			want := tt.fields.wantCount
			got := sb.IndexOf(tt.fields.name)
			assert.Equal(
				t,
				want,
				got,
			)
		})
	}
}
