package service

import (
	"context"
	"sync"

	"github.com/coreos/go-oidc"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

var (
	m       sync.Mutex
	channel = make(map[*websocket.Conn]struct{})
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

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*oidc.IDToken, error)
}

type ChattyService struct {
	auth Authenticator
}

func NewChattyService(auth Authenticator) *ChattyService {
	return &ChattyService{
		auth: auth,
	}
}

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func (s *ChattyService) HandleWS(ctx context.Context, ws *websocket.Conn) {
	var msg Message

	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to receive from ws")
	}
	if msg.Event != "auth" {
		_ = websocket.JSON.Send(ws, &Message{
			Event: "error",
			Data:  "invalid auth",
		})
		return
	}

	_, err = s.auth.Authenticate(ctx, msg.Data)
	if err != nil {
		_ = websocket.JSON.Send(ws, &Message{
			Event: "error",
			Data:  err.Error(),
		})
		return
	}

	joinChannel(ws)
	defer exitChannel(ws)

	for {
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			return
		}
		for client := range channel {
			err = websocket.JSON.Send(client, &msg)
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to send msg to client")
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}
