package lexer_test

import (
	"nox/internals/lexer"
	. "nox/internals/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
    fn main() {
        print(1 + 1);
    }
    `
	toks := []Token{
		{Literal: "fn", Type: FUNC},

		{Literal: "main", Type: IDENT},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: ")", Type: CLOSE_PARAN},

		{Literal: "{", Type: OPEN_CURLY},

		{Literal: "print", Type: IDENT},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: "1", Type: INT},
		{Literal: "+", Type: BIN_PLUS},
		{Literal: "1", Type: INT},
		{Literal: ")", Type: CLOSE_PARAN},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "}", Type: CLOSE_CURLY},

		{Literal: "", Type: EOF},
	}

	test_tokens(t, input, toks)
}

func test_tokens(t *testing.T, input string, toks []Token) {
	l := lexer.NewLexerFromString(input)

	for _, eTok := range toks {
		aTok := l.NextToken()

		if aTok.Literal != eTok.Literal {
			t.Fatalf("expected Token Literal [%s], got [%s]", eTok.Literal, aTok.Literal)
		}

		if aTok.Type != eTok.Type {
			t.Fatalf("expected Token Type [%s], got [%s]", eTok.Type.String(), aTok.Type.String())
		}
	}
}
