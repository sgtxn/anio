package localapp

import (
	"context"

	"anio/config"
	"anio/input/localapp/windowtitle"
	"anio/input/shared"
)

type CustomPoller interface {
	Start(context.Context)
}

type LocalProcessPoller struct {
	WindowTitlePoller *windowtitle.Poller
	CustomPollers     []CustomPoller
}

func New(cfg *config.LocalAppConfig, outputChan chan<- shared.InputFileInfo) *LocalProcessPoller {
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
