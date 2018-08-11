FROM golang:latest

COPY . /go/src/golang-test-task

WORKDIR /go/src/golang-test-task

RUN go get ./... && go build

CMD ["golang-test-task"]