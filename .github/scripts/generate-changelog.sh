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

# component name | repo | version
COMPONENTS=(
  "Ledger|ledger|${LEDGER_VERSION}"
  "Payments|payments|${PAYMENTS_VERSION}"
  "Wallets|wallets|${WALLETS_VERSION}"
  "Webhooks|webhooks|${WEBHOOKS_VERSION}"
  "Auth|auth|${AUTH_VERSION}"
  "Search|search|${SEARCH_VERSION}"
  "Orchestration|flows|${ORCHESTRATION_VERSION}"
  "Reconciliation|reconciliation|${RECONCILIATION_VERSION}"
  "Gateway|gateway|${GATEWAY_VERSION}"
)

# Fetch the release notes body for a given component repo/tag from the public GitHub API.
fetch_release_notes() {
  local repo="$1" version="$2"
  curl -fsSL \
    -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/formancehq/${repo}/releases/tags/${version}" \
    | jq -r '.body // ""' 2>/dev/null || true
}

{
  echo "## Component Versions"
  echo ""
  echo "| Component | Version | Changelog |"
  echo "|-----------|---------|-----------|"
  for entry in "${COMPONENTS[@]}"; do
    IFS='|' read -r name repo version <<< "$entry"
    echo "| ${name} | \`${version}\` | [Release notes](https://github.com/formancehq/${repo}/releases/tag/${version}) |"
  done

  echo ""
  echo "## Product Changelog"

  for entry in "${COMPONENTS[@]}"; do
    IFS='|' read -r name repo version <<< "$entry"
    echo ""
    echo "### ${name} \`${version}\`"
    echo ""
    notes=$(fetch_release_notes "$repo" "$version")
    if [ -n "$notes" ]; then
      echo "$notes"
    else
      echo "_No release notes available._"
    fi
  done
} > "$OUTPUT"

echo "Changelog written to $OUTPUT"
