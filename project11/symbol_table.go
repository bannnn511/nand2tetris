package main

import "fmt"

type VariableKind int

const (
	Undefined VariableKind = iota
	Static
	Field
	Arg
	Var
)

type Symbol struct {
	name  string
	sType Token  // primitive token
	uType string // user defined type
	kind  VariableKind
	index uint
}

type SymbolTable struct {
	m     map[string]Symbol
	count map[VariableKind]uint
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		m:     make(map[string]Symbol),
		count: make(map[VariableKind]uint),
	}
}

func (sb *SymbolTable) Define(
	tok string,
	name string,
	kind VariableKind,
) {
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

	token := Lookup2(tok)

	sb.m[name] = Symbol{
		name:  name,
		sType: token,
		uType: tok,
		kind:  kind,
		index: count,
	}
	sb.count[kind] = count
}

// IndexOf returns count index variable name
func (sb *SymbolTable) IndexOf(name string) uint {
	return sb.m[name].index
}

// TypeOf returns Token type of variable name
func (sb *SymbolTable) TypeOf(name string) (string, Token) {
	v, ok := sb.m[name]
	if !ok {
		return "", ILLEGAL
	}

	if v.sType == USR {
		return v.uType, v.sType
	}

	return v.sType.String(), v.sType
}

// KindOf returns variable kind
func (sb *SymbolTable) KindOf(name string) VariableKind {
	return sb.m[name].kind
}

// GetSegment returns variable kind and index
func (sb *SymbolTable) GetSegment(name string) (VariableKind, uint) {
	kind := sb.KindOf(name)
	idx := sb.IndexOf(name)

	return kind, idx
}

// VarCount returns number of variable of the given kind
func (sb *SymbolTable) VarCount(kind VariableKind) uint {
	v, ok := sb.count[kind]
	if !ok {
		return 0
	}

	return v + 1
}

func (sb *SymbolTable) IsExists(name string) bool {
	_, ok := sb.m[name]
	return ok
}

func (sb *SymbolTable) Print() {
	fmt.Printf("%-10s", "name")
	fmt.Printf("%-10s", "type")
	fmt.Printf("%-10s", "kind")
	fmt.Printf("%-10s\n", "index")
	fmt.Println("-----------------------------------")

	for _, v := range sb.m {
		vType := v.sType.String()
		if v.sType == USR {
			vType = v.uType
		}
		fmt.Printf("%-10s", v.name)
		fmt.Printf("%-10s", vType)
		fmt.Printf("%-10v", v.kind)
		fmt.Printf("%-10d\n", v.index)
	}

	fmt.Println()
	fmt.Println()
}

func (v VariableKind) String() string {
	switch v {
	case Static:
		return "static"
	case Field:
		return "field"
	case Var:
		return "local"
	case Arg:
		return "argument"
	}

	return ""
}

func WhichKind(str string) VariableKind {
	if str == "field" {
		return Field
	}

	if str == "static" {
		return Static
	}

	if str == "var" {
		return Var
	}

	if str == "argument" {
		return Arg
	}

	return Undefined
}
