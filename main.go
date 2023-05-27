package main

import (
	"context"
	"fmt"
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
		fmt.Println(data)
	}

	// authCtx, cancel := context.WithTimeout(ctx, anilistConsts.RequestTimeout)
	// defer cancel()
	// auth, err := anilist.Authenticate(authCtx, &cfg.AnilistConfig.Auth)
	// if err != nil {
	// 	log.Fatal().Err(err).Msgf("anilist authentication failure: %s", err)
	// }

	// sampleQuery := `query ($id: Int) {
	// 	Media (id: $id, type: ANIME) {
	// 	  id
	// 	  title {
	// 		romaji
	// 		english
	// 		native
	// 	  }
	// 	}
	//   }`

	// sampleVariables := map[string]interface{}{"id": 15125}
	// anilistApiURL := "https://graphql.anilist.co"

	// requestPayload := struct {
	// 	Query     string                 `json:"query"`
	// 	Variables map[string]interface{} `json:"variables"`
	// }{
	// 	Query:     sampleQuery,
	// 	Variables: sampleVariables,
	// }

	// reqBody, err := json.Marshal(requestPayload)
	// if err != nil {
	// 	log.Fatal().Err(err).Msgf("couldn't marshal the request payload to a json: %s", err)
	// }

	// fmt.Println(string(reqBody))

	// resp, err := auth.Client.Post(anilistApiURL, "application/json", bytes.NewReader(reqBody))
	// if err != nil {
	// 	log.Fatal().Err(err).Msgf("anilist sample query bad response: %s", err)
	// }
	// defer resp.Body.Close()

	// respBody, _ := io.ReadAll(resp.Body)

	// if resp.StatusCode != http.StatusOK {
	// 	log.Fatal().Err(err).Msgf("anilist sample query expected a 200 response, but got %d with body \n%s", resp.StatusCode, respBody)
	// }

	// log.Info().Msgf("got an anilist response with body: %s", respBody)
}
