package api

import (
	"testing"
)

func TestCleanProfane(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "kerfuffle sharbert fornax",
			expected: "**** **** ****",
		},
		{
			input:    "kerfuffle! sharbert! fornax!",
			expected: "kerfuffle! sharbert! fornax!",
		},
		{
			input:    "kerfuffle sharbert! forna",
			expected: "**** sharbert! forna",
		},
	}

	for _, case_ := range cases {
		actual := cleanProfane(case_.input)
		if case_.expected != actual {
			t.Errorf("not matcing %s vs %s", actual, case_.expected)
		}
	}
}
