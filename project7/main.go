package main

import (
	"fmt"
	"os"
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
}

var arithCmds = [9]string{
	"add", "sub", "neq",
	"eq", "gt", "lt",
	"and", "or", "not",
}

func hasComment(line string) bool {
	return strings.Contains(line, "//")
}

func removeComment(line string) string {
	offset := strings.Index(line, "//")
	return strings.TrimSpace(line[0:offset])
}

func printErr(err string) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}
