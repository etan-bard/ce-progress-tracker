FROM golang:1.26-alpine

# Install git as it is required by go install for remote modules
RUN apk add --no-cache git

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Set the working directory for migrations
WORKDIR /migrations

# Default entrypoint
ENTRYPOINT ["goose"]
