package lexer

func isIdentifierStart(c byte) bool {
	isCapital := c >= 'A' && c <= 'Z'
	isLower := c >= 'a' && c <= 'z'
	isOther := c == '_'
	return isCapital || isLower || isOther
}

func isIdentifierMiddle(c byte) bool {
	return isIdentifierStart(c) || isNumber(c, 10)
}

func isNumber(c byte, radix int) bool {
	switch radix {
	case 2:
		return c == '0' || c == '1'
	case 8:
		return c >= '0' && c <= '7'
	case 10:
		return c >= '0' && c <= '9'
	case 16:
		return (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'f') ||
			(c >= 'A' && c <= 'F')
  default: panic("Unreachable: No other bases are implemented")
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

func charToRadix(c byte) int {
	switch c {
	case 'b':
		return 2
	case 'o':
		return 8
	case 'x':
		return 16
	default:
		return 0
	}
}
