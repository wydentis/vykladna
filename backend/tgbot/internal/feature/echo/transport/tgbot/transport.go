package echo_transport_tgbot

import (
	"context"
	"fmt"

	core_tgbot_client "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/client"
	core_tgbot_server "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/server"
	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

type Client interface {
	SendMessage(ctx context.Context, p core_tgbot_client.SendMessageParams) error
	AnswerCallbackQuery(ctx context.Context, callbackQueryID, text string) error
}

type EchoTransport struct {
	client Client
}

func NewEchoTransport(client Client) *EchoTransport {
	return &EchoTransport{client: client}
}

func (t *EchoTransport) Routes() []core_tgbot_server.Route {
	return []core_tgbot_server.Route{
		core_tgbot_server.NewRoute(t.OnText),
		core_tgbot_server.NewRoute(t.OnCallback),
	}
}

func (t *EchoTransport) OnText(ctx context.Context, m core_tgbot_types.TextMessage) error {
	return t.client.SendMessage(ctx, core_tgbot_client.SendMessageParams{
		ChatID:   m.ChatID(),
		ThreadID: m.MessageThreadID,
		Text:     "Hi! Your chat id: " + fmt.Sprint(m.ChatID()) + " " + m.Text,
	})
}

func (t *EchoTransport) OnCallback(ctx context.Context, c core_tgbot_types.CallbackQuery) error {
	return t.client.AnswerCallbackQuery(ctx, c.ID, "got: "+c.Data)
}
