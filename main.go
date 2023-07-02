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

	outputChan := make(chan shared.InputFileInfo)
	inputBlock := input.New(cfg.Inputs, outputChan)
	inputBlock.Start(ctx)

	for data := range outputChan {
		log.Info().Any("data", data).Send()
	}
}
