package lowerer

import (
	"fmt"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/type_checker/ir"
)

type lowerer struct {
	diagnostics diagnostics.Manager
}

func Lower(pkg *ir.Package, diagnostics diagnostics.Manager) (*ir.Package, diagnostics.Manager) {
	lowerer := lowerer{
		diagnostics: diagnostics,
	}

	lowered := &ir.Package{
		Modules: map[string]*ir.Module{},
	}

	for name, module := range pkg.Modules {
		mod := &ir.Module{
			Name:       name,
			Statements: []ir.Statement{},
		}
		lowered.Modules[name] = mod

		for _, stmt := range module.Statements {
			mod.Statements = append(mod.Statements, lowerer.lower(stmt))
		}
	}
	return lowered, lowerer.diagnostics
}

func (l *lowerer) lower(statement ir.Statement) ir.Statement {
	switch stmt := statement.(type) {
	case *ir.VariableDeclaration:
		return l.lowerVariableDeclaration(stmt)
	case *ir.FunctionDeclaration:
		return l.lowerFunctionDeclaration(stmt)
	case *ir.ReturnStatement:
		return l.lowerReturnStatement(stmt)
	case *ir.BreakStatement:
		return l.lowerBreakStatement(stmt)
	case *ir.YieldStatement:
		return l.lowerYieldStatement(stmt)
	case *ir.ImportStatement:
		return l.lowerImportStatement(stmt)
	case *ir.TypeDeclaration:
		return l.lowerTypeDeclaration(stmt)

	case ir.Expression:
		return l.lowerExpression(stmt)

	default:
		panic(fmt.Sprintf("TODO: lower %T", stmt))
	}
}

func (l *lowerer) lowerExpression(expression ir.Expression) ir.Expression {
	switch expr := expression.(type) {
	case *ir.IntegerLiteral:
		return l.lowerIntegerLiteral(expr)
	case *ir.FloatLiteral:
		return l.lowerFloatLiteral(expr)
	case *ir.BooleanLiteral:
		return l.lowerBooleanLiteral(expr)
	case *ir.StringLiteral:
		return l.lowerStringLiteral(expr)
	case *ir.VariableExpression:
		return l.lowerVariableExpression(expr)
	case *ir.BinaryExpression:
		return l.lowerBinaryExpression(expr)
	case *ir.UnaryExpression:
		return l.lowerUnaryExpression(expr)
	case *ir.Conversion:
		return l.lowerConversion(expr)
	case *ir.InvalidExpression:
		return l.lowerInvalidExpression(expr)
	case *ir.ArrayExpression:
		return l.lowerArrayExpression(expr)
	case *ir.IndexExpression:
		return l.lowerIndexExpression(expr)
	case *ir.MapExpression:
		return l.lowerMapExpression(expr)
	case *ir.Assignment:
		return l.lowerAssignment(expr)
	case *ir.TupleExpression:
		return l.lowerTupleExpression(expr)
	case *ir.TypeCheck:
		return l.lowerTypeCheck(expr)
	case *ir.FunctionCall:
		return l.lowerFunctionCall(expr)
	case *ir.StructExpression:
		return l.lowerStructExpression(expr)
	case *ir.TupleStructExpression:
		return l.lowerTupleStructExpression(expr)
	case *ir.MemberExpression:
		return l.lowerMemberExpression(expr)
	case *ir.Block:
		return l.lowerBlock(expr)
	case *ir.IfExpression:
		return l.lowerIfExpression(expr)
	case *ir.WhileLoop:
		return l.lowerWhileLoop(expr)
	case *ir.ForLoop:
		return l.lowerForLoop(expr)
	case *ir.TypeExpression:
		return l.lowerTypeExpression(expr)
	case *ir.FunctionExpression:
		return l.lowerFunctionExpression(expr)
	case *ir.RefExpression:
		return l.lowerRefExpression(expr)
	case *ir.DerefExpression:
		return l.lowerDerefExpression(expr)

	default:
		panic(fmt.Sprintf("TODO: lower %T", expr))
	}
}
