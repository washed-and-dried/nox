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
			val := eval_ast(st.Value, ctx)

			check_type_compliant(val, st.Type.Type) // check if type and value aligns

			ctx.objs[st.Ident] = val
			return EVAL_NULL_OBJ
		}
	case parser.VarUpdation:
		{
			updatedValue := eval_ast(st.Value, ctx)

			if ctx.Get(st.Var.Name) == EVAL_NULL_OBJ {
				panic("No variable named: " + st.Var.Name)
			}

			ctx.objs[st.Var.Name] = updatedValue
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
	case parser.ForStmt:
		{
			return eval_for_stmt(st, ctx.CreateNewEnclosedCtx())
		}
	case parser.IfStmt:
		{
			return eval_if_stmt(st, ctx.CreateNewEnclosedCtx())
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
			return eval_bin_expr(expr.Value.AsBinOp, ctx)
		}
	case parser.EXPR_TYPE_STR: // FIXME: idk about this, please check
		{
			return StrObj{
				Value: expr.Value.AsStr.Value,
			}
		}
	case parser.EXPR_TYPE_BOOL:
		{
			return BoolObj{
				Value: expr.Value.AsBool.Value,
			}
		}
	case parser.EXPR_TYPE_VAR:
		{
			return eval_ast(expr.Value.AsVar, ctx)
		}
	case parser.EXPR_TYPE_SUBSCRIPT:
		{
			sbs_expr := expr.Value.AsSubscript

			ident := eval_ast(sbs_expr.Ident, ctx)
            index := eval_ast(*sbs_expr.Index, ctx)

            // FIXME: for now we will just return an StrObj because I am too lazy to implement print and equals for CharObj

            if ident.Type() != EVAL_STR && index.Type() != EVAL_INT {
                return EVAL_ERROR_OBJ
            }

            strObj, _ := ident.(StrObj)
            intObj, _ := index.(IntObj)

            return StrObj{
                Value: string(strObj.Value[intObj.Value]),
            }
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

func eval_for_stmt(stmt parser.ForStmt, ctx *EvalContext) EvalObj {
	eval_ast(stmt.Init, ctx) // introduce the loop variable into the context

	for ifBoolTrue(eval_ast(stmt.Cond, ctx)) {
		if ret := eval_block_stmts(stmt.Body, ctx); ret.Type() == EVAL_RETURN {
			return ret
		}

		eval_ast(stmt.Updation, ctx) // update the variable or other shit
	}

	for key := range ctx.objs {
		if _, ok := ctx.outer.objs[key]; ok {
			ctx.outer.objs[key] = ctx.objs[key]
		}
	}

	return EVAL_NULL_OBJ
}

func eval_if_stmt(stmt parser.Statement, ctx *EvalContext) EvalObj {
	switch stmt := stmt.(type) {
	case parser.IfStmt:
		{
			if ifBoolTrue(eval_ast(stmt.Cond, ctx)) {
				return eval_block_stmts(stmt.Body, ctx)
			} else {
				return eval_if_stmt(stmt.Else, ctx)
			}
		}
	case parser.BodyStatement:
		{
			return eval_block_stmts(stmt, ctx)
		}
	}

	return EVAL_NULL_OBJ
}

func ifBoolTrue(boolObj EvalObj) bool {
	obj, ok := boolObj.(BoolObj)

	if !ok {
		panic("Value is not a bool value: " + boolObj.Type())
	}

	return obj.Value
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

func eval_bin_expr(bin_expr parser.BinaryExpr, ctx *EvalContext) EvalObj {
    if bin_expr.Left == nil || bin_expr.Right == nil {
        fmt.Println("Left or right BinaryExpr was nil")
        return EVAL_NULL_OBJ
    }

	left := eval_ast(*bin_expr.Left, ctx) // NOTE: we would have to dereference both left and right since they are pointers due to recursive types
	right := eval_ast(*bin_expr.Right, ctx)

	if left.Type() == EVAL_INT && right.Type() == EVAL_INT {
		return perform_bin_operation_int(left, right, bin_expr.Operator)
	} else if left.Type() == EVAL_BOOL && right.Type() == EVAL_BOOL {
		return perform_bin_operation_bool(left, right, bin_expr.Operator)
	} else if left.Type() == EVAL_STR && right.Type() == EVAL_STR {
		return perform_bin_operation_str(left, right, bin_expr.Operator)
	} else {
		// FIXME: for now only allow operations between ints
		panic("Illegal operation between: " + left.Type() + " and " + right.Type())
	}
}

func perform_bin_operation_int(left EvalObj, right EvalObj, operator token.Token) EvalObj {
	boolObj := func(b bool) BoolObj { // FIXME: only temporary for satisfaction
		return BoolObj{Value: b}
	}

	// be careful with the types since we might introduce 32 and 64 bit ints separately

	lval := left.(IntObj).Value
	rval := right.(IntObj).Value
	var res int

	switch operator.Type {
	case token.BIN_PLUS:
		res = lval + rval
	case token.BIN_MINUS:
		res = lval - rval
	case token.BIN_ASTERIC:
		res = lval * rval
	case token.BIN_DIVIDE:
		res = lval / rval
	case token.BIN_MODULO:
		res = lval % rval
	case token.BIN_EQUAL: // FIXME: organise this shit, make it generic
		return boolObj(lval == rval)
	case token.BIN_NOT_EQUAL:
		return boolObj(lval != rval)
	case token.BIN_GREATER_THAN:
		return boolObj(lval > rval)
	case token.BIN_LESS_THAN:
		return boolObj(lval < rval)
	case token.BIN_GREATER_THAN_EQUAL:
		return boolObj(lval >= rval)
	case token.BIN_LESS_THAN_EQUAL:
		return boolObj(lval <= rval)
	default:
		panic("Unhandled operator")
	}

	return IntObj{
		Value: res,
	}
}

func perform_bin_operation_bool(left EvalObj, right EvalObj, operator token.Token) EvalObj {
	var res bool
	lval := left.(BoolObj).Value
	rval := right.(BoolObj).Value

	switch operator.Type {
	case token.BIN_AND:
		res = lval && rval
	case token.BIN_OR:
		res = lval || rval
	case token.BIN_NOT: // FIXME: handle this once we implement prefix operators
	default:
		panic("Unhandled operator: " + operator.Type.String())
	}

	return BoolObj{
		Value: res,
	}
}
func perform_bin_operation_str(left EvalObj, right EvalObj, operator token.Token) EvalObj {
    var res bool
	lval := left.(StrObj).Value
	rval := right.(StrObj).Value

	switch operator.Type {
	case token.BIN_EQUAL:
		res = lval == rval
    case token.BIN_NOT_EQUAL:
		res = lval != rval
	default:
		panic("Unhandled operator: " + operator.Type.String())
	}

	return BoolObj{
		Value: res,
	}
}

func check_type_compliant(val EvalObj, dt token.TokenType) {
	// for now our types are string, and int
	if val.Type() == EVAL_INT && dt == token.TYPE_INT {
		return
	} else if val.Type() == EVAL_STR && dt == token.TYPE_STR {
		return
	}

	panic("Unmatched data type and value: want " + dt.String() + " got: " + string(val.Type()))
}
