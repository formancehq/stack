VERSION 0.8
PROJECT FormanceHQ/stack

IMPORT github.com/formancehq/earthly:tags/v0.15.0 AS core
IMPORT github.com/formancehq/ledger:v2.1.2 AS ledger
IMPORT github.com/formancehq/webhooks:v2.1.0 AS webhooks
IMPORT github.com/formancehq/stack/components/payments:v2.0.25 AS payments
IMPORT github.com/formancehq/wallets:v2.1.0 AS wallets
IMPORT github.com/formancehq/gateway:main AS gateway
IMPORT github.com/formancehq/stack/ee/auth:v2.0.24 AS auth
IMPORT github.com/formancehq/stack/ee/search:v2.0.24 AS search
IMPORT github.com/formancehq/stack/ee/stargate:v2.0.24 AS stargate
IMPORT github.com/formancehq/stack/ee/orchestration:v2.0.24 AS orchestration
IMPORT github.com/formancehq/stack/ee/reconciliation:v2.0.24 AS reconciliation

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

    COPY (ledger+openapi/openapi.yaml) /src/components/ledger/
    COPY (payments+openapi/openapi.yaml) /src/components/payments/
    COPY (gateway+openapi/openapi.yaml) /src/ee/gateway/
    COPY (auth+openapi/openapi.yaml) /src/ee/auth/
    COPY (search+openapi/openapi.yaml) /src/ee/search/
    COPY (webhooks+openapi/openapi.yaml) /src/ee/webhooks/
    COPY (wallets+openapi/openapi.yaml) /src/ee/wallets/
    COPY (reconciliation+openapi/openapi.yaml) /src/ee/reconciliation/
    COPY (orchestration+openapi/openapi.yaml) /src/ee/orchestration/

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

pre-commit: # Generate the final spec and run all the pre-commit hooks
    LOCALLY
    BUILD +build-final-spec