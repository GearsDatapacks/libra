package lexer

import "strings"

func isNumeric(char rune) bool {
	return char <= '9' && char >= '0'
}

func isWhitespace(char rune) bool {
	str := string(char)
	return len(str) != len(strings.TrimSpace(str))
}
