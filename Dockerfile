FROM golang:1.20.2-alpine3.17 as BUILD

WORKDIR /app
COPY main.go .
COPY go.mod .
RUN GOINSECURE="*" go mod tidy
RUN GOOS=linux go build -o scanner .

FROM alpine:3.17 AS FINAL
WORKDIR /app
COPY --from=BUILD /app .

CMD ["./scanner"]