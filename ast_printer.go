// /*
// * Class to help print a syntax tree into Lisp formater
// * to make debugging the parser easier, but not
// * technically needed
// * Created: 9/10
// * Modified: 9/13
//  */

package main

// import "fmt"

// type AstPrinter struct {
// }

// func newAstPrinter() *AstPrinter {
// 	return &AstPrinter{}
// }

// // basic print function for provided expression
// func (a AstPrinter) print(e Expr) string {
// 	return e.accept(a).(string)
// }

// /**VISITORS**/
// func (a AstPrinter) visitLiteralExpr(expr LiteralExpr) interface{} {
// 	if expr.value == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprint(expr.value)
// }

// func (a AstPrinter) visitGroupingExpr(expr GroupingExpr) interface{} {
// 	return a.parenthesize("group", expr.expression)
// }

// func (a AstPrinter) visitUnaryExpr(expr UnaryExpr) interface{} {
// 	return a.parenthesize(expr.operator.lexeme, expr.right)
// }

// func (a AstPrinter) visitBinaryExpr(expr BinaryExpr) interface{} {
// 	return a.parenthesize(expr.operator.lexeme, expr.left, expr.right)
// }

// func (a AstPrinter) visitVariableExpr(expr VariableExpr) interface{} {
// 	return expr.name
// }

// // Builds a parenthesized expression in a readable form, similar to a Lisp expression
// func (a AstPrinter) parenthesize(name string, exprList ...Expr) string {
// 	var str string
// 	str = "(" + name
// 	for _, expr := range exprList {
// 		str += " " + a.print(expr)
// 	}
// 	str += ")"
// 	return str
// }

// // func main() {
// // 	var l UnaryExpr;
// // 		l.operator = Token{kind: MINUS, lexeme: "-", literal:nil, line:1}
// // 		l.right = LiteralExpr{value:123}
// // 	var tok Token = Token{kind: STAR, lexeme:"*", literal:nil, line:1}
// // 	var r GroupingExpr;
// // 		r.expression = LiteralExpr{value:45.67}
// // 	var expr BinaryExpr;
// // 		expr.left = l
// // 		expr.right = r
// // 		expr.operator = tok

// // 	a := new(AstPrinter)
// // 	fmt.Println(a.print(expr))
// // }
