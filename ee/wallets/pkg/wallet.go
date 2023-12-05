package wallet

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/google/uuid"
)

type ListWallets struct {
	Metadata metadata.Metadata
	Name     string
}

type PatchRequest struct {
	Metadata metadata.Metadata `json:"metadata"`
}

func (c *PatchRequest) Bind(r *http.Request) error {
	return nil
}

type CreateRequest struct {
	PatchRequest
	Name string `json:"name"`
}

func (c *CreateRequest) Bind(r *http.Request) error {
	return nil
}

type Wallet struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Metadata  metadata.Metadata `json:"metadata"`
	CreatedAt time.Time         `json:"createdAt"`
	Ledger    string            `json:"ledger"`
}

type WithBalances struct {
	Wallet
	Balances map[string]*big.Int `json:"balances"`
}

func (w *WithBalances) UnmarshalJSON(data []byte) error {
	type view struct {
		Wallet
		Balances struct {
			Main ExpandedBalance `json:"main"`
		} `json:"balances"`
	}
	v := view{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*w = WithBalances{
		Wallet:   v.Wallet,
		Balances: v.Balances.Main.Assets,
	}
	return nil
}

func (w WithBalances) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Wallet
		Balances struct {
			Main ExpandedBalance `json:"main"`
		} `json:"balances"`
	}{
		Wallet: w.Wallet,
		Balances: struct {
			Main ExpandedBalance `json:"main"`
		}{
			Main: ExpandedBalance{
				Assets: w.Balances,
			},
		},
	})
}

func (w Wallet) LedgerMetadata() metadata.Metadata {
	return metadata.Metadata{
		MetadataKeyWalletSpecType: PrimaryWallet,
		MetadataKeyWalletName:     w.Name,
		MetadataKeyWalletID:       w.ID,
		MetadataKeyWalletBalance:  TrueValue,
		MetadataKeyBalanceName:    MainBalance,
		MetadataKeyCreatedAt:      w.CreatedAt.UTC().Format(time.RFC3339Nano),
	}.Merge(EncodeCustomMetadata(w.Metadata))
}

func NewWallet(name, ledger string, m metadata.Metadata) Wallet {
	if m == nil {
		m = metadata.Metadata{}
	}
	return Wallet{
		ID:        uuid.NewString(),
		Metadata:  m,
		Name:      name,
		CreatedAt: time.Now().UTC().Round(time.Nanosecond),
		Ledger:    ledger,
	}
}

func FromAccount(ledger string, account Account) Wallet {
	createdAt, err := time.Parse(time.RFC3339Nano, GetMetadata(account, MetadataKeyCreatedAt))
	if err != nil {
		panic(err)
	}

	return Wallet{
		ID:        GetMetadata(account, MetadataKeyWalletID),
		Name:      GetMetadata(account, MetadataKeyWalletName),
		Metadata:  ExtractCustomMetadata(account),
		CreatedAt: createdAt,
		Ledger:    ledger,
	}
}

func WithBalancesFromAccount(ledger string, account AccountWithVolumesAndBalances) WithBalances {
	return WithBalances{
		Wallet:   FromAccount(ledger, account.Account),
		Balances: account.GetBalances(),
	}
}
