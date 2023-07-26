package scanner

import (
	"bufio"
	"io"

	"github.com/patrickhuber/go-earley/capture"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexeme"
	"github.com/patrickhuber/go-earley/parser"
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
	position int
	line     int
	column   int
	parser   parser.Parser
	reader   *bufio.Reader
	lexemes  []lexeme.Lexeme
	capture  *capture.Builder
	registry map[string]lexeme.Factory
}

// New creates a new scanner from the given parser and io reader
func New(p parser.Parser, reader io.Reader) Scanner {
	return &scanner{
		parser:   p,
		reader:   bufio.NewReader(reader),
		position: -1,
		line:     0,
		column:   0,
		capture:  capture.FromSlice(),
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

	if s.matchesExistingLexemes() {
		if s.EndOfStream() {
			if !s.tryParseExistingLexemes() {
				return false, nil
			}
			return true, nil
		}
	}

	if s.anyExistingLexemes() {
		if !s.tryParseExistingLexemes() {
			return false, nil
		}
	}

	matched, err := s.matchesNewLexemes()
	if err != nil {
		return false, err
	}

	if matched {
		if s.EndOfStream() {
			if !s.tryParseExistingLexemes() {
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
	s.capture.Append(ch)
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

func (s *scanner) matchesExistingLexemes() bool {
	return false
}

func (s *scanner) tryParseExistingLexemes() bool {
	size := len(s.lexemes)
	anyLexemes := size > 0
	if !anyLexemes {
		return false
	}
	i := 0
	for i < size {
		tok := s.lexemes[i]
		if tok.Accepted() {
			i++
			continue
		}
	}
	return false
}

func (s *scanner) anyExistingLexemes() bool {
	return false
}

func (s *scanner) matchesNewLexemes() (bool, error) {
	return s.matchLexerRules(s.parser.Expected(), s.lexemes)
}

func (s *scanner) matchLexerRules(lexerRules []grammar.LexerRule, lexemes []lexeme.Lexeme) (bool, error) {
	anyMatches := false

	ch := s.capture.RuneAt(s.position)

	for _, lexerRule := range lexerRules {
		if !lexerRule.CanApply(ch) {
			continue
		}

		factory, ok := s.registry[lexerRule.Type()]
		if !ok {
			continue
		}
		span := capture.SpanLength(s.capture, s.capture.Offset(), 0)
		tok, err := factory.Create(lexerRule, span, s.position)
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
