// /*
// * Attempt 2 for the resolver
// * Created: 10/22 (yikes)
//  */

package main

// import (
// 	"fmt"
// )

// //Data structure for the scopes
// type scope map[string]bool

// func newScope() scope {
// 	return make(map[string]bool)
// }

// type ScopeStack struct {
// 	scopes []scope
// 	size int
// }

// func newScopeStack() *ScopeStack {
// 	return &ScopeStack{scopes: make([]scope, 0), size: 0}
// }

// func (s *ScopeStack) top() scope {
// 	if s.size == 0 {
// 		return nil
// 	}

// 	return s.scopes[s.size - 1]
// }

// func (s *ScopeStack) push(sc scope) {
// 	s.scopes = append(s.scopes, sc)
// 	s.size = s.size + 1
// }

// func (s *ScopeStack) pop() scope {
// 	if s.size == 0 {return nil}

// 	top := s.top()
// 	s.size = s.size - 1
// 	s.scopes = s.scopes[:s.size]
// 	return top;
// }

// func (s *ScopeStack) get(i int) scope {
// 	return s.scopes[i]
// }

// func (s *ScopeStack) isEmpty() bool {
// 	return s.size == 0
// }

// //Function Type
// type FunctionType int
// const (
// 	NOFUNCTION = iota
// 	FUNCTION
// 	INITIALIZER
// 	METHOD
// )

// /*RESOLVER STRUCTURE*/
// type Resolver struct {
// 	interpreter *Interpreter
// 	scopes *ScopeStack
// 	currentFunction FunctionType
// 	currentClass ClassType
// 	hadError bool
// }

// func newResolver(itpr *Interpreter) *Resolver {
// 	return &Resolver{interpreter: itpr, scopes: newScopeStack(), currentFunction: NOFUNCTION, currentClass: NOCLASS, hadError: false}
// }

// func (r *Resolver) Error(token Token, msg string) {
// 	var loc string
// 	if (token.kind == EOF) {loc = "end"
// 	} else {loc = token.lexeme}

// 	fmt.Printf("[line %d] Resolution Error at \"%s\": %s\n", token.line, loc, msg)
// 	r.hadError = true
// }

// //Resolve a list of statements
// func (r *Resolver) resolveStmts(stmts []Stmt) {
// 	for _, s := range stmts {
// 		r.resolveStmt(s)
// 	}
// }

// /**VISITORS**/
// func (r *Resolver) visitBlockStmt(stmt BlockStmt) interface{} {
// 	r.beginScope()
// 	r.resolveStmts(stmt.statements)
// 	r.endScope()
// 	return nil
// }

// func (r *Resolver) visitClassStmt(stmt ClassStmt) interface{} {
// 	enclosingClass := r.currentClass
// 	r.currentClass = CLASS

// 	r.declare(stmt.name)
// 	r.define(stmt.name)

// 	if stmt.superclass != nil {
// 		if stmt.name.lexeme == stmt.superclass.name.lexeme {
// 			r.Error(stmt.superclass.name, "A class cannot inherit from itself.")
// 		}
// 	}

// 	if stmt.superclass != nil {
// 		r.currentClass = SUBCLASS
// 		r.resolveExpr(stmt.superclass)
// 		r.beginScope()
// 		r.scopes.top()["super"] = true
// 	}

// 	r.beginScope()
// 	r.scopes.top()["this"] = true

// 	for _, method := range stmt.methods {
// 		var declaration FunctionType = METHOD
// 		if method.name.lexeme == "init" {
// 			declaration = INITIALIZER
// 		}

// 		r.resolveFunction(method, declaration)
// 	}

// 	r.endScope()
// 	if stmt.superclass != nil {r.endScope()}

// 	r.currentClass = enclosingClass
// 	return nil
// }

// func (r *Resolver) visitExpressionStmt(stmt ExpressionStmt) interface{} {
// 	r.resolveExpr(stmt.expression)
// 	return nil
// }

// func (r *Resolver) visitFunctionStmt(stmt FunctionStmt) interface{} {
// 	r.declare(stmt.name)
// 	r.define(stmt.name)

// 	r.resolveFunction(stmt, FUNCTION)
// 	return nil
// }

// func (r *Resolver) visitIfStmt(stmt IfStmt) interface{} {
// 	r.resolveExpr(stmt.condition)
// 	r.resolveStmt(stmt.thenBranch)
// 	if stmt.elseBranch != nil {
// 		r.resolveStmt(stmt.elseBranch)
// 	}
// 	return nil
// }

// func (r *Resolver) visitPrintStmt(stmt PrintStmt) interface{} {
// 	r.resolveExpr(stmt.expression)
// 	return nil
// }

// func (r *Resolver) visitReturnStmt(stmt ReturnStmt) interface{} {
// 	if r.currentFunction == NOFUNCTION {
// 		r.Error(stmt.keyword, "Can't return from top-level code.")
// 	}

// 	if stmt.value != nil {
// 		if r.currentFunction == INITIALIZER {
// 			r.Error(stmt.keyword, "Can't return a value from an initializer.")
// 		}
// 		r.resolveExpr(stmt.value)
// 	}
// 	return nil
// }

