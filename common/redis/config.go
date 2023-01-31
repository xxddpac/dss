package redis

import (
	"fmt"
	"runtime"
)

const (
	MaxDialTimeout  = 1000
	MaxReadTimeout  = 1000
	MaxWriteTimeout = 3000
	MaxPoolTimeout  = 2
	MinIdleConnS    = 3
	MaxRetries      = 1
)

type Config struct {
	Network              string
	Addr                 string
	Passwd               string
	DB                   int
	DialTimeout          int
	ReadTimeout          int
	WriteTimeout         int
	PoolSize             int
	PoolTimeout          int
	MinIdleConnS         int
	MaxRetries           int
	TraceIncludeNotFound bool
}

func (c *Config) Name() string {
	return fmt.Sprintf("%s(%s/%d)", c.Network, c.Addr, c.DB)
}

func (c *Config) FillWithDefaults() {
	if c.DialTimeout <= 0 {
		c.DialTimeout = MaxDialTimeout
	}

	if c.ReadTimeout <= 0 {
		c.ReadTimeout = MaxReadTimeout
	}

	if c.WriteTimeout <= 0 {
		c.WriteTimeout = MaxWriteTimeout
	}

	if c.PoolSize <= 0 {
		c.PoolSize = 10 * runtime.NumCPU()
	}

	if c.PoolTimeout <= 0 {
		c.PoolTimeout = MaxPoolTimeout
	}

	if c.MinIdleConnS <= 0 {
		c.MinIdleConnS = MinIdleConnS
	}

	if c.MaxRetries < 0 {
		c.MaxRetries = MaxRetries
	}
}
