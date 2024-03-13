package text

import (
	"os"
	"strings"
)

type SourceFile struct {
	FileName string
	Text     string
	Lines    []string
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
	return &SourceFile{
		FileName: fileName,
		Text:     text,
		Lines:    strings.Split(text, "\n"),
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
	StartLine, EndLine, Column, End int
}

func NewSpan(startLine, endLine, col, end int) Span {
	return Span{
		StartLine: startLine,
		EndLine:   endLine,
		Column:    col,
		End:       end,
	}
}

func (s Span) To(other Span) Span {
	return NewSpan(s.StartLine, other.EndLine, s.Column, other.End)
}
