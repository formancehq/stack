VERSION --pass-args --arg-scope-and-set 0.7

ARG core=github.com/formancehq/earthly:v0.5.2
IMPORT $core AS core
IMPORT ../.. AS stack

FROM core+base-image

sources:
    FROM core+builder-image
    WORKDIR /src
    COPY (stack+sources/out --LOCATION=components/search) components/search
    #COPY (stack+sources/out --LOCATION=components/search/benthos) components/search:benthos
    COPY (stack+sources/out --LOCATION=libs/go-libs) libs/go-libs
    WORKDIR components/operator
    COPY --dir apis internal pkg .
    COPY main.go go.* .
    SAVE ARTIFACT /src

compile:
    FROM core+builder-image
    COPY (+sources/*) /src
    COPY --pass-args (+generate/*) /src/components/operator
    WORKDIR /src/components/operator
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
    DO --pass-args core+GO_INSTALL --package=sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.2

manifests:
    FROM --pass-args +controller-gen
    COPY (+sources/*) /src
    WORKDIR /src/components/operator
    COPY --dir config .
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        controller-gen rbac:roleName=manager-role crd webhook paths="./apis/..." output:crd:artifacts:config=config/crd/bases
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        controller-gen rbac:roleName=manager-role crd webhook paths="./internal/controllers/..." output:crd:artifacts:config=config/crd/bases

    SAVE ARTIFACT config AS LOCAL config

generate:
    FROM --pass-args +controller-gen
    COPY +sources/* /src
    WORKDIR /src/components/operator
    COPY --dir hack .
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./apis/..."
    RUN --mount=type=cache,id=gomod,target=${GOPATH}/pkg/mod \
        --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
        controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./internal/..."
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT apis AS LOCAL apis

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
    RUN rm -f helm/templates/gen/apiextensions.k8s.io_v1_customresourcedefinition_*.components.formance.com.yaml

    SAVE ARTIFACT helm AS LOCAL helm

deploy:
    COPY (+sources/*) /src
    LET tag=$(tar cf - /src | sha1sum | awk '{print $1}')
    WAIT
        BUILD --pass-args +build-image --tag=$tag
    END
    FROM --pass-args core+vcluster-deployer-image
    COPY --pass-args (+helm-update/helm) helm
    RUN rm -rf helm/templates/gen/admissionregistration.k8s.io_v1_mutatingwebhookconfiguration_formance-system-mutating-webhook-configuration.yaml
    WORKDIR helm
    RUN helm upgrade --namespace formance-system --install formance-operator \
        --wait \
        --create-namespace \
        --create-namespace \
        --set operator.disableWebhooks=true \
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
    COPY --pass-args (stack+tidy/go.* --component=operator) .
    WORKDIR /src/components/operator
    DO --pass-args stack+GO_LINT
    SAVE ARTIFACT apis AS LOCAL apis
    SAVE ARTIFACT internal AS LOCAL internal
    SAVE ARTIFACT pkg AS LOCAL pkg
    SAVE ARTIFACT main.go AS LOCAL main.go

tests:
    FROM core+builder-image
    RUN apk update && apk add bash
    DO --pass-args core+GO_INSTALL --package=sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
    ENV ENVTEST_VERSION 1.24.1
    RUN setup-envtest use $ENVTEST_VERSION -p path
    ENV KUBEBUILDER_ASSETS /root/.local/share/kubebuilder-envtest/k8s/$ENVTEST_VERSION-linux-$(go env GOHOSTARCH)
    DO --pass-args core+GO_INSTALL --package=github.com/onsi/ginkgo/v2/ginkgo@v2.6.0
    COPY (+sources/*) /src
    COPY --pass-args (+manifests/config) /src/components/operator/config
    COPY --pass-args (+generate/internal) /src/components/operator/internal
    COPY --pass-args (+generate/apis) /src/components/operator/apis
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

pre-commit:
    BUILD --pass-args +generate
    BUILD --pass-args +manifests
    BUILD --pass-args +helm-update