#!/usr/bin/env bash

dir=$(dirname "$0")
source ${dir}/common.sh

cmdLine="go test -v "
for mod in $(find-updated-modules $@); do
  cmdLine="${cmdLine} ./${mod}/..."
done

$cmdLine
