/*
* Parses together given strings into the correct grammar with correct precedence
* Created: 9/16
* Modified: 10/7
 */

package main

import (
	"fmt"
)

type ParseError struct {
	token Token
	msg   string
}

// Displays full Parse Error
func (pe *ParseError) error() string {
	var errLoc string
	if pe.token.kind == EOF {
		errLoc = "at end"
	} else {
		errLoc = "at '" + pe.token.lexeme + "'"
	}
	return fmt.Sprintf("[line %d] Error %v: %v", pe.token.line, errLoc, pe.msg)
}

type Parser struct {
	tokens   []Token
	cur      int
	hadError bool
}

// Constructor
func newParser(t []Token) *Parser {
	return &Parser{tokens: t, cur: 0, hadError: false}
}

// starts the parser
func (p *Parser) parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

/**GRAMMAR DEFINITIONS**/
// expression -> assignment
func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment → ( call "." )? IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		_, ok := expr.(VariableExpr)
		if ok {
			name := expr.(VariableExpr).name
			return AssignExpr{name: name, value: value}
		} else if _, ok := expr.(GetExpr); ok {
			get := expr.(GetExpr)
			return SetExpr{object: get.object, name: get.name, value: value}
		}

		p.error(&ParseError{token: equals, msg: "Invalid assignment target."})
	}

	return expr
}

//logic_or → logic_and ( "or" logic_and )* ;
func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = LogicalExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

//logic_and → equality ( "and" equality )* ;
func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = LogicalExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// equality -> comparison ( ("!-=" | "==") comparison)*
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}

	return expr
}

// unary → ( "!" | "-" ) unary | call ;
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{operator: operator, right: right}
	}

	return p.call()
}

//call → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = GetExpr{object: expr, name: name}
		} else {
			break
		}
	}

	return expr
}

//Basically the argument part of the grammar
//argument → expression ( "," expression )* ;
func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		for {
			//do
			if (len(arguments) >= 255) {
				p.error(&ParseError{token: p.peek(), msg: "Can't have more than 255 arguments."})
			}
			arguments = append(arguments, p.expression())
			//while
			if !p.match(COMMA) {break}
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return CallExpr{callee: callee, paren: paren, arguments: arguments}
}

// primary → "true" | "false" | "nil" | "this"
//			| NUMBER | STRING | IDENTIFIER | "(" expression ")"
//			| "super" "." IDENTIFIER ;
func (p *Parser) primary() Expr {
	if p.match(FALSE) {return LiteralExpr{value: false}}
	if p.match(TRUE) {return LiteralExpr{value: true}}
	if p.match(NIL) {return LiteralExpr{value: nil}}

	if p.match(NUMBER, STRING) {
		return LiteralExpr{value: p.previous().literal}
	}
	if p.match(SUPER) {
		keyword := p.previous()
		p.consume(DOT, "Expect '.' after 'super'")
		method := p.consume(IDENTIFIER, "Expect superclass method name.")
		return SuperExpr{keyword: keyword, method: method}
	}
	if p.match(THIS) {
		return ThisExpr{keyword: p.previous()}
	}
	if p.match(IDENTIFIER) {
		return VariableExpr{name: p.previous()}
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return GroupingExpr{expression: expr}
	}

	//throws a parse error
	p.error(&ParseError{token: p.peek(), msg: "Error: expected an expression"})
	return nil
}

// declaration → classDecl | funDecl | varDecl | statement ;
func (p *Parser) declaration() Stmt {
	defer func() {
		//only goes here if theres a parse error
		if err := recover(); err != nil {
			p.hadError = true
			p.synchronize()
		}
	}()

	//match the stmt type
	if p.match(CLASS) {return p.classDeclaration()}
	if p.match(FUN) {return p.function("function")}
	if p.match(VAR) {return p.varDeclaration()}
	return p.statement()
}

