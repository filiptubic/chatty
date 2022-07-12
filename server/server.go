package server

import (
	"chatty/config"
	"chatty/middleware"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type Server struct {
	cfg    *config.Config
	engine *gin.Engine
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		engine: gin.New(),
	}
}

func (s *Server) Start() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLvl, err := zerolog.ParseLevel(s.cfg.Server.Log.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLvl)

	authMiddleware, err := middleware.NewAuthMiddleware(s.cfg)
	if err != nil {
		return err
	}

	s.engine.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"Origin", "Authorization"},
	}))

	v1 := s.engine.Group("/v1")
	v1.Use(authMiddleware.Middleware)

	// websockets
	type Message struct {
		Event string `json:"event"`
		Data  string `json:"data"`
	}
	s.engine.GET("/ws", func(ctx *gin.Context) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			var msg Message
			// first authenticate
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to receive from ws")
			}
			if msg.Event != "auth" {
				websocket.JSON.Send(ws, &Message{
					Event: "error",
					Data:  "invalid auth",
				})
				return
			}

			_, err = authMiddleware.OIDCTokenVerifier.Verify(ctx, msg.Data)
			if err != nil {
				websocket.JSON.Send(ws, &Message{
					Event: "error",
					Data:  err.Error(),
				})
				return
			}

			websocket.JSON.Send(ws, &Message{
				Event: "joined",
				Data:  "hello",
			})
			// authenticated
			for {
				websocket.JSON.Receive(ws, &msg)
				websocket.JSON.Send(ws, &msg)
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}

			}
		}).ServeHTTP(ctx.Writer, ctx.Request)
	})

	return s.engine.Run(fmt.Sprintf(":%d", s.cfg.Server.Port))
}
