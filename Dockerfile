# Use the official Go image as the base image
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Build the client binary
RUN go build -o word-of-wisdom-client cmd/main.go

# Use a smaller base image for the final image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy only the necessary files from the builder image
COPY --from=builder /app/word-of-wisdom-client /app/word-of-wisdom-client

# Run the client
CMD ["./word-of-wisdom-client", "--address", "host.docker.internal:8080"]
