# SDK contract composition

`.speakeasy/workflow.yaml` combines component OpenAPI contracts into
`releases/build/generate.json`. Component models are namespaced; routes address
services through the stack gateway. `releases/base.yaml` owns the global
metadata, OAuth2 scheme, servers and SDK tag groups.

## Build and validate

```sh
nix develop --command just build-openapi
nix develop --command just validate-openapi
```

The build downloads and prepares contracts, generates the composition overlay,
merges the documents, and validates the result. It does not upload the contract
to the Speakeasy registry. `just pre-commit` also regenerates the event schemas.

## Connectivity source and authentication

`CONNECTIVITY_REF` in `Justfile` follows Connectivity `main`. Each build fetches
`misc/connectivity-api/openapi.yaml` at the current branch head; no root-level
copy is required. Connectivity API releases do not currently include an OpenAPI
asset. A change on Connectivity main is included in the next stack composition,
so rebuilding an older stack revision can consume a newer Connectivity contract.

Connectivity is private: local builds require GitHub CLI authentication with
read access to `formancehq/connectivity` (`gh auth login`, or a `GH_TOKEN` with
that access). Anonymous contributors cannot reproduce this part of the build.
The release build uses the existing `NUMARY_GITHUB_TOKEN` secret. Before merging
this integration, verify that the release credential can read the source file;
a successful local build does not establish CI credential access. Never print
or commit tokens.

Preparation adds `/api/connectivity` to routes, groups operations under
`connectivity.v1`, and prefixes operation IDs to avoid collisions with other
components. Models use the `connectivity` namespace. As for Ledger, the source
bearer authentication and operation descriptions remain unchanged. The API
server enforces Connectivity scopes; the stack does not rewrite these into
OAuth2 scope requirements. Composition reuses the existing bearer scheme; tokens
still need the scopes accepted by the service being called.

Connectivity must provide an OpenAPI 3.0.3 contract, including `allOf` wrappers
for schema references with annotations such as `readOnly` and descriptions.
The stack does not convert its dialect. This keeps existing modules’ `nullable`
fields meaningful in the composed OpenAPI 3.0.3 document. The compatibility change was integrated through
[Connectivity #493](https://github.com/formancehq/connectivity/pull/493).

## Publication and downstream SDKs

After review and merge, an authorized `vX.Y.Z` stack tag triggers
`.github/workflows/releases.yml`. It publishes `openapi.json` and `events.json`,
then publishes `stacks-source` to the Speakeasy registry with the version and
`LATEST_RELEASE` tags. This includes the API contract fetched from private
Connectivity in the public stack release assets and SDK input; review that
publication boundary before release.

The SDK repositories consume `stacks-source@LATEST_RELEASE`. For example,
`formance-sdk-go` runs its Generate workflow on demand or on its daily schedule,
and opens a generation PR. Publishing that SDK and updating consumers such as
fctl are separate steps; changing Connectivity main alone performs none of them.
