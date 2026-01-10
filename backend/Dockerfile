FROM golang:1.23-alpine

WORKDIR /app

# Install git (required for go mod download)
RUN apk add --no-cache git

# Copy go mod files FIRST
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# goose up migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build binary
RUN go build -o app

EXPOSE 8080
CMD ["./app"]
