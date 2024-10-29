
Welcome to my implementation of the Crafting Interpreters Tree-Walk Interpreter in Go, aka
"GLOX." In this file you should find all of the information to run this program and test
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
testing installed in order to run these tests. Unfortunately, there is no way I could find
to format my own Pass/Fail formatting, so the extension's own output formatting is used.
However, the repl tests do follow my own format: "[test name]: PASSED/FAILED". The output
of my own unit tests can be found in the file titled "Unit_Test_Output.txt".

Each file in this project has a comment block at the top of the file containing a small
description of what it contributes to the rest of the project, as well as a creation and
last modified date.

Don't forget to use semicolons! -- Go does not require semicolons, and I forgot mine several 
times during testing!

An sample of a Lox program and its output can be found in the file titled "sample_run.txt".

Snapshots of the code after completing the chapters "Representing Code", "Statements and
State", and "Resolving and Binding" can be found in the subfolders titled "Snapshot" 1, 2,
and 3 respectively.


Created by Campbell Thompson for CS 403-001
