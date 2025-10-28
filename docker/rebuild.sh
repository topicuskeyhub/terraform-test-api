#!/bin/bash

set -xe

if [ $# -eq 0 ]; then
    BRANCH=main
else
	BRANCH=$1
fi

cd /work/sdk-go
git checkout $BRANCH || echo "Using branch 'main'"
go generate .
go install .

cd /work/terraform-provider-keyhub-generator
git checkout $BRANCH || echo "Using branch 'main'"

cd /work/terraform-provider-keyhub
git checkout $BRANCH || echo "Using branch 'main'"
go generate ./...
go install .
