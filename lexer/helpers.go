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
  return c == ' ' || c == '\t' || c == '\r'
}

func escape(c byte) (char byte, ok bool) {
	switch c {
  // TODO: \x and \u
	case '\\':
		char = '\\'
	case '"':
		char = '"'
	case 'a':
		char = '\a'
	case 'b':
		char = '\b'
	case 'f':
		char = '\f'
	case 'n':
		char = '\n'
	case 'r':
		char = '\r'
	case 't':
		char = '\t'
	case 'v':
		char = '\v'
	default:
    return c, false
	}

  ok = true
  return
}

