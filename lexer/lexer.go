package lexer

import (
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	Lexemes     []Lexeme
	source      string
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

	LT_EQUALS

	//Logical
	LT_COMPARE
	LT_AND
	LT_OR
	LT_NAND
	LT_NOR
	LT_XOR
	LT_XAND
	LT_XNOR
	LT_XNAND

	//Parentheses
	LT_LPAREN
	LT_RPAREN
	LT_LCURLY
	LT_RCURLY
	LT_LCHEVRON
	LT_RCHEVRON
	LT_LBRACKET
	LT_RBRACKET

	//Binary operations
	LT_BANG

	LT_IDENTIFIER

	// Literals
	LT_LITERAL_NUMBER
	LT_LITERAL_FLOAT
	LT_LITERAL_STRING
	LT_LITERAL_BOOL

	//Keywords
	LT_CONST
	LT_VAL
	LT_IF
	LT_ELSE
	LT_ELSEIF
	LT_INTERFACE
	LT_STRUCT
	LT_OF

	//Misc operators
	LT_LAMBDA
	LT_SEMICOLON
	LT_COLON
	LT_MACRO

	LT_COMMA
	LT_PERIOD

	LT_NONE
	LT_UNKNOWN
	LT_END
)

var LexemeTypeLabels = map[LexemeType]string{
	//Binary operators
	LT_PLUS:     "LT_PLUS",
	LT_MINUS:    "LT_MINUS",
	LT_MULTIPLY: "LT_MULTIPLY",
	LT_DIVIDE:   "LT_DIVIDE",
	LT_MODULO:   "LT_MODULO",
	LT_POWER:    "LT_POWER",

	LT_EQUALS: "LT_EQUALS",

	//Logical
	LT_COMPARE: "LT_COMPARE",
	LT_AND:     "LT_AND",
	LT_OR:      "LT_OR",
	LT_NAND:    "LT_NAND",
	LT_NOR:     "LT_NOR",
	LT_XOR:     "LT_XOR",
	LT_XAND:    "LT_XAND",
	LT_XNOR:    "LT_XNOR",
	LT_XNAND:   "LT_XNAND",

	//Parentheses
	LT_LPAREN:   "LT_LPAREN",
	LT_RPAREN:   "LT_RPAREN",
	LT_LCURLY:   "LT_LCURLY",
	LT_RCURLY:   "LT_RCURLY",
	LT_LCHEVRON: "LT_LCHEVRON",
	LT_RCHEVRON: "LT_RCHEVRON",
	LT_LBRACKET: "LT_LBRACKET",
	LT_RBRACKET: "LT_RBRACKET",

	LT_BANG: "LT_BANG",

	LT_IDENTIFIER: "LT_IDENTIFIER",

	// Literals
	LT_LITERAL_NUMBER: "LT_LITERAL_NUMBER",
	LT_LITERAL_FLOAT:  "LT_LITERAL_FLOAT",
	LT_LITERAL_STRING: "LT_LITERAL_STRING",
	LT_LITERAL_BOOL:   "LT_LITERAL_BOOL",

	//Keywords
	LT_CONST:     "LT_CONST",
	LT_VAL:       "LT_VAL",
	LT_IF:        "LT_IF",
	LT_ELSE:      "LT_ELSE",
	LT_ELSEIF:    "LT_ELSEIF",
	LT_INTERFACE: "LT_INTERFACE",
	LT_STRUCT:    "LT_STRUCT",
	LT_OF:        "LT_OF",

	//Misc operators
	LT_LAMBDA:    "LT_LAMBDA",
	LT_SEMICOLON: "LT_SEMICOLON",
	LT_COLON:     "LT_COLON",
	LT_MACRO:     "LT_MACRO",

	LT_COMMA:  "LT_COMMA",
	LT_PERIOD: "LT_PERIOD",

	LT_NONE:    "LT_NONE",
	LT_UNKNOWN: "LT_UNKNOWN",
	LT_END:     "LT_END",
}

