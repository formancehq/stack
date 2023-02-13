# Load testing

## Structure

Tests are located in tests folder.
Libraries are located in libs folder.

A task file help to start tests.

## Run a test

```
task run:ledger
```

This will :
* Start an influxdb container
* Start a grafana container
* Create the influxdb database associated to the test
* Create the datasource on grafana for the influxdb database
* Import the k6 dashboard into grafana using the previously created datasource

Enjoy on http://localhost:3000

Base grafana credentials are admin/admin, please keep these to allow automated script to do their job.
