package client

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/serg-kovalev/word-of-wisdom-tcp-client/pkg/powsolver"
)

const bufferMaxSize = 64

// SolveAndSendPoWSolution solves the Proof of Work challenge and sends the solution to the server.
func SolveAndSendPoWSolution(conn net.Conn) (string, error) {
	// Receive the Proof of Work challenge from the server.
	challenge, difficulty := receivePoWChallenge(conn)
	nonce := powsolver.SolvePoWChallenge(challenge, difficulty)

	// Send the Proof of Work solution to the server.
	_, err := conn.Write([]byte(nonce))
	if err != nil {
		log.Printf("error sending solution to server: %v", err)
	}

	// Receive the server's response.
	response := make([]byte, bufferMaxSize)
	_, err = conn.Read(response)
	if errors.Is(err, io.EOF) {
		log.Println("verification failed: server closed the connection")
		return "", err
	} else if err != nil {
		log.Println("error receiving successful response from server:", err)
	}

	return string(response), err
}

// ProcessQuote receives and processes the random quote from the server.
func ProcessQuote(quote string) {
	log.Println("received quote from server:", quote)
}

func receivePoWChallenge(conn net.Conn) (string, int) {
	// Receive nonce and difficulty from the server.
	buffer := make([]byte, bufferMaxSize)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("error receiving challenge from server:", err)
		return "", 0
	}

	challenge := string(buffer[:n])
	difficulty, _ := difficultyFromChallenge(challenge)

	return challenge, difficulty
}

func difficultyFromChallenge(challenge string) (int, error) {
	difficulty, err := strconv.Atoi(strings.Split(challenge, ":")[0])
	if err != nil {
		log.Println("error parsing difficulty from challenge:", err)
		return 1, err
	}

	return difficulty, nil
}
