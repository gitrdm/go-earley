package terminal

import "github.com/patrickhuber/go-earley/grammar"

type Character struct {
	grammar.SymbolImpl
	Value rune
}

func NewCharacter(ch rune) *Character {
	return &Character{Value: ch}
}

// IsMatch implements grammar.Terminal.
func (c *Character) IsMatch(ch rune) bool {
	return c.Value == ch
}

func (c Character) String() string {
	return string(c.Value)
}
