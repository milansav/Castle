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
	ET_EXPRESSION_ARRAY
	ET_FUNCTION_CALL
)

const (
	ST_STATEMENT_ARRAY StatementType = iota
	ST_STATEMENT
	ST_EXPRESSION
	ST_FUNCTION
	ST_DECLARATION
	ST_STRUCT
	ST_IF
)

type AST_Expression struct {
	EType         ExpressionType
	Lhs           *AST_Expression
	Operator      lexer.LexemeType
	Rhs           *AST_Expression
	Value         string
	RhsExpression *AST_Expression
	FunctionCall  *AST_FunctionCall
}

type AST_FunctionCall struct {
	name   string
	params []*AST_Expression
}

type AST_If struct {
	Condition  *AST_Expression
	Statements []*AST_Statement
}

type AST_Function struct {
	Name      string
	Props     []string
	Statement *AST_Statement
}

type AST_Declaration struct {
	Name  string
	Value *AST_Expression
}

type AST_Statement struct {
	SType       StatementType
	Statements  []*AST_Statement
	Statement   *AST_Statement
	Expression  *AST_Expression
	Function    *AST_Function
	Declaration *AST_Declaration
	If          *AST_If
}

type AST_Program struct {
	Statements []*AST_Statement
}

func createDeclarationNode(name string, value *AST_Expression) *AST_Declaration {
	decl := &AST_Declaration{
		Name:  name,
		Value: value,
	}

	return decl
}