// func (r *Resolver) visitVarStmt(stmt VarStmt) interface{} {
// 	r.declare(stmt.name)
// 	if stmt.initializer != nil {
// 		r.resolveExpr(stmt.initializer)
// 	}
// 	r.define(stmt.name)
// 	return nil;
// }

// func (r *Resolver) visitWhileStmt(stmt WhileStmt) interface{} {
// 	r.resolveExpr(stmt.condition)
// 	r.resolveStmt(stmt.body)
// 	return nil
// }

// //Expression visitors
// func  (r *Resolver) visitAssignExpr(expr AssignExpr) interface{} {
// 	r.resolveExpr(expr.value)
// 	r.resolveLocal(expr, expr.name)
// 	return nil
// }

// func (r *Resolver) visitBinaryExpr(expr BinaryExpr) interface{} {
// 	r.resolveExpr(expr.left)
// 	r.resolveExpr(expr.right)
// 	return nil
// }

// func (r *Resolver) visitCallExpr(expr CallExpr) interface{} {
// 	r.resolveExpr(expr.callee)
// 	for _, arg := range expr.arguments {
// 		r.resolveExpr(arg)
// 	}
// 	return nil
// }

// func (r *Resolver) visitGetExpr(expr GetExpr) interface{} {
// 	r.resolveExpr(expr.object)
// 	return nil
// }

// func (r *Resolver) visitGroupingExpr(expr GroupingExpr) interface{} {
// 	r.resolveExpr(expr.expression)
// 	return nil
// }

// func (r *Resolver) visitLiteralExpr(expr LiteralExpr) interface{} {
// 	return nil
// }

// func (r *Resolver) visitLogicalExpr(expr LogicalExpr) interface{} {
// 	r.resolveExpr(expr.left)
// 	r.resolveExpr(expr.right)
// 	return nil
// }

// func (r *Resolver) visitSetExpr(expr SetExpr) interface{} {
// 	r.resolveExpr(expr.value)
// 	r.resolveExpr(expr.object)
// 	return nil
// }

// func (r *Resolver) visitSuperExpr(expr SuperExpr) interface{} {
// 	if r.currentClass == NOCLASS {
// 		r.Error(expr.keyword, "Can't use 'super' outside of a class.")
// 	} else if r.currentClass == SUBCLASS {
// 		r.Error(expr.keyword, "Can't use 'super' in a class with no superclass.")
// 	}

// 	r.resolveLocal(expr, expr.keyword)
// 	return nil
// }

// func (r *Resolver) visitThisExpr(expr ThisExpr) interface{} {
// 	if r.currentClass == NOCLASS {
// 		r.Error(expr.keyword, "Can't use 'this' outside of a class.")
// 	}

// 	r.resolveLocal(expr, expr.keyword)
// 	return nil
// }

// func (r *Resolver) visitUnaryExpr(expr UnaryExpr) interface{} {
// 	r.resolveExpr(expr.right)
// 	return nil
// }

// func (r *Resolver) visitVariableExpr(expr VariableExpr) interface{} {
// 	//check the scope exists
// 	s := r.scopes.top()
// 	if s != nil {
// 		//check the var exists & is defined
// 		defined, exists := s[expr.name.lexeme]
// 		if exists && !defined {
// 			r.Error(expr.name, "can't read local variable in its own initializer.")
// 		}
// 	}

// 	r.resolveLocal(expr, expr.name)
// 	return nil
// }

// /**HELPERS**/
// func (r *Resolver) resolveStmt(stmt Stmt) {
// 	stmt.accept(r)
// }

// func (r *Resolver) resolveExpr(expr Expr) {
// 	expr.accept(r)
// }

// func (r *Resolver) resolveFunction(function FunctionStmt, ftype FunctionType) {
// 	enclosingFunction := r.currentFunction
// 	r.currentFunction = ftype
// 	r.beginScope()
// 	for _, param := range function.params {
// 		r.declare(param)
// 		r.define(param)
// 	}
// 	r.resolveStmts(function.body)
// 	r.endScope()
// 	r.currentFunction = enclosingFunction
// }

// //Start a new scope
// func (r *Resolver) beginScope() {
// 	r.scopes.push(newScope())
// }

// func (r *Resolver) endScope() {
// 	r.scopes.pop()
// }

// //Declare a new var/func
// func (r *Resolver) declare(name Token) {
// 	sc := r.scopes.top()
// 	if sc == nil {return}

// 	//if already exists
// 	_, exists := sc[name.lexeme]
// 	if exists {r.Error(name, "Already a value within this scope.")}

// 	//else
// 	sc[name.lexeme] = false
// }

// func (r *Resolver) define(name Token) {
// 	sc := r.scopes.top()
// 	if sc == nil {return}
// 	//else flag as defined
// 	sc[name.lexeme] = true
// }

// func (r *Resolver) resolveLocal(expr Expr, name Token) {
// 	for i := r.scopes.size - 1; i >= 0; i-- {
// 		s := r.scopes.get(i)
// 		if s != nil {
// 			_, exists := s[name.lexeme]
// 			if exists {
// 				r.interpreter.resolve(expr, r.scopes.size - 1 - i)
// 				return
// 			}
// 		}
// 	}
// }

