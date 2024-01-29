package client

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockConn struct {
	readData                   []string
	writeData                  string
	closeBeforeChallengeSent   bool
	closeAfterChallengeSent    bool
	closeAfterSolutionReceived bool
}

func TestSolveAndSendPoWSolution(t *testing.T) {
	tests := []struct {
		name             string
		mockConn         *mockConn
		expectedResponse string
		expectedError    error
	}{
		{
			name: "SuccessfulChallengeAndResponse",
			mockConn: &mockConn{
				readData:  []string{"1:challenge123", "server response"},
				writeData: "solution123",
			},
			expectedResponse: "server response",
		},
		{
			name: "ErrorReceivingChallenge",
			mockConn: &mockConn{
				closeBeforeChallengeSent: true,
			},
			expectedError: errors.New("EOF"),
		},
		{
			name: "ErrorSendingSolution",
			mockConn: &mockConn{
				readData:                []string{"1:challenge456", ""},
				closeAfterChallengeSent: true,
			},
			expectedError: errors.New("io: read/write on closed pipe"),
		},
		{
			name: "ErrorReceivingResponse",
			mockConn: &mockConn{
				readData:                   []string{"5:challenge789"},
				writeData:                  "solution123",
				closeAfterSolutionReceived: true,
			},
			expectedError: errors.New("io: read/write on closed pipe"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server, client := net.Pipe()
			client.SetDeadline(time.Now().Add(time.Second))
			defer client.Close()
			go func() {
				defer server.Close()
				if test.mockConn.closeBeforeChallengeSent {
					return
				}
				server.Write([]byte(test.mockConn.readData[0]))
				if test.mockConn.closeAfterChallengeSent {
					return
				}
				if len(test.mockConn.readData) > 1 {
					server.SetReadDeadline(time.Now().Add(time.Second * 5))
					server.Read([]byte(test.mockConn.writeData))
					if test.mockConn.closeAfterSolutionReceived {
						return
					}
					client.SetDeadline(time.Now().Add(time.Second))
					server.Write([]byte(test.mockConn.readData[1]))
				}
			}()

			c := New(client)
			response, err := c.SolveAndSendPoWSolution()
			assert.Equal(t, test.expectedResponse, response, test.name)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestReceivePoWChallenge(t *testing.T) {
	tests := []struct {
		name               string
		mockConn           *mockConn
		expectedChallenge  string
		expectedDifficulty int
		expectedError      error
	}{
		{
			name: "SuccessfulChallenge",
			mockConn: &mockConn{
				readData:                []string{"2:challenge123"},
				closeAfterChallengeSent: true,
			},
			expectedChallenge:  "2:challenge123",
			expectedDifficulty: 2,
		},
		{
			name: "ErrorReceivingChallenge",
			mockConn: &mockConn{
				closeBeforeChallengeSent: true,
			},
			expectedChallenge:  "",
			expectedDifficulty: 0,
			expectedError:      errors.New("EOF"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server, client := net.Pipe()
			client.SetDeadline(time.Now().Add(time.Second))
			defer client.Close()
			go func() {
				defer server.Close()
				if test.mockConn.closeBeforeChallengeSent {
					return
				}
				server.Write([]byte(test.mockConn.readData[0]))
				if test.mockConn.closeAfterChallengeSent {
					return
				}
				if len(test.mockConn.readData) > 1 {
					server.SetReadDeadline(time.Now().Add(time.Second * 5))
					server.Read([]byte(test.mockConn.writeData))
					if test.mockConn.closeAfterSolutionReceived {
						return
					}
					client.SetDeadline(time.Now().Add(time.Second))
					server.Write([]byte(test.mockConn.readData[1]))
				}
			}()

			c := New(client)
			challenge, difficulty, err := c.receivePoWChallenge()
			assert.Equal(t, test.expectedChallenge, challenge)
			assert.Equal(t, test.expectedDifficulty, difficulty)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestDifficultyFromChallenge(t *testing.T) {
	tests := []struct {
		name               string
		challenge          string
		expectedDifficulty int
		expectedError      error
	}{
		{
			name:               "SuccessfulParsing",
			challenge:          "5:example123",
			expectedDifficulty: 5,
			expectedError:      nil,
		},
		{
			name:               "ErrorParsing",
			challenge:          "invalidChallenge",
			expectedDifficulty: 0,
			expectedError:      errors.New("invalid syntax"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			difficulty, err := difficultyFromChallenge(test.challenge)
			assert.Equal(t, test.expectedDifficulty, difficulty)
			assert.Equal(t, test.expectedError, errors.Unwrap(err))
		})
	}
}
