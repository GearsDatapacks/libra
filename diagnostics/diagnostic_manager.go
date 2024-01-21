package diagnostics

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type Manager struct {
  Diagnostics []Diagnostic
  file string
  lines []string
}

func New(file string, src string) Manager {
  return Manager{
    Diagnostics: []Diagnostic{},
    file: file,
    lines: strings.Split(src, "\n"),
  }
}

func (m *Manager) reportError(msg string, span token.Span) {
  m.Diagnostics = append(m.Diagnostics, new(Error, msg, span, m.file, m.lines))
}

func (m *Manager) ReportInvalidCharacter(span token.Span, char byte) {
  msg := fmt.Sprintf("Invalid character: %q", char)
  m.reportError(msg, span)
}

func (m *Manager) ReportUnterminatedString(span token.Span) {
  msg := "Unterminated string"
  m.reportError(msg, span)
}

func (m *Manager) ReportInvalidEscapeSequence(span token.Span, char byte) {
  msg := fmt.Sprintf("Invalid escape sequence: '\\%c'", char)
  m.reportError(msg, span)
}

func (m *Manager) ReportNumbersCannotEndWithSeparator(span token.Span) {
  msg := "Numbers cannot end with numeric separators"
  m.reportError(msg, span)
}

