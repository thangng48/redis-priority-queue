package redis_priority_queue

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redisPQueue struct {
	rcl   *redis.Client
	queue string
}

func NewRedisPQueue(rcl *redis.Client, queueKey string) PQueue {
	return &redisPQueue{
		rcl:   rcl,
		queue: queueKey,
	}
}

func (p *redisPQueue) Push(ctx context.Context, members ...Element) error {
	zMembers := make([]*redis.Z, len(members))
	for i, m := range members {
		zMembers[i] = &redis.Z{
			Score:  float64(m.Score),
			Member: m.Name,
		}
	}
	return p.rcl.ZAdd(ctx, p.queue, zMembers...).Err()
}

func (p *redisPQueue) BPop(ctx context.Context) (*Element, error) {
	member := p.rcl.BZPopMin(ctx, 0, p.queue)
	err := member.Err()
	if err != nil {
		return nil, err
	}
	zResult, err := member.Result()
	if err != nil {
		return nil, err
	}
	return &Element{
		Name:  zResult.Member.(string),
		Score: int64(zResult.Score),
	}, nil
}

// Pop pull "highest" elements from the queue
// lower score value means higher priority
func (p *redisPQueue) Pop(ctx context.Context, quantity int64) ([]*Element, error) {
	members := p.rcl.ZPopMin(ctx, p.queue, quantity)
	err := members.Err()
	if err != nil {
		return nil, err
	}
	zResults, err := members.Result()
	if err != nil {
		return nil, err
	}

	elements := make([]*Element, len(zResults))
	for i, z := range zResults {
		elements[i] = &Element{
			Name:  z.Member.(string),
			Score: int64(z.Score),
		}
	}
	return elements, nil
}

func (p *redisPQueue) Size(ctx context.Context) (int64, error) {
	zCount := p.rcl.ZCount(ctx, p.queue, "-inf", "+inf")
	err := zCount.Err()
	if err != nil {
		return 0, err
	}
	return zCount.Result()
}

func (p *redisPQueue) Get(ctx context.Context, quantity int64) ([]*Element, error) {
	members := p.rcl.ZRangeWithScores(ctx, p.queue, 0, quantity-1)
	err := members.Err()
	if err != nil {
		return nil, err
	}
	zResults, err := members.Result()
	if err != nil {
		return nil, err
	}

	qMembers := make([]*Element, len(zResults))
	for i, z := range zResults {
		qMembers[i] = &Element{
			Name:  z.Member.(string),
			Score: int64(z.Score),
		}
	}
	return qMembers, nil
}
