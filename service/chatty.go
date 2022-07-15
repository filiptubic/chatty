package service

import (
	repository "chatty/pkg/repostory"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

var (
	m              sync.Mutex
	channel        = make(map[*websocket.Conn]struct{})
	defaultChannel repository.Channel
)

func joinChannel(ws *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	channel[ws] = struct{}{}
}

func exitChannel(ws *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	delete(channel, ws)
	log.Info().Msg("client exited")
}

type Repository interface {
	SaveInHistory(ch repository.Channel, m repository.Message)
	LoadHistory(ch repository.Channel) []repository.Message
}

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*oidc.IDToken, error)
}

type ChattyService struct {
	auth Authenticator
	repo Repository
}

func NewChattyService(auth Authenticator, repo Repository) (*ChattyService, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	defaultChannel = repository.Channel(uid)

	return &ChattyService{
		auth: auth,
		repo: repo,
	}, nil
}

func (s *ChattyService) Join(ws *websocket.Conn) {
	joinChannel(ws)
	go func() {
		history := s.repo.LoadHistory(defaultChannel)
		for _, msg := range history {
			err := websocket.JSON.Send(ws, msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to send message")
			}
		}
	}()
}

func (s *ChattyService) Send(ws *websocket.Conn, m repository.Message) error {
	m.ID = uuid.New()
	m.SendAt = time.Now().UTC()

	if m.Event == repository.MessageEvent {
		s.repo.SaveInHistory(defaultChannel, m)
	}

	for client := range channel {
		go func(client *websocket.Conn) {
			err := websocket.JSON.Send(client, &m)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to send msg to client")
			}
		}(client)
	}
	return nil
}

func (s *ChattyService) Route(ws *websocket.Conn, m repository.Message) {
	switch m.Event {
	case repository.MessageEvent, repository.TypingEvent:
		err := s.Send(ws, m)
		if err != nil {
			log.Error().Err(err).Msg("failed to send message")
		}
	default:
		log.Info().Interface("msg", m).Msg("discarding unkown message")
	}
}

func (s *ChattyService) HandleWS(ctx context.Context, ws *websocket.Conn) {
	var msg repository.Message

	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to receive from ws")
	}
	if msg.Event != repository.AuthEvent {
		_ = websocket.JSON.Send(ws, repository.ErrorMessage(errors.New("invalid auth")))
		return
	}

	_, err = s.auth.Authenticate(ctx, msg.Data.(string))
	if err != nil {
		_ = websocket.JSON.Send(ws, repository.ErrorMessage(errors.New("invalid auth")))
		return
	}

	s.Join(ws)
	defer exitChannel(ws)

	for {
		var msg repository.Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			return
		}
		s.Route(ws, msg)
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}
