# Events

This repository centralize events schemas across Formance Stack.
For each stack version, a repository "vX" contains all related events.

Each "vX" folder contains a "base" folder containing the base event format, common to all services.
This base event contains a "type" property, and an "app" property, denoting the format of the "payload" property.

For example, an event with "type" == "SAVED_PAYMENT" and app == "payments", must have a payload matching schema in file payments/SAVED_PAYMENT.yaml.
