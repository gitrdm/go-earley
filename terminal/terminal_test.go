package terminal_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/terminal"
)

func TestTerminal(t *testing.T) {
	t.Run("whitespace", func(t *testing.T) {
		term := terminal.NewWhitespace()
		whitespaces := []rune{' ', '\f', '\t', '\r', '\n'}
		if _, ok := matches(term, whitespaces); !ok {
			t.Fatalf("expected all whitespace to match")
		}
	})
	t.Run("number", func(t *testing.T) {
		term := terminal.NewNumber()
		numbers := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
		if _, ok := matches(term, numbers); !ok {
			t.Fatalf("expected all numbers to match")
		}
	})
	t.Run("letter", func(t *testing.T) {
		term := terminal.NewLetter()
		var letters []rune
		for letter := 'a'; letter <= 'Z'; letter++ {
			letters = append(letters, letter)
		}
		if _, ok := matches(term, letters); !ok {
			t.Fatalf("expected all letters to match")
		}
	})
}

func matches(term grammar.Terminal, runes []rune) (rune, bool) {
	for _, r := range runes {
		if !term.IsMatch(r) {
			return r, false
		}
	}
	var zero rune
	return zero, true
}
