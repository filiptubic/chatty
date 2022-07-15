package api

import (
	"chatty/config"
	"chatty/middleware"
	"chatty/pkg/auth"
	"chatty/pkg/repostory/inmemory"
	"chatty/service"

	"github.com/gin-gonic/gin"
)

func Mount(cfg *config.Config) (*gin.Engine, error) {
	authenticator, err := auth.NewAuthenticator(cfg)
	if err != nil {
		return nil, err
	}
	repo := inmemory.NewInMemoryRepository()

	engine := gin.New()
	engine.Use(middleware.CorsMiddleware())

	v1 := engine.Group("/v1")
	v1.Use(middleware.AuthMiddleware(authenticator))

	chattyService, err := service.NewChattyService(authenticator, repo)
	if err != nil {
		return nil, err
	}
	chattyHandler := NewChattyHandler(chattyService)

	engine.GET("/ws", chattyHandler.handleWS)

	return engine, nil
}
