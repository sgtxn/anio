package config

import (
	"regexp"

	"anio/pkg/duration"
)

type InputsConfig struct {
	WebListener  *WebListenerConfig `json:"webListener,omitempty"`
	WebPollers   *WebPollersConfig  `json:"webPollers,omitempty"`
	LocalPollers *LocalAppConfig    `json:"localPollers,omitempty"`
}

type (
	WebListenerConfig struct{}
	WebPollersConfig  struct{}
)

type LocalAppConfig struct {
	PollingInterval duration.Duration `json:"pollingInterval"`
	MpvConfig       *MpvConfig        `json:"mpvConfig"`
}

type PolledAppConfig struct {
	AppName            string
	AppExecutable      string
	FilenameMatchRegex *regexp.Regexp
}
