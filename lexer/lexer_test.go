package lexer_test

import (
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lexer/token"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestSingleTokens(t *testing.T) {
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
	}

	for _, p := range tokens {
		lexer := lexer.New(p.src)
		tokens := lexer.Tokenise()

		utils.AssertEq(t, len(tokens), 2)
		utils.AssertEq(t, tokens[0].Kind, p.kind)
		utils.AssertEq(t, tokens[1].Kind, token.EOF)
	}
}
