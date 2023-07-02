package inputs

import (
	"regexp"
	"runtime"
)

type MpvConfig struct {
	Enabled       bool `json:"enabled"`
	UseJSONRPCAPI bool `json:"useJsonRpcApi,omitempty"`
}

func (cfg *MpvConfig) GetProcessPollerConfig() PolledAppConfig {
	polledAppCfg := PolledAppConfig{
		AppName:            "mpv",
		FilenameMatchRegex: regexp.MustCompile(`(.+) - mpv`),
	}

	switch runtime.GOOS {
	case "windows":
		polledAppCfg.AppExecutable = "mpv.exe"
	case "linux":
		polledAppCfg.AppExecutable = "mpv"
	}

	return polledAppCfg
}
