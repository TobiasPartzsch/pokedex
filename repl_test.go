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
			expected: []string{"hello", "world"},
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
			input:    "item1	tabSeparated",
			expected: []string{"item1", "tabseparated"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			actualWord := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if actualWord != expectedWord {
				t.Errorf("Expected \"%s\" but found \"%s\"!", expectedWord, actualWord)
			}
		}
	}
}
