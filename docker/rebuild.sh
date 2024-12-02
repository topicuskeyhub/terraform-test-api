#!/bin/bash

set -xe

cd /work/sdk-go
go generate .
go install .

cd /work/terraform-provider-keyhub
go generate ./...
go install .
