package parser

import (
	"errors"
	"fmt"
	"log"

	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/util"
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

func prev(parser *Parser) lexer.Lexeme {
	if parser.currentStep-1 > 0 {
		return parser.lexemes[parser.currentStep-1]
	}

	return lexer.Lexeme{Type: lexer.LT_NONE}

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

func is(parser *Parser, symbol lexer.LexemeType) bool {
	if parser.currentSym == symbol {
		fmt.Printf("Accepted symbol: %s\n", lexer.LexemeTypeLabels[parser.currentSym])
		return true
	}

	return false
}

func expect(parser *Parser, symbol lexer.LexemeType) bool {

	expectingMessage := fmt.Sprintf(
		"%sExpecting symbol: %s%s\n",
		util.Green,
		lexer.LexemeTypeLabels[symbol],
		util.Reset)

	fmt.Println(expectingMessage)

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
	ST_STATEMENT_ARRAY StatementType = iota
	ST_STATEMENT
	ST_EXPRESSION
	ST_FUNCTION
	ST_DECLARATION
	ST_STRUCT
)

type AST_Expression struct {
	eType    ExpressionType
	lhs      *AST_Expression
	operator lexer.LexemeType
	rhs      *AST_Expression
	value    string
}

type AST_Function struct {
	name      string
	props     []string
	Statement *AST_Statement
}

type AST_Declaration struct {
	name  string
	value *AST_Expression
}

type AST_Statement struct {
	sType       StatementType
	Statements  []*AST_Statement
	Statement   *AST_Statement
	Expression  *AST_Expression
	Function    *AST_Function
	Declaration *AST_Declaration
}

type AST_Program struct {
	Statements []*AST_Statement
}

func createFunctionNode(name string, props []string, statement *AST_Statement) *AST_Function {
	fn := &AST_Function{
		name:      name,
		props:     props,
		Statement: statement,
	}

	return fn
}

func createProgramNode() *AST_Program {
	prog := &AST_Program{
		Statements: make([]*AST_Statement, 0),
	}

	return prog
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

	log := func(_prefix string, field string) {
		fmt.Printf("%s%s[ %s ]%s\n", _prefix, util.Yellow, field, util.Reset)
	}

	for i := 0; i < depth; i++ {
		prefix += prefixChar
	}

	if tree.eType == ET_GROUP {
		log(prefix, "GROUP")

		//log(prefix, "LHS")

		PrintTree(tree.lhs, depth+1)
	} else if tree.eType == ET_BINARY {

		switch tree.operator {
		case lexer.LT_PLUS:
			log(prefix, "ADD")
		case lexer.LT_MINUS:
			log(prefix, "SUBTRACT")
		case lexer.LT_MULTIPLY:
			log(prefix, "MULTIPLY")
		case lexer.LT_DIVIDE:
			log(prefix, "DIVIDE")
		}

		//log(prefix, "LHS")

		PrintTree(tree.lhs, depth+1)

		//log(prefix, "RHS")

		PrintTree(tree.rhs, depth+1)
	} else if tree.eType == ET_LITERAL {
		log(prefix, "VALUE")
		fmt.Println(prefix + prefixChar + tree.value)
	}
}

func Create(lexer lexer.Lexer) Parser {
	return Parser{lexemes: lexer.Lexemes, currentLexeme: lexer.Lexemes[0], currentSym: lexer.Lexemes[0].Type}
}

func StartExpressionParser(parser *Parser) []*AST_Expression {

	expressions := make([]*AST_Expression, 0)

	for hasNext(parser) {

		if currentLexeme(parser).Type == lexer.LT_NUMBER {
			expr := expression(parser)
			expressions = append(expressions, expr)
			continue
		} else {
			next(parser)
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

		// fmt.Print("Current token: ")
		// fmt.Println(lexer.LexemeTypeLabels[parser.currentLexeme.Type])

		sttmnt := statement(parser)

		program.Statements = append(program.Statements, sttmnt)

		continue
	}

	if accept(parser, lexer.LT_END) {
		fmt.Println("Finished")
	}

	return program
}

func statement(parser *Parser) *AST_Statement {

	currentStatement := &AST_Statement{}

	if accept(parser, lexer.LT_VAL) || accept(parser, lexer.LT_CONST) { // LET / CONST

		// variableType := prev(parser)

		// Declare variable type here to be used later
		if expect(parser, lexer.LT_IDENTIFIER) { // LET / CONST {name}

			identifier := prev(parser)

			if expect(parser, lexer.LT_EQUALS) { // LET / CONST {name} =
				if accept(parser, lexer.LT_LPAREN) { // LET / CONST {name} = (

					params := make([]string, 0)

					// Parse function
					for { // n1, n2, .. nx )
						if accept(parser, lexer.LT_IDENTIFIER) {

							params = append(params, prev(parser).Label)

							if accept(parser, lexer.LT_RPAREN) {
								break
							} else {
								expect(parser, lexer.LT_COMMA)
							}
						}
					}

					expect(parser, lexer.LT_LAMBDA) // LET / CONST {name} = ((params)) =>

					currentStatement.sType = ST_FUNCTION

					functionStatements := &AST_Statement{sType: ST_STATEMENT_ARRAY, Statements: make([]*AST_Statement, 0)}

					if accept(parser, lexer.LT_LCURLY) { // LET / CONST {name} = ((params)) => { (statement) }

						for {
							if accept(parser, lexer.LT_RCURLY) {
								break
							}

							functionStatements.Statements = append(functionStatements.Statements, statement(parser))

						}

					} else {
						functionStatements.sType = ST_STATEMENT
						functionStatements.Statement = statement(parser) // LET / CONST {name} = ((params)) => {statement}
					}

					createFunctionNode(identifier.Label, params, functionStatements)
					expect(parser, lexer.LT_SEMICOLON) // LET / CONST {name} = ((params)) => {statement};

					return currentStatement
				} else if accept(parser, lexer.LT_STRING) {
					fmt.Println("Hello")
					expect(parser, lexer.LT_SEMICOLON)
				} else {
					expression(parser)
					expect(parser, lexer.LT_SEMICOLON)
				}
			}
		}
	} else if accept(parser, lexer.LT_IF) { // IF
		expect(parser, lexer.LT_LPAREN)

		condition(parser)

		expect(parser, lexer.LT_RPAREN)

		if accept(parser, lexer.LT_LCURLY) {
			for {
				if accept(parser, lexer.LT_RCURLY) {
					break
				}

				currentStatement.Statements = append(currentStatement.Statements, statement(parser))
			}
		}
	}
	return &AST_Statement{}
}

// logicop -> LT_AND | LT_OR | LT_NOR | LT_NAND | LT_XOR | LT_XAND | LT_XNOR | LT_XNAND | LT_
// condition -> expression (logicop expression)*
func condition(parser *Parser) {
	expression(parser)

	isLogicop := func(c lexer.LexemeType) bool {
		switch c {
		case lexer.LT_AND:
			fallthrough
		case lexer.LT_OR:
			fallthrough
		case lexer.LT_NAND:
			fallthrough
		case lexer.LT_NOR:
			fallthrough
		case lexer.LT_XAND:
			fallthrough
		case lexer.LT_XOR:
			fallthrough
		case lexer.LT_XNAND:
			fallthrough
		case lexer.LT_XNOR:
			return true

		default:
			return false
		}
	}

	for {
		if !isLogicop(curr(parser).Type) {
			break
		}

		expression(parser)
	}
}

/*

Expression Grammar

expression -> term

//                                                                                 TODO 👇
primary -> LT_NUMBER | LT_FLOAT | LT_LPAREN expression LT_RPAREN | LT_IDENTIFIER | "function call"

term -> factor (( LT_PLUS | LT_MINUS ) factor)*

factor -> unary (( LT_DIVIDE | LT_MULTIPLY ) unary)*

unary -> ( LT_BANG | LT_MINUS | LT_PLUS ) unary | primary
*/

func primary(parser *Parser) *AST_Expression {
	// fmt.Println("Primary")

	c := currentLexeme(parser).Type

	fmt.Println(lexer.LexemeTypeLabels[c])

	if c == lexer.LT_NUMBER || c == lexer.LT_FLOAT {
		rhs := currentLexeme(parser).Label
		//fmt.Println(rhs)

		expr := expressionLiteral(rhs)
		return expr
	} else if c == lexer.LT_IDENTIFIER {
		name := currentLexeme(parser).Label

		if accept(parser, lexer.LT_LPAREN) {
			for {
				if accept(parser, lexer.LT_RPAREN) {
					break
				}
				expect(parser, lexer.LT_IDENTIFIER)
				accept(parser, lexer.LT_COMMA)
			}
		}

		expr := expressionLiteral(name)
		return expr
	} else if c == lexer.LT_LPAREN {
		next(parser)
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

	// fmt.Println("Expression")

	lhs := term(parser)

	condition := func() bool {
		one := hasNext(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MINUS
		three := currentLexeme(parser).Type == lexer.LT_PLUS

		return (two || three)
	}

	for condition() {
		operator := currentLexeme(parser).Type
		// fmt.Println(operator)
		next(parser)
		rhs := factor(parser)
		lhs = expressionBinary(lhs, operator, rhs)
	}

	return lhs
}

func term(parser *Parser) *AST_Expression {

	// fmt.Println("Term")

	lhs := factor(parser)

	condition := func() bool {
		one := hasNext(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MINUS
		three := currentLexeme(parser).Type == lexer.LT_PLUS

		return (two || three)
	}

	for condition() {
		operator := currentLexeme(parser).Type
		// fmt.Println(operator)
		next(parser)
		rhs := factor(parser)
		lhs = expressionBinary(lhs, operator, rhs)
	}

	return lhs
}

func factor(parser *Parser) *AST_Expression {

	// fmt.Println("Factor")

	lhs := unary(parser)

	isCondition := func() bool {
		one := hasNext(parser)
		if !one {
			return false
		}
		two := currentLexeme(parser).Type == lexer.LT_MULTIPLY
		three := currentLexeme(parser).Type == lexer.LT_DIVIDE

		return (two || three)
	}

	for isCondition() {
		operator := currentLexeme(parser).Type
		//fmt.Println(operator)
		next(parser)
		rhs := unary(parser)
		lhs = expressionBinary(lhs, operator, rhs)
	}

	return lhs
}

func unary(parser *Parser) *AST_Expression {

	// fmt.Println("Unary")

	if currentLexeme(parser).Type == lexer.LT_BANG || currentLexeme(parser).Type == lexer.LT_MINUS {
		operator := currentLexeme(parser).Type
		//fmt.Println(operator)
		next(parser)
		rhs := unary(parser)
		return expressionUnary(operator, rhs)
	}

	rhs := primary(parser)
	next(parser)
	return rhs
}

func currentLexeme(parser *Parser) lexer.Lexeme {
	currLexeme := parser.lexemes[parser.currentStep]

	//fmt.Printf("Current Lexeme: %s\n", currLexeme.Label)

	return currLexeme
}
