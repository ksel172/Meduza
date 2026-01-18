FROM golang:1.24

WORKDIR /app

# Copy dependency files first for caching
COPY teamserver/go.mod teamserver/go.sum ./teamserver/

# Download Go dependencies
WORKDIR /app/teamserver
RUN go mod download

WORKDIR /app
COPY . .

# Add a non-root user and switch to that user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["go", "test", "-C", "./teamserver", "./tests/...", "-v"]
