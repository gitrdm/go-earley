package scanner_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexrule"
	"github.com/patrickhuber/go-earley/parser"
	"github.com/patrickhuber/go-earley/scanner"
	"github.com/patrickhuber/go-earley/terminal"
	"github.com/patrickhuber/go-earley/token"
	"github.com/stretchr/testify/require"
)

func TestScanner(t *testing.T) {
	t.Run("reads input", func(t *testing.T) {
		scanner := NewScanner(" ")
		result, err := scanner.Read()
		require.Nil(t, err)
		require.True(t, result)
		require.True(t, scanner.EndOfStream())
	})
	t.Run("updates position", func(t *testing.T) {
		scanner := NewScanner(" ")
		require.Equal(t, -1, scanner.Position())
		result, err := scanner.Read()
		require.Nil(t, err)
		require.True(t, result)
		require.Equal(t, 0, scanner.Position())
	})
	t.Run("resets column", func(t *testing.T) {
		scanner := NewScanner("test\nfile")
		for {
			result, err := scanner.Read()
			require.Nil(t, err)
			if !result {
				break
			}
			if scanner.Position() < 4 {
				require.Equal(t, scanner.Position()+1, scanner.Column())
			} else {
				require.Equal(t, scanner.Position()-4, scanner.Column())
			}
		}
	})
	t.Run("reads words", func(t *testing.T) {
		S := grammar.NewNonTerminal("S")
		word := lexrule.NewString("word")
		ws := terminal.NewWhitespace()
		// S -> ws
		// S -> word
		g := grammar.New(
			S,
			// S -> S S
			grammar.NewProduction(S, S, S),
			// S -> S
			grammar.NewProduction(S, S),
			// S -> word
			grammar.NewProduction(S, word),
			// S -> ws
			grammar.NewProduction(S, ws),
		)
		p := parser.New(g)
		s := scanner.New(p, strings.NewReader("word word word word"))
		ok, err := scanner.RunToEnd(s)
		require.Nil(t, err)
		require.True(t, ok)
	})
}

func NewScanner(text string) scanner.Scanner {
	reader := strings.NewReader(text)
	parser := NewFakeParser(len(text))
	return scanner.New(parser, reader)
}

type FakeParser struct {
	pulseCount   int
	currentCount int
}

func NewFakeParser(pulseCount int) parser.Parser {
	return &FakeParser{
		pulseCount:   pulseCount,
		currentCount: 0,
	}
}

func (p *FakeParser) Pulse(tokens ...token.Token) bool {
	if p.currentCount >= p.pulseCount {
		return false
	}
	p.currentCount++
	return true
}

func (p *FakeParser) Accepted() bool {
	return p.currentCount >= p.pulseCount
}

func (p *FakeParser) Location() int {
	return p.currentCount
}

func (p *FakeParser) Expected() []grammar.LexerRule {
	return []grammar.LexerRule{
		lexrule.NewTerminal(
			terminal.NewAny(),
		),
	}
}
