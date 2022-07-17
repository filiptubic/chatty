package api

import (
	"chatty/config"
	"chatty/middleware"
	"chatty/pkg/auth"
	"chatty/pkg/client/keycloak"
	"chatty/pkg/repostory/inmemory"
	"chatty/service"

	"github.com/gin-gonic/gin"
)

func Mount(cfg *config.Config) (*gin.Engine, error) {
	keycloak := keycloak.New(cfg)

	authenticator, err := auth.NewAuthenticator(cfg)
	if err != nil {
		return nil, err
	}
	repo := inmemory.NewInMemoryRepository()
	chattyService, err := service.NewChattyService(
		authenticator,
		repo,
		keycloak,
	)
	if err != nil {
		return nil, err
	}
	chattyHandler := NewChattyHandler(chattyService)

	engine := gin.New()
	engine.Use(middleware.CorsMiddleware())

	v1 := engine.Group("/v1")
	v1.Use(middleware.AuthMiddleware(authenticator))
	v1.GET("/users", chattyHandler.listUsers)

	engine.GET("/ws/:chatID", chattyHandler.handleChat)

	return engine, nil
}
