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

type Partial struct {
	Kind    DiagnosticKind
	Message string
}

func partial(kind DiagnosticKind, message string) *Partial {
	return &Partial{
		Kind:    kind,
		Message: message,
	}
}

func (p *Partial) Location(location text.Location) Diagnostic {
	return new(p.Kind, p.Message, location)
}

type Diagnostic struct {
	Kind     DiagnosticKind
	Message  string
	Location text.Location
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
	fmt.Println(span)

	SetColour(White)
	fmt.Printf("%s:%d:%d:\n", fileName, span.StartLine+1, span.Column+1)
	ResetColour()

	spanLines := lines[span.StartLine : span.EndLine+1]
	numLines := len(spanLines)

	fmt.Print(spanLines[0][:span.Column])

	SetColour(colour)
	if numLines == 1 {
		fmt.Print(spanLines[0][span.Column:span.End])
	} else {
		for i, line := range spanLines {
			if i == 0 {
				fmt.Println(line[span.Column:])
			} else if i == numLines-1 {
				fmt.Print(line[:span.End])
			} else {
				fmt.Println(line)
			}
		}
	}

	ResetColour()

	fmt.Println(spanLines[numLines-1][span.End:])

	column := span.Column
	if numLines > 1 {
		column = 0
	}
	arrow := strings.Repeat(" ", column) + "^"
	fmt.Print(arrow + " ")

	SetColour(colour)
	fmt.Println(d.Message)
	fmt.Println()

	ResetColour()
}
