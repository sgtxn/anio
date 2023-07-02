package config

type AnilistConfig struct {
	Auth AnilistAuthConfig `json:"auth"`
}

type AnilistAuthConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}
