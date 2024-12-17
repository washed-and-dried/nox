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
    COLON

	IDENT
    RETURN
    LET
    FOR

	OPEN_PARAN
	CLOSE_PARAN
	OPEN_CURLY
	CLOSE_CURLY
    ASSIGN

	// types
	INT
    TYPE_INT
    STR
    TYPE_STR
	FUNC

	// operator
    BIN_OP_START

	BIN_PLUS
	BIN_MINUS
    BIN_ASTERIC
    BIN_DIVIDE
    BIN_MODULO

    BIN_LESS_THAN
    BIN_LESS_THAN_EQUAL
    BIN_GREATER_THAN
    BIN_GREATER_THAN_EQUAL
    BIN_EQUAL
    BIN_NOT_EQUAL

    BIN_AND
    BIN_OR
    BIN_NOT

    BIN_BITWISE_AND
    BIN_BITWISE_OR

    BIN_OP_END

)

func (tokType TokenType) String() string {
    switch tokType {
    case IDENT: return "ident"
    case RETURN: return "return"
    case LET: return "let"
    case FOR: return "for"

    case OPEN_PARAN: return "("
    case CLOSE_PARAN: return ")"
    case OPEN_CURLY: return "{"
    case CLOSE_CURLY: return "}"
    case ASSIGN: return "="

    // types
    case INT: return "int"
    case TYPE_INT: return "type_int"
    case STR: return "str"
    case TYPE_STR: return "type_str"
    case FUNC: return "fn"

    // operators
    case BIN_PLUS: return "+"
    case BIN_MINUS: return "-"
    case BIN_ASTERIC: return "*"
    case BIN_DIVIDE: return "/"
    case BIN_MODULO: return "%"

    case BIN_LESS_THAN   : return "<"
    case BIN_LESS_THAN_EQUAL: return "<="
    case BIN_GREATER_THAN: return ">"
    case BIN_GREATER_THAN_EQUAL: return ">="
    case BIN_EQUAL       : return "=="
    case BIN_NOT_EQUAL: return "!="

    case BIN_AND         : return "&&"
    case BIN_OR          : return "||"
    case BIN_NOT         : return "!"

    case BIN_BITWISE_AND : return "&"
    case BIN_BITWISE_OR  : return "|"

    case SEMICOLON: return ";"
    case COLON: return ":"
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
    case "let": return LET, true
    case "for": return FOR, true
    case "int": return TYPE_INT, true
    case "string": return TYPE_STR, true
    case "fn": return FUNC, true
    case "return": return RETURN, true
    default:
        return EOF, false
    }
}
