package re

import (
	"fmt"

	"github.com/patrickhuber/go-earley/forest"
	"github.com/patrickhuber/go-earley/parser"
	"github.com/patrickhuber/go-earley/scanner"
)

func Parse(input string) (*Definition, error) {
	g := Grammar()
	p := parser.New(g)
	s := scanner.New(p, input)
	for {
		ok, err := s.Read()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
	}
	accpeted := s.Parser().Accepted()
	if !accpeted {
		return nil, fmt.Errorf("failed to parse")
	}
	root, ok := s.Parser().GetForestRoot()
	if !ok {
		return nil, fmt.Errorf("failed to get forest root")
	}
	return transform(root)
}

func transform(node forest.Node) (*Definition, error) {
	return &Definition{}, nil
}
