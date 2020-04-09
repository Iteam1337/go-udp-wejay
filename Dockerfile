FROM golang:1.13-stretch AS builder

WORKDIR /go/src/udp

COPY . .

RUN make release

FROM golang:1.13-alpine AS main

WORKDIR /app

COPY --from=builder /go/src/udp /app

ENV SPOTIFY_ID= \
    SPOTIFY_SECRET= \
    ADDR=:8090 \
    STORE_STATE=1 \
    SAVE_STATE_LOCATION=/tmp/wejay \
    GEN_COVER=

EXPOSE 8090/udp

CMD /app/bin
