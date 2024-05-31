_Version 14 or higher is required._

The recommended way to spin up a PostgreSQL deployment is to use your cloud provider's managed services. It is recommended to start with a small instance and scale up as needed, starting with for example a `db.m4.large` sizing on AWS.


## Create the database settings

### Option 1: Use the same server for all modules

In this example, you'll set up a configuration for a PostgreSQL cluster that will be used only by the `formance-dev` stack but will apply to all the modules of this stack.
Thus, the different modules of the Stack will use this PostgreSQL cluster while being isolated in their own database.

:::info
This database is created following the format: `{stackName}-{module}`
:::

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-postgres-uri
spec:
  key: postgres.*.uri
  stacks:
    - 'formance-dev'
  value: postgresql://formance:formance@postgresql.formance-system.svc:5432?disableSSLMode=true
```

### Option 2: Use different servers for each module

In this example you'll set up a configuration dedicated Postgres cluster for the `ledger` and `payments` modules of the `formance-dev` stack.

```yaml
---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-ledger-postgres-uri
spec:
  key: postgres.ledger.uri
  stacks:
    - 'formance-dev'
  value: postgresql://formance:formance@postgresql-ledger.formance-system.svc:5432?disableSSLMode=true
---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-payments-postgres-uri
spec:
  key: postgres.payments.uri
  stacks:
    - 'formance-dev'
  value: postgresql://formance:formance@postgresql-payments.formance-system.svc:5432?disableSSLMode=true
```

### Option 3: Use PostgreSQL on AWS RDS with an IAM role

In this example, you'll use an AWS IAM role to connect to the database. The `formance-dev` stack will use this configuration.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aws-rds-access-role
  namespace: formance-system
  labels:
    formance.com/stack: any
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::AWS_ACCOUNT_ID:role/AWS_ROLE_NAME
---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-aws-service-account
spec:
  stacks:
    - formance-dev
  key: aws.service-account
  value: aws-rds-access-role
---
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: formance-dev-postgres-uri
spec:
  key: postgres.*.uri
  stacks:
    - 'formance-dev'
  value: postgresql://formance@postgresql.formance-system.svc:5432
 ```
