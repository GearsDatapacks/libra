package diagnostics

import (
	"fmt"
	"strings"

	"github.com/gearsdatapacks/libra/lexer/token"
)

type DiagnosticKind int

const (
	Error DiagnosticKind = iota
	Warning
	Info
)

type Diagnostic struct {
	Kind    DiagnosticKind
	Message string
	Span    token.Span
	file    string
  lines []string
}

func new(kind DiagnosticKind, message string, span token.Span, file string, lines []string) Diagnostic {
	return Diagnostic{
		Kind:    kind,
		Message: message,
		Span:    span,
		file:    file,
    lines: lines,
	}
}

func (d *Diagnostic) Print() {
  colour := Reset
	switch d.Kind {
	case Error:
		colour = Red
	case Warning:
		colour = Yellow
	case Info:
		colour = Cyan
	}

  SetColour(White)
  fmt.Printf("%s:%d:%d:\n", d.file, d.Span.Line+1, d.Span.Col+1)
  ResetColour()

  line := d.lines[d.Span.Line]

  fmt.Print(line[:d.Span.Col])

  SetColour(colour)
  fmt.Print(line[d.Span.Col:d.Span.End])
  ResetColour()

  fmt.Println(line[d.Span.End:])

  arrow := strings.Repeat(" ", d.Span.Col) + "^"
  fmt.Print(arrow + " ")

  SetColour(colour)
	fmt.Println(d.Message)
  fmt.Println()

  ResetColour()
}
