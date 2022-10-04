package gateway_plugin_auth

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Config struct {
	Issuer           string
	SigningMethodRSA string
	RefreshTimeError time.Duration
	RefreshTime      time.Duration
}

func CreateConfig() *Config {
	return &Config{}
}

type Plugin struct {
	next             http.Handler
	issuer           string
	signingMethodRSA *jwt.SigningMethodRSA
	jwksURI          string
	publicKey        *rsa.PublicKey
	refreshTimeError time.Duration
	refreshTime      time.Duration
}

const (
	refreshTimeErrorDefault = 10 * time.Second
	refreshTimeDefault      = 15 * time.Minute
)

var (
	signingMethodDefault = jwt.SigningMethodRS256
)

func New(ctx context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	fmt.Printf("New with config: %+v\n", config)
	var err error
	p := &Plugin{
		next:             next,
		issuer:           config.Issuer,
		refreshTimeError: config.RefreshTimeError,
		refreshTime:      config.RefreshTime,
	}

	switch config.SigningMethodRSA {
	case "":
		p.signingMethodRSA = signingMethodDefault
	case jwt.SigningMethodRS256.Alg():
		p.signingMethodRSA = jwt.SigningMethodRS256
	case jwt.SigningMethodRS384.Alg():
		p.signingMethodRSA = jwt.SigningMethodRS384
	case jwt.SigningMethodRS512.Alg():
		p.signingMethodRSA = jwt.SigningMethodRS512
	default:
		err := fmt.Errorf("ERROR: unsupported config signing method: %s", config.SigningMethodRSA)
		fmt.Println(err)
		return p, err
	}

	if p.refreshTimeError == 0 {
		p.refreshTimeError = refreshTimeErrorDefault
	}
	if p.refreshTime == 0 {
		p.refreshTime = refreshTimeDefault
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("New: context is done: Plugin: %+v\n", p)
			return p, err
		default:
			if err := p.fetchPublicKeys(ctx); err != nil {
				fmt.Printf("ERROR: Plugin.fetchPublicKeys (first): %s\n", err)
				time.Sleep(p.refreshTimeError)
			} else {
				go p.backgroundRefresh(ctx)
				fmt.Printf("New: Plugin: %+v\n", p)
				return p, nil
			}
		}
	}
}

func (p *Plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenString, err := p.extractToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	parser := new(jwt.Parser)
	token, parts, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if token.Method.Alg() != p.signingMethodRSA.Alg() {
		http.Error(w, fmt.Errorf("unsupported token signing method: %s", token.Method.Alg()).Error(), http.StatusUnauthorized)
		return
	}

	vErr := &jwt.ValidationError{}
	token.Signature = parts[2]
	if err = verifyRSA(strings.Join(parts[0:2], "."), token.Signature, p.publicKey, p.signingMethodRSA); err != nil {
		vErr.Inner = err
		vErr.Errors |= jwt.ValidationErrorSignatureInvalid
		http.Error(w, fmt.Errorf("unverified JWT token: %w", vErr).Error(), http.StatusUnauthorized)
		return
	}

	p.next.ServeHTTP(w, r)
}

func verifyRSA(signingString, signature string, rsaKey *rsa.PublicKey, m *jwt.SigningMethodRSA) error {
	var err error

	var sig []byte
	if sig, err = base64.RawURLEncoding.DecodeString(signature); err != nil {
		return err
	}

	if !m.Hash.Available() {
		return errors.New("the requested hash function is unavailable")
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig)
}

func (p *Plugin) backgroundRefresh(ctx context.Context) {
	time.Sleep(p.refreshTime)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := p.fetchPublicKeys(ctx); err != nil {
				fmt.Printf("ERROR: Plugin.fetchPublicKeys: %s\n", err)
				time.Sleep(p.refreshTimeError)
			} else {
				time.Sleep(p.refreshTime)
			}
		}
	}
}

const (
	bearerToken  = "Bearer"
	prefixBearer = bearerToken + " "
)

var (
	ErrHeaderAuthMissing   = errors.New("missing authorization header")
	ErrHeaderAuthMalformed = errors.New("malformed authorization header")
	ErrTokenInvalidFormat  = "invalid token format"
	ErrIssuerInvalid       = errors.New("issuer does not match")
)

const (
	discoveryEndpoint = "/.well-known/openid-configuration"
)

type discoveryConfig struct {
	// Issuer is the identifier of the OP and is used in the tokens as `iss` claim.
	Issuer string `json:"issuer"`

	// JwksURI is the URL of the JSON Web Key Set. This site contains the signing keys that RPs can use to validate the signature.
	// It may also contain the OP's encryption keys that RPs can use to encrypt request to the OP.
	JwksURI string `json:"jwks_uri,omitempty"`
}

