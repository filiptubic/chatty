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
	m        sync.Mutex
	channels = make(map[uuid.UUID]map[*websocket.Conn]struct{})
	// TODO
	directChannel = make(map[uuid.UUID]map[uuid.UUID]uuid.UUID)
)

type Repository interface {
	SaveInHistory(ch model.Channel, m model.Message)
	LoadHistory(ch model.Channel) []model.Message
}

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*oidc.IDToken, error)
}

type UserClient interface {
	ListUsers(firstName, lastName, email, search, username string) (keycloak.UserList, error)
}

type ChattyService struct {
	auth  Authenticator
	repo  Repository
	users UserClient
}

func NewChattyService(auth Authenticator, repo Repository, users UserClient) (*ChattyService, error) {
	return &ChattyService{
		auth:  auth,
		repo:  repo,
		users: users,
	}, nil
}

func (s *ChattyService) ListUsers(firstName, lastName, email, search string) (keycloak.UserList, error) {
	return s.users.ListUsers(firstName, lastName, email, search, "")
}

func (s *ChattyService) createChat(user1, user2 uuid.UUID, newChatID uuid.UUID) (uuid.UUID, error) {
	// TODO mutex

	chats, ok := directChannel[user1]
	if !ok {
		chats = make(map[uuid.UUID]uuid.UUID)
		directChannel[user1] = chats
	}

	chatID, ok := chats[user2]
	if !ok {
		chatID = newChatID
		chats[user2] = chatID
	}

	return chatID, nil
}

func (s *ChattyService) CreateChat(current, other uuid.UUID) (uuid.UUID, error) {
	newChatID := uuid.New()

	chatID1, err := s.createChat(current, other, newChatID)
	if err != nil {
		return uuid.UUID{}, err
	}

	chatID2, err := s.createChat(other, current, newChatID)
	if err != nil {
		return uuid.UUID{}, err
	}

	if chatID1 != chatID2 {
		return uuid.UUID{}, errors.New("different chat IDs")
	}

	return chatID1, nil
}

func (s *ChattyService) JoinChat(ws *websocket.Conn, chatID uuid.UUID) {
	m.Lock()
	defer m.Unlock()

	if _, ok := channels[chatID]; !ok {
		channels[chatID] = make(map[*websocket.Conn]struct{})
	}
	channels[chatID][ws] = struct{}{}

	go func() {
		history := s.repo.LoadHistory(model.Channel(chatID))
		for _, msg := range history {
			err := websocket.JSON.Send(ws, msg)
			if err != nil {
				log.Error().Err(err).Msg("failed to send message")
			}
		}
	}()
	log.Info().Str("chatID", chatID.String()).Msg("connected to chat")
}

func (s *ChattyService) ExitChat(ws *websocket.Conn, chatID uuid.UUID) {
	m.Lock()
	defer m.Unlock()

	delete(channels[chatID], ws)
}

func (s *ChattyService) Send(ws *websocket.Conn, chatID uuid.UUID, m model.Message) error {
	m.ID = uuid.New()
	m.SendAt = time.Now().UTC()

	if m.Event == model.MessageEvent {
		s.repo.SaveInHistory(model.Channel(chatID), m)
	}

	for client := range channels[chatID] {
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

func (s *ChattyService) Route(ws *websocket.Conn, chatID uuid.UUID, m model.Message) {
	switch m.Event {
	case model.MessageEvent, model.TypingEvent:
		err := s.Send(ws, chatID, m)
		if err != nil {
			log.Error().Err(err).Msg("failed to send message")
		}
	default:
		log.Info().Interface("msg", m).Msg("discarding unkown message")
	}
}

func (s *ChattyService) Authenticate(ctx context.Context, ws *websocket.Conn) error {
	var msg model.Message

	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		return err
	}
	if msg.Event != model.AuthEvent {
		return err
	}

	_, err = s.auth.Authenticate(ctx, msg.Data.(string))
	if err != nil {
		return err
	}
	return nil
}

func (s *ChattyService) HandleChat(ctx context.Context, ws *websocket.Conn, chatID uuid.UUID) {
	err := s.Authenticate(ctx, ws)
	if err != nil {
		_ = websocket.JSON.Send(ws, model.ErrorMessage(errors.New("invalid auth")))
		return
	}

	s.JoinChat(ws, chatID)
	defer s.ExitChat(ws, chatID)

	for {
		var msg model.Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			return
		}
		s.Route(ws, chatID, msg)
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}
