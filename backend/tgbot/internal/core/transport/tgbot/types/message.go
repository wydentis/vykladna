package core_tgbot_types

import "strings"

type Base struct {
	MessageID       int64 `json:"message_id"`
	MessageThreadID int64 `json:"message_thread_id"`
	From            *User `json:"from"`
	Chat            Chat  `json:"chat"`
	Date            int64 `json:"date"`
}

func (b Base) ChatID() int64 {
	return b.Chat.ID
}

func (b Base) UserID() int64 {
	if b.From == nil {
		return 0
	}
	return b.From.ID
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type TextMessage struct {
	Base
	Text     string          `json:"text"`
	Entities []MessageEntity `json:"entities"`
}

func (m TextMessage) Command() (name, botUsername, args string, ok bool) {
	if len(m.Entities) == 0 {
		return
	}
	e := m.Entities[0]
	if e.Type != "bot_command" || e.Offset != 0 || e.Length > len(m.Text) {
		return
	}
	cmd := strings.TrimPrefix(m.Text[:e.Length], "/")
	name, botUsername, _ = strings.Cut(cmd, "@")
	return name, botUsername, strings.TrimSpace(m.Text[e.Length:]), true
}

func (m TextMessage) IsCommand() bool {
	name, _, _, ok := m.Command()
	return (ok && name != "")
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size"`
}

type PhotoMessage struct {
	Base
	Photo        []PhotoSize `json:"photo"`
	Caption      string      `json:"caption"`
	MediaGroupID string      `json:"media_group_id"`
}

func (m PhotoMessage) Largest() PhotoSize {
	var best PhotoSize
	for _, p := range m.Photo {
		if p.Width*p.Height > best.Width*best.Height {
			best = p
		}
	}
	return best
}

type Document struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type DocumentMessage struct {
	Base
	Document Document `json:"document"`
	Caption  string   `json:"caption"`
}

type Sticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Emoji        string `json:"emoji"`
	SetName      string `json:"set_name"`
}

type StickerMessage struct {
	Base
	Sticker Sticker `json:"sticker"`
}

type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type VoiceMessage struct {
	Base
	Voice Voice `json:"voice"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationMessage struct {
	Base
	Location Location `json:"location"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int64  `json:"user_id"`
}

type ContactMessage struct {
	Base
	Contact Contact `json:"contact"`
}

type OtherMessage struct{ Base }
