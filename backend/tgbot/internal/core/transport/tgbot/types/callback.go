package core_tgbot_types

type CallbackQuery struct {
	ID      string `json:"id"`
	From    User   `json:"from"`
	Message *Base  `json:"message"`
	Data    string `json:"data"`
}

func (c CallbackQuery) ChatID() int64 {
	if c.Message != nil {
		return c.Message.Chat.ID
	}
	return c.From.ID
}

func (c CallbackQuery) UserID() int64 {
	return c.From.ID
}

func (c CallbackQuery) Accessible() bool {
	return c.Message != nil && c.Message.Date != 0
}
