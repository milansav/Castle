package astprinter

import (
	"fmt"

	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/parser"
	"github.com/milansav/Castle/util"
)

type ASTPrinter struct {
	indentation int
}

func Prefix(printer *ASTPrinter) string {
	prefix := ""
	// prefixChar := "≫ "
	prefixChar := "  "

	for i := 0; i < printer.indentation; i++ {
		prefix += prefixChar
	}

	return prefix
}

func (printer *ASTPrinter) Group(name string) {
	fmt.Printf("%s%s[ %s ]%s\n", Prefix(printer), util.Yellow, name, util.Reset)
}

func (printer *ASTPrinter) Value(name string, value string) {
	printer.In()
	fmt.Printf("%s- %s%s: %s%s\n", Prefix(printer), util.Yellow, name, util.Reset, value)
	printer.Out()
}

func (printer *ASTPrinter) Info(description string) {
	printer.In()
	fmt.Printf("%s- %s%s%s\n", Prefix(printer), util.Yellow, description, util.Reset)
	printer.Out()

}

func (printer *ASTPrinter) In() {
	printer.indentation++
}

func (printer *ASTPrinter) Out() {
	printer.indentation--
}

func PrintAST(program *parser.AST_Program) {
	printer := &ASTPrinter{
		indentation: 0,
	}

	printer.Group("Program")
	printer.In()

	for _, value := range program.Statements {
		printer.PrintStatement(value)
	}

	printer.Out()
}

func (printer *ASTPrinter) PrintStatement(statement *parser.AST_Statement) {
	// Group(printer, "Statement")
	// Value(printer, "type", parser.StatementTypeLabels[statement.SType])

	switch statement.SType {
	case parser.ST_STATEMENT_ARRAY:
		// Group(printer, "Block")
		printer.In()

		for _, value := range statement.Statements {
			printer.PrintStatement(value)
		}

		printer.Out()
	case parser.ST_STATEMENT:
		printer.Group("Statement")
		printer.In()

		printer.PrintStatement(statement.Statement)

		printer.Out()
	case parser.ST_EXPRESSION:
		// Group(printer, "Expression")
		// printer.indentation++
		printer.PrintExpression(statement.Expression)
		// printer.indentation--
	case parser.ST_FUNCTION:
		printer.Group("Function")

		printer.In()

		printer.PrintFunction(statement.Function)

		printer.Out()

	case parser.ST_DECLARATION:
		printer.Group("Declaration")
		printer.In()
		PrintDeclaration(printer, statement.Declaration)
		printer.Out()
	case parser.ST_STRUCT:
		printer.Group("Struct")
	case parser.ST_IF:
		printer.Group("If")

		printer.Info("Condition")

		printer.In()

		printer.PrintExpression(statement.If.Condition)

		printer.Out()

		printer.Info("Body")

		printer.In()

		for _, value := range statement.If.Statements {
			printer.PrintStatement(value)
		}

		printer.Out()
	case parser.ST_RETURN:
		printer.Group("Return")
		printer.Info("Value")
		printer.In()
		printer.PrintExpression(statement.Expression)
		printer.Out()
	}
}

func (printer *ASTPrinter) PrintFunction(function *parser.AST_Function) {
	printer.Value("Name", function.Name)
	printer.Info("Args")
	printer.In()

	for _, value := range function.Props {
		printer.Info(value)
	}

	printer.Out()

	printer.Info("Body")

	printer.In()

	printer.PrintStatement(function.Statement)

	printer.Out()
}

func PrintDeclaration(printer *ASTPrinter, declaration *parser.AST_Declaration) {
	printer.Value("Name", declaration.Name)
	printer.Info("Value")

	printer.In()
	printer.PrintExpression(declaration.Value)
	printer.Out()

}

func (printer *ASTPrinter) PrintExpression(expression *parser.AST_Expression) {
	// Group(printer, "Expression")

	switch expression.EType {
	case parser.ET_GROUP:
		printer.Group("Group")
		printer.In()
		printer.PrintExpression(expression.Lhs)
		printer.Out()
	case parser.ET_BINARY:
		switch expression.Operator {
		case lexer.LT_PLUS:
			printer.Group("ADD")
		case lexer.LT_MINUS:
			printer.Group("SUBTRACT")
		case lexer.LT_MULTIPLY:
			printer.Group("MULTIPLY")
		case lexer.LT_DIVIDE:
			printer.Group("DIVIDE")
		case lexer.LT_MODULO:
			printer.Group("MODULO")
		case lexer.LT_POWER:
			printer.Group("POWER")
		case lexer.LT_EQ:
			printer.Group("EQUALS")
		case lexer.LT_NEQ:
			printer.Group("NOT EQUALS")
		case lexer.LT_GEQ:
			printer.Group("GREATER EQUAL")
		case lexer.LT_LEQ:
			printer.Group("LESS EQUAL")
		case lexer.LT_AND:
			printer.Group("AND")
		case lexer.LT_OR:
			printer.Group("OR")
		case lexer.LT_NAND:
			printer.Group("NAND")
		case lexer.LT_NOR:
			printer.Group("NOR")
		case lexer.LT_XAND:
			printer.Group("XAND")
		case lexer.LT_XOR:
			printer.Group("XOR")
		case lexer.LT_XNAND:
			printer.Group("XNAND")
		case lexer.LT_XNOR:
			printer.Group("XNOR")
		case lexer.LT_LCHEVRON:
			printer.Group("LESS THAN")
		case lexer.LT_RCHEVRON:
			printer.Group("GREATER THAN")
		default:
			printer.Group("UNKNOWN")
		}
		printer.In()
		printer.PrintExpression(expression.Lhs)
		printer.Out()

		printer.In()
		printer.PrintExpression(expression.Rhs)
		printer.Out()
	case parser.ET_LITERAL:
		printer.Group("Literal")
		printer.Value("Value", expression.Literal.Value)
		printer.Value("Type", parser.LiteralTypeLabels[expression.Literal.Type])
	case parser.ET_IDENTIFIER:
		printer.Group("Identifier")
		printer.Value("Name", expression.Value)
	case parser.ET_FUNCTION_CALL:
		printer.Group("Call")
		printer.Value("Name", expression.FunctionCall.Name)

		printer.Info("Args")
		printer.In()
		for index, value := range expression.FunctionCall.Params {
			printer.Info(fmt.Sprintf("%d", index))
			printer.In()
			printer.PrintExpression(value)
			printer.Out()
		}
		printer.Out()
	case parser.ET_EXPRESSION_ARRAY:
		printer.Group("Expressions")
		printer.In()
		printer.PrintExpression(expression.Lhs)
		printer.PrintExpression(expression.Rhs)
		printer.Out()
	}
}
