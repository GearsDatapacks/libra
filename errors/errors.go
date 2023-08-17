package errors

import (
	"log"

	"github.com/gearsdatapacks/libra/parser/ast"
)

func DevError(message string, errorNodes ...ast.Node) {
	logError("Error", message+"\nIf you're seeing this, some feature has not been implemented properly", errorNodes...)
}

func logError(prefix, message string, errorNodes ...ast.Node) {
	if len(errorNodes) == 0 {
		log.Fatalf("%s: %s", prefix, message)
	}
	errorNode := errorNodes[0]
	log.Fatalf("%s at line %d, column %d: %s", prefix, errorNode.GetToken().Line, errorNode.GetToken().Column, message)
}

func RuntimeError(message string, errorNodes ...ast.Node) {
	logError("RuntimeError", message, errorNodes...)
}

func TypeError(message string, errorNodes ...ast.Node) {
	logError("TypeError", message, errorNodes...)
}
