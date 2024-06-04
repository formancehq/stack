# Demo deployment

:::warning
This guide assumes that you have already installed the Formance Operator on your cluster. If you haven't done so, please refer to the [installation guide](02-Installation.md) before proceeding.
:::

If you want to have a sample deployment, you can use our Demo Helm chart which will deploy all the necessary resources to install Formance on your cluster.
This will install the following components:
- Gateway
- Ledger
- Payments
- PostgreSQL

```bash
helm upgrade --install demo oci://ghcr.io/formancehq/helm/demo \
--version 2.0.0 \
--namespace formance-system \
--create-namespace
```

You can verify that the installation was successful by using the following command:
```bash
kubectl get pods,svc -n formance-dev
```

### How to access to Demo environement
```BASH
kubectl port-forward -n formance-dev svc/gateway 8080:8080
```
You can now check the different versions of the components installed on your cluster using the following command:
```bash
curl http://localhost:8080/versions
```
