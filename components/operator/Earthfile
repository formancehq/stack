VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.9.0
IMPORT $core AS core
IMPORT ../.. AS stack
IMPORT .. AS components

FROM core+base-image

sources:
    FROM core+builder-image
    WORKDIR /src
    COPY (stack+sources/out --LOCATION=ee/search) ee/search
    COPY (stack+sources/out --LOCATION=libs/go-libs) libs/go-libs
    WORKDIR /src/components/operator
    COPY --dir api internal pkg cmd .
    COPY go.* .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args (+generate/*) /src/components/operator
    WORKDIR /src/components/operator/cmd
	DO --pass-args core+GO_COMPILE

build-image:
    FROM core+final-image
    ENTRYPOINT ["/usr/bin/operator"]
    COPY --pass-args (+compile/main) /usr/bin/operator
    ARG REPOSITORY=ghcr.io
    ARG tag=latest
    DO --pass-args core+SAVE_IMAGE --COMPONENT=operator --TAG=$tag

controller-gen:
    FROM core+builder-image
    DO --pass-args core+GO_INSTALL --package=sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0

manifests:
    FROM --pass-args +controller-gen
    COPY (+sources/*) /src
    WORKDIR /src/components/operator
    COPY --dir config .
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

    SAVE ARTIFACT config AS LOCAL config

generate:
    FROM --pass-args +controller-gen
    COPY +sources/* /src
    WORKDIR /src/components/operator
    COPY --dir hack .
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT api AS LOCAL api

helm-update:
    FROM core+builder-image
    RUN curl -s https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh | bash -s -- 4.5.4 /bin

    WORKDIR /src
    COPY --pass-args (+manifests/config) /src/config
    COPY --dir helm hack .

    RUN rm -f helm/templates/gen/*
    RUN kustomize build config/default --output helm/templates/gen
    RUN rm -f helm/templates/gen/v1_namespace*.yaml
    RUN rm -f helm/templates/gen/apps_v1_deployment_*.yaml

    SAVE ARTIFACT helm AS LOCAL helm

deploy:
    COPY (+sources/*) /src
    LET tag=$(tar cf - /src | sha1sum | awk '{print $1}')
    WAIT
        BUILD --pass-args +build-image --tag=$tag
    END
    FROM --pass-args core+vcluster-deployer-image
    COPY --pass-args (+helm-update/helm) helm
    WORKDIR helm
    RUN helm upgrade --namespace formance-system --install formance-operator \
        --wait \
        --create-namespace \
        --create-namespace \
        --set image.tag=$tag .
    WORKDIR /
    COPY .earthly .earthly
    WORKDIR .earthly
    RUN kubectl get versions default || kubectl apply -f k8s-versions.yaml
    ARG user
    RUN --secret tld helm upgrade --install operator-configuration ./configuration \
        --namespace formance-system \
        --set gateway.fallback=https://console.$user.$tld

lint:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args +tidy/go.* .
    WORKDIR /src/components/operator
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT api AS LOCAL api
    SAVE ARTIFACT internal AS LOCAL internal

tests:
    FROM core+builder-image
    RUN apk update && apk add bash
    DO --pass-args core+GO_INSTALL --package=sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
    ENV ENVTEST_VERSION 1.28.0
    RUN setup-envtest use $ENVTEST_VERSION -p path
    ENV KUBEBUILDER_ASSETS /root/.local/share/kubebuilder-envtest/k8s/$ENVTEST_VERSION-linux-$(go env GOHOSTARCH)
    DO --pass-args core+GO_INSTALL --package=github.com/onsi/ginkgo/v2/ginkgo@v2.14.0
    COPY (+sources/*) /src
    COPY --pass-args (+manifests/config) /src/components/operator/config
    COPY --pass-args (+generate/internal) /src/components/operator/internal
    COPY --pass-args (+generate/api) /src/components/operator/api
    WORKDIR /src/components/operator
    COPY --dir hack .
    ARG GOPROXY
    ARG updateTestData=0
    ENV UPDATE_TEST_DATA=$updateTestData
    ARG focus
    RUN --mount=type=cache,id=gomod,target=$GOPATH/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        ginkgo --focus=$focus -p ./...
    IF [ "$updateTestData" = "1" ]
        SAVE ARTIFACT internal/controllers/stack/testdata AS LOCAL internal/controllers/stack/testdata
    END

generate-docs:
    FROM core+builder-image
    COPY (+sources/*) /src
    RUN go install github.com/elastic/crd-ref-docs@latest
    WORKDIR /src/components/operator
    COPY docs.config.yaml .
    COPY --dir docs api .
    RUN crd-ref-docs --source-path=api/formance.com/v1beta1 --renderer=markdown --output-path=./docs/crd.md --config=./docs.config.yaml
    SAVE ARTIFACT docs/crd.md AS LOCAL docs/crd.md

pre-commit:
    WAIT
      BUILD --pass-args +tidy
    END
    BUILD --pass-args +lint
    BUILD --pass-args +generate
    BUILD --pass-args +manifests
    BUILD --pass-args +helm-update
    BUILD +generate-docs

openapi:
    RUN echo "not implemented"

tidy:
    FROM core+builder-image
    COPY --pass-args (+sources/src) /src
    WORKDIR /src/components/operator
    DO --pass-args stack+GO_TIDY