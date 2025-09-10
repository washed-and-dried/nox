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

type Program struct {
	Stmts []Statement
}

const (
	STMTS = "stmts"

	EXPR     = "expr"
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

type Identifier struct {
	Name string
}

func (i Identifier) String() string {
	return fmt.Sprintf("Identifier [%s]", i.Name)
}

type FuncDefStmt struct {
	Body  *BodyStatement
	Ident Identifier
	// Params []*ExpressionStmt
}

func (f FuncDefStmt) String() string {
	return fmt.Sprintf("FuncDefExpr [%s]", f.Ident)
}

type ReturnStmt struct {
	ExprStmt ExpressionStmt
	Void     bool
}

func (r ReturnStmt) String() string {
	return fmt.Sprintf("ReturnStmtExpr [%s]", r.ExprStmt.String())
}

type AssignStmt struct {
	Type  token.Token
	Ident string
	Value ExpressionStmt
}

func (a AssignStmt) String() string {
	return fmt.Sprintf("AssignStmt type [%s] with value [%s]", a.Type.Type.String(), a.Value.String())
}

type BodyStatement struct {
	Stmts []Statement
}

func (b BodyStatement) String() string {
	return fmt.Sprintf("BodyStatement with [%d] statements", len(b.Stmts))
}

type FuncCallExpr struct {
	Ident Identifier
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

type StrExpr struct {
	Value string
}

func (s StrExpr) String() string {
	return fmt.Sprintf("StrAST with value [%s]", s.Value)
}

type BoolExpr struct {
	Value bool
}

func (b BoolExpr) String() string {
	return fmt.Sprintf("BoolExpr with value as %t", b.Value)
}

type SubscriptExpr struct {
	Index *ExpressionStmt
    Ident Identifier
}

func (s SubscriptExpr) String() string {
    return fmt.Sprintf("SubScriptExpr: %s[%d]", s.Ident.Name, s.Index)
}

type VarUpdation struct {
	Var   Identifier
	Value ExpressionStmt
}

func (v VarUpdation) String() string {
	return fmt.Sprintf("VarUpdation: [%s] = [%s]", v.Var, v.Value)
}

type LoopStmt struct {
	Init     Statement // we are putting Statement since it  could be an empty statement like for(;;){...}
	Cond     Statement
	Updation Statement
	Body     BodyStatement
}

func (f LoopStmt) String() string {
	return fmt.Sprintf("LoopStmt")
}

type NullStmt struct {}

func (f NullStmt) String() string {
	return fmt.Sprintf("NullStmt")
}

type IfStmt struct {
	Cond Statement
	Init Statement
	Else Statement
	Body BodyStatement
}

func (f IfStmt) String() string {
	return fmt.Sprintf("IfStmt")
}
