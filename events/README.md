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

## Gateway AuditLog

`gateway/v2.0.0/AUDIT` documents the Gateway's `app: gateway`, `version: v2`,
`type: AUDIT` envelope. The catalogue key groups the v2 contract; it is not a
claim that every Gateway v2 release emits every optional field. Native gRPC
capture is introduced by [gateway#157](https://github.com/formancehq/gateway/pull/157).
Envelope `version: v2` first appeared in Gateway v2.2.0; Gateway v2.0 and v2.1
published a different `version: v1` envelope that this schema does not describe.
This schema does not change the Gateway version pinned by the Stack build.

Both existing HTTP audit events and native gRPC events retain `payload.id`,
`trace_id`, `actor` and `http`. HTTP bodies remain text, with optional truncation
flags. Header maps contain arrays of strings; absent maps or values may be null.
Claims are an optional extensible object. The schema remains open to additional
fields for additive evolution.

Native gRPC adds the optional `payload.grpc` object:

| Field | Meaning |
| --- | --- |
| `service`, `method` | Inbound RPC service and method. |
| `status_code` | Observed gRPC status from 0 through 16; omitted when unknown, never inferred from HTTP 200. |
| `status_message` | Optional percent-decoded gRPC message. |
| `termination` | `completed`, `incomplete`, `gateway_error`, `canceled`, `deadline_exceeded` or `aborted`. A completed RPC can carry a server error status. |
| `gateway_status_code` | Optional Caddy handler error status, separate from the observed HTTP status. |
| `request`, `response` | Bounded binary capture with `encoding: base64`, optional `data`, `bytes_total`, `truncated` and `complete`. |

Each binary direction stores at most 65536 decoded bytes, including gRPC message
envelopes and compression. `bytes_total` counts observed bytes before finalization;
`truncated` means the stored prefix is shorter than the observed data. Request
`complete` means EOF was observed; response `complete` means normal handler
completion with a valid gRPC status and no local interruption. Neither means the
client consumed the bytes. Empty data is omitted. These binary prefixes are not
decoded protobuf values and are not stored in the HTTP text body fields.

The Gateway emits one audit event at RPC termination. The envelope date is the
publication time. If Caddy generates its error response outside the audit
middleware, `http.response.status_code` is `0` and `grpc.gateway_status_code`
records the proposed error status. Consumers must read the gRPC status and
termination to classify RPC outcomes. Adding the schema does not update
HTTP-status-based downstream billing or reporting views.

Validators must support JSON Schema draft-07, including `allOf` and `if`/`then`,
to enforce the binary capture limit.

`testdata/gateway/` contains synthetic events captured from the real Gateway
producer at commit `aadadb2ce8bc705cb1eb9111b93038918c9263cd`, with Caddy, grpc-go
and NATS: HTTP, unary success, unary NotFound and an unavailable upstream.
Tests validate these fixtures against the embedded schema, standalone generated
schema and both generated catalogues, and reject malformed gRPC captures.

To regenerate and validate the event catalogue without rebuilding OpenAPI:

```sh
just generate-events
cd events
go test -race ./...
```

## Payments Versions

We decided to go with stack releases starting at v2.0.x.

Before that, the last payments version was v0.9.7.

This is why we do not have a v1.0.0 directory with events, since payments
was never released with a v1.x version.
