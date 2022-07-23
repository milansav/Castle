package main

import (
	"fmt"
	"github.com/milansav/Castle/lexer"
	"github.com/milansav/Castle/parser"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Hello, World")

	contents, err := ioutil.ReadFile("test.cst")

	if err != nil {
		_ = fmt.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded file contents: \n%s\n", string(contents))

	mainLexer := lexer.Create(string(contents))
	lexer.Start(&mainLexer)

	mainParser := parser.Create(mainLexer)
	result := parser.Start(&mainParser)

	fmt.Println(string(contents))
	for _, value := range(result) {
		parser.PrintTree(value, 0)
	}
}
