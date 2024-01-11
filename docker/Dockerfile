FROM golang:latest

ADD build.sh /

RUN /build.sh

EXPOSE 8080/tcp
ENTRYPOINT ["/go/bin/terraform-test-api"]
