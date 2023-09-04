package errors

import (
	"fmt"
	"os"

	"github.com/gearsdatapacks/libra/parser/ast"
)

type LanguageError struct {
	ErrorType string
	Message string
	Line int
	Column int
}

func (err LanguageError) Error() string {
	if err.Line == -1 || err.Column == -1 {
		return fmt.Sprintf("%s: %s", err.ErrorType, err.Message)
	}
	return fmt.Sprintf("%s at line %d, column %d: %s", err.ErrorType, err.Line, err.Column, err.Message)
}

func DevError(message string, errorNodes ...ast.Node) error {
	return makeError("Error", message+"\nIf you're seeing this, some feature has not been implemented properly", errorNodes...)
}

func makeError(prefix, message string, errorNodes ...ast.Node) error {
	if len(errorNodes) == 0 {
		return LanguageError{
			Line: -1,
			Column: -1,
			Message: message,
			ErrorType: prefix,
		}
	}

	errorNode := errorNodes[0]
	return LanguageError{
		Line: errorNode.GetToken().Line,
		Column: errorNode.GetToken().Column,
		Message: message,
		ErrorType: prefix,
	}
}

func RuntimeError(message string, errorNodes ...ast.Node) error {
	return makeError("RuntimeError", message, errorNodes...)
}

func LogError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
