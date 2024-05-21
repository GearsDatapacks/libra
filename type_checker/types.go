package typechecker

import (
	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/ir"
	"github.com/gearsdatapacks/libra/type_checker/types"
	"github.com/gearsdatapacks/libra/type_checker/values"
)

/*
	case *ast.ArrayType:
		elemType := t.typeFromAst(ty.Type)

		if ty.Count == nil {
			return &types.ListType{
				ElemType: elemType,
			}
		}

		length := -1

		if _, ok := ty.Count.(*ast.InferredExpression); !ok {
			expr := convert(t.typeCheckExpression(ty.Count), types.Int, implicit)
			if expr == nil {
				t.Diagnostics.Report(diagnostics.CountMustBeInt(ty.Count.Location()))
			} else if expr.IsConst() {
				value := expr.ConstValue().(values.IntValue)
				length = int(value.Value)
			} else {
				t.Diagnostics.Report(diagnostics.NotConst(ty.Count.Location()))
			}
		}

		return &types.ArrayType{
			ElemType: elemType,
			Length:   length,
			CanInfer: false,
		}
	case *ast.MapType:
		keyType := t.typeFromAst(ty.KeyType)
		valueType := t.typeFromAst(ty.ValueType)
		return &types.MapType{
			KeyType:   keyType,
			ValueType: valueType,
		}
	case *ast.TupleType:
		dataTypes := []types.Type{}
		for _, ty := range ty.Types {
			dataTypes = append(dataTypes, t.typeFromAst(ty))
		}
		return &types.TupleType{
			Types: dataTypes,
		}
	default:
		panic(fmt.Sprintf("TODO: Types from %T", ty))
	}
}*/

func (t *typeChecker) typeCheckType(expression ast.Expression) types.Type {
	return t.typeFromExpr(t.typeCheckExpression(expression), expression.Location())
}

func (t *typeChecker) typeFromExpr(expr ir.Expression, location text.Location) types.Type {
	if expr.Type() == types.Invalid {
		return types.Invalid
	}
	if expr.Type() != types.RuntimeType {
		t.Diagnostics.Report(diagnostics.ExpressionNotType(location, expr.Type()))
		return types.Invalid
	}
	if !expr.IsConst() {
		t.Diagnostics.Report(diagnostics.NotConst(location))
		return types.Invalid
	}
	return expr.ConstValue().(values.TypeValue).Type.(types.Type)
}

func (t *typeChecker) lookupType(tok token.Token) types.Type {
	symbol := t.symbols.Lookup(tok.Value)
	if symbol == nil {
		t.Diagnostics.Report(diagnostics.UndefinedType(tok.Location, tok.Value))
		return types.Invalid
	}
	if symbol.GetType() != types.RuntimeType {
		t.Diagnostics.Report(diagnostics.ExpressionNotType(tok.Location, symbol.GetType()))
		return types.Invalid
	}
	if symbol.Value() == nil {
		t.Diagnostics.Report(diagnostics.NotConst(tok.Location))
		return types.Invalid
	}
	return symbol.Value().(values.TypeValue).Type.(types.Type)
}
