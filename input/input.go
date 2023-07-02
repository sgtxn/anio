package input

import (
	"context"

	"anio/config"
	"anio/input/localapp"
	"anio/input/shared"

	"github.com/rs/zerolog/log"
)

type Input struct {
	LocalAppPoller *localapp.Poller
}

func New(cfg *config.InputsConfig, outputChan chan<- shared.InputFileInfo) *Input {
	var input Input

	if cfg == nil {
		return &input
	}

	if cfg.LocalPollers != nil {
		input.LocalAppPoller = localapp.New(cfg.LocalPollers, outputChan)
	}

	return &input
}

func (input *Input) Start(ctx context.Context) {
	log.Info().Msg("initializing inputs...")

	if input.LocalAppPoller != nil {
		input.LocalAppPoller.Start(ctx)
	}
}
