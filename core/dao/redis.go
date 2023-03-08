package dao

import (
	"context"
	"dss/common/redis"
	"dss/common/utils"
	_redis "github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
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

func (r *_Redis) Lock(key string, timeout time.Duration) (identifier string, b bool) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	identifier = uuid.NewV4().String()
	end := time.Now().Add(queryTimeout)
	for time.Now().Before(end) {
		if conn.SetNX(key, identifier, timeout).Val() {
			b = true
			return
		} else if conn.TTL(key).Val().Seconds() == -1 {
			conn.Expire(key, timeout)
		}
		time.Sleep(time.Microsecond)
	}
	return
}

func (r *_Redis) UnLock(key, identifier string) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	script := _redis.NewScript(`
		if redis.call("get",KEYS[1]) == ARGV[1] then
			return redis.call("del",KEYS[1])
		else
			return 0
		end
	`)
	return script.Run(conn, []string{key}, identifier).Err()
}

func (r *_Redis) LuaRun(key string, count int) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	script := _redis.NewScript(`
		if redis.call("exists",KEYS[1]) == 0 then
			redis.call("set",KEYS[1],ARGV[1])
			redis.call("expire",KEYS[1],60)
		else
			local value = redis.call("get",KEYS[1])
			redis.call("set",KEYS[1],tonumber(value) + ARGV[1])
			redis.call("expire",KEYS[1],60)
		end
	`)
	if err := script.Run(conn, []string{key}, count).Err(); err != nil && err != _redis.Nil {
		return err
	}
	return nil
}

func (r *_Redis) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	resp := conn.Get(key)
	return resp.Result()
}

func (r *_Redis) IsExistsOrGet(key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()
	conn := redis.GetConn(ctx)
	script := _redis.NewScript(` 
		if redis.call("exists",KEYS[1]) == 0 then
			return "0"
		else
			local value = redis.call("get",KEYS[1])
			return value
		end
	`)
	resp, err := script.Run(conn, []string{key}).Result()
	if err != nil {
		return -1, err
	}
	return utils.StrToInt(resp.(string)), nil
}
