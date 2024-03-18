FROM golang:latest

WORKDIR /go/src/

COPY . .

RUN go mod download

RUN go build -o /app/portfolio

CMD ["/app/portfolio"]
