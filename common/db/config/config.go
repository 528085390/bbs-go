package config

import "temp/common/env"

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	Dbname   string
}

type Config struct {
	Database DatabaseConfig
}

func (c *DatabaseConfig) LoadFromEnv() {
	c.Host = env.GetEnv("DB_HOST", c.Host)
	c.Port = env.GetEnv("DB_PORT", c.Port)
	c.User = env.GetEnv("DB_USER", c.User)
	c.Password = env.GetEnv("DB_PASSWORD", c.Password)
	c.Dbname = env.GetEnv("DB_NAME", c.Dbname)
}

func (c *DatabaseConfig) SetDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == "" {
		c.Port = "5432"
	}
	if c.User == "" {
		c.User = "postgres"
	}
	if c.Password == "" {
		c.Password = "123456"
	}
	if c.Dbname == "" {
		c.Dbname = "bbs-go"
	}
}
