package localapp

import (
	"context"

	"anio/input/localapp/mpv"
	"anio/input/localapp/windowtitle"
	"anio/input/shared"
	"anio/pkg/duration"
)

type CustomPoller interface {
	Start(context.Context)
}

type Config struct {
	PollingInterval duration.Duration `json:"pollingInterval"`
	MpvConfig       *mpv.Config       `json:"mpvConfig"`
}

type LocalProcessPoller struct {
	WindowTitlePoller *windowtitle.Poller
	CustomPollers     []CustomPoller
}

func New(cfg *Config, outputChan chan<- shared.InputFileInfo) *LocalProcessPoller {
	windowTitlePoller := windowtitle.New(cfg.PollingInterval.Duration, outputChan)

	if cfg.MpvConfig != nil && cfg.MpvConfig.Enabled {
		windowTitlePoller.AddApplication(cfg.MpvConfig.GetProcessPollerConfig())
	}

	processInfoPoller := LocalProcessPoller{
		WindowTitlePoller: windowTitlePoller,
	}

	return &processInfoPoller
}

func (poller *LocalProcessPoller) Start(ctx context.Context) {
	if poller.WindowTitlePoller != nil {
		go poller.WindowTitlePoller.Start(ctx)
	}

	for _, poller := range poller.CustomPollers {
		go poller.Start(ctx)
	}
}
