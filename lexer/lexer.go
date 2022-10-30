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
	LT_PLUS LexemeType = iota
	LT_MINUS
	LT_MULTIPLY
	LT_DIVIDE
	LT_MODULO
	LT_LPAREN
	LT_RPAREN
	LT_BANG
	LT_KEYWORD
	LT_LITERAL
	LT_NUMBER
	LT_NONE
	LT_END
)

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

	for unicode.IsDigit(currentRune(lexer)) {
		step(lexer)
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

	if target == '+' {
		lexeme.Type = LT_PLUS
		//fmt.Println("Is plus")
	} else if target == '-' {
		lexeme.Type = LT_MINUS
		//fmt.Println("Is minus")
	} else if target == '*' {
		lexeme.Type = LT_MULTIPLY
		//fmt.Println("Is multiply")
	} else if target == '/' {
		lexeme.Type = LT_DIVIDE
		//fmt.Println("Is divide")
	} else if target == '(' {
		lexeme.Type = LT_LPAREN
		//fmt.Println("Is left parentheses")
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
