package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printErr("invalid number of arguments")
	}

	// Open file
	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)
	if err != nil {
		printErr(fmt.Sprintf("%s file not exists\n", os.Args[1]))
	}
	defer file.Close()

	// init symbol table
	table := NewSymbolTable()

	r := bufio.NewScanner(file)

	// First pass
	// remove comments and replace labels
	lines := make([]string, 0, 50)
	for r.Scan() {
		line := r.Text()
		for i, word := range line {
			if word == '/' {
				line = line[0:i]
				break
			}
		}
		if len(line) == 0 {
			continue
		}
		if line[0] == '(' {
			table.AddSymbol(line[1 : len(line)-1])
			break
		}
		lines = append(lines, strings.TrimSpace(line))
	}

	for _, line := range lines {
		fmt.Println(parse(table, line))
	}

	// for symbol, value := range table.m {
	// 	fmt.Println(symbol, value)
	// }
}

// parse parses Hack language to Hack machine instruction
//  1. if the instruction is @symbol; look up the symbol in the symbol table
//     if symbol value is found, use value to complete the instruction ’s translation
//     if not found:
//     add (symbol,n) to the symbol table
//     use n to complete the instruction translation
//     n++
func parse(sb *SymbolTable, instruction string) string {
	machine := "0"
	if instruction[0] == '@' {
		symbol := instruction[1:]
		value, err := strconv.Atoi(symbol)
		if err != nil {
			v, ok := sb.m[symbol]
			if !ok {
				v = sb.n
				sb.AddSymbol(symbol)
			}
			value = int(v)
		}
		machine = fmt.Sprintf("%s%015b", machine, value)
	}

	return machine
}

type SymbolTable struct {
	m map[string]uint32
	n uint32
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		n: 16,
		m: map[string]uint32{
			"R0":     0,
			"R1":     1,
			"R2":     2,
			"R3":     3,
			"R4":     4,
			"R5":     5,
			"R6":     6,
			"R7":     7,
			"R8":     8,
			"R9":     9,
			"R10":    10,
			"R11":    11,
			"R12":    12,
			"R13":    13,
			"R14":    14,
			"R15":    15,
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
			"SCREEN": 16384,
			"KBD":    24576,
		},
	}
}

func (sb *SymbolTable) AddSymbol(symb string) {
	_, ok := sb.m[symb]
	if ok {
		return
	}
	sb.m[symb] = sb.n
	sb.n++
}

func (sb *SymbolTable) Setn(n uint32) {
	sb.n = n
}

func printErr(err string) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}
