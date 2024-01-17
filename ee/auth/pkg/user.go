package auth

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID      string `json:"id" bun:",pk"`
	Subject string `json:"subject" bun:",unique"`
	Email   string `json:"email"`
}
