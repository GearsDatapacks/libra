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

type parameter struct {
	mutable  bool
	name     string
	dataType string
}

func TestFunctionDeclaration(t *testing.T) {
	tests := []struct {
		src        string
		methodOf   string
		thisMut    bool
		memberOf   string
		name       string
		params     []parameter
		returnType string
		bodyValue  any
	}{
		{`fn hello() { "Hello, world!" }`, "", false, "", "hello", []parameter{}, "", "Hello, world!"},
		{"fn (i32) print() {\nthis\n}", "i32", false, "", "print", []parameter{}, "", "$this"},
		{"fn (i32) add(\nother: i32\n,\n)\n:i32\n{ 7 }", "i32", false, "", "add", []parameter{{false, "other", "i32"}}, "i32", 7},
		{"fn u8.zero(): u8 {0}", "", false, "u8", "zero", []parameter{}, "u8", 0},
		{"fn sum(a,b,c:f64) : usize{ 3.14 }", "", false, "", "sum", []parameter{{false, "a", ""}, {false, "b", ""}, {false, "c", "f64"}}, "usize", 3.14},
		{"fn inc(mut x: u32): u32 { x }", "", false, "", "inc", []parameter{{true, "x", "u32"}}, "u32", "$x"},
		{"fn (mut foo) bar(): foo { this }", "foo", true, "", "bar", []parameter{}, "foo", "$this"},
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

			if test.thisMut {
				utils.Assert(t, fn.MethodOf.Mutable != nil, "Expected this to be mutable")
			} else {
				utils.Assert(t, fn.MethodOf.Mutable == nil, "Expected this not to be mutable")
			}
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
			utils.AssertEq(t, fnParam.Name.Value, param.name)

			if param.dataType == "" {
				utils.Assert(t, fnParam.Type == nil, "Expected no type annotation")
			} else {
				utils.Assert(t, fnParam.Type != nil, "Expected a type annotation")
				name, ok := fnParam.Type.Type.(*ast.TypeName)
				utils.Assert(t, ok, "Param type is not a type name")
				utils.AssertEq(t, name.Name.Value, param.dataType)
			}

			if param.mutable {
				utils.Assert(t, fnParam.Mutable != nil, "Expected param to be mutable")
			} else {
				utils.Assert(t, fnParam.Mutable == nil, "Expected param not to be mutable")
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

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		src   string
		value any
	}{
		{"return", nil},
		{"return 7", 7},
		{"return false", false},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		ret := getStmt[*ast.ReturnStatement](t, program)

		if test.value == nil {
			utils.Assert(t, ret.Value == nil, "Expected no return value")
		} else {
			testLiteral(t, ret.Value, test.value)
		}
	}
}

func TestTypeDeclaration(t *testing.T) {
	tests := []struct {
		src      string
		name     string
		typeName string
	}{
		{"type foo = bar", "foo", "bar"},
		{"type int=i32", "int", "i32"},
		{"type boolean\n =\n bool", "boolean", "bool"},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		td := getStmt[*ast.TypeDeclaration](t, program)

		utils.AssertEq(t, td.Name.Value, test.name)
		typeName, ok := td.Type.(*ast.TypeName)
		utils.Assert(t, ok, "Type is not a type name")
		utils.AssertEq(t, typeName.Name.Value, test.typeName)
	}
}

func TestStructDeclaration(t *testing.T) {
	tests := []struct {
		src    string
		name   string
		fields [][2]string
	}{
		{"struct Empty {}", "Empty", [][2]string{}},
		{"struct Rect { w, h: i32 }", "Rect", [][2]string{{"w"}, {"h", "i32"}}},
		{"struct Vec2{x:f32,y:f32,}", "Vec2", [][2]string{{"x", "f32"}, {"y", "f32"}}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		sd := getStmt[*ast.StructDeclaration](t, program)

		utils.AssertEq(t, sd.Name.Value, test.name)

		utils.AssertEq(t, len(sd.Fields), len(test.fields), "Field lengths do not match")
		for i, field := range test.fields {
			structField := sd.Fields[i]
			utils.AssertEq(t, field[0], structField.Name.Value)

			if field[1] == "" {
				utils.Assert(t, structField.Type == nil, "Expected no type annotation")
			} else {
				utils.Assert(t, structField.Type != nil, "Expected a type annotation")
				typeName, ok := structField.Type.Type.(*ast.TypeName)
				utils.Assert(t, ok, "Type is not a type name")
				utils.AssertEq(t, typeName.Name.Value, field[1])
			}
		}
	}
}
