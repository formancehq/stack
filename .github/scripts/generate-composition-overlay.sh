#!/usr/bin/env bash

set -euo pipefail

workflow_file="${1:-.speakeasy/workflow.yaml}"
output_file="${2:-releases/overlays/generated.overlay.json}"

mkdir -p "$(dirname "$output_file")"

{
    while IFS=$'\t' read -r location namespace; do
        # Speakeasy applies modelNamespace to schema names and structural refs,
        # but explicit discriminator mappings need the same transformation.
        yq -o=json "$location" | jq -c --arg namespace "$namespace" '
            def json_path_segment:
                if type == "number"
                then "[" + tostring + "]"
                else "[" + (@json) + "]"
                end;

            .components.schemas as $schemas
            | $schemas
            | to_entries[]
            | .key as $schema_name
            | .value as $schema
            | ($schema | path(.. | objects | select(.discriminator.mapping? != null))) as $path
            | ($schema | getpath($path) | .discriminator.mapping) as $mapping
            | {
                target: (
                    "$.components.schemas["
                    + (($namespace + "_" + $schema_name) | @json)
                    + "]"
                    + ($path | map(json_path_segment) | join(""))
                    + ".discriminator.mapping"
                ),
                update: (
                    $mapping
                    | with_entries(
                        .value |= (
                            . as $ref
                            | if (
                                ($ref | startswith("#/components/schemas/"))
                                and ($schemas[$ref | sub("^#/components/schemas/"; "")] != null)
                            )
                            then (
                                "#/components/schemas/"
                                + $namespace
                                + "_"
                                + ($ref | sub("^#/components/schemas/"; ""))
                            )
                            elif $schemas[$ref] != null
                            then "#/components/schemas/" + $namespace + "_" + $ref
                            else $ref
                            end
                        )
                    )
                )
            }
        '

        # Ledger uses a tagged inline union while discriminator inference is
        # disabled. Convert every singleton resource enum, including future
        # branches, into the constant expected by generated union helpers.
        if [[ "$namespace" == "ledger" ]]; then
            yq -o=json "$location" | jq -c --arg namespace "$namespace" '
                .components.schemas.V2QueryParams.oneOf
                | to_entries[]
                | select(
                    (.value.properties.resource.enum? | type) == "array"
                    and (.value.properties.resource.enum | length) == 1
                )
                | .key as $index
                | .value.properties.resource.enum[0] as $resource
                | {
                    target: (
                        "$.components.schemas[\""
                        + $namespace
                        + "_V2QueryParams\"].oneOf["
                        + ($index | tostring)
                        + "].properties.resource.enum"
                    ),
                    remove: true
                },
                {
                    target: (
                        "$.components.schemas[\""
                        + $namespace
                        + "_V2QueryParams\"].oneOf["
                        + ($index | tostring)
                        + "].properties.resource"
                    ),
                    update: {const: $resource}
                }
            '
        fi
    done < <(
        yq -r '
            .sources[].inputs[]
            | select(.modelNamespace != null)
            | [.location, .modelNamespace]
            | @tsv
        ' "$workflow_file"
    )
} | jq -s '
    {
        overlay: "1.0.0",
        info: {
            title: "Generated namespace fixes for the Formance Stack OpenAPI spec",
            version: "0.0.1"
        },
        actions: .
    }
' > "$output_file"
