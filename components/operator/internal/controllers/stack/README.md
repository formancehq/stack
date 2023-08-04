
# Integration Tests

## Init k3d-cluster: Setup the cluster, use other projects

1. `garden create-cluster`
2. `garden deploy`



# Run the operator

> From project root

```bash

go run --disable-webhooks


kubectl apply -f ./deploy/example-v1beta3.yaml

```

## Forward Setup

```sh
# Forward Minio Api port
kubectl port-forward --namespace default svc/minio 9000:9000

# Forward ELK
kubectl port-forward --namespace default svc/elasticsearch-master 9200:9200

# Forward Nats
kubectl port-forward --namespace default svc/nats 4222:4222


# Postgres is forwarded using NodePort config in garden configuration
```


## Run tests

> Now the configuration is deployed as `stacks`

Run test with `go test` in ./


## Uninstall everything

1. Open [Project Operator](github.com/formancehq/stacks/components/operator)
2. Then with `Garden`: `garden delete-cluster`