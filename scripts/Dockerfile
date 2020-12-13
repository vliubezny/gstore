ARG ALPINE_VERSION=3.10
ARG GOLANG_VERSION=1.15.5

FROM golang:${GOLANG_VERSION}
WORKDIR /go/src/github.com/vliubezny/gstore
COPY . .
RUN make linux

FROM alpine:${ALPINE_VERSION}
RUN apk update && apk add --no-cache ca-certificates
COPY --from=0 /go/src/github.com/vliubezny/gstore/build/gstore-linux-amd64 /gstored
ENTRYPOINT [ "/gstored" ]