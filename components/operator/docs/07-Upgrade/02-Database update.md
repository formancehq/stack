# Database update

This page describes how the operator updates the database schema when deploying a new version of a module.

## Database updates management

The database migrations are handled automatically by the Formance Kubernetes Operator. This applies for all the services that are deployed by the operator.

:::info
The services come with the guarantee that the database schema will be updated in a backward-compatible way. This means that the database schema will be updated in a way that the old version of the service will still be able to work with the new schema.
:::

## Database update process

The process starts when the operator detects that a new version of a module is set to be deployed. 

First, the operator starts a Kubernetes Job that runs the database migration. The job is responsible for updating the database schema to the new version. Most of the time, the database migration is done using the `migrate` command embeded in the service.

If the database migration fails, the operator will stop the deployment of the new version of the module and will keep the old version running. The database migration is applied as a transaction, so if it fails, the database is not updated and the old version of the module is still able to work with the old schema.

When the database migration is complete, the operator will start the deployment of the new version of the module. It ensures that the new version is running properly before sending traffic to it.

Once the new version is running, the operator will eventually stop the old version of the module and clean up the resources that are not needed anymore.
