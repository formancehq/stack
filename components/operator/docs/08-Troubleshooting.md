# Troubleshooting

There are two places where you can encounter issues with your deployment.

## Option 1: Operator failures

If the operator itself encounters an error, it is very likely that there is either a bug in the operator or an issue with your Kubernetes. In this case, you should contact the Formance support team.

## Option 2: Stack failures

If the operator is running fine, but the stack is failing, you can troubleshoot the stack by activating the open telemetry publication as described in the [observability section](06-Observability/01-Configure OpenTelemetry.md).
