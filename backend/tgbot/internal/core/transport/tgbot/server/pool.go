package core_tgbot_server

import (
	"context"
	"errors"
	"sync"
	"time"

	core_logger "github.com/wydentis/vykladna/shared/core/logger"
	core_tgbot_middleware "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/middleware"
	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

var ErrQueueFull = errors.New("telegram update queue is full")

type pool struct {
	handle core_tgbot_middleware.HandlerFunc
	log    *core_logger.Logger
	shards []chan core_tgbot_types.Update
	wg     sync.WaitGroup
}

func newPool(h core_tgbot_middleware.HandlerFunc, log *core_logger.Logger, workers, queueSize int) *pool {
	if workers < 1 {
		workers = 1
	}
	if queueSize < 1 {
		queueSize = 1
	}

	p := &pool{
		handle: h,
		log:    log,
		shards: make([]chan core_tgbot_types.Update, workers),
	}

	for i := range p.shards {
		p.shards[i] = make(chan core_tgbot_types.Update, queueSize)
	}

	return p
}

func (p *pool) start() {
	for _, ch := range p.shards {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for u := range ch {
				p.run(u)
			}
		}()
	}
}

func (p *pool) run(u core_tgbot_types.Update) {
	l := p.log.With("update_id", u.UpdateID)
	ctx := l.ToContext(context.Background())

	defer func() {
		if r := recover(); r != nil {
			l.Error("panic in telegram handler", "panic", r)
		}
	}()

	if err := p.handle(ctx, u); err != nil {
		l.Error("telegram handler failed", "err", err)
	}
}

func (p *pool) stop(timeout time.Duration) {
	for _, ch := range p.shards {
		close(ch)
	}

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		p.log.Warn("telegram pool drain timed out")
	}
}

func (p *pool) Submit(u core_tgbot_types.Update) error {
	select {
	case p.shards[shardOf(u.Payload, len(p.shards))] <- u:
		return nil
	default:
		return ErrQueueFull
	}
}

type chatter interface {
	ChatID() int64
}

func shardOf(payload any, n int) int {
	if c, ok := payload.(chatter); ok {
		id := c.ChatID()
		if id < 0 {
			id = -id
		}

		return int(id % int64(n))
	}

	return 0
}
