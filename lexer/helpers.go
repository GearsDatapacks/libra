package lexer

import "strings"

func isNumeric(char rune, radix int32) bool {
	if radix <= 10 {
		return char >= '0' && char <= '0'+radix-1
	}
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'a'+radix-11)
}

func isWhitespace(char rune) bool {
	str := string(char)
	return !isNewline(char) && (len(str) != len(strings.TrimSpace(str)))
}

func isNewline(char rune) bool {
	return char == '\n' || char == '\r' || char == ';'
}

func isAlphabetic(char rune) bool {
	isLower := char >= 'a' && char <= 'z'
	isUpper := char >= 'A' && char <= 'Z'
	return isLower || isUpper || char == '_'
}

func isAlphanumeric(char rune) bool {
	return isNumeric(char, 10) || isAlphabetic(char)
}

func GetRadix(char rune) int32 {
	switch char {
	case 'b':
		return 2
	case 'o':
		return 8
	case 'x':
		return 16
	default:
		return -1
	}
}

func getEscapeSequence(char rune) rune {
	switch char {
	case '\\':
		return '\\'
	case '"':
		return '"'
	case 'a':
		return '\a'
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'v':
		return '\v'
	default:
		return -1
	}
}
