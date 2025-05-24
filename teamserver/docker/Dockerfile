FROM golang:1.24

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

# Add these lines to your teamserver Dockerfile
# Install Docker CLI only (not the full engine)
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    && curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list \
    && apt-get update \
    && apt-get install -y docker-ce-cli

# Copy the entire project into the container (assuming the project files are in the build context)
COPY ../.. .

# Set the entry point for the container
ENTRYPOINT ["./docker/init_ts.sh"]
