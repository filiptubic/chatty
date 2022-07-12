package middleware

import (
	"chatty/config"
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

const AuthorizationCookie = "jwt_token"

type AuthMiddleware struct {
	cfg               *config.Config
	OIDCConfig        *oidc.Config
	OIDCTokenVerifier *oidc.IDTokenVerifier
}

func NewAuthMiddleware(cfg *config.Config) (*AuthMiddleware, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.Auth.Issuer)
	if err != nil {
		return nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.Auth.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	return &AuthMiddleware{
		cfg:               cfg,
		OIDCConfig:        oidcConfig,
		OIDCTokenVerifier: verifier,
	}, nil
}

func (m *AuthMiddleware) Middleware(ctx *gin.Context) {
	r := ctx.Request

	rawAccessToken := r.Header.Get("Authorization")

	if rawAccessToken == "" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := m.OIDCTokenVerifier.Verify(ctx, parts[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	var claims map[string]interface{}

	err = token.Claims(&claims)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.Set("user", claims)

	ctx.Next()
}
