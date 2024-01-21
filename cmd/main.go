package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/serg-kovalev/word-of-wisdom-tcp-client/pkg/client"

	cli "github.com/jawher/mow.cli"
)

func main() {
	cliApp := cli.App("word-of-wisdom-tcp-client", "TCP client for word-of-wisdom-server")
	cliApp.LongDesc = "Example: word-of-wisdom-tcp-client --address host.docker.internal:8080"
	hostOpt := cliApp.StringOpt("a address", "localhost:8080", "hostname:port for word-of-wisdom-server")

	cliApp.Action = func() {
		host := strings.TrimSpace(*hostOpt)
		log.Printf("calling %s\n", host)

		// Connect to the server.
		conn, err := net.Dial("tcp", host)
		if err != nil {
			log.Fatal("error connecting to server:", err)
		}
		defer conn.Close()

		// Solve the Proof of Work challenge and send the solution to the server.
		response, err := client.SolveAndSendPoWSolution(conn)
		if err == nil {
			// Receive and process the random quote from the server.
			client.ProcessQuote(response)
		}
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatalf("can't run CLI app: %v", err)
	}
}
