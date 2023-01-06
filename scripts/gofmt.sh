#!/bin/sh

dir=$(dirname "$0")
source ${dir}/common.sh

for i in $(gomodules); do
  echo "Run go fmt on ${i}"
  pushd ${i} >/dev/null
  gofmt -w .
  popd >/dev/null
done
