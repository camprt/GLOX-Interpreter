/*
* File to create the syntax tree for statements
* Similar structure to the expression file
* Created: 9/25
* Modified: 10/7
 */

package main

// general type
type Stmt interface {
	accept(visitor StmtVisitor) interface{}
}

// specific types
type BlockStmt struct {
	statements []Stmt
}

type ClassStmt struct {
	name Token
	superclass *VariableExpr
	methods []FunctionStmt
}

type ExpressionStmt struct {
	expression Expr
}

//Im tired of dealing with that damn syntatic sugar!!!!!!!!!!!!!!
type ForStmt struct {
	initializer Stmt
	condition Expr
	increment Expr
	body Stmt
}

type FunctionStmt struct {
	name Token
	params []Token
	body []Stmt
}

type IfStmt struct {
	condition Expr
	thenBranch Stmt
	elseBranch Stmt
}

type PrintStmt struct {
	expression Expr
}

type ReturnStmt struct {
	keyword Token
	value Expr
}

type VarStmt struct {
	name        Token
	initializer Expr
}

type WhileStmt struct {
	condition Expr
	body Stmt
}

/**VISITOR**/
type StmtVisitor interface {
	visitBlockStmt(stmt BlockStmt) interface{}
	visitClassStmt(stmt ClassStmt) interface{}
	visitExpressionStmt(stmt ExpressionStmt) interface{}
	visitForStmt(stmt ForStmt) interface{}
	visitFunctionStmt(stmt FunctionStmt) interface{}
	visitIfStmt(stmt IfStmt) interface{}
	visitPrintStmt(stmt PrintStmt) interface{}
	visitReturnStmt(stmt ReturnStmt) interface{}
	visitVarStmt(stmt VarStmt) interface{}
	visitWhileStmt(stmt WhileStmt) interface{}
}

/**Accept funcs**/
func (s BlockStmt) accept(v StmtVisitor) interface{} {
	return v.visitBlockStmt(s)
}

func (s ClassStmt) accept(v StmtVisitor) interface{} {
	return v.visitClassStmt(s)
}

func (s ExpressionStmt) accept(v StmtVisitor) interface{} {
	return v.visitExpressionStmt(s)
}

func (s ForStmt) accept(v StmtVisitor) interface{} {
	return v.visitForStmt(s)
}

func (s FunctionStmt) accept(v StmtVisitor) interface{} {
	return v.visitFunctionStmt(s)
}

func (s IfStmt) accept(v StmtVisitor) interface{} {
	return v.visitIfStmt(s)
}

func (s PrintStmt) accept(v StmtVisitor) interface{} {
	return v.visitPrintStmt(s)
}

func (s ReturnStmt) accept(v StmtVisitor) interface{} {
	return v.visitReturnStmt(s)
}

func (s VarStmt) accept(v StmtVisitor) interface{} {
	return v.visitVarStmt(s)
}

func (s WhileStmt) accept(v StmtVisitor) interface{} {
	return v.visitWhileStmt(s)
}
