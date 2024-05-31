It is good practice to keep the operator up to date with the latest [releases](https://github.com/formancehq/stack/releases) to anticipate for any future component upgrades. As there is no immediate impact on the deployed stacks when upgrading the operator, it is advised to simply ensure the operator is running the latest version most of the time.

You can download the latest version of Formance directly from [our github repo](https://github.com/formancehq/stack/releases).
```bash
helm upgrade regions oci://ghcr.io/formancehq/helm/regions \
--version LATEST_RELEASE \
--namespace formance-system
```

Once the upgrade is complete, you can verify the operator is running the latest version with the usual `kubectl` commands:

```bash
kubectl -n formance-system describe deployments operator
```

## Components upgrade
The upgrade process is managed by the operator, who will upgrade the components one by one, as specified in the `versions` CRD. Any migration that needs to be carried out will also be managed by the operator.

When the Operator is upgraded to the latest version, the patch versions will be automatically applied to all the stacks managed by the Operator.

If you wish to change the Minor or Major version, you need to modify your Stack object and specify the new version. The version is represented by the `versionFromFile` object.

```bash
kubectl edit stack stack1
```

:::warning
Do not update the versions CRD unless you are ready to upgrade your deployment.
:::

Once updated, the operator will start the upgrade process. You can follow the process in the `status` section of the `stack` CRD. The status will be updated in real time, and will indicate the health status of the various services, such as whether they are `healthy` or not.

For each component, a [rolling upgrade](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/) will be performed, meaning that the component will be upgraded one by one, and the service will be kept running during the upgrade process.

:::info
Note that the rollout process is managed by Kubernetes, and as such, the operator has no control over it. As a result, the upgrade process may be interrupted by Kubernetes, and a service may be unavailable for a short period of time.
:::



## Updating from Operator v1 to Operator v2

If you wish to update the operator from version 1 to version 2, we have managed the major steps of the migration for you.
Thus, the Configuration CRD from v1 will be migrated to the Settings CRD of v2, and the same will happen for the Versions CRD.

All the Stack CRDs will also be split and migrated to the different Module CRDs of v2.

However, it is important to note that the Stack CRDs from v1 will not be deleted, and the Operator will replicate each change to the v2 objects.
Therefore, after the migration, you will have both v1 and v2 objects in your cluster. We recommend deleting the v1 objects once you are sure that the migration went smoothly.

To migrate the operator from version 1 to version 2, you simply need to update the operator with the following command:

```bash
helm upgrade regions oci://ghcr.io/formancehq/helm/regions \
--version LATEST_RELEASE \
--namespace formance-system
```

As soon as Operator v2 starts, it will begin migrating v1 objects to v2 objects, then create the new associated resources. All steps have been designed to be carried out without any service interruption. However, we recommend performing this update in Staging before doing it in Production.

If you encounter any issues, do not hesitate to contact us. We will be happy to help you with the migration process.