package diagnostics

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gearsdatapacks/libra/colour"
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

func (p *Partial) Location(location text.Location) *Diagnostic {
	return new(p.Kind, p.Message, location)
}

type Diagnostic struct {
	Kind     DiagnosticKind
	Message  string
	Location text.Location
}

func new(kind DiagnosticKind, message string, location text.Location) *Diagnostic {
	return &Diagnostic{
		Kind:     kind,
		Message:  message,
		Location: location,
	}
}

func (d *Diagnostic) Print() {
	d.WriteTo(os.Stdout, true)
}

func (d *Diagnostic) WriteTo(to io.Writer, printColour bool) {
	colour.UseColour = printColour
	colour.Writer = to

	diagColour := colour.Reset
	switch d.Kind {
	case Error:
		diagColour = colour.Error
	case Warning:
		diagColour = colour.Warning
	case Info:
		diagColour = colour.Info
	}

	fileName := d.Location.File.FileName
	span := d.Location.Span.ToLineSpan(d.Location.File)
	lines := d.Location.File.Lines

	colour.SetColour(colour.White)
	fmt.Fprintf(to, "%s:%d:%d:\n", fileName, span.StartLine+1, span.StartColumn+1)
	colour.ResetColour()

	spanLines := lines[span.StartLine : span.EndLine+1]
	numLines := len(spanLines)

	fmt.Fprint(to, spanLines[0].Text[:span.StartColumn])

	colour.SetColour(diagColour)
	if numLines == 1 {
		fmt.Fprint(to, spanLines[0].Text[span.StartColumn:span.EndColumn])
	} else {
		for i, line := range spanLines {
			line := line.Text
			if i == 0 {
				fmt.Fprintln(to, line[span.StartColumn:])
			} else if i == numLines-1 {
				fmt.Fprint(to, line[:span.EndColumn])
			} else {
				fmt.Fprintln(to, line)
			}
		}
	}

	colour.ResetColour()

	fmt.Fprintln(to, spanLines[numLines-1].Text[span.EndColumn:])

	column := span.StartColumn
	if numLines > 1 {
		column = 0
	}
	arrow := strings.Repeat(" ", column) + "^"
	fmt.Fprint(to, arrow+" ")

	colour.SetColour(diagColour)
	fmt.Fprintln(to, d.Message)
	fmt.Fprintln(to)

	colour.ResetColour()
}
