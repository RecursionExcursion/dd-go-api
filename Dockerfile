FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

#  compiles stdlib for target OS/ARCH combos and caches them
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go install std && \
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go install std && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install std

COPY . .

RUN go build -o dd-go-api

COPY .env .env

EXPOSE 8080
CMD ["./dd-go-api"]
