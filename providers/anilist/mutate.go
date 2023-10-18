package anilist

import (
	"context"
	"errors"
	"fmt"

	"anio/shared"
)

func (anilist *Client) UpdateAnime(params *shared.AnimeUpdateParams) error {
	if params == nil {
		return errors.New("nothing to do")
	}

	requestVars := make(map[string]any)
	requestVars["id"] = params.ID.Int64

	if params.Progress.Valid {
		requestVars["progress"] = params.Progress.Int64
	} else {
		requestVars["progress"] = (*int)(nil)
	}

	if params.Status.Valid {
		requestVars["status"] = params.Status.String
	} else {
		requestVars["status"] = (*MediaListStatus)(nil)
	}

	if params.Score.Valid {
		requestVars["scoreRaw"] = params.Score.Int64
	} else {
		requestVars["scoreRaw"] = 0
	}

	err := anilist.graphqlClient.Mutate(context.Background(), &updateEntryQuery{}, requestVars)
	if err != nil {
		return fmt.Errorf("error calling an anilist update: %w", err)
	}
	return nil
}
