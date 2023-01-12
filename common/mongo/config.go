package mongo

const (
	MaxDialTimeout      = 1
	MaxReadTimeout      = 3
	MaxWriteTimeout     = 5
	MaxPoolSize         = 4096
	MaxPoolTimeout      = 300
	MinSafeWriteAck     = 1
	MaxSafeWriteTimeout = 5000
	MaxSyncTimeout      = 5
)

type AuthConfig struct {
	User     string
	Passwd   string
	Database string
}

func (c *AuthConfig) IsValid() bool {
	return c != nil && c.User != "" && c.Passwd != ""
}

type Config struct {
	Host             string
	Database         string
	Auth             *AuthConfig
	Mode             string
	DialTimeout      int
	ReadTimeout      int
	WriteTimeout     int
	PoolSize         int
	PoolTimeout      int
	SyncTimeout      int
	SafeWriteAck     int
	SafeWriteTimeout int
	SafeJournal      bool
	ReplicaSetName   string
}

func (c *Config) FillWithDefaults() {
	if c == nil {
		return
	}
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
		c.PoolSize = MaxPoolSize
	}

	if c.PoolTimeout <= 0 {
		c.PoolTimeout = MaxPoolTimeout
	}

	if c.SafeWriteAck <= 0 {
		c.SafeWriteAck = MinSafeWriteAck
	}

	if c.SafeWriteTimeout <= 0 {
		c.SafeWriteTimeout = MaxSafeWriteTimeout
	}

	if c.SyncTimeout <= 0 {
		c.SyncTimeout = MaxSyncTimeout
	}
}

func (c *Config) HasAuth() bool {
	return c != nil && c.Auth.IsValid()
}

func (c *Config) Copy() *Config {
	config := *c

	config.Auth = &AuthConfig{}
	if c.HasAuth() {
		config.Auth.User = c.Auth.User
		config.Auth.Passwd = c.Auth.Passwd
		config.Auth.Database = c.Auth.Database
	}

	return &config
}
