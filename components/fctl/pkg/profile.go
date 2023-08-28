package fctl

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/zitadel/oidc/v2/pkg/client"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
)

type ErrInvalidAuthentication struct {
	err error
}

func (e ErrInvalidAuthentication) Error() string {
	return e.err.Error()
}

func (e ErrInvalidAuthentication) Unwrap() error {
	return e.err
}

func (e ErrInvalidAuthentication) Is(err error) bool {
	_, ok := err.(*ErrInvalidAuthentication)
	return ok
}

func IsInvalidAuthentication(err error) bool {
	return errors.Is(err, &ErrInvalidAuthentication{})
}

func newErrInvalidAuthentication(err error) *ErrInvalidAuthentication {
	return &ErrInvalidAuthentication{
		err: err,
	}
}

const AuthClient = "fctl"

type persistedProfile struct {
	MembershipURI       string                    `json:"membershipURI"`
	Token               *oidc.AccessTokenResponse `json:"token"`
	DefaultOrganization string                    `json:"defaultOrganization"`
}

type Profile struct {
	membershipURI       string
	token               *oidc.AccessTokenResponse
	defaultOrganization string
	config              *Config
}

func (p *Profile) ServicesBaseUrl(stack *membershipclient.Stack) *url.URL {
	baseUrl, err := url.Parse(stack.Uri)
	if err != nil {
		panic(err)
	}
	return baseUrl
}

func (p *Profile) ApiUrl(stack *membershipclient.Stack, service string) *url.URL {
	url := p.ServicesBaseUrl(stack)
	url.Path = "/api/" + service
	return url
}

func (p *Profile) UpdateToken(token *oidc.AccessTokenResponse) {
	p.token = token
}

func (p *Profile) SetMembershipURI(v string) {
	p.membershipURI = v
}

func (p *Profile) MarshalJSON() ([]byte, error) {
	return json.Marshal(persistedProfile{
		MembershipURI:       p.membershipURI,
		Token:               p.token,
		DefaultOrganization: p.defaultOrganization,
	})
}

func (p *Profile) UnmarshalJSON(data []byte) error {
	cfg := &persistedProfile{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}
	*p = Profile{
		membershipURI:       cfg.MembershipURI,
		token:               cfg.Token,
		defaultOrganization: cfg.DefaultOrganization,
	}
	return nil
}

func (p *Profile) GetMembershipURI() string {
	return p.membershipURI
}

func (p *Profile) GetDefaultOrganization() string {
	return p.defaultOrganization
}

func (p *Profile) GetToken(ctx context.Context, httpClient *http.Client) (*oauth2.Token, error) {
	if p.token == nil {
		return nil, errors.New("not authenticated")
	}
	if p.token != nil {
		claims := &oidc.AccessTokenClaims{}
		_, err := oidc.ParseToken(p.token.AccessToken, claims)
		if err != nil {
			return nil, newErrInvalidAuthentication(errors.Wrap(err, "parsing token"))
		}
		if claims.Expiration.AsTime().Before(time.Now()) {
			relyingParty, err := GetAuthRelyingParty(httpClient, p.membershipURI)
			if err != nil {
				return nil, err
			}

			newToken, err := rp.RefreshAccessToken(relyingParty, p.token.RefreshToken, "", "")
			if err != nil {
				return nil, newErrInvalidAuthentication(errors.Wrap(err, "refreshing token"))
			}

			p.UpdateToken(&oidc.AccessTokenResponse{
				AccessToken:  newToken.AccessToken,
				TokenType:    newToken.TokenType,
				RefreshToken: newToken.RefreshToken,
				IDToken:      newToken.Extra("id_token").(string),
			})
			if err := p.config.Persist(); err != nil {
				return nil, err
			}
		}
	}
	claims := &oidc.AccessTokenClaims{}
	_, err := oidc.ParseToken(p.token.AccessToken, claims)
	if err != nil {
		return nil, newErrInvalidAuthentication(err)
	}
	return &oauth2.Token{
		AccessToken:  p.token.AccessToken,
		TokenType:    p.token.TokenType,
		RefreshToken: p.token.RefreshToken,
		Expiry:       claims.Expiration.AsTime(),
	}, nil
}

func (p *Profile) GetClaims() (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	_, _, err := parser.ParseUnverified(p.token.AccessToken, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (p *Profile) GetUserInfo() (*userClaims, error) {
	claims := &userClaims{}
	if p.token != nil && p.token.IDToken != "" {
		_, err := oidc.ParseToken(p.token.IDToken, claims)
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}

func (p *Profile) GetStackToken(ctx context.Context, httpClient *http.Client, stack *membershipclient.Stack) (string, error) {

	form := url.Values{
		"grant_type":         []string{string(oidc.GrantTypeTokenExchange)},
		"audience":           []string{fmt.Sprintf("stack://%s/%s", stack.OrganizationId, stack.Id)},
		"subject_token":      []string{p.token.AccessToken},
		"subject_token_type": []string{"urn:ietf:params:oauth:token-type:access_token"},
	}

	membershipDiscoveryConfiguration, err := client.Discover(p.membershipURI, httpClient)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, membershipDiscoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(AuthClient, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if ret.StatusCode != http.StatusOK {
		data, err := io.ReadAll(ret.Body)
		if err != nil {
			panic(err)
		}
		return "", errors.New(string(data))
	}

	securityToken := oauth2.Token{}
	if err := json.NewDecoder(ret.Body).Decode(&securityToken); err != nil {
		return "", err
	}

	apiUrl := p.ApiUrl(stack, "auth")
	form = url.Values{
		"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"assertion":  []string{securityToken.AccessToken},
		"scope":      []string{"openid email"},
	}

	stackDiscoveryConfiguration, err := client.Discover(apiUrl.String(), httpClient)
	if err != nil {
		return "", err
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, stackDiscoveryConfiguration.TokenEndpoint,
		bytes.NewBufferString(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ret, err = httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if ret.StatusCode != http.StatusOK {
		data, err := io.ReadAll(ret.Body)
		if err != nil {
			panic(err)
		}
		return "", errors.New(string(data))
	}

	stackToken := oauth2.Token{}
	if err := json.NewDecoder(ret.Body).Decode(&stackToken); err != nil {
		return "", err
	}

	return stackToken.AccessToken, nil
}

func (p *Profile) SetDefaultOrganization(o string) {
	p.defaultOrganization = o
}

func (p *Profile) IsConnected() bool {
	return p.token != nil
}

type CurrentProfile Profile

func ListProfiles(flags *flag.FlagSet, toComplete string) ([]string, error) {
	config, err := GetConfig(flags)
	if err != nil {
		return []string{}, nil
	}

	ret := make([]string, 0)
	for p := range config.GetProfiles() {
		if strings.HasPrefix(p, toComplete) {
			ret = append(ret, p)
		}
	}
	sort.Strings(ret)
	return ret, nil
}
