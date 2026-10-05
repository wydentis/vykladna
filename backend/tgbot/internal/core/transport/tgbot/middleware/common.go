package core_tgbot_middleware

import (
	"context"
	"fmt"
	"time"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

func Trace() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, u core_tgbot_types.Update) error {
			log := core_logger.FromContext(ctx)
			before := time.Now()
			log.Debug(
				"-> incoming telegram update",
				"timestamp", before.UTC(),
				"payload", fmt.Sprint(u.Type()),
			)

			err := next(ctx, u)

			now := time.Now()

			log.Debug(
				"<- done telegram update",
				"timestamp", now.UTC(),
				"latency", now.Sub(before),
				"err", err,
			)

			return err
		}
	}
}
