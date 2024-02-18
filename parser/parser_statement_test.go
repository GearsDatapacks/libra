package parser_test

import (
	"fmt"
	"testing"

	"github.com/gearsdatapacks/libra/parser/ast"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestVariableDeclaration(t *testing.T) {
	f32 := "f32"
	str := "string"

	tests := []struct {
		src     string
		keyword string
		ident   string
		ty      *string
		value   any
	}{
		{"let x = 1", "let", "x", nil, 1},
		{"mut y: f32 = 7", "mut", "y", &f32, 7},
		{`const message: string = "Hi"`, "const", "message", &str, "Hi"},
		{"mut isCool = true", "mut", "isCool", nil, true},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		stmt := getStmt[*ast.VariableDeclaration](t, program)

		utils.AssertEq(t, stmt.Keyword.Value, tt.keyword)
		utils.AssertEq(t, stmt.Identifier.Value, tt.ident)

		if tt.ty != nil {
			utils.Assert(t, stmt.Type != nil, "Expected no type, but got one")
			typeName, ok := stmt.Type.Type.(*ast.TypeName)
			utils.Assert(t, ok, "Type is not a type name")
			utils.AssertEq(t, typeName.Name.Value, *tt.ty)
		}

		testLiteral(t, stmt.Value, tt.value)
	}
}

type elseBranch struct {
	condition  any
	bodyValue  any
	elseBranch *elseBranch
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		src        string
		condition  any
		bodyValue  any
		elseBranch *elseBranch
	}{
		{"if a { 10 }", "$a", 10, nil},
		{"if false { 10 } else { 20 }", false, 10, &elseBranch{nil, 20, nil}},
		{`if 69
		{"Nice"}
		else if 42 { "UATLTUAE" }else{
			"Boring"
		}`, 69, "Nice", &elseBranch{42, "UATLTUAE", &elseBranch{nil, "Boring", nil}}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		stmt := getStmt[*ast.IfStatement](t, program)
		testIfStatement(t, stmt, tt.condition, tt.bodyValue, tt.elseBranch)
	}
}

func testIfStatement(t *testing.T, stmt *ast.IfStatement, condition any, bodyValue any, elseBranch *elseBranch) {
	testLiteral(t, stmt.Condition, condition)
	bodyStmt := utils.AssertSingle(t, stmt.Body.Statements)
	exprStmt, ok := bodyStmt.(*ast.ExpressionStatement)
	utils.Assert(t, ok, "Body is not an expression statement")
	testLiteral(t, exprStmt.Expression, bodyValue)

	if elseBranch != nil {
		utils.Assert(t, stmt.ElseBranch != nil, "Expected else branch")
		testElseBranch(t, stmt.ElseBranch, elseBranch)
	} else {
		utils.Assert(t, stmt.ElseBranch == nil, "Expected no else branch")
	}
}

func testElseBranch(t *testing.T, branch *ast.ElseBranch, expected *elseBranch) {
	if expected.condition == nil {
		block, ok := branch.Statement.(*ast.BlockStatement)
		utils.Assert(t, ok, "Else branch is not a block")
		bodyStmt := utils.AssertSingle(t, block.Statements)
		exprStmt, ok := bodyStmt.(*ast.ExpressionStatement)
		utils.Assert(t, ok, "Body is not an expression statement")
		testLiteral(t, exprStmt.Expression, expected.bodyValue)
	} else {
		ifStmt, ok := branch.Statement.(*ast.IfStatement)
		utils.Assert(t, ok, "Else branch is not an if statement")
		testIfStatement(t, ifStmt, expected.condition, expected.bodyValue, expected.elseBranch)
	}
}

func getStmt[T ast.Statement](t *testing.T, program *ast.Program) T {
	t.Helper()

	utils.AssertEq(t, len(program.Statements), 1,
		fmt.Sprintf("Program does not contain one statement. (has %d)",
			len(program.Statements)))

	stmt, ok := program.Statements[0].(T)
	utils.Assert(t, ok, fmt.Sprintf(
		"Statement is not an %T (is %T)", struct{ t T }{}.t, program.Statements[0]))

	return stmt
}

func TestWhileLoop(t *testing.T) {
	tests := []struct {
		src       string
		condition any
		bodyValue any
	}{
		{"while true { nop }", true, "$nop"},
		{`while thing { "Hi" }`, "$thing", "Hi"},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		loop := getStmt[*ast.WhileLoop](t, program)

		testLiteral(t, loop.Condition, tt.condition)
		bodyStmt := utils.AssertSingle(t, loop.Body.Statements)
		exprStmt, ok := bodyStmt.(*ast.ExpressionStatement)
		utils.Assert(t, ok, "Body is not an expression statement")
		testLiteral(t, exprStmt.Expression, tt.bodyValue)
	}
}
