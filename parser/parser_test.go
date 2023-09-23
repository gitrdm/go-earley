package parser_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/go-earley/forest"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexrule"
	"github.com/patrickhuber/go-earley/parser"
	"github.com/patrickhuber/go-earley/token"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("s series", func(t *testing.T) {

		S := grammar.NewNonTerminal("S")
		s := lexrule.NewString("s")

		// a series of S's
		// S -> S S | S | 's'
		g := grammar.New(
			S,
			// S -> S S
			grammar.NewProduction(S, S, S),
			// S -> S
			grammar.NewProduction(S, S),
			// S -> 's'
			grammar.NewProduction(S, s),
		)

		p := parser.New(g)

		for i := 0; i < 10; i++ {
			tok := token.FromString("s", i, s.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err, "loop %d", i)
			require.True(t, ok, "loop %d", i)
		}
		require.True(t, p.Accepted())
	})
	
	t.Run("underscore_binary", func(t *testing.T) {
		// this is a grammar from pliant that caused some problems
		file := grammar.NewNonTerminal("file")
		ws := grammar.NewNonTerminal("ws")
		ows := lexrule.NewString("_")
		directives := grammar.NewNonTerminal("directives")
		directivesRepeat := grammar.NewNonTerminal("directives_repeat")
		directive := grammar.NewNonTerminal("directive")
		zero := lexrule.NewString("0")
		one := lexrule.NewString("1")

		g := grammar.New(file,
			// file = ws directives ws
			grammar.NewProduction(file, ws, directives, ws),
			// ws = [ ows ]
			// ows = "_"
			grammar.NewProduction(ws, ows),
			grammar.NewProduction(ws),
			// directives = directive { ows directive }
			grammar.NewProduction(directives, directive, directivesRepeat),
			grammar.NewProduction(directivesRepeat, ows, directive, directivesRepeat),
			grammar.NewProduction(directivesRepeat),
			// directive = "0" | "1"
			grammar.NewProduction(directive, zero),
			grammar.NewProduction(directive, one),
		)
		p := parser.New(g)
		input := "_0_1_0_0_1_1_"
		var tokens []*lexrule.String
		for _, r := range input {
			var tok *lexrule.String
			switch r {
			case '_':
				tok = ows
			case '0':
				tok = zero
			case '1':
				tok = one
			}
			tokens = append(tokens, tok)
		}
		RunParse(t, p, tokens)
	})
}

func TestAycockHorspool(t *testing.T) {
	/*
		S' -> S
		S  -> A A A A
		A  -> a | E
		E  ->
	*/
	SPrime := grammar.NewNonTerminal("S'")
	S := grammar.NewNonTerminal("S")
	A := grammar.NewNonTerminal("A")
	E := grammar.NewNonTerminal("E")
	a := lexrule.NewString("a")

	g := grammar.New(SPrime,
		grammar.NewProduction(SPrime, S),
		grammar.NewProduction(S, A, A, A, A),
		grammar.NewProduction(A, a),
		grammar.NewProduction(A, E),
		grammar.NewProduction(E),
	)

	p := parser.New(g)
	tok := token.FromString("a", 0, a.TokenType())
	ok, err := p.Pulse(tok)

	require.NoError(t, err)
	require.True(t, ok)
	require.True(t, p.Accepted())
}

func TestLeo(t *testing.T) {
	t.Run("leo 1", func(t *testing.T) {
		A := grammar.NewNonTerminal("A")
		a := lexrule.NewString("a")

		// A -> A 'a'
		// A ->
		g := grammar.New(A,
			grammar.NewProduction(A, a, A),
			grammar.NewProduction(A),
		)

		p := parser.New(g)
		for i := 0; i < 10; i++ {
			tok := token.FromString("a", i, a.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err, "loop %d", i)
			require.True(t, ok, "loop %d", i)
		}
		require.True(t, p.Accepted())
	})

	t.Run("leo 2", func(t *testing.T) {
		// S -> 'a' S
		// S -> C
		// C -> 'a' C 'b'
		// C ->
		S := grammar.NewNonTerminal("S")
		C := grammar.NewNonTerminal("C")
		a := lexrule.NewString("a")
		b := lexrule.NewString("b")

		g := grammar.New(S,
			grammar.NewProduction(S, a, S),
			grammar.NewProduction(S, C),
			grammar.NewProduction(C, a, C, b),
			grammar.NewProduction(C))

		p := parser.New(g)
		for i := 0; i < 10; i++ {
			lr := a
			if i > 5 {
				lr = b
			}
			tok := token.FromString(lr.Value, i, lr.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err)
			require.True(t, ok)
		}
		require.True(t, p.Accepted())
	})
}

