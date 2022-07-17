package api

import (
	"chatty/pkg/client/keycloak"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type Service interface {
	HandleWS(ctx context.Context, ws *websocket.Conn)
	ListUsers(firstName, lastName, email, search string) (keycloak.UserList, error)
}

type ChattyHandler struct {
	service Service
}

func NewChattyHandler(service Service) *ChattyHandler {
	return &ChattyHandler{
		service: service,
	}
}

func (h *ChattyHandler) handleWS(ctx *gin.Context) {
	websocket.Handler(func(ws *websocket.Conn) {
		h.service.HandleWS(ctx, ws)
		defer ws.Close()
	}).ServeHTTP(ctx.Writer, ctx.Request)
}

func (h *ChattyHandler) listUsers(ctx *gin.Context) {
	email := ctx.Query("email")
	firstName := ctx.Query("first_name")
	lastName := ctx.Query("last_name")
	search := ctx.Query("search")
	users, err := h.service.ListUsers(firstName, lastName, email, search)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to fetch users")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, users)
}
