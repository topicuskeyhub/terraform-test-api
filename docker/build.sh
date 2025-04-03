#!/bin/bash

set -xe

TF_VERSION="1.11.3"

apt update
apt -y install unzip libicu72
wget https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip
unzip -d /go/bin terraform_${TF_VERSION}_linux_amd64.zip

cd /work
go work init
go work use ./sdk-go ./terraform-provider-keyhub-generator ./terraform-provider-keyhub ./terraform-test-api

cd terraform-provider-keyhub
go generate ./...
go install .

cd ../terraform-test-api
go install .
