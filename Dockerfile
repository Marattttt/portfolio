FROM golang:latest

WORKDIR /go/src/

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/portfolio

CMD ["/app/portfolio"]
