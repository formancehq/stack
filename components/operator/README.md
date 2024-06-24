# operator

## Description

The operator allow to install formance components on a k8s cluster.

## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- k3d

### To Deploy on the cluster

**Create a dev cluster**

```sh
k3d cluster create internal
```

**Install CRDs**

```sh
make install
```

**Run operator**

```sh
make run
```

**Install a PostgresServer**

```sh
helm install postgres oci://registry-1.docker.io/bitnamicharts/postgresql \
  --set auth.username=formance \
  --set auth.password=formance \
  --set auth.database=formance
```

**Create a Settings object for database connection**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: databaseconfiguration0
spec:
  stacks:
  - "*"
  key: postgres.*.uri
  value: postgresql://formance:formance@postgres-postgresql.default:5432?disableSSLMode=true
EOF
```

**Create a stack**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Stack
metadata:
  name: stack0
spec: {}
EOF
```

**Create a ledger**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Ledger
metadata:
  name: stack0-ledger
spec:
  stack: stack0
EOF
```

**Create a Gateway**

```sh
cat <<EOF | kubectl create -f -
apiVersion: formance.com/v1beta1
kind: Gateway
metadata:
  name: gateway0
spec:
  stack: stack0
#  ingress:
#    host: testing.formance.dev
#    scheme: https
EOF
```

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/operator:tag
```

**Run locally without building/pushing image**

```sh
make run
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Settings

See defined settings in [Settings](docs/09-Configuration%20reference/01-Settings.md)

## Contributing

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

