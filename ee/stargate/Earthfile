VERSION 0.8

IMPORT github.com/formancehq/earthly:tags/v0.15.0 AS core
IMPORT ../.. AS stack
IMPORT .. AS ee

FROM core+base-image

sources:
    WORKDIR src
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    WORKDIR /src/ee/stargate
    COPY go.* .
    COPY --dir cmd internal .
    COPY main.go .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/stargate
    ARG VERSION=latest
    DO --pass-args core+GO_COMPILE --VERSION=$VERSION

build-image:
    FROM core+final-image
    ENTRYPOINT ["/bin/stargate"]
    CMD ["client"]
    COPY (+compile/main) /bin/stargate
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO core+SAVE_IMAGE --COMPONENT=stargate --REPOSITORY=${REPOSITORY} --TAG=$tag

tests:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/stargate
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
    RUN kubectl patch Versions.formance.com default -p "{\"spec\":{\"stargate\": \"${tag}\"}}" --type=merge

deploy-staging:
    BUILD --pass-args stack+deployer-module --MODULE=stargate

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/ee/stargate
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT main.go AS LOCAL main.go

pre-commit:
    WAIT
        BUILD --pass-args +tidy
    END
    BUILD --pass-args +lint

openapi:
    RUN echo "not implemented"

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/ee/stargate
    DO --pass-args stack+GO_TIDY

grpc-generate:
    FROM core+grpc-base
    LET protoName=stargate.proto
    COPY $protoName .
    DO core+GRPC_GEN --protoName=$protoName
    SAVE ARTIFACT generated AS LOCAL internal/generated

release:
    BUILD --pass-args stack+goreleaser --path=ee/stargate