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

{
  echo "## Component Versions"
  echo ""
  echo "| Component | Version | Changelog |"
  echo "|-----------|---------|-----------|"
  echo "| Ledger | \`${LEDGER_VERSION}\` | [Release notes](https://github.com/formancehq/ledger/releases/tag/${LEDGER_VERSION}) |"
  echo "| Payments | \`${PAYMENTS_VERSION}\` | [Release notes](https://github.com/formancehq/payments/releases/tag/${PAYMENTS_VERSION}) |"
  echo "| Wallets | \`${WALLETS_VERSION}\` | [Release notes](https://github.com/formancehq/wallets/releases/tag/${WALLETS_VERSION}) |"
  echo "| Webhooks | \`${WEBHOOKS_VERSION}\` | [Release notes](https://github.com/formancehq/webhooks/releases/tag/${WEBHOOKS_VERSION}) |"
  echo "| Auth | \`${AUTH_VERSION}\` | [Release notes](https://github.com/formancehq/auth/releases/tag/${AUTH_VERSION}) |"
  echo "| Search | \`${SEARCH_VERSION}\` | [Release notes](https://github.com/formancehq/search/releases/tag/${SEARCH_VERSION}) |"
  echo "| Orchestration | \`${ORCHESTRATION_VERSION}\` | [Release notes](https://github.com/formancehq/flows/releases/tag/${ORCHESTRATION_VERSION}) |"
  echo "| Reconciliation | \`${RECONCILIATION_VERSION}\` | [Release notes](https://github.com/formancehq/reconciliation/releases/tag/${RECONCILIATION_VERSION}) |"
  echo "| Gateway | \`${GATEWAY_VERSION}\` | [Release notes](https://github.com/formancehq/gateway/releases/tag/${GATEWAY_VERSION}) |"
} > "$OUTPUT"

echo "Changelog written to $OUTPUT"
