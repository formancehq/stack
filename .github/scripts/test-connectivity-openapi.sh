#!/usr/bin/env bash
set -euo pipefail

component=$1
composed=$2
script_dir=$(cd "$(dirname "$0")" && pwd)
test_dir=$(mktemp -d)
trap 'rm -rf "$test_dir"' EXIT

bash "$script_dir/validate-connectivity-openapi.sh" "$component" "$composed"
cp "$component" "$test_dir/component.yaml"
bash "$script_dir/prepare-connectivity-openapi.sh" "$test_dir/component.yaml"
yq -o=json -I=0 '.' "$component" > "$test_dir/before.json"
yq -o=json -I=0 '.' "$test_dir/component.yaml" > "$test_dir/after.json"
cmp "$test_dir/before.json" "$test_dir/after.json"

# Each mutation remains parseable JSON but loses a supported API guarantee.
# The integration validator must reject it rather than merely count routes.
for mutation in \
    'del(.paths["/api/connectivity/connectors"].get)' \
    '.paths["/api/connectivity/connectorinstances"].post.security = []' \
    'del(.paths["/api/connectivity/connectors"].get.responses["200"].content)' \
    'del(.components.parameters.connectivity_Cursor)' \
    '.components.schemas.connectivity_ConnectorInstanceSpec.properties.connectivityRef.readOnly = false' \
    '.components.schemas.connectivity_ConnectorInstancePatch.additionalProperties = false' \
    '.components.schemas.connectivity_ConnectorInstanceSpec.properties.initialCursor.type = "integer"'
do
    jq "$mutation" "$composed" > "$test_dir/broken.json"
    if bash "$script_dir/validate-connectivity-openapi.sh" "$component" "$test_dir/broken.json" 2> "$test_dir/error"; then
        echo "Connectivity validation accepted a broken contract: $mutation" >&2
        exit 1
    fi
done
