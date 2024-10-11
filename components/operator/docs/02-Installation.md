The deployment of resources is orchestrated by the Formance Kubernetes Operator, which is driven by the CRDs present in the Kubernetes Cluster and reconciles them against the current state of the cluster.

This essentially means that the operator will be creating and maintaining pods, services, and other resources on your cluster on your behalf. It is also the last mile of this setup guide before you get to actually use your Formance Stack — let's get started!


## Setup


With the emergence of multiple tools and services that use helm templatization. We have choosen to separate `Custom Resource Definitions`(CRDs) separatelty from operator deployment. As noticed within [Helm CRDs best practicies](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/),
this allow us to have more control over the CRDs lifecycle and also avoid deletion of the CRDs when the operator is uninstalled.

### Helm installation

If you don't have Helm installed, you can follow the [official Helm installation guide](https://helm.sh/docs/intro/install/).
Dependending of your needs, you can install the CRDs and the operator separatly or together.

### Operator CRDs
To install the Formance Operator CRDs, you can use the following command:

```bash
helm upgrade --install operator-crds oci://ghcr.io/formancehq/helm/operator-crds \
--version v2.0.19 \
--namespace formance-system \
--create-namespace
```

As noticed above, the version will always be the same as the operator version. CRDs **must always** be upgraded before the operator.

### Operator Deployment

From version v2.0.19, CRDs are now packaged with `helm.sh/resource-policy: keep` to avoid deletion of the CRDs when the operator is uninstalled.

You can deploy Formance Operator using Helm:

```bash
helm upgrade --install regions oci://ghcr.io/formancehq/helm/regions \
--version v2.0.19 \
--namespace formance-system \
--create-namespace \
--set operator.operator-crds.create=false
```

### Migrating from Operator chart with CRDs to dedicated CRDs chart

First make **sure** you've already upgraded the operator chart with crds `operator-crds.create=true` to make sure the `helm.sh/resource-policy: keep` is present.
You can verify this by running:

```bash
kubectl get crds authclients.formance.com -o json | jq .metadata.annotations                                                                                                         
{
  "controller-gen.kubebuilder.io/version": "v0.14.0",
  "helm.sh/resource-policy": "keep", <---- This should be present
  "meta.helm.sh/release-name": "formance-operator",
  "meta.helm.sh/release-namespace": "formance-system"
}
```

If using Helm Release, you might want to set appropriate helm release ownership annotations to the new chart using kubectl:

#### Migration bash script

```bash
#!/bin/bash

# Set the required namespace and release name
NAMESPACE=formance-system
RELEASE_NAME=operator-crds

for GROUP in stack.formance.com formance.com
do
    NAMES=$(kubectl api-resources --api-group=$GROUP --verbs=list -o name)
    for CRD in $NAMES
    do
        echo "----- Annotating $CRD -----"
        kubectl annotate crd $CRD meta.helm.sh/release-name=$RELEASE_NAME --overwrite
        kubectl annotate crd $CRD meta.helm.sh/release-namespace=$NAMESPACE --overwrite
        echo "--------- Done ----------"
    done
done

```

Then you will be able to disable `operator-crds.create: false` and install the operator-crds chart separately.

#### Disable CRDs creation

```bash
helm upgrade --install regions oci://ghcr.io/formancehq/helm/regions \
--version v2.0.19 \
--namespace $NAMESPACE \
--create-namespace \
--set operator.operator-crds.create=false

```

#### Install the new charts

```bash
helm upgrade --install $RELEASE_NAME oci://ghcr.io/formancehq/helm/operator-crds \
--version v2.0.19 \
--namespace $NAMESPACE \
--create-namespace
```