package wallet

import (
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/go-chi/chi/v5"
)

var balanceNameRegex = regexp.MustCompile("[0-9A-Za-z_-]+")

type CreateBalance struct {
	WalletID  string     `json:"walletID"`
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	Priority  int        `json:"priority"`
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
	Priority  int        `json:"priority"`
}

func (b Balance) LedgerMetadata(walletID string) metadata.Metadata {
	m := metadata.Metadata{
		MetadataKeyWalletID:         walletID,
		MetadataKeyWalletBalance:    TrueValue,
		MetadataKeyBalanceName:      b.Name,
		MetadataKeyBalanceExpiresAt: "",
		MetadataKeyBalancePriority:  fmt.Sprint(b.Priority),
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
	case b[i].Name == "main":
		return false
	case b[j].Name == "main":
		return true
	case b[i].ExpiresAt == nil && b[j].ExpiresAt != nil:
		return false
	case b[i].ExpiresAt != nil && b[j].ExpiresAt == nil:
		return true
	case b[i].ExpiresAt != nil && b[j].ExpiresAt != nil:
		return b[i].ExpiresAt.Before(*b[j].ExpiresAt)
	case b[i].ExpiresAt == nil && b[j].ExpiresAt == nil:
		return b[i].Priority < b[j].Priority
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
	priorityRaw := GetMetadata(account, MetadataKeyBalancePriority)
	var (
		priority int64
		err      error
	)
	if priorityRaw != "" {
		priority, err = strconv.ParseInt(priorityRaw, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	return Balance{
		Name:      GetMetadata(account, MetadataKeyBalanceName),
		ExpiresAt: expiresAt,
		Priority:  int(priority),
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
