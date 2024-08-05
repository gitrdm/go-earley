package re

import (
	"github.com/patrickhuber/go-earley/automata/dfa"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/terminal"
)

func Grammar() *grammar.Grammar {
	definition := nonTerminal("definition")
	expression := nonTerminal("expression")
	term := nonTerminal("term")
	factor := nonTerminal("factor")
	atom := nonTerminal("atom")
	set := nonTerminal("set")
	positiveSet := nonTerminal("positive_set")
	negativeSet := nonTerminal("negative_Set")
	characterClass := nonTerminal("character_class")
	characterRange := nonTerminal("character_range")
	character := nonTerminal("character")
	characterClassCharacter := nonTerminal("character_class_character")

	upCaret := oneOf('^')
	dollar := oneOf('$')
	pipe := oneOf('|')
	iterator := oneOf('*', '+', '?')
	openBracket := oneOf('[')
	closeBracket := oneOf(']')
	dash := oneOf('-')
	notMeta := not(oneOf('^', '.', '$', '(', ')', '[', ']', '+', '*', '?', '\\', '/'))
	escape := sequence("escape_sequence", oneOf('\\'), anyCh())
	notCloseBracket := not(closeBracket)
	dot := oneOf('.')

	productions := []*grammar.Production{
		// definition
		production(definition, expression),
		production(definition, upCaret, expression),
		production(definition, expression, dollar),
		production(definition, upCaret, expression, dollar),
		// expression
		production(expression, term),
		production(expression, term, pipe, expression),
		// term
		production(term, factor),
		production(term, factor, term),
		// factor
		production(factor, atom),
		production(factor, atom, iterator),
		// atom
		production(atom, character),
		production(atom, expression),
		production(atom, dot),
		// set
		production(set, positiveSet),
		production(set, negativeSet),
		// positive_set
		production(positiveSet, openBracket, characterClass, closeBracket),
		// negative_set
		production(negativeSet, openBracket, upCaret, characterClass, closeBracket),
		// character_class
		production(characterClass, characterRange),
		production(characterClass, characterRange, characterClass),
		// character_range
		production(characterRange, characterClassCharacter),
		production(characterRange, characterClassCharacter, dash, characterClassCharacter),
		// character
		production(character, notMeta),
		production(character, escape),
		// character_class_character
		production(characterClassCharacter, notCloseBracket),
		production(characterClassCharacter, escape),
	}
	return grammar.New(definition, productions...)
}

func production(lhs grammar.NonTerminal, rhs ...grammar.Symbol) *grammar.Production {
	return grammar.NewProduction(lhs, rhs...)
}

func nonTerminal(name string) grammar.NonTerminal {
	return grammar.NewNonTerminal(name)
}

func oneOf(runes ...rune) grammar.Terminal {
	if len(runes) == 1 {
		return terminal.NewCharacter(runes[0])
	}
	var terminals []grammar.Terminal
	for _, ch := range runes {
		t := terminal.NewCharacter(ch)
		terminals = append(terminals, t)
	}
	return terminal.NewSet(terminals)
}

func not(t grammar.Terminal) grammar.Terminal {
	return terminal.NewNegate(t)
}

func anyCh() grammar.Terminal {
	return terminal.NewAny()
}

func sequence(name string, terminals ...grammar.Terminal) grammar.LexerRule {
	start := &dfa.State{
		Final: false,
	}
	current := start
	for _, term := range terminals {
		state := &dfa.State{
			Final: false,
		}
		current.Transitions = append(current.Transitions, dfa.Transition{
			Target:   state,
			Terminal: term,
		})
		current = state
	}
	current.Final = true
	return dfa.NewDfa(start, name)
}
