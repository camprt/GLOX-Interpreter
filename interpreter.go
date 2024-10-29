/*
* Evaluates expressions generated by the parser, with Interpreter Object acting as a visitor
* as established in the file expression.go.
* Where textbook relies on Java's generic Object classification, I use Go's interface{}
* Created: 9/23
* Modified: 10/7
 */

package main

import (
	"fmt"
)

type RuntimeError struct {
	token Token
	msg string
}

// assuming runtime error
func (rt RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Runtime Error: %v\n", rt.token.line, rt.msg)
}

type Interpreter struct {
	globals *Environment
	environment *Environment
	locals map[Expr]int
	hadRuntimeError bool
}

func newInterpreter() *Interpreter { //creates a nil enclosing env because this should be the global
	g := newEnvironment(nil)
	//I don't think go can do nested functions???? so that's gonna go in its own file
	g.define("clock", clock{})

	return &Interpreter{globals: g, environment: g, locals: make(map[Expr]int), hadRuntimeError: false}
}

func (itpr *Interpreter) interpret(statments []Stmt) {
	defer func() {
		//read as "err from recovered after failure"
		//			"if error occured, print the stuff"
		if err := recover(); err != nil {
			//check if its a runtime error
			fmt.Println("Had an error: ")
			itpr.hadRuntimeError = true
		}
	}()

	//if no errors,
	for _, statement := range statments {
		itpr.execute(statement)
		//fmt.Printf("Executed line %d\n", i)
	}
}

/**STATEMENT VISITORS**/
//Block Stmt
func (itpr *Interpreter) visitBlockStmt(stmt BlockStmt) interface{} {
	itpr.executeBlock(stmt.statements, newEnvironment(itpr.environment))
	return nil
}

//Class Stmt
func (itpr *Interpreter) visitClassStmt(stmt ClassStmt) interface{} {
	//extract superclass
	var super *LoxClass
	if stmt.superclass != nil {
		//object := itpr.evaluate(stmt.superclass)
		
		//check that result is a class
		object, ok := itpr.evaluate(stmt.superclass).(LoxClass)
		if !ok {
			itpr.error(&RuntimeError{token: stmt.superclass.name, msg: "Superclass must be a class"})
		}
		super = &object
	}
	
	//extract methods
	itpr.environment.define(stmt.name.lexeme, nil)

	//super methods
	if stmt.superclass != nil {
		itpr.environment = newEnvironment(itpr.environment)
		itpr.environment.define("super", super)
	}

	methods := make(map[string]LoxFunction)
	for _, m := range stmt.methods {
		function := LoxFunction{declaration: m, closure: itpr.environment, isInitializer: (m.name.lexeme == "init")}
		methods[m.name.lexeme] = function
	}

	class := LoxClass{name: stmt.name.lexeme, superclass: super, methods: methods}
	
	if stmt.superclass != nil {
		itpr.environment = itpr.environment.enclosing //enclsoing is nil?
	}
	
	itpr.environment.assign(stmt.name, class)
	return nil
}

//Expression Stmt
func (itpr *Interpreter) visitExpressionStmt(stmt ExpressionStmt) interface{} {
	itpr.evaluate(stmt.expression)
	return nil
}

//For Stmt AAAAAAAAAAAAhahhhhhhhhHHHHHHHHH
func (itpr *Interpreter) visitForStmt(stmt ForStmt) interface{} {
	if stmt.initializer != nil {
		itpr.execute(stmt.initializer)
	}

	//the loop!!!!!!!!!!!!
	for {
		if stmt.condition != nil {
			//eval the condition
			if !itpr.isTruthy(itpr.evaluate(stmt.condition)) {break}
		}

		itpr.execute(stmt.body)

		if stmt.increment != nil {
			itpr.evaluate(stmt.increment)
		}
	}

	// for itpr.isTruthy(itpr.evaluate(stmt.condition)) {
	// 	itpr.execute(stmt.body)
	// 	if stmt.increment != nil {
	// 		itpr.evaluate(stmt.increment)
	// 	}
	// }

	return nil
}

//Function Stmt
func (itpr *Interpreter) visitFunctionStmt(stmt FunctionStmt) interface{} {
	function := LoxFunction{declaration: stmt, closure: itpr.environment, isInitializer: false}
	itpr.environment.define(stmt.name.lexeme, function)
	return nil
}

//If Stmt
func (itpr *Interpreter) visitIfStmt(stmt IfStmt) interface{} {
	if itpr.isTruthy(itpr.evaluate(stmt.condition)) {
		itpr.execute(stmt.thenBranch)
	} else if (stmt.elseBranch != nil) {
		itpr.execute(stmt.elseBranch)
	}
	return nil
}

//Print Stmt
func (itpr *Interpreter) visitPrintStmt(stmt PrintStmt) interface{} {
	value := itpr.evaluate(stmt.expression)
	fmt.Println(itpr.stringify(value))
	return nil
}

