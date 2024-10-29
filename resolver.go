// /*
// * Looks ahead at variable usage to help with scoping.
// * As of 10/7: issue seems to be that maps aren't initialized properly before adding things
// * Created: 10/7
//  */

package main

import (
	"fmt"
)

//Go does not have a prebuilt stack structure, so doing it ourselves
type Stack struct {
	items []interface{}
	size int
}

func (s *Stack) push(data interface{}) {
	s.items = append(s.items, data)
	s.size += 1
}

func (s *Stack) pop() interface{} {
	if s.isEmpty() {
		return nil
	}

	data := s.items[s.size - 1]
	s.size -= 1
	s.items = s.items[:s.size]
	return data
}

func (s *Stack) isEmpty() bool {
	return (s.size == 0)
}

func (s *Stack) peek() interface{} {
	if s.isEmpty() {
		return nil
	}
	return s.items[s.size - 1]
}

func (s *Stack) getAt(i int) interface{} {
	if (i >= s.size) {
		return nil
	}
	return s.items[i]
}

/**Function Type**/
type FunctionType int
const (
	NOFUNC = iota
	FUNCTION
	INITIALIZER
	METHOD
)

/**Actual Resolver stuff now**/
type Resolver struct {
	interpreter *Interpreter
	scopes Stack
	curFunction FunctionType
	curClass ClassType
	hadError bool
}

func newResolver(itpr *Interpreter) *Resolver {
	return &Resolver{interpreter: itpr, curFunction: NOFUNC, curClass: NOCLASS, hadError: false}
}

/**VISITORS**/
func (r *Resolver) visitBlockStmt(stmt BlockStmt) interface{} {
	r.beginScope()
	r.resolveStmts(stmt.statements)
	r.endScope()
	return nil
}

