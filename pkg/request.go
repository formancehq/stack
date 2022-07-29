package auth

import (
	"time"

	"github.com/zitadel/oidc/pkg/oidc"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type Request struct {
	ID            string `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ApplicationID string
	CallbackURI   string
	TransferState string
	Prompt        Array[string]       `gorm:"type:text"`
	UiLocales     Array[language.Tag] `gorm:"type:text"`
	LoginHint     string
	MaxAuthAge    *time.Duration
	Scopes        Array[string] `gorm:"type:text"`
	ResponseType  oidc.ResponseType
	Nonce         string
	CodeChallenge *OIDCCodeChallenge `gorm:"embedded"`
	Subject       string
	AuthTime      time.Time
	Code          string
}

func (a *Request) GetID() string {
	return a.ID
}

func (a *Request) GetACR() string {
	return "" //we won't handle acr in this example
}

func (a *Request) GetAMR() []string {
	return nil
}

func (a *Request) GetAudience() []string {
	return []string{a.ApplicationID}
}

func (a *Request) GetAuthTime() time.Time {
	return a.AuthTime
}

func (a *Request) GetClientID() string {
	return a.ApplicationID
}

func (a *Request) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(a.CodeChallenge)
}

func (a *Request) GetNonce() string {
	return a.Nonce
}

func (a *Request) GetRedirectURI() string {
	return a.CallbackURI
}

func (a *Request) GetResponseType() oidc.ResponseType {
	return a.ResponseType
}

func (a *Request) GetResponseMode() oidc.ResponseMode {
	return "" //we won't handle response mode in this example
}

func (a *Request) GetScopes() []string {
	return a.Scopes
}

func (a *Request) GetState() string {
	return a.TransferState
}

func (a *Request) GetSubject() string {
	return a.Subject
}

func (a *Request) Done() bool {
	return a.Subject != ""
}

func PromptToInternal(oidcPrompt oidc.SpaceDelimitedArray) []string {
	prompts := make([]string, len(oidcPrompt))
	for _, oidcPrompt := range oidcPrompt {
		switch oidcPrompt {
		case oidc.PromptNone,
			oidc.PromptLogin,
			oidc.PromptConsent,
			oidc.PromptSelectAccount:
			prompts = append(prompts, oidcPrompt)
		}
	}
	return prompts
}

func MaxAgeToInternal(maxAge *uint) *time.Duration {
	if maxAge == nil {
		return nil
	}
	dur := time.Duration(*maxAge) * time.Second
	return &dur
}
