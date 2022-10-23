package util

import "unicode/utf8"

func GetRune(input string, position int) rune {
	currentRune, _ := utf8.DecodeRuneInString(input[position:])

	return currentRune
}
