FROM golang:latest

COPY . /go/src/app

WORKDIR /go/src/app

RUN go get ./...

CMD ["go", "run", "main.go"]