FROM golang:alpine AS builder
WORKDIR $GOPATH/src/ymatsiuk/hello
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/hello
FROM scratch
COPY --from=builder /go/bin/hello /go/bin/hello
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/hello"]
