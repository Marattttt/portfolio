FROM golang:latest

WORKDIR /app/src/

COPY . .

RUN go mod download

RUN go build -o /app/portfolio

CMD ["/app/portfolio"]
