VERSION 0.8
PROJECT FormanceHQ/stack

IMPORT github.com/formancehq/earthly:tags/v0.15.0 AS core

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

    RUN mkdir -p /src/components/ledger/ /src/components/payments/ /src/ee/gateway/ /src/ee/auth/ /src/ee/search/ /src/ee/webhooks/ /src/ee/wallets/ /src/ee/reconciliation/ /src/ee/orchestration/
    
    RUN curl -L -o /src/components/ledger/openapi.yaml https://github.com/formancehq/ledger/releases/download/v2.3.0/openapi.yaml
    RUN curl -L -o /src/components/payments/openapi.yaml https://github.com/formancehq/payments/releases/download/v3.0.18/openapi.yaml
    RUN curl -L -o /src/ee/gateway/openapi.yaml https://github.com/formancehq/gateway/releases/download/v2.1.0/openapi.yaml
    RUN curl -L -o /src/ee/auth/openapi.yaml https://github.com/formancehq/auth/releases/download/v2.4.0/openapi.yaml
    RUN curl -L -o /src/ee/search/openapi.yaml https://github.com/formancehq/search/releases/download/v2.1.0/openapi.yaml
    RUN curl -L -o /src/ee/webhooks/openapi.yaml https://github.com/formancehq/webhooks/releases/download/v2.2.0/openapi.yaml
    RUN curl -L -o /src/ee/wallets/openapi.yaml https://github.com/formancehq/wallets/releases/download/v2.1.5/openapi.yaml
    RUN curl -L -o /src/ee/reconciliation/openapi.yaml https://github.com/formancehq/reconciliation/releases/download/v2.2.0/openapi.yaml
    RUN curl -L -o /src/ee/orchestration/openapi.yaml https://github.com/formancehq/flows/releases/download/v2.4.0/openapi.yaml

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