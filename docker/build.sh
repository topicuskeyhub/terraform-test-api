#!/bin/bash

set -xe

TF_VERSION="1.14.8"

apt update
apt -y install unzip libicu76
wget https://releases.hashicorp.com/terraform/${TF_VERSION}/terraform_${TF_VERSION}_linux_amd64.zip
unzip -d /go/bin terraform_${TF_VERSION}_linux_amd64.zip

# Cached so a rebuild does not depend on GitHub serving the release download
KIOTA_URL=$(grep -om1 'https://github.com/microsoft/kiota/releases/download/[^ ]*' /work/sdk-go/generate.sh)
mkdir -p /opt/kiota
wget -O /opt/kiota/linux-x64.zip "${KIOTA_URL}"
echo "${KIOTA_URL}" > /opt/kiota/url
echo 'no_clobber = on' > /opt/kiota/wgetrc

cd /work
go work init
go work use ./sdk-go ./terraform-provider-keyhub-generator ./terraform-provider-keyhub ./terraform-test-api

cd terraform-provider-keyhub
go generate ./...
go install .

cd ../terraform-test-api
go install .
