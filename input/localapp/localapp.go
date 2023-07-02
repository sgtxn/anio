package localapp

import (
	"context"

	"anio/config/inputs"
	"anio/input/localapp/windowtitle"
	"anio/input/shared"
)

type CustomPoller interface {
	Start(context.Context)
}

type Poller struct {
	WindowTitlePoller *windowtitle.Poller
	CustomPollers     []CustomPoller
}

func New(cfg *inputs.LocalAppConfig, outputChan chan<- shared.InputFileInfo) *Poller {
	windowTitlePoller := windowtitle.New(cfg.PollingInterval.Duration, outputChan)

	if cfg.MpvConfig != nil && cfg.MpvConfig.Enabled {
		windowTitlePoller.AddApplication(cfg.MpvConfig.GetProcessPollerConfig())
	}

	processInfoPoller := Poller{
		WindowTitlePoller: windowTitlePoller,
	}

	return &processInfoPoller
}

func (poller *Poller) Start(ctx context.Context) {
	if poller.WindowTitlePoller != nil {
		go poller.WindowTitlePoller.Start(ctx)
	}

	for _, poller := range poller.CustomPollers {
		go poller.Start(ctx)
	}
}
