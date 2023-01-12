package dao

import (
	"context"
	"goportscan/common/redis"
	"sync"
	"time"
)

var Redis = &_Redis{}

const (
	queryTimeout = 3 * time.Second
)

type _Redis struct {
	mutex sync.Mutex
}

func (r *_Redis) LPush(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	return conn.LPush(key, value).Err()
}

func (r *_Redis) BRPop(key string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	messages, err := conn.BRPop(timeout, key).Result()
	if nil != err {
		return "", err
	} else {
		return messages[1], nil
	}
}
