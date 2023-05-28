package config

import (
	"regexp"
)

type MpvConfig struct {
	Enabled       bool `json:"enabled"`
	UseJSONRPCAPI bool `json:"useJsonRpcApi,omitempty"`
}

func (cfg *MpvConfig) GetProcessPollerConfig() PolledAppConfig {
	return PolledAppConfig{
		AppName:            "mpv",
		AppExecutable:      "mpv.exe",
		FilenameMatchRegex: regexp.MustCompile(`(.+) - mpv`),
	}
}
