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
		utils.AssertEq(t, stmt.Name, tt.ident)

		if tt.ty != nil {
			utils.Assert(t, stmt.Type != nil, "Expected no type, but got one")
			typeName, ok := stmt.Type.(*ast.Identifier)
			utils.Assert(t, ok, "Type is not a type name")
			utils.AssertEq(t, typeName.Name, *tt.ty)
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
		stmt := getStmt[*ast.IfExpression](t, program)
		testIfStatement(t, stmt, tt.condition, tt.bodyValue, tt.elseBranch)
	}
}

func testIfStatement(t *testing.T, stmt *ast.IfExpression, condition any, bodyValue any, elseBranch *elseBranch) {
	testLiteral(t, stmt.Condition, condition)
	bodyStmt := utils.AssertSingle(t, stmt.Body.Statements)
	expr, ok := bodyStmt.(ast.Expression)
	utils.Assert(t, ok, "Body is not an expression")
	testLiteral(t, expr, bodyValue)

	if elseBranch != nil {
		utils.Assert(t, stmt.ElseBranch != nil, "Expected else branch")
		testElseBranch(t, stmt.ElseBranch, elseBranch)
	} else {
		utils.Assert(t, stmt.ElseBranch == nil, "Expected no else branch")
	}
}

func testElseBranch(t *testing.T, branch ast.Statement, expected *elseBranch) {
	if expected.condition == nil {
		block, ok := branch.(*ast.Block)
		utils.Assert(t, ok, "Else branch is not a block")
		bodyStmt := utils.AssertSingle(t, block.Statements)
		expr, ok := bodyStmt.(ast.Expression)
		utils.Assert(t, ok, "Body is not an expression")
		testLiteral(t, expr, expected.bodyValue)
	} else {
		ifStmt, ok := branch.(*ast.IfExpression)
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
		expr, ok := bodyStmt.(ast.Expression)
		utils.Assert(t, ok, "Body is not an expression")
		testLiteral(t, expr, tt.bodyValue)
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

		utils.AssertEq(t, loop.Variable, tt.ident)
		testLiteral(t, loop.Iterator, tt.iterator)
		bodyStmt := utils.AssertSingle(t, loop.Body.Statements)
		expr, ok := bodyStmt.(ast.Expression)
		utils.Assert(t, ok, "Body is not an expression")
		testLiteral(t, expr, tt.bodyValue)
	}
}

