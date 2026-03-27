FROM golang:1.26-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy application source
COPY . .

# Expose port (will be overridden by API_PORT)
EXPOSE 8080

CMD ["go", "run", "main.go"]
