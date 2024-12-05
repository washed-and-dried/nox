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

	IDENT

	OPEN_PARAN
	CLOSE_PARAN
	OPEN_CURLY
	CLOSE_CURLY

	// types
	INT
	FUNC

	// operator
	BIN_PLUS
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

    case EOF: return "eof"
    default:
        panic("unhandled token type")
    }
}

func IsKeyword(lit string) (TokenType, bool) {
    switch lit {
    case "fn": return FUNC, true
    default:
        return EOF, false
    }
}
