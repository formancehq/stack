#!/usr/bin/env bash

dir=$(dirname "$0")
source ${dir}/common.sh

cmdLine="go test -p 8 "
for mod in $(find-updated-modules $@); do
  cmdLine="${cmdLine} ./${mod}/..."
done

echo "Run $cmdLine"
$cmdLine