func createFunctionNode(name string, props []string, statement *AST_Statement) *AST_Function {
	fn := &AST_Function{
		Name:      name,
		Props:     props,
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

func createExpressionLiteralNode(value string) *AST_Expression {
	expr := &AST_Expression{
		EType: ET_LITERAL,
		Value: value,
	}

	return expr
}

func createExpressionUnaryNode(operator lexer.LexemeType, rhs *AST_Expression) *AST_Expression {
	_rhs := &AST_Expression{
		EType:    rhs.EType,
		Lhs:      rhs.Lhs,
		Operator: rhs.Operator,
		Rhs:      rhs.Rhs,
		Value:    rhs.Value,
	}

	expr := &AST_Expression{
		EType:    ET_UNARY,
		Operator: operator,
		Rhs:      _rhs,
	}

	return expr
}

func createExpressionBinaryNode(lhs *AST_Expression, operator lexer.LexemeType, rhs *AST_Expression) *AST_Expression {
	_lhs := &AST_Expression{
		EType:    lhs.EType,
		Lhs:      lhs.Lhs,
		Operator: lhs.Operator,
		Rhs:      lhs.Rhs,
		Value:    lhs.Value,
	}

	_rhs := &AST_Expression{
		EType:    rhs.EType,
		Lhs:      rhs.Lhs,
		Operator: rhs.Operator,
		Rhs:      rhs.Rhs,
		Value:    rhs.Value,
	}

	expr := &AST_Expression{
		EType:    ET_BINARY,
		Lhs:      _lhs,
		Operator: operator,
		Rhs:      _rhs,
		Value:    "",
	}

	return expr
}

func createExpressionGroupNode(lhs *AST_Expression) *AST_Expression {

	expr := &AST_Expression{
		EType:    ET_GROUP,
		Lhs:      lhs,
		Operator: lexer.LT_NONE,
		Rhs:      nil,
		Value:    "",
	}

	return expr
}

func createExpressionFunctionCallNode(name string, params []*AST_Expression) *AST_Expression {
	expr := &AST_Expression{
		EType: ET_FUNCTION_CALL,
		FunctionCall: &AST_FunctionCall{
			name:   name,
			params: params,
		},
	}

	return expr
}

// TODO create dedicated AST printer
func PrintTree(tree *AST_Expression, depth int) {

	prefixChar := "≫ "
	prefix := ""

	log := func(_prefix string, field string) {
		fmt.Printf("%s%s[ %s ]%s\n", _prefix, util.Yellow, field, util.Reset)
	}

	for i := 0; i < depth; i++ {
		prefix += prefixChar
	}

	if tree.EType == ET_GROUP {
		log(prefix, "GROUP")

		//log(prefix, "LHS")

		PrintTree(tree.Lhs, depth+1)
	} else if tree.EType == ET_BINARY {

		switch tree.Operator {
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

		PrintTree(tree.Lhs, depth+1)

		//log(prefix, "RHS")

		PrintTree(tree.Rhs, depth+1)
	} else if tree.EType == ET_LITERAL {
		log(prefix, "VALUE")
		fmt.Println(prefix + prefixChar + tree.Value)
	}
}

func Create(lexer lexer.Lexer) Parser {
	return Parser{lexemes: lexer.Lexemes, currentLexeme: lexer.Lexemes[0], currentSym: lexer.Lexemes[0].Type}
}

func Start(parser *Parser) *AST_Program {

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

	program := createProgramNode()

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

	fmt.Println("statement")

	if accept(parser, lexer.LT_VAL) || accept(parser, lexer.LT_CONST) { // LET / CONST

		// variableType := prev(parser)

		if accept(parser, lexer.LT_MACRO) {
			// TODO Is macro
		}

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

					currentStatement.SType = ST_FUNCTION

					functionStatements := &AST_Statement{}

					if accept(parser, lexer.LT_LCURLY) { // LET / CONST {name} = ((params)) => { (statement) }

						functionStatements.SType = ST_STATEMENT_ARRAY
						functionStatements.Statements = make([]*AST_Statement, 0)

						for {
							if accept(parser, lexer.LT_RCURLY) {
								break
							}

							functionStatements.Statements = append(functionStatements.Statements, statement(parser))

						}

					} else {
						functionStatements.SType = ST_STATEMENT
						functionStatements.Statement = statement(parser) // LET / CONST {name} = ((params)) => {statement}
					}

					currentStatement.Function = createFunctionNode(identifier.Label, params, functionStatements)
					expect(parser, lexer.LT_SEMICOLON) // LET / CONST {name} = ((params)) => {statement};

					return currentStatement
				} else {

					fmt.Println("Declaration")

					expr := expression(parser)

					decl := createDeclarationNode(identifier.Label, expr)

					currentStatement.SType = ST_DECLARATION
					currentStatement.Declaration = decl

					expect(parser, lexer.LT_SEMICOLON)
					return currentStatement
				}
			}
		}
	} else if accept(parser, lexer.LT_IF) { // IF
		expect(parser, lexer.LT_LPAREN)

		currentStatement.SType = ST_IF

		expr := expression(parser)

		condition := &AST_If{
			Condition:  expr,
			Statements: make([]*AST_Statement, 0),
		}

		expect(parser, lexer.LT_RPAREN)

		if accept(parser, lexer.LT_LCURLY) {
			for {
				if accept(parser, lexer.LT_RCURLY) {
					break
				}

				condition.Statements = append(condition.Statements, statement(parser))
			}
		}

		currentStatement.If = condition

		return currentStatement

	} else if accept(parser, lexer.LT_RETURN) { // RETURN
		currentStatement.SType = ST_EXPRESSION
		currentStatement.Expression = expression(parser)

		return currentStatement
	} else {
		fmt.Println("Caught expression")
		expr := expression(parser)

		currentStatement.SType = ST_EXPRESSION
		currentStatement.Expression = expr

		return currentStatement
	}

	log.Panic("Unexpected path")
	return nil
}

/*

Expression Grammar

expression -> compare (LT_COMMA compare)*

compare -> term (( LT_LESS | LT_GREATER | LT_LEQ | LT_LGE | LT_EQ | LT_NE ) term)*

term -> factor (( LT_PLUS | LT_MINUS ) factor)*

factor -> unary (( LT_DIVIDE | LT_MULTIPLY ) unary)*

unary -> ( LT_BANG | LT_MINUS | LT_PLUS ) unary | primary

primary -> LT_NUMBER | LT_FLOAT | LT_LPAREN expression LT_RPAREN | LT_IDENTIFIER | "function call"
*/

func expression(parser *Parser) *AST_Expression {

	fmt.Println("Expression")

	lhs := compareOR(parser)

	if accept(parser, lexer.LT_COMMA) {
		expression := &AST_Expression{
			EType:         ET_EXPRESSION_ARRAY,
			Lhs:           lhs,
			RhsExpression: expression(parser),
		}

		return expression
	}

	return lhs
}

