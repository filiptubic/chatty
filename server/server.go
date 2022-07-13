package server

import (
	"chatty/api"
	"chatty/config"
	"fmt"

	"github.com/rs/zerolog"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLvl, err := zerolog.ParseLevel(s.cfg.Server.Log.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLvl)

	engine, err := api.Mount(s.cfg)
	if err != nil {
		return err
	}

	return engine.Run(fmt.Sprintf(":%d", s.cfg.Server.Port))
}
