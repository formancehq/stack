VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.5.2
IMPORT $core AS core
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
    RUN kubectl patch Versions default -p "{\"spec\":{\"stargate\": \"${tag}\"}}" --type=merge
    COPY .earthly/stargate-values.yaml stargate-values.yaml
    COPY helm helm
    ARG --required user
    RUN --secret tld helm upgrade --namespace formance-system --create-namespace --install formance-stargate ./helm \
        -f stargate-values.yaml \
        --set config.auth_issuer_url=https://$user.$tld/api \
        --set image.tag=$tag
    COPY .earthly/ingress ingress-chart
    RUN --secret tld helm upgrade --install stargate-ingress ./ingress-chart \
        --namespace formance-system \
        --set user=$user \
        --set tld=$tld

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
    BUILD --pass-args +tidy
    BUILD --pass-args +lint

openapi:
    RUN echo "not implemented"

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/ee/stargate
    DO --pass-args stack+GO_TIDY