package main

import "fmt"

type VariableKind int

const (
	Static VariableKind = iota
	Field
	Arg
	Var
)

type Symbol struct {
	name  string
	sType Token
	kind  VariableKind
	n     uint32
}

type SymbolTable struct {
	m     map[string]Symbol
	count map[VariableKind]uint32
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		m:     make(map[string]Symbol),
		count: make(map[VariableKind]uint32),
	}
}

func (sb *SymbolTable) Define(tok Token, name string, kind VariableKind) {
	count, ok := sb.count[kind]
	if !ok {
		sb.count[kind] = 0
	} else {
		count++
	}

	_, ok = sb.m[name]
	if ok {
		return
	}

	fmt.Println("here", count)

	sb.m[name] = Symbol{
		name:  name,
		sType: tok,
		kind:  kind,
		n:     count,
	}
	sb.count[kind] = count
}

func (sb *SymbolTable) IndexOf(name string) uint32 {
	return sb.m[name].n
}
