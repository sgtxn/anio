package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"anio/config"
	"anio/providers/anilist"
	anilistConsts "anio/providers/anilist/consts"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	log.Info().Msg("Loading config")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't load config.")
	}

	log.Info().Msgf("Hello, %s. This line demonstrates that config was successfully loaded.", cfg.Name)

	authCtx, cancel := context.WithTimeout(ctx, anilistConsts.AnilistRequestTimeout)
	defer cancel()
	auth, err := anilist.Authenticate(authCtx, &cfg.AnilistConfig.Auth)
	if err != nil {
		log.Fatal().Err(err).Msgf("Anilist authentication failure: %s", err)
	}

	sampleQuery := `query ($id: Int) {
		Media (id: $id, type: ANIME) {
		  id
		  title {
			romaji
			english
			native
		  }
		}
	  }`

	sampleVariables := map[string]interface{}{"id": 15125}
	anilistApiURL := "https://graphql.anilist.co"

	requestPayload := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     sampleQuery,
		Variables: sampleVariables,
	}

	reqBody, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't marshal the request payload to a json: %s", err)
	}

	fmt.Println(string(reqBody))

	resp, err := auth.Client.Post(anilistApiURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal().Err(err).Msgf("Anilist sample query bad response: %s", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatal().Err(err).Msgf("Anilist sample query expected a 200 response, but got %d with body \n%s", resp.StatusCode, respBody)
	}

	log.Info().Msgf("Got an Anilist response with this body: %s", respBody)
}
