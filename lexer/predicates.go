package lexer

import "strings"

func isNumeric(char rune) bool {
	return char >= '0' && char <= '9'
}

func isWhitespace(char rune) bool {
	str := string(char)
	return len(str) != len(strings.TrimSpace(str))
}
