package main

import (
	"bufio"
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

	filepath := args[1]

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)

	p := parser.NewParser(r)
	program := p.Parse_program()

	eval.Eval_program(program)
}
