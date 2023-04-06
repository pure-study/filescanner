FROM golang:1.20.2-alpine3.17 as BUILD

WORKDIR /app
COPY . .
RUN GOINSECURE="*" go mod tidy
RUN go build -o scanner .

CMD [ "./scanner" ]