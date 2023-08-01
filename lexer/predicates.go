package lexer

import "strings"

func isNumeric(char rune) bool {
	return char >= '0' && char <= '9'
}

func isWhitespace(char rune) bool {
	str := string(char)
	return len(str) != len(strings.TrimSpace(str))
}

func isAlphabetic(char rune) bool {
	isLower := char >= 'a' && char <= 'z'
	isUpper := char >= 'A' && char <= 'Z'
	return isLower || isUpper || char == '_'
}

func isAlphanumeric(char rune) bool {
	return isNumeric(char) || isAlphabetic(char)
}
