package anilist

import (
	"context"

	"anio/shared"
)

func (anilist *Client) GetUserList() ([]shared.UserAnimeInfo, error) {
	var query userListQuery

	if err := anilist.graphqlClient.Query(
		context.Background(),
		&query,
		map[string]any{"userId": anilist.userID},
	); err != nil {
		return nil, err
	}

	result := make([]shared.UserAnimeInfo, 0, len(query.MediaListCollection.Lists))
	for _, medialist := range query.MediaListCollection.Lists {
		for _, item := range medialist.Entries {
			result = append(result, shared.UserAnimeInfo{
				UserID: int64(anilist.userID),
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

func (anilist *Client) FindAnimeByTitle(title string, pageID int) ([]shared.AnimeInfo, error) {
	var query findAnimeByTitleQuery

	if err := anilist.graphqlClient.Query(
		context.Background(),
		&query,
		map[string]any{
			"title":  title,
			"pageId": pageID,
		},
	); err != nil {
		return nil, err
	}

	result := make([]shared.AnimeInfo, 0, len(query.Page.Media))
	for _, item := range query.Page.Media {
		result = append(result, shared.AnimeInfo{
			SourceID: anilistSourceID,
			ID:       int64(item.ID),
			Title:    string(item.Title.Romaji),
			Episodes: uint16(item.Episodes),
		})
	}

	return result, nil
}
