// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type Hold struct {
	Description string   `json:"description"`
	Destination *Subject `json:"destination,omitempty"`
	// The unique ID of the hold.
	ID string `json:"id"`
	// Metadata associated with the hold.
	Metadata map[string]string `json:"metadata"`
	// The ID of the wallet the hold is associated with.
	WalletID string `json:"walletID"`
}
