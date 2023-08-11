# Formance-operator

This operator is in charge of deploying a full or partial Formance OSS Stack.
It aims to simplify deployment and releases management of different parts of the Formance ecosystem.

## Getting Started

You'll need a Kubernetes cluster to run against.
Scripts of this repository are using [K3D](https://k3d.io/v5.4.6/). You have to install it.
Also, we use [Garden](https://docs.garden.io/) for management.

### Running on the cluster

1. Create the cluster:

```sh
garden create-cluster
```

2. Build the operator image yourself or skip and deploy:

> Add an entry for `k3d-registry.host.k3d.internal` inside /etc/hosts file, pointing to 127.0.0.1.

```sh
    1. BUILD: `make docker-build`
    2. PUSH: `make docker-push
    3. BUILD Helm: `make helm-update`
```

3. Deploy:

```sh
garden deploy
```

4. Create a stack

```sh
kubectl apply -f garden/example-v1beta3.yaml
```

5. Stop the cluster

```sh
garden stop
```

6. Start the cluster

```sh
garden start
```

Add an entry for `host.k3d.internal` inside /etc/hosts file, pointing to 127.0.0.1.
Go to http://host.k3d.internal.
Login with admin@formance.com / password

### Push to local registry

In order to be able to pull and push the image in the internal-registry named `k3d-registry.host.k3d.internal` 
on fixed port `12345` defined in `garden/k3d.yaml` 


Add an entry for `k3d-registry.host.k3d.internal` inside /etc/hosts file, pointing to 127.0.0.1.

Then in order to build and publish your image
    1. BUILD: `make docker-build`
    2. PUSH: `make docker-push`
    3. BUILD CRD: `make kustomize`
    4. DEPLOY HELM:`make helm-local-install`
    5. REDEPLOY HELM: `make helm-local-upgrade`
   


### Push to local registry

In order to be able to pull and push the image in the internal-registry named `k3d-registry.host.k3d.internal` 
on fixed port `12345` defined in `garden/k3d.yaml` 


Add an entry for `k3d-registry.host.k3d.internal` inside /etc/hosts file, pointing to 127.0.0.1.

Then in order to build and publish your image
    1. BUILD: `make docker-build`
    2. PUSH: `make docker-push`
    3. BUILD Helm: `make helm-update`
At this step you can use `garden deploy`
    1. DEPLOY Helm:`make helm-local-install`
    2. REDEPLOY Helm: `make helm-local-upgrade`


### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster


#### Create a stack

```sh
kubectl apply -f garden/example-v1beta3.yaml
### Tests

Run command :
```sh
make test
```

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.