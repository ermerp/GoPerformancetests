FROM golang:1.23.3 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN  CGO_ENABLED=0 go build -o "bankGo"

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bankGo ./bankGo

CMD ["./bankGo"]