func TestForest(t *testing.T) {

	t.Run("Scott2008_sec4_ex2", func(t *testing.T) {
		S := grammar.NewNonTerminal("S")
		b := lexrule.NewString("b")

		productions := []*grammar.Production{
			grammar.NewProduction(S, S, S),
			grammar.NewProduction(S, b)}

		g := grammar.New(S,
			productions...,
		)

		p := parser.New(g)
		input := []*lexrule.String{b, b, b}
		RunParse(t, p, input)

		/*
			(S,0,3)	->
				(S->S*S,0,2) (S,2,3)
			|	(S->S*S,0,1) (S,1,3)

			(S->S*S,0,2) ->
				(S,0,2)

			(S->S*S,0,1) ->
				(S,0,1)

			(S,0,1) ->
				(b,0,1)

			(S,0,2) ->
				(S->S*S,0,1) (S,1,2)

			(S,1,2) ->
				(b,1,2)

			(S,1,3) ->
				(S->S*S,1,2) (S,2,3)

			(S->S*S,1,2) ->
				(S, 1, 2)

			(S,2,3) ->
				(b,2,3)
		*/
		root, ok := p.GetForestRoot()
		require.True(t, ok)

		S_0_3 := Symbol(S, 0, 3)
		S_0_1 := Symbol(S, 0, 1)
		S_SS_0_1 := Intermediate(Rule(productions[0], 1), 0, 1)
		S_1_2 := Symbol(S, 1, 2)
		S_SS_1_2 := Intermediate(Rule(productions[0], 1), 1, 2)
		S_0_2 := Symbol(S, 0, 2)
		S_SS_0_2 := Intermediate(Rule(productions[0], 1), 0, 2)
		S_2_3 := Symbol(S, 2, 3)
		S_1_3 := Symbol(S, 1, 3)
		b_0_1 := Token(b, 0, 1)
		b_1_2 := Token(b, 1, 2)
		b_2_3 := Token(b, 2, 3)

		Edge(S_0_3, S_SS_0_2, S_2_3)
		Edge(S_0_3, S_SS_0_1, S_1_3)
		Edge(S_SS_0_2, S_0_2)
		Edge(S_SS_0_1, S_0_1)
		Edge(S_0_1, b_0_1)
		Edge(S_0_2, S_SS_0_1, S_1_2)
		Edge(S_1_2, b_1_2)
		Edge(S_1_3, S_SS_1_2, S_2_3)
		Edge(S_SS_1_2, S_1_2)
		Edge(S_2_3, b_2_3)

		Equal(t, S_0_3, root)
	})

	t.Run("Scott2008_sec4_ex3", func(t *testing.T) {
		// S -> AT | aT
		// A -> a | BA
		// B ->
		// T -> b b b
		S := grammar.NewNonTerminal("S")
		A := grammar.NewNonTerminal("A")
		T := grammar.NewNonTerminal("T")
		B := grammar.NewNonTerminal("B")
		a := lexrule.NewString("a")
		b := lexrule.NewString("b")
		productions := []*grammar.Production{
			grammar.NewProduction(S, A, T),
			grammar.NewProduction(S, a, T),
			grammar.NewProduction(A, a),
			grammar.NewProduction(A, B, A),
			grammar.NewProduction(B),
			grammar.NewProduction(T, b, b, b),
		}
		g := grammar.New(S, productions...)
		p := parser.New(g)
		input := []*lexrule.String{a, b, b, b}
		RunParse(t, p, input)

		root, ok := p.GetForestRoot()
		require.True(t, ok)
		require.NotNil(t, root)

		/*
			(S,0,4) ->
				(S->a*T,0,1) (T,1,4)
			|	(S->A*T,0,1) (T,1,4)

			(S->a*T,0,1) ->
				(a,0,1)

			(T,1,4) ->
				(T->bb*b,1,3) (b,3,4)

			(T->bb*b,1,3) ->
				(T->b*bb,1,2) (b,2,3)

			(T->b*bb,1,2) ->
				(b,1,2)

			(S->A*T,0,1) ->
				(A,0,1)

			(A,0,1) ->
				(a,0,1)
			|	(A->B*A,0,0) (A,0,1)

			(A->B*A,0,0) ->
				(B,0,0)

			(B,0,0)->
		*/
		S_0_4 := Symbol(S, 0, 4)
		S_aT_0_1 := Intermediate(Rule(productions[1], 1), 0, 1)
		T_1_4 := Symbol(T, 1, 4)
		T_bbb_1_3 := Intermediate(Rule(productions[5], 2), 1, 3)
		T_bbb_1_2 := Intermediate(Rule(productions[5], 1), 1, 2)
		S_AT_0_1 := Intermediate(Rule(productions[0], 1), 0, 1)
		A_0_1 := Symbol(A, 0, 1)
		B_0_0 := Symbol(B, 0, 0)
		A_BA_0_0 := Intermediate(Rule(productions[3], 1), 0, 0)
		a_0_1 := Token(a, 0, 1)
		b_1_2 := Token(b, 1, 2)
		b_2_3 := Token(b, 2, 3)
		Edge(S_0_4, S_aT_0_1, T_1_4)
		Edge(S_0_4, S_AT_0_1, T_1_4)
		Edge(S_aT_0_1, a_0_1)
		Edge(T_1_4, T_bbb_1_3, b_2_3)
		Edge(T_bbb_1_3, T_bbb_1_2, b_2_3)
		Edge(T_bbb_1_2, b_1_2)
		Edge(S_AT_0_1, A_0_1)
		Edge(A_0_1, a_0_1)
		Edge(A_0_1, A_BA_0_0, A_0_1)
		Edge(A_BA_0_0, B_0_0)

		Equal(t, S_0_4, root)
	})

	t.Run("regex stub", func(t *testing.T) {
		R := grammar.NewNonTerminal("R")
		E := grammar.NewNonTerminal("E")
		T := grammar.NewNonTerminal("T")
		F := grammar.NewNonTerminal("F")
		A := grammar.NewNonTerminal("A")
		pipe := lexrule.NewString("|")
		a := lexrule.NewString("a")

		productions := []*grammar.Production{
			grammar.NewProduction(R, E),
			grammar.NewProduction(E, T),
			grammar.NewProduction(E, T, pipe, E),
			grammar.NewProduction(E),
			grammar.NewProduction(T, F, T),
			grammar.NewProduction(T, F),
			grammar.NewProduction(F, A),
			grammar.NewProduction(A, a),
		}
		g := grammar.New(R, productions...)
		p := parser.New(g, parser.OptimizeRightRecursion(true))
		input := []*lexrule.String{a, a, a, a}
		RunParse(t, p, input)
		root, ok := p.GetForestRoot()
		require.True(t, ok)
		printer := forest.NewPrinter(os.Stdout)
		acceptor, ok := root.(forest.Acceptor)
		require.True(t, ok)
		acceptor.Accept(printer)

		/*
			(R, 0, 4) ->
				(E, 0, 4)
			(E, 0, 4) ->
				(T, 0, 4)
			(T, 0, 4) ->
				(T -> F•T, 0, 1) (T, 1, 4)
			(T -> F•T, 0, 1) ->
				(F, 0, 1)
			(F, 0, 1) ->
				(A, 0, 1)
			(A, 0, 1) ->
				(a, 0, 1)
			(T, 1, 4) ->
				(T -> F•T, 1, 2) (T, 2, 4)
			(T -> F•T, 1, 2) ->
				(F, 1, 2)
			(F, 1, 2) ->
				(A, 1, 2)
			(A, 1, 2) ->
				(a, 1, 2)
			(T, 2, 4) ->
				(T -> F•T, 2, 3) (T, 3, 4)
			(T -> F•T, 2, 3) ->
				(F, 2, 3)
			(F, 2, 3) ->
				(A, 2, 3)
			(A, 2, 3) ->
				(a, 2, 3)
			(T, 3, 4) ->
				(F, 3, 4)
			(F, 3, 4) ->
				(A, 3, 4)
			(A, 3, 4) ->
				(a, 3, 4)
		*/
		R_0_4 := Symbol(R, 0, 4)
		E_0_4 := Symbol(E, 0, 4)
		T_0_4 := Symbol(T, 0, 4)
		T_FT_0_1 := Intermediate(Rule(productions[4], 1), 0, 1)
		T_1_4 := Symbol(T, 1, 4)
		F_0_1 := Symbol(F, 0, 1)
		A_0_1 := Symbol(A, 0, 1)
		a_0_1 := Token(a, 0, 1)
		T_FT_1_2 := Intermediate(Rule(productions[4], 1), 1, 2)
		T_2_4 := Symbol(T, 2, 4)
		F_1_2 := Symbol(F, 1, 2)
		A_1_2 := Symbol(A, 1, 2)
		a_1_2 := Token(a, 1, 2)
		T_FT_2_3 := Intermediate(Rule(productions[4], 1), 2, 3)
		T_3_4 := Symbol(T, 3, 4)
		F_2_3 := Symbol(F, 2, 3)
		A_2_3 := Symbol(A, 2, 3)
		a_2_3 := Token(a, 2, 3)
		F_3_4 := Symbol(F, 3, 4)
		A_3_4 := Symbol(A, 3, 4)
		a_3_4 := Token(a, 3, 4)
		Edge(R_0_4, E_0_4)
		Edge(E_0_4, T_0_4)
		Edge(T_0_4, T_FT_0_1, T_1_4)
		Edge(T_FT_0_1, F_0_1)
		Edge(F_0_1, A_0_1)
		Edge(A_0_1, a_0_1)
		Edge(T_1_4, T_FT_1_2, T_2_4)
		Edge(T_FT_1_2, F_1_2)
		Edge(F_1_2, A_1_2)
		Edge(A_1_2, a_1_2)
		Edge(T_2_4, T_FT_2_3, T_3_4)
		Edge(T_FT_2_3, F_2_3)
		Edge(F_2_3, A_2_3)
		Edge(A_2_3, a_2_3)
		Edge(T_3_4, F_3_4)
		Edge(F_3_4, A_3_4)
		Edge(A_3_4, a_3_4)

		Equal(t, R_0_4, root)
	})

	t.Run("leo 1", func(t *testing.T) {
		A := grammar.NewNonTerminal("A")
		a := lexrule.NewString("a")

		// A -> A 'a'
		// A ->
		g := grammar.New(A,
			grammar.NewProduction(A, a, A),
			grammar.NewProduction(A),
		)

		p := parser.New(g, parser.OptimizeRightRecursion(true))
		for i := 0; i < 4; i++ {
			tok := token.FromString("a", i, a.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err, "loop %d", i)
			require.True(t, ok, "loop %d", i)
		}
		require.True(t, p.Accepted())

		root, ok := p.GetForestRoot()
		require.True(t, ok)

		printer := forest.NewPrinter(os.Stdout)
		acceptor, ok := root.(forest.Acceptor)
		require.True(t, ok)
		acceptor.Accept(printer)
	})
}

