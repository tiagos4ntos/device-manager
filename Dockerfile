FROM golang:1.25-alpine

# Install necessary build tools
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/server.go

EXPOSE 8080

CMD ["./main"]
