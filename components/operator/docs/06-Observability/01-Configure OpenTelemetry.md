# Configure OpenTelemetry

## Overview

Formance uses OpenTelemetry to collect and send telemetry data to your monitoring system. [OpenTelemetry](https://opentelemetry.io/) is an open-source observability framework for cloud-native software. It provides a single set of APIs, libraries, agents, and instrumentation to capture distributed traces and metrics from your applications.

## Configuration

### Prerequisites

This guide assumes that you already have an OpenTelemetry server running and that you have the necessary credentials to connect to it.

### Configuration

In this example, you'll set up an OpenTelemetry configuration for the `formance-dev` stack. This configuration will apply to all the modules of this stack.

```yaml
apiVersion: formance.com/v1beta1
kind: Settings
metadata:
  name: stacks-otel-collector
spec:
  key: opentelemetry.*.dsn
  stacks:
    - "formance-dev"
  value: grpc://opentelemetry-collector.formance-system.svc:4317?insecure=true
```
