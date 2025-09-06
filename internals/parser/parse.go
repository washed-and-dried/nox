package parser

import (
	"bufio"
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

func NewParser(file *bufio.Reader) *Parser {
	p := &Parser{
		l: lexer.NewLexer(file),
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

	program := &Program{
		Stmts: stmts,
	}

	p.append_main_call(program)

	return program
}

func (p *Parser) parse_statement() Statement {
	switch p.tok.Type {
	case token.FUNC:
		return p.Parse_func_def()
	case token.RETURN:
		{
			p.expect_token_type(token.RETURN)

			if p.tok.Type == token.SEMICOLON {
				p.expect_token_type(token.SEMICOLON)
				return ReturnStmt{
					Void: true,
				}
			} else {
				expr_stmt := p.parse_expr()
				p.expect_token_type(token.SEMICOLON)
				return ReturnStmt{
					Void:     false,
					ExprStmt: expr_stmt,
				}
			}
		}
	case token.IDENT:
		{
			if p.peekTok.Type == token.OPEN_PARAN {
				expr_stmt := p.parse_expr()
				p.expect_token_type(token.SEMICOLON) // expr stmts must end with a semicolon
				return expr_stmt
			} else { // handle var updation: a = a + 1
				varName := p.expect_token_type(token.IDENT).Literal

				p.expect_token_type(token.ASSIGN)

				value := p.parse_expr()
				if p.tok.Type != token.CLOSE_PARAN {
					p.expect_token_type(token.SEMICOLON)
				}

				return VarUpdation{
					Var: Identifier{
						Name: varName,
					},
					Value: value,
				}
			}
		}
	case token.LET:
		{
			p.expect_token_type(token.LET)
			ident := p.expect_token_type(token.IDENT).Literal
			p.expect_token_type(token.COLON)
			tok := p.tok
			p.next_token()
			p.expect_token_type(token.ASSIGN)
			expr := p.parse_expr()
			p.expect_token_type(token.SEMICOLON)
			return AssignStmt{
				Type:  tok,
				Value: expr,
				Ident: ident,
			}
		}
	case token.FOR:
		return p.parse_for_stmt()
	case token.IF:
		return p.parse_if_stmt(true)
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

func (p *Parser) parse_for_stmt() ForStmt {
	p.expect_token_type(token.FOR)        // for
	p.expect_token_type(token.OPEN_PARAN) // (

	// FIXME: check for for statements such as: for(i; ;)
	init := p.parse_statement() // let i: int = 10;

	cond := p.parse_expr()
	p.expect_token_type(token.SEMICOLON) // semicolon after expression: i < 10;

	updation := p.parse_statement() // i = i + 1

	p.expect_token_type(token.CLOSE_PARAN) // )

	body := p.parse_body() // {...body...}

	return ForStmt{
		Init:     init,
		Cond:     cond,
		Updation: updation,
		Body:     body,
	}
}

func (p *Parser) parse_if_stmt(first bool) IfStmt {
	cond := ExpressionStmt{}
	// isElse := p.tok.Type == token.ELSE
	if first || p.tok.Type == token.IF {
		p.expect_token_type(token.IF)
		p.expect_token_type(token.OPEN_PARAN)
		cond = p.parse_expr()
		p.expect_token_type(token.CLOSE_PARAN)
	}

	body := p.parse_body()

	var elseStmt Statement
	if p.tok.Type == token.ELSE {
		// if isElse {
		//     panic(fmt.Sprintf("Redundant Else at %d", p.tok.Pos)) // TODO: better panic log
		// }
		p.next_token()
		if p.tok.Type == token.IF {
			elseStmt = p.parse_if_stmt(false)
		} else {
			elseStmt = p.parse_body()
		}
	}

	return IfStmt{
		Else: elseStmt,
		Cond: cond,
		Body: body,
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

func (p *Parser) append_main_call(program *Program) {
	main_call := ExpressionStmt{
		Value: ExprValue{
			AsFuncCall: FuncCallExpr{
				Ident: Identifier{
					Name: "main",
				},
				Args: []ExpressionStmt{},
			},
		},
		Type: EXPR_TYPE_FUNC,
	}

	program.Stmts = append(program.Stmts, main_call)
}
