package parser

import (
	"nox/internals/lexer"
	"nox/internals/token"
)

// this will be the interface which other parts will interact with to use the parser.
// It will just parse the file into a single global AST node called a Program.

type Parser struct {
	l       *lexer.Lexer
	peekTok token.Token
	tok     token.Token
}

// TODO: func NewParser(filePath string)
func NewParserFromString(inputStr string) *Parser {
	p := &Parser{
		l: lexer.NewLexerFromString(inputStr),
	}

	p.next_token() // set peekToken, tok
	p.next_token()

	return p
}

func (p *Parser) Parse_func_def() FuncDefStmt {
	p.expect_token_type(token.FUNC)

	ident := p.expect_token_type(token.IDENT).Literal

	p.expect_token_type(token.OPEN_PARAN)
	// TODO: parse params
	p.expect_token_type(token.CLOSE_PARAN)

    body := p.parse_body()

	return FuncDefStmt{
		Ident: ident,
		Body:  &body,
	}
}

func (p *Parser) parse_body() BodyStatement {
	p.expect_token_type(token.OPEN_CURLY)

	stmts := []Statement{}

	// parse function call expression
	funcCall := p.parse_expr()
	stmts = append(stmts, funcCall)

	// FIXME: for now only parsing function call statement
	p.expect_token_type(token.SEMICOLON)
	p.expect_token_type(token.CLOSE_CURLY)

	return BodyStatement{
		Stmts: stmts,
	}
}

func (p *Parser) next_token() {
	p.tok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) expect_token_type(tokType token.TokenType) token.Token {
	if p.tok.Type != tokType {
		panic("Expected: " + tokType.String() + " got: " + p.tok.Type.String())
	}

	tok := p.tok
	p.next_token()

	return tok
}

func (p *Parser) expect_peek(tokType token.TokenType) bool { // returns if peek = expected tok, does not consume the current one
	return p.peekTok.Type == tokType
}
