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

func Eval_program(program *parser.Program) {
	ctx := EvalContext{
		objs: map[string]EvalObj{},
	}

	for _, stmt := range program.Stmts {
		eval_ast(stmt, &ctx)
	}
}

func eval_ast(stmt parser.Statement, ctx *EvalContext) EvalObj {
	switch st := stmt.(type) {
	case parser.FuncDefStmt:
		{
			return eval_func_def(st, ctx)
		}
	case parser.ExpressionStmt:
		{
			return eval_expr(st, ctx)
		}
	case parser.AssignStmt:
		{
			ctx.objs[st.Ident] = eval_ast(st.Value, ctx)
			return EVAL_NULL_OBJ
		}
	case parser.Identifier:
		{
			return eval_identifier(st, ctx)
		}
	case parser.ReturnStmt:
		{
			if st.Void { // if the return statement is empty "return;"
				return ReturnObj{
					Value: EVAL_NULL_OBJ,
				}
			}

			return ReturnObj{
				Value: eval_ast(st.ExprStmt, ctx),
			}
		}
	case parser.BodyStatement:
		{
			return eval_block_stmts(st, ctx)
		}
	default:
		panic("Unhandled ast type: " + stmt.String())
	}
}

func eval_identifier(ident parser.Identifier, ctx *EvalContext) EvalObj {
	if val, ok := builtins[ident.Name]; ok {
		return val
	}

	return ctx.Get(ident.Name)
}

func eval_func_def(fn parser.FuncDefStmt, ctx *EvalContext) EvalObj {
	// TODO: handle params, don't return anything just add the function definition onto the context
	fn_def := FuncDefObj{
		Ident: fn.Ident,
		Body:  *fn.Body,
		// FIXME: handle params
		// FIXME: add ctx if we want to go C like route. For now we use a global context!
	}

	ctx.objs[fn.Ident.Name] = fn_def // FIXME: encapsulate this shit into context only!

	return EVAL_NULL_OBJ
}

func eval_expr(expr parser.ExpressionStmt, ctx *EvalContext) EvalObj {
	switch expr.Type {
	case parser.EXPR_TYPE_FUNC:
		{
			fn_expr := expr.Value.AsFuncCall

			func_def := eval_ast(fn_expr.Ident, ctx)

			args := eval_expressions(&fn_expr.Args, ctx)

			return eval_func_call(func_def, args, ctx)
		}
	case parser.EXPR_TYPE_INT:
		{
			return IntObj{
				Value: int(expr.Value.AsInt.Value), // FIXME: fix this shit
			}
		}
	case parser.EXPR_TYPE_BIN:
		{
			return IntObj{
				Value: int(eval_bin_expr(expr.Value.AsBinOp)), // FIXME: fix this shit
			}
		}
	case parser.EXPR_TYPE_STR: // FIXME: idk about this, please check
		{
			return StrObj{
				Value: expr.Value.AsStr.Value,
			}
		}
	case parser.EXPR_TYPE_VAR:
		{
			uwu := eval_ast(expr.Value.AsVar, ctx)
			fmt.Println(uwu)
			return uwu // FIXME: fix this shit
		}
	default:
		{
			panic("You fucked up, unhandled expression type: " + expr.Type)
		}
	}
}

func eval_expressions(exprs *[]parser.ExpressionStmt, ctx *EvalContext) *[]EvalObj {
	res := []EvalObj{}

	for _, expr := range *exprs {
		res = append(res, eval_ast(expr, ctx))
	}

	return &res
}

func eval_func_call(fn_def EvalObj, args *[]EvalObj, ctx *EvalContext) EvalObj {
	switch fn := fn_def.(type) {
	case FuncDefObj:
		{
			extendedEnv := extendFunctionEnv(fn, args, ctx)
			evaluated := eval_ast(fn.Body, extendedEnv)

			return unwrapReturnValue(evaluated)
		}
	case BuiltinFuncObj:
		{
			return fn.fn(*args...)
		}
	default:
		panic("Not a function: " + fn_def.Type())
	}
}

func unwrapReturnValue(obj EvalObj) EvalObj {
	if ret, ok := obj.(ReturnObj); ok {
		return ret.Value
	}

	return obj
}

func eval_block_stmts(body parser.BodyStatement, ctx *EvalContext) EvalObj {
	// TODO: [DISCUSS] We would consider that if there are no return statements, the function will return NULL
	// var res EvalObj = EVAL_NULL_OBJ

	for _, stmt := range body.Stmts {
		res := eval_ast(stmt, ctx)

		if res.Type() == EVAL_RETURN {
			return res
		}
	}

	return EVAL_NULL_OBJ
}

func extendFunctionEnv(fn FuncDefObj, args *[]EvalObj, ctx *EvalContext) *EvalContext {
	extended := ctx.CreateNewEnclosedCtx()

	// FIXME: extend the context with params => args mapping

	return extended
}

func eval_bin_expr(bin_expr parser.BinaryExpr) int64 { // FIXME: proper return type, this ain't gonna cut it bro
    left := bin_expr.Left.Value.AsInt.Value // FIXME: fixmeplz
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
