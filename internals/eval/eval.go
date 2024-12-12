package eval

import (
	"fmt"
	"nox/internals/parser"
	"nox/internals/token"
)

// FIXME: add struct Eval with context
func Eval_program(program *parser.Program) {
	for _, stmt := range program.Stmts {
		eval_stmt(stmt)
	}
}

func eval_stmt(stmt parser.Statement) {
	switch st := stmt.(type) {
	case parser.FuncDefStmt:
		{
			Eval_func_def(st)
		}
	case parser.ExpressionStmt:
		{
			eval_expr(st)
		}
	default:
		panic("Unhandled type: " + stmt.String())
	}
}

func Eval_func_def(fn parser.FuncDefStmt) /*TODO:*/ {
	// fn_call := fn.Body.Stmts[0].(parser.ExpressionStmt) // FIXME: remove hard coding
	//
	// if fn_call.Type != parser.EXPR_FUNC_CALL {
	// 	panic("Expected function call got: " + fn_call.Type)
	// }
	// eval_func_call(fn_call.Value.AsFuncCall)

	for _, stmt := range fn.Body.Stmts {
		eval_stmt(stmt)
	}
}

func eval_expr(expr parser.ExpressionStmt) int64 { // FIXME: add proper return system using structs prolly
	switch expr.Type {
	case parser.EXPR_TYPE_FUNC:
		{
			return eval_func_call(expr.Value.AsFuncCall)
		}
	case parser.EXPR_TYPE_INT:
		{
			return eval_bin_expr(expr.Value.AsBinOp)
		}
	default:
		{
			panic("You fucked up, unhandled expression type: " + expr.Type)
		}
	}
}

func eval_func_call(fnc parser.FuncCallExpr) int64 {
	expr := fnc.Args[0]
	value := eval_bin_expr(expr.Value.AsBinOp)

	if fnc.Ident != "print" {
		panic("Something else popped up (other then ur cheerry heheheh): " + fnc.Ident)
	}

	// call print
	fmt.Println(value)

	return 0
}

func eval_bin_expr(bin_expr parser.BinaryExpr) int64 {
	left := bin_expr.Left.Value.AsInt.Value
	right := bin_expr.Right.Value.AsInt.Value

	switch bin_expr.Operator.Type {
	case token.BIN_PLUS:
		return left + right
	case token.BIN_MINUS:
		return left - right
	case token.BIN_ASTERIC:
		return left * right
	case token.BIN_DIVIDE:
		return left / right
	case token.BIN_MODULO:
		return left % right
	default:
		panic("Unhandled operator")
	}
}
