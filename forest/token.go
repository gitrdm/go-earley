package forest

import (
	"fmt"

	"github.com/patrickhuber/go-earley/token"
)

type Token struct {
	Token    token.Token
	origin   int
	location int
}

func NewToken(tok token.Token, origin int, location int) *Token {
	return &Token{
		origin:   origin,
		location: location,
		Token:    tok,
	}
}

func (Token) node() {}

func (t Token) Origin() int { return t.origin }

func (t Token) Location() int { return t.location }

func (t Token) String() string {
	return fmt.Sprintf("(%s, %d, %d)", t.Token.Type(), t.origin, t.location)
}

func (t *Token) Accept(v Visitor) {
	v.VisitToken(t)
}
