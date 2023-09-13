package parser_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/forest"
	"github.com/patrickhuber/go-earley/grammar"
	"github.com/patrickhuber/go-earley/lexrule"
	"github.com/patrickhuber/go-earley/parser"
	"github.com/patrickhuber/go-earley/token"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
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
		for i, sym := range input {
			tok := token.FromString(sym.Value, i, sym.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err)
			require.True(t, ok)
		}
		require.True(t, p.Accepted())

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

		S_0_3.Internal.Alternatives = append(S_0_3.Internal.Alternatives,
			Alternative(S_SS_0_2, S_2_3),
			Alternative(S_SS_0_1, S_1_3))
		S_SS_0_2.Internal.Alternatives = append(S_SS_0_2.Internal.Alternatives,
			Alternative(S_0_2))
		S_SS_0_1.Internal.Alternatives = append(S_SS_0_1.Internal.Alternatives,
			Alternative(S_0_1))
		S_0_1.Internal.Alternatives = append(S_0_1.Internal.Alternatives,
			Alternative(b_0_1))
		S_0_2.Internal.Alternatives = append(S_0_2.Internal.Alternatives,
			Alternative(S_SS_0_1, S_1_2))
		S_1_2.Internal.Alternatives = append(S_1_2.Internal.Alternatives,
			Alternative(b_1_2))
		S_1_3.Internal.Alternatives = append(S_1_3.Internal.Alternatives,
			Alternative(S_SS_1_2, S_2_3))
		S_SS_1_2.Internal.Alternatives = append(S_SS_1_2.Internal.Alternatives,
			Alternative(S_1_2))
		S_2_3.Internal.Alternatives = append(S_2_3.Internal.Alternatives,
			Alternative(b_2_3))

		Equal(t, root, S_0_3)
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
		for i, sym := range input {
			tok := token.FromString(sym.Value, i, sym.TokenType())
			ok, err := p.Pulse(tok)
			require.NoError(t, err)
			require.True(t, ok)
		}
		require.True(t, p.Accepted())

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

		S_0_4.Internal.Alternatives = append(S_0_4.Internal.Alternatives,
			Alternative(S_aT_0_1, T_1_4),
			Alternative(S_AT_0_1, T_1_4))
		S_aT_0_1.Internal.Alternatives = append(S_aT_0_1.Internal.Alternatives,
			Alternative(a_0_1))
		T_1_4.Internal.Alternatives = append(T_1_4.Internal.Alternatives,
			Alternative(T_bbb_1_3, b_2_3))
		T_bbb_1_3.Internal.Alternatives = append(T_bbb_1_3.Internal.Alternatives,
			Alternative(T_bbb_1_2, b_2_3))
		T_bbb_1_2.Internal.Alternatives = append(T_bbb_1_2.Internal.Alternatives,
			Alternative(b_1_2))
		S_AT_0_1.Internal.Alternatives = append(S_AT_0_1.Internal.Alternatives,
			Alternative(A_0_1))
		A_0_1.Internal.Alternatives = append(A_0_1.Internal.Alternatives,
			Alternative(a_0_1),
			Alternative(A_BA_0_0, A_0_1))
		A_BA_0_0.Internal.Alternatives = append(A_BA_0_0.Internal.Alternatives,
			Alternative(B_0_0))
		Equal(t, S_0_4, root)
	})
}

func Symbol(sym grammar.Symbol, origin, location int, alternatives ...*forest.Group) *forest.Symbol {
	return &forest.Symbol{
		Symbol:   sym,
		Origin:   origin,
		Location: location,
		Internal: &forest.Internal{
			Alternatives: alternatives,
		},
	}
}

func Intermediate(rule *grammar.DottedRule, origin, location int, alternatives ...*forest.Group) *forest.Intermediate {
	return &forest.Intermediate{
		Origin:   origin,
		Location: location,
		Rule:     rule,
		Internal: &forest.Internal{
			Alternatives: alternatives,
		},
	}
}

func Token(rule grammar.LexerRule, origin, location int) *forest.Token {
	return &forest.Token{
		Token:    token.FromString(rule.String(), origin, rule.TokenType()),
		Origin:   origin,
		Location: location,
	}
}

func Alternative(nodes ...forest.Node) *forest.Group {
	return &forest.Group{
		Children: nodes,
	}
}

