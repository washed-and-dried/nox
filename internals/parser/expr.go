package parser

import (
	"nox/internals/token"
	"strconv"
)

// this file will contain the expression parser. The meats and potatoes that will be used basically everywhere.

type ExprType string

const (
	EXPR_TYPE_FUNC = "func_call"
	EXPR_TYPE_BIN  = "bin_op"
	EXPR_TYPE_INT  = "int"
)

type ExprValue struct {
	AsFuncCall FuncCallExpr
	AsBinOp    BinaryExpr
	AsInt      IntExpr
}

func (p *Parser) parse_expr() ExpressionStmt {
	return p.parse_bin_op(0)
}

func (p *Parser) parse_bin_op(precedence int) ExpressionStmt {
	if precedence > MAX_PRECEDENCE {
		// parse primary exprs
		return p.parse_primary_exprs()
	}

	left := p.parse_bin_op(precedence + 1)
	for p.tok.Type != token.EOF && token.IsBinaryOperator(p.tok.Type) && p.get_op_precedence(p.tok.Type) == precedence {
		operator := p.tok
		p.next_token() // consume the binary operator

		right := p.parse_bin_op(precedence + 1)

		bin_left := left // go shinanigans
		left = ExpressionStmt{
			Type: EXPR_TYPE_BIN,
			Value: ExprValue{
				AsBinOp: BinaryExpr{
					Right:    &right,
					Left:     &bin_left,
					Operator: operator,
				},
			},
		}
	}

	return left
}

func (p *Parser) parse_primary_exprs() ExpressionStmt {
	switch p.tok.Type {
	case token.INT:
		{
			value, err := strconv.ParseInt(p.tok.Literal, 10, 64)
			if err != nil {
				panic(err)
			}
			p.next_token()

			return ExpressionStmt{
				Type: EXPR_TYPE_INT,
				Value: ExprValue{
					AsInt: IntExpr{
						Value: value,
					},
				},
			}
		}
	case token.IDENT:
		{
			if p.expect_peek(token.OPEN_PARAN) { // function calls
				return p.parse_func_calls()
			} else {
				return ExpressionStmt{} // TODO: handle variables
			}
		}
	case token.OPEN_PARAN: // (..Expr..)
		{
			p.next_token()
			if p.expect_peek(token.CLOSE_PARAN) {
				return ExpressionStmt{}
			}

			expr := p.parse_expr()

			p.expect_token_type(token.CLOSE_PARAN)

			return expr
		}
	default: // FIXME: add guard with all non-primary tokens rather than panic
		panic("Not an primary expression: " + p.tok.Type.String())
	}
}

func (p *Parser) parse_func_calls() ExpressionStmt {
	name := p.expect_token_type(token.IDENT).Literal
	p.expect_token_type(token.OPEN_PARAN)

	args := []ExpressionStmt{p.parse_expr()} // TODO: handle args, assume just "1 + 1" for now

	p.expect_token_type(token.CLOSE_PARAN)

	return ExpressionStmt{
		Type: EXPR_TYPE_FUNC,
		Value: ExprValue{
			AsFuncCall: FuncCallExpr{
				Ident: name,
				Args:  args,
			},
		},
	}
}

const MAX_PRECEDENCE = 2

func (p *Parser) get_op_precedence(tokType token.TokenType) int {
	switch tokType {
	case token.BIN_PLUS, token.BIN_MINUS:
		return 1
	case token.BIN_ASTERIC, token.BIN_DIVIDE, token.BIN_MODULO:
		return 2
	default:
		panic("UNHANDLED OPERATOR PRECEDENCE: " + tokType.String())
	}
}
