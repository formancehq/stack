# Component versions
LEDGER_VERSION := "v2.4.0"
PAYMENTS_VERSION := "v3.2.0"
WALLETS_VERSION := "v2.1.5"
WEBHOOKS_VERSION := "v2.2.0"
AUTH_VERSION := "v2.4.3"
SEARCH_VERSION := "v2.1.0"
ORCHESTRATION_VERSION := "v2.4.1"
RECONCILIATION_VERSION := "v2.2.2"
GATEWAY_VERSION := "v2.2.0"

# Download all component OpenAPI specs from GitHub releases
download-specs:
    mkdir -p components
    wget -q https://github.com/formancehq/ledger/releases/download/{{LEDGER_VERSION}}/openapi.yaml -O components/ledger.openapi.yaml
    wget -q https://github.com/formancehq/payments/releases/download/{{PAYMENTS_VERSION}}/openapi.yaml -O components/payments.openapi.yaml
    wget -q https://github.com/formancehq/gateway/releases/download/{{GATEWAY_VERSION}}/openapi.yaml -O components/gateway.openapi.yaml
    wget -q https://github.com/formancehq/auth/releases/download/{{AUTH_VERSION}}/openapi.yaml -O components/auth.openapi.yaml
    wget -q https://github.com/formancehq/search/releases/download/{{SEARCH_VERSION}}/openapi.yaml -O components/search.openapi.yaml
    wget -q https://github.com/formancehq/webhooks/releases/download/{{WEBHOOKS_VERSION}}/openapi.yaml -O components/webhooks.openapi.yaml
    wget -q https://github.com/formancehq/wallets/releases/download/{{WALLETS_VERSION}}/openapi.yaml -O components/wallets.openapi.yaml
    wget -q https://github.com/formancehq/reconciliation/releases/download/{{RECONCILIATION_VERSION}}/openapi.yaml -O components/reconciliation.openapi.yaml
    wget -q https://github.com/formancehq/flows/releases/download/{{ORCHESTRATION_VERSION}}/openapi.yaml -O components/orchestration.openapi.yaml

# Prepend API path prefix to each component spec
prepend-paths: download-specs
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/auth" + .key) | from_entries)' components/auth.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/ledger" + .key) | from_entries)' components/ledger.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/payments" + .key) | from_entries)' components/payments.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/search" + .key) | from_entries)' components/search.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/webhooks" + .key) | from_entries)' components/webhooks.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/wallets" + .key) | from_entries)' components/wallets.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/orchestration" + .key) | from_entries)' components/orchestration.openapi.yaml
    yq -i '.paths |= (to_entries | map(select(.key == "/*").key = "/api/reconciliation" + .key) | from_entries)' components/reconciliation.openapi.yaml

# Build the merged OpenAPI spec using Speakeasy
build-openapi version="v0.0.0": prepend-paths
    mkdir -p releases/build
    speakeasy run -s all --skip-upload-spec
    cd releases && sed -i'' -e 's/SDK_VERSION/{{version}}/g' build/generate.json

# Generate event schemas
generate-events:
    cd events && npm install
    cd events && node index.js

# Build everything (OpenAPI spec + events)
build version="v0.0.0": (build-openapi version) generate-events

# Publish OpenAPI spec to Speakeasy Registry
publish-speakeasy version: prepend-paths
    speakeasy run -s all --registry-tags {{version}},LATEST_RELEASE

# Pre-commit: build spec and generate events
pre-commit: (build-openapi) generate-events
