FROM golang:1.23.3

ARG DLV_PORT=2345

# Set the working directory
WORKDIR /app

# Switch to root user
USER root

# Copy dependency files first for caching
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

USER root
# Install .NET SDK
RUN apt-get update && apt-get install -y wget apt-transport-https && \
    wget https://packages.microsoft.com/config/debian/11/packages-microsoft-prod.deb && \
    dpkg -i packages-microsoft-prod.deb && \
    apt-get update && apt-get install -y dotnet-sdk-8.0 && \
    rm -rf /var/lib/apt/lists/*

# Copy the entire project into the container (assuming the project files are in the build context)
COPY ../.. .

# Set the entry point for the container
ENTRYPOINT ["./docker/init_ts.sh"]