type parameter struct {
	mutable  bool
	name     string
	dataType string
	value    any
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
		{"fn (i32) add(\nother: i32\n,\n)\n:i32\n{ 7 }", "i32", false, "", "add", []parameter{{false, "other", "i32", nil}}, "i32", 7},
		{"fn u8.zero(): u8 {0}", "", false, "u8", "zero", []parameter{}, "u8", 0},
		{"fn sum(a,b,c:f64) : usize{ 3.14 }", "", false, "", "sum", []parameter{{false, "a", "", nil}, {false, "b", "", nil}, {false, "c", "f64", nil}}, "usize", 3.14},
		{"fn inc(mut x: u32): u32 { x }", "", false, "", "inc", []parameter{{true, "x", "u32", nil}}, "u32", "$x"},
		{"fn (mut foo) bar(): foo { this }", "foo", true, "", "bar", []parameter{}, "foo", "$this"},
		{"fn add(a = 1, mut b: i64 = 2): i64 { c }", "", false, "", "add", []parameter{{false, "a", "", 1}, {true, "b", "i64", 2}}, "i64", "$c"},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		fn := getStmt[*ast.FunctionDeclaration](t, program)

		if test.methodOf == "" {
			utils.Assert(t, fn.MethodOf == nil, "Expected not to be a method")
		} else {
			utils.Assert(t, fn.MethodOf != nil, "Expected to be a method")
			name, ok := fn.MethodOf.Type.(*ast.Identifier)
			utils.Assert(t, ok, "MethodOf is not a type name")
			utils.AssertEq(t, name.Name, test.methodOf)

			if test.thisMut {
				utils.Assert(t, fn.MethodOf.Mutable, "Expected this to be mutable")
			} else {
				utils.Assert(t, !fn.MethodOf.Mutable, "Expected this not to be mutable")
			}
		}

		if test.memberOf == "" {
			utils.Assert(t, fn.MemberOf == nil, "Expected not to be a member")
		} else {
			utils.Assert(t, fn.MemberOf != nil, "Expected to be a member")
			utils.AssertEq(t, fn.MemberOf.Name, test.memberOf)
		}

		utils.AssertEq(t, fn.Name, test.name)

		utils.AssertEq(t, len(fn.Parameters), len(test.params))
		for i, param := range test.params {
			fnParam := fn.Parameters[i]
			utils.AssertEq(t, *fnParam.Name, param.name)

			if param.dataType == "" {
				utils.Assert(t, fnParam.Type == nil, "Expected no type annotation")
			} else {
				utils.Assert(t, fnParam.Type != nil, "Expected a type annotation")
				name, ok := fnParam.Type.(*ast.Identifier)
				utils.Assert(t, ok, "Param type is not a type name")
				utils.AssertEq(t, name.Name, param.dataType)
			}

			if param.value == nil {
				utils.Assert(t, fnParam.Default == nil, "Expected no default value")
			} else {
				utils.Assert(t, fnParam.Default != nil, "Expected a default value")
				testLiteral(t, fnParam.Default, param.value)
			}

			if param.mutable {
				utils.Assert(t, fnParam.Mutable, "Expected param to be mutable")
			} else {
				utils.Assert(t, !fnParam.Mutable, "Expected param not to be mutable")
			}
		}

		if test.returnType == "" {
			utils.Assert(t, fn.ReturnType == nil, "Expected no return type")
		} else {
			utils.Assert(t, fn.ReturnType != nil, "Expected a return type")
			name, ok := fn.ReturnType.(*ast.Identifier)
			utils.Assert(t, ok, "ReturnType is not a type name")
			utils.AssertEq(t, name.Name, test.returnType)
		}

		bodyStmt := utils.AssertSingle(t, fn.Body.Statements)
		expr, ok := bodyStmt.(ast.Expression)
		utils.Assert(t, ok, "Body is not an expression")
		testLiteral(t, expr, test.bodyValue)
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

func TestBreakStatement(t *testing.T) {
	tests := []struct {
		src   string
		value any
	}{
		{"break", nil},
		{"break true", true},
		{"break [1,2,3]", []any{1, 2, 3}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		brk := getStmt[*ast.BreakStatement](t, program)

		if test.value == nil {
			utils.Assert(t, brk.Value == nil, "Expected no value value")
		} else {
			testLiteral(t, brk.Value, test.value)
		}
	}
}

func TestYieldStatement(t *testing.T) {
	tests := []struct {
		src   string
		value any
	}{
		{"yield 73", 73},
		{`yield "foo"`, "foo"},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		yield := getStmt[*ast.YieldStatement](t, program)

		testLiteral(t, yield.Value, test.value)
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

		utils.AssertEq(t, td.Name, test.name)
		typeName, ok := td.Type.(*ast.Identifier)
		utils.Assert(t, ok, "Type is not a type name")
		utils.AssertEq(t, typeName.Name, test.typeName)
	}
}

func TestStructDeclaration(t *testing.T) {
	tests := []struct {
		src    string
		name   string
		fields any
	}{
		{"struct Unit", "Unit", nil},
		{"struct Mything123", "Mything123", nil},
		{"struct Wrapper { value }", "Wrapper", []string{"value"}},
		{"struct Three{a,b,c,}", "Three", []string{"a", "b", "c"}},
		{"struct Empty {}", "Empty", [][2]string{}},
		{"struct Rect { w, h: i32 }", "Rect", [][2]string{{"w"}, {"h", "i32"}}},
		{"struct Vec2{x:f32,y:f32,}", "Vec2", [][2]string{{"x", "f32"}, {"y", "f32"}}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		sd := getStmt[*ast.StructDeclaration](t, program)

		utils.AssertEq(t, sd.Name, test.name)

		switch fields := test.fields.(type) {
		case [][2]string:
			utils.Assert(t, sd.Body != nil, "Expected a non-unit struct")
			utils.AssertEq(t, len(sd.Body), len(fields), "Field lengths do not match")
			for i, field := range fields {
				structField := sd.Body[i]
				utils.Assert(t, structField.Name != nil, "Expected named fields")
				utils.AssertEq(t, field[0], *structField.Name)

				if field[1] == "" {
					utils.Assert(t, structField.Type == nil, "Expected no type annotation")
				} else {
					utils.Assert(t, structField.Type != nil, "Expected a type annotation")
					typeName, ok := structField.Type.(*ast.Identifier)
					utils.Assert(t, ok, "Type is not a type name")
					utils.AssertEq(t, typeName.Name, field[1])
				}
			}

		case []string:
			utils.Assert(t, sd.Body != nil, "Expected a non-unit struct")
			utils.AssertEq(t, len(sd.Body), len(fields), "Type lengths do not match")

			for i, ty := range fields {
				structField := sd.Body[i]
				var typeName string
				if structField.Name != nil {
					typeName = *structField.Name
					utils.Assert(t, structField.Type == nil, "Expected unnamed fields")
				} else {
					name, ok := structField.Type.(*ast.Identifier)
					utils.Assert(t, ok, "Type is not a type name")
					typeName = name.Name
				}
				utils.AssertEq(t, typeName, ty)
			}

		case nil:
			utils.Assert(t, sd.Body == nil, "Expected a unit struct")
		default:
			panic("Invalid test")
		}
	}
}

type interfaceField struct {
	name       string
	params     []string
	returnType string
}

func TestInterfaceDeclaration(t *testing.T) {
	tests := []struct {
		src    string
		name   string
		fields []interfaceField
	}{
		{"interface Any {}", "Any", []interfaceField{}},
		{"interface Fooer { foo(bar): baz }", "Fooer", []interfaceField{{"foo", []string{"bar"}, "baz"}}},
		{`interface Order {
			less ( i32 , f64 ) : bool , 
			greater(u32,i32,):f16
		}`, "Order",
			[]interfaceField{
				{"less", []string{"i32", "f64"}, "bool"},
				{"greater", []string{"u32", "i32"}, "f16"},
			},
		},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		intDecl := getStmt[*ast.InterfaceDeclaration](t, program)

		utils.AssertEq(t, intDecl.Name, test.name)
		utils.AssertEq(t, len(intDecl.Members), len(test.fields), "Incorrect number of fields")

		for i, field := range test.fields {
			intField := intDecl.Members[i]
			utils.AssertEq(t, intField.Name, field.name)
			utils.AssertEq(t, len(intField.Parameters), len(field.params), "Incorrect number of parameters")
			for i, param := range field.params {
				intParam := intField.Parameters[i]
				name, ok := intParam.(*ast.Identifier)
				utils.Assert(t, ok, "Parameter is not a type name")
				utils.AssertEq(t, name.Name, param)
			}

			if field.returnType == "" {
				utils.Assert(t, intField.ReturnType == nil, "Expected no return type")
			} else {
				utils.Assert(t, intField.ReturnType != nil, "Expected return type")
				name, ok := intField.ReturnType.(*ast.Identifier)
				utils.Assert(t, ok, "Return type is not a type name")
				utils.AssertEq(t, name.Name, field.returnType)
			}
		}
	}
}

func TestImportStatement(t *testing.T) {
	tests := []struct {
		src     string
		symbols []string
		all     bool
		module  string
		alias   string
	}{
		{`import "fs"`, nil, false, "fs", ""},
		{`import ".././foo/bar"`, nil, false, ".././foo/bar", ""},
		{`import * from "helpers"`, nil, true, "helpers", ""},
		{`import { read, write } from "io"`, []string{"read", "write"}, false, "io", ""},
		{`import "42" as life_universe_everything`, nil, false, "42", "life_universe_everything"},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		stmt := getStmt[*ast.ImportStatement](t, program)

		if test.symbols == nil {
			utils.Assert(t, stmt.Symbols == nil, "Expected no imported symbols")
		} else {
			utils.Assert(t, stmt.Symbols != nil, "Expected imported symbols")
			utils.AssertEq(t, len(stmt.Symbols), len(test.symbols))
			for i, symbol := range test.symbols {
				imported := stmt.Symbols[i]
				utils.AssertEq(t, imported.Name, symbol)
			}
		}

		if test.all {
			utils.Assert(t, stmt.All, "Expected to import all symbols")
		} else {
			utils.Assert(t, !stmt.All, "Expected not to import all symbols")
		}

		utils.AssertEq(t, stmt.Module.ExtraValue, test.module)

		if test.alias == "" {
			utils.Assert(t, stmt.Alias == nil, "Expected not to import as an alias")
		} else {
			utils.Assert(t, stmt.Alias != nil, "Expected to import as an alias")
			utils.AssertEq(t, *stmt.Alias, test.alias)
		}
	}
}

type enumField struct {
	name  string
	value any
}

func TestEnumDeclaration(t *testing.T) {
	tests := []struct {
		src      string
		name     string
		dataType string
		members  []enumField
	}{
		{"enum Empty {}", "Empty", "", []enumField{}},
		{"enum Colour: u64 { Invalid, red = 100, green = 783, blue = 1.5 }", "Colour", "u64", []enumField{
			{"Invalid", nil},
			{"red", 100},
			{"green", 783},
			{"blue", 1.5},
		}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		stmt := getStmt[*ast.EnumDeclaration](t, program)

		utils.AssertEq(t, stmt.Name, test.name)

		if test.dataType == "" {
			utils.Assert(t, stmt.ValueType == nil, "Expected no type annotation")
		} else {
			utils.Assert(t, stmt.ValueType != nil, "Expected type annotation")

			name, ok := stmt.ValueType.(*ast.Identifier)
			utils.Assert(t, ok, "Type is not a type name")
			utils.AssertEq(t, name.Name, test.dataType)
		}

		utils.AssertEq(t, len(stmt.Members), len(test.members))

		for i, expected := range test.members {
			enumMember := stmt.Members[i]

			utils.AssertEq(t, enumMember.Name, expected.name)

			if expected.value == nil {
				utils.Assert(t, enumMember.Value == nil, "Expected no value")
			} else {
				utils.Assert(t, enumMember.Value != nil, "Expected a value")

				testLiteral(t, enumMember.Value, expected.value)
			}

		}
	}
}

type unionField struct {
	name string
	ty   any
}

func TestUnionDeclaration(t *testing.T) {
	tests := []struct {
		src     string
		name    string
		members []unionField
	}{
		{"union AOrB { a, b }", "AOrB", []unionField{
			{"a", nil},
			{"b", nil},
		}},
		{"union Int { i8, i16, i32, i64 ,}", "Int", []unionField{
			{"i8", nil},
			{"i16", nil},
			{"i32", nil},
			{"i64", nil},
		}},
		{"union Property { Age: i32, Height: f32, Weight:f32,string}", "Property", []unionField{
			{"Age", "i32"},
			{"Height", "f32"},
			{"Weight", "f32"},
			{"string", nil},
		}},
		{"union Shape { Square { f32, f32 }, Circle { radius: f32 } }", "Shape", []unionField{
			{"Square", []string{"f32", "f32"}},
			{"Circle", [][2]string{{"radius", "f32"}}},
		}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		stmt := getStmt[*ast.UnionDeclaration](t, program)

		utils.AssertEq(t, stmt.Name, test.name)

		utils.AssertEq(t, len(stmt.Members), len(test.members))

		for i, expected := range test.members {
			unionMember := stmt.Members[i]

			utils.AssertEq(t, unionMember.Name, expected.name)
			switch ty := expected.ty.(type) {
			case string:
				utils.Assert(t, unionMember.Type != nil)
				utils.Assert(t, unionMember.Compound == nil)

				name, ok := unionMember.Type.(*ast.Identifier)
				utils.Assert(t, ok, "Type is not a type name")
				utils.AssertEq(t, name.Name, ty)

			case [][2]string:
				utils.Assert(t, unionMember.Type == nil)
				utils.Assert(t, unionMember.Compound != nil)
				utils.AssertEq(t, len(unionMember.Compound), len(ty), "Field lengths do not match")
				for i, field := range ty {
					structField := unionMember.Compound[i]
					utils.Assert(t, structField.Name != nil, "Expected named fields")
					utils.AssertEq(t, field[0], *structField.Name)

					if field[1] == "" {
						utils.Assert(t, structField.Type == nil, "Expected no type annotation")
					} else {
						utils.Assert(t, structField.Type != nil, "Expected a type annotation")
						typeName, ok := structField.Type.(*ast.Identifier)
						utils.Assert(t, ok, "Type is not a type name")
						utils.AssertEq(t, typeName.Name, field[1])
					}
				}

			case []string:
				utils.Assert(t, unionMember.Type == nil)
				utils.Assert(t, unionMember.Compound != nil)
				utils.AssertEq(t, len(unionMember.Compound), len(ty), "Type lengths do not match")

				for i, ty := range ty {
					structField := unionMember.Compound[i]
					var typeName string
					if structField.Name != nil {
						typeName = *structField.Name
						utils.Assert(t, structField.Type == nil, "Expected unnamed fields")
					} else {
						name, ok := structField.Type.(*ast.Identifier)
						utils.Assert(t, ok, "Type is not a type name")
						typeName = name.Name
					}
					utils.AssertEq(t, typeName, ty)
				}

			case nil:
				utils.Assert(t, unionMember.Type == nil)
				utils.Assert(t, unionMember.Compound == nil)
			}
		}
	}
}

func TestTagDeclaration(t *testing.T) {
	tests := []struct {
		src  string
		name string
	}{
		{"tag MyTag", "MyTag"},
		{"tag Test124", "Test124"},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		stmt := getStmt[*ast.TagDeclaration](t, program)

		utils.AssertEq(t, stmt.Name, test.name)
	}
}
