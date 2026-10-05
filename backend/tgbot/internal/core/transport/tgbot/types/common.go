package core_tgbot_types

type User struct {
	ID           int64   `json:"id"`
	IsBot        bool    `json:"is_bot"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Username     *string `json:"username"`
	LanguageCode *string `json:"language_code"`
}

type Chat struct {
	ID               int64   `json:"id"`
	Type             string  `json:"type"`
	Title            *string `json:"title"`
	Username         *string `json:"username"`
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	IsForum          *bool   `json:"is_forum"`
	IsDirectMessages *bool   `json:"is_direct_messages"`
}
