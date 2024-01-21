package powsolver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

// SolvePoWChallenge solves the Proof of Work challenge.
func SolvePoWChallenge(challenge string, difficulty int) string {
	var nonce, curHash string
	var i int
	for {
		i++
		nonce = fmt.Sprintf("%d", i)
		curHash = calculateHash(challenge + nonce)

		// Check if the solution meets the required difficulty.
		if isValidHash(curHash, difficulty) {
			break
		}
	}
	log.Printf("iterations: '%d' challenge: '%s' nonce: '%s', solution: '%s'", i, challenge, nonce, curHash)

	return nonce
}

// calculateHash calculates the SHA-256 hash of a string.
func calculateHash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))

	return hex.EncodeToString(hash.Sum(nil))
}

// isValidHash checks if the hash meets the required difficulty.
func isValidHash(hash string, difficulty int) bool {
	prefix := fmt.Sprintf("%0*d", difficulty, 0)

	return hash[:difficulty] == prefix
}
