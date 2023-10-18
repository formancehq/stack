VERSION --arg-scope-and-set --pass-args 0.7

ARG core=github.com/formancehq/earthly:v0.3.0
IMPORT $core AS core

sources:
    FROM core+base-image
    ARG --required LOCATION
    COPY ${LOCATION} out
    SAVE ARTIFACT out

speakeasy:
    FROM core+base-image
    RUN apk update && apk add yarn jq unzip curl
    ARG VERSION=v1.78.4
    ARG TARGETARCH
    RUN curl -fsSL https://github.com/speakeasy-api/speakeasy/releases/download/${VERSION}/speakeasy_linux_$TARGETARCH.zip -o /tmp/speakeasy_linux_$TARGETARCH.zip
    RUN unzip /tmp/speakeasy_linux_$TARGETARCH.zip speakeasy
    RUN chmod +x speakeasy
    SAVE ARTIFACT speakeasy

build-final-spec:
    FROM core+base-image
    RUN apk update && apk add yarn nodejs npm jq
    WORKDIR /src/openapi
    RUN mkdir build
    COPY openapi/base.yaml .
    COPY openapi/package.* .
    RUN npm install
    WORKDIR /src/components
    FOR c IN auth ledger payments webhooks search wallets orchestration
        COPY components/$c/openapi.yaml $c/openapi.yaml
    END
    WORKDIR /src/openapi
    COPY openapi/openapi-merge.json .
    RUN npm run build
    ENV VERSION v1.0.$(date +%Y%m%d)
    RUN jq '.info.version = "${VERSION}"' build/generate.json > build/generate-with-version.json
    SAVE ARTIFACT build/generate-with-version.json
    SAVE ARTIFACT build/generate-with-version.json AS LOCAL build/generate.json

build-sdk:
    FROM core+base-image
    WORKDIR /src
    RUN apk update && apk add yq
    COPY (+speakeasy/speakeasy) /bin/speakeasy
    COPY (+build-final-spec/generate-with-version.json) final-spec.json
    ARG LANG=go
    COPY --dir openapi/templates/${LANG} sdks/${LANG}
    ENV VERSION v1.0.$(date +%Y%m%d)
    RUN yq e -i ".${LANG}.version = \"${VERSION}\"" sdks/${LANG}/gen.yaml
    RUN --secret SPEAKEASY_API_KEY speakeasy generate sdk -s ./final-spec.json -o ./sdks/${LANG} -l ${LANG}
    ARG export
    IF [ "$export" == "1" ]
        SAVE ARTIFACT sdks/${LANG} AS LOCAL ./sdks/${LANG}
    ELSE
        SAVE ARTIFACT sdks/${LANG}
    END

build-all-sdk:
    LOCALLY
      FOR lang IN $(ls openapi/templates)
          BUILD +build-sdk --LANG=${lang}
      END

goreleaser:
    FROM core+builder-image
    ARG --required component
    ARG mode=local
    COPY --pass-args (./components/$component+sources/*) /src
    COPY ./components/$component/.goreleaser.yml /src/components/$component/.goreleaser.yml
    COPY --if-exists ./components/$component/scripts/completions.sh /src/components/$component/scripts/completions.sh
    COPY --if-exists ./components/$component/build.Dockerfile /src/components/$component/build.Dockerfile
    COPY ./.goreleaser.default.yaml /src/.goreleaser.default.yaml
    COPY .git /src/.git
    WORKDIR /src/components/$component
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
      RUN --secret GORELEASER_KEY --secret GITHUB_TOKEN --secret SPEAKEASY_API_KEY --secret FURY_TOKEN --secret SEGMENT_WRITE_KEY goreleaser release -f .goreleaser.yml $buildArgs
    END

all-local-goreleaser:
    LOCALLY
    FOR component IN $(ls components)
        BUILD --pass-args +goreleaser --component=$component --mode=local
    END

build-images:
    LOCALLY
    FOR component IN $(ls components)
        BUILD --pass-args ./components/$component+build-image
    END

deploy:
    FROM core+base-image
    ARG --required component
    BUILD --pass-args ./components/$component+deploy

deploy-all:
    FROM core+base-image
    WAIT
        BUILD --pass-args +deploy --component=operator
    END
    COPY components components
    FOR component IN $(ls components)
        IF [ "$component" != "operator" ]
            BUILD --pass-args +deploy --component=$component
        END
    END

lint-all:
    LOCALLY
    BUILD --pass-args +tidy-all
    FOR component IN $(ls components)
        BUILD --pass-args ./components/${component}+lint
    END

tidy:
    FROM core+builder-image
    ARG --required component
    COPY --pass-args (./components/$component+sources/*) /src
    ARG GOPROXY
    WORKDIR /src/components/$component
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        go mod tidy
    SAVE ARTIFACT go.* AS LOCAL components/$component/

tidy-all:
    LOCALLY
    FOR component IN $(ls components)
        BUILD --pass-args +tidy --component=$component
    END

tests:
    LOCALLY
    ARG --required component
    BUILD --pass-args ./components/${component}+tests

tests-all:
    LOCALLY
    BUILD --pass-args +tidy-all
    FOR component IN $(ls components)
        BUILD --pass-args +tests --component=$component
    END

integration-tests:
    FROM core+base-image
    BUILD --pass-args ./tests/integration+tests

go-dep-updates:
    FROM core+builder-image
    ARG --required component
    COPY --pass-args ./components/$component+sources/* /src
    WORKDIR /src/components/$component
    ARG GOPROXY
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        go get -u=patch
    SAVE ARTIFACT go.mod AS LOCAL components/$component/go.mod
    SAVE ARTIFACT go.sum AS LOCAL components/$component/go.sum

go-dep-updates-all:
    LOCALLY
    FOR component IN $(ls components)
        BUILD --pass-args +go-dep-updates --component=$component
    END

pre-commit:
    LOCALLY
    BUILD --pass-args +tidy-all
    BUILD --pass-args +lint-all
    BUILD --pass-args +tests-all
    FOR component IN $(ls components)
        BUILD --pass-args ./components/$component+pre-commit
    END
    BUILD --pass-args +integration-tests

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