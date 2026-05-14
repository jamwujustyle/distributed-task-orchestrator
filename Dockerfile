# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS lambda-builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda

FROM golang:1.25-alpine
RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY --from=lambda-builder /build/bootstrap /app/bin/bootstrap