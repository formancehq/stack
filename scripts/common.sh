#!/bin/sh

set -e

gomodules() {
  find services libs \( -name vendor -o -name '[._].*' -o -name node_modules \) -prune -o -name go.mod -print | sed 's:/go.mod$::'
}
