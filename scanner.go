/*
* Scanner file to scan through a given input file and break
* down its lexical grammar. Not using Lex!
* Created: 9/6
* Modified: 10/2
 */

package main

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source            string
	tokens            []Token
	start, curr, line int
	hadError          bool
}

// Constructer
func newScanner(src string) *Scanner {
	return &Scanner{source: src, start: 0, curr: 0, line: 1, hadError: false}
}

// Hash Map for reserved words
var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// Scans the tokens in the raw source code string
func (s *Scanner) scanTokens() []Token {
	//loop to scan all chars in source
	for !s.isAtEnd() {
		//determine lexeme
		s.start = s.curr
		s.scanToken()
	}

	//add null token to list
	s.tokens = append(s.tokens, Token{kind: EOF, lexeme: "", literal: nil, line: s.line})
	return s.tokens
}

// Determines the token kind (w/ receiver s)
func (s *Scanner) scanToken() {
	c := s.advance()

	//all the kinds you can imagine
	switch c {
	//single chars
	case '(':
		s.addBasicToken(LEFT_PAREN)
	case ')':
		s.addBasicToken(RIGHT_PAREN)
	case '{':
		s.addBasicToken(LEFT_BRACE)
	case '}':
		s.addBasicToken(RIGHT_BRACE)
	case ',':
		s.addBasicToken(COMMA)
	case '.':
		s.addBasicToken(DOT)
	case '-':
		s.addBasicToken(MINUS)
	case '+':
		s.addBasicToken(PLUS)
	case ';':
		s.addBasicToken(SEMICOLON)
	case '*':
		s.addBasicToken(STAR)

	//check for 2 char lexemes
	case '!':
		//check for !=
		if s.match('=') {
			s.addBasicToken(BANG_EQUAL)
			//just !
		} else {
			s.addBasicToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addBasicToken(EQUAL_EQUAL) //==
		} else {
			s.addBasicToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addBasicToken(LESS_EQUAL)
		} else {
			s.addBasicToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addBasicToken(GREATER_EQUAL)
		} else {
			s.addBasicToken(GREATER)
		}

	//special case for /
	case '/':
		if s.match('/') {
			//consume comment until new line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}

			//challenge: implement /**/ comments
		} else if s.match('*') {
			for s.peek() != '*' && s.peekNext() != '/' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addBasicToken(SLASH)
		}

	//whitespace chars (ignore)
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++

	//strings
	case '"':
		s.addString()

	//report unknown chars
	default:
		//numbers
		if s.isDigit(c) {
			s.addNumber()

			//reserved words
		} else if s.isAlpha(c) {
			s.addIdentifier()

		} else {
			s.error(s.curr, "Unexpected char")
		}

	}
}

// Return next character in source, advance curr iterater
func (s *Scanner) advance() rune {
	next := rune(s.source[s.curr])
	s.curr++
	return next
}

/**ADDERS**/

// Create new token for the current lexeme with a "nil" ltieral
func (s *Scanner) addBasicToken(kind TokenType) {
	s.addToken(kind, nil)
}

// Create new token for current lexeme w/ specified literal
func (s *Scanner) addToken(kind TokenType, literal interface{}) {
	//extract lexeme
	text := s.source[s.start:s.curr]
	s.tokens = append(s.tokens, Token{kind: kind, lexeme: text, literal: literal, line: s.line})
}

// Adds a string token
func (s *Scanner) addString() {
	//while still in string & not at end, keep consuming
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	//if didn't close "" before end of line, throw error
	if s.isAtEnd() {
		s.error(s.curr, "Unterminated string")
		return
	}

	s.advance() //consumes closing "

	//add whole string to token list
	value := s.source[s.start+1 : s.curr-1]
	s.addToken(STRING, value)

}

// Adds a number token
func (s *Scanner) addNumber() {
	//keep advancing until no more numbers
	for s.isDigit(s.peek()) {
		s.advance()
	}

	//check for decimal
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		//consume .
		s.advance()
		//get fractional parts
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	//convert number string to float
	value, _ := strconv.ParseFloat(s.source[s.start:s.curr], 64)
	s.addToken(NUMBER, value)
}

// Adds a longer lexeme token
func (s *Scanner) addIdentifier() {
	//extract lexeme
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	//determine tokenType
	text := s.source[s.start:s.curr]
	kind, exists := keywords[text]
	if !exists {
		kind = IDENTIFIER
	}

	//add non literal token
	s.addBasicToken(kind)
}

/**CHECKERS & LOOKAHEAD**/

// Check if scanTokens() func has reached end of source string
func (s *Scanner) isAtEnd() bool {
	return s.curr >= len(s.source)
}

// Checks next char to determine multi-char lexemes
func (s *Scanner) match(expected rune) bool {
	//short circuit if last char
	if s.isAtEnd() {
		return false
	}

	//or if chars dont match
	if rune(s.source[s.curr]) != expected {
		return false
	}

	//must be true otherwise, consume next char
	s.curr++
	return true
}

// Lookahead function to help ignore comments
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.curr])
}

// Lookahead for next next character
func (s *Scanner) peekNext() rune {
	//return nil if reach out of bounds
	if (s.curr + 1) >= len(s.source) {
		return '\000'
	}

	//return 2 ahead
	return rune(s.source[s.curr+1])
}

// Checks if input is a number literal
func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// Checks if character is a letter
func (s *Scanner) isAlpha(c rune) bool {
	//only allows a-z, A-Z, or underscores
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

// Checks if is a number or a digit
func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

/**Errors**/
func (s *Scanner) error(line int, msg string) {
	fmt.Println("[line ", line, "] Error: ", msg)
	s.hadError = true
}
