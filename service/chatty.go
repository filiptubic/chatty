package service

import (
	"context"

	"github.com/coreos/go-oidc"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

var (
	channel = make(map[*websocket.Conn]struct{})
)

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
	channel[ws] = struct{}{}

	var msg Message

	// first authenticate
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

	err = websocket.JSON.Send(ws, &Message{
		Event: "joined",
		Data:  "hello",
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to send join event")
	}

	// authenticated
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
		err = websocket.JSON.Send(ws, &msg)
		if err != nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}
