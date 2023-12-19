VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.6.0
IMPORT $core AS core
IMPORT ../.. AS stack
IMPORT .. AS ee

FROM core+base-image

sources:
    WORKDIR src
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    WORKDIR /src/ee/gateway
    COPY go.* .
    COPY --dir internal .
    COPY --dir pkg .
    COPY main.go Caddyfile .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/gateway
    ARG VERSION=latest
    DO --pass-args core+GO_COMPILE --VERSION=$VERSION

build-image:
    FROM core+final-image
    ENTRYPOINT ["/usr/bin/caddy"]
    CMD ["run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
    COPY Caddyfile /etc/caddy/Caddyfile
    COPY (+compile/main) /usr/bin/caddy
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO core+SAVE_IMAGE --COMPONENT=gateway --REPOSITORY=${REPOSITORY} --TAG=$tag

deploy:
    COPY (+sources/*) /src
    LET tag=$(tar cf - /src | sha1sum | awk '{print $1}')
    WAIT
        BUILD --pass-args +build-image --tag=$tag
    END
    FROM --pass-args core+vcluster-deployer-image
    RUN kubectl patch Versions default -p "{\"spec\":{\"gateway\": \"${tag}\"}}" --type=merge

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/ee/gateway
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT main.go AS LOCAL main.go

tests:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/gateway
    DO --pass-args core+GO_TESTS

pre-commit:
    BUILD --pass-args +tidy
    BUILD --pass-args +lint

openapi:
    RUN echo "not implemented"

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/ee/gateway
    DO --pass-args stack+GO_TIDY