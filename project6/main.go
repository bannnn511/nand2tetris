package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Errors return by program.
var (
	ErrInvalidArguments    = errors.New("invalid number of arguments")
	ErrInvalidCInstruction = errors.New("parseCInstruction: invalid C instruction")
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
//  1. if instruction is @symbol; look up the symbol in the symbol table
//     if symbol value is found, use value to complete the instruction ’s translation
//     if not found:
//     add (symbol,n) to the symbol table
//     use n to complete the instruction translation
//     n++
//  2. if instruction is C-instruction, uses parseCInstruction
func parse(sb *SymbolTable, instruction string) string {
	if sb == nil {
		return ""
	}

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
	} else {
		dest, comp, jump, _ := parseCInstruction(instruction)
		machine = code(dest, comp, jump)
	}

	return machine
}

// parseCInstruction parses the Hack C-instruction
// symbolic: dest = comp; jump
// comp is mandatory
// if dest is empty; the = is ommited
// if jump is empty; the ; is ommited
//
// destIndx and jmpIdx is used to separate dest, comp and jump
func parseCInstruction(instruction string) (string, string, string, error) {
	var (
		dest    string
		comp    string
		jmp     string
		destIdx int
		jmpIdx  int
	)
	destIdx = strings.Index(instruction, "=")
	if destIdx == -1 {
		destIdx = 0
	} else {
		dest = instruction[0:destIdx]
		destIdx++
	}

	jmpIdx = strings.Index(instruction, ";")
	if jmpIdx == -1 {
		jmpIdx = len(instruction)
	} else {
		jmp = instruction[jmpIdx+1:]
	}

	comp = instruction[destIdx:jmpIdx]

	return dest, comp, jmp, nil
}

var compMap = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"M":   "1110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"!M":  "1110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"-M":  "1110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"M+1": "1110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"M-1": "1110010",
	"D+A": "0000010",
	"D+M": "1000010",
	"D-A": "0010011",
	"D-M": "1010011",
	"A-D": "0000111",
	"M-D": "1000111",
	"D&A": "0000000",
	"D&M": "1000000",
	"D|A": "0010101",
	"D|M": "1010101",
}

var destMap = map[string]string{
	"":    "000",
	"M":   "001",
	"D":   "010",
	"DM":  "011",
	"A":   "100",
	"AM":  "101",
	"AD":  "11O",
	"ADM": "111",
}

var jumpMap = map[string]string{
	"":    "000",
	"JGT": "001",
	"JEQ": "010",
	"JGE": "011",
	"JLT": "100",
	"JNE": "101",
	"JLE": "110",
	"JMP": "111",
}

func code(dest, comp, jmp string) string {
	return "111" + compMap[comp] + destMap[dest] + jumpMap[jmp]
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
