[Formance Stack](https://formance.com) • [![Ledger Stars](https://img.shields.io/github/stars/formancehq/ledger?label=Ledger%20stars)](https://github.com/formancehq/ledger/stargazers) [![License MIT](https://img.shields.io/badge/license-mit-purple)](https://github.com/formancehq/ledger/blob/main/LICENSE) [![YCombinator](https://img.shields.io/badge/Backed%20by-Y%20Combinator-%23f26625)](https://www.ycombinator.com/companies/formance-fka-numary) [![slack](https://img.shields.io/badge/slack-formance-brightgreen.svg?logo=slack)](https://bit.ly/formance-slack)

Formance is a highly modular developer platform to build and operate complex money flows of any size and shapes. It comes with several components, that can be used as a whole as the Formance Stack or separately as standalone micro-services and libraries:

- **Formance Ledger** - Programmable double-entry, immutable source of truth to record internal financial transactions and money movements
- **Formance Payments** - Unified API and data layer for payments processing
- **Formance Numscript** - DSL and virtual machine monetary computations and transactions modeling

## ⚡️ Getting started with Formance Cloud Sandbox

### Install Formance CLI

```SHELL
# macOS
brew tap formancehq/tap
brew install fctl
```

###
```SHELL
# login to formance cloud
fctl login

# create a sandbox stack deployment
# please note: sandbox are made available for testing and not made for production usage
# read more in the docs [1]
fctl stack create foobar

# commit your first ledger transaction
fctl ledger send world foo 100 EUR/2 --ledger=demo

# checkout the control dashboard
fctl ui
```

[1] https://docs.formance.com/guides/newSandbox

## 💻 Getting started locally

### Requirements
1. Make sure docker is installed on your machine.
2. Ensure your docker daemon has at least 5GB uf usable RAM available. Otherwise you will run into random crashes.
3. Make sure Docker Compose is installed and available (it should be the case if you have chosen to install Docker via Docker Desktop); and
4. Make sure Git is also installed on your machine.


### Run the app
To start using Formance Stack, run the following commands in a shell:

```
# Get the code
git clone https://github.com/formancehq/stack.git

# Go to the cloned stack directory
cd stack

# Start the stack containers
docker compose up
```

The Stack's API is exposed at http://localhost/api.

You can run :
````
curl http://localhost/api/ledger/_info
````

## ☁️ Cloud Native Deployment

The Formance Stack is distributed as a collection of binaries, with optional packaging as Docker images and configuration support through command line options and environment variables. The recommended, standard way to deploy the collection of services is to a Kubernetes cluster through our Formance official Helm charts, which repository is available at [helm.formance.com](https://helm.formance.com/).

## 📚 Documentation

The full documentation for the formance stack can be found at [docs.formance.com](https://docs.formance.com)

## 💽 Codebase

Formance is transitioning to a unified public facing monorepo (this one) that imports versioned services submodules and provides a common infrastructure layer. As we are finalizing this transition, this monorepo is structured as below

### Technologies

The Formance Stack is built on open-source, battle tested technologies including:

- **PostgreSQL** - Main storage backend
- **Kafka/NATS** - Cross-services async communication
- **Traefik** - Main HTTP gateway
