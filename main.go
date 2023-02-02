package main

import (
	"fmt"
	"os"

	"github.com/milansav/Castle/cli"
	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/parser"
)

func main() {

	settings := cli.ParseArguments()

	for _, file := range settings.Files {
		contents, err := os.ReadFile(file)

		if err != nil {
			_ = fmt.Errorf("Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Loaded file contents: \n%s\n", string(contents))

		mainLexer := lexer.Create(string(contents))
		lexer.Start(&mainLexer)

		for _, element := range mainLexer.Lexemes {
			fmt.Printf("Label: %s, Type: %s\n", element.Label, lexer.LexemeTypeLabels[element.Type])
		}

		mainParser := parser.Create(mainLexer)
		program := parser.StartNew(&mainParser)

		fmt.Println(len(program.Statements))
		fmt.Println(program.Statements[0].SType)
		fmt.Println(program.Statements[0].Function.Name)
		fmt.Println(program.Statements[0].Function.Statement.Statements[0].Declaration.Value)

		//expressionParser := parser.Create(mainLexer)
		//result := parser.StartExpressionParser(&expressionParser)

		//fmt.Println(string(contents))
		//for _, value := range result {
		//	parser.PrintTree(value, 0)
		//}
	}
}
