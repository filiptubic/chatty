package server

import (
	"chatty/config"
	"chatty/middleware"
	"fmt"
	"net/http"
	"time"

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
	s.engine.GET("/callback", authMiddleware.Callback())

	s.engine.LoadHTMLGlob("web/*.html")
	s.engine.Static("/assets", "./web/assets")

	s.engine.GET("/redirect", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	v1 := s.engine.Group("/v1")
	v1.Use(authMiddleware.Middleware)
	v1.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	s.engine.GET("/ws", func(ctx *gin.Context) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				err := websocket.Message.Send(ws, "Hello World!")
				if err != nil {
					log.Error().Err(err).Msg("failed to send msg")
				}
				time.Sleep(time.Second * 2)
			}
		}).ServeHTTP(ctx.Writer, ctx.Request)
	})

	return s.engine.Run(fmt.Sprintf(":%d", s.cfg.Server.Port))
}
