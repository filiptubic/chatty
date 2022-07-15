package auth

import (
	"chatty/config"
	"context"

	"github.com/coreos/go-oidc"
)

type Authenticator struct {
	cfg               *config.Config
	OIDCConfig        *oidc.Config
	OIDCTokenVerifier *oidc.IDTokenVerifier
}

func NewAuthenticator(cfg *config.Config) (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.Auth.Issuer)
	if err != nil {
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.Auth.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	return &Authenticator{
		cfg:               cfg,
		OIDCConfig:        oidcConfig,
		OIDCTokenVerifier: verifier,
	}, nil
}

func (a *Authenticator) Authenticate(ctx context.Context, token string) (*oidc.IDToken, error) {
	return a.OIDCTokenVerifier.Verify(ctx, token)
}
