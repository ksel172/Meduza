FROM golang:1.23.3

ARG DLV_PORT=2345

# Set the working directory
WORKDIR /app
USER root

# Copy dependency files first for caching
COPY ../../go.mod go.sum ./

# Download dependencies
RUN go mod download

## Copy the entire project into the container
COPY ../.. .

ENTRYPOINT ["./docker/init_ts.sh"]