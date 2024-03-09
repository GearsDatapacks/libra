package diagnostics

import (
	"fmt"

	"github.com/gearsdatapacks/libra/lexer/token"
	"github.com/gearsdatapacks/libra/text"
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

	if location.Span.Line != fnLocation.Span.Line {
		info := "Parameter of this function"
		m.reportInfo(info, fnLocation)
	}
}

func (m *Manager) ReportLastStructFieldMustHaveType(location text.Location, structLoc text.Location) {
	errMsg := "The last field of a struct must have a type annotation"
	m.reportError(errMsg, location)

	if location.Span.Line != structLoc.Span.Line {
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