func RunParse(t *testing.T, p parser.Parser, input []*lexrule.String) {
	for i, sym := range input {
		tok := token.FromString(sym.Value, i, sym.TokenType())
		ok, err := p.Pulse(tok)
		require.NoError(t, err)
		require.True(t, ok)
	}
	require.True(t, p.Accepted())
}

func Symbol(sym grammar.Symbol, origin, location int, alternatives ...forest.Group) *forest.Symbol {
	return forest.NewSymbol(sym, origin, location, alternatives...)
}

func Intermediate(rule *grammar.DottedRule, origin, location int, alternatives ...forest.Group) *forest.Intermediate {
	return forest.NewIntermediate(rule, origin, location, alternatives...)
}

func Token(rule grammar.LexerRule, origin, location int) *forest.Token {
	return forest.NewToken(token.FromString(rule.String(), origin, rule.TokenType()), origin, location)
}

func Alternative(nodes ...forest.Node) forest.Group {
	return forest.NewGroup(nodes...)
}

func Equal(t *testing.T, expectedNode forest.Node, actualNode forest.Node) {
	cache := map[forest.Node]struct{}{}
	expectedWork := []forest.Node{expectedNode}
	actualWork := []forest.Node{actualNode}

	for len(expectedWork) > 0 && len(actualWork) > 0 {

		expectedNode = expectedWork[0]
		expectedWork = expectedWork[1:]

		actualNode = actualWork[0]
		actualWork = actualWork[1:]

		_, ok := cache[expectedNode]
		if ok {
			continue
		}
		cache[expectedNode] = struct{}{}

		var expectedInternal forest.Internal
		var actualInternal forest.Internal

		switch n := expectedNode.(type) {
		case *forest.Intermediate:
			i, ok := actualNode.(*forest.Intermediate)
			require.True(t, ok, "%s != %s", n, actualNode)
			IntermediateEqual(t, n, i)
			expectedInternal = n
			actualInternal = i
		case *forest.Symbol:
			s, ok := actualNode.(*forest.Symbol)
			require.True(t, ok, "%s != %s", n, actualNode)
			SymbolEqual(t, n, s)
			expectedInternal = n
			actualInternal = s
		case *forest.Token:
			tok, ok := actualNode.(*forest.Token)
			require.True(t, ok, "%s != %s", n, actualNode)
			TokenEqual(t, n, tok)
			return
		}

		InternalEqual(t, expectedInternal, actualInternal)
		for g := 0; g < len(expectedInternal.Alternatives()); g++ {
			alt1 := expectedInternal.Alternatives()[g]
			alt2 := actualInternal.Alternatives()[g]
			for c := 0; c < len(alt1.Children()); c++ {
				c1 := alt1.Children()[c]
				expectedWork = append(expectedWork, c1)
				c2 := alt2.Children()[c]
				actualWork = append(actualWork, c2)
			}
		}
	}
	require.Equal(t, len(expectedWork), len(actualWork))
}

