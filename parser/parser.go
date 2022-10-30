package parser

import (
	"fmt"
	"github.com/milansav/Castle/lexer"
	"log"
)

type Parser struct {
	lexemes     []lexer.Lexeme
	currentStep int
}

type ExpressionType int

const (
	ET_BINARY ExpressionType = iota
	ET_UNARY
	ET_LITERAL
	ET_GROUP
)

type AST_Expression struct {
	eType    ExpressionType
	lhs      *AST_Expression
	operator lexer.LexemeType
	rhs      *AST_Expression
	value    string
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
	return Parser{lexemes: lexer.Lexemes}
}

func Start(parser *Parser) []*AST_Expression {

	expressions := make([]*AST_Expression, 0)

	for canStep(parser) {

		/*if currentLexeme(parser).Type == lexer.LT_NUMBER {
			expr := expression(parser)
			expressions = append(expressions, expr)
			continue
		} else {
			step(parser)
		}*/

		expr := expression(parser)
		expressions = append(expressions, expr)
		continue
	}

	return expressions
}

/*

Grammar

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

	if c == lexer.LT_NUMBER {
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
