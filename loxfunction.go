/*
* File to specify which expressions are "callable" and which aren't
* Created: 10/4
* Modified: 10/8
 */

package main

import (
	"fmt"
	"time"
)

//"Interface"
type LoxCallable interface {
	arity() int
	call(itpr *Interpreter, arguments []interface{}) interface{}
}

/**List of Native Callables**/

//Clock
type clock struct {}

func (c clock) arity() int {return 0}

func (c clock) call(itpr *Interpreter, args []interface{}) interface{} {
	return float64(time.Now().UnixMilli())
}

func (c clock) String() string {
	return "<native fn>"
}

//User defined functions
type LoxFunction struct {
	declaration FunctionStmt
	closure *Environment
	isInitializer bool
}

func (f LoxFunction) bind(instance *LoxInstance) LoxFunction {
	env := newEnvironment(f.closure)
	env.define("this", instance)
	return LoxFunction{declaration: f.declaration, closure: env, isInitializer: f.isInitializer}
}

func (f LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f LoxFunction) call(itpr *Interpreter, arguments []interface{}) (returnValue interface{}) {
	//catch
	defer func() {
		if err := recover(); err != nil {
			if rv, ok := err.(Return); ok {
				if (f.isInitializer) {returnValue = f.closure.getAt(0, "this")
				//have the value set in the func declaration so return works properly ig
				} else {returnValue = rv.value}
				return
			}
			panic(err)
		}
	}()
	
	env := newEnvironment(f.closure)
	for i := 0; i < len(f.declaration.params); i++ {
		env.define(f.declaration.params[i].lexeme, arguments[i])
	}

	//try
	itpr.executeBlock(f.declaration.body, env)

	//Force any initializer to return "this"
	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}

	return nil
}

func (f LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}

