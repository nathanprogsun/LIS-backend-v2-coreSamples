# build the executable named as coresamples in the golang env aliased as builder
FROM golang:1.19.4-alpine3.16 AS builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/coresamples
# use apk as package manager under alpine
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev
# RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates gcc libtool make musl-dev

COPY Makefile go.mod go.sum ./
RUN go mod download
COPY . .
RUN make tidy build
RUN chmod +x /go/src/coresamples/coresamples

# migrate the built coresamples executable from builder image to the minimal image(scratch)
FROM scratch
COPY --from=builder /go/src/coresamples/etc /etc/ssl/certs/
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /go/src/coresamples/coresamples /coresamples
ENTRYPOINT ["/coresamples"]
CMD []