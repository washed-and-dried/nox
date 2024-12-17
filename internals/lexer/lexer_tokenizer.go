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
		tok = token.Token{Literal: string(l.ch), Type: token.OPEN_PARAN, Pos: l.pos}
	case ')':
		tok = token.Token{Literal: string(l.ch), Type: token.CLOSE_PARAN, Pos: l.pos}
	case '{':
		tok = token.Token{Literal: string(l.ch), Type: token.OPEN_CURLY, Pos: l.pos}
	case '}':
		tok = token.Token{Literal: string(l.ch), Type: token.CLOSE_CURLY, Pos: l.pos}
	case '=':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string("=="), Type: token.BIN_EQUAL, Pos: l.pos}
				l.read_char()
			} else {
				tok = token.Token{Literal: string(l.ch), Type: token.ASSIGN, Pos: l.pos}
			}
		}
	case ';':
		tok = token.Token{Literal: string(l.ch), Type: token.SEMICOLON, Pos: l.pos}
	case ':':
		tok = token.Token{Literal: string(l.ch), Type: token.COLON, Pos: l.pos}
	case NULL_CHAR:
		tok = token.Token{Literal: "", Type: token.EOF, Pos: l.pos}

	case '"', '`', '\'':
		tok = l.read_string()

	case '+':
		{
			// TODO: handle prefix ++ in here
			tok = token.Token{Literal: string(l.ch), Type: token.BIN_PLUS, Pos: l.pos}
		}
	case '-':
		tok = token.Token{Literal: string(l.ch), Type: token.BIN_MINUS, Pos: l.pos}
	case '*':
		tok = token.Token{Literal: string(l.ch), Type: token.BIN_ASTERIC, Pos: l.pos}
	case '/':
		tok = token.Token{Literal: string(l.ch), Type: token.BIN_DIVIDE, Pos: l.pos}
	case '%':
		tok = token.Token{Literal: string(l.ch), Type: token.BIN_MODULO, Pos: l.pos}

	case '<':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string("<="), Type: token.BIN_LESS_THAN_EQUAL, Pos: l.pos}
				l.read_char()
			} else {
				tok = token.Token{Literal: string(l.ch), Type: token.BIN_LESS_THAN, Pos: l.pos}
			}
		}

	case '>':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string(">="), Type: token.BIN_GREATER_THAN_EQUAL, Pos: l.pos}
				l.read_char()
			} else {
				tok = token.Token{Literal: string(l.ch), Type: token.BIN_GREATER_THAN, Pos: l.pos}
			}
		}

	case '&':
		{
			if c, _ := l.peek_next_char(); c == '&' {
				tok = token.Token{Literal: string("&&"), Type: token.BIN_AND, Pos: l.pos}
				l.read_char()
			} else {
				tok = token.Token{Literal: string(l.ch), Type: token.BIN_BITWISE_AND, Pos: l.pos}
			}
		}

	case '|':
		{
			if c, _ := l.peek_next_char(); c == '|' {
				tok = token.Token{Literal: string("||"), Type: token.BIN_OR, Pos: l.pos}
				l.read_char()
			} else {
				tok = token.Token{Literal: string(l.ch), Type: token.BIN_BITWISE_OR, Pos: l.pos}
			}
		}

	case '!':
		{
			if c, _ := l.peek_next_char(); c == '=' {
				tok = token.Token{Literal: string("!="), Type: token.BIN_NOT_EQUAL, Pos: l.pos}
				l.read_char()
			} else {
				tok = token.Token{Literal: string(l.ch), Type: token.BIN_NOT, Pos: l.pos}
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
		} else if len(lit) >= MAX_STRING_SIZE {
			panic("len(str) > 256 or string wasn't closed | bro really thought this was javascript")
		} else {
			lit += string(l.ch)
		}
	}

	return token.Token{Literal: lit, Type: token.STR, Pos: pos}
}

func (l *Lexer) read_number() token.Token {
	pos := l.pos
	lit := ""

	for unicode.IsDigit(l.ch) {
		lit += string(l.ch)
		l.read_char()
	}

	return token.Token{Literal: lit, Type: token.INT, Pos: pos}
}

func (l *Lexer) read_ident_or_keyword() token.Token {
	pos := l.pos
	lit := ""

	for unicode.IsLetter(l.ch) {
		lit += string(l.ch)
		l.read_char()
	}

	if tokType, ok := token.IsKeyword(lit); ok {
		return token.Token{Literal: lit, Type: tokType, Pos: pos}
	} else {
		return token.Token{Literal: lit, Type: token.IDENT, Pos: pos}
	}
}

func (l *Lexer) skip_bloat_chars() {
	for unicode.IsSpace(l.ch) {
		l.read_char()
	}
}
