# SDK contract composition

`.speakeasy/workflow.yaml` combines component OpenAPI contracts into
`releases/build/generate.json`. Component models are namespaced; routes address
services through the stack gateway. `releases/base.yaml` owns the global
metadata, OAuth2 scheme, servers and SDK tag groups.

## Build and validate

```sh
nix develop --command just build-openapi
nix develop --command just validate-openapi
nix develop --command just test-openapi
```

The build downloads and prepares contracts, generates the composition overlay,
merges the documents, and validates the result. It does not upload the contract
to the Speakeasy registry. `just pre-commit` also regenerates the event schemas.

## Connectivity source and authentication

`CONNECTIVITY_REF` in `Justfile` pins a commit from Connectivity main. The source
is `misc/connectivity-api/openapi.yaml`; no root-level copy is required. Update
the pin deliberately when adopting a newer contract. Connectivity API releases
do not currently include an OpenAPI asset.

Connectivity is private: local builds require GitHub CLI authentication with
read access to `formancehq/connectivity` (`gh auth login`, or a `GH_TOKEN` with
that access). Anonymous contributors cannot reproduce this part of the build.
The release build uses the existing `NUMARY_GITHUB_TOKEN` secret. Before merging
this integration, verify that the release credential can read the pinned file;
a successful local build does not establish CI credential access. Never print
or commit tokens.

Preparation adds `/api/connectivity` to routes, groups operations under
`connectivity.v1`, and keeps models under the `connectivity` namespace. The
stack OAuth2 scheme replaces the standalone bearer scheme. Public probes remain
public; writes declare `connectivity:write`.

Speakeasy rejects two alternatives for one OAuth2 scheme with different scopes.
Reads therefore declare OAuth2 authentication without a formal scope list and
carry `x-formance-required-scopes-any-of` plus a description stating the server
requirement: `connectivity:read` **or** `connectivity:write`. This avoids excluding
write-only tokens. The server continues to enforce these scopes.

Composition validation checks the Connectivity operations and their security,
response codes, parameters, request/response bodies, all component models,
merge-patch media type and base64 cursor representation. `test-openapi` also
checks preparation idempotence and rejects deliberately damaged contracts.

The pinned Connectivity input declares OpenAPI 3.1.0 but uses schema constructs
compatible with 3.0.3 after wrapping schema references with sibling annotations
in `allOf`. This preserves `readOnly` and descriptions that 3.0 Reference Objects
would otherwise ignore. Preparation declares that component as 3.0.3, and the
merged contract is validated in that dialect. This preserves the meaning of
existing modules’ `nullable` fields. Reassess dialect compatibility whenever
updating the source pin; do not adopt 3.1-only schema constructs by changing
the version string alone.

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
