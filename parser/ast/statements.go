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
	DataType string
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

type FunctionDeclaration struct {
	*BaseNode
	*BaseStatment
	Name string
	Parameters [][2]string // [name, type]
	ReturnType string
	Body []Statement
}

func (funcDec *FunctionDeclaration) Type() NodeType { return "FunctionDeclaration" }

func (funcDec *FunctionDeclaration) String() string {
	result := "function "

	result += funcDec.Name
	result += "("

	for i, parameter := range funcDec.Parameters {
		result += parameter[0]
		result += " "
		result += parameter[1]

		if i != len(funcDec.Parameters) - 1 {
			result += ", "
		}
	}

	result += ") {\n"

	for _, statement := range funcDec.Body {
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