type Return struct {
	value interface{}
}

//Return Stmt
func (itpr *Interpreter) visitReturnStmt(stmt ReturnStmt) interface{} {
	var value interface{} = nil
	if (stmt.value != nil) {value = itpr.evaluate(stmt.value)}
	panic(Return{value: value})
}

//Var Stmt
func (itpr *Interpreter) visitVarStmt(stmt VarStmt) interface{} {
	var value interface{} //default sets to nil
	if stmt.initializer != nil {
		value = itpr.evaluate(stmt.initializer)
	}

	//match the pair
	itpr.environment.define(stmt.name.lexeme, value)
	return nil
}

// While Stmt
func (itpr *Interpreter) visitWhileStmt(stmt WhileStmt) interface{} {
	for itpr.isTruthy(itpr.evaluate(stmt.condition)) {
		itpr.execute(stmt.body)
	}
	return nil
}

/**EXPRESSION VISITORS**/
// Assign
func (itpr *Interpreter) visitAssignExpr(expr AssignExpr) interface{} {
	value := itpr.evaluate(expr.value)

	//check thing exists first
	distance, ok := itpr.locals[expr]
	if ok {
		itpr.environment.assignAt(distance, expr.name, value)
	} else {
		itpr.globals.assign(expr.name, value)
	}
	return value
}

// Binary
func (itpr *Interpreter) visitBinaryExpr(expr BinaryExpr) interface{} {
	left := itpr.evaluate(expr.left)
	right := itpr.evaluate(expr.right)

	//perform operations
	switch expr.operator.kind {
	case MINUS:
		//makes sure that left and right are both numbers first
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) - right.(float64))
		} else {
			itpr.error(err)
		} //throws the error (kind of?)

	case SLASH:
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) / right.(float64))
		} else {
			itpr.error(err)
		}

	case STAR:
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) * right.(float64))
		} else {
			itpr.error(err)
		}

	case PLUS: //need to determine if adding nums or strings
		//check if floats?
		_, rType := right.(float64)
		_, lType := left.(float64)
		if rType && lType {
			return (left.(float64) + right.(float64))
		}

		//or strings?
		_, rType = right.(string)
		_, lType = left.(string)
		if rType && lType {
			return (left.(string) + right.(string))
		}

		//else "throw" (?) an error
		err := &RuntimeError{token: expr.operator, msg: "Operands must be 2 numbers or strings"}
		itpr.error(err)

	case GREATER:
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) > right.(float64))
		} else {
			itpr.error(err)
		}

	case GREATER_EQUAL:
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) >= right.(float64))
		} else {
			itpr.error(err)
		}

	case LESS:
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) < right.(float64))
		} else {
			itpr.error(err)
		}

	case LESS_EQUAL:
		err := itpr.checkNumberOperands(expr.operator, left, right)
		if err == nil {
			return (left.(float64) <= right.(float64))
		} else {
			itpr.error(err)
		}

	case BANG_EQUAL:
		return !itpr.isEqual(left, right)

	case EQUAL_EQUAL:
		return itpr.isEqual(left, right)
	} //end of switch case

	return nil
}

// Call
func (itpr *Interpreter) visitCallExpr(expr CallExpr) interface{} {
	callee := itpr.evaluate(expr.callee)

	var arguments []interface{}

	for _, arg := range expr.arguments {
		arguments = append(arguments, itpr.evaluate(arg))
	}

	//this is how we label the name to "callable" level priority - i think???
	function, ok := (callee).(LoxCallable)
	if !ok { //throw runtime error if not callable
		itpr.error(&RuntimeError{token: expr.paren, msg: "Can only call functions and classes."})
	}

	//check arity
	if (len(arguments) != function.arity()) {
		itpr.error(&RuntimeError{token: expr.paren, msg: fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments))})
	}

	return function.call(itpr, arguments)
}

//Get
func (itpr *Interpreter) visitGetExpr(expr GetExpr) interface{} {
	object := itpr.evaluate(expr.object)
	//check if its an instance
	inst, ok := object.(*LoxInstance) 
	if ok {
		val, err := inst.get(expr.name)
		//"Throw" the error
		if (err != nil) {
			panic(err)
		}
		return val
	}
	itpr.error(&RuntimeError{token: expr.name, msg: "Only instances have properties"})
	return nil
}

//Grouping
func (itpr *Interpreter) visitGroupingExpr(expr GroupingExpr) interface{} {
	return itpr.evaluate(expr.expression)
}

//Literal
func (itpr *Interpreter) visitLiteralExpr(expr LiteralExpr) interface{} {
	return expr.value
}

//Logical
func (itpr *Interpreter) visitLogicalExpr(expr LogicalExpr) interface{} {
	left := itpr.evaluate(expr.left)

	if (expr.operator.kind == OR) {
		if itpr.isTruthy(left) {return left}
	} else {
		if !itpr.isTruthy(left) {return left}
	}

	return itpr.evaluate(expr.right)
}