//classDecl → "class" IDENTIFIER ("<" IDENTIFIER)? "{" function* "}" ;
func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect class name.")

	var super *VariableExpr
	if p.match(LESS) {
		p.consume(IDENTIFIER, "Expect superclass name.")
		super = &VariableExpr{name: p.previous()}
	}

	p.consume(LEFT_BRACE, "Expect '{' before class body.")

	var methods []FunctionStmt
	for (!p.check(RIGHT_BRACE) && !p.isAtEnd()) {
		methods = append(methods, p.function("method"))
	}
	p.consume(RIGHT_BRACE, "Expect '}' after class body.")

	return ClassStmt{name: name, superclass: super, methods: methods}
}

//function → IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(kind string) FunctionStmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind)) 

	//parse parameters
	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			//do
			if (len(parameters) >= 255) {
				p.error(&ParseError{token: p.peek(), msg: "Can't have more than 255 parameters"})
			}

			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			//while
			if !p.match(COMMA) {break}
		}
	} 
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	//parse body
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return FunctionStmt{name: name, params: parameters, body: body}
}

//varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect a variable name.")

	var initializer Expr

	//get intial value if it exists, otherwise leave nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return VarStmt{name: name, initializer: initializer}
}


/**STATEMENTS**/
//statement → exprStmt | forStmt | ifStmt | printStmt | returnStmt | whileStmt | block;
func (p *Parser) statement() Stmt {
	//check statement type & call correct method
	if p.match(FOR) {return p.forStatement()}
	if p.match(IF) {return p.ifStatement()}
	if p.match(PRINT) {return p.printStatement()}
	if p.match(RETURN) {return p.returnStatement()}
	if p.match(WHILE) {return p.whileStatement()}
	if p.match(LEFT_BRACE) {return BlockStmt{statements: p.block()}}

	return p.expressionStatement()
}

//forStmt → "for" "(" ( varDecl | exprStmt | ";" )
//			expression? ";"
//	 		expression? ")" statement ;
func (p *Parser) forStatement() Stmt {
	p.consume (LEFT_PAREN, "Expect '(' after 'for'.")

	//initializer clause
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	//condition clause
	var condition Expr = nil
	if !p.check(SEMICOLON) {condition = p.expression()}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	//increment clause
	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {increment = p.expression()}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	//body
	body := p.statement()

	return ForStmt{initializer: initializer, condition: condition, increment: increment, body: body}
}

//ifStmt → "if" "(" expression ")" statement 
//		( "else" statement )? ;
func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

//exprStmt → expression ";" ;
func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return ExpressionStmt{expression: expr}
}

//block → "{" declaration* "}"
func (p *Parser) block() []Stmt {
	var statements []Stmt
	//get all the stuff inside the block
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

//printStmt → "print" expression ";" ;
func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return PrintStmt{expression: value}
}

//returnStmt → "return" expression? ";" ;
func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr = nil
	if !p.check(SEMICOLON) {
		value = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' aftern return value.")
	return ReturnStmt{keyword: keyword, value: value}
}

//while → "while" "(" expression ")" statement ;
func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return WhileStmt{condition: condition, body: body}
}


/**HELPER FUNCTIONS**/
// Checks if the current token matches any of the given types, then consumes if true
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

// Returns true if cur token is of type t
func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().kind == t
}

// Consumes the current token & returns it
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.cur++
	}
	return p.previous()
}

// Checks if reached EOF
func (p *Parser) isAtEnd() bool {
	return p.peek().kind == EOF
}

// Returns current token
func (p *Parser) peek() Token {
	return p.tokens[p.cur]
}

// Returns previous token
func (p *Parser) previous() Token {
	return p.tokens[p.cur-1]
}

/**ERROR HANDLING**/

// Works same as match, but throws an error if the checked token doesn't match expected
func (p *Parser) consume(tokType TokenType, message string) Token {
	if p.check(tokType) {
		return p.advance()
	}

	//creates a new ParseError//
	p.error(&ParseError{token: p.peek(), msg: message})
	return Token{}
}

// passes the error to the main class
func (p *Parser) error(err *ParseError) {
	p.hadError = true
	fmt.Println(err.error())
	panic(err)
}

// Helps reset the parser's state
func (p *Parser) synchronize() {
	p.advance()

	//Discards tokens until reaches beginning of next full statement
	for !p.isAtEnd() {
		if p.previous().kind == SEMICOLON {
			return
		}

		switch p.peek().kind {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
