#!/bin/bash

set -xe

apt update
apt -y install unzip
wget https://releases.hashicorp.com/terraform/1.6.6/terraform_1.6.6_linux_amd64.zip
unzip -d /go/bin terraform_1.6.6_linux_amd64.zip

cd /work
go work init
go work use ./sdk-go ./terraform-provider-keyhub-generator ./terraform-provider-keyhub ./terraform-test-api

cd terraform-provider-keyhub
go generate ./...
go install .

cd ../terraform-test-api
go install .
