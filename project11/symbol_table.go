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
	index uint32
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

	sb.m[name] = Symbol{
		name:  name,
		sType: tok,
		kind:  kind,
		index: count,
	}
	fmt.Println(sb.m[name].kind)
	sb.count[kind] = count
}

// IndexOf returns count index variable name
func (sb *SymbolTable) IndexOf(name string) uint32 {
	return sb.m[name].index
}

// TypeOf returns Token type of variable name
func (sb *SymbolTable) TypeOf(name string) Token {
	return sb.m[name].sType
}

// KindOf returns variable kind
func (sb *SymbolTable) KindOf(name string) VariableKind {
	return sb.m[name].kind
}

// VarCount returns number of variable of the given kind
func (sb *SymbolTable) VarCount(kind VariableKind) uint32 {
	return sb.count[kind] + 1
}
