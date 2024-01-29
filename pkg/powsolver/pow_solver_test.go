package powsolver

import (
	"testing"
)

func TestSolvePoWChallenge(t *testing.T) {

	tests := []struct {
		input      string
		difficulty int
		wanted     string
	}{
		{
			input:      "3:zaxj4Ezou2naaVouBEUFlF6sdf32dfsxfsdfv423iRmM30kIUu9ykLYwjTl7Bnc",
			difficulty: 3,
			wanted:     "2800",
		},
		{
			input:      "3:zaxj4Ezou2naaVouBEUFlF6z423iRmM30kIUu9ykLYwjTl7Bnc",
			difficulty: 3,
			wanted:     "355",
		},
		{
			input:      "4:zaxj4Ezou2naaVouBEUFlF6z423iRmM30kIUu9ykLYwjTl7Bnc",
			difficulty: 4,
			wanted:     "53095",
		},
	}

	for _, test := range tests {
		result := SolvePoWChallenge(test.input, test.difficulty)

		if result != test.wanted {
			t.Errorf("SolvePoWChallenge(%s, %d) = %s, expected %s", test.input, test.difficulty, result, test.wanted)
		}
	}
}

func TestCalculateHash(t *testing.T) {
	input := "test_input"
	expectedHash := "952822de6a627ea459e1e7a8964191c79fccfb14ea545d93741b5cf3ed71a09a"

	hash := calculateHash(input)

	if hash != expectedHash {
		t.Errorf("calculateHash(%s) = %s, expected %s", input, hash, expectedHash)
	}
}

func TestIsValidHash(t *testing.T) {
	tests := []struct {
		input      string
		difficulty int
		wanted     bool
	}{
		{
			input:      "000000fc09ec1eca86fcc549b9ce85f27b8dc38d0b669990d3407fdaf3daad7f",
			difficulty: 6,
			wanted:     true,
		},
		{
			input:      "000001fc09ec1eca86fcc549b9ce85f27b8dc38d0b669990d3407fdaf3daad7f",
			difficulty: 6,
			wanted:     false,
		},
		{
			input:      "000001fc09ec1eca86fcc549b9ce85f27b8dc38d0b669990d3407fdaf3daad7f",
			difficulty: 5,
			wanted:     true,
		},
	}

	for _, test := range tests {
		result := isValidHash(test.input, test.difficulty)

		if result != test.wanted {
			t.Errorf("isValidHash(%s, %d) = %t, expected %t", test.input, test.difficulty, result, test.wanted)
		}
	}
}
