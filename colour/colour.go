package colour

import (
	"fmt"
	"io"
)

type Colour string

// Ansi colours
const (
	Reset  Colour = "\033[0m"
	Red    Colour = "\033[31m"
	Green  Colour = "\033[32m"
	Yellow Colour = "\033[33m"
	Blue   Colour = "\033[34m"
	Purple Colour = "\033[35m"
	Cyan   Colour = "\033[36m"
	Gray   Colour = "\033[37m"
	White  Colour = "\033[97m"
)

// Special colours
const (
	// Diagnostic colours
	Error = Red
	Warning = Yellow
	Info = Cyan

	// AST print colours
	NodeName = Blue
	Literal = Green
	Name = Purple
	Symbol = White
	Attribute = Yellow
	Location = Cyan
)

var UseColour = true
var Writer io.Writer

func SetColour(colour Colour) {
	if UseColour {
		fmt.Fprint(Writer, colour)
	}
}

func ResetColour() {
	if UseColour {
		fmt.Print(Reset)
	}
}
