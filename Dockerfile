FROM golang:1.12 as builder

ENV GOPROXY https://goproxy.io
ENV GO111MODULE=on

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/src/go-vrs

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o app entry/produce/produce.go

FROM scratch

ENV GIN_MODE=release

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/go-vrs/config.yml /config.yml
COPY --from=builder /go/src/go-vrs/app /

EXPOSE 8080

ENTRYPOINT ["/app"]