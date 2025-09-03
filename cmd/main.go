package main

import (
	"nox/internals/eval"
	"nox/internals/parser"
)

func main() {
	file := "examples/triangle.nox"
	p := parser.NewParser(file)
	program := p.Parse_program()
	eval.Eval_program(program)
}
