FROM golang:1.19.1-alpine3.16 AS builder

WORKDIR /app  

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]
