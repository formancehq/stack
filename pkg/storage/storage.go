package storage

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	auth "github.com/formancehq/auth/pkg"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
	"golang.org/x/text/language"
	"gopkg.in/square/go-jose.v2"
	"gorm.io/gorm"
)

type Storage interface {
	op.Storage
	op.ClientCredentialsStorage
	MarkAuthRequestAsDone(ctx context.Context, id, subject string) error
	CreateUser(ctx context.Context, user *auth.User) error
	FindUserByEmail(ctx context.Context, email string) (*auth.User, error)
}

var _ Storage = (*storage)(nil)

type storage struct {
	signingKey    signingKey
	db            *gorm.DB
	relyingParty  rp.RelyingParty
	staticClients []*auth.Client
}

func (s *storage) ClientCredentialsTokenRequest(ctx context.Context, clientID string, scopes []string) (op.TokenRequest, error) {
	client := &auth.Client{}
	err := s.db.
		WithContext(ctx).
		Preload("Scopes").
		First(client, "id = ?", clientID).
		Error
	if err != nil {
		return nil, oidc.ErrInvalidClient().WithDescription("client not found")
	}
	allowedScopes := auth.Array[string]{}
	verifiedScopes := make(map[string]any)
	scopesToCheck := client.Scopes

l:
	for _, scope := range scopes {
		for {
			if len(scopesToCheck) == 0 {
				break l
			}
			clientScope := scopesToCheck[0]
			if len(scopesToCheck) == 1 {
				scopesToCheck = make([]auth.Scope, 0)
			} else {
				scopesToCheck = scopesToCheck[1:]
			}
			if clientScope.Label == scope {
				allowedScopes = append(allowedScopes, scope)
				continue l
			}
			verifiedScopes[clientScope.ID] = struct{}{}

			triggeredScopes := make([]auth.Scope, 0)
			if err := s.db.
				WithContext(ctx).
				Model(clientScope).
				Association("TransientScopes").
				Find(&triggeredScopes); err != nil {
				return nil, err
			}
			scopesToCheck = append(scopesToCheck, triggeredScopes...)
		}
	}
	return &auth.Request{
		ID:            uuid.NewString(),
		CreatedAt:     time.Now(),
		ApplicationID: clientID,
		Scopes:        allowedScopes,
	}, nil
}

func (s *storage) CreateUser(ctx context.Context, user *auth.User) error {
	return s.db.Where(ctx).Create(user).Error
}

func (s *storage) FindUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	user := &auth.User{}
	return user, s.db.
		WithContext(ctx).
		First(user, "email = ?", email).
		Error
}

type signingKey struct {
	ID        string
	Algorithm string
	Key       *rsa.PrivateKey
}

func New(db *gorm.DB, relyingParty rp.RelyingParty, key *rsa.PrivateKey, opts []auth.ClientOptions) *storage {
	var staticClients []*auth.Client
	for _, c := range opts {
		staticClients = append(staticClients, auth.NewClient(c))
	}

	return &storage{
		relyingParty: relyingParty,
		signingKey: signingKey{
			ID:        "id",
			Algorithm: "RS256",
			Key:       key,
		},
		staticClients: staticClients,
		db:            db,
	}
}

// CreateAuthRequest implements the op.Storage interface
// it will be called after parsing and validation of the authentication request
func (s *storage) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	request := &auth.Request{
		CreatedAt:     time.Now(),
		ApplicationID: authReq.ClientID,
		CallbackURI:   authReq.RedirectURI,
		TransferState: authReq.State,
		Prompt:        auth.PromptToInternal(authReq.Prompt),
		UiLocales:     auth.Array[language.Tag](authReq.UILocales),
		LoginHint:     authReq.LoginHint,
		MaxAuthAge:    auth.MaxAgeToInternal(authReq.MaxAge),
		Scopes:        auth.Array[string](authReq.Scopes),
		ResponseType:  authReq.ResponseType,
		Nonce:         authReq.Nonce,
		CodeChallenge: &auth.OIDCCodeChallenge{
			Challenge: authReq.CodeChallenge,
			Method:    string(authReq.CodeChallengeMethod),
		},
		ID: uuid.NewString(),
	}
	return request, s.db.
		WithContext(ctx).
		Create(request).
		Error
}

