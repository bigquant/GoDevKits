#!/bin/bash -ex

BASE_DIR=$(cd $(dirname $0); pwd)
source "${BASE_DIR}"/common.sh

for file in $(find ${GIT_ROOT}/tests -type f -name "*_test.go"); do
    # go test -v -cover -coverprofile="${file}".cover "${file}"
    go test -v "${file}"
done
