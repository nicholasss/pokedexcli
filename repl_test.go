package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {

	// test cases
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "helloworld",
			expected: []string{"helloworld"},
		},
		{
			input:    "this is a long string",
			expected: []string{"this", "is", "a", "long", "string"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// check the length of the actual slice
		// if they do not match, use t.Errorf to print an error msg
		// and fail the test
		fmt.Println("    ")
		fmt.Println("a:", actual, "len:", len(actual))
		fmt.Println("e:", c.expected, "len:", len(c.expected))

		if len(actual) != len(c.expected) {
			t.Errorf("length of actual different from length of expected.")
			return
		}

		for i := range actual {
			actualWord := actual[i]
			expectedWord := c.expected[i]
			// check each word in the slice
			// if they do not match, use t.Errorf to print an error msg
			// and fail the test

			if actualWord != expectedWord {
				t.Errorf("actual word different from expected word")
			}
		}

	}

}
