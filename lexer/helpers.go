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

