#!/bin/bash

set -xe

if [ $# -eq 0 ]; then
    BRANCH=main
else
	BRANCH=$1
fi

cd /work/sdk-go
git fetch
git checkout $BRANCH || echo "Using branch 'main'"
if grep -qF "$(cat /opt/kiota/url)" generate.sh; then
	mkdir -p tmp
	cp /opt/kiota/linux-x64.zip tmp/
fi
WGETRC=/opt/kiota/wgetrc go generate .
go install .

cd /work/terraform-provider-keyhub-generator
git fetch
git checkout $BRANCH || echo "Using branch 'main'"

cd /work/terraform-provider-keyhub
git fetch
git checkout $BRANCH || echo "Using branch 'main'"
go generate ./...
go install .
