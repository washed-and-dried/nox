package eval

import "nox/internals/parser"

type ObjType string

const (
	EVAL_FUNC_DEF = "func_def"
	EVAL_INT      = "int"
	EVAL_STR      = "string"
	EVAL_NULL     = "null"
	EVAL_ERROR    = "error"
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

type FunDefObj struct {
	Ident string
	Body  parser.BodyStatement
	// params
}

func (f FunDefObj) Type() ObjType {
	return EVAL_FUNC_DEF
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
