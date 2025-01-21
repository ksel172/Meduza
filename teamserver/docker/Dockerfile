FROM golang:1.23.3

ARG DLV_PORT=2345

# Set the working directory
WORKDIR /app
USER root

# Copy dependency files first for caching
COPY ../../go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN apt-get update && apt-get install -y docker.io
# Install .NET SDK
RUN apt-get update && apt-get install -y wget apt-transport-https && \
    wget https://packages.microsoft.com/config/debian/11/packages-microsoft-prod.deb && \
    dpkg -i packages-microsoft-prod.deb && \
    apt-get update && apt-get install -y dotnet-sdk-8.0 && \
    rm -rf /var/lib/apt/lists/*

## Copy the entire project into the container
COPY ../.. .

ENTRYPOINT ["./docker/init_ts.sh"]
