package lexer

import (
	"nox/internals/token"
	"unicode"
)

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skip_bloat_chars()

	switch l.ch {
	case '(':
		tok = l.token_for(token.OPEN_PARAN)
	case ')':
		tok = l.token_for(token.CLOSE_PARAN)
	case '{':
		tok = l.token_for(token.OPEN_CURLY)
	case '}':
		tok = l.token_for(token.CLOSE_CURLY)
	case '[':
		tok = l.token_for(token.OPEN_SQUARE)
	case ']':
		tok = l.token_for(token.CLOSE_SQUARE)
	case '=':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string("=="), Type: token.BIN_EQUAL, Pos: l.pos, Line: l.line_no}
				l.read_char()
			} else {
				tok = l.token_for(token.ASSIGN)
			}
		}
	case ';':
		tok = l.token_for(token.SEMICOLON)
	case ':':
		tok = l.token_for(token.COLON)
	case NULL_CHAR:
		tok = token.Token{Literal: "", Type: token.EOF, Pos: l.pos, Line: l.line_no}

	case '"', '`', '\'':
		tok = l.read_string()

	case '+':
		{
			// TODO: handle prefix ++ in here
			tok = l.token_for(token.BIN_PLUS)
		}
	case '-':
		tok = l.token_for(token.BIN_MINUS)
	case '*':
		tok = l.token_for(token.BIN_ASTERIC)
	case '/':
		{
			if c, _ := l.peek_next_char(); c == '/' {
				l.skip_comment()
				return l.NextToken()
			} else {
				tok = l.token_for(token.BIN_DIVIDE)
			}
		}
	case '%':
		tok = l.token_for(token.BIN_MODULO)

	case '<':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string("<="), Type: token.BIN_LESS_THAN_EQUAL, Pos: l.pos, Line: l.line_no}
				l.read_char()
			} else {
				tok = l.token_for(token.BIN_LESS_THAN)
			}
		}

	case '>':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string(">="), Type: token.BIN_GREATER_THAN_EQUAL, Pos: l.pos, Line: l.line_no}
				l.read_char()
			} else {
				tok = l.token_for(token.BIN_GREATER_THAN)
			}
		}

	case '&':
		{
			if c, _ := l.peek_next_char(); c == '&' {
				tok = token.Token{Literal: string("&&"), Type: token.BIN_AND, Pos: l.pos, Line: l.line_no}
				l.read_char()
			} else {
				tok = l.token_for(token.BIN_BITWISE_AND)
			}
		}

	case '|':
		{
			if c, _ := l.peek_next_char(); c == '|' {
				tok = token.Token{Literal: string("||"), Type: token.BIN_OR, Pos: l.pos, Line: l.line_no}
				l.read_char()
			} else {
				tok = l.token_for(token.BIN_BITWISE_OR)
			}
		}

	case '!':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string("!="), Type: token.BIN_NOT_EQUAL, Pos: l.pos, Line: l.line_no}
				l.read_char()
			} else {
				tok = l.token_for(token.BIN_NOT)
			}
		}

	default:
		{
			if unicode.IsDigit(l.ch) { // TODO: handle floating point and hexadecimal numbers
				return l.read_number()
			} else {
				return l.read_ident_or_keyword()
			}
		}
	}

	l.read_char()

	return tok
}

func (l *Lexer) read_string() token.Token {
	strType := l.ch
	pos := l.pos
	lit := ""
	MAX_STRING_SIZE := 256 // FIXME: Move this outside?

	for len(lit) <= MAX_STRING_SIZE {
		l.read_char()
		if l.ch == '\n' {
			panic("newline char in regular str | your retarded ass doesn't understand common fucking strings")
		} else if l.ch == strType {
			break
		} else if l.ch == '\\' {
			l.read_char()
			switch l.ch {
			case 'n':
				lit += string('\n')
			case 't':
				lit += string('\t')
			case 'r':
				lit += string('\r')
			default:
				lit += string(l.ch)
			}
		} else if len(lit) >= MAX_STRING_SIZE {
			panic("len(str) > 256 or string wasn't closed | bro really thought this was javascript")
		} else {
			lit += string(l.ch)
		}
	}

	return token.Token{Literal: lit, Type: token.STR, Pos: pos, Line: l.line_no}
}

func (l *Lexer) read_number() token.Token {
	pos := l.pos
	lit := ""

	for unicode.IsDigit(l.ch) {
		lit += string(l.ch)
		l.read_char()
	}

	return token.Token{Literal: lit, Type: token.INT, Pos: pos, Line: l.line_no}
}

func (l *Lexer) read_ident_or_keyword() token.Token {
	pos := l.pos
	lit := ""

	for unicode.IsLetter(l.ch) {
		lit += string(l.ch)
		l.read_char()
	}

	if tokType, ok := token.IsKeyword(lit); ok {
		return token.Token{Literal: lit, Type: tokType, Pos: pos, Line: l.line_no}
	} else {
		return token.Token{Literal: lit, Type: token.IDENT, Pos: pos, Line: l.line_no}
	}
}

func (l *Lexer) skip_bloat_chars() {
	for unicode.IsSpace(l.ch) {
		l.read_char()
	}
}

func (l *Lexer) skip_comment() {
	l.read_char()
	l.read_char() // skip both the /

	for l.ch != '\n' {
		l.read_char()
	}
	l.read_char()
}

func (l *Lexer) token_for(Type token.TokenType) token.Token {
	return token.Token{Literal: string(l.ch), Type: Type, Pos: l.pos, Line: l.line_no}
}
