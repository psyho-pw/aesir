FROM golang:1.20
WORKDIR /go/src/aesir
COPY . .
RUN make build
EXPOSE 3000
RUN cd /go/src/aesir
CMD ["./out/aesir"]