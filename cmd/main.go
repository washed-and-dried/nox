package main

import (
	"fmt"
	"nox/internals/parser"
	"os"
)

func main() {
	inputStr, err := os.ReadFile("main.nox")
	if err != nil {
		panic(err)
	}

	p := parser.NewParserFromString(string(inputStr))
    fn := p.Parse_func_def()
    fmt.Printf("%v", fn)
}
