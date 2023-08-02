package scanner

import (
	"strings"

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
	lexemes  []lexeme.Lexeme
	input    string
	reader   *strings.Reader
	registry map[string]lexeme.Factory
}

// New creates a new scanner from the given parser and io reader
func New(p parser.Parser, input string) Scanner {
	return &scanner{
		parser:   p,
		position: -1,
		line:     0,
		column:   0,
		input:    input,
		reader:   strings.NewReader(input),
	}
}

// Column implements Scanner.
func (s *scanner) Column() int {
	return s.column
}

// EndOfStream implements Scanner.
func (s *scanner) EndOfStream() bool {
	// if unread runes are zero, we are at the end
	return s.reader.Len() == 0
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
	ch, err := s.read()
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

	matched, err := s.matchesNewLexemes(ch)
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

func (s *scanner) read() (rune, error) {
	ch, n, err := s.reader.ReadRune()
	if err != nil {
		var zero rune
		return zero, err
	}
	s.position += n
	return ch, nil
}

func (s *scanner) update(ch rune) {
	if s.isEndOfLineCharacter(ch) {
		s.column = 0
		s.line++
	} else {
		s.column++
	}
}

func (s *scanner) matchesExistingLexemes() bool {
	if len(s.lexemes) == 0 {
		return false
	}
	var matched []lexeme.Lexeme
	for _, lexeme := range s.lexemes {
		if lexeme.Scan() {
			matched = append(matched, lexeme)
		}
	}

	s.lexemes = matched
	return len(s.lexemes) > 0
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
	return len(s.lexemes) > 0
}

func (s *scanner) matchesNewLexemes(ch rune) (bool, error) {
	return s.matchLexerRules(ch, s.parser.Expected())
}

func (s *scanner) matchLexerRules(ch rune, lexerRules []grammar.LexerRule) (bool, error) {
	anyMatches := false

	for _, lexerRule := range lexerRules {
		if !lexerRule.CanApply(ch) {
			continue
		}

		factory, ok := s.registry[lexerRule.Type()]
		if !ok {
			continue
		}

		tok, err := factory.Create(lexerRule, s.input, s.position)
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
