package parser

import "github.com/gearsdatapacks/libra/lexer/token"

func parseDelemitedList[Elem any](p *parser, delim token.Kind, elemFn func() Elem) ( result []Elem, delimToken token.Token) {
	result = []Elem{}

	for !p.eof() && p.next().Kind != delim {
		result = append(result, elemFn())

		if p.next().Kind == token.COMMA{
			p.consume()
		} else {
			break
		}
	}

	delimToken = p.expect(delim)

	return result, delimToken
}
