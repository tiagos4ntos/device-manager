FROM golang:1.25-alpine

# Install necessary build tools
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Generate Swagger docs
RUN swag init -g ./cmd/server.go

# Build the Go binary
RUN go build -a -tags netgo -o main ./cmd/server.go

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]
