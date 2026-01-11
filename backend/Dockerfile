FROM golang:1.24-alpine

WORKDIR /app

# Install git (required for go mod download)
# gcc and musl-dev are needed to run testings
RUN apk add --no-cache gcc musl-dev

# Copy go mod files FIRST
COPY go.mod ./
RUN go mod tidy
RUN go mod download

# Copy the rest of the source
COPY . .

# goose up migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build binary
RUN go build -o app

EXPOSE 8080
CMD ["./app"]
