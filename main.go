package main

import (
	"fmt"
	"os"

	"github.com/milansav/Castle/cli"
	"github.com/milansav/Castle/codegen"
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
		program := parser.Start(&mainParser)

		mainCodegen := codegen.Create()
		codegen.Start(mainCodegen)

		fmt.Println(len(program.Statements))

		//fmt.Println(string(contents))
		//for _, value := range result {
		//	parser.PrintTree(value, 0)
		//}
	}
}
