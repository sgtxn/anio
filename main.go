package main

import (
	"context"
	"os"
	"time"

	"anio/config"
	"anio/input"
	"anio/input/shared"
	"anio/storage"
	"anio/titleparser"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	projectFolderName = "anio"
)

func main() {
	ctx := context.Background()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})

	log.Info().Msg("loading config...")
	cfgFilePath, err := config.GetConfigFilePath(projectFolderName)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't get config file path")
	}

	cfg, err := config.Load(cfgFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't load config")
	}

	log.Info().Msg("config loaded successfully")

	zerolog.SetGlobalLevel(cfg.LogLevel)

	inputsChan := make(chan shared.PlaybackFileInfo)
	inputs := input.New(cfg.Inputs, inputsChan)
	inputs.Start(ctx)

	parsedTitlesChan := make(chan shared.PlaybackAnimeDetails)
	titleParser := titleparser.New(inputsChan, parsedTitlesChan)
	titleParser.Start()

	_, err = storage.New(".")
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't initialize the database")
	}

	for data := range parsedTitlesChan {
		log.Info().Any("parsedTitle", data).Send()
	}
}
