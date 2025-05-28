FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o log-service ./cmd

EXPOSE 8083

CMD ["./log-service"]
