package anilist

import (
	"context"

	"anio/shared"
)

func (anilist *Client) GetUserList() ([]shared.UserAnimeInfo, error) {
	var result []shared.UserAnimeInfo

	var query userListQuery

	if err := anilist.graphqlClient.Query(
		context.Background(),
		&query,
		map[string]any{"userId": 123},
	); err != nil {
		return nil, err
	}

	for _, medialist := range query.MediaListCollection.Lists {
		for _, item := range medialist.Entries {
			result = append(result, shared.UserAnimeInfo{
				AnimeInfo: shared.AnimeInfo{
					SourceID: anilistSourceID,
					ID:       int64(item.ID),
					Title:    string(item.Media.Title.Romaji),
					Episodes: uint16(item.Media.Episodes),
				},
				Progress:   uint16(item.Progress),
				Score:      uint8(item.Score),
				Rewatches:  uint16(item.Repeat),
				StartDate:  item.StartedAt.ToDateTime(),
				FinishDate: item.CompletedAt.ToDateTime(),
				Status:     string(medialist.Status),
			})
		}
	}

	return result, nil
}
