# word-of-wisdom-tcp-client

## Problem/Purpose

Design and implement “Word of Wisdom” tcp server.

- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## It's a client, not a server

See [Server's README](https://github.com/serg-kovalev/word-of-wisdom-tcp-server/blob/main/README.md)

## How to run a client?

At first, please ensure that server is up and running (see README above).

Have a look into Dockerfile. It creates light-weight docker image on top of alpine where you have compiled binary `/app/word-of-wisdom-tcp-client`.

You can customize hostname and port if needed. Defaults: 0.0.0.0:8080

```shell
go run cmd/main.go -h

Usage: word-of-wisdom-tcp-client [OPTIONS]

Example: word-of-wisdom-tcp-client --address host.docker.internal:8080

Options:
  -a, --address   hostname:port for word-of-wisdom-server (default "localhost:8080")
```

### Run a server in Docker container

```shell
# Build docker image
$ docker build -t word-of-wisdom-client .

# Run an app in docker container
$ docker run word-of-wisdom-client
# 2024/01/29 14:19:46 calling host.docker.internal:8080
# 2024/01/29 14:19:47 iterations: '4281404' challenge: '6:J8sauUO5bweuZ6Qb9RWDzKMta4Uk3uzT8chTExYOpdZG6ewRIE' nonce: '4281404', solution: '000000c8e60bc39e2471da1d5f4601f0a41eb28605a738a7bad5e8d98f40c3e7'
# 2024/01/29 14:19:47 received quote from server: Quote 10
```

### Example

```shell
$ go run cmd/main.go
2024/01/29 14:19:46 calling localhost:8080
2024/01/29 14:19:47 iterations: '4281404' challenge: '6:J8sauUO5bweuZ6Qb9RWDzKMta4Uk3uzT8chTExYOpdZG6ewRIE' nonce: '4281404', solution: '000000c8e60bc39e2471da1d5f4601f0a41eb28605a738a7bad5e8d98f40c3e7'
2024/01/29 14:19:47 received quote from server: Quote 10
```
