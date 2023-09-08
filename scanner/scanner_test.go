package scanner_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/forest"
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
		scanner := NewScanner(" ", NewFakeParser(lexrule.NewTerminal(terminal.NewWhitespace())))
		result, err := scanner.Read()
		require.NoError(t, err)
		require.True(t, result)
		require.True(t, scanner.EndOfStream())
	})
	t.Run("updates position", func(t *testing.T) {
		scanner := NewScanner(" ", NewFakeParser(lexrule.NewTerminal(terminal.NewWhitespace())))
		require.Equal(t, -1, scanner.Position())
		result, err := scanner.Read()
		require.NoError(t, err)
		require.True(t, result)
		require.Equal(t, 0, scanner.Position())
	})
	t.Run("resets column", func(t *testing.T) {
		parser := NewFakeParser(
			lexrule.NewString("test"),
			lexrule.NewString("\n"),
			lexrule.NewString("file"),
		)
		scanner := NewScanner("test\nfile", parser)
		for {
			result, err := scanner.Read()
			require.NoError(t, err)
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
}

func NewScanner(text string, parser parser.Parser) scanner.Scanner {
	return scanner.New(parser, text)
}

type FakeParser struct {
	rules []grammar.LexerRule
	index int
}

func NewFakeParser(rules ...grammar.LexerRule) parser.Parser {
	return &FakeParser{
		rules: rules,
		index: 0,
	}
}

func (p *FakeParser) Pulse(tokens token.Token) (bool, error) {
	if p.index >= len(p.rules) {
		return false, nil
	}
	p.index++
	return true, nil
}

func (p *FakeParser) Accepted() bool {
	return p.index >= len(p.rules)
}

func (p *FakeParser) Location() int {
	return p.index
}

func (p *FakeParser) Expected() []grammar.LexerRule {
	if p.index >= len(p.rules) {
		return nil
	}
	return p.rules[p.index : p.index+1]
}

// GetForestRoot implements parser.Parser.
func (*FakeParser) GetForestRoot() (forest.Node, bool) {
	return nil, false
}
