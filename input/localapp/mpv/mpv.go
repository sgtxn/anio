package mpv

import (
	"regexp"

	"anio/input/localapp/windowtitle"
)

type Config struct {
	Enabled       bool `json:"enabled"`
	UseJSONRPCAPI bool `json:"useJsonRpcApi,omitempty"`
}

func (cfg *Config) GetProcessPollerConfig() windowtitle.PolledAppConfig {
	return windowtitle.PolledAppConfig{
		AppName:            "mpv",
		AppExecutable:      "mpv.exe",
		FilenameMatchRegex: regexp.MustCompile(`(.+) - mpv`),
	}
}
