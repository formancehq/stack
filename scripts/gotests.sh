#!/bin/sh

dir=$(dirname "$0")
source ${dir}/common.sh

for i in $(gomodules); do
  echo "Run go test on ${i}"
  pushd ${i} >/dev/null
  go test ./...
  popd >/dev/null
done
