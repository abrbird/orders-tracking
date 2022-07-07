# https://hub.docker.com/_/golang
FROM golang:1.18

WORKDIR /usr/src/orders_tracking

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./cmd/...  # build all main.go files in ./cmd folder
