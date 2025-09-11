package parser

import (
	"fmt"
	"strconv"

	"nox/internals/token"
)

// this file will contain the expression parser. The meats and potatoes that will be used basically everywhere.

type ExprType string

const (
	EXPR_TYPE_FUNC      = "func_call"
	EXPR_TYPE_BIN       = "bin_op"
	EXPR_TYPE_INT       = "int"
	EXPR_TYPE_STR       = "str"
	EXPR_TYPE_VAR       = "var"
	EXPR_TYPE_BOOL      = "bool"
	EXPR_TYPE_SUBSCRIPT = "subscript"
)

type ExprValue struct {
	AsStr       StrExpr
	AsVar       Identifier
	AsFuncCall  FuncCallExpr
	AsBinOp     BinaryExpr
	AsInt       IntExpr
	AsBool      BoolExpr
	AsSubscript SubscriptExpr
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
	case token.STR:
		{
			value := p.tok.Literal
			p.next_token()

			return ExpressionStmt{
				Type: EXPR_TYPE_STR,
				Value: ExprValue{
					AsStr: StrExpr{
						Value: value,
					},
				},
			}
		}
	case token.BOOL_TRUE, token.BOOL_FALSE:
		{
			value := true
			if p.tok.Type == token.BOOL_FALSE {
				value = false
			}
			p.next_token()

			return ExpressionStmt{
				Type: EXPR_TYPE_BOOL,
				Value: ExprValue{
					AsBool: BoolExpr{
						Value: value,
					},
				},
			}
		}
	case token.IDENT:
		{
			if p.expect_peek(token.OPEN_PARAN) { // function calls
				return p.parse_func_calls()
			} else if p.expect_peek(token.OPEN_SQUARE) { // subscript str[0] //FIXME: Make this into a statement since we can have mutability str[0] = 'b'
				return p.parse_subscript()
			} else {
				return ExpressionStmt{
					Type: EXPR_TYPE_VAR,
					Value: ExprValue{
						AsVar: Identifier{p.expect_token_type(token.IDENT).Literal},
					},
				}
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

func (p *Parser) prep_default_val(t token.TokenType) ExpressionStmt {
	switch t {
	case token.TYPE_INT:
		{
			var value int64

			return ExpressionStmt{
				Type: EXPR_TYPE_INT,
				Value: ExprValue{
					AsInt: IntExpr{
						Value: value,
					},
				},
			}
		}
	case token.TYPE_STR:
		{
			var value string

			return ExpressionStmt{
				Type: EXPR_TYPE_STR,
				Value: ExprValue{
					AsStr: StrExpr{
						Value: value,
					},
				},
			}
		}
	default:
		{
			panic("Invalid Type: " + t.String() + " in line: " + fmt.Sprint(p.tok.Line-1))
		}
	}
}

func (p *Parser) parse_func_calls() ExpressionStmt {
	name := p.expect_token_type(token.IDENT).Literal
	p.expect_token_type(token.OPEN_PARAN)

	// FIXME: this shit is fking hard coded, fix it!!!
	args := []ExpressionStmt{}
	if p.tok.Type != token.CLOSE_PARAN { // if there are args to the functions!
		args = []ExpressionStmt{p.parse_expr()} // TODO: handle args, assume just "1 + 1" for now
	}

	p.expect_token_type(token.CLOSE_PARAN)

	return ExpressionStmt{
		Type: EXPR_TYPE_FUNC,
		Value: ExprValue{
			AsFuncCall: FuncCallExpr{
				Ident: Identifier{Name: name},
				Args:  args,
			},
		},
	}
}

func (p *Parser) parse_subscript() ExpressionStmt {
	name := p.expect_token_type(token.IDENT).Literal

	p.expect_token_type(token.OPEN_SQUARE)
	index := p.parse_expr()
	p.expect_token_type(token.CLOSE_SQUARE)

	return ExpressionStmt{
		Type: EXPR_TYPE_SUBSCRIPT,
		Value: ExprValue{
			AsSubscript: SubscriptExpr{
				Ident: Identifier{Name: name},
				Index: &index,
			},
		},
	}
}

const MAX_PRECEDENCE = 3

func (p *Parser) get_op_precedence(tokType token.TokenType) int {
	switch tokType {
	case token.BIN_AND, token.BIN_OR, token.BIN_NOT:
		return 0
	case token.BIN_GREATER_THAN_EQUAL, token.BIN_LESS_THAN_EQUAL,
		token.BIN_LESS_THAN, token.BIN_GREATER_THAN,
		token.BIN_EQUAL, token.BIN_NOT_EQUAL:
		return 1
	case token.BIN_PLUS, token.BIN_MINUS:
		return 2
	case token.BIN_ASTERIC, token.BIN_DIVIDE, token.BIN_MODULO:
		return 3
	default:
		panic("UNHANDLED OPERATOR PRECEDENCE: " + tokType.String())
	}
}
