package gateway_plugin_auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
)

func TestPlugin_ServeHTTP(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == discoveryEndpoint {
			cfg := discoveryConfig{
				Issuer:  "http://" + l.Addr().String(),
				JwksURI: "http://" + l.Addr().String() + "/keys",
			}
			by, _ := json.Marshal(cfg)
			_, _ = w.Write(by)
		} else if r.URL.String() == "/keys" {
			keys := jsonWebKeySet{Keys: []jsonWebKey{
				{
					KeyID:     "id",
					Key:       publicKey,
					Algorithm: signingMethodDefault.Alg(),
					Use:       "sig",
				},
			}}
			by, _ := json.Marshal(keys)
			_, _ = w.Write(by)
		}
	}))

	_ = ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	ctx := context.Background()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	config := CreateConfig()
	config.Issuer = "http://" + l.Addr().String()
	handler, err := New(ctx, next, config, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ERR missing header", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatal(recorder.Code)
		}
		b := recorder.Body.String()
		if !strings.HasPrefix(b, ErrHeaderAuthMissing.Error()) {
			t.Fatal(b)
		}
	})

	t.Run("ERR malformed header", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
		req.Header.Set("Authorization", "malformed")
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatal(recorder.Code)
		}
		b := recorder.Body.String()
		if !strings.HasPrefix(b, ErrHeaderAuthMalformed.Error()) {
			t.Fatal(b)
		}
	})

	t.Run("ERR invalid format token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
		req.Header.Set("Authorization", prefixBearer+"invalid format")
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatal(recorder.Code)
		}
		b := recorder.Body.String()
		if !strings.HasPrefix(b, ErrTokenInvalidFormat) {
			t.Fatal(b)
		}
	})

	t.Run("ERR unverified token", func(t *testing.T) {
		invalidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ"
		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
		req.Header.Set("Authorization", prefixBearer+invalidToken)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(recorder, req)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatal(recorder.Code)
		}
		b := recorder.Body.String()
		if !strings.HasPrefix(b, "unverified JWT token") {
			t.Fatal(b)
		}
	})

	t.Run("verified token", func(t *testing.T) {
		jwtToken := jwt.NewWithClaims(signingMethodDefault, jwt.StandardClaims{
			Issuer: l.Addr().String(),
		})

		jwtTokenString, err := signedString(jwtToken, privateKey, signingMethodDefault)
		if err != nil {
			t.Fatal(fmt.Errorf("signedString: %w", err))
		}

		parser := new(jwt.Parser)
		token, parts, err := parser.ParseUnverified(jwtTokenString, jwt.MapClaims{})
		if err != nil {
			t.Fatal(fmt.Errorf("parser.ParseUnverified: %w", err))
		}

		vErr := &jwt.ValidationError{}
		token.Signature = parts[2]
		if err = verifyRSA(strings.Join(parts[0:2], "."), token.Signature, publicKey, signingMethodDefault); err != nil {
			vErr.Inner = err
			vErr.Errors |= jwt.ValidationErrorSignatureInvalid
			t.Fatal(fmt.Errorf("verifyRSA: %w", err))
		}

		recorder := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
		req.Header.Set("Authorization", prefixBearer+jwtTokenString)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(recorder, req)
		b := recorder.Body.String()
		if recorder.Code != http.StatusOK {
			t.Fatal(recorder.Code, b)
		}
	})
}

// Get the complete, signed token
func signedString(t *jwt.Token, privateKey *rsa.PrivateKey, m *jwt.SigningMethodRSA) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = t.SigningString(); err != nil {
		return "", fmt.Errorf("t.SigningString: %w", err)
	}

	if sig, err = signRSA(sstr, privateKey, m); err != nil {
		return "", fmt.Errorf("signRSA: %w", err)
	}

	return strings.Join([]string{sstr, sig}, "."), nil
}

func signRSA(signingString string, privateKey *rsa.PrivateKey, m *jwt.SigningMethodRSA) (string, error) {
	if !m.Hash.Available() {
		return "", errors.New("the requested hash function is unavailable")
	}

	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, m.Hash, hasher.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("rsa.SignPKCS1v15: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(sigBytes), nil
}
