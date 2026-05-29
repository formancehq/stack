#!/usr/bin/env bash
set -euo pipefail

JUSTFILE="${JUSTFILE_PATH:-Justfile}"
OUTPUT="${OUTPUT_FILE:-/tmp/changelog.md}"

extract_version() {
  grep "^${1}" "$JUSTFILE" | sed 's/.*"\(.*\)".*/\1/'
}

LEDGER_VERSION=$(extract_version "LEDGER_VERSION")
PAYMENTS_VERSION=$(extract_version "PAYMENTS_VERSION")
WALLETS_VERSION=$(extract_version "WALLETS_VERSION")
WEBHOOKS_VERSION=$(extract_version "WEBHOOKS_VERSION")
AUTH_VERSION=$(extract_version "AUTH_VERSION")
SEARCH_VERSION=$(extract_version "SEARCH_VERSION")
ORCHESTRATION_VERSION=$(extract_version "ORCHESTRATION_VERSION")
RECONCILIATION_VERSION=$(extract_version "RECONCILIATION_VERSION")
GATEWAY_VERSION=$(extract_version "GATEWAY_VERSION")

get_release_notes() {
  local repo=$1
  local version=$2
  local notes
  notes=$(gh api "repos/formancehq/${repo}/releases/tags/${version}" --jq '.body' 2>/dev/null) || true
  if [ -z "$notes" ] || [ "$notes" = "null" ]; then
    echo "_No release notes available._"
  else
    echo "$notes"
  fi
}

{
  echo "## Component Versions"
  echo ""
  echo "| Component | Version |"
  echo "|-----------|---------|"
  echo "| Ledger | [\`${LEDGER_VERSION}\`](https://github.com/formancehq/ledger/releases/tag/${LEDGER_VERSION}) |"
  echo "| Payments | [\`${PAYMENTS_VERSION}\`](https://github.com/formancehq/payments/releases/tag/${PAYMENTS_VERSION}) |"
  echo "| Wallets | [\`${WALLETS_VERSION}\`](https://github.com/formancehq/wallets/releases/tag/${WALLETS_VERSION}) |"
  echo "| Webhooks | [\`${WEBHOOKS_VERSION}\`](https://github.com/formancehq/webhooks/releases/tag/${WEBHOOKS_VERSION}) |"
  echo "| Auth | [\`${AUTH_VERSION}\`](https://github.com/formancehq/auth/releases/tag/${AUTH_VERSION}) |"
  echo "| Search | [\`${SEARCH_VERSION}\`](https://github.com/formancehq/search/releases/tag/${SEARCH_VERSION}) |"
  echo "| Orchestration | [\`${ORCHESTRATION_VERSION}\`](https://github.com/formancehq/flows/releases/tag/${ORCHESTRATION_VERSION}) |"
  echo "| Reconciliation | [\`${RECONCILIATION_VERSION}\`](https://github.com/formancehq/reconciliation/releases/tag/${RECONCILIATION_VERSION}) |"
  echo "| Gateway | [\`${GATEWAY_VERSION}\`](https://github.com/formancehq/gateway/releases/tag/${GATEWAY_VERSION}) |"
  echo ""
  echo "---"
  echo ""

  declare -a COMPONENTS=(
    "ledger:Ledger:${LEDGER_VERSION}"
    "payments:Payments:${PAYMENTS_VERSION}"
    "wallets:Wallets:${WALLETS_VERSION}"
    "webhooks:Webhooks:${WEBHOOKS_VERSION}"
    "auth:Auth:${AUTH_VERSION}"
    "search:Search:${SEARCH_VERSION}"
    "flows:Orchestration:${ORCHESTRATION_VERSION}"
    "reconciliation:Reconciliation:${RECONCILIATION_VERSION}"
    "gateway:Gateway:${GATEWAY_VERSION}"
  )

  for entry in "${COMPONENTS[@]}"; do
    IFS=':' read -r repo label version <<< "$entry"
    notes=$(get_release_notes "$repo" "$version")
    echo "## ${label} \`${version}\`"
    echo ""
    echo "${notes}"
    echo ""
    echo "---"
    echo ""
  done
} > "$OUTPUT"

echo "Changelog written to $OUTPUT"
