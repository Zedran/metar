package metar

import "testing"

/* Tests PrepareCodes function, assessing length and content against expected output. */
func TestPrepareCodes(t *testing.T) {
	type testCase struct {
		input    []string
		expected []string
	}

	cases := []testCase{
		{input: []string{""},                     expected: []string{}              },
		{input: []string{"abc", "abcd"},          expected: []string{"ABCD"}        },
		{input: []string{"abcD", "ABCD", "abcd"}, expected: []string{"ABCD"}        },
		{input: []string{"abcD", "ABCD", "EFGH"}, expected: []string{"ABCD", "EFGH"}},
	}

	for _, c := range cases {
		output := PrepareCodes(c.input...)

		if len(output) != len(c.expected) {
			t.Errorf("Failed for %v (length) :: output: %v -- expected: %v", c.input, output, c.expected)
		}

		for i := range output {
			if output[i] != c.expected[i] {
				t.Errorf("Failed for %v (content) :: output: %v -- expected: %v", c.input, output, c.expected)
			}
		}		
	}
}
