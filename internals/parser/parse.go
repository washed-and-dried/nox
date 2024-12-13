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

func NewParser(filePath string) *Parser {
	p := &Parser{
		l: lexer.NewLexer(filePath),
	}

	p.next_token() // set peekToken, tok
	p.next_token()

	return p
}

func NewParserFromString(inputStr string) *Parser {
	p := &Parser{
		l: lexer.NewLexerFromString(inputStr),
	}

	p.next_token() // set peekToken, tok
	p.next_token()

	return p
}

func (p *Parser) Parse_program() *Program {
	stmts := []Statement{}

	for p.tok.Type != token.EOF {
		stmts = append(stmts, p.parse_statement())
		// if p.tok.Type == token.SEMICOLON { // FIXME: remove if not necessary
		// 	p.next_token()
		// }
	}

	return &Program{
		Stmts: stmts,
	}
}

func (p *Parser) parse_statement() Statement {
	switch p.tok.Type {
	case token.FUNC:
		return p.Parse_func_def()
	case token.IDENT:
		{
			if p.peekTok.Type == token.OPEN_PARAN {
				expr_stmt := p.parse_expr()
				p.expect_token_type(token.SEMICOLON) // expr stmts must end with a semicolon
				return expr_stmt
			} else {
				return nil // TODO: handle variable assigments statements
			}
		}
	default:
		panic("Unhandled statement type: " + p.tok.Type.String())
	}
}

func (p *Parser) Parse_func_def() FuncDefStmt {
	p.expect_token_type(token.FUNC)

	ident := p.expect_token_type(token.IDENT).Literal

	p.expect_token_type(token.OPEN_PARAN)
	// TODO: parse params
	p.expect_token_type(token.CLOSE_PARAN)

	body := p.parse_body()

	return FuncDefStmt{
		Ident: Identifier{Name: ident},
		Body:  &body,
	}
}

func (p *Parser) parse_body() BodyStatement {
	p.expect_token_type(token.OPEN_CURLY)

	stmts := []Statement{}

	for p.tok.Type != token.CLOSE_CURLY {
		stmts = append(stmts, p.parse_statement())
	}

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
