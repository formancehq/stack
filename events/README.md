# Events

This repository centralizes event schemas across the Formance Stack. For each stack version, a repository named "vX" contains all related events.

Each "vX" folder contains a "base" folder that contains the base event format, which is common to all services. This base event includes a "type" property and an "app" property, which denote the format of the "payload" property.

For example, an event with "type" == "SAVED_PAYMENT" and "app" == "payments" must have a payload matching schema in the file "payments/SAVED_PAYMENT.yaml".

## Event Envelopes

By default, event schemas are payload schemas composed with `base.yaml`, which
describes the historical Stack envelope (`app`, `version`, `type`, `date`, and
`payload`).

A service version can define a complete, version-specific envelope in
`services/<service>/<version>/base.yaml`. The generator composes each sibling
event schema with that base through JSON Schema `allOf`; `base.yaml` is not
published as an event type.

Ledger v3 uses this mechanism because its NATS and Kafka sinks publish the same
new envelope directly: `type`, `ledger`, `date`, `logSequence`, and the complete
`log`. It does not contain the historical `app`, `version`, or `payload`
properties.

## Payments Versions

We decided to go with stack releases starting at v2.0.x.

Before that, the last payments version was v0.9.7.

This is why we do not have a v1.0.0 directory with events, since payments
was never released with a v1.x version.
