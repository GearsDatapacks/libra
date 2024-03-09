package diagnostics

import (
	"fmt"
	"strings"
	"github.com/gearsdatapacks/libra/text"
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
	Location    text.Location
}

func new(kind DiagnosticKind, message string, location text.Location) Diagnostic {
	return Diagnostic{
		Kind:     kind,
		Message:  message,
		Location: location,
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

	fileName := d.Location.File.FileName
	span := d.Location.Span
	lines := d.Location.File.Lines

	SetColour(White)
	fmt.Printf("%s:%d:%d:\n", fileName, span.Line+1, span.Column+1)
	ResetColour()

	line := lines[span.Line]

	fmt.Print(line[:span.Column])

	SetColour(colour)
	fmt.Print(line[span.Column:span.End])
	ResetColour()

	fmt.Println(line[span.End:])

	arrow := strings.Repeat(" ", span.Column) + "^"
	fmt.Print(arrow + " ")

	SetColour(colour)
	fmt.Println(d.Message)
	fmt.Println()

	ResetColour()
}
