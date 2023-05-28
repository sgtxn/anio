package input

import (
	"context"

	"anio/config"
	"anio/input/localapp"
	"anio/input/shared"

	"github.com/rs/zerolog/log"
)

type Input struct {
	LocalAppPoller localapp.LocalProcessPoller
}

func New(cfg *config.InputsConfig, outputChan chan<- shared.InputFileInfo) *Input {
	return &Input{
		LocalAppPoller: *localapp.New(cfg.LocalPollers, outputChan),
	}
}

func (input *Input) Start(ctx context.Context) {
	log.Info().Msg("initializing inputs...")
	input.LocalAppPoller.Start(ctx)
}
