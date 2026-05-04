# syntax=docker/dockerfile:1
FROM golang:1.25-alpine

# air for hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# source copied in by compose watch at runtime -
# this COPY is only for plain `docker build` fallback
COPY . .