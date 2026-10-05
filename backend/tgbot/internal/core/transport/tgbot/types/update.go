package core_tgbot_types

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Update struct {
	UpdateID int64
	Payload  any
}

func (u Update) Type() reflect.Type {
	return reflect.TypeOf(u.Payload)
}

var kinds = map[string]reflect.Type{
	"message":                    nil,
	"edited_message":             reflect.TypeFor[EditedMessage](),
	"channel_post":               reflect.TypeFor[ChannelPost](),
	"edited_channel_post":        reflect.TypeFor[EditedChannelPost](),
	"business_message":           reflect.TypeFor[BusinessMessage](),
	"edited_business_message":    reflect.TypeFor[EditedBusinessMessage](),
	"guest_message":              reflect.TypeFor[GuestMessage](),
	"callback_query":             reflect.TypeFor[CallbackQuery](),
	"business_connection":        reflect.TypeFor[BusinessConnection](),
	"deleted_business_messages":  reflect.TypeFor[DeletedBusinessMessages](),
	"message_reaction":           reflect.TypeFor[MessageReaction](),
	"message_reaction_count":     reflect.TypeFor[MessageReactionCount](),
	"inline_query":               reflect.TypeFor[InlineQuery](),
	"chosen_inline_result":       reflect.TypeFor[ChosenInlineResult](),
	"shipping_query":             reflect.TypeFor[ShippingQuery](),
	"pre_checkout_query":         reflect.TypeFor[PreCheckoutQuery](),
	"purchased_paid_media":       reflect.TypeFor[PurchasedPaidMedia](),
	"poll":                       reflect.TypeFor[PollUpdate](),
	"poll_answer":                reflect.TypeFor[PollAnswer](),
	"my_chat_member":             reflect.TypeFor[MyChatMember](),
	"chat_member":                reflect.TypeFor[ChatMember](),
	"chat_join_request":          reflect.TypeFor[ChatJoinRequest](),
	"chat_boost":                 reflect.TypeFor[ChatBoost](),
	"removed_chat_boost":         reflect.TypeFor[RemovedChatBoost](),
	"managed_bot":                reflect.TypeFor[ManagedBot](),
	"subscription":               reflect.TypeFor[Subscription](),
	"stopped_message_generation": reflect.TypeFor[StoppedMessageGeneration](),
}

type MessageEnvelope struct {
	Message any
}

func (e *MessageEnvelope) UnmarshalJSON(b []byte) error {
	m, err := decodeMessage(b)
	if err != nil {
		return err
	}
	e.Message = m
	return nil
}

func (e MessageEnvelope) ChatID() int64 {
	if c, ok := e.Message.(interface{ ChatID() int64 }); ok {
		return c.ChatID()
	}
	return 0
}

type (
	EditedMessage         struct{ MessageEnvelope }
	ChannelPost           struct{ MessageEnvelope }
	EditedChannelPost     struct{ MessageEnvelope }
	BusinessMessage       struct{ MessageEnvelope }
	EditedBusinessMessage struct{ MessageEnvelope }
	GuestMessage          struct{ MessageEnvelope }
)

type (
	BusinessConnection       struct{ json.RawMessage }
	DeletedBusinessMessages  struct{ json.RawMessage }
	MessageReaction          struct{ json.RawMessage }
	MessageReactionCount     struct{ json.RawMessage }
	InlineQuery              struct{ json.RawMessage }
	ChosenInlineResult       struct{ json.RawMessage }
	ShippingQuery            struct{ json.RawMessage }
	PreCheckoutQuery         struct{ json.RawMessage }
	PurchasedPaidMedia       struct{ json.RawMessage }
	PollUpdate               struct{ json.RawMessage }
	PollAnswer               struct{ json.RawMessage }
	MyChatMember             struct{ json.RawMessage }
	ChatMember               struct{ json.RawMessage }
	ChatJoinRequest          struct{ json.RawMessage }
	ChatBoost                struct{ json.RawMessage }
	RemovedChatBoost         struct{ json.RawMessage }
	ManagedBot               struct{ json.RawMessage }
	Subscription             struct{ json.RawMessage }
	StoppedMessageGeneration struct{ json.RawMessage }
)

func (u *Update) UnmarshalJSON(b []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(b, &obj); err != nil {
		return fmt.Errorf("decode update envelope: %w", err)
	}
	if v, ok := obj["update_id"]; ok {
		if err := json.Unmarshal(v, &u.UpdateID); err != nil {
			return fmt.Errorf("decode update_id: %w", err)
		}
		delete(obj, "update_id")
	}

	for key, body := range obj {
		if !present(body) {
			continue
		}
		t, ok := kinds[key]
		if !ok {
			return fmt.Errorf("unknown update kind %q", key)
		}

		var err error
		if t == nil {
			u.Payload, err = decodeMessage(body)
		} else {
			u.Payload, err = decodeType(t, body)
		}

		if err != nil {
			return fmt.Errorf("decode %q update: %w", key, err)
		}

		return nil
	}

	return nil
}

var messageKinds = []struct {
	key string
	typ reflect.Type
}{
	{"animation", reflect.TypeFor[OtherMessage]()},
	{"live_photo", reflect.TypeFor[OtherMessage]()},
	{"venue", reflect.TypeFor[OtherMessage]()},
	{"sticker", reflect.TypeFor[StickerMessage]()},
	{"voice", reflect.TypeFor[VoiceMessage]()},
	{"photo", reflect.TypeFor[PhotoMessage]()},
	{"document", reflect.TypeFor[DocumentMessage]()},
	{"location", reflect.TypeFor[LocationMessage]()},
	{"contact", reflect.TypeFor[ContactMessage]()},
	{"text", reflect.TypeFor[TextMessage]()},
}

func decodeMessage(b []byte) (any, error) {
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(b, &keys); err != nil {
		return nil, fmt.Errorf("decode message keys: %w", err)
	}
	for _, k := range messageKinds {
		if present(keys[k.key]) {
			return decodeType(k.typ, b)
		}
	}

	return decodeType(reflect.TypeFor[OtherMessage](), b)
}

func decodeType(t reflect.Type, b []byte) (any, error) {
	p := reflect.New(t)
	if err := json.Unmarshal(b, p.Interface()); err != nil {
		return nil, fmt.Errorf("decode %s: %w", t, err)
	}
	return p.Elem().Interface(), nil
}

func present(r json.RawMessage) bool {
	return len(r) > 0 && string(r) != "null"
}
