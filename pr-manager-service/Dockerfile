FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr-manager-service ./cmd/

CMD ["./pr-manager-service"]