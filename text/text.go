package text

import (
	"os"
	"strings"
)

type Line struct {
	Text string
	Span Span
}

type SourceFile struct {
	FileName string
	Text     string
	Lines    []Line
}

func LoadFile(fileName string) *SourceFile {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	text := string(bytes)
	return NewFile(fileName, text)
}

func NewFile(fileName string, text string) *SourceFile {
	lines := []Line{}

	position := 0
	for _, line := range strings.Split(text, "\n") {
		start := position
		position += len(line)
		lines = append(lines, Line{
			Text: line,
			Span: NewSpan(start, position),
		})
	}

	return &SourceFile{
		FileName: fileName,
		Text:     text,
		Lines:    lines,
	}
}

type Location struct {
	Span Span
	File *SourceFile
}

func (l Location) To(other Location) Location {
	if l.File != other.File {
		panic("Must join locations from the same file")
	}
	return Location{
		Span: l.Span.To(other.Span),
		File: l.File,
	}
}

type Span struct {
	Start, End int
}

func NewSpan(start, end int) Span {
	return Span{
		Start: start,
		End:   end,
	}
}

func (s Span) To(other Span) Span {
	return NewSpan(s.Start, other.End)
}

func (s Span) ToLineSpan(file *SourceFile) LineSpan {
	var startLine, endLine, startColumn, endColumn int

	for i, line := range file.Lines {
		if s.Start > line.Span.Start && s.Start < line.Span.End {
			startLine = i
			startColumn = s.Start - line.Span.Start
		}
		if s.End > line.Span.Start && s.End < line.Span.End {
			endLine = i
			endColumn = s.End - line.Span.Start
		}
	}

	return LineSpan{
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
	}
}

type LineSpan struct {
	StartLine, EndLine, StartColumn, EndColumn int
}
