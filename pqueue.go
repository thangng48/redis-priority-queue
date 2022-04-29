package redis_priority_queue

import (
	"context"
)

type PQueue interface {
	// Push insert element with priority to the queue
	Push(ctx context.Context, members ...Element) error

	// Pop pull "highest" priority elements from the queue
	// higher priority definition is up to the implementation
	Pop(ctx context.Context, quantity int64) ([]*Element, error)

	// BPop wait util receiving any Element
	BPop(ctx context.Context) (*Element, error)

	// Size returns queue's cardinality
	Size(ctx context.Context) (int64, error)

	// Get retrieve "highest" priority elements from the queue
	Get(ctx context.Context, quantity int64) ([]*Element, error)
}
