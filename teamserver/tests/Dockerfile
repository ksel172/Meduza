FROM golang:1.24

WORKDIR /app

# Copy dependency files first for caching
COPY teamserver/go.mod teamserver/go.sum ./teamserver/

# Download Go dependencies
WORKDIR /app/teamserver
RUN go mod download

WORKDIR /app
COPY . .

# Create a non-root user
RUN useradd -m nonrootuser
USER nonrootuser

ENTRYPOINT ["go", "test", "-C", "./teamserver", "./tests/...", "-v"]
