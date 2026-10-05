package core_tgbot_client

import "context"

type SendMessageParams struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ThreadID  int64  `json:"message_thread_id,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func (c *Client) SendMessage(ctx context.Context, p SendMessageParams) error {
	return c.call(ctx, "sendMessage", p, nil)
}

func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackQueryID, text string) error {
	return c.call(ctx, "answerCallbackQuery", map[string]any{
		"callback_query_id": callbackQueryID,
		"text":              text,
	}, nil)
}
