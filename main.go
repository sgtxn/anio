package main

import (
	"context"
	"os"
	"time"

	"anio/config"
	"anio/pkg/userdirs"
	"anio/providers/anilist"
	"anio/shared"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func main() {
	ctx := context.Background()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})

	log.Info().Msg("loading config...")
	anioCfgDir, err := userdirs.GetProjectConfigDirectory()
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't get project config file path")
	}

	cfg, err := config.Load(anioCfgDir)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't load config")
	}

	log.Info().Msg("config loaded successfully")

	zerolog.SetGlobalLevel(cfg.LogLevel)

	// inputsChan := make(chan shared.PlaybackFileInfo)
	// inputs := input.New(cfg.Inputs, inputsChan)
	// inputs.Start(ctx)

	// parsedTitlesChan := make(chan shared.PlaybackAnimeDetails)
	// titleParser := titleparser.New(inputsChan, parsedTitlesChan)
	// titleParser.Start()

	// anioDataDir, err := userdirs.GetProjectDataDirectory()
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("couldn't get user data directory")
	// }

	// _, err = storage.New(anioDataDir)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("couldn't initialize the database")
	// }

	// for data := range parsedTitlesChan {
	// 	log.Info().Any("parsedTitle", data).Send()
	// }

	client, err := anilist.New(ctx, &cfg.Outputs.Anilist.Auth)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't connect to anilist")
	}

	err = client.UpdateAnime(&shared.AnimeUpdateParams{
		ID:       null.IntFrom(147103),
		Progress: null.IntFrom(3),
	})
	if err != nil {
		log.Error().Err(err).Send()
	}
}
