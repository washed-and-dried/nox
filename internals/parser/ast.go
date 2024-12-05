package parser

import (
	"fmt"
	"nox/internals/token"
)

// this file will contain the definitions of different nodes of AST.
// [ASTNode] Statement -> Type, StatementObject
// StatementObject -> AssigmentStatement n = 1;, ExpressionStatement 1 + 2;, ...
// Program will contain the list of such statements

type AST struct {
	Stmt Statement
	Type AstType
}

type AstType string

const (
	STMTS = "stmts"

	EXPR           = "expr"
	EXPR_FUNC_CALL = "expr_func_call"
	EXPR_BIN       = "expr_bin"

	ASSIGN   = "assign"
	FUNC_DEF = "func_def"
)

type Statement interface {
	String() string
}

type ExpressionStmt struct {
	Type  ExprType
	Value ExprValue
}

func (e ExpressionStmt) String() string {
	return fmt.Sprintf("ExpressionStmt of type [%s]", e.Type)
}

type FuncDefStmt struct {
	Body  *BodyStatement
	Ident string
	// Params []*ExpressionStmt
}

func (f FuncDefStmt) String() string {
	return fmt.Sprintf("FuncDefExpr [%s]", f.Ident)
}

type BodyStatement struct {
	Stmts []Statement
}

func (b BodyStatement) String() string {
	return fmt.Sprintf("BodyStatement with [%d] statements", len(b.Stmts))
}

type FuncCallExpr struct {
	Ident string
	Args  []ExpressionStmt
}

func (f FuncCallExpr) String() string {
	return fmt.Sprintf("FuncCallExpr [%s]", f.Ident)
}

type BinaryExpr struct {
	Left     *ExpressionStmt
	Right    *ExpressionStmt
	Operator token.Token
}

func (b BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr: [%s]", b.Operator.Type.String())
}

type IntExpr struct {
	Value int64 // change this later
}

func (i IntExpr) String() string {
	return fmt.Sprintf("IntAST with value [%d]", i.Value)
}
