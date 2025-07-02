package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "1234   blaBLub",
			expected: []string{"1234", "blablub"},
		},
		{
			input:    "item1\ttabSeparated",
			expected: []string{"item1", "tabseparated"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)

		actualLen := len(actual)
		expectedLen := len(c.expected)
		if actualLen != expectedLen {
			t.Fatalf(`
For input "%s":
Length mismatch!
Expected length: %d, Actual length: %d.
Expected words: %v
Actual words: %v
`,
				c.input,
				expectedLen,
				actualLen,
				c.expected,
				actual,
			)
		}
		for i := range actual {
			actualWord := actual[i]
			expectedWord := c.expected[i]

			if actualWord != expectedWord {
				t.Errorf(
					`
For input \"%s\",
at index %d: Expected \"%s\" but found \"%s\"!"
`,
					c.input,
					i,
					expectedWord,
					actualWord,
				)
			}
		}
	}
}
