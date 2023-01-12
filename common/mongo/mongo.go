package mongo

import (
	"github.com/globalsign/mgo"
	"strings"
	"sync"
	"time"
)

var (
	client *Client
)

func GetConn(name string) *Client {
	return client.C(name)
}

type Client struct {
	session    *mgo.Session
	config     *Config
	collection *mgo.Collection
	mux        sync.RWMutex
}

func Init(config *Config) (err error) {
	config.FillWithDefaults()

	info := &mgo.DialInfo{
		Addrs:          strings.Split(config.Host, ","),
		Timeout:        time.Duration(config.DialTimeout) * time.Second,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		ReplicaSetName: config.ReplicaSetName,
	}

	if len(config.Database) > 0 {
		info.Database = config.Database
	}

	if config.HasAuth() {
		info.Username = config.Auth.User
		info.Password = config.Auth.Passwd

		if len(config.Auth.Database) > 0 {
			info.Database = config.Auth.Database
		}
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return
	}

	switch config.Mode {
	case "Strong":
		session.SetMode(mgo.Strong, true)
	case "Monotonic":
		session.SetMode(mgo.Monotonic, true)
	case "Eventual":
		session.SetMode(mgo.Eventual, true)
	default:
		session.SetMode(mgo.Strong, true)
	}

	session.SetSafe(&mgo.Safe{
		W:        config.SafeWriteAck,
		WTimeout: config.SafeWriteTimeout,
		FSync:    config.SafeJournal,
	})

	session.SetPoolLimit(config.PoolSize)
	session.SetPoolTimeout(time.Duration(config.PoolTimeout) * time.Millisecond)
	session.SetSyncTimeout(time.Duration(config.SyncTimeout) * time.Second)

	client = &Client{
		session: session,
		config:  config,
	}

	return
}

func (c *Client) Database() string {
	return c.config.Database
}

func (c *Client) Session() *mgo.Session {
	return c.session
}

func (c *Client) Collection() *mgo.Collection {
	return c.collection
}

func (c *Client) Close() {
	c.session.Close()
}

func (c *Client) Copy() *Client {
	return &Client{
		session: c.session.Copy(),
		config:  c.config.Copy(),
	}
}

func (c *Client) C(name string) *Client {
	c.mux.Lock()
	defer c.mux.Unlock()

	copied := c.Copy()
	copied.collection = copied.session.DB(c.Database()).C(name)

	return copied
}
