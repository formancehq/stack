# formance-operator

This operator is in charge of deploying a full or partial Formance OSS Stack.
It aims to simplify deployment and releases management of different parts of the Formance ecosystem.

## Getting Started

You’ll need a Kubernetes cluster to run against.
Scripts of this repository are using [K3D](https://k3d.io/v5.4.6/). You have to install it.
Also, we use [Garden](https://docs.garden.io/) for management.

### Running on the cluster
1. Create the cluster:

```sh
garden create-cluster
```

2. Deploy:

```sh
garden  deploy
```

This will automatically install all the stack.
When developing, use following command to update the operator code :
```sh
go run main.go --disable-webhooks
```

3. Create a stack

```sh
kubectl apply -f garden/example-v1beta3.yaml
```

Add an entry for `host.k3d.internal` inside /etc/hosts file, pointing to 127.0.0.1.
Go to http://host.k3d.internal.
Login with admin@formance.com / password

### Push to local registry

In order to be able to pull and push the image a registry named `k3d-registry.host.k3d.internal` 
on fixed port `12345` defined in `garden/k3d.yaml` 


Add an entry for `k3d-registry.host.k3d.internal` inside /etc/hosts file, pointing to 127.0.0.1.

Then in order to build and publish your image
    1. `make docker-build`
    2. `make docker-push`
    3. `make helm-local-install`
    4. `make helm-local-upgrade`

### Tests

Run command :
```sh
make test
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster

### Test It Out

You can install a full stack using the command:
```sh
kubectl apply -f example.yaml
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
