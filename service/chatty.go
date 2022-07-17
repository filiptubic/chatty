package service

import (
	"chatty/pkg/client/keycloak"
	"chatty/pkg/model"
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
	defaultChannel model.Channel
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
	SaveInHistory(ch model.Channel, m model.Message)
	LoadHistory(ch model.Channel) []model.Message
}

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*oidc.IDToken, error)
}

type UserClient interface {
	ListUsers(firstName, lastName, email string) (keycloak.UserList, error)
}

type ChattyService struct {
	auth  Authenticator
	repo  Repository
	users UserClient
}

func NewChattyService(auth Authenticator, repo Repository, users UserClient) (*ChattyService, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	defaultChannel = model.Channel(uid)

	return &ChattyService{
		auth:  auth,
		repo:  repo,
		users: users,
	}, nil
}

func (s *ChattyService) ListUsers(firstName, lastName, email string) (keycloak.UserList, error) {
	return s.users.ListUsers(firstName, lastName, email)
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

func (s *ChattyService) Send(ws *websocket.Conn, m model.Message) error {
	m.ID = uuid.New()
	m.SendAt = time.Now().UTC()

	if m.Event == model.MessageEvent {
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

func (s *ChattyService) Route(ws *websocket.Conn, m model.Message) {
	switch m.Event {
	case model.MessageEvent, model.TypingEvent:
		err := s.Send(ws, m)
		if err != nil {
			log.Error().Err(err).Msg("failed to send message")
		}
	default:
		log.Info().Interface("msg", m).Msg("discarding unkown message")
	}
}

func (s *ChattyService) HandleWS(ctx context.Context, ws *websocket.Conn) {
	var msg model.Message

	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to receive from ws")
	}
	if msg.Event != model.AuthEvent {
		_ = websocket.JSON.Send(ws, model.ErrorMessage(errors.New("invalid auth")))
		return
	}

	_, err = s.auth.Authenticate(ctx, msg.Data.(string))
	if err != nil {
		_ = websocket.JSON.Send(ws, model.ErrorMessage(errors.New("invalid auth")))
		return
	}

	s.Join(ws)
	defer exitChannel(ws)

	for {
		var msg model.Message
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
