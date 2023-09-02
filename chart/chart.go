package chart

import (
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/state"
)

type Chart struct {
	Sets []*Set
}

func New() *Chart {
	return &Chart{}
}

func (c *Chart) Contains(index int, ty state.Type, rule *grammar.DottedRule, origin int) bool {
	set := c.getOrCreateSet(index)
	return set.Contains(ty, rule, origin)
}

func (c *Chart) Enqueue(index int, s state.State) bool {
	set := c.getOrCreateSet(index)
	return set.Enqueue(s)
}

func (c *Chart) GetOrCreate(index int, rule *grammar.DottedRule, origin int) *state.Normal {
	set := c.getOrCreateSet(index)
	return set.GetOrCreate(rule, origin)
}

func (c *Chart) getOrCreateSet(index int) *Set {
	if len(c.Sets) <= index {
		return c.create(index)
	}
	return c.Sets[index]
}

func (c *Chart) create(index int) *Set {
	set := &Set{}
	c.Sets = append(c.Sets, set)
	return set
}
