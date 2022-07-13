package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type Service interface {
	HandleWS(ctx context.Context, ws *websocket.Conn)
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
