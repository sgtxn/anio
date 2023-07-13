package input

import (
	"context"

	"anio/config/inputs"
	"anio/input/localapp"
	"anio/input/shared"

	"github.com/rs/zerolog/log"
)

type Input struct {
	LocalAppPoller *localapp.Poller
}

func New(cfg *inputs.Config, outputChan chan<- shared.PlaybackFileInfo) *Input {
	var input Input

	if cfg.LocalPollers != nil {
		input.LocalAppPoller = localapp.New(cfg.LocalPollers, outputChan)
	}

	return &input
}

func (input *Input) Start(ctx context.Context) {
	log.Info().Msg("listening to the inputs...")

	if input.LocalAppPoller != nil {
		input.LocalAppPoller.Start(ctx)
	}
}
