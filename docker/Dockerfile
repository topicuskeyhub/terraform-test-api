FROM golang:latest

ADD build.sh /
ADD work /work
RUN /build.sh
ADD terraformrc /root/.terraformrc

EXPOSE 8080/tcp
ENTRYPOINT ["/go/bin/terraform-test-api"]
