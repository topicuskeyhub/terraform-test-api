FROM golang:latest as builder

ADD build.sh /
ADD work /work
RUN /build.sh

FROM golang:latest

COPY --from=builder /go/bin /go/bin/
ADD terraformrc /root/.terraformrc

EXPOSE 8080/tcp
ENTRYPOINT ["/go/bin/terraform-test-api"]
