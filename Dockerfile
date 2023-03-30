FROM golang:latest

WORKDIR /go/src/app/URL-Shortening-Service-Go

COPY . .

RUN go build -o main .

CMD ["/go/src/app/URL-Shortening-Service-Go/main"]
