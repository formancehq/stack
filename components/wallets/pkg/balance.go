package wallet

import (
	"math/big"
	"net/http"
	"regexp"
	"time"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/go-chi/chi/v5"
)

var balanceNameRegex = regexp.MustCompile("[0-9A-Za-z_-]+")

type CreateBalance struct {
	WalletID  string     `json:"walletID"`
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

func (c *CreateBalance) Validate() error {
	if !balanceNameRegex.MatchString(c.Name) {
		return ErrInvalidBalanceName
	}
	if c.Name == MainBalance {
		return ErrReservedBalanceName
	}
	return nil
}

func (c *CreateBalance) Bind(r *http.Request) error {
	c.WalletID = chi.URLParam(r, "walletID")
	return nil
}

type Balance struct {
	Name      string     `json:"name,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

func (b Balance) LedgerMetadata(walletID string) metadata.Metadata {
	m := metadata.Metadata{
		MetadataKeyWalletID:         walletID,
		MetadataKeyWalletBalance:    TrueValue,
		MetadataKeyBalanceName:      b.Name,
		MetadataKeyBalanceExpiresAt: "",
	}
	if b.ExpiresAt != nil {
		m[MetadataKeyBalanceExpiresAt] = b.ExpiresAt.Format(time.RFC3339Nano)
	}
	return m
}

func NewBalance(name string, expiresAt *time.Time) Balance {
	return Balance{
		Name:      name,
		ExpiresAt: expiresAt,
	}
}

type Balances []Balance

func (b Balances) Len() int {
	return len(b)
}

func (b Balances) Less(i, j int) bool {
	switch {
	case b[i].ExpiresAt == nil && b[j].ExpiresAt != nil:
		return false
	case b[i].ExpiresAt != nil && b[j].ExpiresAt == nil:
		return true
	case b[i].ExpiresAt != nil && b[j].ExpiresAt != nil:
		return b[i].ExpiresAt.Before(*b[j].ExpiresAt)
	case b[i].ExpiresAt == nil && b[j].ExpiresAt == nil:
		return b[i].Name < b[j].Name
	}
	panic("Should not happen")
}

func (b Balances) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func BalanceFromAccount(account Account) Balance {
	expiresAtRaw := GetMetadata(account, MetadataKeyBalanceExpiresAt)
	var expiresAt *time.Time
	if expiresAtRaw != "" {
		parsedExpiresAt, err := time.Parse(time.RFC3339Nano, expiresAtRaw)
		if err != nil {
			panic(err)
		}
		expiresAt = &parsedExpiresAt
	}
	return Balance{
		Name:      GetMetadata(account, MetadataKeyBalanceName),
		ExpiresAt: expiresAt,
	}
}

type ExpandedBalance struct {
	Balance
	Assets map[string]*big.Int `json:"assets"`
}

func ExpandedBalanceFromAccount(account AccountWithVolumesAndBalances) ExpandedBalance {
	return ExpandedBalance{
		Balance: BalanceFromAccount(account.Account),
		Assets:  account.GetBalances(),
	}
}
