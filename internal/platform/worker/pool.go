package worker

import (
	"context"
)

// Pool is a simple worker group that runs a number of
// tasks at a configured concurrency
type Pool struct {
	taskCh     chan func(ctx context.Context)
	cancel     context.CancelFunc
	workersNum int
}

// NewPool initializes a new pool with a given concurrency
func NewPool(workersNum int) Pool {
	return Pool{
		taskCh:     make(chan func(ctx context.Context), workersNum),
		workersNum: workersNum,
	}
}

// Run spawns all workers within the pool
func (p Pool) Run(ctx context.Context) {
	_, p.cancel = context.WithCancel(ctx)

	for i := 0; i < p.workersNum; i++ {
		go func(ch chan func(ctx context.Context)) {
			for t := range ch {
				t(ctx)
			}
		}(p.taskCh)
	}
}

// Enqueue adds new task to the tasks queue
func (p Pool) Enqueue(cb func(ctx context.Context)) {
	p.taskCh <- cb
}

// Close cancels all tasks and closes tasks channel
func (p Pool) Close() {
	p.cancel()
	close(p.taskCh)
}
