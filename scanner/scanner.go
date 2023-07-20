package scanner

import (
	"bufio"
	"io"

	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/parser"
	"github.com/patrickhuber/go-earley/token"
)

type Scanner interface {
	Read() (bool, error)
	Position() int
	Line() int
	Column() int
	EndOfStream() bool
	Parser() parser.Parser
}

type scanner struct {
	position    int
	line        int
	column      int
	parser      parser.Parser
	reader      *bufio.Reader
	tokens      []token.Token
	accumulator token.Accumulator
	registry    map[string]token.Factory
}

// New creates a new scanner from the given parser and io reader
func New(p parser.Parser, reader io.Reader) Scanner {
	return &scanner{
		parser:      p,
		reader:      bufio.NewReader(reader),
		position:    -1,
		line:        0,
		column:      0,
		accumulator: token.NewAccumulator(),
	}
}

// Column implements Scanner.
func (s *scanner) Column() int {
	return s.column
}

// EndOfStream implements Scanner.
func (s *scanner) EndOfStream() bool {
	_, err := s.reader.Peek(1)
	return err != nil
}

// Line implements Scanner.
func (s *scanner) Line() int {
	return s.line
}

// Position implements Scanner.
func (s *scanner) Position() int {
	return s.position
}

// Parser implements Scanner.Parser
func (s *scanner) Parser() parser.Parser {
	return s.parser
}

// Read implements Scanner.
func (s *scanner) Read() (bool, error) {
	if s.EndOfStream() {
		return false, nil
	}
	ch, err := s.readRune()
	if err != nil {
		return false, err
	}

	s.update(ch)

	if s.matchesExistingTokens() {
		if s.EndOfStream() {
			if !s.tryParseExistingTokens() {
				return false, nil
			}
			return true, nil
		}
	}

	if s.anyExistingTokens() {
		if !s.tryParseExistingTokens() {
			return false, nil
		}
	}

	matched, err := s.matchesNewTokens()
	if err != nil {
		return false, err
	}

	if matched {
		if s.EndOfStream() {
			if !s.tryParseExistingTokens() {
				return false, nil
			}
			return s.parser.Accepted(), nil
		}
		return true, nil
	}

	if !s.isEndOfLineCharacter(ch) {
		return true, nil
	}

	return true, nil
}

func (s *scanner) readRune() (rune, error) {
	ch, _, err := s.reader.ReadRune()
	var zero rune
	if err != nil {
		return zero, err
	}
	s.accumulator.Accumulate(ch)
	return ch, nil
}

func (s *scanner) update(ch rune) {
	if s.isEndOfLineCharacter(ch) {
		s.column = 0
		s.line++
	} else {
		s.column++
	}
	s.position++
}

func (s *scanner) matchesExistingTokens() bool {
	return false
}

func (s *scanner) tryParseExistingTokens() bool {
	size := len(s.tokens)
	anyTokens := size > 0
	if !anyTokens {
		return false
	}
	i := 0
	for i < size {
		tok := s.tokens[i]
		if tok.Accepted() {
			i++
			continue
		}
	}
	return false
}

func (s *scanner) anyExistingTokens() bool {
	return false
}

func (s *scanner) matchesNewTokens() (bool, error) {
	return s.matchLexerRules(s.parser.Expected(), s.tokens)
}

func (s *scanner) matchLexerRules(lexerRules []grammar.LexerRule, lexemes []token.Token) (bool, error) {
	anyMatches := false

	ch := s.accumulator.RuneAt(s.position)

	for _, lexerRule := range lexerRules {
		if !lexerRule.CanApply(ch) {
			continue
		}

		factory, ok := s.registry[lexerRule.Type()]
		if !ok {
			continue
		}
		span := s.accumulator.Span(s.position, 0)
		tok, err := factory.Create(lexerRule, span)
		if err != nil {
			return false, err
		}
		if !tok.Scan() {
			err = factory.Free(tok)
			if err != nil {
				return false, err
			}
			continue
		}
	}
	return anyMatches, nil
}

func (s *scanner) isEndOfLineCharacter(ch rune) bool {
	return ch == '\n'
}