func (r *Resolver) visitClassStmt(stmt ClassStmt) interface{} {
	enclosingClass := r.curClass
	r.curClass = REGCLASS

	r.declare(stmt.name)
	r.define(stmt.name)

	if (stmt.superclass != nil && stmt.name.lexeme == stmt.superclass.name.lexeme) {
		r.error(stmt.superclass.name, "A class can't inherit from itself.")
	}

	if stmt.superclass != nil {
		r.curClass = SUBCLASS
		r.resolveExpr(stmt.superclass)
		r.beginScope()
		r.peekScopes()["super"] = true
	}

	r.beginScope()
	r.peekScopes()["this"] = true

	for _, method := range stmt.methods {
		var declaration FunctionType = METHOD
		if method.name.lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()
	if stmt.superclass != nil {
		r.endScope()
	}

	r.curClass = enclosingClass
	return nil
}

func (r *Resolver) visitExpressionStmt(stmt ExpressionStmt) interface{} {
	r.resolveExpr(stmt.expression)
	return nil
}

func (r *Resolver) visitForStmt(stmt ForStmt) interface{} {
	if stmt.initializer != nil {
		r.resolveStmt(stmt.initializer)
	}
	if stmt.condition != nil {
		r.resolveExpr(stmt.condition)
	}
	if stmt.increment != nil {
		r.resolveExpr(stmt.increment)
	}
	r.resolveStmt(stmt.body)
	
	return nil
}

func (r *Resolver) visitFunctionStmt(stmt FunctionStmt) interface{} {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) visitIfStmt(stmt IfStmt) interface{} {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if (stmt.elseBranch != nil) {
		r.resolveStmt(stmt.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(stmt PrintStmt) interface{} {
	r.resolveExpr(stmt.expression)
	return nil
}

func (r *Resolver) visitReturnStmt(stmt ReturnStmt) interface{} {
	if (r.curFunction == NOFUNC) {
		r.error(stmt.keyword, "Can't return from top-level code.")
	}

	if (stmt.value != nil) {
		if r.curFunction == INITIALIZER {
			r.error(stmt.keyword, "Can't return a value from an initializer.")
		}
		r.resolveExpr(stmt.value)
	}
	return nil
}

func (r *Resolver) visitVarStmt(stmt VarStmt) interface{} {
	r.declare(stmt.name)
	if (stmt.initializer != nil) {
		r.resolveExpr(stmt.initializer)
	}
	r.define(stmt.name)
	return nil
}

func (r *Resolver) visitWhileStmt(stmt WhileStmt) interface{} {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return nil
}

func (r *Resolver) visitAssignExpr(expr AssignExpr) interface{} {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) visitBinaryExpr(expr BinaryExpr) interface{} {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r *Resolver) visitCallExpr(expr CallExpr) interface{} {
	r.resolveExpr(expr.callee)

	for _, arg := range expr.arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) visitGetExpr(expr GetExpr) interface{} {
	r.resolveExpr(expr.object)
	return nil
}

func (r *Resolver) visitGroupingExpr(expr GroupingExpr) interface{} {
	r.resolveExpr(expr.expression)
	return nil
}

func (r *Resolver) visitLiteralExpr(expr LiteralExpr) interface{} {
	return nil //no vars to resolve
}

func (r *Resolver) visitLogicalExpr(expr LogicalExpr) interface{} {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
	return nil
}

func (r *Resolver) visitSetExpr(expr SetExpr) interface{} {
	r.resolveExpr(expr.value)
	r.resolveExpr(expr.object)
	return nil
}

func (r *Resolver) visitSuperExpr(expr SuperExpr) interface{} {
	if r.curClass == NOCLASS {
		r.error(expr.keyword, "Can't use 'super' outside of a class.")
	} else if r.curClass != SUBCLASS {
		r.error(expr.keyword, "Can't use 'super' in a class with no superclass")
	}

	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r *Resolver) visitThisExpr(expr ThisExpr) interface{} {
	if r.curClass == NOCLASS {
		r.error(expr.keyword, "Can't use 'this' outside of a class.")
		return nil
	}

	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r *Resolver) visitUnaryExpr(expr UnaryExpr) interface{} {
	r.resolveExpr(expr.right)
	return nil
}

func (r *Resolver) visitVariableExpr(expr VariableExpr) interface{} {
	if !r.scopes.isEmpty() {
		defined, declared := r.peekScopes()[expr.name.lexeme]
		//if it exists but isn't yet defined, throw an error
		if (declared && !defined) {
			r.error(expr.name, "Can't read local variable in its own intiializer")
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil
}

/**RESOLVE FUNCTIONS**/
//Visits list of statements
func (r *Resolver) resolveStmts(statements []Stmt) {
	for _, stmt := range statements {
		r.resolveStmt(stmt)
	}
}

//Turn resolver into stmt visitor
func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}

//Turn resolver into expr visitor
func (r *Resolver) resolveExpr(expr Expr) {
	expr.accept(r)
}

//Resolves function stuff in its own scope
func (r *Resolver) resolveFunction(fun FunctionStmt, ftype FunctionType) {
	enclosingFunction := r.curFunction
	r.curFunction = ftype

	r.beginScope()
	for _, p := range fun.params {
		r.declare(p)
		r.define(p)
	}
	r.resolveStmts(fun.body)
	r.endScope()
	r.curFunction = enclosingFunction
}

//Resolves the given variable
func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.size - 1; i >= 0; i-- {
		//"contains key"
		// _, exists := r.scopes.getAt(i).(map[string]bool)[name.lexeme]
		// if exists {
		// 	r.interpreter.resolve(expr, r.scopes.size - 1 - i)
		// 	return
		// }
		s := r.scopes.getAt(i).(map[string]bool)
		_, exists := s[name.lexeme]
		if exists {
			varDepth := r.scopes.size - 1 - i
			r.interpreter.resolve(expr, varDepth)
			return
		}
	}
}

/**HELPERS**/
//Written myself to avoid keep doing type assertions
func (r *Resolver) peekScopes() map[string]bool{
	return r.scopes.peek().(map[string]bool)
}

//prints the error messages
func (r *Resolver) error(token Token, msg string) {
	var loc string
	if (token.kind == EOF) {loc = "end"
	} else {loc = token.lexeme}

	fmt.Printf("[line %d] Resolution Error at \"%s\": %s\n", token.line, loc, msg)
	r.hadError = true
}

//Creates new block scope
func (r *Resolver) beginScope() {
	r.scopes.push(make(map[string]bool))
}

//Exits the current scope
func (r *Resolver) endScope() {
	r.scopes.pop()
}

//Adds new var to innermost scope so it takes precendence over those in outer
func (r *Resolver) declare(name Token) {
	if r.scopes.isEmpty() {return}

	//add to a map
	s := r.peekScopes()
	if (s != nil) {
		//don't add if already exists
		_, exists := s[name.lexeme]
		if (exists) {
			r.error(name, "Already a variable with this name in this scope.")
		}
		//else
		s[name.lexeme] = false
	}
}

//sets variable to true after initializer has been resolved
func (r *Resolver) define(name Token) {
	if r.scopes.isEmpty() {return}
	s := r.peekScopes()
	s[name.lexeme] = true
}
