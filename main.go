package main

import (
	"context"
	"os"
	"time"

	"anio/config"
	"anio/input"
	"anio/input/shared"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})

	ctx := context.Background()

	log.Info().Msg("loading config...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't load config")
	}

	log.Info().Msg("config loaded successfully")

	outputChan := make(chan shared.InputFileInfo)
	inputBlock := input.New(cfg.Inputs, outputChan)
	inputBlock.Start(ctx)

	for data := range outputChan {
		log.Info().Any("data", data).Send()
	}
}
