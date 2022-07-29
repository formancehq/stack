package auth

import (
	"time"
)

type RefreshToken struct {
	ID            string `gorm:"primarykey"`
	Token         string
	AuthTime      time.Time
	AMR           Array[string] `gorm:"type:text"`
	Audience      Array[string] `gorm:"type:text"`
	UserID        string
	ApplicationID string
	Expiration    time.Time
	Scopes        Array[string] `gorm:"type:text"`
}
