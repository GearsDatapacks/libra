package parser_test

import (
	"testing"

	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func getProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input, "test.lb")
	tokens := l.Tokenise()

	p := parser.New(tokens, l.Diagnostics)
	program := p.Parse()
	utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 0)

	return program
}

func getExpr[T ast.Expression](t *testing.T, program *ast.Program) T {
	utils.AssertEq(t, len(program.Statements), 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	utils.Assert(t, ok)

	expr, ok := stmt.Expression.(T)
	utils.Assert(t, ok)
  return expr
}

func TestIdentifierExpression(t *testing.T) {
	input := "hello_World123"
  program := getProgram(t, input)

  ident := getExpr[*ast.Identifier](t, program)

	utils.AssertEq(t, ident.Name, input)
	utils.AssertEq(t, ident.Token,
		token.New(token.IDENTIFIER, input,
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestIntegerExpression(t *testing.T) {
	input := "156_098"
	val := 156098

  program := getProgram(t, input)

  integer := getExpr[*ast.IntegerLiteral](t, program)

	utils.AssertEq(t, integer.Value, int64(val))
	utils.AssertEq(t, integer.Token,
		token.New(token.INTEGER, "156098",
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestFloatExpression(t *testing.T) {
	input := "3.1415"
	val := 3.1415

  program := getProgram(t, input)

  float := getExpr[*ast.FloatLiteral](t, program)

	utils.AssertEq(t, float.Value, val)
	utils.AssertEq(t, float.Token,
		token.New(token.FLOAT, input,
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestBooleanExpression(t *testing.T) {
	input := "true"

  program := getProgram(t, input)

  boolean := getExpr[*ast.BooleanLiteral](t, program)

	utils.AssertEq(t, boolean.Value, true)
	utils.AssertEq(t, boolean.Token,
		token.New(token.IDENTIFIER, input,
			token.NewSpan(0, 0, len(input)),
		),
	)
}

func TestErrorExpression(t *testing.T) {
  input := ")"

  l := lexer.New(input, "test.lb")
  tokens := l.Tokenise()
  p := parser.New(tokens, l.Diagnostics)
  program := p.Parse()

  utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 1)
  diag := p.Diagnostics.Diagnostics[0]
  utils.AssertEq(t, diag.Message, "Expected expression, got `)`")
  utils.AssertEq(t, diag.Span, token.NewSpan(0, 0, 1))

  getExpr[*ast.ErrorExpression](t, program)
}

func TestMissingNewlineError(t *testing.T) {
  input := "1 2"

  l := lexer.New(input, "test.lb")
  tokens := l.Tokenise()
  p := parser.New(tokens, l.Diagnostics)
  p.Parse()

  utils.AssertEq(t, len(p.Diagnostics.Diagnostics), 1)
  diag := p.Diagnostics.Diagnostics[0]
  utils.AssertEq(t, diag.Message, "Expected newline after statement, got integer")
  utils.AssertEq(t, diag.Span, token.NewSpan(0, 2, 3))
}

