VERSION --arg-scope-and-set --pass-args 0.7

ARG core=github.com/formancehq/earthly:v0.5.2
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
    RUN mkdir build
    COPY libs/clients/base.yaml .
    COPY libs/clients/package.* .
    RUN npm install
    WORKDIR /src/components
    FOR c IN ledger payments
        COPY components/$c/openapi*.yaml $c/
    END
    WORKDIR /src/ee
    FOR c IN auth webhooks search wallets orchestration reconciliation
        COPY ee/$c/openapi*.yaml $c/
    END
    WORKDIR /src/libs/clients
    COPY libs/clients/openapi-merge.json .
    RUN npm run build
    LET VERSION=$(date +%Y%m%d)
    RUN jq '.info.version = "v1.0.${VERSION}"' build/generate.json > build/generate-with-version.json
    SAVE ARTIFACT build/generate-with-version.json
    SAVE ARTIFACT build/generate-with-version.json AS LOCAL libs/clients/build/generate.json

build-sdk:
    BUILD --pass-args +build-final-spec # Force output of the final spec
    FROM core+base-image
    WORKDIR /src
    RUN apk update && apk add yq
    COPY (+speakeasy/speakeasy) /bin/speakeasy
    COPY (+build-final-spec/generate-with-version.json) final-spec.json
    COPY --dir libs/clients/go ./sdks/go
    RUN --secret SPEAKEASY_API_KEY speakeasy generate sdk -s ./final-spec.json -o ./sdks/go -l go
    RUN rm -rf ./libs/clients/go
    SAVE ARTIFACT sdks/go AS LOCAL ./libs/clients/go
    SAVE ARTIFACT sdks/go

openapi:
    LOCALLY
    FOR component IN $(cd ./components && ls -d */)
      BUILD ./components/$component+openapi
    END
    FOR component IN $(cd ./ee && ls -d */)
      BUILD ./ee/$component+openapi
    END

goreleaser:
    FROM core+builder-image
    ARG --required components
    ARG --required type
    WAIT
      BUILD +openapi
    END
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

tidy-all:
    LOCALLY
    WAIT
      BUILD --pass-args ./libs/go-libs+tidy
    END
    #BUILD --pass-args ./tests/integration+tidy
    FOR component IN $(cd ./components && ls -d */)
      BUILD --pass-args ./components+tidy --components=$component
    END
    FOR component IN $(cd ./ee && ls -d */)
      BUILD --pass-args ./ee+tidy --components=$component
    END

tests-all:
    LOCALLY
    WAIT
      BUILD --pass-args +tidy-all
    END
    FOR component IN $(cd ./components && ls -d */)
      BUILD --pass-args ./components/${component}+tests
    END
    FOR component IN $(cd ./ee && ls -d */)
      BUILD --pass-args ./ee/${component}+tests
    END

lint-all:
    LOCALLY
    WAIT
      BUILD --pass-args +tidy-all
    END
    FOR component IN $(cd ./components && ls -d */)
      BUILD --pass-args ./components/${component}+lint
    END
    FOR component IN $(cd ./ee && ls -d */)
      BUILD --pass-args ./ee/${component}+lint
    END

tests-integration:
    FROM core+base-image
    BUILD --pass-args ./tests/integration+tests

pre-commit:
    LOCALLY
    WAIT
      BUILD --pass-args +tidy-all
      BUILD +openapi
    END
    BUILD --pass-args +build-final-spec
    BUILD --pass-args +lint-all
    BUILD --pass-args +tests-all
    FOR component IN $(cd ./components && ls -d */)
      BUILD --pass-args ./components/${component}+pre-commit
    END
    FOR component IN $(cd ./ee && ls -d */)
      BUILD --pass-args ./ee/${component}+pre-commit
    END
    ARG skipIntegrationTests=0
    IF [ "$skipIntegrationTests" = "0" ]
        BUILD --pass-args +tests-integration
    END

pr:
    LOCALLY
    WAIT
      BUILD --pass-args +tidy-all
    END
    BUILD --pass-args +lint-all
    BUILD --pass-args +tests-all
    BUILD --pass-args +tests-integration

INCLUDE_GO_LIBS:
    COMMAND
    ARG --required LOCATION
    COPY (+sources/out --LOCATION=libs/go-libs) ${LOCATION}

GO_LINT:
    COMMAND
    COPY (+sources/out --LOCATION=.golangci.yml) .golangci.yml
    ARG GOPROXY
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        --mount=type=cache,id=golangci,target=/root/.cache/golangci-lint \
        golangci-lint run --fix ./...