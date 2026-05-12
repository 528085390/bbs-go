package config

type Config struct {
	JwtSecret string
}

func NewConfig() *Config {
	return &Config{}
}
