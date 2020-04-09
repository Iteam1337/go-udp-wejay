FROM golang:1.13-stretch AS builder

WORKDIR /build

COPY . .

RUN make release

# # #

FROM golang:1.13-stretch

WORKDIR /app

COPY --from=builder /build/release/udp /app

ENV \
    # SPOTIFY_ID= \
    SPOTIFY_ID= \
    # SPOTIFY_SECRET= \
    SPOTIFY_SECRET= \
    # ADDR=:8090 \
    ADDR= \
    # STORE_STATE=1 \
    STORE_STATE= \
    # SAVE_STATE_LOCATION=/tmp/wejay \
    SAVE_STATE_LOCATION= \
    # GEN_COVER=localhost:8091
    GEN_COVER=

EXPOSE 8090/udp

RUN adduser --disabled-password --gecos '' wejay && \
    chmod -R g+rwX         /app && \
    chgrp -R wejay         /app && \
    chown -R wejay:wejay   /app

USER wejay

ENTRYPOINT [ "/app/bin" ]
CMD [  ]

