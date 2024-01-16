---
title: Upgrade existing deployment
---

# Upgrade an existing deployment

## Supported upgrade paths

| Target version | Minimum current version | Zero-downtime supported | Instructions                                                                                                           |
|----------------|-------------------------|-------------------------|------------------------------------------------------------------------------------------------------------------------|
| 1.9.x          | 1.4.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)    •    [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime) |
| 1.8.x          | 1.4.x                   | No                      | [Single-node](#upgrade-a-single-node-deployment)    •    [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime) |
| 1.8.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)   •   [Multi-nodes](#upgrade-a-multi-nodes-deployment)                |
| 1.7.x          | 1.4.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)    •    [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime) |
| 1.7.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)   •   [Multi-nodes](#upgrade-a-multi-nodes-deployment)                |
| 1.6.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)   •   [Multi-nodes](#upgrade-a-multi-nodes-deployment)                |
| 1.6.x          | 1.4.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)    •    [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime) |
| 1.5.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment)                  |
| 1.5.x          | 1.4.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)   •   [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime)  |
| 1.4.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment) • [Multi-nodes](#upgrade-a-multi-nodes-deployment)                    |
| 1.4.x          | 1.4.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime)    |
| 1.3.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment)                  |
| 1.3.x          | 1.3.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime)    |
| 1.2.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment)                  |
| 1.2.x          | 1.2.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime)    |
| 1.1.x          | 1.x                     | No                      | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment)                  |
| 1.1.x          | 1.1.x                   | Yes                     | [Single-node](#upgrade-a-single-node-deployment)  •  [Multi-nodes](#upgrade-a-multi-nodes-deployment-zero-downtime)    |

## Prerequisites

### Compatibility version
Verify that you meet the upgrade path requirements above for your desired version and review the relevant [release notes](https://github.com/formancehq/ledger/releases) for any preparation instructions.

## Upgrade a multi-nodes deployment

### Preparation

* Prior to upgrading, confirm that all your `ledger` instances were cleanly shut down
* Create a backup of your database (e.g. using pg_dump if you're using Postgres as your storage backend, or by copying the SQLite files otherwise)
* Choose a node where you'll be able to execute numary cli commands

### Upgrade instructions

1. Download the desired target version on the node used for upgrade
2. Using the new binary, run:
```
numary storage scan
numary storage list
numary storage upgrade {LEDGER_NAME}
```
3. Restart all your nodes with the upgraded ledger binary

## Upgrade a multi-nodes deployment (zero-downtime)

### Upgrade instructions

1. Download the desired target version
2. Deploy new nodes with the new upgraded version
3. Decomission old nodes with the previous version

## Upgrade a single-node deployment

### Preparation

* Prior to upgrading, confirm that your `ledger` instance was cleanly shut down
* Create a backup of your database (e.g. using pg_dump if you're using Postgres as your storage backend, or by copying the SQLite files otherwise)

### Upgrade instructions

1. Download the desired version of Formance Ledger
2. Using the new binary, run:
```
numary storage scan
numary storage list
numary storage upgrade {LEDGER_NAME}
```
3. Start your `ledger` instance with new version of ledger binary
