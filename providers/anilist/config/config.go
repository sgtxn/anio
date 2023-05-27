package config

import "anio/providers/anilist/consts"

type Config struct {
	Auth AuthConfig
}

type AuthConfig struct {
	ClientID     string
	ClientSecret string
}

func GetDefaultConfig() *Config {
	return &Config{
		Auth: AuthConfig{
			ClientID:     consts.ClientID,
			ClientSecret: consts.ClientSecret,
		},
	}
}
