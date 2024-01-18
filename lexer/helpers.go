package lexer

func isIdentifierStart(c byte) bool {
  isCapital := c >= 'A' && c <= 'Z'
  isLower := c >= 'a' && c <= 'z'
  isOther := c == '_'
  return isCapital || isLower || isOther
}

func isIdentifierMiddle(c byte) bool {
  return isIdentifierStart(c) || isNumber(c)
}

func isNumber(c byte) bool {
  return c >= '0' && c <= '9'
}

func isWhitespace(c byte) bool {
  return c == ' ' || c == '\t'
}

func escape(c byte) byte {
	switch c {
  // TODO: \x and \u
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
  // TODO: Error for invalid sequence
	default:
		return c
	}
}

