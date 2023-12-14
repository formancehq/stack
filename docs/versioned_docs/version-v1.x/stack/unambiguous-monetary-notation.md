---
title: Unambiguous Monetary Notation
---

# Unambiguous Monetary Notation

The Formance Platform uses a unified, safe-by-design way of representing monetary values across all its services and components. We call this representation the _Unambiguous Monetary Notation_, or UMN for short.

While you can use any `[A-Z]{1,16}(\/\d{1,6})` asset in your ledger transactions, it is encouraged to always use UMN, especially if you're dealing with any of the standardized [ISO-4217](https://en.wikipedia.org/wiki/ISO_4217) currencies.

## Specification

An UMN value is represented as:

```text
[ASSET/SCALE AMOUNT]
```


Where:
* `ASSET` is a string of 1 to 16 uppercase letters, representing the currency code of the asset, either standardized or fictional.
* `SCALE` represents the negative power of ten to multiply the amount with to obtain the decimal value in the given asset
* `AMOUNT` is an unsigned integer.

As an example `[USD/2 30]` is equivalent to `USD 30*1E-2`, i.e `USD 0.30`, i.e 30 USD cents.

:::info
For values where the amount already represents the amount of said asset, a scale of zero should not be represented, e.g. `[JPY 100]`.
:::

## Precision

The UMN specification does not enforce a specific precision of the amount, beyond the fact that it must be represented as an unsigned integer. Decisions on the precision of the amount are left to the implementation when implemented by a third party. Internally, Formance Stack components all use arbitrary precision unsigned integers to represent amounts.

## Examples

| UMN | Human Readable | ISO-4217 code |
| --- | --- | --- |
| `[USD/2 30]` | `$0.30` | `USD` |
| `[JPY 100]` | `¥100` | `JPY` |
| `[BTC/8 100000000]` | `1 BTC` | `BTC` |
| `[GBP/2 100]` | `£1.00` | `GBP` |
| `[EUR/2 100]` | `€1.00` | `EUR` |
| `[INR/2 100]` | `₹1.00` | `INR` |
| `[CNY/2 100]` | `¥1.00` | `CNY` |
| `[CAD/2 100]` | `CA$1.00` | `CAD` |

:::info
While `USD/2` is a reasonable notation for most USD-handling use-cases, nothing prevents you from using `USD/4` or `USD/6` if you need to represent smaller amounts and subdivisions of USD in your system. The same applies to other currencies, e.g. `JPY/2` or `JPY/4` for Japanese Yen and while such a coin is not in circulation, it is still a valid notation when these amounts are used in a context where they will end up being floored or ceiled to the nearest whole unit later down the line.
:::

## Rationale

The reason behind this recommendation is that using non explicitly scaled currencies like `USD` is inherently ambiguous, with interpretation of the scale left as an exercise to the reader.

If you receive from a payment processor an API response as follows:

```json
{
  "amount": 100,
  "currency": "USD"
}
```

Without more context, it is unfortunately impossible to tell whether the amount is in cents, or in dollars. While best practices dictate that the amount should be denominated in the smallest unit of the currency, this is not always the case as this interpretation is not standardized across payments services providers.

Some services will inevitably use different formats and encoding rules, resulting in situations where both `100`, and `100.30` are happily parsed, leaving the door open to catastrophic consequences.

As you start to scale your business and deal with multiple and specialized payment services providers, the risk of different formats making their way to your internal representation increases along with the risk of misinterpreting the amount.

As Formance components are designed to be used in a variety of contexts and find themselves dealing with a variety formats from different providers, we decided to explicitly specify the scale of the amount in the notation, making UMN really hard to misinterpret by design.
