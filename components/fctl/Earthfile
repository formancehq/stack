VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.5.2
IMPORT $core AS core
IMPORT ../.. AS stack
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
    COPY --pass-args (stack+build-sdk/go --LANG=go) sdks/go
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    WORKDIR /src/components/fctl
    COPY go.* .
    COPY --dir cmd pkg membershipclient .
    COPY main.go .
    SAVE ARTIFACT /src

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args (components+tidy/go.* --components=fctl) .
    WORKDIR /src/components/fctl
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT main.go AS LOCAL main.go

pre-commit:
    RUN echo "not implemented"