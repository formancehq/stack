VERSION 0.8
PROJECT FormanceHQ/stack

IMPORT github.com/formancehq/earthly:tags/v0.15.0 AS core

ARG --global LEDGER_VERSION=v2.4.0
ARG --global PAYMENTS_VERSION=v3.2.0
ARG --global WALLETS_VERSION=v2.1.5
ARG --global WEBHOOKS_VERSION=v2.2.0
ARG --global AUTH_VERSION=v2.4.3
ARG --global SEARCH_VERSION=v2.1.0
ARG --global ORCHESTRATION_VERSION=v2.4.1
ARG --global RECONCILIATION_VERSION=v2.2.2
ARG --global GATEWAY_VERSION=v2.2.0

sources:
    FROM core+base-image
    ARG --required LOCATION
    COPY ${LOCATION} out
    SAVE ARTIFACT out

build-final-spec:
    FROM core+base-image
    RUN apk update && apk add yarn nodejs npm jq
    WORKDIR /src/releases
    COPY releases/package.* .
    RUN npm install

    WORKDIR /src/releases
    COPY releases/base.yaml .
    COPY releases/openapi-overlay.json .
    COPY releases/openapi-merge.json .
    RUN mkdir ./build

    WORKDIR /src/components
    # Download OpenAPI specs from GitHub releases
    RUN wget -q https://github.com/formancehq/ledger/releases/download/${LEDGER_VERSION}/openapi.yaml -O ledger.openapi.yaml
    RUN wget -q https://github.com/formancehq/payments/releases/download/${PAYMENTS_VERSION}/openapi.yaml -O payments.openapi.yaml
    RUN wget -q https://github.com/formancehq/gateway/releases/download/${GATEWAY_VERSION}/openapi.yaml -O gateway.openapi.yaml
    RUN wget -q https://github.com/formancehq/auth/releases/download/${AUTH_VERSION}/openapi.yaml -O auth.openapi.yaml
    RUN wget -q https://github.com/formancehq/search/releases/download/${SEARCH_VERSION}/openapi.yaml -O search.openapi.yaml
    RUN wget -q https://github.com/formancehq/webhooks/releases/download/${WEBHOOKS_VERSION}/openapi.yaml -O webhooks.openapi.yaml
    RUN wget -q https://github.com/formancehq/wallets/releases/download/${WALLETS_VERSION}/openapi.yaml -O wallets.openapi.yaml
    RUN wget -q https://github.com/formancehq/reconciliation/releases/download/${RECONCILIATION_VERSION}/openapi.yaml -O reconciliation.openapi.yaml
    RUN wget -q https://github.com/formancehq/flows/releases/download/${ORCHESTRATION_VERSION}/openapi.yaml -O orchestration.openapi.yaml

    WORKDIR /src/releases
    RUN npm run build
    RUN jq -s '.[0] * .[1] | del(.components.schemas.V2QueryParams.properties.sort."$ref")' build/generate.json openapi-overlay.json > build/latest.json
    ARG version=v0.0.0
    IF [ "$version" = "v0.0.0" ]
        RUN sed -i 's/SDK_VERSION/v0.0.0/g' build/latest.json
        SAVE ARTIFACT build/latest.json AS LOCAL releases/build/latest.json
    ELSE
        RUN sed -i 's/SDK_VERSION/'$version'/g' build/latest.json
        SAVE ARTIFACT build/latest.json AS LOCAL releases/build/generate.json
    END
    SAVE ARTIFACT build/latest.json


build:
    LOCALLY
    BUILD --pass-args +build-final-spec
    BUILD --pass-args ./events+generate

pre-commit: # Generate the final spec and run all the pre-commit hooks
    LOCALLY
    BUILD +build-final-spec
    BUILD ./events+pre-commit
