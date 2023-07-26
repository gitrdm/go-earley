package chart

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/state"
)

type Chart struct {
	Sets []Set
}

func New() *Chart {
	return &Chart{}
}

func (c *Chart) Contains(index int, ty state.Type, rule *grammar.DottedRule, origin int) bool {
	set := c.Sets[index]
	return set.Contains(ty, rule, origin)
}