func (s *storage) MarkAuthRequestAsDone(ctx context.Context, id, subject string) error {
	return s.db.
		Model(&auth.Request{}).
		Where("id = ?", id).
		Update("subject", subject).
		Error
}

// AuthRequestByID implements the op.Storage interface
// it will be called after the Login UI redirects back to the OIDC endpoint
func (s *storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	request := &auth.Request{}
	return request, s.db.
		WithContext(ctx).
		First(request, "id = ?", id).
		Error
}

func (s *storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	request := &auth.Request{}
	return request, s.db.
		WithContext(ctx).
		First(request, "code = ?", code).
		Error
}

func (s *storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	return s.db.
		WithContext(ctx).
		Model(&auth.Request{}).
		Where("id = ?", id).
		Update("code", code).
		Error
}

func (s *storage) DeleteAuthRequest(ctx context.Context, id string) error {
	return s.db.
		WithContext(ctx).
		Delete(&auth.Request{}, "id = ?", id).
		Error
}

func (s *storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (string, time.Time, error) {
	var applicationID string
	authReq, ok := request.(*auth.Request)
	if ok {
		applicationID = authReq.ApplicationID
	}
	token, err := s.saveAccessToken(s.db.WithContext(ctx), applicationID, request.GetSubject(), request.GetAudience(), request.GetScopes())
	if err != nil {
		return "", time.Time{}, err
	}

	return token.ID, token.Expiration, nil
}

// CreateAccessAndRefreshTokens implements the op.Storage interface
// it will be called for all requests able to return an access and refresh token (Authorization Code Flow, Refresh Token Request)
func (s *storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error) {
	//get the information depending on the request type / implementation
	applicationID, authTime, amr := getInfoFromRequest(request)

	var (
		accessToken  = &auth.Token{}
		refreshToken = &auth.RefreshToken{}
	)

	//if currentRefreshToken is empty (Code Flow) we will have to create a new refresh token
	if currentRefreshToken == "" {
		if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if accessToken, err = s.saveAccessToken(tx, applicationID, request.GetSubject(), request.GetAudience(), request.GetScopes()); err != nil {
				return err
			}

			if refreshToken, err = s.createRefreshToken(tx, accessToken, amr, authTime); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return "", "", time.Time{}, err
		}

		return accessToken.ID, refreshToken.ID, accessToken.Expiration, nil
	}

	//if we get here, the currentRefreshToken was not empty, so the call is a refresh token request
	//we therefore will have to check the currentRefreshToken and renew the refresh token

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.First(refreshToken, "id = ?", currentRefreshToken).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", refreshToken.Token).Delete(&auth.Token{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", refreshToken.ID).Delete(&auth.RefreshToken{}).Error; err != nil {
			return err
		}

		if accessToken, err = s.saveAccessToken(tx, applicationID, request.GetSubject(), request.GetAudience(), request.GetScopes()); err != nil {
			return err
		}

		if refreshToken, err = s.createRefreshToken(tx, accessToken, amr, authTime); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken.ID, refreshToken.ID, accessToken.Expiration, nil
}

// TokenRequestByRefreshToken implements the op.Storage interface
// it will be called after parsing and validation of the refresh token request
func (s *storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	refresh := &auth.RefreshToken{}
	return newRefreshTokenRequest(refresh), s.db.
		WithContext(ctx).
		Find(refresh, "id = ?", refreshToken).
		Error
}

// TerminateSession implements the op.Storage interface
// it will be called after the user signed out, therefore the access and refresh token of the user of this client must be removed
func (s *storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	return s.db.
		WithContext(ctx).
		Where("application_id = ? AND subject = ?", clientID, userID).
		Delete(&auth.Token{}).
		Error
}

// RevokeToken implements the op.Storage interface
// it will be called after parsing and validation of the token revocation request
func (s *storage) RevokeToken(ctx context.Context, token string, userID string, clientID string) *oidc.Error {
	err := s.db.
		WithContext(ctx).
		Where("id = ? AND application_id = ? AND subject = ?", token, clientID, userID).
		Delete(&auth.Token{}).
		Error
	if err != nil {
		return oidc.ErrServerError().WithDescription(err.Error())
	}
	return nil
}

