FROM golang:1.24

WORKDIR /app

# Copy dependency files first for caching
COPY teamserver/go.mod teamserver/go.sum ./teamserver/

# Download Go dependencies
RUN go mod download

COPY . .

ENTRYPOINT ["go", "test", "./tests/...", "-v"]