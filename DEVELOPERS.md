# Docs for Developers

The entire monorepo is managed by [Earthly](https://earthly.dev).
And you need [direnv](https://direnv.net/).

## How to start developing
To use Earthly and dev environments in Kubernetes, you need the following environment variables (in .direnv/.env):
Environment variables are split into several pieces:
- Project config: License of the various tools used in the repository
    - SPEAKEASY_API_KEY: License for [Speakeasy](https://www.speakeasyapi.dev/)
    - GITHUB_TOKEN: Github token for actions
    - GORELEASER_KEY: Github token for releases
- vCluster config:
    - FORMANCE_DEV_KUBE_TOKEN: Kubernetes token for dev cluster access
    - FORMANCE_DEV_KUBE_API_SERVER_ADDRESS: dev cluster address
    - FORMANCE_DEV_TLD: dev domain to access the stack
    - FORMANCE_DEV_GHCR_REGISTRY: Registry for images
    - FORMANCE_DEV_USER: User for images
- Admin config: only for our DNS configuration with Route53
    - FORMANCE_AWS_KEY_ID: AWS Key ID
    - FORMANCE_AWS_SECRET_KEY: AWS Secret Key

Once you have configured all the following elements, you can run the following commands to deploy the various components:
```bash
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component traefik
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component cert-manager
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component postgres
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component elasticsearch
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component otel-collector
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component nats
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component jaeger
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component goproxy
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component ghcr-registry
earthly -P --push github.com/formancehq/earthly/bootstrap/+deploy-base-component --component buildkitd
```
You can now start deploying all components.
```bash
earthly -P --push +deploy-all
```

## How to deploy a new version of components
There are several ways to deploy a component:
- Deploy all components
- Deploy a specific component
### Deploy all components
```bash
earthly -P --push +deploy-all
```

### Deploy a specific component
```bash
earthly -P --push +deploy --component <component_name>
```