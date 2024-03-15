VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly

IMPORT $core AS core
IMPORT ../.. AS stack
IMPORT ../../releases AS releases
IMPORT .. AS components

FROM core+base-image

build-image:
    RUN echo "not implemented"

deploy:
    RUN echo "not implemented"

tests:
    RUN echo "not implemented"

sources:
    WORKDIR src
    COPY --pass-args (releases+sdk-generate/go) /src/releases/sdks/go
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    WORKDIR /src/components/fctl
    COPY go.* .
    COPY --dir cmd pkg membershipclient .
    COPY main.go .
    SAVE ARTIFACT /src

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/components/fctl
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT main.go AS LOCAL main.go

pre-commit:
    WAIT
        BUILD --pass-args +tidy
    END
    BUILD --pass-args +lint
    BUILD --pass-args +completions

completions:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/components/fctl
    RUN mkdir -p ./completions
    RUN go run main.go completion bash > "./completions/fctl.bash"
    RUN go run main.go completion zsh > "./completions/fctl.zsh"
    RUN go run main.go completion fish > "./completions/fctl.fish"
    SAVE ARTIFACT ./completions AS LOCAL completions

openapi:
    RUN echo "not implemented"

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/components/fctl
    DO --pass-args stack+GO_TIDY

generate-membership-client:
    FROM openapitools/openapi-generator-cli:v6.6.0
    WORKDIR /src
    COPY membership-swagger.yaml .
    RUN docker-entrypoint.sh generate \
        -i ./membership-swagger.yaml \
        -g go \
        -o ./membershipclient \
        --git-user-id=formancehq \
        --git-repo-id=fctl \
        -p packageVersion=latest \
        -p isGoSubmodule=true \
        -p packageName=membershipclient
    RUN rm -rf ./membershipclient/test
    SAVE ARTIFACT ./membershipclient AS LOCAL membershipclient