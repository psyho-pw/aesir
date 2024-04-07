FROM golang:alpine as build

RUN apk add --no-cache go gcc g++
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /go/src/aesir
COPY go.mod go.sum main.go ./
RUN go mod download

COPY Makefile wire.go wire_gen.go ./
COPY src/ ./src
RUN CGO_ENABLED=1 GOOS=linux go build -o out/aesir .
#RUN ldd /go/src/aesir/out/aesir | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /go/src/aesir/out/%
#RUN ln -s ld-musl-x86_64.so.1 /go/src/aesir/out/lib/libc.musl-x86_64.so.1

FROM alpine:edge as prod

ENV APP_ENV="production"
ENV PORT=8000

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Seoul

COPY public/ ./public
#COPY --from=build /go/src/aesir/.env ./
COPY --from=build /go/src/aesir/out/ ./

#USER 65534
ENTRYPOINT ["/aesir"]