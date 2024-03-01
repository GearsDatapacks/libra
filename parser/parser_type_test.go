package parser_test

import (
	"fmt"
	"testing"

	"github.com/gearsdatapacks/libra/parser/ast"
	utils "github.com/gearsdatapacks/libra/test_utils"
)

func TestTypeName(t *testing.T) {
	tests := []struct {
		src  string
		name string
	}{
		{"i32", "i32"},
		{"string", "string"},
		{"Option", "Option"},
	}

	for _, test := range tests {
		ty := parseType[*ast.TypeName](t, test.src)
		utils.AssertEq(t, ty.Name.Value, test.name)
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		src   string
		types []string
	}{
		{"a | b | c", []string{"a", "b", "c"}},
		{"i32 | f32 | u32", []string{"i32", "f32", "u32"}},
		{"string | cstring", []string{"string", "cstring"}},
	}

	for _, test := range tests {
		ty := parseType[*ast.Union](t, test.src)
		utils.AssertEq(t, len(ty.Types), len(test.types))
		for i, expected := range test.types {
			name, ok := ty.Types[i].(*ast.TypeName)
			utils.Assert(t, ok, "Type is not a type name")
			utils.AssertEq(t, name.Name.Value, expected)
		}
	}
}

func parseType[T ast.TypeExpression](t *testing.T, src string) T {
	program := getProgram(t, "type _ = "+src)
	ty := getType[T](t, program)
	return ty
}

func getType[T ast.TypeExpression](t *testing.T, program *ast.Program) T {
	t.Helper()

	stmt := getStmt[*ast.TypeDeclaration](t, program)
	ty, ok := stmt.Type.(T)
	utils.Assert(t, ok, fmt.Sprintf(
		"Type is not %T (is %T)", struct{ t T }{}.t, program.Statements[0]))

	return ty
}