func SymbolEqual(t *testing.T, expected, actual *forest.Symbol) {
	require.Equal(t, expected.Origin(), actual.Origin(), "%s != %s", expected.String(), actual.String())
	require.Equal(t, expected.Location(), actual.Location(), "%s != %s", expected.String(), actual.String())
	require.Equal(t, expected.Symbol, actual.Symbol, "%s != %s", expected.String(), actual.String())
}

func IntermediateEqual(t *testing.T, expected, actual *forest.Intermediate) {
	require.Equal(t, expected.Origin(), actual.Origin(), "%s != %s", expected.String(), actual.String())
	require.Equal(t, expected.Location(), actual.Location(), "%s != %s", expected.String(), actual.String())
	RuleEqual(t, expected.Rule, actual.Rule)
}

func InternalEqual(t *testing.T, expected, actual forest.Internal) {
	if expected == nil && actual == nil {
		return
	}
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	require.Equal(t, len(expected.Alternatives()), len(actual.Alternatives()), "%s != %s", expected, actual)
	for i := 0; i < len(expected.Alternatives()); i++ {
		expectedAlternative := expected.Alternatives()[i]
		actualAlternative := actual.Alternatives()[i]
		require.Equal(t, len(expectedAlternative.Children()), len(actualAlternative.Children()))
	}
}

