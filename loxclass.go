/*
* Stores all the class setup stuff for lox classes, including class instances
* Created: 10/8
 */

package main

import (
	"fmt"
)

type ClassType int
const (
	NOCLASS = iota
	REGCLASS
	SUBCLASS
)

/**CLASS OBJECT**/
type LoxClass struct {
	name string
	superclass *LoxClass
	methods map[string]LoxFunction
}

//Returns given method
func (c LoxClass) findMethod(name string) *LoxFunction {
	m, exists := c.methods[name]
	if exists {
		return &m
	}

	//or an inherited method
	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil
}

func (c LoxClass) String() string {
	return c.name
}

//"Implements loxcallable" stuff
func (c LoxClass) call(itpr *Interpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{class: c}
	intializer := c.findMethod("init")
	if intializer != nil {
		intializer.bind(instance).call(itpr, arguments)
	}

	return instance
}

func (c LoxClass) arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {return 0}
	return initializer.arity()
}

/**INSTANCE**/
type LoxInstance struct {
	class LoxClass
	fields map[string]interface{}
}

//Getter
func (inst *LoxInstance) get(name Token) (interface{}, error) {
	value, exists := inst.fields[name.lexeme]
	if exists {
		return value, nil
	}

	method := inst.class.findMethod(name.lexeme)
	if (method != nil) {return method.bind(inst), nil}

	//if doesn't exist, throw runtime error
	return nil, RuntimeError{token: name, msg: fmt.Sprintf("Undefined property '%s'.", name.lexeme)}
}

//Setter
func (inst *LoxInstance) set(name Token, value interface{}) {
	//initialize the map if nil
	if (inst.fields == nil) {
		inst.fields = make(map[string]interface{})
	}
	inst.fields[name.lexeme] = value
}

func (inst LoxInstance) String() string {
	return fmt.Sprintf("%s instance", inst.class.name)
}