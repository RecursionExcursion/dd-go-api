# Stage 1 (builds bin)
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o dd-go-api

# Stage 2
# FROM alpine:latest
FROM ubuntu:22.04

# RUN apk --no-cache add ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


WORKDIR /app

COPY --from=builder /app/dd-go-api .
COPY .env .env

EXPOSE 8080

CMD ["./dd-go-api"]
