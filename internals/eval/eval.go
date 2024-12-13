package eval

import (
	"fmt"
	"nox/internals/parser"
	"nox/internals/token"
)

var (
	EVAL_NULL_OBJ  = NullObj{}
	EVAL_ERROR_OBJ = ErrorObj{}
)

// FIXME: add struct Eval with context
func Eval_program(program *parser.Program) {
	for _, stmt := range program.Stmts {
		eval_ast(stmt)
	}
}

func eval_ast(stmt parser.Statement) EvalObj {
	switch st := stmt.(type) {
	case parser.FuncDefStmt:
		{
			return Eval_func_def(st)
		}
	case parser.ExpressionStmt:
		{
			return eval_expr(st)
		}
	default:
		panic("Unhandled ast type: " + stmt.String())
	}
}

func Eval_func_def(fn parser.FuncDefStmt) EvalObj {
	// TODO: handle params, don't return anything just add the function definition onto the context
	for _, stmt := range fn.Body.Stmts {
		eval_ast(stmt)
	}

	return nil // FIXME: fix this shit nigger
}

func eval_expr(expr parser.ExpressionStmt) EvalObj { // FIXME: add proper return system using structs prolly
	switch expr.Type {
	case parser.EXPR_TYPE_FUNC:
		{
			return eval_func_call(expr.Value.AsFuncCall)
		}
	case parser.EXPR_TYPE_INT:
		{
			return IntObj{
				Value: int(eval_bin_expr(expr.Value.AsBinOp)), // FIXME: fix this shit
			}
		}
	default:
		{
			panic("You fucked up, unhandled expression type: " + expr.Type)
		}
	}
}

func eval_func_call(fnc parser.FuncCallExpr) EvalObj {
	// FIXME: replace with fetching the function definition, extending the context and thus evaluating the block statements
	expr := fnc.Args[0]
	value := eval_bin_expr(expr.Value.AsBinOp)

	if fnc.Ident.Name != "print" {
		panic("Something else popped up (other then ur cheerry heheheh): " + fnc.Ident.Name)
	}

	// call print
	fmt.Println(value)

	return EVAL_NULL_OBJ
}

func eval_bin_expr(bin_expr parser.BinaryExpr) int64 { // FIXME: proper return type, this ain't gonna cut it bro
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
