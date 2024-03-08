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

			if param.value == nil {
				utils.Assert(t, fnParam.Default == nil, "Expected no default value")
			} else {
				utils.Assert(t, fnParam.Default != nil, "Expected a default value")
				testLiteral(t, fnParam.Default.Value, param.value)
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
		fields any
	}{
		{"struct Unit", "Unit", nil},
		{"struct Mything123", "Mything123", nil},
		{"struct Wrapper ( value )", "Wrapper", []string{"value"}},
		{"struct Three(a,b,c,)", "Three", []string{"a", "b", "c"}},
		{"struct Empty {}", "Empty", [][2]string{}},
		{"struct Rect { w, h: i32 }", "Rect", [][2]string{{"w"}, {"h", "i32"}}},
		{"struct Vec2{x:f32,y:f32,}", "Vec2", [][2]string{{"x", "f32"}, {"y", "f32"}}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		sd := getStmt[*ast.StructDeclaration](t, program)

		utils.AssertEq(t, sd.Name.Value, test.name)

		switch fields := test.fields.(type) {
		case [][2]string:
			utils.Assert(t, sd.StructType != nil, "Expected a curly-brace struct")
			utils.Assert(t, sd.TupleType == nil, "Cannot be both curly-brace and tupe struct")
			utils.AssertEq(t, len(sd.StructType.Fields), len(fields), "Field lengths do not match")
			for i, field := range fields {
				structField := sd.StructType.Fields[i]
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

		case []string:
			utils.Assert(t, sd.TupleType != nil, "Expected a tuple struct")
			utils.Assert(t, sd.StructType == nil, "Cannot be both curly-brace and tupe struct")
			utils.AssertEq(t, len(sd.TupleType.Types), len(fields), "Type lengths do not match")

			for i, ty := range fields {
				tupleType := sd.TupleType.Types[i]
				typeName, ok := tupleType.(*ast.TypeName)
				utils.Assert(t, ok, "Type is not a type name")
				utils.AssertEq(t, typeName.Name.Value, ty)
			}

		case nil:
			utils.Assert(t, sd.StructType == nil, "Expected a unit struct")
			utils.Assert(t, sd.TupleType == nil, "Expected a unit struct")
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

		utils.AssertEq(t, intDecl.Name.Value, test.name)
		utils.AssertEq(t, len(intDecl.Members), len(test.fields), "Incorrect number of fields")

		for i, field := range test.fields {
			intField := intDecl.Members[i]
			utils.AssertEq(t, intField.Name.Value, field.name)
			utils.AssertEq(t, len(intField.Parameters), len(field.params), "Incorrect number of parameters")
			for i, param := range field.params {
				intParam := intField.Parameters[i]
				name, ok := intParam.(*ast.TypeName)
				utils.Assert(t, ok, "Parameter is not a type name")
				utils.AssertEq(t, name.Name.Value, param)
			}

			if field.returnType == "" {
				utils.Assert(t, intField.ReturnType == nil, "Expected no return type")
			} else {
				utils.Assert(t, intField.ReturnType != nil, "Expected return type")
				name, ok := intField.ReturnType.Type.(*ast.TypeName)
				utils.Assert(t, ok, "Return type is not a type name")
				utils.AssertEq(t, name.Name.Value, field.returnType)
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
			utils.AssertEq(t, len(stmt.Symbols.Symbols), len(test.symbols))
			for i, symbol := range test.symbols {
				imported := stmt.Symbols.Symbols[i]
				utils.AssertEq(t, imported.Value, symbol)
			}
		}

		if test.all {
			utils.Assert(t, stmt.All != nil, "Expected to import all symbols")
		} else {
			utils.Assert(t, stmt.All == nil, "Expected not to import all symbols")
		}

		utils.AssertEq(t, stmt.Module.Value, test.module)

		if test.alias == "" {
			utils.Assert(t, stmt.Alias == nil, "Expected not to import as an alias")
		} else {
			utils.Assert(t, stmt.Alias != nil, "Expected to import as an alias")
			utils.AssertEq(t, stmt.Alias.Alias.Value, test.alias)
		}
	}
}

type enumField struct {
	name string
	data any
}

func TestEnumDeclaration(t *testing.T) {
	tests := []struct {
		src      string
		name     string
		dataType string
		members  []enumField
	}{
		{"enum Empty {}", "Empty", "", []enumField{}},
		{"enum Option { None, Some(i32) }", "Option", "", []enumField{
			{"None", nil},
			{"Some", []string{"i32"}},
		}},
		{"union AOrB { a {b: c}, b {c: d, e:f,}, }", "AOrB", "", []enumField{
			{"a", [][2]string{{"b", "c"}}},
			{"b", [][2]string{{"c", "d"}, {"e", "f"}}},
		}},
		{"enum Colour: u64 { Invalid, red = 100, green = 783, blue = 1.5, custom(i32, u32, f32) }", "Colour", "u64", []enumField{
			{"Invalid", nil},
			{"red", 100},
			{"green", 783},
			{"blue", 1.5},
			{"custom", []string{"i32", "u32", "f32"}},
		}},
		{"union number { i32, u32, f32 }", "number", "", []enumField{
			{"i32", nil},
			{"u32", nil},
			{"f32", nil},
		}},
	}

	for _, test := range tests {
		program := getProgram(t, test.src)
		stmt := getStmt[*ast.EnumDeclaration](t, program)

		utils.AssertEq(t, stmt.Name.Value, test.name)

		if test.dataType == "" {
			utils.Assert(t, stmt.ValueType == nil, "Expected no type annotation")
		} else {
			utils.Assert(t, stmt.ValueType != nil, "Expected type annotation")

			name, ok := stmt.ValueType.Type.(*ast.TypeName)
			utils.Assert(t, ok, "Type is not a type name")
			utils.AssertEq(t, name.Name.Value, test.dataType)
		}

		utils.AssertEq(t, len(stmt.Members), len(test.members))

		for i, expected := range test.members {
			enumMember := stmt.Members[i]

			utils.AssertEq(t, enumMember.Name.Value, expected.name)

			switch member := expected.data.(type) {
			case nil:
				utils.Assert(t, enumMember.Types == nil, "Expected a unit enum member")
				utils.Assert(t, enumMember.Struct == nil, "Expected a unit enum member")
				utils.Assert(t, enumMember.Value == nil, "Expected a unit enum member")

			case []string:
				utils.Assert(t, enumMember.Types != nil, "Expected a tuple enum member")
				utils.Assert(t, enumMember.Struct == nil, "Expected a tuple enum member")
				utils.Assert(t, enumMember.Value == nil, "Expected a tuple enum member")
				types := enumMember.Types.Types

				utils.AssertEq(t, len(types), len(member))
				for i, ty := range member {
					name, ok := types[i].(*ast.TypeName)
					utils.Assert(t, ok, "Type is not a type name")
					utils.AssertEq(t, name.Name.Value, ty)
				}

			case [][2]string:
				utils.Assert(t, enumMember.Types == nil, "Expected a struct enum member")
				utils.Assert(t, enumMember.Struct != nil, "Expected a struct enum member")
				utils.Assert(t, enumMember.Value == nil, "Expected a struct enum member")
				fields := enumMember.Struct.Fields

				utils.AssertEq(t, len(fields), len(member))
				for i, field := range member {
					utils.AssertEq(t, fields[i].Name.Value, field[0])

					name, ok := fields[i].Type.Type.(*ast.TypeName)
					utils.Assert(t, ok, "Type is not a type name")
					utils.AssertEq(t, name.Name.Value, field[1])
				}

			default:
				utils.Assert(t, enumMember.Types == nil, "Expected a value enum member")
				utils.Assert(t, enumMember.Struct == nil, "Expected a value enum member")
				utils.Assert(t, enumMember.Value != nil, "Expected a value enum member")

				testLiteral(t, enumMember.Value.Value, member)
			}
		}
	}
}