// rawJSONWebKey represents a public or private key in JWK format, used for parsing/serializing.
type rawJSONWebKey struct {
	Use string `json:"use,omitempty"`
	Kty string `json:"kty,omitempty"`
	Kid string `json:"kid,omitempty"`
	Crv string `json:"crv,omitempty"`
	Alg string `json:"alg,omitempty"`
	K   []byte `json:"k,omitempty"`
	X   []byte `json:"x,omitempty"`
	Y   []byte `json:"y,omitempty"`
	N   []byte `json:"n,omitempty"`
	E   []byte `json:"e,omitempty"`
	D   []byte `json:"d,omitempty"`
	P   []byte `json:"p,omitempty"`
	Q   []byte `json:"q,omitempty"`
}

// jsonWebKeySet represents a JWK Set object.
type jsonWebKeySet struct {
	Keys []jsonWebKey `json:"keys"`
}

func (p *Plugin) fetchPublicKeys(ctx context.Context) error {
	c := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.issuer+discoveryEndpoint, nil)
	if err != nil {
		return fmt.Errorf("new discovery request: %w", err)
	}

	response, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("get discovery: %w", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll discovery: %w", err)
	}

	var cfg discoveryConfig
	if err = json.Unmarshal(body, &cfg); err != nil {
		return fmt.Errorf("json.Unmarshal discovery: %w", err)
	}

	if cfg.JwksURI == "" {
		return errors.New("could not fetch JWKS URI")
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, cfg.JwksURI, nil)
	if err != nil {
		return fmt.Errorf("new jwks request: %w", err)
	}

	response, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("get jwks: %w", err)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll jwks: %w", err)
	}

	var jwksKeys jsonWebKeySet
	if err = json.Unmarshal(body, &jwksKeys); err != nil {
		return fmt.Errorf("json.Unmarshal jwks: %w", err)
	}

	fmt.Printf("DEBUG: JwksURI (%s) body: %s\n", cfg.JwksURI, string(body))

	if len(jwksKeys.Keys) > 1 {
		return fmt.Errorf("multiple public keys is not supported: %+v", jwksKeys.Keys)
	}

	key := jwksKeys.Keys[0]
	if key.Algorithm != jwt.SigningMethodRS256.Alg() {
		return fmt.Errorf("only RS256 algorithm is supported: %+v", key)
	}

	if key.Use != "sig" {
		return fmt.Errorf("only sig use is supported: %+v", key)
	}

	if cfg.Issuer != p.issuer {
		fmt.Printf("ERROR: %s: discovery: %s plugin: %s\n", ErrIssuerInvalid, cfg.Issuer, p.issuer)
		return ErrIssuerInvalid
	}

	p.publicKey = key.Key

	return nil
}

func (p *Plugin) extractToken(request *http.Request) (string, error) {
	authHeader, ok := request.Header["Authorization"]
	if !ok {
		return "", ErrHeaderAuthMissing
	}
	auth := authHeader[0]
	if !strings.HasPrefix(auth, prefixBearer) {
		return "", ErrHeaderAuthMalformed
	}
	parts := strings.Split(auth[len(prefixBearer):], ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("%s: should have 3 parts", ErrTokenInvalidFormat)
	}

	return auth[len(prefixBearer):], nil
}

// jsonWebKey represents a public or private key in JWK format.
type jsonWebKey struct {
	// Cryptographic key, asymmetric key.
	Key *rsa.PublicKey
	// Key identifier, parsed from `kid` header.
	KeyID string
	// Key algorithm, parsed from `alg` header.
	Algorithm string
	// Key use, parsed from `use` header.
	Use string
}

func (k *jsonWebKey) MarshalJSON() ([]byte, error) {
	var raw *rawJSONWebKey

	raw = fromRsaPublicKey(k.Key)
	raw.Kid = k.KeyID
	raw.Alg = k.Algorithm
	raw.Use = k.Use

	return json.Marshal(raw)
}

func (k *jsonWebKey) UnmarshalJSON(data []byte) error {
	var raw rawJSONWebKey
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.Kty != "RSA" {
		return fmt.Errorf("unknown json web key type '%s'", raw.Kty)
	}

	key, err := raw.rsaPublicKey()
	if err != nil {
		return fmt.Errorf("rsaPublicKey: %w", err)
	}

	*k = jsonWebKey{Key: key, KeyID: raw.Kid, Algorithm: raw.Alg, Use: raw.Use}
	return nil
}

func fromRsaPublicKey(pub *rsa.PublicKey) *rawJSONWebKey {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(pub.E))
	return &rawJSONWebKey{
		Kty: "RSA",
		N:   pub.N.Bytes(),
		E:   bytes.TrimLeft(data, "\x00"),
	}
}

func (k rawJSONWebKey) rsaPublicKey() (*rsa.PublicKey, error) {
	if k.N == nil || k.E == nil {
		return nil, fmt.Errorf("invalid RSA key, missing n/e values")
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(k.N),
		E: int(new(big.Int).SetBytes(k.E).Int64()),
	}, nil
}
