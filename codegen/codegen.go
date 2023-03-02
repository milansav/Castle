package codegen

import (
	"fmt"

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

func (c Codegen) Out(text string) {
	c.OutBuffer += text
}

func (codegen *Codegen) Start() {
	for _, programStatement := range codegen.Program.Statements {
		fmt.Println(programStatement.SType)

		switch programStatement.SType {
		case parser.ST_STATEMENT_ARRAY:
		case parser.ST_STATEMENT:
		case parser.ST_EXPRESSION:
		case parser.ST_FUNCTION:
		case parser.ST_DECLARATION:
		case parser.ST_STRUCT:
		case parser.ST_IF:
		}
	}
}
