FROM golang:latest
 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0 \
    go build -ldflags="-s -w" -o certificates-plugin.wasm ./cmd