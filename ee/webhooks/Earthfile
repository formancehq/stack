VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly

IMPORT $core AS core
IMPORT ../.. AS stack
IMPORT .. AS ee

FROM core+base-image

sources:
    WORKDIR src
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    WORKDIR /src/ee/webhooks
    COPY go.* .
    COPY --dir pkg cmd .
    COPY main.go .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/webhooks
    ARG VERSION=latest
    DO --pass-args core+GO_COMPILE --VERSION=$VERSION

build-image:
    FROM core+final-image
    ENTRYPOINT ["/bin/webhooks"]
    CMD ["serve"]
    COPY (+compile/main) /bin/webhooks
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO core+SAVE_IMAGE --COMPONENT=webhooks --REPOSITORY=${REPOSITORY} --TAG=$tag

tests:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/webhooks
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
    RUN kubectl patch Versions.formance.com default -p "{\"spec\":{\"webhooks\": \"${tag}\"}}" --type=merge

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/ee/webhooks
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT main.go AS LOCAL main.go

pre-commit:
    BUILD --pass-args +tidy
    BUILD --pass-args +lint

openapi:
    COPY ./openapi.yaml .
    SAVE ARTIFACT ./openapi.yaml

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/ee/webhooks
    DO --pass-args stack+GO_TIDY