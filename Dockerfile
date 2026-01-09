FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o geo-alert ./cmd/api

EXPOSE 8080

CMD ["./geo-alert"]
