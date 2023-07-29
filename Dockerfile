FROM golang:1.20 as build

WORKDIR /go/src/aesir
COPY go.mod go.sum main.go ./
RUN go mod download

COPY .env Makefile wire.go wire_gen.go ./
COPY src/ ./src
RUN make build

FROM golang:1.20 as prod

WORKDIR /go/src/aesir

COPY public/ ./public
COPY --from=build /go/src/aesir/.env ./
COPY --from=build /go/src/aesir/out/ ./

ENTRYPOINT ["./aesir"]