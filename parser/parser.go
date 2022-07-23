package parser

import (
	"github.com/milansav/Castle/lexer"
)

type Parser struct {
	lexemes []lexer.Lexeme
}

func Create(lexer lexer.Lexer) Parser {
	return Parser{lexemes: lexer.Lexemes}
}

func Start(parser *Parser) {
	
}

func expression(parser *Parser) {

}

func term(parser *Parser) {

}

func factor(parser *Parser) {

}
