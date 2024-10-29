
Welcome to my implementation of the Crafting Interpreters Tree-Walk Interpreter in Go, aka
"GLOX". In this file you should find all of the information to run this program and test
your own output.

The main file is titled "lox.go." If using VSCode, use "go run ." to run the repl,
or use "go run . [path to file]" if you want to run a Lox file. To exit the repl, 
press Control-D or Control-C.

My pre-built tests are in the subfolder titled "tests", with all the files following
the format "[filename].lox". The expected results of these files are in the subfolder 
"test_results", with all the corresponding files titled "[original filname]_results.txt".
These are the folders that my unit tests will look for test cases and expected results in,
so if you would like to add any test cases of your own, please add them in this format.

My unit tests use the "testing" tool in VSCode, so you must have the Go extension for
testing installed in order to run these tests. The repl tests follow the own format: 
"[test name]: PASSED/FAILED", and the run file tests follow Go's own unit testing ouput.

Each file in this project has a comment block at the top of the file containing a small
description of what it contributes to the rest of the project, as well as a creation and
last modified date.

Don't forget to use semicolons! -- Go does not require semicolons, and I forgot mine several 
times during testing!


Created by Campbell Thompson for CS 403-001
