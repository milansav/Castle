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
	ET_MEMBER_ACCESS
)

var ExpressionTypeLabels = map[ExpressionType]string{
	ET_BINARY:           "ET_BINARY",
	ET_UNARY:            "ET_UNARY",
	ET_LITERAL:          "ET_LITERAL",
	ET_GROUP:            "ET_GROUP",
	ET_EXPRESSION_ARRAY: "ET_EXPRESSION_ARRAY",
	ET_FUNCTION_CALL:    "ET_FUNCTION_CALL",
	ET_MEMBER_ACCESS:    "ET_MEMBER_ACCESS",
}

const (
	ST_STATEMENT_ARRAY StatementType = iota
	ST_STATEMENT
	ST_EXPRESSION
	ST_FUNCTION
	ST_DECLARATION
	ST_STRUCT
	ST_IF
	ST_RETURN
)

var StatementTypeLabels = map[StatementType]string{
	ST_STATEMENT_ARRAY: "ST_STATEMENT_ARRAY",
	ST_STATEMENT:       "ST_STATEMENT",
	ST_EXPRESSION:      "ST_EXPRESSION",
	ST_FUNCTION:        "ST_FUNCTION",
	ST_DECLARATION:     "ST_DECLARATION",
	ST_STRUCT:          "ST_STRUCT",
	ST_IF:              "ST_IF",
	ST_RETURN:          "ST_RETURN",
}

type AST_Expression struct {
	EType         ExpressionType
	Lhs           *AST_Expression
	Operator      lexer.LexemeType
	Rhs           *AST_Expression
	Value         string
	RhsExpression *AST_Expression
	FunctionCall  *AST_FunctionCall
	Member        *AST_Expression
}

type AST_FunctionCall struct {
	Name   string
	Params []*AST_Expression
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
			Name:   name,
			Params: params,
		},
	}

	return expr
}

func createExpressionMemberAccessNode(lhs *AST_Expression, member *AST_Expression) *AST_Expression {
	expr := &AST_Expression{
		EType:  ET_MEMBER_ACCESS,
		Lhs:    lhs,
		Member: member,
	}

	return expr
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
		currentStatement.SType = ST_RETURN
		currentStatement.Expression = expression(parser)
		expect(parser, lexer.LT_SEMICOLON)

		return currentStatement
	} else {
		fmt.Println("Caught expression")
		expr := expression(parser)

		currentStatement.SType = ST_EXPRESSION
		currentStatement.Expression = expr

		expect(parser, lexer.LT_SEMICOLON)

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

// Same effect as calling compareOR, prevents catching LT_COMMA inside arguments etc..
func safeExpression(parser *Parser) *AST_Expression {
	return compareOR(parser)
}

func compareOR(parser *Parser) *AST_Expression {

	lhs := compareAND(parser)

	for accept(parser, lexer.LT_OR) || accept(parser, lexer.LT_NOR) || accept(parser, lexer.LT_XOR) || accept(parser, lexer.LT_XNOR) {
		operator := prev(parser).Type

		rhs := compareAND(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func compareAND(parser *Parser) *AST_Expression {

	lhs := compareNEQEQ(parser)

	for accept(parser, lexer.LT_AND) || accept(parser, lexer.LT_NAND) || accept(parser, lexer.LT_XAND) || accept(parser, lexer.LT_XNAND) {
		operator := prev(parser).Type

		rhs := compareNEQEQ(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func compareNEQEQ(parser *Parser) *AST_Expression {

	lhs := compareLEQGEQLTGT(parser)

	for accept(parser, lexer.LT_NEQ) || accept(parser, lexer.LT_EQ) {
		operator := prev(parser).Type

		rhs := compareLEQGEQLTGT(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func compareLEQGEQLTGT(parser *Parser) *AST_Expression {

	lhs := term(parser)

	for accept(parser, lexer.LT_LEQ) || accept(parser, lexer.LT_GEQ) || accept(parser, lexer.LT_LCHEVRON) || accept(parser, lexer.LT_RCHEVRON) {
		operator := prev(parser).Type

		rhs := term(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func term(parser *Parser) *AST_Expression {

	lhs := factor(parser)

	for accept(parser, lexer.LT_MINUS) || accept(parser, lexer.LT_PLUS) {
		operator := prev(parser).Type

		rhs := factor(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func factor(parser *Parser) *AST_Expression {

	lhs := unary(parser)

	for accept(parser, lexer.LT_MULTIPLY) || accept(parser, lexer.LT_DIVIDE) {
		operator := prev(parser).Type

		rhs := unary(parser)
		lhs = createExpressionBinaryNode(lhs, operator, rhs)
	}

	return lhs
}

func unary(parser *Parser) *AST_Expression {

	if accept(parser, lexer.LT_BANG) || accept(parser, lexer.LT_MINUS) {
		operator := prev(parser).Type

		rhs := unary(parser)
		return createExpressionUnaryNode(operator, rhs)
	}

	rhs := memberAccess(parser)
	return rhs
}

// memberAccess -> primary ((LT_PERIOD primary)*)
func memberAccess(parser *Parser) *AST_Expression {
	lhs := primary(parser)

	for {
		if accept(parser, lexer.LT_PERIOD) {
			member := primary(parser)

			lhs = createExpressionMemberAccessNode(lhs, member)

			continue
		}

		break
	}

	return lhs
}

func primary(parser *Parser) *AST_Expression {

	if accept(parser, lexer.LT_LITERAL_NUMBER) || accept(parser, lexer.LT_LITERAL_FLOAT) || accept(parser, lexer.LT_LITERAL_STRING) || accept(parser, lexer.LT_LITERAL_BOOL) {
		rhs := prev(parser).Label

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
				expressions = append(expressions, safeExpression(parser))

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

		// fmt.Printf("After group, current lexeme: %s\n", currentLexeme(parser).Label)
		return expr
	} else {
		log.Panic("Unexpected path")
		return nil
	}
}
