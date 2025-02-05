#!/bin/bash

set -xe

apt update
apt -y install unzip libicu72
wget https://releases.hashicorp.com/terraform/1.10.5/terraform_1.10.5_linux_amd64.zip
unzip -d /go/bin terraform_1.10.5_linux_amd64.zip

cd /work
go work init
go work use ./sdk-go ./terraform-provider-keyhub-generator ./terraform-provider-keyhub ./terraform-test-api

cd terraform-provider-keyhub
go generate ./...
go install .

cd ../terraform-test-api
go install .
