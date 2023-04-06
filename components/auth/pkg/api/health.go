package api

//func delegatedOIDCServerAvailability(rp rp.RelyingParty) health.NamedCheck {
//	return health.NewNamedCheck("Delegated OIDC server", health.CheckFn(func(ctx context.Context) error {
//		_, err := client.Discover(rp.Issuer(), http.DefaultClient)
//		return err
//	}))
//}
