package parser

import (
	"errors"
	"fmt"
	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/util"
	"log"
)

type Parser struct {
	lexemes       []lexer.Lexeme
	currentStep   int
	currentSym    lexer.LexemeType
	currentLexeme lexer.Lexeme
}

func hasNext(parser *Parser) bool {
	return parser.lexemes[parser.currentStep].Type != lexer.LT_END
}

func next(parser *Parser) {
	if hasNext(parser) {
		parser.currentStep++
		parser.currentSym = parser.lexemes[parser.currentStep].Type
		parser.currentLexeme = parser.lexemes[parser.currentStep]
	}
}

func curr(parser *Parser) lexer.Lexeme {
	return parser.currentLexeme
}

func peek(parser *Parser) lexer.Lexeme {
	if hasNext(parser) {
		return parser.lexemes[parser.currentStep+1]
	}

	return lexer.Lexeme{Type: lexer.LT_NONE}
}

func accept(parser *Parser, symbol lexer.LexemeType) bool {

	if parser.currentSym == symbol {
		next(parser)
		fmt.Printf("Accepted symbol: %s\n", lexer.LexemeTypeLabels[parser.currentSym])
		return true
	}

	return false
}

func expect(parser *Parser, symbol lexer.LexemeType) bool {
	if accept(parser, symbol) {
		return true
	}

	errorMessage := fmt.Sprintf(
		"%sUnexpected symbol: %s, expected: %s%s\n",
		util.Red,
		lexer.LexemeTypeLabels[parser.currentSym],
		lexer.LexemeTypeLabels[symbol],
		util.Reset)

	err := errors.New(errorMessage)

	fmt.Println(err)

	return false
}

type ExpressionType int
type StatementType int

const (
	ET_BINARY ExpressionType = iota
	ET_UNARY
	ET_LITERAL
	ET_GROUP
)

const (
	ET_STATEMENT StatementType = iota
)

type AST_Expression struct {
	eType    ExpressionType
	lhs      *AST_Expression
	operator lexer.LexemeType
	rhs      *AST_Expression
	value    string
}

type AST_Statement struct {
}

type AST_Program struct {
	Statements []*AST_Statement
}

func expressionLiteral(value string) *AST_Expression {
	expr := &AST_Expression{
		eType: ET_LITERAL,
		value: value,
	}

	return expr
}

func expressionUnary(operator lexer.LexemeType, rhs *AST_Expression) *AST_Expression {
	_rhs := &AST_Expression{
		eType:    rhs.eType,
		lhs:      rhs.lhs,
		operator: rhs.operator,
		rhs:      rhs.rhs,
		value:    rhs.value,
	}

	expr := &AST_Expression{
		eType:    ET_UNARY,
		operator: operator,
		rhs:      _rhs,
	}

	return expr
}

func expressionBinary(lhs *AST_Expression, operator lexer.LexemeType, rhs *AST_Expression) *AST_Expression {
	_lhs := &AST_Expression{
		eType:    lhs.eType,
		lhs:      lhs.lhs,
		operator: lhs.operator,
		rhs:      lhs.rhs,
		value:    lhs.value,
	}

	_rhs := &AST_Expression{
		eType:    rhs.eType,
		lhs:      rhs.lhs,
		operator: rhs.operator,
		rhs:      rhs.rhs,
		value:    rhs.value,
	}

	expr := &AST_Expression{
		eType:    ET_BINARY,
		lhs:      _lhs,
		operator: operator,
		rhs:      _rhs,
		value:    "",
	}

	return expr
}

func expressionGroup(lhs *AST_Expression) *AST_Expression {

	expr := &AST_Expression{
		eType:    ET_GROUP,
		lhs:      lhs,
		operator: lexer.LT_NONE,
		rhs:      nil,
		value:    "",
	}

	return expr
}

func PrintTree(tree *AST_Expression, depth int) {

	prefixChar := "≫ "
	prefix := ""

	for i := 0; i < depth; i++ {
		prefix += prefixChar
	}

	if tree.eType == ET_GROUP {
		fmt.Println(prefix + "[ GROUP ]")

		//fmt.Println(prefix + prefixChar + "[ LHS ]")

		PrintTree(tree.lhs, depth+1)
	} else if tree.eType == ET_BINARY {

		switch tree.operator {
		case lexer.LT_PLUS:
			fmt.Println(prefix + "[ ADD ]")
		case lexer.LT_MINUS:
			fmt.Println(prefix + "[ SUBTRACT ]")
		case lexer.LT_MULTIPLY:
			fmt.Println(prefix + "[ MULTIPLY ]")
		case lexer.LT_DIVIDE:
			fmt.Println(prefix + "[ DIVIDE ]")
		}

		//fmt.Println(prefix + prefixChar + "[ LHS ]")

		PrintTree(tree.lhs, depth+1)

		//fmt.Println(prefix + prefixChar + "[ RHS ]")

		PrintTree(tree.rhs, depth+1)
	} else if tree.eType == ET_LITERAL {
		fmt.Println(prefix + "[ VALUE ]")
		fmt.Println(prefix + prefixChar + tree.value)
	}
}

func Create(lexer lexer.Lexer) Parser {
	return Parser{lexemes: lexer.Lexemes, currentLexeme: lexer.Lexemes[0], currentSym: lexer.Lexemes[0].Type}
}

