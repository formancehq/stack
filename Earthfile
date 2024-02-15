VERSION --arg-scope-and-set --pass-args --use-function-keyword 0.7

ARG core=github.com/formancehq/earthly:v0.9.0
IMPORT $core AS core

sources:
    FROM core+base-image
    ARG --required LOCATION
    COPY ${LOCATION} out
    SAVE ARTIFACT out

speakeasy:
    FROM core+base-image
    RUN apk update && apk add yarn jq unzip curl
    ARG VERSION=v1.147.0
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
    ARG --required components
    ARG --required type
    COPY . /src
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
    FOR component IN $(cd ./tools && ls -d */)
        BUILD --pass-args +goreleaser --type=components --components=$component --mode=ci
    END
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args +goreleaser --type=components --components=$component --mode=ci
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args +goreleaser --type=ee --components=$component --mode=ci
    END

build-all:
    LOCALLY
    FOR component IN $(cd ./tools && ls -d */)
        BUILD --pass-args ./tools/${component}+build-image
    END
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
    FOR component IN $(cd ./tools && ls -d */)
        BUILD --pass-args ./tools/$component+deploy
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
    BUILD --pass-args ./releases+sdk-generate
    FOR component IN $(cd ./tools && ls -d */)
        BUILD --pass-args ./tools/${component}+pre-commit
    END
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args ./components/${component}+pre-commit
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/${component}+pre-commit
    END
    BUILD --pass-args ./tests/integration+pre-commit

tidy: # Run tidy on all the components
    LOCALLY
    FOR component IN $(cd ./tools && ls -d */)
            BUILD --pass-args ./tools/${component}+tidy
        END
    FOR component IN $(cd ./components && ls -d */)
        BUILD --pass-args ./components/${component}+tidy
    END
    FOR component IN $(cd ./ee && ls -d */)
        BUILD --pass-args ./ee/${component}+tidy
    END
    BUILD --pass-args ./tests/integration+tidy

tests:
    LOCALLY
    BUILD --pass-args +tests-all
    BUILD --pass-args +tests-integration

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
