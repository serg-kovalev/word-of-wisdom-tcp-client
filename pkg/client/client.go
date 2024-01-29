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

// Client represents a client
type Client struct {
	conn net.Conn
}

// New creates a new client
func New(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// SolveAndSendPoWSolution solves the Proof of Work challenge and sends the solution to the server.
func (c Client) SolveAndSendPoWSolution() (string, error) {
	// Receive the Proof of Work challenge from the server.
	challenge, difficulty, err := c.receivePoWChallenge()
	if err != nil {
		return "", err
	}

	nonce := powsolver.SolvePoWChallenge(challenge, difficulty)
	// Send the Proof of Work solution to the server.
	if _, err = c.sendPoWSolution(nonce); err != nil {
		log.Printf("error sending solution to server: %v", err)
		return "", err
	}

	// Receive the server's response.
	response, err := c.receiveServerVerificationResponse()
	if errors.Is(err, io.EOF) {
		log.Println("verification failed: server closed the connection")
		return "", err
	} else if err != nil {
		log.Println("error receiving successful response from server:", err)
	}

	return response, err
}

// ProcessQuote receives and processes the random quote from the server.
func (c Client) ProcessQuote(quote string) {
	log.Println("received quote from server:", quote)
}

func (c Client) receivePoWChallenge() (string, int, error) {
	// Receive nonce and difficulty from the server.
	buffer := make([]byte, bufferMaxSize)
	n, err := c.conn.Read(buffer)
	if err != nil {
		log.Println("error receiving challenge from server:", err)
		return "", 0, err
	}

	challenge := string(buffer[:n])
	difficulty, _ := difficultyFromChallenge(challenge)

	return challenge, difficulty, nil
}

func difficultyFromChallenge(challenge string) (int, error) {
	difficulty, err := strconv.Atoi(strings.Split(challenge, ":")[0])
	if err != nil {
		log.Println("error parsing difficulty from challenge:", err)
		return 0, err
	}

	return difficulty, nil
}

func (c Client) sendPoWSolution(nonce string) (int, error) {
	// Send the Proof of Work solution to the server.
	return c.conn.Write([]byte(nonce))
}

func (c Client) receiveServerVerificationResponse() (string, error) {
	// Receive the server's response.
	response := make([]byte, bufferMaxSize)
	n, err := c.conn.Read(response)

	return string(response[:n]), err
}
