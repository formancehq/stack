VERSION 0.8
PROJECT FormanceHQ/stack

IMPORT github.com/formancehq/earthly:tags/v0.12.0 AS core

sources:
    FROM core+base-image
    ARG --required LOCATION
    COPY ${LOCATION} out
    SAVE ARTIFACT out

speakeasy:
    FROM core+base-image
    RUN apk update && apk add yarn jq unzip curl
    ARG VERSION=v1.346.0
    ARG TARGETARCH
    RUN curl -fsSL https://github.com/speakeasy-api/speakeasy/releases/download/${VERSION}/speakeasy_linux_$TARGETARCH.zip -o /tmp/speakeasy_linux_$TARGETARCH.zip
    RUN unzip /tmp/speakeasy_linux_$TARGETARCH.zip speakeasy
    RUN chmod +x speakeasy
    SAVE ARTIFACT speakeasy

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

    FOR c IN payments ledger
        COPY (./components/$c+openapi/openapi.yaml) /src/components/$c/
    END
    FOR c IN auth webhooks search wallets reconciliation orchestration
        COPY (./ee/$c+openapi/openapi.yaml) /src/ee/$c/
    END

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
    SAVE ARTIFACT /src

goreleaser:
    FROM core+builder-image
    ARG --required path
    COPY . /src
    WORKDIR /src/$path
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
            goreleaser release -f .goreleaser.yml $buildArgs
    END

all-ci-goreleaser:
    LOCALLY
    FOR service IN $(cd ./components && ls -d */)
        BUILD --pass-args ./components/$service+release --mode=ci
    END
    FOR service IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/$service+release --mode=ci
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

deployer-module:
    FROM --pass-args core+base-image
    ARG --required MODULE
    ARG --required TAG
    
    LET ARGS="--parameter=versions.files.default.$MODULE=$TAG"
    FROM --pass-args core+deployer-module --ARGS=$ARGS --TAG=$TAG

staging-application-set:
    LOCALLY
    ARG TAG=latest
    LET PARAMETERS=""

    WAIT
        FOR component IN $(cd ./components && ls -d */ | sed 's/.$//')
            IF [ "$component" != "operator"  ]  && [ "$component" != "fctl" ]
                SET PARAMETERS="$PARAMETERS --parameter versions.files.default.$component=$TAG"
            END
        END
        
        FOR component IN $(cd ./ee && ls -d */ | sed 's/.$//')
            IF [ "$component" != "agent"  ]
                SET PARAMETERS="$PARAMETERS --parameter versions.files.default.$component=$TAG"
            END
        END

        SET PARAMETERS="$PARAMETERS --parameter agent.image.tag=$TAG"
        SET PARAMETERS="$PARAMETERS --parameter operator.image.tag=$TAG"
    END
    BUILD --pass-args core+application-set --ARGS=$PARAMETERS --WITH_SYNC=false
    

staging-application-sync:
    BUILD core+application-sync --APPLICATION=staging-eu-west-1-hosting-regions

tests:
    LOCALLY
    BUILD ./components+run --TARGET=tests
    BUILD ./ee+run --TARGET=tests

tests-integration:
    FROM core+base-image
    BUILD ./tests/integration+tests

pre-commit: # Generate the final spec and run all the pre-commit hooks
    LOCALLY
    BUILD ./releases+sdk-generate
    BUILD ./libs+run --TARGET=pre-commit
    BUILD ./components+run --TARGET=pre-commit
    BUILD ./ee+run --TARGET=pre-commit
    BUILD ./helm/+pre-commit
    BUILD ./tests/integration+pre-commit

tidy: # Run tidy on all the components
    LOCALLY
    BUILD ./components+run --TARGET=tidy
    BUILD ./ee+run --TARGET=tidy
    BUILD ./tests/integration+tidy

tests-all:
    LOCALLY
    BUILD +tests
    BUILD +tests-integration

helm-publish:
    LOCALLY
    BUILD ./helm/+publish
    BUILD ./components/operator+helm-publish

HELM_PUBLISH:
    FUNCTION
    ARG --required path
    WITH DOCKER
        RUN --secret GITHUB_TOKEN echo $GITHUB_TOKEN | docker login ghcr.io -u NumaryBot --password-stdin
    END
    WITH DOCKER
        RUN helm push ${path} oci://ghcr.io/formancehq/helm
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
