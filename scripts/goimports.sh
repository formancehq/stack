#!/bin/sh

dir=$(dirname "$0")
source ${dir}/common.sh

for i in $(gomodules); do
  echo "Run goimports on ${i}"
  pushd ${i} >/dev/null
  goimports -w .
  popd >/dev/null
done
