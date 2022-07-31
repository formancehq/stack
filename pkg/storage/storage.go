package storage

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/numary/auth-membership-gateway/pkg"
	"github.com/numary/auth/pkg"
	"golang.org/x/text/language"
	"gopkg.in/square/go-jose.v2"
	"gorm.io/gorm"

	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
)

var (
	//serviceKey1 is a public key which will be used for the JWT Profile Authorization Grant
	//the corresponding private key is in the service-key1.json (for demonstration purposes)
	serviceKey1 = &rsa.PublicKey{
		N: func() *big.Int {
			n, _ := new(big.Int).SetString("00f6d44fb5f34ac2033a75e73cb65ff24e6181edc58845e75a560ac21378284977bb055b1a75b714874e2a2641806205681c09abec76efd52cf40984edcf4c8ca09717355d11ac338f280d3e4c905b00543bdb8ee5a417496cb50cb0e29afc5a0d0471fd5a2fa625bd5281f61e6b02067d4fe7a5349eeae6d6a4300bcd86eef331", 16)
			return n
		}(),
		E: 65537,
	}
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
	services   map[string]pkg.Service
	signingKey signingKey
	db         *gorm.DB
}

func (s *storage) ClientCredentialsTokenRequest(ctx context.Context, clientID string, scopes []string) (op.TokenRequest, error) {
	client, err := s.getClientByClientID(ctx, clientID)
	if err != nil {
		return nil, oidc.ErrInvalidClient().WithDescription("client not found")
	}
	allowedScopes := auth.Array[string]{}
	for _, scope := range scopes {
		if client.Scopes.Contains(scope) {
			allowedScopes = append(allowedScopes, scope)
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

func New(db *gorm.DB) *storage {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	return &storage{
		services: map[string]pkg.Service{
			"service": {
				Keys: map[string]*rsa.PublicKey{
					"key1": serviceKey1,
				},
			},
		},
		signingKey: signingKey{
			ID:        "id",
			Algorithm: "RS256",
			Key:       key,
		},
		db: db,
	}
}

//CreateAuthRequest implements the op.Storage interface
//it will be called after parsing and validation of the authentication request
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

//AuthRequestByID implements the op.Storage interface
//it will be called after the Login UI redirects back to the OIDC endpoint
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

//CreateAccessAndRefreshTokens implements the op.Storage interface
//it will be called for all requests able to return an access and refresh token (Authorization Code Flow, Refresh Token Request)
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

//TokenRequestByRefreshToken implements the op.Storage interface
//it will be called after parsing and validation of the refresh token request
func (s *storage) TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (op.RefreshTokenRequest, error) {
	refresh := &auth.RefreshToken{}
	return newRefreshTokenRequest(refresh), s.db.
		WithContext(ctx).
		Find(refresh, "id = ?", refreshToken).
		Error
}

//TerminateSession implements the op.Storage interface
//it will be called after the user signed out, therefore the access and refresh token of the user of this client must be removed
func (s *storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	return s.db.
		WithContext(ctx).
		Where("application_id = ? AND subject = ?", clientID, userID).
		Delete(&auth.Token{}).
		Error
}

//RevokeToken implements the op.Storage interface
//it will be called after parsing and validation of the token revocation request
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

//GetSigningKey implements the op.Storage interface
//it will be called when creating the OpenID Provider
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

//GetKeySet implements the op.Storage interface
//it will be called to get the current (public) keys, among others for the keys_endpoint or for validating access_tokens on the userinfo_endpoint, ...
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

//GetClientByClientID implements the op.Storage interface
//it will be called whenever information (type, redirect_uris, ...) about the client behind the client_id is needed
func (s *storage) getClientByClientID(ctx context.Context, clientID string) (*auth.Client, error) {
	client := &auth.Client{}
	return client, s.db.
		WithContext(ctx).
		First(client, "id = ?", clientID).
		Error
}

//GetClientByClientID implements the op.Storage interface
//it will be called whenever information (type, redirect_uris, ...) about the client behind the client_id is needed
func (s *storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	client, err := s.getClientByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return newClientFacade(client), nil
}

//AuthorizeClientIDSecret implements the op.Storage interface
//it will be called for validating the client_id, client_secret on token or introspection requests
func (s *storage) AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error {
	client, err := s.getClientByClientID(ctx, clientID)
	if err != nil {
		return err
	}

	if client.Secret != clientSecret {
		return fmt.Errorf("invalid secret")
	}
	return nil
}

//SetUserinfoFromScopes implements the op.Storage interface
//it will be called for the creation of an id_token, so we'll just pass it to the private function without any further check
func (s *storage) SetUserinfoFromScopes(ctx context.Context, userinfo oidc.UserInfoSetter, userID, clientID string, scopes []string) error {
	return s.setUserinfo(ctx, userinfo, userID, clientID, scopes)
}

//SetUserinfoFromToken implements the op.Storage interface
//it will be called for the userinfo endpoint, so we read the token and pass the information from that to the private function
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

//SetIntrospectionFromToken implements the op.Storage interface
//it will be called for the introspection endpoint, so we read the token and pass the information from that to the private function
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

//GetPrivateClaimsFromScopes implements the op.Storage interface
//it will be called for the creation of a JWT access token to assert claims for custom scopes
func (s *storage) GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (claims map[string]interface{}, err error) {
	return claims, nil
}

//GetKeyByIDAndUserID implements the op.Storage interface
//it will be called to validate the signatures of a JWT (JWT Profile Grant and Authentication)
func (s *storage) GetKeyByIDAndUserID(ctx context.Context, keyID, userID string) (*jose.JSONWebKey, error) {
	service, ok := s.services[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	key, ok := service.Keys[keyID]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	return &jose.JSONWebKey{
		KeyID: keyID,
		Use:   "sig",
		Key:   key,
	}, nil
}

//ValidateJWTProfileScopes implements the op.Storage interface
//it will be called to validate the scopes of a JWT Profile Authorization Grant request
func (s *storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	allowedScopes := make([]string, 0)
	for _, scope := range scopes {
		if scope == oidc.ScopeOpenID {
			allowedScopes = append(allowedScopes, scope)
		}
	}
	return allowedScopes, nil
}

//Health implements the op.Storage interface
func (s *storage) Health(ctx context.Context) error {
	// TODO
	return nil
}

//createRefreshToken will store a refresh_token in-memory based on the provided information
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

//renewRefreshToken checks the provided refresh_token and creates a new one based on the current
func (s *storage) renewRefreshToken(tx *gorm.DB, currentRefreshToken string, newToken string) error {
	return tx.
		Model(&auth.RefreshToken{}).
		Where("id = ?", currentRefreshToken).
		Update("token", newToken).Error
}

//accessToken will store an access_token in-memory based on the provided information
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

func (s *storage) deleteAccessToken(tx *gorm.DB, id string) error {
	return tx.Where("id = ?", id).Delete(&auth.Token{}).Error
}

//setUserinfo sets the info based on the user, scopes and if necessary the clientID
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
			// TODO: Support that
		case oidc.ScopePhone:
			// TODO: Support that ?
		}
	}
	return nil
}

//getInfoFromRequest returns the clientID, authTime and amr depending on the op.TokenRequest type / implementation
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
