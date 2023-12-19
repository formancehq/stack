VERSION --arg-scope-and-set --pass-args --use-function-keyword 0.7

ARG core=github.com/formancehq/earthly:v0.6.0
IMPORT $core AS core

sources:
    FROM core+base-image
    ARG --required LOCATION
    COPY ${LOCATION} out
    SAVE ARTIFACT out

speakeasy:
    FROM core+base-image
    RUN apk update && apk add yarn jq unzip curl
    ARG VERSION=v1.109.1
    ARG TARGETARCH
    RUN curl -fsSL https://github.com/speakeasy-api/speakeasy/releases/download/${VERSION}/speakeasy_linux_$TARGETARCH.zip -o /tmp/speakeasy_linux_$TARGETARCH.zip
    RUN unzip /tmp/speakeasy_linux_$TARGETARCH.zip speakeasy
    RUN chmod +x speakeasy
    SAVE ARTIFACT speakeasy

build-final-spec:
    FROM core+base-image
    RUN apk update && apk add yarn nodejs npm jq
    WORKDIR /src/libs/clients
    COPY libs/clients/package.* .
    RUN npm install
    WORKDIR /src/components
    FOR c IN payments ledger
        COPY (./components/$c+openapi/openapi.yaml) /src/components/$c/
    END
    WORKDIR /src/ee
    FOR c IN auth webhooks search wallets reconciliation orchestration
        COPY (./ee/$c+openapi/openapi.yaml) /src/ee/$c/
    END

    WORKDIR /src/libs/clients
    COPY libs/clients/base.yaml .
    COPY libs/clients/openapi-overlay.json .
    COPY libs/clients/openapi-merge.json .
    RUN mkdir ./build
    RUN npm run build
    RUN jq -s '.[0] * .[1]' build/generate.json openapi-overlay.json > build/final.json
    RUN sed -i 's/SDK_VERSION/INTERNAL/g' build/final.json
    SAVE ARTIFACT build/final.json AS LOCAL libs/clients/build/generate.json

build-sdk:
    BUILD --pass-args +build-final-spec # Force output of the final spec
    FROM core+base-image
    WORKDIR /src
    RUN apk update && apk add yq
    COPY (+speakeasy/speakeasy) /bin/speakeasy
    COPY (+build-final-spec/final.json) final-spec.json
    COPY --dir libs/clients/go ./sdks/go
    IF [ "SPEAKEASY_API_KEY" != "" ]
        RUN --secret SPEAKEASY_API_KEY speakeasy generate sdk -s ./final-spec.json -o ./sdks/go -l go
    END
    RUN rm -rf ./libs/clients/go
    SAVE ARTIFACT sdks/go AS LOCAL ./libs/clients/go
    SAVE ARTIFACT sdks/go

openapi:
    FROM core+base-image
    COPY . /src
    WORKDIR /src
    FOR component IN $(cd ./components && ls -d */)
        COPY (./components/$component+openapi/src/components/$component) /src/components/$component
    END
    FOR component IN $(cd ./ee && ls -d */)
        COPY (./ee/$component+openapi/src/ee/$component) /src/ee/$component
    END
    RUN toto
    SAVE ARTIFACT /src

goreleaser:
    FROM core+builder-image
    ARG --required components
    ARG --required type
    COPY . /src
    COPY (+build-sdk/go --LANG=go) /src/libs/clients/go
    WORKDIR /src/$type/$components
    ARG mode=local
    LET buildArgs = --clean
    IF [ "$mode" = "local" ]
        SET buildArgs = --nightly --skip=publish --clean
    ELSE IF [ "$mode" = "ci" ]
        SET buildArgs = --nightly --clean
    END
    IF [ "$mode" != "local" ]
        WITH DOCKER
            RUN --secret GITHUB_TOKEN echo $GITHUB_TOKEN | docker login ghcr.io -u NumaryBot --password-stdin
        END
    END
    WITH DOCKER
       RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
           --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
           --secret GORELEASER_KEY \
           --secret GITHUB_TOKEN \
           --secret SPEAKEASY_API_KEY \
           --secret FURY_TOKEN \
           --secret SEGMENT_WRITE_KEY \
           goreleaser release -f .goreleaser.yml $buildArgs
    END

all-ci-goreleaser:
    LOCALLY
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args +goreleaser --type=components --components=$component --mode=ci
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args +goreleaser --type=ee --components=$component --mode=ci
    END

build-all:
    LOCALLY
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args ./components/${component}+build-image
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/${component}+build-image
    END

deploy-all:
    LOCALLY
    WAIT
        BUILD --pass-args ./components/+deploy --components=operator
    END
    FOR component IN $(cd ./components && ls -d */)
        IF [ "$component" != "operator" ]
            BUILD --pass-args ./components/+deploy --components=$component
        END
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/+deploy --components=$component
    END

tests-all:
    LOCALLY
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args ./components/${component}+tests
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/${component}+tests
    END

tests-integration:
    FROM core+base-image
    BUILD --pass-args ./tests/integration+tests

pre-commit: # Generate the final spec and run all the pre-commit hooks
    LOCALLY
    BUILD --pass-args +build-sdk
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args ./components/${component}+pre-commit
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/${component}+pre-commit
    END
    BUILD --pass-args ./tests/integration+pre-commit

pr:
    LOCALLY
    BUILD --pass-args +tests-all
    BUILD --pass-args +tests-integration

deploy-staging:
    FROM core+base-image
    RUN apk update && apk add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community helm git jq kubectl kustomize aws-cli
    COPY ./.kubeconfig /root/.kube/config
    RUN kubectl config use-context arn:aws:eks:eu-west-1:955332203423:cluster/staging-eu-west-1-hosting
    COPY . /src
    WORKDIR /src
    FOR COMPONENT IN $(cd ./components && ls -d */)
        RUN --secret AWS_ACCESS_KEY_ID \
            --secret AWS_SECRET_ACCESS_KEY \
            --secret AWS_SESSION_TOKEN \
             kubectl patch Versions default -p "{\"spec\":{\"${COMPONENT}\": \"${GITHUB_SHA}\"}}" --type=merge
    END
    FOR COMPONENT IN $(cd ./ee && ls -d */)
        RUN --secret AWS_ACCESS_KEY_ID \
            --secret AWS_SECRET_ACCESS_KEY \
            --secret AWS_SESSION_TOKEN \
             kubectl patch Versions default -p "{\"spec\":{\"${COMPONENT}\": \"${GITHUB_SHA}\"}}" --type=merge
    END

INCLUDE_GO_LIBS:
    FUNCTION
    ARG --required LOCATION
    COPY (+sources/out --LOCATION=libs/go-libs) ${LOCATION}

GO_LINT:
    FUNCTION
    COPY (+sources/out --LOCATION=.golangci.yml) .golangci.yml
    ARG GOPROXY
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        --mount=type=cache,id=golangci,target=/root/.cache/golangci-lint \
        golangci-lint run --fix ./...

GO_TIDY:
    FUNCTION
    ARG GOPROXY
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        go mod tidy
    SAVE ARTIFACT go.* AS LOCAL .