func Start(parser *Parser) []*AST_Expression {

	expressions := make([]*AST_Expression, 0)

	for hasNext(parser) {

		if currentLexeme(parser).Type == lexer.LT_NUMBER {
			expr := expression(parser)
			expressions = append(expressions, expr)
			continue
		} else {
			step(parser)
		}

		expr := expression(parser)
		expressions = append(expressions, expr)

		continue
	}

	return expressions
}

func StartNew(parser *Parser) *AST_Program {

	program := program(parser)

	return program
}

/*

Language Gramamr

program -> (statement) END

statement 	-> IF "(" expression ")" "{" statement "}"
			-> WHILE "(" expression ")" "{" statement "}"
			-> IMPORT STRING
			-> LET IDENTIFIER ( ";" | "=" expression ";" )

*/

func program(parser *Parser) *AST_Program {

	program := &AST_Program{Statements: make([]*AST_Statement, 0)}

	for hasNext(parser) {

		//if accept(parser, lexer.LT_END) {
		//	fmt.Println("Finished")
		//}

		fmt.Println(lexer.LexemeTypeLabels[parser.currentLexeme.Type])

		sttmnt := statement(parser)

		program.Statements = append(program.Statements, sttmnt)

		continue
	}

	return program
}

func statement(parser *Parser) *AST_Statement {
	if accept(parser, lexer.LT_VAL) {
		if expect(parser, lexer.LT_IDENTIFIER) {
			if expect(parser, lexer.LT_EQUALS) {
				fmt.Println("Here")
			}
		}
	}
	return &AST_Statement{}
}

func condition(parser *Parser) {
}

/*

Expression Grammar

expression -> term

primary -> LT_NUMBER | LT_LPAREN expression LT_RPAREN

term -> factor (( LT_PLUS | LT_MINUS ) factor)*

factor -> unary (( LT_DIVIDE | LT_MULTIPLY ) unary)*

unary -> ( LT_BANG | LT_MINUS | LT_PLUS ) unary | primary
*/

// TODO: Finish to satisfy grammar
func primary(parser *Parser) *AST_Expression {
	//fmt.Println("Primary")

	c := currentLexeme(parser).Type

	if c == lexer.LT_NUMBER || c == lexer.LT_FLOAT {
		rhs := currentLexeme(parser).Label
		//fmt.Println(rhs)

		expr := expressionLiteral(rhs)
		return expr
	} else if c == lexer.LT_LPAREN {
		step(parser)
		expr := expression(parser)

		expr = expressionGroup(expr)

		//TODO: Check if current symbol is )
		//fmt.Printf("After group, current lexeme: %s\n", currentLexeme(parser).Label)
		return expr
	} else {
		log.Panic("Unexpected path")
		return nil
	}
}

func expression(parser *Parser) *AST_Expression {

	//fmt.Println("Expression")

	lhs := term(parser)

	condition := func() bool {
		one := canStep(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MINUS
		three := currentLexeme(parser).Type == lexer.LT_PLUS

		return (two || three)
	}

	for condition() {
		operator := currentLexeme(parser).Type
		//fmt.Println(operator)
		step(parser)
		rhs := factor(parser)
		lhs = expressionBinary(lhs, operator, rhs)
	}

	return lhs
}

func term(parser *Parser) *AST_Expression {

	//fmt.Println("Term")

	lhs := factor(parser)

	condition := func() bool {
		one := canStep(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MINUS
		three := currentLexeme(parser).Type == lexer.LT_PLUS

		return (two || three)
	}

	for condition() {
		operator := currentLexeme(parser).Type
		//fmt.Println(operator)
		step(parser)
		rhs := factor(parser)
		lhs = expressionBinary(lhs, operator, rhs)
	}

	return lhs
}

func factor(parser *Parser) *AST_Expression {

	//fmt.Println("Factor")

	lhs := unary(parser)

	condition := func() bool {
		one := canStep(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MULTIPLY
		three := currentLexeme(parser).Type == lexer.LT_DIVIDE

		return (two || three)
	}

	for condition() {
		operator := currentLexeme(parser).Type
		//fmt.Println(operator)
		step(parser)
		rhs := unary(parser)
		lhs = expressionBinary(lhs, operator, rhs)
	}

	return lhs
}

func unary(parser *Parser) *AST_Expression {

	//fmt.Println("Unary")

	if currentLexeme(parser).Type == lexer.LT_BANG || currentLexeme(parser).Type == lexer.LT_MINUS {
		operator := currentLexeme(parser).Type
		//fmt.Println(operator)
		step(parser)
		rhs := unary(parser)
		return expressionUnary(operator, rhs)
	}

	rhs := primary(parser)
	step(parser)
	return rhs
}

func canStep(parser *Parser) bool {
	c := parser.currentStep < len(parser.lexemes)
	/*fmt.Printf("Can step: ")
	fmt.Println(c)*/
	return c
}

func step(parser *Parser) {
	parser.currentStep += 1
}

func currentLexeme(parser *Parser) lexer.Lexeme {
	currLexeme := parser.lexemes[parser.currentStep]

	//fmt.Printf("Current Lexeme: %s\n", currLexeme.Label)

	return currLexeme
}

func nextLexeme(parser *Parser) lexer.Lexeme {
	return parser.lexemes[parser.currentStep+1]
}
