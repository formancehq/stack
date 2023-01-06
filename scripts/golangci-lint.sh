#!/bin/sh

dir=$(dirname "$0")
source ${dir}/common.sh

for i in $(gomodules); do
  echo "Run golangci-lint on ${i}"
  pushd ${i} >/dev/null
  golangci-lint run --fix >/dev/null
  popd >/dev/null
done