func RuleEqual(t *testing.T, expected *grammar.DottedRule, actual *grammar.DottedRule) {
	require.Equal(t, expected.Position, actual.Position)
	ProductionEqual(t, expected.Production, actual.Production)
}

func ProductionEqual(t *testing.T, expected *grammar.Production, actual *grammar.Production) {
	require.Equal(t, expected.LeftHandSide, actual.LeftHandSide)
	require.Equal(t, len(expected.RightHandSide), len(actual.RightHandSide))
	for i := 0; i < len(expected.RightHandSide); i++ {
		require.Equal(t, expected.RightHandSide[i], actual.RightHandSide[i])
	}
}

func TokenEqual(t *testing.T, expected, actual *forest.Token) {
	require.Equal(t, expected.Location(), actual.Location(), "%s != %s", expected.String(), actual.String())
	require.Equal(t, expected.Origin(), actual.Origin(), "%s != %s", expected.String(), actual.String())
	require.Equal(t, expected.Token.Type(), actual.Token.Type(), "%s != %s", expected.String(), actual.String())
	require.Equal(t, expected.Token.Position(), actual.Token.Position(), "%s != %s", expected.String(), actual.String())
}

func Rule(production *grammar.Production, position int) *grammar.DottedRule {
	return &grammar.DottedRule{
		Production: production,
		Position:   position,
	}
}

func Edge(internal forest.Internal, nodes ...forest.Node) {
	if len(nodes) == 2 {
		internal.AddUniqueFamily(nodes[0], nodes[1])
	} else if len(nodes) == 1 {
		internal.AddUniqueFamily(nodes[0], nil)
	}
}
