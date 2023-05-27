package config

import "anio/providers/anilist/consts"

type Config struct {
	Auth AuthConfig `json:"auth"`
}

type AuthConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func GetDefaultConfig() *Config {
	return &Config{
		Auth: AuthConfig{
			ClientID:     consts.ClientID,
			ClientSecret: consts.ClientSecret,
		},
	}
}
