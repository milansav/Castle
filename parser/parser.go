package parser

import (
	"fmt"
	"github.com/milansav/Castle/lexer"
)

type Parser struct {
	lexemes []lexer.Lexeme
	currentStep int
}

type ExpressionType int

const (

	ET_BINARY ExpressionType = iota
)

func Create(lexer lexer.Lexer) Parser {
	return Parser{lexemes: lexer.Lexemes}
}

func Start(parser *Parser) {
	for(canStep(parser)) {

		if(currentLexeme(parser).Type == lexer.LT_NUMBER) {
			expression(parser)
			continue
		} else {
			step(parser)
		}
	}
}

/*

Grammar 

expression -> term

primary -> LT_NUMBER | LT_STRING | LT_TRUE | LT_FALSE | LT_NULL 
				| LT_LPAREN expression LT_RPAREN

term -> factor (( LT_PLUS | LT_MINUS ) factor)*

factor -> unary (( LT_DIVIDE | LT_MULTIPLY ) unary)*

unary -> ( LT_NEG | LT_MINUS | LT_PLUS ) unary | primary
*/

func primary(parser *Parser) {
	rhs := currentLexeme(parser).Label
	fmt.Println(rhs)
}

func expression(parser *Parser) {

	fmt.Println("Expression")

	/*lhs := */term(parser)

	condition := func() bool {
		one := canStep(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MINUS 
		three := currentLexeme(parser).Type == lexer.LT_PLUS

		return (two || three)
	}

	for(condition()) {
		operator := currentLexeme(parser).Label
		fmt.Println(operator)
		step(parser)
		/*rhs := */factor(parser)
	}
}

func term(parser *Parser) {

	fmt.Println("Term")
	
	/*lhs := */factor(parser)

	condition := func() bool {
		one := canStep(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MINUS 
		three := currentLexeme(parser).Type == lexer.LT_PLUS

		return (two || three)
	}

	for(condition()) {
		operator := currentLexeme(parser).Label
		fmt.Println(operator)
		step(parser)
		/*rhs := */factor(parser)
	}
}

func factor(parser *Parser) {

	fmt.Println("Factor")
	
	/*lhs := */unary(parser)

	condition := func() bool {
		one := canStep(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MULTIPLY 
		three := currentLexeme(parser).Type == lexer.LT_DIVIDE

		return (two || three)
	}

	for(condition()) {
		operator := currentLexeme(parser).Label
		fmt.Println(operator)
		step(parser)
		/*rhs := */unary(parser)
	}
}

func unary(parser *Parser) {

	fmt.Println("Unary")
	
	if(currentLexeme(parser).Type == lexer.LT_BANG || currentLexeme(parser).Type == lexer.LT_MINUS) {
		operator := currentLexeme(parser).Label
		fmt.Println(operator)
		step(parser)
		/*rhs := */unary(parser)
	}

	/*return */primary(parser)
	step(parser)
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
	return parser.lexemes[parser.currentStep]
}

func nextLexeme(parser *Parser) lexer.Lexeme {
	return parser.lexemes[parser.currentStep+1]
}
