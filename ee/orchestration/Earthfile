VERSION 0.8

IMPORT github.com/formancehq/earthly:tags/v0.14.0 AS core
IMPORT ../.. AS stack
IMPORT ../../releases AS releases
IMPORT .. AS ee

FROM core+base-image

sources:
    WORKDIR src
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    COPY --pass-args (releases+sdk-generate/go) /src/releases/sdks/go
    WORKDIR /src/ee/orchestration
    COPY go.* .
    COPY --dir pkg cmd internal .
    COPY main.go .
    SAVE ARTIFACT /src

generate:
    FROM core+builder-image
    RUN apk update && apk add openjdk11
    DO --pass-args core+GO_INSTALL --package=go.uber.org/mock/mockgen@latest
    COPY (+sources/*) /src
    WORKDIR /src/ee/orchestration
    DO --pass-args core+GO_GENERATE
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT cmd AS LOCAL cmd

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/orchestration
    ARG VERSION=latest
    DO --pass-args core+GO_COMPILE --VERSION=$VERSION

build-image:
    FROM core+final-image
    ENTRYPOINT ["/bin/orchestration"]
    CMD ["serve"]
    COPY (+compile/main) /bin/orchestration
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO core+SAVE_IMAGE --COMPONENT=orchestration --REPOSITORY=${REPOSITORY} --TAG=$tag

tests:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/orchestration
    WITH DOCKER --pull=postgres:15-alpine
        DO --pass-args core+GO_TESTS
    END

deploy:
    COPY (+sources/*) /src
    LET tag=$(tar cf - /src | sha1sum | awk '{print $1}')
    WAIT
        BUILD --pass-args +build-image --tag=$tag
    END
    FROM --pass-args core+vcluster-deployer-image
    RUN kubectl patch Versions.formance.com default -p "{\"spec\":{\"orchestration\": \"${tag}\"}}" --type=merge

deploy-staging:
    BUILD --pass-args stack+deployer-module --MODULE=orchestration

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/ee/orchestration
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT main.go AS LOCAL main.go

pre-commit:
    WAIT
      BUILD --pass-args +tidy
    END
    BUILD --pass-args +lint

openapi:
    FROM node:20-alpine
    RUN apk update && apk add yq
    RUN npm install -g openapi-merge-cli
    WORKDIR /src/ee/orchestration
    COPY --dir openapi openapi
    RUN openapi-merge-cli --config ./openapi/openapi-merge.json
    RUN yq -oy ./openapi.json > openapi.yaml
    SAVE ARTIFACT ./openapi.yaml AS LOCAL ./openapi.yaml

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/ee/orchestration
    DO --pass-args stack+GO_TIDY

release:
    BUILD --pass-args stack+goreleaser --path=ee/orchestration