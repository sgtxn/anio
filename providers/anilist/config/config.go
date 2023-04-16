package config

type Config struct {
	Auth AuthConfig
}

type AuthConfig struct {
	ClientID     string
	ClientSecret string
}