// GetSigningKey implements the op.Storage interface
// it will be called when creating the OpenID Provider
func (s *storage) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey) {
	//in this example the signing key is a static rsa.PrivateKey and the algorithm used is RS256
	//you would obviously have a more complex implementation and store / retrieve the key from your database as well
	//
	//the idea of the signing key channel is, that you can (with what ever mechanism) rotate your signing key and
	//switch the key of the signer via this channel
	keyCh <- jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(s.signingKey.Algorithm), //always tell the signer with algorithm to use
		Key: jose.JSONWebKey{
			KeyID: s.signingKey.ID, //always give the key an id so, that it will include it in the token header as `kid` claim
			Key:   s.signingKey.Key,
		},
	}
}

// GetKeySet implements the op.Storage interface
// it will be called to get the current (public) keys, among others for the keys_endpoint or for validating access_tokens on the userinfo_endpoint, ...
func (s *storage) GetKeySet(ctx context.Context) (*jose.JSONWebKeySet, error) {
	//as mentioned above, this example only has a single signing key without key rotation,
	//so it will directly use its public key
	//
	//when using key rotation you typically would store the public keys alongside the private keys in your database
	//and give both of them an expiration date, with the public key having a longer lifetime (e.g. rotate private key every
	return &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{
			KeyID:     s.signingKey.ID,
			Algorithm: s.signingKey.Algorithm,
			Use:       oidc.KeyUseSignature,
			Key:       &s.signingKey.Key.PublicKey,
		}},
	}, nil
}

// GetClientByClientID implements the op.Storage interface
// it will be called whenever information (type, redirect_uris, ...) about the client behind the client_id is needed
func (s *storage) getClientByClientID(ctx context.Context, clientID string) (*auth.Client, error) {
	for _, c := range s.staticClients {
		if c.Id == clientID {
			return c, nil
		}
	}

	client := &auth.Client{}
	return client, s.db.
		WithContext(ctx).
		First(client, "id = ?", clientID).
		Error
}

// GetClientByClientID implements the op.Storage interface
// it will be called whenever information (type, redirect_uris, ...) about the client behind the client_id is needed
func (s *storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	client, err := s.getClientByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return newClientFacade(client, s.relyingParty, s.db), nil
}

// AuthorizeClientIDSecret implements the op.Storage interface
// it will be called for validating the client_id, client_secret on token or introspection requests
func (s *storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	client, err := s.getClientByClientID(ctx, clientID)
	if err != nil {
		return err
	}

	if client.HasSecret(clientSecret) {
		return nil
	}

	return fmt.Errorf("invalid secret")
}

// SetUserinfoFromScopes implements the op.Storage interface
// it will be called for the creation of an id_token, so we'll just pass it to the private function without any further check
func (s *storage) SetUserinfoFromScopes(ctx context.Context, userinfo oidc.UserInfoSetter, userID, clientID string, scopes []string) error {
	return s.setUserinfo(ctx, userinfo, userID, clientID, scopes)
}

// SetUserinfoFromToken implements the op.Storage interface
// it will be called for the userinfo endpoint, so we read the token and pass the information from that to the private function
func (s *storage) SetUserinfoFromToken(ctx context.Context, userinfo oidc.UserInfoSetter, tokenID, subject, origin string) error {
	token := &auth.Token{}
	if err := s.db.Find(token, "id = ?", tokenID).Error; err != nil {
		return err
	}

	//the userinfo endpoint should support CORS. If it's not possible to specify a specific origin in the CORS handler,
	//and you have to specify a wildcard (*) origin, then you could also check here if the origin which called the userinfo endpoint here directly
	//note that the origin can be empty (if called by a web client)
	//
	//if origin != "" {
	//	client, ok := s.clients[token.ApplicationID]
	//	if !ok {
	//		return fmt.Errorf("client not found")
	//	}
	//	if err := checkAllowedOrigins(client.allowedOrigins, origin); err != nil {
	//		return err
	//	}
	//}
	return s.setUserinfo(ctx, userinfo, token.Subject, token.ApplicationID, token.Scopes)
}

