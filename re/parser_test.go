package re_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-earley/re"
)

func TestParser(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected *re.Definition
	}
	tests := []test{
		{
			name:  "any",
			input: ".",
			expected: &re.Definition{
				Expression: re.ExpressionTerm{
					Term: re.TermFactor{
						Factor: re.FactorAtom{
							Atom: re.AtomAny{},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := re.Parse(test.input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Fatalf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
