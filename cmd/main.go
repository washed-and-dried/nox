package main

import (
	"fmt"
	"nox/internals/eval"
	"nox/internals/parser"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Invalid number of arguments: %d\n", len(args))
		return
	}

	file := "examples/triangle.nox"
	p := parser.NewParser(file)
	program := p.Parse_program()

	eval.Eval_program(program)
}
