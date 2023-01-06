#!/bin/sh

dir=$(dirname "$0")
source ${dir}/common.sh

for i in $(gomodules); do
  echo "Run go mod tidy on ${i}"
  pushd ${i} >/dev/null
  go mod tidy
  popd >/dev/null
done
