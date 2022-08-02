package auth

import (
	"github.com/google/uuid"
)

type Scope struct {
	ID       string  `gorm:"primarykey" json:"id"`
	Label    string  `json:"label" gorm:"unique"`
	Triggers []Scope `json:"triggers" gorm:"many2many:scopes_triggers;"`
}

func (s *Scope) AddTrigger(scope *Scope) *Scope {
	s.Triggers = append(s.Triggers, *scope)
	return s
}

func NewScope(value string) *Scope {
	return &Scope{
		ID:    uuid.NewString(),
		Label: value,
	}
}
