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
		AppName:              "mpv",
		AppExecutableWindows: "mpv.exe",
		AppExecutableLinux:   "mpv",
		FilenameMatchRegex:   regexp.MustCompile(`(.+) - mpv`),
	}
}
