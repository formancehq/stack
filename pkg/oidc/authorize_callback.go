package oidc

import (
	"context"
	"net/http"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/delegatedauth"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
)

func authorizeCallbackHandler(
	provider op.OpenIDProvider,
	storage Storage,
	relyingParty rp.RelyingParty,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO: error handling
		state, err := delegatedauth.DecodeDelegatedState(r.URL.Query().Get("state"))
		if err != nil {
			panic(err)
		}

		authRequest, err := storage.FindAuthRequest(context.Background(), state.AuthRequestID)
		if err != nil {
			panic(err)
		}

		tokens, err := rp.CodeExchange(context.Background(), r.URL.Query().Get("code"), relyingParty)
		if err != nil {
			panic(err)
		}

		user, err := storage.FindUserBySubject(r.Context(), tokens.IDTokenClaims.GetSubject())
		if err != nil {
			user = &auth.User{
				ID:      uuid.NewString(),
				Subject: tokens.IDTokenClaims.GetSubject(),
				Email:   tokens.IDTokenClaims.GetEmail(),
			}
			if err := storage.SaveUser(r.Context(), *user); err != nil {
				panic(err)
			}
		}

		authRequest.UserID = user.ID

		if err := storage.UpdateAuthRequest(r.Context(), *authRequest); err != nil {
			panic(err)
		}

		w.Header().Set("Location", op.AuthCallbackURL(provider)(state.AuthRequestID))
		w.WriteHeader(http.StatusFound)
	}
}
