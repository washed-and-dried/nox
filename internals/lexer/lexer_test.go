package lexer_test

import (
	"nox/internals/lexer"
	. "nox/internals/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
    fn main() {
        int a = 1;
        print(1 + 1);
        print("Hello World!");
        return;
        return 1 % 2;
    }
    -*/%;
    `
	toks := []Token{
		{Literal: "fn", Type: FUNC},

		{Literal: "main", Type: IDENT},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: ")", Type: CLOSE_PARAN},

		{Literal: "{", Type: OPEN_CURLY},

		{Literal: "int", Type: TYPE_INT},
		{Literal: "a", Type: IDENT},
		{Literal: "=", Type: ASSIGN},
		{Literal: "1", Type: INT},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "print", Type: IDENT},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: "1", Type: INT},
		{Literal: "+", Type: BIN_PLUS},
		{Literal: "1", Type: INT},
		{Literal: ")", Type: CLOSE_PARAN},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "print", Type: IDENT},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: "Hello World!", Type: STR},
		{Literal: ")", Type: CLOSE_PARAN},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "return", Type: RETURN},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "return", Type: RETURN},
		{Literal: "1", Type: INT},
		{Literal: "%", Type: BIN_MODULO},
		{Literal: "2", Type: INT},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "}", Type: CLOSE_CURLY},

		{Literal: "-", Type: BIN_MINUS},
		{Literal: "*", Type: BIN_ASTERIC},
		{Literal: "/", Type: BIN_DIVIDE},
		{Literal: "%", Type: BIN_MODULO},
		{Literal: ";", Type: SEMICOLON},

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
