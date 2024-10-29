// /*
//   - Runs unit tests for GLOX interpreter
//     */
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

//Struct to store all the file tests
func TestRunner(t *testing.T) {
	tests := []struct{
		testName string
		srcCode string
		expectedOutput string
 	}{
		//name of test		command					output			error
		{"string test", "print \"hello world\";", "hello world\n"},
		{"multi-line string", "print \"hello\nworld\";", "hello\nworld\n"},
		{"number test", "print 678.098;", "678.098\n"},
		{"nil test", "print nil and 34;", "nil\n"},
		{"no semicolon error", "print 89", "[line 1] Error at end: Expect ';' after value.\n"},
		{"comment after code", "print 1 + 1; //comment!", "2\n"},
		{"full line comment", `//hello
			print "hello";`, "hello\n"},
		{"arithmetic test", "print 3.8 - 9 * 2.3;", "-16.9\n"},
		{"concatenation test", "print \"hello\" + \" \" + \"world\";", "hello world\n"},
		{"logical operators", "print true and (true and !false);", "true\n"},
		{"gte test", "print 9 >= 3;", "true\n"},
		{"gte failure", "print 2 >= 3;", "false\n"},
		{"lte test", "print 3 <= 7;", "true\n"},
		{"lte failure", "print 9 <= 7;", "false\n"},
		{"equality test", "print 9.2 == 9.2;", "true\n"},
		{"not equal test", "print 9 != 0;", "true\n"},
		{"var declaration", "var a = 3; print a;", "3\n"},
		{"var assignment", "var b = 3; b = 5; print b;", "5\n"},
		{"var reassignment", "var c = 2; print c; c = 5; print 5;", "2\n5\n"},
		{"block scoping", `var a = "global a";
			var b = "global b";
			var c = "global c";
			{
			var a = "outer a";
			var b = "outer b";
			{
				var a = "inner a";
				print a;
				print b;
				print c;
			}
			print a;
			print b;
			print c;
			}
			print a;
			print b;
			print c;`, "inner a\nouter b\nglobal c\nouter a\nouter b\nglobal c\nglobal a\nglobal b\nglobal c\n"},
		{"if test", "if (true) {if (false) {print \"don't print\";} else {print \"print this!\";}}", "print this!\n"},
		{"for loop test", `var a = 0;
			var temp;

			for (var b = 1; a < 10; b = temp + b) {
			print a;
			temp = a;
			a = b;
			}`, "0\n1\n1\n2\n3\n5\n8\n"},
		{"while loop test", `var a = 0;
			var temp;
			var b = 1;

			while (a < 10) {
				print a;
				temp = a;
				a = b;
				b = temp + b;
			}`, "0\n1\n1\n2\n3\n5\n8\n"},
		{"function declaration", `fun sayHi(first, last) {
			print "Hi, " + first + " " + last + "!";
			}

			sayHi("There", "Man");`, "Hi, There Man!\n"},

		{"return test", `fun sayHi(first, last) {
			return "Hello, " + first + " " + last;
		  	}

		  	print sayHi("Dear", "Reader");`, "Hello, Dear Reader\n",},

		{"nested function test", `fun makeCounter() {
				var i = 0;
				fun count() {
				i = i + 1;
				print i;
				}

				return count;
			}

			var counter = makeCounter();
			counter(); // "1".
			counter(); // "2".`, "1\n2\n"},
		{"scope test", `//Passed!
			var a = "global";
			{
			fun showA() {
				print a;
			}

			showA();
			var a = "block";
			showA();
			}`, "global\nglobal\n"},

		{"class methods", `class Bacon {
			eat() {
			print "Crunch crunch crunch!";
			}
		}

		Bacon().eat(); // Prints "Crunch crunch crunch!".
		`, "Crunch crunch crunch!\n"},

		{"this test", `class Cake {
			taste() {
				var adjective = "delicious";
				print "The " + this.flavor + " cake is " + adjective + "!";
				}
			}

			var cake = Cake();
			cake.flavor = "German chocolate";
			cake.taste();`, "The German chocolate cake is delicious!\n"},

		{"inheritance test", `class Doughnut {
				cook() {
				print "Fry until golden brown.";
				}
			}

			class BostonCream < Doughnut {}

			BostonCream().cook();`, "Fry until golden brown.\n"},
		{"superclass test", `class Doughnut {
				cook() {
				print "Fry until golden brown.";
				}
			}

			class BostonCream < Doughnut {
				cook() {
				super.cook();
				print "Pipe full of custard and coat with chocolate.";
				}
			}

			BostonCream().cook();`, "Fry until golden brown.\nPipe full of custard and coat with chocolate.\n"},
	}

	for _, testCase := range tests {
		//need to redirect the standard output
		ogOs := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		runner := newRunner()
		runner.run(testCase.srcCode)

		w.Close()
		output, _ := ioutil.ReadAll(r)
		os.Stdout = ogOs

		if string(output) != testCase.expectedOutput {
			t.Errorf("Output error at Test %s: got %s, expected %s", testCase.testName, strconv.Quote(string(output)), strconv.Quote(testCase.expectedOutput))
		} else {
			fmt.Printf("%s: \tPASSED\n", testCase.testName)
		}
	}
	fmt.Print("\n")
}

func TestFileRunner(t *testing.T) {
	//get the files stored in the test folder
	filepaths, err := filepath.Glob(filepath.Join("tests", "*.lox"))
	if err != nil {
		t.Fatal(err)
	}

	//go thru each file
	for _, f := range filepaths {
		//extract the file
		_, file := filepath.Split(f)
		testName := file[:len(file)-len(filepath.Ext(f))]

		//Create a new function to test whole contents of each file
		t.Run(testName, func(t *testing.T) {
			test, err := os.ReadFile(f)
			if err != nil {
				t.Fatal("Error reading test file:", err)
			}

			//get the correct output from text file
			correctFile := filepath.Join("test_results", testName+"_results.txt")
			expected, err := os.ReadFile(correctFile)
			if err != nil {
				t.Fatal("Error reading expected output:", err)
			}

			//now run the real code, same as in testRun
			ogOs := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			runner := newRunner()
			runner.run(string(test))

			w.Close()
			output, _ := ioutil.ReadAll(r)
			os.Stdout = ogOs

			//compare
			if string(output) != string(expected) {
				t.Errorf("Error at test %s: got %s expected %s", testName, strconv.Quote(string(output)), strconv.Quote(string(expected)))
			}
		})
	}
}
