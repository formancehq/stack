VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.9.0
IMPORT $core AS core

IMPORT ../.. AS stack
IMPORT .. AS ee

FROM core+base-image

sources:
    WORKDIR src
    DO stack+INCLUDE_GO_LIBS --LOCATION libs/go-libs
    COPY (../../components/operator+sources/*) /src
    WORKDIR /src/ee/agent
    COPY go.* .
    COPY --dir cmd internal tests .
    COPY main.go .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    WORKDIR /src/ee/agent
    ARG VERSION=latest
    DO --pass-args core+GO_COMPILE --VERSION=$VERSION

build-image:
    FROM core+final-image
    ENTRYPOINT ["/bin/agent"]
    COPY (+compile/main) /bin/agent
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO core+SAVE_IMAGE --COMPONENT=agent --REPOSITORY=${REPOSITORY} --TAG=$tag

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/ee/agent
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT cmd AS LOCAL cmd
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT main.go AS LOCAL main.go

deploy:
    COPY (+sources/*) /src
    LET tag=$(tar cf - /src | sha1sum | awk '{print $1}')
    WAIT
        BUILD --pass-args +build-image --tag=$tag
    END
    FROM --pass-args core+vcluster-deployer-image
    COPY helm helm
    COPY .earthly .earthly
    ARG --required user
    RUN --secret tld helm upgrade --namespace formance-system \
        --create-namespace \
        --install \
        -f .earthly/values.yaml \
        --set image.tag=$tag \
        --set agent.baseUrl=https://$user.$tld \
        --set server.address=$user.$tld:443 \
        formance-membership-agent ./helm

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
    WORKDIR /src/ee/agent
    DO --pass-args stack+GO_TIDY

grpc-generate:
    FROM core+builder-image
    RUN apk add --no-cache protobuf git protobuf-dev && \
        go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    WORKDIR /src
    RUN mkdir generated
    COPY server.proto .
    RUN protoc --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative server.proto
    SAVE ARTIFACT generated AS LOCAL internal/generated

tests:
    FROM core+builder-image
    RUN apk update && apk add bash
    DO --pass-args core+GO_INSTALL --package=sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
    ENV ENVTEST_VERSION 1.28.0
    RUN setup-envtest use $ENVTEST_VERSION -p path
    ENV KUBEBUILDER_ASSETS /root/.local/share/kubebuilder-envtest/k8s/$ENVTEST_VERSION-linux-$(go env GOHOSTARCH)
    DO --pass-args core+GO_INSTALL --package=github.com/onsi/ginkgo/v2/ginkgo@v2.14.0
    COPY --pass-args +sources/* /src
    COPY --pass-args ../../components/operator+manifests/config /src/components/operator/config
    WORKDIR /src/ee/agent
    COPY tests tests
    ARG GOPROXY
    ARG focus
    RUN --mount=type=cache,id=gomod,target=$GOPATH/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        ginkgo --focus=$focus ./tests/...