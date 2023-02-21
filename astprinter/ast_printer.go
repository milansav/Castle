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

func Group(printer *ASTPrinter, name string) {
	fmt.Printf("%s%s[ %s ]%s\n", Prefix(printer), util.Yellow, name, util.Reset)
}

func Value(printer *ASTPrinter, name string, value string) {
	printer.indentation++
	fmt.Printf("%s- %s%s: %s%s\n", Prefix(printer), util.Yellow, name, util.Reset, value)
	printer.indentation--
}

func Info(printer *ASTPrinter, description string) {
	printer.indentation++
	fmt.Printf("%s- %s%s%s\n", Prefix(printer), util.Yellow, description, util.Reset)
	printer.indentation--

}

func PrintAST(program *parser.AST_Program) {
	printer := &ASTPrinter{
		indentation: 0,
	}

	Group(printer, "Program")
	printer.indentation++

	for _, value := range program.Statements {
		PrintStatement(printer, value)
	}

	printer.indentation--
}

func PrintStatement(printer *ASTPrinter, statement *parser.AST_Statement) {
	// Group(printer, "Statement")
	// Value(printer, "type", parser.StatementTypeLabels[statement.SType])

	switch statement.SType {
	case parser.ST_STATEMENT_ARRAY:
		// Group(printer, "Block")
		printer.indentation++

		for _, value := range statement.Statements {
			PrintStatement(printer, value)
		}

		printer.indentation--
	case parser.ST_STATEMENT:
		Group(printer, "Statement")
	case parser.ST_EXPRESSION:
		Group(printer, "Expression")
		printer.indentation++
		PrintExpression(printer, statement.Expression)
		printer.indentation--
	case parser.ST_FUNCTION:
		Group(printer, "Function")

		printer.indentation++

		PrintFunction(printer, statement.Function)

		printer.indentation--

	case parser.ST_DECLARATION:
		Group(printer, "Declaration")
		printer.indentation++
		PrintDeclaration(printer, statement.Declaration)
		printer.indentation--
	case parser.ST_STRUCT:
		Group(printer, "Struct")
	case parser.ST_IF:
		Group(printer, "If")
	}
}

func PrintFunction(printer *ASTPrinter, function *parser.AST_Function) {
	Value(printer, "Name", function.Name)
	Info(printer, "Args")
	printer.indentation++

	for _, value := range function.Props {
		Info(printer, value)
	}

	printer.indentation--

	Info(printer, "Body")

	printer.indentation++

	PrintStatement(printer, function.Statement)

	printer.indentation--
}

func PrintDeclaration(printer *ASTPrinter, declaration *parser.AST_Declaration) {
	Value(printer, "Name", declaration.Name)
	Info(printer, "Value")

	printer.indentation++
	PrintExpression(printer, declaration.Value)
	printer.indentation--

}

func PrintExpression(printer *ASTPrinter, expression *parser.AST_Expression) {
	// Group(printer, "Expression")

	switch expression.EType {
	case parser.ET_GROUP:
		Group(printer, "Group")
		printer.indentation++
		PrintExpression(printer, expression.Lhs)
		printer.indentation--
	case parser.ET_BINARY:
		switch expression.Operator {
		case lexer.LT_PLUS:
			Group(printer, "ADD")
		case lexer.LT_MINUS:
			Group(printer, "SUBTRACT")
		case lexer.LT_MULTIPLY:
			Group(printer, "MULTIPLY")
		case lexer.LT_DIVIDE:
			Group(printer, "DIVIDE")
		default:
			Group(printer, "UNKNOWN")
		}
		printer.indentation++
		PrintExpression(printer, expression.Lhs)
		printer.indentation--

		printer.indentation++
		PrintExpression(printer, expression.Rhs)
		printer.indentation--
	case parser.ET_LITERAL:
		Group(printer, "Literal")
		Value(printer, "Value", expression.Value)
	case parser.ET_FUNCTION_CALL:
		Group(printer, "Call")
		Value(printer, "Name", expression.FunctionCall.Name)

		Info(printer, "Args")
		printer.indentation++
		for index, value := range expression.FunctionCall.Params {
			Info(printer, fmt.Sprintf("%d", index))
			printer.indentation++
			PrintExpression(printer, value)
			printer.indentation--
		}
		printer.indentation--
	case parser.ET_EXPRESSION_ARRAY:
		Group(printer, "Expressions")
		printer.indentation++
		if expression.Lhs != nil {
			PrintExpression(printer, expression.Lhs)
		}
		if expression.RhsExpression != nil {
			PrintExpression(printer, expression.RhsExpression)
		}
		printer.indentation--
	}
}
