package middleware

import (
	"chatty/pkg/auth"
	"chatty/pkg/client/keycloak"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	CtxUserKey          = "user"
)

func AuthMiddleware(auth *auth.Authenticator, keycloak *keycloak.Keycloak) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r := ctx.Request

		rawAccessToken := r.Header.Get(AuthorizationHeader)

		if rawAccessToken == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := auth.Authenticate(ctx, parts[1])
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		// TODO
		var claims map[string]interface{}

		err = token.Claims(&claims)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		username, ok := claims["preferred_username"]
		if !ok {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		users, err := keycloak.ListUsers("", "", "", "", username.(string))
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ctx.Set(CtxUserKey, users[0])

		ctx.Next()
	}
}
