package lexer

import (
	"bufio"
	"io"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	file     *bufio.Reader // the file to read from
	ch       rune          // the current character
	pos      int           // current pos
	pos_next int           // next position to peek next char
	line_no  int           // current line number
}

func NewLexer(file *bufio.Reader) *Lexer {
	l := &Lexer{
		file: file,
		line_no: 0,
	}

	// sets ch, pos, pos_next
	l.read_char()

	return l
}

func NewLexerFromString(inputStr string) *Lexer {
	r := bufio.NewReader(strings.NewReader(inputStr))
	l := &Lexer{
		file: r,
	}

	// sets ch, pos, pos_next
	l.read_char()

	return l
}

const NULL_CHAR = '\x00'

func (l *Lexer) read_char() {
	char, eof := l.read_next_from_file()
	if eof {
		l.ch = NULL_CHAR
	} else {
		l.ch = char
	}

	if l.ch == '\n' {
		l.line_no++;
	}

	l.pos = l.pos_next
	l.pos_next++
}

func (l *Lexer) read_next_from_file() (rune, bool) {
	if c, _, err := l.file.ReadRune(); err != nil {
		if err == io.EOF {
			// returns ('\0', true) if we have reached EOF
			return NULL_CHAR, true
		} else {
			panic(err)
		}
	} else {
		// returns (next_char, false) if not EOF
		return c, false
	}
}

func (l *Lexer) peek_next_char() (rune, bool) {
	// Peek a sufficient number of bytes to decode at least one rune
	const maxPeekBytes = utf8.UTFMax // Maximum size of a UTF-8 encoded rune is 4 bytes
	buf, err := l.file.Peek(maxPeekBytes)
	if err != nil {
		if err == io.EOF {
			// returns ('\0', true) if we have reached EOF
			return NULL_CHAR, true
		} else {
			panic(err)
		}
	}

	// Decode the first rune from the buffer
	r, size := utf8.DecodeRune(buf)
	if r == utf8.RuneError && size == 1 {
		panic("Could not decode a single fking rune")
	}

	return r, false
}
