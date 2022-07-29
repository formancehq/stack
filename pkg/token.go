package auth

import (
	"time"
)

type Token struct {
	ID            string `gorm:"primarykey"`
	ApplicationID string
	Subject       string
	Audience      Array[string] `gorm:"type:text"`
	Expiration    time.Time
	Scopes        Array[string] `gorm:"type:text"`
}
