package errors

import (
	"log"

	"github.com/gearsdatapacks/libra/parser/ast"
)

func DevError(message string, errorNodes ...ast.Node) {
	RuntimeError(message+"\nIf you're seeing this, some feature has not been implemented properly", errorNodes...)
}

func RuntimeError(message string, errorNodes ...ast.Node) {
	if len(errorNodes) == 0 {
		log.Fatalf("RuntimeError: %s", message)
	}
	errorNode := errorNodes[0]
	log.Fatalf("RuntimeError at line %d, column %d: %s", errorNode.GetToken().Line, errorNode.GetToken().Column, message)
}

func TypeError(message string, errorNodes ...ast.Node) {
	if len(errorNodes) == 0 {
		log.Fatalf("TypeError: %s", message)
	}
	errorNode := errorNodes[0]
	log.Fatalf("TypeError at line %d, column %d: %s", errorNode.GetToken().Line, errorNode.GetToken().Column, message)
}
