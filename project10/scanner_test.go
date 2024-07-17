package main_test

import (
	"fmt"
	prj "project10"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	src := []byte("class Main {")
	var s prj.Scanner
	s.Init(src)

	tok, lit := s.Scan()
	fmt.Println(tok, lit)

	tok, lit = s.Scan()
	fmt.Println(tok, lit)

	tok, lit = s.Scan()
	fmt.Println(tok, lit)
}