// SetIntrospectionFromToken implements the op.Storage interface
// it will be called for the introspection endpoint, so we read the token and pass the information from that to the private function
func (s *storage) SetIntrospectionFromToken(ctx context.Context, introspection oidc.IntrospectionResponse, tokenID, subject, clientID string) error {
	token := &auth.Token{}
	if err := s.db.Find(token, "id = ?", tokenID).Error; err != nil {
		return err
	}

	//check if the client is part of the requested audience
	for _, aud := range token.Audience {
		if aud == clientID {
			//the introspection response only has to return a boolean (active) if the token is active
			//this will automatically be done by the library if you don't return an error
			//you can also return further information about the user / associated token
			//e.g. the userinfo (equivalent to userinfo endpoint)
			err := s.setUserinfo(ctx, introspection, subject, clientID, token.Scopes)
			if err != nil {
				return err
			}
			//...and also the requested scopes...
			introspection.SetScopes(token.Scopes)
			//...and the client the token was issued to
			introspection.SetClientID(token.ApplicationID)
			return nil
		}
	}
	return fmt.Errorf("token is not valid for this client")
}

// GetPrivateClaimsFromScopes implements the op.Storage interface
// it will be called for the creation of a JWT access token to assert claims for custom scopes
func (s *storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	return map[string]any{
		"scp": scopes,
	}, nil
}

// GetKeyByIDAndUserID implements the op.Storage interface
// it will be called to validate the signatures of a JWT (JWT Profile Grant and Authentication)
func (s *storage) GetKeyByIDAndUserID(ctx context.Context, keyID, userID string) (*jose.JSONWebKey, error) {
	return nil, errors.New("not implemented")
}

// ValidateJWTProfileScopes implements the op.Storage interface
// it will be called to validate the scopes of a JWT Profile Authorization Grant request
func (s *storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	allowedScopes := make([]string, 0)
	for _, scope := range scopes {
		if scope == oidc.ScopeOpenID {
			allowedScopes = append(allowedScopes, scope)
		}
	}
	return allowedScopes, nil
}

// Health implements the op.Storage interface
func (s *storage) Health(ctx context.Context) error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

// createRefreshToken will store a refresh_token in-memory based on the provided information
func (s *storage) createRefreshToken(tx *gorm.DB, accessToken *auth.Token, amr []string, authTime time.Time) (*auth.RefreshToken, error) {
	token := &auth.RefreshToken{
		ID:            uuid.NewString(),
		Token:         accessToken.ID,
		AuthTime:      authTime,
		AMR:           amr,
		ApplicationID: accessToken.ApplicationID,
		UserID:        accessToken.Subject,
		Audience:      accessToken.Audience,
		Expiration:    time.Now().Add(5 * time.Hour),
		Scopes:        accessToken.Scopes,
	}
	return token, tx.
		Create(token).
		Error
}

// accessToken will store an access_token in-memory based on the provided information
func (s *storage) saveAccessToken(tx *gorm.DB, applicationID, subject string, audience, scopes []string) (*auth.Token, error) {
	token := &auth.Token{
		ID:            uuid.NewString(),
		ApplicationID: applicationID,
		Subject:       subject,
		Audience:      audience,
		Expiration:    time.Now().Add(5 * time.Minute),
		Scopes:        scopes,
	}
	return token, tx.
		Create(token).
		Error
}

// setUserinfo sets the info based on the user, scopes and if necessary the clientID
func (s *storage) setUserinfo(ctx context.Context, userInfo oidc.UserInfoSetter, userID, clientID string, scopes []string) (err error) {

	user := &auth.User{}
	err = s.db.WithContext(ctx).Find(user, "subject = ?", userID).Error
	if err != nil {
		return err
	}

	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.SetSubject(userID)
		case oidc.ScopeEmail:
			userInfo.SetEmail(user.Email, true)
		case oidc.ScopeProfile:
			// TODO: Support that ?
		case oidc.ScopePhone:
			// TODO: Support that ?
		}
	}
	return nil
}

// getInfoFromRequest returns the clientID, authTime and amr depending on the op.TokenRequest type / implementation
func getInfoFromRequest(req op.TokenRequest) (clientID string, authTime time.Time, amr []string) {
	authReq, ok := req.(*auth.Request) //Code Flow (with scope offline_access)
	if ok {
		return authReq.ApplicationID, authReq.AuthTime, authReq.GetAMR()
	}

	refreshReq, ok := req.(*RefreshTokenRequest) //Refresh Token Request
	if ok {
		return refreshReq.ApplicationID, refreshReq.AuthTime, refreshReq.AMR
	}

	return "", time.Time{}, nil
}
