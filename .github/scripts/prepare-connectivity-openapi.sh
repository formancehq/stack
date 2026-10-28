#!/usr/bin/env bash
set -euo pipefail

source_json=$(mktemp)
prepared_json=$(mktemp)
trap 'rm -f "$source_json" "$prepared_json"' EXIT
yq -o=json '.' "$1" > "$source_json"

# Adapt existing operations only. Public probes keep security: []; read scope
# and write scope are alternatives. Speakeasy rejects OAuth2 alternatives for
# one scheme, so retain authentication and describe the server scope check.
# The pinned Connectivity schema uses the OpenAPI 3.0-compatible subset. Keep
# the composed dialect at 3.0.3 so existing nullable schemas retain their meaning.
jq '
    "Requires connectivity:read or connectivity:write. The SDK supplies the OAuth2 bearer token; the server enforces scopes." as $read_note
    | .paths |= with_entries(
        .key = "/api/connectivity" + (.key | sub("^/api/connectivity"; ""))
        | .value |= with_entries(
            .key as $method
            | if ($method | IN("get", "post", "put", "patch", "delete", "options", "head", "trace")) then
                .value |= (
                    .operationId = "connectivity_" + (.operationId | sub("^connectivity_"; ""))
                    | .tags = ["connectivity.v1"]
                    | if .security == [] then .
                      elif ($method | IN("get", "head")) then
                        .security = [{"Authorization": []}]
                        | .["x-formance-required-scopes-any-of"] = ["connectivity:read", "connectivity:write"]
                        | .description = ((.description // "") | if endswith($read_note) then . elif . == "" then $read_note else . + "\n\n" + $read_note end)
                      else .security = [{"Authorization": ["connectivity:write"]}]
                      end
                )
              else . end
        )
    )
    # OpenAPI 3.0 ignores siblings of Reference Objects. Preserve the 3.1
    # schema annotations (notably connectivityRef.readOnly) on allOf wrappers.
    | walk(if type == "object" and has("$ref") then
        if (."$ref" | startswith("#/components/schemas/")) and (keys | length > 1) then
            .allOf = [{"$ref": ."$ref"}] | del(."$ref")
        else . end
      else . end)
    | .openapi = "3.0.3"
    | del(.security, .components.securitySchemes.bearerAuth)
    | .tags = [{"name": "connectivity.v1", "description": "Connectivity tenant API"}]
' "$source_json" > "$prepared_json"
yq -oy -P -p=json '.' "$prepared_json" > "$1"
