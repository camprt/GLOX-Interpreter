/*
* Main file for Crafting Interpreters glox
* Tree-Walk Version
* Created 9/5
* Modified: 10/7
 */

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//"bufio"
)

func main() {
	if len(os.Args) > 2 {
		log.Fatal("Usage: glox [script]")
	} else if (len(os.Args) == 2) {
		runner := newRunner()
		runner.runFile(os.Args[1])
	} else {
		runner := newRunner()
		runner.runPrompt()
	}
}

func (r *Runner) runFile(path string) {
	//get the file from the path
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error: could not read file path")
		panic(err)
	}

	//send to runner
	r.run(string(file))

	if r.hadError {
		os.Exit(65)
	}
	if r.hadRuntimeError {
		os.Exit(70)
	}
}

func (r *Runner) runPrompt() {
	for {
		fmt.Print("> ") //delim? i think that's the word

		//read in user input
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error: Could not read user input")
			log.Fatal(err)
		}

		r.run(input)

		//reset for next round
		r.hadError = false
		r.hadRuntimeError = false
	}
}

//doesn't really do much
type Runner struct {
	hadError bool
	hadRuntimeError bool
	interpreter *Interpreter
}

//"Constructor"
func newRunner() *Runner {
	return &Runner{hadError: false, hadRuntimeError: false, interpreter: newInterpreter()}
}

/**Runs inputted Lox statement from given stream "source"*/
func (r *Runner) run(source string) {

	//scan source stream into tokens
	scanner := newScanner(source)
	tokens := scanner.scanTokens()

	parser := newParser(tokens)
	statements := parser.parse()

	//Stop if a syntax error
	r.hadError = parser.hadError
	if r.hadError {return}

	resolver := newResolver(r.interpreter)
	resolver.resolveStmts(statements)

	//stop if resolution error
	r.hadError = resolver.hadError
	if r.hadError {return}

	r.interpreter.interpret(statements)
	r.hadRuntimeError = r.interpreter.hadRuntimeError
	if (r.hadRuntimeError) {return}

}
