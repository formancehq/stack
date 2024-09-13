VERSION 0.8

IMPORT github.com/formancehq/earthly:tags/v0.14.0 AS core

IMPORT ../.. AS stack
IMPORT ../../helm/libs AS helm-libs 
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

deploy-staging:
    FROM --pass-args core+base-argocd 

    ARG --required TAG

    ARG APPLICATION=staging-eu-west-1-hosting-regions
    LET SERVER=argocd.internal.formance.cloud
    
    RUN --secret AUTH_TOKEN \
        argocd app set $APPLICATION \ 
        --parameter agent.image.tag=$TAG \
        --auth-token=$AUTH_TOKEN --server=$SERVER --grpc-web

    BUILD --pass-args core+deploy-staging

pre-commit:
    BUILD --pass-args +helm-validate
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

generate:
    FROM core+builder-image
    DO --pass-args core+GO_INSTALL --package=go.uber.org/mock/mockgen@latest
    COPY (+sources/*) /src
    WORKDIR /src/ee/agent    
    RUN go generate -run mockgen ./...
    SAVE ARTIFACT internal AS LOCAL internal


grpc-generate:
    FROM core+grpc-base
    LET protoName=agent.proto
    COPY $protoName .
    DO core+GRPC_GEN --protoName=$protoName
    SAVE ARTIFACT generated AS LOCAL internal/generated

tests:
    FROM core+builder-image
    RUN apk update && apk add bash
    DO --pass-args core+GO_INSTALL --package=sigs.k8s.io/controller-runtime/tools/setup-envtest@v0.0.0-20240320141353-395cfc7486e6
    ENV ENVTEST_VERSION 1.28.0
    RUN setup-envtest use $ENVTEST_VERSION -p path
    ENV KUBEBUILDER_ASSETS /root/.local/share/kubebuilder-envtest/k8s/$ENVTEST_VERSION-linux-$(go env GOHOSTARCH)
    DO --pass-args core+GO_INSTALL --package=github.com/onsi/ginkgo/v2/ginkgo@v2.14.0
    COPY --pass-args +sources/* /src
    COPY --pass-args ../../components/operator+manifests/config /src/components/operator/config
    WORKDIR /src/ee/agent
    COPY tests tests
    COPY internal internal
    ARG GOPROXY
    ARG focus
    
    RUN --mount=type=cache,id=gomod,target=$GOPATH/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        go test ./internal/...

    RUN --mount=type=cache,id=gomod,target=$GOPATH/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        ginkgo --focus=$focus -p ./tests/...
    
helm-validate:
    FROM core+helm-base
    WORKDIR /src

    COPY (helm-libs+sources/*) helm/libs/
    COPY --dir helm ee/agent/

    WORKDIR /src/ee/agent/helm
    RUN helm dependencies update
    DO --pass-args core+HELM_VALIDATE
    SAVE ARTIFACT /src/ee/agent/helm AS LOCAL helm

helm-package:
    FROM +helm-validate
    RUN helm package .
    SAVE ARTIFACT /src

    SAVE ARTIFACT . AS LOCAL helm

release:
    BUILD --pass-args stack+goreleaser --path=ee/agent