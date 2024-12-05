package token

type Token struct {
	Literal string
	Type    TokenType
	Pos     int
}

type TokenType int

// token types
const (
	EOF = iota
	SEMICOLON

	IDENT

	OPEN_PARAN
	CLOSE_PARAN
	OPEN_CURLY
	CLOSE_CURLY

	// types
	INT
	FUNC

	// operator
    BIN_OP_START
	BIN_PLUS
    BIN_OP_END
)

func (tokType TokenType) String() string {
    switch tokType {
    case IDENT: return "ident"

    case OPEN_PARAN: return "("
    case CLOSE_PARAN: return ")"
    case OPEN_CURLY: return "{"
    case CLOSE_CURLY: return "}"

    // types
    case INT: return "int"
    case FUNC: return "fn"

    // operators
    case BIN_PLUS: return "+"

    case SEMICOLON: return ";"
    case EOF: return "eof"
    default:
        panic("unhandled token type")
    }
}

func IsBinaryOperator(tokType TokenType) bool {
    return tokType > BIN_OP_START && tokType < BIN_OP_END
}

func IsKeyword(lit string) (TokenType, bool) {
    switch lit {
    case "fn": return FUNC, true
    default:
        return EOF, false
    }
}
