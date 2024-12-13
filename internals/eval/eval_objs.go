package eval

import (
	"fmt"
	"nox/internals/parser"
)

type ObjType string

const (
	EVAL_FUNC_DEF     = "func_def"
	EVAL_BUILINT_FUNC = "builtin_func"
	EVAL_INT          = "int"
	EVAL_STR          = "string"
	EVAL_NULL         = "null"
	EVAL_ERROR        = "error"
)

type EvalObj interface {
	Type() ObjType
}

type NullObj struct{}

func (n NullObj) Type() ObjType {
	return EVAL_NULL
}

type ErrorObj struct{}

func (n ErrorObj) Type() ObjType {
	return EVAL_ERROR
}

type FuncDefObj struct {
	Ident parser.Identifier
	Body  parser.BodyStatement
	// params
}

func (f FuncDefObj) Type() ObjType {
	return EVAL_FUNC_DEF
}

type BuiltinFuncObj struct {
	fn func(args ...EvalObj) EvalObj
}

func (b BuiltinFuncObj) Type() ObjType {
	return EVAL_BUILINT_FUNC
}

var builtins = map[string]BuiltinFuncObj{
	"print": {
		fn: func(args ...EvalObj) EvalObj {
			if len(args) < 1 {
				return EVAL_ERROR_OBJ
			}

			for _, arg := range args {
				fmt.Println(arg)
			}

			return EVAL_NULL_OBJ
		},
	},
}

type IntObj struct {
	Value int
}

func (i IntObj) Type() ObjType {
	return EVAL_INT
}

type StrObj struct {
	Value string
}

func (s StrObj) Type() ObjType {
	return EVAL_STR
}

type EvalContext struct {
	objs map[string]EvalObj

	outer *EvalContext
}

func (e *EvalContext) Get(key string) EvalObj {
	val, ok := e.objs[key]

	if !ok {
		if e.outer == nil {
			return EVAL_ERROR_OBJ
		} else {
			return e.outer.Get(key)
		}
	}

	return val
}

func (e *EvalContext) CreateNewEnclosedCtx() *EvalContext {
	return &EvalContext{
		objs:  map[string]EvalObj{},
		outer: e,
	}
}
