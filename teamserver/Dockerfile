FROM golang:1.23.3

WORKDIR /app
USER root

# https://antonio-si.medium.com/a-few-tips-on-remote-debugging-golang-applications-running-in-an-m1-docker-container-68606326e83e
RUN go install github.com/go-delve/delve/cmd/dlv@latest
COPY go.mod go.sum ./

COPY . .

RUN go build -gcflags="-N -l" -a -o server .

CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./server"]