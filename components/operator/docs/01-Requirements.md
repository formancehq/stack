Before creating our region and deploying the Formance Operator, let's make sure that our cluster is ready to host our operator and have the required following items below in place.

## Ingress controller (optional)

If you want to expose your stacks outside of your Kubernetes cluster, you'll need an Ingress Controller. At Formance, we use Traefik as our Ingress Controller, but you can use another Ingress Controller without any issues.

SSL certificate management can be done either at the level of your LoadBalancer upstream of the Ingress Controller or directly by your Ingress Controller.

:::info
As the Formance Operator will create standard `Ingress` objects to be picked up, alternative ingress controllers can work just as great but might require additional configuration not covered in this setup guide.
:::

## Stateful dependencies

To function properly, the Formance Modules deployed in your Kubernetes cluster will need certain dependencies among the following stateful dependencies.

The list of dependencies changes according to the modules you wish to deploy.
In every module installation guide, you will find a list of required dependencies as well as a guide on how to integrate them with your Formance deployment.
