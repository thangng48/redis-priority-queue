package redis_priority_queue

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var pqueue = getRedisPqueue()

func getRedisPqueue() *redisPQueue {
	var redisUrl = os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}

	var (
		rcl = redis.NewClient(&redis.Options{
			Addr: redisUrl,
		})
		queueKey = "queue"
	)

	if err := rcl.Del(context.TODO(), queueKey).Err(); err != nil {
		panic(err)
	}

	return &redisPQueue{
		rcl:   rcl,
		queue: queueKey,
	}
}

func TestRedisPQueue(t *testing.T) {
	err := pqueue.Push(context.TODO(), Element{Name: "t3", Score: 3},
		Element{Name: "t4", Score: 3},
		Element{Name: "t2", Score: 2},
		Element{Name: "t5", Score: 5})
	assert.NoError(t, err)
	size, err := pqueue.Size(context.TODO())
	assert.NoError(t, err)
	assert.EqualValues(t, 4, size)
	elements, err := pqueue.Get(context.TODO(), 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(elements))
	assert.EqualValues(t, 2, elements[0].Score)
	assert.Equal(t, "t2", elements[0].Name)
	assert.EqualValues(t, 3, elements[1].Score)
	assert.Equal(t, "t3", elements[1].Name)

	// update element's priority
	err = pqueue.Push(context.TODO(), Element{Name: "t2", Score: 1})
	assert.NoError(t, err)
	size, err = pqueue.Size(context.TODO())
	assert.NoError(t, err)
	assert.EqualValues(t, 4, size)
	elements, err = pqueue.Get(context.TODO(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(elements))
	assert.EqualValues(t, 1, elements[0].Score)
	assert.Equal(t, "t2", elements[0].Name)

	// pop element
	elements, err = pqueue.Pop(context.TODO(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(elements))
	assert.EqualValues(t, 1, elements[0].Score)
	assert.Equal(t, "t2", elements[0].Name)
	size, err = pqueue.Size(context.TODO())
	assert.NoError(t, err)
	assert.EqualValues(t, 3, size)
}
