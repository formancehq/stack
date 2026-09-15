#!/usr/bin/env bash

set -euo pipefail

workflow_file="${1:-.speakeasy/workflow.yaml}"
output_file="${2:-releases/overlays/generated.overlay.json}"

mkdir -p "$(dirname "$output_file")"

base_input_count=$(yq -r '
    [.sources[].inputs[] | select(.modelNamespace == null)]
    | length
' "$workflow_file")

if [[ "$base_input_count" != "1" ]]; then
    echo "expected exactly one base OpenAPI input, found $base_input_count" >&2
    exit 1
fi

base_location=$(yq -r '
    .sources[].inputs[]
    | select(.modelNamespace == null)
    | .location
' "$workflow_file")

{
    # Component metadata is merged with the base metadata by Speakeasy. Restore
    # the canonical info object from the sole non-namespaced base input.
    yq -o=json -I=0 '{"target": "$.info", "update": .info}' "$base_location"

    while IFS=$'\t' read -r location namespace; do
        # Speakeasy applies modelNamespace to schema names and structural refs,
        # but explicit discriminator mappings need the same transformation.
        # Extract only discriminator metadata before converting to JSON. Some
        # component specs contain valid uint64 bounds that yq cannot marshal
        # through its signed JSON integer representation.
        yq -o=json -I=0 '
            .components.schemas as $schemas
            | $schemas
            | ..
            | select(kind == "map" and .discriminator.mapping != null)
            | {
                "path": path,
                "mapping": .discriminator.mapping,
                "schema_names": ($schemas | keys)
            }
            | select(.mapping != null)
        ' "$location" | jq -c --arg namespace "$namespace" '
            def json_path_segment:
                if type == "number"
                then "[" + tostring + "]"
                else "[" + (@json) + "]"
                end;

            .path[2] as $schema_name
            | .path[3:] as $path
            | .mapping as $mapping
            | .schema_names as $schema_names
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
                                and (
                                    ($schema_names | index($ref | sub("^#/components/schemas/"; "")))
                                    != null
                                )
                            )
                            then (
                                "#/components/schemas/"
                                + $namespace
                                + "_"
                                + ($ref | sub("^#/components/schemas/"; ""))
                            )
                            elif ($schema_names | index($ref)) != null
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
            yq -o=json -I=0 '.components.schemas.V2QueryParams.oneOf // []' "$location" \
                | jq -c --arg namespace "$namespace" '
                to_entries[]
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
            title: "Generated composition fixes for the Formance Stack OpenAPI spec",
            version: "0.0.1"
        },
        actions: .
    }
' > "$output_file"
