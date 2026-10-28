#!/usr/bin/env bash
set -euo pipefail

source_json=$(mktemp)
trap 'rm -f "$source_json"' EXIT
yq -o=json '.' "$1" > "$source_json"

# Check the actual composition against the prepared component rather than a
# hard-coded operation count, so newly added upstream operations are covered.
jq -e --slurpfile source "$source_json" '
    def namespace_refs:
        walk(if type == "object" and has("$ref") then
            ."$ref" |= sub("^#/components/(?<category>[^/]+)/(?<name>.*)$"; "#/components/\(.category)/connectivity_\(.name)")
        else . end);
    def without_composition_metadata:
        walk(if type == "object" then del(.["x-speakeasy-name-override"], .["x-speakeasy-model-namespace"]) else . end);
    def operations:
        .paths | to_entries[] | .key as $path | .value | to_entries[]
        | select(.key | IN("get", "post", "put", "patch", "delete", "options", "head", "trace"))
        | {path: $path, method: .key, operation: .value};
    . as $merged
    | [$source[0] | operations] as $expected
    | [operations | select(.path | startswith("/api/connectivity/"))] as $actual
    | (.openapi == "3.0.3")
      and ($expected | length > 0)
      and (($expected | length) == ($actual | length))
      and all($expected[];
        . as $entry
        | $merged.paths[$entry.path][$entry.method] as $operation
        | ($operation.operationId == $entry.operation.operationId)
          and ($operation.tags == ["connectivity.v1"])
          and ($operation.security == $entry.operation.security)
          and (($operation.responses | keys) == ($entry.operation.responses | keys))
          and (($operation.requestBody.content // {} | keys) == ($entry.operation.requestBody.content // {} | keys))
      )
      and all($source[0].paths | to_entries[];
        . as $entry
        | ($merged.paths[$entry.key] | without_composition_metadata) == ($entry.value | namespace_refs | without_composition_metadata)
      )
      and all($source[0].components | to_entries[];
        .key as $category
        | all(.value | to_entries[];
            . as $entry
            | ($merged.components[$category]["connectivity_" + $entry.key] | without_composition_metadata) == ($entry.value | namespace_refs | without_composition_metadata)
        )
      )
      and ($merged.components.securitySchemes.Authorization.flows.clientCredentials.scopes | has("connectivity:read") and has("connectivity:write"))
      and ($merged.paths["/api/connectivity/_healthcheck"].get.security == [])
      and ($merged.paths["/api/connectivity/_info"].get.security == [])
      and ($merged.components.securitySchemes.Authorization.type == "oauth2")
      and ($merged.components.securitySchemes.bearerAuth == null)
      and ($merged.components.schemas.connectivity_ConnectorInstanceSpec.properties.connectivityRef.readOnly == true)
      and ($merged.components.schemas.connectivity_ConnectorInstanceSpec.properties.connectivityRef.allOf[0]."$ref" == "#/components/schemas/connectivity_LocalObjectRef")
      and ($merged.components.schemas.connectivity_ConnectorInstanceSpec.properties.connectorRef.description | length > 0)
      and ($merged.components.schemas.connectivity_ConnectorInstanceSpec.properties.initialCursor.type == "string")
      and ($merged.components.schemas.connectivity_ConnectorInstanceSpec.properties.initialCursor.format == "byte")
      and ($merged.paths["/api/connectivity/connectorinstances/{name}"].patch.requestBody.content["application/merge-patch+json"] != null)
' "$2" > /dev/null || {
    echo 'Connectivity operations, security, errors or configuration contract changed during composition' >&2
    exit 1
}
