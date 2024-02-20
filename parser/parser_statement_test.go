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

func TestForLoop(t *testing.T) {
	tests := []struct {
		src       string
		ident     string
		iterator  any
		bodyValue any
	}{
		{"for i in [1,2,3] { i }", "i", []any{1, 2, 3}, "$i"},
		{"for foo in 93\n{[foo,bar,]}", "foo", 93, []any{"$foo", "$bar"}},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.src)
		loop := getStmt[*ast.ForLoop](t, program)

		utils.AssertEq(t, loop.Variable.Value, tt.ident)
		testLiteral(t, loop.Iterator, tt.iterator)
		bodyStmt := utils.AssertSingle(t, loop.Body.Statements)
		exprStmt, ok := bodyStmt.(*ast.ExpressionStatement)
		utils.Assert(t, ok, "Body is not an expression statement")
		testLiteral(t, exprStmt.Expression, tt.bodyValue)
	}
}

func TestFunctionDeclaration(t *testing.T) {
	tests := []struct {
		src        string
		methodOf   string
		memberOf   string
		name       string
		params     [][2]string
		returnType string
		bodyValue  any
	}{
		{`fn hello() { "Hello, world!" }`, "", "", "hello", [][2]string{}, "", "Hello, world!"},
		{"fn (i32) print() {\nthis\n}", "i32", "", "print", [][2]string{}, "", "$this"},
		{"fn (i32) add(\nother: i32\n,\n)\n:i32\n{ 7 }", "i32", "", "add", [][2]string{{"other", "i32"}}, "i32", 7},
		{"fn u8.zero(): u8 {0}", "", "u8", "zero", [][2]string{}, "u8", 0},
		{"fn sum(a,b,c:f64) : usize{ 3.14 }", "", "", "sum", [][2]string{{"a"}, {"b"}, {"c", "f64"}}, "usize", 3.14},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		fn := getStmt[*ast.FunctionDeclaration](t, program)

		if test.methodOf == "" {
			utils.Assert(t, fn.MethodOf == nil, "Expected not to be a method")
		} else {
			utils.Assert(t, fn.MethodOf != nil, "Expected to be a method")
			name, ok := fn.MethodOf.Type.(*ast.TypeName)
			utils.Assert(t, ok, "MethodOf is not a type name")
			utils.AssertEq(t, name.Name.Value, test.methodOf)
		}

		if test.memberOf == "" {
			utils.Assert(t, fn.MemberOf == nil, "Expected not to be a member")
		} else {
			utils.Assert(t, fn.MemberOf != nil, "Expected to be a member")
			utils.AssertEq(t, fn.MemberOf.Name.Value, test.memberOf)
		}

		utils.AssertEq(t, fn.Name.Value, test.name)

		utils.AssertEq(t, len(fn.Parameters), len(test.params))
		for i, param := range test.params {
			fnParam := fn.Parameters[i]
			utils.AssertEq(t, fnParam.Name.Value, param[0])

			if param[1] == "" {
				utils.Assert(t, fnParam.Type == nil, "Expected no type annotation")
			} else {
				utils.Assert(t, fnParam.Type != nil, "Expected a type annotation")
				name, ok := fnParam.Type.Type.(*ast.TypeName)
				utils.Assert(t, ok, "Param type is not a type name")
				utils.AssertEq(t, name.Name.Value, param[1])
			}
		}

		if test.returnType == "" {
			utils.Assert(t, fn.ReturnType == nil, "Expected no return type")
		} else {
			utils.Assert(t, fn.ReturnType != nil, "Expected a return type")
			name, ok := fn.ReturnType.Type.(*ast.TypeName)
			utils.Assert(t, ok, "ReturnType is not a type name")
			utils.AssertEq(t, name.Name.Value, test.returnType)
		}

		bodyStmt := utils.AssertSingle(t, fn.Body.Statements)
		exprStmt, ok := bodyStmt.(*ast.ExpressionStatement)
		utils.Assert(t, ok, "Body is not an expression statement")
		testLiteral(t, exprStmt.Expression, test.bodyValue)
	}
}
