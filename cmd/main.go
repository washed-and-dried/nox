package main

import (
	"nox/internals/eval"
	"nox/internals/parser"
)

func main() {
	file := "main.nox"
	p := parser.NewParser(file)
	program := p.Parse_program()
	eval.Eval_program(program)
}
