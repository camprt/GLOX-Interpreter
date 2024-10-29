/*
* Class to format expression grammar
* Automated in text, done manually here because generate made no sense to me
* Created 9/10
* Modified: 10/8
 */

package main

//import "fmt"

// The "super" class
type Expr interface {
	accept(v Visitor) interface{}
}

/**Basic Expression Types**/
type AssignExpr struct {
	name  Token
	value Expr
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

type CallExpr struct {
	callee Expr
	paren Token
	arguments []Expr
}

type GetExpr struct {
	object Expr
	name Token
}

type GroupingExpr struct {
	expression Expr
}

type LiteralExpr struct {
	value interface{}
}

type LogicalExpr struct {
	left Expr
	operator Token
	right Expr
}

type SetExpr struct {
	object Expr
	name Token
	value Expr
}

type SuperExpr struct {
	keyword Token
	method Token
}

type ThisExpr struct {
	keyword Token
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

type VariableExpr struct {
	name Token
}

/**Visitor struct/class**/
type Visitor interface {
	visitAssignExpr(expr AssignExpr) interface{}
	visitBinaryExpr(expr BinaryExpr) interface{}
	visitCallExpr(expr CallExpr) interface{}
	visitGetExpr(expr GetExpr) interface{}
	visitGroupingExpr(expr GroupingExpr) interface{}
	visitLiteralExpr(expr LiteralExpr) interface{}
	visitLogicalExpr(expr LogicalExpr) interface{}
	visitSetExpr(expr SetExpr) interface{}
	visitSuperExpr(expr SuperExpr) interface{}
	visitThisExpr(expr ThisExpr) interface{}
	visitUnaryExpr(expr UnaryExpr) interface{}
	visitVariableExpr(expr VariableExpr) interface{}
}

/**Accept funcs**/
func (expr AssignExpr) accept(v Visitor) interface{} {
	return v.visitAssignExpr(expr)
}

func (expr BinaryExpr) accept(v Visitor) interface{} {
	return v.visitBinaryExpr(expr)
}

func (expr CallExpr) accept(v Visitor) interface{} {
	return v.visitCallExpr(expr)
}

func (expr GetExpr) accept(v Visitor) interface{} {
	return v.visitGetExpr(expr)
}

func (expr GroupingExpr) accept(v Visitor) interface{} {
	return v.visitGroupingExpr(expr)
}

func (expr LiteralExpr) accept(v Visitor) interface{} {
	return v.visitLiteralExpr(expr)
}

func (expr LogicalExpr) accept(v Visitor) interface{} {
	return v.visitLogicalExpr(expr)
}

func (expr SetExpr) accept(v Visitor) interface{} {
	return v.visitSetExpr(expr)
}

func (expr SuperExpr) accept(v Visitor) interface{} {
	return v.visitSuperExpr(expr)
}

func (expr ThisExpr) accept(v Visitor) interface{} {
	return v.visitThisExpr(expr)
}

func (expr UnaryExpr) accept(v Visitor) interface{} {
	return v.visitUnaryExpr(expr)
}

func (expr VariableExpr) accept(v Visitor) interface{} {
	return v.visitVariableExpr(expr)
}
