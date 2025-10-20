VERSION 0.8
PROJECT FormanceHQ/stack

IMPORT github.com/formancehq/earthly:tags/v0.15.0 AS core

ARG LEDGER_VERSION=v2.3.1
ARG PAYMENTS_VERSION=v3.0.18
ARG WALLETS_VERSION=v2.1.5
ARG WEBHOOKS_VERSION=v2.2.0
ARG AUTH_VERSION=v2.4.0
ARG SEARCH_VERSION=v2.1.0
ARG ORCHESTRATION_VERSION=v2.4.0
ARG RECONCILIATION_VERSION=v2.2.1
ARG GATEWAY_VERSION=v2.1.0

sources:
    FROM core+base-image
    ARG --required LOCATION
    COPY ${LOCATION} out
    SAVE ARTIFACT out

build-final-spec:
    FROM core+base-image
    RUN apk update && apk add yarn nodejs npm jq curl
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
    RUN curl -L https://github.com/formancehq/ledger/releases/download/${LEDGER_VERSION}/openapi.yaml -o ledger.openapi.yaml
    RUN curl -L https://github.com/formancehq/payments/releases/download/${PAYMENTS_VERSION}/openapi.yaml -o payments.openapi.yaml
    RUN curl -L https://github.com/formancehq/gateway/releases/download/${GATEWAY_VERSION}/openapi.yaml -o gateway.openapi.yaml
    RUN curl -L https://github.com/formancehq/auth/releases/download/${AUTH_VERSION}/openapi.yaml -o auth.openapi.yaml
    RUN curl -L https://github.com/formancehq/search/releases/download/${SEARCH_VERSION}/openapi.yaml -o search.openapi.yaml
    RUN curl -L https://github.com/formancehq/webhooks/releases/download/${WEBHOOKS_VERSION}/openapi.yaml -o webhooks.openapi.yaml
    RUN curl -L https://github.com/formancehq/wallets/releases/download/${WALLETS_VERSION}/openapi.yaml -o wallets.openapi.yaml
    RUN curl -L https://github.com/formancehq/reconciliation/releases/download/${RECONCILIATION_VERSION}/openapi.yaml -o reconciliation.openapi.yaml
    RUN curl -L https://github.com/formancehq/flows/releases/download/${ORCHESTRATION_VERSION}/openapi.yaml -o orchestration.openapi.yaml

    WORKDIR /src/releases
    RUN npm run build
    RUN jq -s '.[0] * .[1]' build/generate.json openapi-overlay.json > build/latest.json
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
