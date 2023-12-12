VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.5.2
IMPORT $core AS core
IMPORT ../.. AS stack
IMPORT .. AS ee

FROM core+base-image

sources:
    WORKDIR src
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    COPY --pass-args (stack+build-sdk/go) libs/clients/go
    WORKDIR /src/ee/reconciliation
    COPY go.* .
    COPY --dir cmd internal .
    COPY main.go .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/reconciliation
    ARG VERSION=latest
    DO --pass-args core+GO_COMPILE --VERSION=$VERSION

build-image:
    FROM core+final-image
    ENTRYPOINT ["/bin/reconciliation"]
    CMD ["serve"]
    COPY (+compile/main) /bin/reconciliation
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO core+SAVE_IMAGE --COMPONENT=reconciliation --REPOSITORY=${REPOSITORY} --TAG=$tag

tests:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/reconciliation
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
    RUN kubectl patch Versions default -p "{\"spec\":{\"reconciliation\": \"${tag}\"}}" --type=merge

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args (ee+tidy/go.* --components=reconciliation) .
    WORKDIR /src/ee/reconciliation
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT main.go AS LOCAL main.go

pre-commit:
    RUN echo "not implemented"

openapi:
    RUN echo "not implemented"