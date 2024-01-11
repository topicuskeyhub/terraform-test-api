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

cat << EOF > ~/.terraformrc
provider_installation {

  dev_overrides {
      "registry.terraform.io/topicuskeyhub/keyhub" = "/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF

cd ../terraform-test-api
go install .