//Set
func (itpr *Interpreter) visitSetExpr(expr SetExpr) interface{} {
	object := itpr.evaluate(expr.object)

	objectInstance, ok := object.(*LoxInstance)
	if !ok {
		itpr.error(&RuntimeError{token: expr.name, msg: "Only instances have fields."})
	}

	value := itpr.evaluate(expr.value)
	objectInstance.set(expr.name, value)
	return value
}

//Super
func (itpr *Interpreter) visitSuperExpr(expr SuperExpr) interface{} {
	distance := itpr.locals[expr]
	superclass := itpr.environment.getAt(distance, "super").(*LoxClass)

	object := itpr.environment.getAt(distance - 1, "this").(*LoxInstance)

	method := superclass.findMethod(expr.method.lexeme)

	if method == nil {
		itpr.error(&RuntimeError{token: expr.method, msg: "Undefied property '"+expr.method.lexeme+"'."})
	}
	return method.bind(object)
}

//This
func (itpr *Interpreter) visitThisExpr(expr ThisExpr) interface{} {
	return itpr.lookUpVariable(expr.keyword, expr)
}

//Unary
func (itpr *Interpreter) visitUnaryExpr(expr UnaryExpr) interface{} {
	right := itpr.evaluate(expr.right)

	//applies minus or negation
	switch expr.operator.kind {
	case MINUS:
		//need to check that right is a num first, else throw an error
		err := itpr.checkNumberOperand(expr.operator, right)
		if err != nil { //error found
			itpr.error(err)
		} else {
			//else all good
			return -right.(float64) //convert to double
		}

	case BANG:
		return !itpr.isTruthy(right)
	}

	//should be unreachable
	return nil
}

// Variable
func (itpr *Interpreter) visitVariableExpr(expr VariableExpr) interface{} {
	return itpr.lookUpVariable(expr.name, expr)
}


/**HELPERS**/
// Passes itpr object to the expr's accept function, utilizing visitor functionality
func (itpr *Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(itpr)
}

// Tests if an object is "truthy" (only nil and false are falsey)
func (itpr *Interpreter) isTruthy(object interface{}) bool {
	//check if nil
	if object == nil {
		return false
	}

	//just return if a bool
	_, ok := object.(bool) //do another typecheck
	if ok {
		return object.(bool)
	}

	//all else is truthy/truthful/wtv
	return true
}

// Checks if 2 values are equal (nulls are equal)
// IDK if Go evaluates this way already but I'll do it by hand anyway
func (itpr *Interpreter) isEqual(l interface{}, r interface{}) bool {
	if l == nil && r == nil {
		return true
	}
	//check if l is still null
	if l == nil {
		return false
	}

	return l == r
}

// Checks that the given operand is a number type
// Used for error checking in unary evaluation
func (itpr *Interpreter) checkNumberOperand(operator Token, operand interface{}) *RuntimeError {
	_, isDouble := operand.(float64)

	//will return an error if not a double
	if !isDouble {
		return &RuntimeError{token: operator, msg: "Operand must be a number in a Unary expression"}
	}

	//else all good
	return nil
}

// Checks that the 2 given variables are numbers
func (itpr *Interpreter) checkNumberOperands(operator Token, left interface{}, right interface{}) *RuntimeError {
	_, lDouble := left.(float64)
	_, rDouble := right.(float64)

	//if both are doubles, proceed
	if lDouble && rDouble {
		return nil
	}

	//else sths not right
	return &RuntimeError{token: operator, msg: "Operands must be numbers in Binary expressions"}
}

// Added myself, just calls error's report and flags in Interpreter
func (itpr *Interpreter) error(err *RuntimeError) {
	fmt.Println(err.Error())
	itpr.hadRuntimeError = true
	panic(err)
}

// Displays the results of an interpreted expression
func (itpr *Interpreter) stringify(object interface{}) string {
	//null?
	if object == nil {
		return "nil"
	}

	//else just sprint
	return fmt.Sprint(object)
}

//Navigates to statement visitor to "execut"
func (itpr *Interpreter) execute(stmt Stmt) {
	stmt.accept(itpr)
}

//Add variables found in the resolver to the locals map
func (itpr *Interpreter) resolve(expr Expr, depth int) {
	itpr.locals[expr] = depth
}

//Executes full block of statements in sub-environment
func (itpr *Interpreter) executeBlock(statements []Stmt, env *Environment) {
	previous := itpr.environment

	defer func() { //idk what a try-finally is but this is it
		itpr.environment = previous
	}()

	//"Try" block
	itpr.environment = env
	for _, stmt := range statements {
		itpr.execute(stmt)
	}
}

//looks up the variable either in the locals or globals
func (itpr *Interpreter) lookUpVariable(name Token, expr Expr) interface{} {
	distance, ok := itpr.locals[expr]
	if ok {
		return itpr.environment.getAt(distance, name.lexeme)
	} else {
		return itpr.globals.get(name)
	}
}
