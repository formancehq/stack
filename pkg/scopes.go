package auth

import (
	"github.com/google/uuid"
)

type Scope struct {
	ID              string  `json:"id" gorm:"primarykey"`
	Label           string  `json:"label" gorm:"unique"`
	TransientScopes []Scope `json:"transient" gorm:"many2many:transient_scopes;"`
}

func (s *Scope) AddTransientScope(scope *Scope) *Scope {
	s.TransientScopes = append(s.TransientScopes, *scope)
	return s
}

func NewScope(value string) *Scope {
	return &Scope{
		ID:    uuid.NewString(),
		Label: value,
	}
}
