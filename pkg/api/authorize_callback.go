package api

import (
	"context"
	"net/http"

	auth "github.com/numary/auth/pkg"
	"github.com/numary/auth/pkg/delegatedauth"
	"github.com/numary/auth/pkg/storage"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
)

func authorizeCallbackHandler(
	provider op.OpenIDProvider,
	storage storage.Storage,
	relyingParty rp.RelyingParty,
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

		tokens, err := rp.CodeExchange(context.Background(), r.URL.Query().Get("code"), relyingParty)
		if err != nil {
			panic(err)
		}

		userInfo, err := rp.Userinfo(tokens.AccessToken, "Bearer", tokens.IDTokenClaims.GetSubject(), relyingParty)
		if err != nil {
			panic(err)
		}

		user, err := storage.FindUserByEmail(r.Context(), userInfo.GetEmail())
		if err != nil {
			user = &auth.User{
				Subject: userInfo.GetSubject(),
				Email:   userInfo.GetEmail(),
			}
			if err := storage.CreateUser(r.Context(), user); err != nil {
				panic(err)
			}
		}

		if err := storage.MarkAuthRequestAsDone(r.Context(), authRequest.GetID(), user.Subject); err != nil {
			panic(err)
		}

		w.Header().Set("Location", op.AuthCallbackURL(provider)(state.AuthRequestID))
		w.WriteHeader(http.StatusFound)
	}
}
