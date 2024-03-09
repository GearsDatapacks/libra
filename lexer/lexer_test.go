package lexer_test

import (
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lexer/token"
	utils "github.com/gearsdatapacks/libra/test_utils"
	"github.com/gearsdatapacks/libra/text"
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
		{"..", token.DOUBLE_DOT},
		{"!", token.BANG},
		{"|", token.PIPE},
		{"->", token.ARROW},
		{"&", token.AMPERSAND},
		{"~", token.TILDE},
		{"\n", token.NEWLINE},
		{";", token.SEMICOLON},
	}

	for _, tok := range tokens {
		lexer := lexer.New(text.NewFile("test.lb", tok.src))
		tokens := lexer.Tokenise()

		utils.AssertEq(t, len(tokens), 2)
		utils.AssertEq(t, tokens[0].Kind, tok.kind)
		utils.AssertEq(t, tokens[1].Kind, token.EOF)
	}
}

func TestVariableTokens(t *testing.T) {
	type tokenData struct {
		src  string
		kind token.Kind
		text string
	}
	tok := func(src string, kind token.Kind, text ...string) tokenData {
		return tokenData{
			src:  src,
			kind: kind,
			text: append(text, "")[0],
		}
	}
	tokens := []tokenData{
		tok("17", token.INTEGER),
		tok("42", token.INTEGER),
		tok("123_456_789", token.INTEGER, "123456789"),
		tok("19.3", token.FLOAT),
		tok("3.141_592_65", token.FLOAT, "3.14159265"),
		tok("foo_bar", token.IDENTIFIER),
		tok("HiThere123", token.IDENTIFIER),
		tok(`"Hi :)"`, token.STRING, "Hi :)"),
		tok(`"\"How are you?\""`, token.STRING, `"How are you?"`),
		tok(`"Hello\nworld"`, token.STRING, "Hello\nworld"),
	}

	for _, data := range tokens {
		lexer := lexer.New(text.NewFile("test.lb", data.src))
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

func TestLexerDiagnostics(t *testing.T) {
	data := []struct {
		src  string
		msg  string
		span text.Span
	}{
		{"foo@bar", "Invalid character: '@'", text.NewSpan(0, 3, 4)},
		{`"Hello`, "Unterminated string", text.NewSpan(0, 0, 6)},
		{`"He\llo"`, "Invalid escape sequence: '\\l'", text.NewSpan(0, 3, 5)},
		{"123_456_", "Numbers cannot end with numeric separators", text.NewSpan(0, 7, 8)},
		{"1_.2", "Numbers cannot end with numeric separators", text.NewSpan(0, 1, 2)},
		{"3.14_", "Numbers cannot end with numeric separators", text.NewSpan(0, 4, 5)},
	}

	for _, data := range data {
		lexer := lexer.New(text.NewFile("test.lb", data.src))
		lexer.Tokenise()

		utils.AssertEq(t, len(lexer.Diagnostics), 1)
		utils.AssertEq(t, lexer.Diagnostics[0].Message, data.msg)
		utils.AssertEq(t, lexer.Diagnostics[0].Location.Span, data.span)
	}
}
