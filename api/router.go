package api

import (
	"chatty/auth"
	"chatty/config"
	"chatty/middleware"
	"chatty/service"

	"github.com/gin-gonic/gin"
)

func Mount(cfg *config.Config) (*gin.Engine, error) {
	authenticator, err := auth.NewAuthenticator(cfg)
	if err != nil {
		return nil, err
	}
	engine := gin.New()
	engine.Use(middleware.CorsMiddleware())

	v1 := engine.Group("/v1")
	v1.Use(middleware.AuthMiddleware(authenticator))

	chattyService := service.NewChattyService(authenticator)
	chattyHandler := NewChattyHandler(chattyService)

	engine.GET("/ws", chattyHandler.handleWS)

	return engine, nil
}
