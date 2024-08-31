package main

import (
	"fmt"
	"os"

	"github.com/milansav/Castle/astprinter"
	"github.com/milansav/Castle/cli"
	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/parser"
)

func main() {

	settings := cli.ParseArguments()

	for _, file := range settings.Files {
		contents, err := os.ReadFile(file)

		if err != nil {
			_ = fmt.Errorf("error: %s", err)
			os.Exit(1)
		}

		fmt.Printf("----------SOURCE----------\n%s\n", string(contents))

		mainLexer := lexer.Create(string(contents))
		mainLexer.Start()

		// fmt.Println("-----LEXICAL ANALYSIS-----")

		// for _, element := range mainLexer.Lexemes {
		// 	fmt.Printf("Label: %s, Type: %s\n", element.Label, lexer.LexemeTypeLabels[element.Type])
		// }

		fmt.Println("------SYNTAX ANALYSIS-----")

		mainParser := parser.Create(mainLexer)
		program := mainParser.Start()

		fmt.Println("---ABSTRACT SYNTAX TREE---")

		astprinter.PrintAST(program)
	}
}
