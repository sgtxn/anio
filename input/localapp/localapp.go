package localapp

import (
	"context"
	"time"

	"anio/input/localapp/mpv"
	"anio/input/localapp/windowtitle"
	"anio/input/shared"
)

type CustomPoller interface {
	Start(context.Context)
}

type Config struct {
	PollingInterval time.Duration
	MpvConfig       *mpv.Config
}

type LocalProcessPoller struct {
	WindowTitlePoller *windowtitle.Poller
	CustomPollers     []CustomPoller
}

func New(cfg *Config, outputChan chan<- shared.InputFileInfo) *LocalProcessPoller {
	windowTitlePoller := windowtitle.New(cfg.PollingInterval, outputChan)

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
