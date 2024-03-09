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

type Span struct {
	Line, Column, End int
}

func NewSpan(line, col, end int) Span {
	return Span{
		Line:   line,
		Column: col,
		End:    end,
	}
}
