FROM golang:1.23.3

# Set the working directory
WORKDIR /app
USER root

# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy dependency files first for caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the binary from the `cmd/main` directory
RUN go build -gcflags="-N -l" -o server ./cmd/main

# Start the server with Delve for debugging
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./server"]