func compareOR(parser *Parser) *AST_Expression {

	fmt.Println("compareOR")

	lhs := compareAND(parser)

	for accept(parser, lexer.LT_OR) || accept(parser, lexer.LT_NOR) || accept(parser, lexer.LT_XOR) || accept(parser, lexer.LT_XNOR) {
		operator := prev(parser).Type

		rhs := compareAND(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func compareAND(parser *Parser) *AST_Expression {

	fmt.Println("compareAND")

	lhs := compareNEQEQ(parser)

	for accept(parser, lexer.LT_AND) || accept(parser, lexer.LT_NAND) || accept(parser, lexer.LT_XAND) || accept(parser, lexer.LT_XNAND) {
		operator := prev(parser).Type

		rhs := compareNEQEQ(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func compareNEQEQ(parser *Parser) *AST_Expression {

	fmt.Println("compareNEQEQ")

	lhs := compareLEQGEQLTGT(parser)

	for accept(parser, lexer.LT_NEQ) || accept(parser, lexer.LT_EQ) {
		operator := prev(parser).Type

		rhs := compareLEQGEQLTGT(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func compareLEQGEQLTGT(parser *Parser) *AST_Expression {

	fmt.Println("compareLEQGEQLTGT")

	lhs := term(parser)

	for accept(parser, lexer.LT_LEQ) || accept(parser, lexer.LT_GEQ) || accept(parser, lexer.LT_LCHEVRON) || accept(parser, lexer.LT_RCHEVRON) {
		operator := prev(parser).Type

		rhs := term(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func term(parser *Parser) *AST_Expression {

	fmt.Println("Term")

	lhs := factor(parser)

	fmt.Printf("term %s\n", curr(parser).Label)

	for accept(parser, lexer.LT_MINUS) || accept(parser, lexer.LT_PLUS) {
		operator := prev(parser).Type

		rhs := factor(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func factor(parser *Parser) *AST_Expression {

	fmt.Println("Factor")

	lhs := unary(parser)

	fmt.Printf("factor %s\n", curr(parser).Label)

	for accept(parser, lexer.LT_MULTIPLY) || accept(parser, lexer.LT_DIVIDE) {
		operator := prev(parser).Type

		fmt.Println("inside")

		rhs := unary(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func unary(parser *Parser) *AST_Expression {

	fmt.Println("Unary")

	if accept(parser, lexer.LT_BANG) || accept(parser, lexer.LT_MINUS) {
		operator := prev(parser).Type

		rhs := unary(parser)
		return createExpressionUnaryNode(operator, rhs)
	}

	rhs := primary(parser)
	return rhs
}

func primary(parser *Parser) *AST_Expression {

	fmt.Println("Primary")

	c := currentLexeme(parser).Type

	fmt.Println(lexer.LexemeTypeLabels[c])

	if accept(parser, lexer.LT_LITERAL_NUMBER) || accept(parser, lexer.LT_LITERAL_FLOAT) || accept(parser, lexer.LT_LITERAL_STRING) || accept(parser, lexer.LT_LITERAL_BOOL) {
		rhs := currentLexeme(parser).Label
		fmt.Println(rhs)

		expr := createExpressionLiteralNode(rhs)
		return expr
	} else if accept(parser, lexer.LT_IDENTIFIER) {
		name := prev(parser).Label

		expressions := make([]*AST_Expression, 0)

		if accept(parser, lexer.LT_LPAREN) {
			for {
				if accept(parser, lexer.LT_RPAREN) {
					break
				}
				expressions = append(expressions, expression(parser))

				accept(parser, lexer.LT_COMMA)
			}

			expr := createExpressionFunctionCallNode(name, expressions)

			return expr
		}

		expr := createExpressionLiteralNode(name)
		return expr
	} else if accept(parser, lexer.LT_LPAREN) {
		expr := expression(parser)

		expr = createExpressionGroupNode(expr)

		expect(parser, lexer.LT_RPAREN)

		fmt.Printf("After group, current lexeme: %s\n", currentLexeme(parser).Label)
		return expr
	} else {
		log.Panic("Unexpected path")
		return nil
	}
}

func currentLexeme(parser *Parser) lexer.Lexeme {
	currLexeme := parser.lexemes[parser.currentStep]

	//fmt.Printf("Current Lexeme: %s\n", currLexeme.Label)

	return currLexeme
}
