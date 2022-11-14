package lexer

import (
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	Lexemes     []Lexeme
	source      string
	currentRune rune
	nextRune    rune
	currentStep int
	row         int
	column      int
}

type Lexeme struct {
	Label  string
	Type   LexemeType
	Row    int
	Column int
}

type LexemeType int

const (
	//Arithmetic operators
	LT_PLUS LexemeType = iota
	LT_MINUS
	LT_MULTIPLY
	LT_DIVIDE
	LT_MODULO
	LT_POWER

	//Parentheses
	LT_LPAREN
	LT_RPAREN
	LT_LCURLY
	LT_RCURLY

	//Binary operations
	LT_BANG

	LT_LITERAL
	LT_NUMBER

	//Keywords
	LT_CONST
	LT_IF
	LT_ELSE
	LT_ELSEIF

	//Misc operators
	LT_LAMBDA

	LT_NONE
	LT_END
)

var LexemeTypeLabels = map[LexemeType]string{
	//Arithmetic operators
	LT_PLUS:     "LT_PLUS",
	LT_MINUS:    "LT_MINUS",
	LT_MULTIPLY: "LT_MULTIPLY",
	LT_DIVIDE:   "LT_DIVIDE",
	LT_MODULO:   "LT_MODULO",
	LT_POWER:    "LT_POWER",

	//Parentheses
	LT_LPAREN: "LT_LPAREN",
	LT_RPAREN: "LT_RPAREN",
	LT_LCURLY: "LT_LCURLY",
	LT_RCURLY: "LT_RCURLY",

	//Binary operations
	LT_BANG: "LT_BANG",

	LT_LITERAL: "LT_LITERAL",
	LT_NUMBER:  "LT_NUMBER",

	//Keywords
	LT_CONST:  "LT_CONST",
	LT_IF:     "LT_IF",
	LT_ELSE:   "LT_ELSE",
	LT_ELSEIF: "LT_ELSEIF",

	//Misc operators
	LT_LAMBDA: "LT_ELSEIF",

	LT_NONE: "LT_NONE",
	LT_END:  "LT_END",
}

func Create(source string) Lexer {

	return Lexer{
		Lexemes:     make([]Lexeme, 0),
		source:      source,
		currentStep: 0,
	}
}

func Start(lexer *Lexer) {
	for canStep(lexer) {

		if unicode.IsSpace(currentRune(lexer)) {
			whitespace(lexer)
			continue
		} else if unicode.IsLetter(currentRune(lexer)) {
			lexeme := literal(lexer)
			lexer.Lexemes = append(lexer.Lexemes, lexeme)
			continue
		} else if unicode.IsDigit(currentRune(lexer)) {
			lexeme := number(lexer)
			lexer.Lexemes = append(lexer.Lexemes, lexeme)
			continue
		} else {
			lexeme := other(lexer)
			lexer.Lexemes = append(lexer.Lexemes, lexeme)
		}

	}
}

func literal(lexer *Lexer) Lexeme {

	start := lexer.currentStep

	for unicode.IsLetter(currentRune(lexer)) || unicode.IsDigit(currentRune(lexer)) {
		step(lexer)
	}

	end := lexer.currentStep

	//length := end - start

	//fmt.Printf("Lexeme processed: Label: \"%s\", Length: %d\n", lexer.source[start:end], length)

	lexeme := Lexeme{Label: lexer.source[start:end], Type: LT_LITERAL}

	return lexeme
}

func number(lexer *Lexer) Lexeme {

	start := lexer.currentStep

	c := currentRune(lexer)

	for unicode.IsDigit(c) || c == ',' {
		step(lexer)
		c = currentRune(lexer)
	}

	end := lexer.currentStep

	//length := end - start

	//fmt.Printf("Lexeme processed: Label: \"%s\", Length: %d\n", lexer.source[start:end], length)

	lexeme := Lexeme{Label: lexer.source[start:end], Type: LT_NUMBER}

	return lexeme
}

func other(lexer *Lexer) Lexeme {
	target := currentRune(lexer)

	lexeme := Lexeme{Label: string(target), Type: LT_NONE}

	switch target {
	case '+':
		lexeme.Type = LT_PLUS
	case '-':
		lexeme.Type = LT_MINUS
	case '*':
		lexeme.Type = LT_MULTIPLY
	case '/':
		lexeme.Type = LT_DIVIDE
	case '%':
		lexeme.Type = LT_MODULO
	case '^':
		lexeme.Type = LT_POWER
	case '(':
		lexeme.Type = LT_LPAREN
	case ')':
		lexeme.Type = LT_RPAREN
	case '{':
		lexeme.Type = LT_LCURLY
	case '}':
		lexeme.Type = LT_RCURLY
	}

	step(lexer)

	return lexeme
}

func whitespace(lexer *Lexer) {
	for unicode.IsSpace(currentRune(lexer)) {
		step(lexer)
	}
}

func canStep(lexer *Lexer) bool {
	return lexer.currentStep < len(lexer.source)
}

func step(lexer *Lexer) {
	_, size := utf8.DecodeRuneInString(lexer.source[lexer.currentStep:])

	lexer.currentStep += size
}

func currentRune(lexer *Lexer) rune {
	currentRune, _ := utf8.DecodeRuneInString(lexer.source[lexer.currentStep:])

	return currentRune
}

func nextRune(lexer *Lexer) rune {
	_, size := utf8.DecodeRuneInString(lexer.source[lexer.currentStep:])
	nextRune, _ := utf8.DecodeRuneInString(lexer.source[lexer.currentStep+size:])

	return nextRune
}