var keywords = map[string]LexemeType{
	"const":     LT_CONST,
	"val":       LT_VAL,
	"if":        LT_IF,
	"else":      LT_ELSE,
	"elseif":    LT_ELSEIF,
	"interface": LT_INTERFACE,
	"struct":    LT_STRUCT,
	"of":        LT_OF,

	"and":   LT_AND,
	"or":    LT_OR,
	"nand":  LT_NAND,
	"nor":   LT_NOR,
	"xor":   LT_XOR,
	"xand":  LT_XAND,
	"xnor":  LT_XNOR,
	"xnand": LT_XNAND,
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

		c := currentRune(lexer)

		if c == '/' {
			switch nextRune(lexer) {
			case '/':
				lineComment(lexer)
				//lexer.Lexemes = append(lexer.Lexemes, lexeme)
				continue
				//case '*':
				//	lexeme := blockComment(lexer)
				//	lexer.Lexemes = append(lexer.Lexemes, lexeme)
				//	continue
			}
		}

		if unicode.IsSpace(c) {
			whitespace(lexer)
			continue
		} else if unicode.IsLetter(c) {
			lexeme := identifier(lexer)
			lexer.Lexemes = append(lexer.Lexemes, lexeme)
			continue
		} else if unicode.IsDigit(c) {
			lexeme := number(lexer)
			lexer.Lexemes = append(lexer.Lexemes, lexeme)
			continue
		} else {
			lexeme := other(lexer)
			lexer.Lexemes = append(lexer.Lexemes, lexeme)
		}

	}

	lexer.Lexemes = append(lexer.Lexemes, Lexeme{Type: LT_END})
}

func lineComment(lexer *Lexer) Lexeme {
	for {
		if currentRune(lexer) != '\n' && canStep(lexer) {
			step(lexer)
			continue
		}
		break
	}

	return Lexeme{}
}

func blockComment(lexer *Lexer) Lexeme {

	return Lexeme{}
}

func identifier(lexer *Lexer) Lexeme {

	start := lexer.currentStep

	for unicode.IsLetter(currentRune(lexer)) || unicode.IsDigit(currentRune(lexer)) {
		step(lexer)
	}

	end := lexer.currentStep

	//length := end - start

	//fmt.Printf("Lexeme processed: Label: \"%s\", Length: %d\n", lexer.source[start:end], length)

	lexeme := Lexeme{Label: lexer.source[start:end], Type: LT_IDENTIFIER}

	lexeme = getKeyword(lexeme)

	return lexeme
}

func getKeyword(lexeme Lexeme) Lexeme {
	if _, ok := keywords[lexeme.Label]; ok {
		lexeme.Type = keywords[lexeme.Label]
	}

	return lexeme
}

func number(lexer *Lexer) Lexeme {

	isFloat := false

	start := lexer.currentStep

	numberType := LT_LITERAL_NUMBER

	c := currentRune(lexer)

	for unicode.IsDigit(c) || c == ',' || (c == '.' && !isFloat) {
		if c == '.' {
			isFloat = true
			numberType = LT_LITERAL_FLOAT
		}
		step(lexer)
		c = currentRune(lexer)
	}

	end := lexer.currentStep

	lexeme := Lexeme{Label: lexer.source[start:end], Type: numberType}

	return lexeme
}

func other(lexer *Lexer) Lexeme {
	target := currentRune(lexer)

	lexeme := Lexeme{Type: LT_NONE}

	start := lexer.currentStep

	switch target {
	case '=':
		lexeme.Type = LT_EQUALS

		if nextRune(lexer) == '>' {
			lexeme.Type = LT_LAMBDA
			step(lexer)
		}
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
	case '<':
		lexeme.Type = LT_LCHEVRON
	case '>':
		lexeme.Type = LT_RCHEVRON
	case '[':
		lexeme.Type = LT_LBRACKET
	case ']':
		lexeme.Type = LT_RBRACKET
	case ',':
		lexeme.Type = LT_COMMA
	case '.':
		lexeme.Type = LT_PERIOD
	case ';':
		lexeme.Type = LT_SEMICOLON
	case ':':
		lexeme.Type = LT_COLON
	case '$':
		if nextRune(lexer) == '$' {
			lexeme.Type = LT_MACRO
			step(lexer)
			break
		}

		lexeme.Type = LT_UNKNOWN
	case '"':
		lexeme.Type = LT_LITERAL_STRING

		start := lexer.currentStep
		step(lexer)
		for {
			if nextRune(lexer) == '"' {
				step(lexer)
				break
			}
			step(lexer)
		}

		end := lexer.currentStep

		lexeme.Label = lexer.source[start:end]
	}

	step(lexer)

	end := lexer.currentStep
	lexeme.Label = lexer.source[start:end]

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
