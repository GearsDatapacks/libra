package ast

type BaseStatment struct {}

func (stmt *BaseStatment) statementNode() {}

type ExpressionStatement struct {
	*BaseNode
	*BaseStatment
	Expression Expression
}

func (es *ExpressionStatement) Type() NodeType { return "ExpressionStatement" }

func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

type VariableDeclaration struct {
	*BaseNode
	*BaseStatment
	Constant bool
	Name string
	Value Expression
	DataType TypeExpression
}

func (varDec *VariableDeclaration) Type() NodeType { return "VariableDeclaration" }

func (varDec *VariableDeclaration) String() string {
	result := ""

	if varDec.Constant { result += "const" } else { result += "var" }
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
	*BaseNode
	*BaseStatment
	Name string
	Parameters []Parameter
	ReturnType TypeExpression
	Body []Statement
}

func (funcDec *FunctionDeclaration) Type() NodeType { return "FunctionDeclaration" }

func (funcDec *FunctionDeclaration) String() string {
	result := "function "

	result += funcDec.Name
	result += "("

	for i, parameter := range funcDec.Parameters {
		result += parameter.Name
		result += " "
		result += parameter.Type.String()

		if i != len(funcDec.Parameters) - 1 {
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
	*BaseNode
	*BaseStatment
	Value Expression
}

func (ret *ReturnStatement) Type() NodeType { return "ReturnStatement" }

func (ret *ReturnStatement) String() string {
	return "return " + ret.Value.String()
}

type IfElseStatement interface { ifElse() }

type IfStatement struct {
	*BaseNode
	*BaseStatment
	Condition Expression
	Body []Statement
	Else IfElseStatement
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
	*BaseNode
	*BaseStatment
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
