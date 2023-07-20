package token

import "github.com/patrickhuber/go-earley/lexrule"

type stringToken struct {
	index int
	rule  lexrule.String
	span  AccumulatorSpan
}

func (tok *stringToken) Reset() {
	tok.index = 0
	// TODO reset span?
}

func (tok *stringToken) Scan() bool {
	if tok.index >= len(tok.rule.Value()) {
		return false
	}
	peek := tok.span.Peek()
	if peek == EOF {
		return false
	}
	if peek != tok.rule.Value()[tok.span.Length()] {
		return false
	}
	tok.index++
	return tok.span.Grow()
}

func (tok *stringToken) Accepted() bool {
	return tok.index == len(tok.rule.Value())
}

func NewString(rule lexrule.String, span AccumulatorSpan) Token {
	return &stringToken{
		span:  span,
		index: 0,
		rule:  rule,
	}
}