func Equal(t *testing.T, n1 forest.Node, n2 forest.Node) {
	cache := map[forest.Node]struct{}{}
	work1 := []forest.Node{n1}
	work2 := []forest.Node{n2}

	for len(work1) > 0 && len(work2) > 0 {

		n1 = work1[0]
		work1 = work1[1:]

		n2 = work2[0]
		work2 = work2[1:]

		_, ok := cache[n1]
		if ok {
			continue
		}
		cache[n1] = struct{}{}

		var internal1 *forest.Internal
		var internal2 *forest.Internal

		switch n := n1.(type) {
		case *forest.Intermediate:
			i, ok := n2.(*forest.Intermediate)
			require.True(t, ok, "%s != %s", n, n2)
			IntermediateEqual(t, n, i)
			internal1 = n.Internal
			internal2 = i.Internal
		case *forest.Symbol:
			s, ok := n2.(*forest.Symbol)
			require.True(t, ok, "%s != %s", n, n2)
			SymbolEqual(t, n, s)
			internal1 = n.Internal
			internal2 = s.Internal
		case *forest.Token:
			tok, ok := n2.(*forest.Token)
			require.True(t, ok, "%s != %s", n, n2)
			TokenEqual(t, n, tok)
			return
		}

		InternalEqual(t, internal1, internal2)
		for g := 0; g < len(internal1.Alternatives); g++ {
			alt1 := internal1.Alternatives[g]
			alt2 := internal2.Alternatives[g]
			for c := 0; c < len(alt1.Children); c++ {
				c1 := alt1.Children[c]
				work1 = append(work1, c1)
				c2 := alt2.Children[c]
				work2 = append(work2, c2)
			}
		}
	}
	require.Equal(t, len(work1), len(work2))
}

func SymbolEqual(t *testing.T, s1, s2 *forest.Symbol) {
	require.Equal(t, s1.Origin, s2.Origin, "%s != %s", s1.String(), s2.String())
	require.Equal(t, s1.Location, s2.Location, "%s != %s", s1.String(), s2.String())
	require.Equal(t, s1.Symbol, s2.Symbol, "%s != %s", s1.String(), s2.String())
}

func IntermediateEqual(t *testing.T, i1, i2 *forest.Intermediate) {
	require.Equal(t, i1.Origin, i2.Origin, "%s != %s", i1.String(), i2.String())
	require.Equal(t, i1.Location, i2.Location, "%s != %s", i1.String(), i2.String())
	RuleEqual(t, i1.Rule, i2.Rule)
}

func InternalEqual(t *testing.T, i1, i2 *forest.Internal) {
	if i1 == nil && i2 == nil {
		return
	}
	require.NotNil(t, i1)
	require.NotNil(t, i2)
	require.Equal(t, len(i1.Alternatives), len(i2.Alternatives))
	for i := 0; i < len(i1.Alternatives); i++ {
		alt1 := i1.Alternatives[i]
		alt2 := i2.Alternatives[i]
		require.Equal(t, len(alt1.Children), len(alt2.Children))
	}
}

func RuleEqual(t *testing.T, one *grammar.DottedRule, two *grammar.DottedRule) {
	require.Equal(t, one.Position, two.Position)
	ProductionEqual(t, one.Production, two.Production)
}

func ProductionEqual(t *testing.T, one *grammar.Production, two *grammar.Production) {
	require.Equal(t, one.LeftHandSide, two.LeftHandSide)
	require.Equal(t, len(one.RightHandSide), len(two.RightHandSide))
	for i := 0; i < len(one.RightHandSide); i++ {
		require.Equal(t, one.RightHandSide[i], two.RightHandSide[i])
	}
}

func TokenEqual(t *testing.T, t1, t2 *forest.Token) {
	require.Equal(t, t1.Location, t2.Location, "%s != %s", t1.String(), t2.String())
	require.Equal(t, t1.Origin, t2.Origin, "%s != %s", t1.String(), t2.String())
	require.Equal(t, t1.Token.Type(), t2.Token.Type(), "%s != %s", t1.String(), t2.String())
	require.Equal(t, t1.Token.Position(), t2.Token.Position(), "%s != %s", t1.String(), t2.String())
}

func Rule(production *grammar.Production, position int) *grammar.DottedRule {
	return &grammar.DottedRule{
		Production: production,
		Position:   position,
	}
}
