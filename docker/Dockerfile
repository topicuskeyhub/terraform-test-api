FROM golang:latest as builder

ADD build.sh /
ADD work /
RUN /build.sh

FROM golang:latest

COPY --from=otel-builder /go/bin /go

EXPOSE 8080/tcp
ENTRYPOINT ["/build-and-run.sh"]
