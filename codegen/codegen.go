package codegen

import (
	"fmt"

	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/parser"
)

type Codegen struct {
	Program   *parser.AST_Program
	OutBuffer string
}

func Create(program *parser.AST_Program) Codegen {
	return Codegen{
		Program: program,
	}
}

func (c *Codegen) Out(text string) {
	c.OutBuffer += text
}

func (codegen *Codegen) Start() {

	// c.Out("int main (int argc, char **argv) {\n")
	// c.Out("return 0;\n")
	// c.Out("}\n")

	codegen.Out("#define bool int\n")
	codegen.Out("#define true 1\n")
	codegen.Out("#define false 0\n")
	codegen.Out("#define NULL 0\n")

	for _, programStatement := range codegen.Program.Statements {

		codegen.PrintStatement(programStatement)
	}
}

func findFirstLiteral(node interface{}) *parser.AST_Literal {
	switch n := node.(type) {
	case *parser.AST_Literal:
		return n
	case *parser.AST_FunctionCall:
		for _, param := range n.Params {
			if result := findFirstLiteral(param); result != nil {
				return result
			}
		}
	case *parser.AST_Expression:
		if n.Lhs != nil {
			if result := findFirstLiteral(n.Lhs); result != nil {
				return result
			}
		}
		if n.Rhs != nil {
			if result := findFirstLiteral(n.Rhs); result != nil {
				return result
			}
		}
		if n.Value != nil {
			return n.Value
		}
		if n.FunctionCall != nil {
			if result := findFirstLiteral(n.FunctionCall); result != nil {
				return result
			}
		}
	}
	return nil
}

func stringifyOperator(lt lexer.LexemeType) string {
	switch lt {
	case lexer.LT_PLUS:
		return "+"
	case lexer.LT_MINUS:
		return "-"
	case lexer.LT_MULTIPLY:
		return "*"
	case lexer.LT_DIVIDE:
		return "/"
	case lexer.LT_MODULO:
		return "%"
	case lexer.LT_GEQ:
		return ">="
	case lexer.LT_LEQ:
		return "<="
	case lexer.LT_NEQ:
		return "!="
	case lexer.LT_EQ:
		return "=="
	case lexer.LT_AND:
		return "&&"
	case lexer.LT_OR:
		return "||"
	case lexer.LT_LCHEVRON:
		return "<"
	case lexer.LT_RCHEVRON:
		return ">"
	default:
		panic(fmt.Errorf("error: unknown operator %s", lexer.LexemeTypeLabels[lt]))
	}
}

func (codegen *Codegen) PrintStatement(statement *parser.AST_Statement) {
	switch statement.SType {
	case parser.ST_STATEMENT_ARRAY:
		for _, statement := range statement.Statements {
			codegen.PrintStatement(statement)
		}
	case parser.ST_STATEMENT:
	case parser.ST_EXPRESSION:
		codegen.PrintExpression(statement.Expression)
		codegen.Out(";\n")
	case parser.ST_FUNCTION:
		codegen.Out("void ")
		codegen.Out(statement.Function.Name)
		codegen.Out("() {\n")
		if statement.Function.Statement.Statement != nil {
			codegen.Out("return ")
			codegen.PrintStatement(statement.Function.Statement.Statement)
			codegen.Out(";\n")
		} else {
			codegen.PrintStatement(statement.Function.Statement)
		}
		codegen.Out("}\n")
	case parser.ST_DECLARATION:
		e := findFirstLiteral(statement.Declaration.Value)

		if e == nil {
			panic(fmt.Errorf("could not find literal in declaration %s", statement.Declaration.Name))
		}

		switch e.Type {
		case parser.TYPE_NUMBER:
			codegen.Out("int ")
		case parser.TYPE_FLOAT:
			codegen.Out("float ")
		case parser.TYPE_STRING:
			codegen.Out("char* ")
		case parser.TYPE_BOOL:
			codegen.Out("bool ")
		case parser.TYPE_UNDEFINED:
			codegen.Out("void* ")
		}

		codegen.Out(statement.Declaration.Name)

		codegen.Out(" = ")

		codegen.PrintExpression(statement.Declaration.Value)

		codegen.Out(";\n")
	case parser.ST_STRUCT:
	case parser.ST_IF:
		codegen.Out("if (")
		codegen.PrintExpression(statement.If.Condition)
		codegen.Out(") {\n")
		for _, statement := range statement.If.Statements {
			codegen.PrintStatement(statement)
		}
		codegen.Out("}\n")
	}
}

func (codegen *Codegen) PrintExpression(expression *parser.AST_Expression) {

	switch expression.EType {
	case parser.ET_BINARY:
		codegen.PrintExpression(expression.Lhs)
		codegen.Out(" ")
		codegen.Out(stringifyOperator(expression.Operator))
		codegen.Out(" ")
		codegen.PrintExpression(expression.Rhs)
	case parser.ET_VALUE:
		codegen.PrintLiteral(expression.Value)
	case parser.ET_FUNCTION_CALL:
		codegen.PrintFunctionCall(expression.FunctionCall)
	case parser.ET_IDENTIFIER:
		codegen.Out(expression.Identifier)
	}

}

func (codegen *Codegen) PrintLiteral(literal *parser.AST_Literal) {
	switch literal.Type {
	case parser.TYPE_NUMBER:
		codegen.Out(literal.Value)
	case parser.TYPE_FLOAT:
		codegen.Out(literal.Value)
	case parser.TYPE_STRING:
		codegen.Out(literal.Value)
	case parser.TYPE_BOOL:
		codegen.Out(literal.Value)
	case parser.TYPE_UNDEFINED:
		codegen.Out("NULL")
	}
}

func (codegen *Codegen) PrintFunctionCall(functionCall *parser.AST_FunctionCall) {
	codegen.Out(functionCall.Name)
	codegen.Out("(")
	for index, param := range functionCall.Params {
		codegen.PrintExpression(param)

		if index < len(functionCall.Params)-1 {
			codegen.Out(",")
		}
	}
	codegen.Out(")")
}
