FROM golang:1.13-buster

WORKDIR /go/src/udp

COPY . .

RUN make release

CMD "release/udp/bin"
