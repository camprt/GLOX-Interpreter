/*
* File to organize where variables are stored in a Hash Map structure
* Created: 10/2
* Modified: 10/7
 */

package main

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func newEnvironment(enc *Environment) *Environment {
	return &Environment{enclosing: enc, values: make(map[string]interface{})}
}

//Retrieves the value of a variable
func (env *Environment) get(name Token) interface{} {
	//check the variable exists first
	//"contains key" method sub
	_, ok := env.values[name.lexeme]
	if ok {
		return env.values[name.lexeme]
	}

	//check for var in eclosing env
	if env.enclosing != nil {
		return env.enclosing.get(name)
	}

	//if doesnt exist anywhere, "throws" the error
	err := RuntimeError{token: name, msg: fmt.Sprintf("Undefined variable '%s'", name.lexeme)}
	fmt.Println(err.Error())
	panic(err)
}

func (env *Environment) assign(name Token, value interface{}) {
	_, ok := env.values[name.lexeme]
	if ok {
		env.define(name.lexeme, value)
		return
	}

	//check in enclosing env
	if env.enclosing != nil {
		env.enclosing.assign(name, value)
		return
	}

	//if nowhere, "throw" and error
	err := RuntimeError{token: name, msg: fmt.Sprintf("Undefined variable '%s'.", name.lexeme)}
	fmt.Println(err.Error())
	panic(err)
}

//Binds a new name-value pair
func (env *Environment) define(name string, value interface{}) {
	env.values[name] = value
}

//Walks up the chain of parents to retrieve the value of a variable
func (env *Environment) ancestor(distance int) *Environment {
	curEnv := env
	for i := 0; i < distance; i++ {
		curEnv = curEnv.enclosing
	}
	return curEnv
}

//Retrieve value of a variable at a distance from local environment
func (env *Environment) getAt(distance int, name string) interface{} {
	return env.ancestor(distance).values[name]
}

//Assigns a new value to a variable at a distance from local environment
func (env *Environment) assignAt(distance int, name Token, value interface{}) {
	env.ancestor(distance).values[name.lexeme] = value
}


