package lexer_test

import (
	"nox/internals/lexer"
	. "nox/internals/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
    fn main() {
        for (let i: int = 0; i == 10 && i < 5; i = i + 1) {}
        let a: int = 1;
        a/2;
        print(1 + 1);
        print("Hello World!");
        let str: string = "something";
        return;
        return 1 % 2;
    }
    -*/%;
    `
    toks := []Token{ // FIXME: handle !=, <=, >= & similar cases
		{Literal: "fn", Type: FUNC},

		{Literal: "main", Type: IDENT},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: ")", Type: CLOSE_PARAN},

		{Literal: "{", Type: OPEN_CURLY},

		{Literal: "for", Type: FOR},
		{Literal: "(", Type: OPEN_PARAN},
		{Literal: "let", Type: LET},
		{Literal: "i", Type: IDENT},
		{Literal: ":", Type: COLON},
		{Literal: "int", Type: TYPE_INT},
		{Literal: "=", Type: ASSIGN},
		{Literal: "0", Type: INT},
		{Literal: ";", Type: SEMICOLON},
		{Literal: "i", Type: IDENT},
		{Literal: "==", Type: BIN_EQUAL},
		{Literal: "10", Type: INT},

		{Literal: "&&", Type: BIN_AND},
		{Literal: "i", Type: IDENT},
		{Literal: "<", Type: BIN_LESS_THAN},
		{Literal: "5", Type: INT},

		{Literal: ";", Type: SEMICOLON},
		{Literal: "i", Type: IDENT},
		{Literal: "=", Type: ASSIGN},
		{Literal: "i", Type: IDENT},
		{Literal: "+", Type: BIN_PLUS},
		{Literal: "1", Type: INT},
		{Literal: ")", Type: CLOSE_PARAN},
		{Literal: "{", Type: OPEN_CURLY},
		{Literal: "}", Type: CLOSE_CURLY},

		{Literal: "let", Type: LET},
		{Literal: "a", Type: IDENT},
		{Literal: ":", Type: COLON},
		{Literal: "int", Type: TYPE_INT},
		{Literal: "=", Type: ASSIGN},
		{Literal: "1", Type: INT},
		{Literal: ";", Type: SEMICOLON},

		{Literal: "a", Type: IDENT},
		{Literal: "/", Type: BIN_DIVIDE},
		{Literal: "2", Type: INT},
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

		{Literal: "let", Type: LET},
		{Literal: "str", Type: IDENT},
		{Literal: ":", Type: COLON},
		{Literal: "string", Type: TYPE_STR},
		{Literal: "=", Type: ASSIGN},
		{Literal: "something", Type: STR},
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
