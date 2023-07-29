FROM golang:1.20-alpine as build

RUN apk add --no-cache build-base tzdata

WORKDIR /go/src/aesir
COPY go.mod go.sum main.go ./
RUN go mod download

COPY .env Makefile wire.go wire_gen.go ./
COPY src/ ./src
RUN go build -ldflags='-s -w' -o out/aesir .
RUN ldd /go/src/aesir/out/aesir | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /go/src/aesir/out%
RUN ln -s ld-musl-x86_64.so.1 /go/src/aesir/out/lib/libc.musl-x86_64.so.1

FROM scratch as prod

COPY public/ ./public
COPY --from=build /go/src/aesir/.env ./
COPY --from=build /go/src/aesir/out/ ./
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Seoul
#USER 65534
ENTRYPOINT ["/aesir"]