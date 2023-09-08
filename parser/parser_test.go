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

		g := grammar.New(S,
			grammar.NewProduction(S, S, S),
			grammar.NewProduction(S, b),
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

			(S,2,3) ->
				(b,2,3)
		*/
		root, ok := p.GetForestRoot()
		require.True(t, ok)

		S_0_3 := SymbolEqual(t, root, S, 0, 3)
		require.NotNil(t, S_0_3)
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
		g := grammar.New(S,
			grammar.NewProduction(S, A, T),
			grammar.NewProduction(S, a, T),
			grammar.NewProduction(A, a),
			grammar.NewProduction(A, B, A),
			grammar.NewProduction(B),
			grammar.NewProduction(T, b, b, b))
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

			(S->A*T,01) ->
				(A,0,1)

			(A,0,1) ->
				(a,0,1)
			|	(A->B*A,0,0) (A,0,1)

			(A->B*A,0,0) ->
				(B,0,0)

			(B,0,0)->
		*/
		S_0_4, ok := root.(*forest.Symbol)
		require.True(t, ok)
		require.Equal(t, S, S_0_4.Symbol)
		require.NotNil(t, S_0_4.Internal)
		require.Equal(t, 2, len(S_0_4.Internal.Alternatives))

		S_0_4_0 := S_0_4.Internal.Alternatives[0]
		require.Equal(t, 2, len(S_0_4_0.Children))

		S_at_0_1, ok := S_0_4_0.Children[0].(*forest.Intermediate)
		require.True(t, ok)
		require.Equal(t, 0, S_at_0_1.Origin)
		require.Equal(t, 1, S_at_0_1.Location)
		require.Equal(t, 1, len(S_at_0_1.Internal.Alternatives))
		require.Equal(t, 1, len(S_at_0_1.Internal.Alternatives[0].Children))

		a_0_1, ok := S_at_0_1.Internal.Alternatives[0].Children[0].(*forest.Token)
		require.True(t, ok)
		require.Equal(t, a.TokenType(), a_0_1.Token.Type())
		require.Equal(t, 0, a_0_1.Origin)
		require.Equal(t, 1, a_0_1.Location)

		T_1_4, ok := S_0_4_0.Children[1].(*forest.Symbol)
		require.True(t, ok)
		require.Equal(t, T, T_1_4.Symbol)
		require.Equal(t, 1, T_1_4.Origin)
		require.Equal(t, 4, T_1_4.Location)
		require.Equal(t, 1, len(T_1_4.Internal.Alternatives))
		require.Equal(t, 2, len(T_1_4.Internal.Alternatives[0].Children))

		T_bbb_1_3, ok := T_1_4.Internal.Alternatives[0].Children[0].(*forest.Intermediate)
		require.True(t, ok)
		require.Equal(t, T, T_bbb_1_3.Rule.Production.LeftHandSide)
		require.Equal(t, 2, T_bbb_1_3.Rule.Position)
		require.Equal(t, 1, T_bbb_1_3.Origin)
		require.Equal(t, 3, T_bbb_1_3.Location)
		require.Equal(t, 1, len(T_bbb_1_3.Internal.Alternatives))
		require.Equal(t, 2, len(T_bbb_1_3.Internal.Alternatives[0].Children))

		T_bbb_1_2, ok := T_bbb_1_3.Internal.Alternatives[0].Children[0].(*forest.Intermediate)
		require.True(t, ok)
		require.Equal(t, T, T_bbb_1_2.Rule.Production.LeftHandSide)
		require.Equal(t, 1, T_bbb_1_2.Rule.Position)
		require.Equal(t, 1, T_bbb_1_2.Origin)
		require.Equal(t, 2, T_bbb_1_2.Location)
		require.Equal(t, 1, len(T_bbb_1_2.Internal.Alternatives))
		require.Equal(t, 1, len(T_bbb_1_2.Internal.Alternatives[0].Children))

		b_1_2, ok := T_bbb_1_2.Internal.Alternatives[0].Children[0].(*forest.Token)
		require.True(t, ok)
		require.Equal(t, b.TokenType(), b_1_2.Token.Type())

		b_2_3, ok := T_bbb_1_3.Internal.Alternatives[0].Children[1].(*forest.Token)
		require.True(t, ok)
		require.Equal(t, b.TokenType(), b_2_3.Token.Type())

		require.Equal(t, 2, len(S_0_4.Internal.Alternatives[1].Children))
		require.Equal(t, T_1_4, S_0_4.Internal.Alternatives[1].Children[1])

		S_AT_0_1, ok := S_0_4.Internal.Alternatives[1].Children[0].(*forest.Intermediate)
		require.True(t, ok)
		require.Equal(t, S, S_at_0_1.Rule.Production.LeftHandSide)
		require.Equal(t, 0, S_AT_0_1.Origin)
		require.Equal(t, 1, S_AT_0_1.Location)
		require.Equal(t, 1, len(S_at_0_1.Internal.Alternatives))
		require.Equal(t, 1, len(S_at_0_1.Internal.Alternatives[0].Children))

		A_0_1, ok := S_AT_0_1.Internal.Alternatives[0].Children[0].(*forest.Symbol)
		require.True(t, ok)
		require.Equal(t, A, A_0_1.Symbol)
		require.Equal(t, 0, A_0_1.Origin)
		require.Equal(t, 1, A_0_1.Location)

		require.Equal(t, 2, len(A_0_1.Internal.Alternatives))
		require.Equal(t, 1, len(A_0_1.Internal.Alternatives[0].Children))
		require.Equal(t, a_0_1, A_0_1.Internal.Alternatives[0].Children[0])

		require.Equal(t, 2, len(A_0_1.Internal.Alternatives[1].Children))
		B_0_0, ok := A_0_1.Internal.Alternatives[1].Children[0].(*forest.Symbol)
		require.True(t, ok)
		require.Equal(t, B, B_0_0.Symbol)
		require.Equal(t, 0, B_0_0.Origin)
		require.Equal(t, 0, B_0_0.Location)

		require.Equal(t, A_0_1, A_0_1.Internal.Alternatives[1].Children[1])
	})
}

func SymbolEqual(t *testing.T, node forest.Node, sym grammar.Symbol, origin, location int) *forest.Symbol {
	symbolNode, ok := node.(*forest.Symbol)
	require.True(t, ok)
	require.Equal(t, sym, symbolNode.Symbol)
	require.Equal(t, origin, symbolNode.Origin)
	require.Equal(t, location, symbolNode.Location)
	return symbolNode
}
