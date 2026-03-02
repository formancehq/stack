# Hanzo Treasury

Programmable financial infrastructure for the Hanzo platform. Double-entry ledger, reconciliation, wallets, payment connectors, and workflow orchestration.

## Architecture

```
hanzo/commerce    Storefront, catalog, orders
       |
hanzo/payments    Payment routing (50+ processors)
       |
hanzo/treasury    Ledger, reconciliation, wallets, flows   <-- you are here
       |
lux/treasury      On-chain treasury, MPC/KMS wallets
```

## Components

| Component | Description |
|-----------|-------------|
| **Ledger** | Programmable double-entry, immutable source of truth for all financial transactions |
| **Payments** | Unified API and data layer for payment processing across providers |
| **Numscript** | DSL for modeling complex monetary computations and transaction flows |
| **Wallets** | Multi-currency virtual wallets with hold/release mechanics |
| **Reconciliation** | Auto-match ledger transactions against payment provider data |
| **Flows** | Workflow orchestration for payment and treasury operations |
| **Webhooks** | Event delivery service for financial state changes |
| **Auth** | Authentication and authorization for treasury APIs |

## Quick Start

```bash
# Install CLI
brew tap hanzoai/tap
brew install hanzo-treasury

# Start local development stack
docker compose up -d

# Create a ledger transaction
curl -X POST http://localhost:3068/v2/ledger/demo/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "postings": [{
      "source": "world",
      "destination": "users:001",
      "amount": 10000,
      "asset": "USD/2"
    }]
  }'
```

## Numscript Example

```numscript
// Transfer with multi-party fee split
send [USD/2 10000] (
  source = @users:001
  destination = {
    90% to @merchants:042
    10% to {
      50% to @platform:fees
      50% to @platform:reserve
    }
  }
)
```

## API

- **Ledger**: `POST /v2/ledger/{name}/transactions` — record transactions
- **Accounts**: `GET /v2/ledger/{name}/accounts` — query balances
- **Payments**: `POST /v2/payments` — initiate payments
- **Wallets**: `POST /v2/wallets` — create/manage wallets
- **Reconciliation**: `POST /v2/reconciliation/policies` — define matching rules

## Integration with Hanzo Stack

Treasury connects to:
- **hanzo/payments** for payment routing
- **hanzo/commerce** for order settlement
- **lux/treasury** for on-chain operations and MPC wallet management
- **hanzo/kms** for secrets and key management

## License

MIT — see [LICENSE](LICENSE)
