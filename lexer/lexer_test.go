package lexer

import (
	"testing"
)

func TestLexerCreate(t *testing.T) {
	input := "1 + 2 * 3 / (145 - 23)"
	lexer := Create(input)

	if len(lexer.Lexemes) != 0 {
		t.Error("lexer.Create lexemes length is not equal to 0")
	}

	if lexer.source != input {
		t.Error("lexer.Create content is not the same as source")
	}

	if lexer.currentStep != 0 {
		t.Error("lexer.Create currentStep is not 0")
	}

}

func TestLexerExpression(t *testing.T) {
	input := "1 + 2 * 3 / (145 - 23)"
	expectedTypes := [11]LexemeType{LT_NUMBER, LT_PLUS, LT_NUMBER, LT_MULTIPLY, LT_NUMBER, LT_DIVIDE, LT_LPAREN, LT_NUMBER, LT_MINUS, LT_NUMBER, LT_RPAREN}
	lexer := Create(input)

	Start(&lexer)

	if len(lexer.Lexemes) != 11 {
		t.Errorf("lexer.Start Lexemes size is incorrect. Expected 5 got %d", len(lexer.Lexemes))
	}

	for index, element := range lexer.Lexemes {
		if element.Type != expectedTypes[index] {
			t.Errorf("lexer.Start Lexeme at index %d incorrect type", index)
			t.FailNow()
		}
	}
}
