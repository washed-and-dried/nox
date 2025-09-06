package lexer_test

import (
	"testing"

	"nox/internals/lexer"
	. "nox/internals/token"
)

func TestLexer(t *testing.T) {
	input := `
    fn main() {
        if (let i: int = 0; i != 10;) {} else {}
        for (let i: int = 0; i != 10 & i <= 5; i = i + 1) {}
        let a: int = 1;
        a/2;
        a = a + 1;
        true == false;
        a >= 5;
        print(1 + 1);
        print("Hello World!");
        let str: string = "something"; // comment hahahha
        return;
        return 1 % 2;
        str[0];
		// this should be ignored
    }
    -*/%;
    `
	toks := []Token{
		{Literal: "fn", Type: FUNC, Line: 1},

		{Literal: "main", Type: IDENT, Line: 1},
		{Literal: "(", Type: OPEN_PARAN, Line: 1},
		{Literal: ")", Type: CLOSE_PARAN, Line: 1},

		{Literal: "{", Type: OPEN_CURLY, Line: 1},

		{Literal: "if", Type: IF, Line: 2},
		{Literal: "(", Type: OPEN_PARAN, Line: 2},

		{Literal: "let", Type: LET, Line: 2},
		{Literal: "i", Type: IDENT, Line: 2},
		{Literal: ":", Type: COLON, Line: 2},
		{Literal: "int", Type: TYPE_INT, Line: 2},
		{Literal: "=", Type: ASSIGN, Line: 2},
		{Literal: "0", Type: INT, Line: 2},
		{Literal: ";", Type: SEMICOLON, Line: 2},
		{Literal: "i", Type: IDENT, Line: 2},
		{Literal: "!=", Type: BIN_NOT_EQUAL, Line: 2},
		{Literal: "10", Type: INT, Line: 2},
		{Literal: ";", Type: SEMICOLON, Line: 2},

		{Literal: ")", Type: CLOSE_PARAN, Line: 2},
		{Literal: "{", Type: OPEN_CURLY, Line: 2},
		{Literal: "}", Type: CLOSE_CURLY, Line: 2},

		{Literal: "else", Type: ELSE, Line: 2},
		{Literal: "{", Type: OPEN_CURLY, Line: 2},
		{Literal: "}", Type: CLOSE_CURLY, Line: 2},

		{Literal: "for", Type: FOR, Line: 3},
		{Literal: "(", Type: OPEN_PARAN, Line: 3},
		{Literal: "let", Type: LET, Line: 3},
		{Literal: "i", Type: IDENT, Line: 3},
		{Literal: ":", Type: COLON, Line: 3},
		{Literal: "int", Type: TYPE_INT, Line: 3},
		{Literal: "=", Type: ASSIGN, Line: 3},
		{Literal: "0", Type: INT, Line: 3},
		{Literal: ";", Type: SEMICOLON, Line: 3},
		{Literal: "i", Type: IDENT, Line: 3},
		{Literal: "!=", Type: BIN_NOT_EQUAL, Line: 3},
		{Literal: "10", Type: INT, Line: 3},

		{Literal: "&", Type: BIN_BITWISE_AND, Line: 3},
		{Literal: "i", Type: IDENT, Line: 3},
		{Literal: "<=", Type: BIN_LESS_THAN_EQUAL, Line: 3},
		{Literal: "5", Type: INT, Line: 3},

		{Literal: ";", Type: SEMICOLON, Line: 3},
		{Literal: "i", Type: IDENT, Line: 3},
		{Literal: "=", Type: ASSIGN, Line: 3},
		{Literal: "i", Type: IDENT, Line: 3},
		{Literal: "+", Type: BIN_PLUS, Line: 3},
		{Literal: "1", Type: INT, Line: 3},
		{Literal: ")", Type: CLOSE_PARAN, Line: 3},
		{Literal: "{", Type: OPEN_CURLY, Line: 3},
		{Literal: "}", Type: CLOSE_CURLY, Line: 3},

		{Literal: "let", Type: LET, Line: 4},
		{Literal: "a", Type: IDENT, Line: 4},
		{Literal: ":", Type: COLON, Line: 4},
		{Literal: "int", Type: TYPE_INT, Line: 4},
		{Literal: "=", Type: ASSIGN, Line: 4},
		{Literal: "1", Type: INT, Line: 4},
		{Literal: ";", Type: SEMICOLON, Line: 4},

		{Literal: "a", Type: IDENT, Line: 5},
		{Literal: "/", Type: BIN_DIVIDE, Line: 5},
		{Literal: "2", Type: INT, Line: 5},
		{Literal: ";", Type: SEMICOLON, Line: 5},

		{Literal: "a", Type: IDENT, Line: 6},
		{Literal: "=", Type: ASSIGN, Line: 6},
		{Literal: "a", Type: IDENT, Line: 6},
		{Literal: "+", Type: BIN_PLUS, Line: 6},
		{Literal: "1", Type: INT, Line: 6},
		{Literal: ";", Type: SEMICOLON, Line: 6},

		{Literal: "true", Type: BOOL_TRUE, Line: 7},
		{Literal: "==", Type: BIN_EQUAL, Line: 7},
		{Literal: "false", Type: BOOL_FALSE, Line: 7},
		{Literal: ";", Type: SEMICOLON, Line: 7},

		{Literal: "a", Type: IDENT, Line: 8},
		{Literal: ">=", Type: BIN_GREATER_THAN_EQUAL, Line: 8},
		{Literal: "5", Type: INT, Line: 8},
		{Literal: ";", Type: SEMICOLON, Line: 8},

		{Literal: "print", Type: IDENT, Line: 9},
		{Literal: "(", Type: OPEN_PARAN, Line: 9},
		{Literal: "1", Type: INT, Line: 9},
		{Literal: "+", Type: BIN_PLUS, Line: 9},
		{Literal: "1", Type: INT, Line: 9},
		{Literal: ")", Type: CLOSE_PARAN, Line: 9},
		{Literal: ";", Type: SEMICOLON, Line: 9},

		{Literal: "print", Type: IDENT, Line: 10},
		{Literal: "(", Type: OPEN_PARAN, Line: 10},
		{Literal: "Hello World!", Type: STR, Line: 10},
		{Literal: ")", Type: CLOSE_PARAN, Line: 10},
		{Literal: ";", Type: SEMICOLON, Line: 10},

		{Literal: "let", Type: LET, Line: 11},
		{Literal: "str", Type: IDENT, Line: 11},
		{Literal: ":", Type: COLON, Line: 11},
		{Literal: "string", Type: TYPE_STR, Line: 11},
		{Literal: "=", Type: ASSIGN, Line: 11},
		{Literal: "something", Type: STR, Line: 11},
		{Literal: ";", Type: SEMICOLON, Line: 11},

		{Literal: "return", Type: RETURN, Line: 12},
		{Literal: ";", Type: SEMICOLON, Line: 12},

		{Literal: "return", Type: RETURN, Line: 13},
		{Literal: "1", Type: INT, Line: 13},
		{Literal: "%", Type: BIN_MODULO, Line: 13},
		{Literal: "2", Type: INT, Line: 13},
		{Literal: ";", Type: SEMICOLON, Line: 13},

		{Literal: "str", Type: IDENT, Line: 14},
		{Literal: "[", Type: OPEN_SQUARE, Line: 14},
		{Literal: "0", Type: INT, Line: 14},
		{Literal: "]", Type: CLOSE_SQUARE, Line: 14},
		{Literal: ";", Type: SEMICOLON, Line: 14},

		// Comment line supposed to be here, ignored by lexer but still is a line in the file, therefore Line_no 15 is skipped

		{Literal: "}", Type: CLOSE_CURLY, Line: 16},

		{Literal: "-", Type: BIN_MINUS, Line: 17},
		{Literal: "*", Type: BIN_ASTERIC, Line: 17},
		{Literal: "/", Type: BIN_DIVIDE, Line: 17},
		{Literal: "%", Type: BIN_MODULO, Line: 17},
		{Literal: ";", Type: SEMICOLON, Line: 17},

		{Literal: "", Type: EOF, Line: 18}, // eof is in a new line, for whatever reason, no matter to us
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

		if eTok.Line != 0 && eTok.Line != aTok.Line {
			t.Fatalf("expected Line Number [%d], got [%d]", eTok.Line, aTok.Line)
		}

	}
}
