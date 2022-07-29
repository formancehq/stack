package api

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	auth "github.com/numary/auth/pkg"
	"github.com/numary/auth/pkg/delegatedauth"
	"github.com/numary/auth/pkg/storage"
	"github.com/zitadel/oidc/pkg/op"
	"golang.org/x/oauth2"
)

func authorizeCallbackHandler(
	provider op.OpenIDProvider,
	storage storage.Storage,
	delegatedOAuth2Config *oauth2.Config,
	delegatedOIDCProvider *oidc.Provider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		state, err := delegatedauth.DecodeDelegatedState(r.URL.Query().Get("state"))
		if err != nil {
			panic(err)
		}

		authRequest, err := storage.AuthRequestByID(context.Background(), state.AuthRequestID)
		if err != nil {
			panic(err)
		}

		token, err := delegatedOAuth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			panic(err)
		}

		idToken, err := delegatedOIDCProvider.Verifier(&oidc.Config{
			ClientID:        "gateway",
			SkipIssuerCheck: true, //TODO: Only on dev as url doesn't match in this case
		}).Verify(context.Background(), token.Extra("id_token").(string))
		if err != nil {
			panic(err)
		}

		claims := &delegatedauth.Claims{}
		if err := idToken.Claims(&claims); err != nil {
			panic(err)
		}

		user, err := storage.FindUserByEmail(r.Context(), claims.Email)
		if err != nil {
			user = &auth.User{
				Subject: claims.Subject,
				Email:   claims.Email,
			}
			if err := storage.CreateUser(r.Context(), user); err != nil {
				panic(err)
			}
		}

		if err := storage.MarkAuthRequestAsDone(r.Context(), authRequest.GetID(), user.Subject); err != nil {
			//TODO: Handle error
			panic(err)
		}

		vars := mux.Vars(r)
		vars["id"] = state.AuthRequestID
		mux.SetURLVars(r, vars)

		op.AuthorizeCallback(w, r, provider)
	}
}
