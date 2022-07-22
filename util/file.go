package util

import "github.com/milansav/Castle/lexer"

type File struct {
	Content string
	Lexer   lexer.Lexer
}
