package middleware

import (
	"chatty/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const AuthorizationCookie = "jwt_token"

type AuthMiddleware struct {
	cfg               *config.Config
	Provider          *oidc.Provider
	OAuth2Config      *oauth2.Config
	OIDCConfig        *oidc.Config
	OIDCTokenVerifier *oidc.IDTokenVerifier
}

func NewAuthMiddleware(cfg *config.Config) (*AuthMiddleware, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.Auth.Issuer)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.Auth.ClientID,
		ClientSecret: cfg.Auth.SecretID,
		RedirectURL:  cfg.Auth.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.Auth.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	return &AuthMiddleware{
		cfg:               cfg,
		Provider:          provider,
		OAuth2Config:      &oauth2Config,
		OIDCConfig:        oidcConfig,
		OIDCTokenVerifier: verifier,
	}, nil
}

func (m *AuthMiddleware) Middleware(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	// TODO encrypt state
	originalURL := fmt.Sprintf("http://%s%s", r.Host, r.URL.String())

	rawAccessToken := ""
	if cookie, err := r.Cookie(AuthorizationCookie); err == nil {
		rawAccessToken = cookie.Value
	}

	if rawAccessToken == "" {
		http.Redirect(w, r, m.OAuth2Config.AuthCodeURL(originalURL), http.StatusFound)
		return
	}

	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := m.OIDCTokenVerifier.Verify(ctx, parts[1])
	if err != nil {
		http.Redirect(w, r, m.OAuth2Config.AuthCodeURL(originalURL), http.StatusFound)
		return
	}

	ctx.Next()
}

func (m *AuthMiddleware) Callback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request
		w := ctx.Writer

		oauth2Token, err := m.OAuth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		idToken, err := m.OIDCTokenVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		originalURLString := r.URL.Query().Get("state")

		originalURL, err := url.Parse(originalURLString)
		if err != nil {
			http.Error(w, "Invalid state: "+err.Error(), http.StatusInternalServerError)
		}

		ctx.SetCookie(
			AuthorizationCookie,
			resp.OAuth2Token.AccessToken,
			10,
			originalURL.Path,
			originalURL.Hostname(),
			true,
			true,
		)
		ctx.Redirect(http.StatusFound, originalURLString)
	}
}
