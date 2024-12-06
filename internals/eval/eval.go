package eval

import (
	"fmt"
	"nox/internals/parser"
	"nox/internals/token"
)

func Eval_func_def(fn parser.FuncDefStmt) /*TODO:*/ {
	fn_call := fn.Body.Stmts[0].(parser.ExpressionStmt) // FIXME: remove hard coding

	if fn_call.Type != parser.EXPR_FUNC_CALL {
		panic("Expected function call got: " + fn_call.Type)
	}
	eval_func_call(fn_call.Value.AsFuncCall)
}

func eval_func_call(fnc parser.FuncCallExpr) {
	expr := fnc.Args[0]
	value := parse_bin_expr(expr)

	if fnc.Ident != "print" {
		panic("Something else popped up (other then ur cheerry heheheh): " + fnc.Ident)
	}

	// call print
	fmt.Println(value)
}

func parse_bin_expr(expr parser.ExpressionStmt) int64 {
	if expr.Type != parser.EXPR_TYPE_BIN {
		panic("Expected binary expression, something else popped up: " + expr.Type)
	}

	bin_expr := expr.Value.AsBinOp

	switch bin_expr.Operator.Type {
	case token.BIN_PLUS:
		{
			l := bin_expr.Left.Value.AsInt.Value
			r := bin_expr.Right.Value.AsInt.Value

			return l + r
		}
	default:
		panic("Unhandled operator")
	}
}
