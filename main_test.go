package main

import (
	"testing"
)

var argumentParsingTests = []struct {
	arguments             []string
	expectedPairLength    int
	expectedCommandLength int
}{
	{[]string{"command"}, 0, 1},
	{[]string{"foo=bar", "command"}, 1, 1},
	{[]string{"foo=bar"}, 1, 0},
	{[]string{"foo=bar", "taco=burrito", "command", "arg"}, 2, 2},
}

func TestArgumentParsing(t *testing.T) {
	for _, tt := range argumentParsingTests {
		pairs, args := parseArguments(tt.arguments)
		if len(pairs) != tt.expectedPairLength {
			t.Errorf("Expected %d keypairs in \"%s\", got %d", tt.expectedPairLength, tt.arguments, len(pairs))
		}
		if len(args) != tt.expectedCommandLength {
			t.Errorf("Expected %d command argument in \"%s\", got %d", tt.expectedCommandLength, tt.arguments, len(args))
		}
	}
}
