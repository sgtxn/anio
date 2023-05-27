package input

import (
	"context"

	"anio/input/localapp"
	"anio/input/shared"

	"github.com/rs/zerolog/log"
)

type Config struct {
	WebListener  any              `json:"webListener,omitempty"`
	WebPollers   any              `json:"webPollers,omitempty"`
	LocalPollers *localapp.Config `json:"localPollers,omitempty"`
}

type Input struct {
	LocalAppPoller localapp.LocalProcessPoller
}

func New(cfg *Config, outputChan chan<- shared.InputFileInfo) *Input {
	return &Input{
		LocalAppPoller: *localapp.New(cfg.LocalPollers, outputChan),
	}
}

func (input *Input) Start(ctx context.Context) {
	log.Info().Msg("initializing inputs...")
	input.LocalAppPoller.Start(ctx)
}
