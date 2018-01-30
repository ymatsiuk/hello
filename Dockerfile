FROM golang:latest as builder
RUN go get -d -v github.com/ymatsiuk/hello 
WORKDIR /go/src/github.com/ymatsiuk/hello
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hello .

FROM gruebel/upx:latest as upx
COPY --from=builder /go/src/github.com/ymatsiuk/hello/hello /hello.orig
RUN upx --best --lzma -o /hello /hello.orig

FROM alpine:latest  
COPY --from=upx /hello /hello
CMD ["/hello"] 

