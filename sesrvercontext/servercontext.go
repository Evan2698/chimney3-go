package servercontext

import (
	"context"
	"sync/atomic"
	"time"
)

type ServerContext interface {
	context.Context
	Interrupted() error
	IsInterrupted() bool
	ClearInterrupted() error
}

type erverContextImpl struct {
	ctx         context.Context
	interrupted int32
}

func NewServerContext(ctx context.Context) ServerContext {
	sc := &erverContextImpl{
		ctx:         ctx,
		interrupted: 0,
	}
	atomic.StoreInt32(&sc.interrupted, 0)
	return sc
}

func (sc *erverContextImpl) Interrupted() error {
	atomic.StoreInt32(&sc.interrupted, 1)
	return nil
}

func (sc *erverContextImpl) IsInterrupted() bool {
	return atomic.LoadInt32(&sc.interrupted) == 1
}

func (sc *erverContextImpl) ClearInterrupted() error {
	atomic.StoreInt32(&sc.interrupted, 0)
	return nil
}

func (sc *erverContextImpl) Deadline() (deadline time.Time, ok bool) {
	return sc.ctx.Deadline()
}

func (sc *erverContextImpl) Done() <-chan struct{} {
	return sc.ctx.Done()
}

func (sc *erverContextImpl) Err() error {
	return sc.ctx.Err()
}

func (sc *erverContextImpl) Value(key interface{}) interface{} {
	return sc.ctx.Value(key)
}
