# Component versions
LEDGER_VERSION := "v2.4.12"
PAYMENTS_VERSION := "v3.4.0"
WALLETS_VERSION := "v2.2.0"
WEBHOOKS_VERSION := "v2.5.0"
AUTH_VERSION := "v2.5.0"
SEARCH_VERSION := "v2.1.0"
ORCHESTRATION_VERSION := "v2.6.0"
RECONCILIATION_VERSION := "v2.4.0"
GATEWAY_VERSION := "v2.3.1"

# Download all component OpenAPI specs from GitHub releases
download-specs:
    mkdir -p components
    wget -q https://github.com/formancehq/ledger/releases/download/{{ LEDGER_VERSION }}/openapi.yaml -O components/ledger.openapi.yaml
    wget -q https://github.com/formancehq/payments/releases/download/{{ PAYMENTS_VERSION }}/openapi.yaml -O components/payments.openapi.yaml
    wget -q https://raw.githubusercontent.com/formancehq/gateway/{{ GATEWAY_VERSION }}/openapi.yaml -O components/gateway.openapi.yaml
    wget -q https://github.com/formancehq/auth/releases/download/{{ AUTH_VERSION }}/openapi.yaml -O components/auth.openapi.yaml
    wget -q https://github.com/formancehq/search/releases/download/{{ SEARCH_VERSION }}/openapi.yaml -O components/search.openapi.yaml
    wget -q https://github.com/formancehq/webhooks/releases/download/{{ WEBHOOKS_VERSION }}/openapi.yaml -O components/webhooks.openapi.yaml
    wget -q https://github.com/formancehq/wallets/releases/download/{{ WALLETS_VERSION }}/openapi.yaml -O components/wallets.openapi.yaml
    wget -q https://github.com/formancehq/reconciliation/releases/download/{{ RECONCILIATION_VERSION }}/openapi.yaml -O components/reconciliation.openapi.yaml
    wget -q https://github.com/formancehq/flows/releases/download/{{ ORCHESTRATION_VERSION }}/openapi.yaml -O components/orchestration.openapi.yaml

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

# Strip component-level servers blocks so only base.yaml servers survive the merge
strip-servers: prepend-paths
    for f in components/*.openapi.yaml; do yq -i 'del(.servers)' "$f"; done

# Generate namespace-aware composition fixes from the configured inputs.
generate-composition-overlay: strip-servers
    mkdir -p releases/build
    bash .github/scripts/generate-composition-overlay.sh

# Build the merged OpenAPI spec using Speakeasy
build-openapi version="v0.0.0": generate-composition-overlay
    mkdir -p releases/build
    speakeasy run -s all
    cd releases && sed -i'' -e 's/SDK_VERSION/{{ version }}/g' build/generate.json
    just validate-openapi

# Validate invariants introduced by composing namespaced OpenAPI documents.
validate-openapi:
    # Every local discriminator mapping must resolve to an existing schema.
    jq -e '.components.schemas as $schemas | [ $schemas | .. | objects | select(.discriminator.mapping? != null) | .discriminator.mapping[] | select(type == "string" and startswith("#/components/schemas/")) | sub("^#/components/schemas/"; "") | select($schemas[.] == null) ] | length == 0' releases/build/generate.json >/dev/null
    # Ledger query-template helpers rely on these resource constants.
    jq -e '[.components.schemas.ledger_V2QueryParams.oneOf[].properties.resource | select(.const == null or .enum != null)] | length == 0' releases/build/generate.json >/dev/null
    # Error responses must not be generated as successful response variants.
    test "$(jq -c '.paths["x-speakeasy-errors"].statusCodes' releases/build/generate.json)" = '["4XX","5XX","default"]'
    # Every operation tag must be declared globally.
    jq -e '[.tags[].name] as $tags | [.paths[] | to_entries[] | select(.key | IN("get", "post", "put", "patch", "delete", "options", "head", "trace")) | .value.tags[]? | select(. as $tag | $tags | index($tag) == null)] | length == 0' releases/build/generate.json >/dev/null
    # Composition cleanups must remain applied when component specs change.
    jq -e '.components.schemas.auth_Scope == null and .components.schemas.auth_ScopeOptions == null' releases/build/generate.json >/dev/null
    jq -e '[.. | objects | select(.enum? != null) | .enum | select(length != (unique | length))] | length == 0' releases/build/generate.json >/dev/null

# Generate event schemas
generate-events:
    cd events && npm install
    cd events && node index.js

# Build everything (OpenAPI spec + events)
build version="v0.0.0": (build-openapi version) generate-events

# Publish OpenAPI spec to Speakeasy Registry
publish-speakeasy version: generate-composition-overlay
    speakeasy run -s all --registry-tags {{ version }},LATEST_RELEASE

# Pre-commit: build spec and generate events
pre-commit: build-openapi generate-events
