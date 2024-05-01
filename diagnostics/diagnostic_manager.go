package diagnostics

import (
	"fmt"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/text"
	"github.com/gearsdatapacks/libra/type_checker/types"
)

type Manager []Diagnostic

func (m *Manager) reportError(msg string, location text.Location) {
	*m = append(*m, new(Error, msg, location))
}

func (m *Manager) reportInfo(msg string, location text.Location) {
	*m = append(*m, new(Info, msg, location))
}

// Lexer Diagnostics

func (m *Manager) ReportInvalidCharacter(location text.Location, char byte) {
	msg := fmt.Sprintf("Invalid character: %q", char)
	m.reportError(msg, location)
}

func (m *Manager) ReportUnterminatedString(location text.Location) {
	msg := "Unterminated string"
	m.reportError(msg, location)
}

func (m *Manager) ReportInvalidEscapeSequence(location text.Location, char byte) {
	msg := fmt.Sprintf("Invalid escape sequence: '\\%c'", char)
	m.reportError(msg, location)
}

func (m *Manager) ReportNumbersCannotEndWithSeparator(location text.Location) {
	msg := "Numbers cannot end with numeric separators"
	m.reportError(msg, location)
}

// Parser Diagnostics

func (m *Manager) ReportExpectedExpression(location text.Location, kind token.Kind) {
	msg := fmt.Sprintf("Expected expression, found %s", kind.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportExpectedNewline(location text.Location, kind token.Kind) {
	msg := fmt.Sprintf("Expected newline after statement, found %s", kind.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportExpectedToken(location text.Location, expected token.Kind, actual token.Kind) {
	msg := fmt.Sprintf("Expected %s, found %s", expected.String(), actual.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportElseStatementWithoutIf(location text.Location) {
	msg := "Else statement not allowed without preceding if"
	m.reportError(msg, location)
}

func (m *Manager) ReportExpectedKeyword(location text.Location, keyword string, foundToken token.Token) {
	tokenValue := foundToken.Kind.String()
	if foundToken.Kind == token.IDENTIFIER {
		tokenValue = foundToken.Value
	}

	msg := fmt.Sprintf("Expected %q keyword, found %s", keyword, tokenValue)
	m.reportError(msg, location)
}

func (m *Manager) ReportKeywordOverwritten(location text.Location, keyword string, declared text.Location) {
	errMsg := fmt.Sprintf(
		"Expected %q keyword, but it has been overwritten by a variable",
		keyword)
	info := "Try removing or renaming this variable"

	m.reportError(errMsg, location)
	m.reportInfo(info, declared)
}

func (m *Manager) ReportLastParameterMustHaveType(location text.Location, fnLocation text.Location) {
	msg := "The last parameter of a function must have a type annotation"
	m.reportError(msg, location)

	if location.Span.StartLine != fnLocation.Span.StartLine {
		info := "Parameter of this function"
		m.reportInfo(info, fnLocation)
	}
}

func (m *Manager) ReportLastStructFieldMustHaveType(location text.Location, structLoc text.Location) {
	errMsg := "The last field of a struct must have a type annotation"
	m.reportError(errMsg, location)

	if location.Span.StartLine != structLoc.Span.StartLine {
		info := "Field in this struct"
		m.reportInfo(info, structLoc)
	}
}

func (m *Manager) ReportMemberAndMethodNotAllowed(location text.Location) {
	msg := "Functions cannot be both methods and static members"

	m.reportError(msg, location)
}

func (m *Manager) ReportExpectedMemberOrStructBody(location text.Location, tok token.Token) {
	msg := fmt.Sprintf("Invalid right-hand side of expression. Expected identifier or struct body, found %s", tok.Kind.String())

	m.reportError(msg, location)
}

func (m *Manager) ReportOneImportModifierAllowed(location text.Location) {
	msg := "Only one import modifier is allowed"

	m.reportError(msg, location)
}

func (m *Manager) ReportOnlyTopLevelStatement(location text.Location, stmtKind string) {
	msg := fmt.Sprintf("%s not allowed here", stmtKind)

	m.reportError(msg, location)
}

func (m *Manager) ReportExpectedType(location text.Location, kind token.Kind) {
	msg := fmt.Sprintf("Expected type, found %s", kind.String())
	m.reportError(msg, location)
}

// Type-checker Diagnostics

func (m *Manager) ReportUndefinedType(location text.Location, name string) {
	msg := fmt.Sprintf("Type %q is not defined", name)
	m.reportError(msg, location)
}

func (m *Manager) ReportNotAssignable(location text.Location, expected, actual types.Type) {
	msg := fmt.Sprintf("Value of type %q is not assignable to type %q", actual.String(), expected.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportVariableDefined(location text.Location, name string) {
	msg := fmt.Sprintf("Variable %q is already defined", name)
	m.reportError(msg, location)
}

func (m *Manager) ReportVariableUndefined(location text.Location, name string) {
	msg := fmt.Sprintf("Variable %q is not defined", name)
	m.reportError(msg, location)
}

func (m *Manager) ReportBinaryOperatorUndefined(location text.Location, operator string, left, right types.Type) {
	msg := fmt.Sprintf("Operator %q is not defined for types %q and %q", operator, left.String(), right.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportUnaryOperatorUndefined(location text.Location, operator string, operand types.Type) {
	msg := fmt.Sprintf("Operator %q is not defined for operand of type %q", operator, operand.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportCannotCast(location text.Location, from, to types.Type) {
	msg := fmt.Sprintf("Cannot cast value of type %q to type %q", from.String(), to.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportCannotIncDec(location text.Location, incDec string) {
	msg := fmt.Sprintf("Cannot %s a non-variable value", incDec)
	m.reportError(msg, location)
}

func (m *Manager) ReportValueImmutable(location text.Location) {
	msg := "Cannot modify value, it is immutable"
	m.reportError(msg, location)
}

func (m *Manager) ReportNotConst(location text.Location) {
	msg := "Value must be known at compile time"
	m.reportError(msg, location)
}

func (m *Manager) ReportCountMustBeInt(location text.Location) {
	msg := "Array length must be an integer"
	m.reportError(msg, location)
}

func (m *Manager) ReportCannotIndex(location text.Location, leftType, indexType types.Type) {
	msg := fmt.Sprintf("Cannot index value of type %q with value of type %q", leftType, indexType)
	m.reportError(msg, location)
}

func (m *Manager) ReportNotHashable(location text.Location, ty types.Type) {
	msg := fmt.Sprintf("Value of type %q cannot be used as a key in a map", ty.String())
	m.reportError(msg, location)
}

func (m *Manager) ReportCannotAssign(location text.Location) {
	msg := "Cannot assign to a non-variable value"
	m.reportError(msg, location)
}
