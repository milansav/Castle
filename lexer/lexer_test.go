package lexer

import "testing"

func TestCreate(t *testing.T) {
	input := "1 + 2 + 3"
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

func TestStart(t *testing.T) {
	input := "1 + 2 + 3"
	lexer := Create(input)

	Start(&lexer)

	if len(lexer.Lexemes) != 5 {
		t.Errorf("lexer.Start Lexemes size is incorrect. Expected 5 got %d", len(lexer.Lexemes))
	}
}
