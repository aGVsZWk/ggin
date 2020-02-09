FROM golang:latest

WORKDIR $GOPATH/src/ggin
COPY . $GOPATH/src/ggin

RUN go build .

EXPOSE 8000

ENTRYPOINT ["./ggin"]