package lexer

import (
	"fmt"
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
	Label string
	Type  LexemeType
}

type LexemeType int

const (
	ADD LexemeType = iota
	SUBTRACT
	MULTIPLY
	DIVIDE
	MODULO
	LPAREN
	RPAREN
	KEYWORD
	LITERAL
	NUMBER
	END
)
func Create(source string) Lexer {

	currentRune, size := utf8.DecodeRuneInString(source)
	nextRune, _ := utf8.DecodeRuneInString(source[size:])

	return Lexer{
		Lexemes:     make([]Lexeme, 0),
		source:      source,
		currentStep: 0,
		currentRune: currentRune,
		nextRune:    nextRune,
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
			//TODO: Process other characters
		}

	}
}

func literal(lexer *Lexer) Lexeme {

	start := lexer.currentStep

	for(unicode.IsLetter(currentRune(lexer)) || unicode.IsDigit(currentRune(lexer))) {
		step(lexer)
	}

	end := lexer.currentStep

	length := end - start

	fmt.Printf("Lexeme processed: Label: \"%s\", Length: %d\n", lexer.source[start:end], length)

	lexeme := Lexeme{Label: lexer.source[start:end], Type: LITERAL}

	return lexeme
}

func number(lexer *Lexer) Lexeme {

	start := lexer.currentStep

	for(unicode.IsDigit(currentRune(lexer))) {
		step(lexer)
	}

	end := lexer.currentStep

	length := end - start

	fmt.Printf("Lexeme processed: Label: \"%s\", Length: %d\n", lexer.source[start:end], length)

	lexeme := Lexeme{Label: lexer.source[start:end], Type: NUMBER}

	return lexeme
}

func whitespace(lexer *Lexer) {
	step(lexer)
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
