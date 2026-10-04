package core_tgbot_middleware

import (
	"context"
	"slices"

	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

type HandlerFunc func(ctx context.Context, u core_tgbot_types.Update) error
type Middleware func(HandlerFunc) HandlerFunc

func Chain(h HandlerFunc, m ...Middleware) HandlerFunc {
	for _, v := range slices.Backward(m) {
		h = v(h)
	}

	return h
}
