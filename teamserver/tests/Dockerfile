FROM golang:1.24

WORKDIR /app

# Copy dependency files first for caching
COPY teamserver/go.mod teamserver/go.sum ./teamserver/

# Download Go dependencies
WORKDIR /app/teamserver
RUN go mod download

# Install .NET SDK
# USER root
# RUN apt-get update && apt-get install -y wget apt-transport-https && \
#     wget https://packages.microsoft.com/config/debian/11/packages-microsoft-prod.deb && \
#     dpkg -i packages-microsoft-prod.deb && \
#     apt-get update && apt-get install -y dotnet-sdk-8.0 && \
#     rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

ENTRYPOINT ["go", "test", "-C", "./teamserver", "./tests/...", "-v"]