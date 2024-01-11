#!/bin/bash

set -xe

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
