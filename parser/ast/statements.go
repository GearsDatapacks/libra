package ast

import "strings"

type BaseStatement struct {
	Exported bool
}

func (stmt *BaseStatement) statementNode() {}
func (stmt *BaseStatement) MarkExport() {
	stmt.Exported = true
}

func (stmt *BaseStatement) IsExport() bool {
	return stmt.Exported
}

type Exportable interface {
	export()
}

type canExport struct{}

func (canExport) export() {}

type ExpressionStatement struct {
	BaseNode
	BaseStatement
	Expression Expression
}

func (es *ExpressionStatement) Type() NodeType { return "ExpressionStatement" }

func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type VariableDeclaration struct {
	BaseNode
	BaseStatement
	Constant bool
	Name     string
	Value    Expression
	DataType TypeExpression
}

func (varDec *VariableDeclaration) Type() NodeType { return "VariableDeclaration" }

func (varDec *VariableDeclaration) String() string {
	result := ""

	if varDec.Constant {
		result += "const"
	} else {
		result += "var"
	}
	result += " "
	result += varDec.Name
	result += " = "
	result += varDec.Value.String()

	return result
}

type Parameter struct {
	Name string
	Type TypeExpression
}

type FunctionDeclaration struct {
	BaseNode
	BaseStatement
	canExport
	Name       string
	MethodOf   TypeExpression
	Parameters []Parameter
	ReturnType TypeExpression
	Body       []Statement
}

func (funcDec *FunctionDeclaration) Type() NodeType { return "FunctionDeclaration" }

func (funcDec *FunctionDeclaration) String() string {
	result := "fn "

	result += funcDec.Name
	result += "("

	for i, parameter := range funcDec.Parameters {
		result += parameter.Name
		result += " "
		result += parameter.Type.String()

		if i != len(funcDec.Parameters)-1 {
			result += ", "
		}
	}

	result += ") {\n"

	for _, statement := range funcDec.Body {
		result += "  "
		result += statement.String()
		result += "\n"
	}

	result += "}"

	return result
}

type ReturnStatement struct {
	BaseNode
	BaseStatement
	Value Expression
}

func (ret *ReturnStatement) Type() NodeType { return "ReturnStatement" }

func (ret *ReturnStatement) String() string {
	return "return " + ret.Value.String()
}

type IfElseStatement interface{ ifElse() }

type IfStatement struct {
	BaseNode
	BaseStatement
	Condition Expression
	Body      []Statement
	Else      IfElseStatement
}

func (ifs *IfStatement) Type() NodeType { return "IfStatement" }

func (ifs *IfStatement) String() string {
	result := "if "
	result += ifs.Condition.String()
	result += " {\n"

	for _, statement := range ifs.Body {
		result += "  "
		result += statement.String()
		result += "\n"
	}

	result += "}"

	return result
}

func (ifs *IfStatement) ifElse() {}

type ElseStatement struct {
	BaseNode
	BaseStatement
	Body []Statement
}

func (elses *ElseStatement) Type() NodeType { return "IfStatement" }

func (elses *ElseStatement) String() string {
	result := "else {\n"

	for _, statement := range elses.Body {
		result += "  "
		result += statement.String()
		result += "\n"
	}

	result += "}"

	return result
}

func (elses *ElseStatement) ifElse() {}

type WhileLoop struct {
	BaseNode
	BaseStatement
	Condition Expression
	Body      []Statement
}

func (while *WhileLoop) Type() NodeType { return "WhileLoop" }

func (while *WhileLoop) String() string {
	result := "while "
	result += while.Condition.String()
	result += " {\n"

	for _, statement := range while.Body {
		result += "  "
		result += statement.String()
		result += "\n"
	}

	result += "}"

	return result
}

type ForLoop struct {
	BaseNode
	BaseStatement
	Initial   Statement
	Condition Expression
	Update    Statement
	Body      []Statement
}

func (forLoop *ForLoop) Type() NodeType { return "ForLoop" }

func (forLoop *ForLoop) String() string {
	result := "for "
	result += forLoop.Initial.String()
	result += "; "
	result += forLoop.Condition.String()
	result += "; "
	result += forLoop.Update.String()
	result += " {\n"

	for _, statement := range forLoop.Body {
		result += "  "
		result += statement.String()
		result += "\n"
	}

	result += "}"

	return result
}

type StructDeclaration struct {
	BaseNode
	BaseStatement
	canExport
	Name    string
	Members map[string]TypeExpression
}

func (structDec *StructDeclaration) Type() NodeType { return "StructDeclaration" }

func (structDec *StructDeclaration) String() string {
	result := "struct "

	result += structDec.Name
	result += " {\n"

	for name, dataType := range structDec.Members {
		result += name
		result += " "
		result += dataType.String()
		result += "\n"
	}

	result += "}"

	return result
}

type TupleStructDeclaration struct {
	BaseNode
	BaseStatement
	canExport
	Name    string
	Members []TypeExpression
}

func (structDec *TupleStructDeclaration) Type() NodeType { return "TupleStructDeclaration" }

func (structDec *TupleStructDeclaration) String() string {
	result := "struct "

	result += structDec.Name
	result += "("

	for i, dataType := range structDec.Members {
		if i != 0 {
			result += ", "
		}
		result += dataType.String()
	}

	result += ")"

	return result
}

type InterfaceMember struct {
	Name       string
	IsFunction bool
	Parameters []TypeExpression
	ResultType TypeExpression
}

type InterfaceDeclaration struct {
	BaseNode
	BaseStatement
	canExport
	Name    string
	Members []InterfaceMember
}

func (intDecl *InterfaceDeclaration) Type() NodeType { return "InterfaceDeclaration" }

func (intDecl *InterfaceDeclaration) String() string {
	return "interface {}"
}

type TypeDeclaration struct {
	BaseNode
	BaseStatement
	canExport
	Name     string
	DataType TypeExpression
}

func (typeDecl *TypeDeclaration) Type() NodeType { return "TypeDeclaration" }

func (typeDecl *TypeDeclaration) String() string {
	return "type " + typeDecl.Name + " = " + typeDecl.DataType.String()
}

type ImportStatement struct {
	BaseNode
	BaseStatement
	Module string
	Alias  string
	ImportAll bool
	ImportedSymbols []string
}

func (*ImportStatement) Type() NodeType { return "ImportStatement" }

func (imp *ImportStatement) String() string {
	if imp.ImportAll {
		return "import * from \"" + imp.Module + "\""
	}
	if imp.Alias != "" {
		return "import \"" + imp.Module + "\" as " + imp.Alias
	}
	if imp.ImportedSymbols != nil {
		return "import {" + strings.Join(imp.ImportedSymbols, ", ") + "} from \"" + imp.Module + "\""
	}
	return "import \"" + imp.Module + "\""
}
