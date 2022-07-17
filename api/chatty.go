package api

import (
	"chatty/middleware"
	"chatty/pkg/client/keycloak"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type Service interface {
	HandleChat(ctx context.Context, ws *websocket.Conn, chatID uuid.UUID)
	ListUsers(firstName, lastName, email, search string) (keycloak.UserList, error)
	CreateChat(current, other uuid.UUID) (uuid.UUID, error)
}

type ChattyHandler struct {
	service Service
}

func NewChattyHandler(service Service) *ChattyHandler {
	return &ChattyHandler{
		service: service,
	}
}

func (h *ChattyHandler) createChat(ctx *gin.Context) {
	// TODO
	currentUser, _ := ctx.Get(middleware.CtxUserKey)
	otherUserIDParam, _ := ctx.Params.Get("userID")

	currentUserID := currentUser.(keycloak.User)
	otherUserID, err := uuid.Parse(otherUserIDParam)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	chatID, err := h.service.CreateChat(currentUserID.ID, otherUserID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, chatID)
}

func (h *ChattyHandler) handleChat(ctx *gin.Context) {
	websocket.Handler(func(ws *websocket.Conn) {
		chatIDParam, ok := ctx.Params.Get("chatID")
		if !ok {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		chatID, err := uuid.Parse(chatIDParam)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		h.service.HandleChat(ctx, ws, chatID)
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
