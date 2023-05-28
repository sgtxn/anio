package config

import (
	"regexp"

	"anio/pkg/duration"
)

type InputsConfig struct {
	WebListener  any             `json:"webListener,omitempty"`
	WebPollers   any             `json:"webPollers,omitempty"`
	LocalPollers *LocalAppConfig `json:"localPollers,omitempty"`
}

type LocalAppConfig struct {
	PollingInterval duration.Duration `json:"pollingInterval"`
	MpvConfig       *MpvConfig        `json:"mpvConfig"`
}

type PolledAppConfig struct {
	AppName            string
	AppExecutable      string
	FilenameMatchRegex *regexp.Regexp
}
