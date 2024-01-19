package lexer_test

import (
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lexer/token"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestFixedTokens(t *testing.T) {
	var tokens = []struct {
		src  string
		kind token.Kind
	}{
		{"(", token.LEFT_PAREN},
		{")", token.RIGHT_PAREN},
		{"{", token.LEFT_BRACE},
		{"}", token.RIGHT_BRACE},
		{"[", token.LEFT_SQUARE},
		{"]", token.RIGHT_SQUARE},
		{",", token.COMMA},
		{".", token.DOT},
    {":", token.COLON},
		{"?", token.QUESTION},
		{"=", token.EQUALS},
		{"+=", token.PLUS_EQUALS},
		{"-=", token.MINUS_EQUALS},
		{"*=", token.STAR_EQUALS},
		{"/=", token.SLASH_EQUALS},
		{"%=", token.PERCENT_EQUALS},
		{"&&", token.DOUBLE_AMPERSAND},
		{"||", token.DOUBLE_PIPE},
		{"<", token.LEFT_ANGLE},
		{">", token.RIGHT_ANGLE},
		{"<=", token.LEFT_ANGLE_EQUALS},
		{">=", token.RIGHT_ANGLE_EQUALS},
		{"==", token.DOUBLE_EQUALS},
		{"!=", token.BANG_EQUALS},
		{"<<", token.DOUBLE_LEFT_ANGLE},
		{">>", token.DOUBLE_RIGHT_ANGLE},
		{"+", token.PLUS},
		{"-", token.MINUS},
		{"*", token.STAR},
		{"/", token.SLASH},
		{"%", token.PERCENT},
		{"**", token.DOUBLE_STAR},
		{"++", token.DOUBLE_PLUS},
		{"--", token.DOUBLE_MINUS},
		{"!", token.BANG},
		{"|", token.PIPE},
    {"->", token.ARROW},
		{"&", token.AMPERSAND},
		{"\n", token.NEWLINE},
		{"\r", token.NEWLINE},
	}

	for _, tok := range tokens {
		lexer := lexer.New(tok.src)
		tokens := lexer.Tokenise()

		utils.AssertEq(t, len(tokens), 2)
		utils.AssertEq(t, tokens[0].Kind, tok.kind)
		utils.AssertEq(t, tokens[1].Kind, token.EOF)
	}
}

func TestVariableTokens(t *testing.T) {
  type tokenData struct{
    src string
    kind token.Kind
    text string
  }
  tok := func (src string, kind token.Kind, text ...string) tokenData {
    return tokenData{
    src: src,
    kind: kind,
    text: append(text, "")[0],
    }
  }
  tokens := []tokenData{
    tok("17", token.INTEGER),
    tok("42", token.INTEGER),
    tok("19.3", token.FLOAT),
    tok("foo_bar", token.IDENTIFIER),
    tok("HiThere123", token.IDENTIFIER),
    tok(`"Hi :)"`, token.STRING, "Hi :)"),
    tok(`"\"How are you?\""`, token.STRING, `"How are you?"`),
    tok(`"Hello\nworld"`, token.STRING, "Hello\nworld"),
  }

  for _, data := range tokens {
    lexer := lexer.New(data.src)
    tokens := lexer.Tokenise()

    utils.AssertEq(t, len(tokens), 2)
    utils.AssertEq(t, tokens[0].Kind, data.kind)
    text := data.src
    if data.text != "" {
      text = data.text
    }
    utils.AssertEq(t, tokens[0].Value, text)
    utils.AssertEq(t, tokens[1].Kind, token.EOF)
  }
}

