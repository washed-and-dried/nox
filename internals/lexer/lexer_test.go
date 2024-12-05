package lexer_test

import (
	"nox/internals/lexer"
	. "nox/internals/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
    1 + 1
    `
	toks := []Token{
		{Literal: "1", Type: INT},
		{Literal: "+", Type: BIN_PLUS},
		{Literal: "1", Type: INT},

		{Literal: "", Type: EOF},
	}

	test_tokens(t, input, toks)
}

func test_tokens(t *testing.T, input string, toks []Token) {
	l := lexer.NewLexerFromString(input)

	for _, eTok := range toks {
		aTok := l.NextToken()

		if aTok.Type != eTok.Type {
			t.Fatalf("expected Token Type [%s], got [%s]", eTok.Type.String(), aTok.Type.String())
		}

		if aTok.Literal != eTok.Literal {
			t.Fatalf("expected Token Literal [%s], got [%s]", eTok.Literal, aTok.Literal)
		}
	}
}
