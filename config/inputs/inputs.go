package inputs

import (
	"regexp"

	"anio/pkg/duration"
)

type Config struct {
	WebListener  *WebListenerConfig `json:"webListener,omitempty"`
	WebPollers   *WebPollersConfig  `json:"webPollers,omitempty"`
	LocalPollers *LocalAppConfig    `json:"local,omitempty"`
}

type (
	WebListenerConfig struct{}
	WebPollersConfig  struct{}
)

type LocalAppConfig struct {
	PollingInterval duration.Duration `json:"pollingInterval"`
	MpvConfig       *MpvConfig        `json:"mpv"`
}

type PolledAppConfig struct {
	AppName            string
	AppExecutable      string
	FilenameMatchRegex *regexp.Regexp
}
