//nolint:revive // can't do without nested structs here
package anilist

import (
	"time"

	"github.com/hasura/go-graphql-client"
	"gopkg.in/guregu/null.v4"
)

type MediaListStatus struct{}

type date struct {
	Year  int
	Month int
	Day   int
}

type animeTitle struct {
	Romaji        graphql.String
	English       graphql.String
	Native        graphql.String
	UserPreferred graphql.String
}

type anilistError struct {
	Message graphql.String
	Status  graphql.Int
}

// -------------------------- queries --------------------------
type userListQuery struct {
	MediaListCollection struct {
		Lists []struct {
			Name                 graphql.String
			IsCustomList         graphql.Boolean
			IsSplitCompletedList graphql.Boolean
			Status               graphql.String
			Entries              []struct {
				ID          graphql.Int
				MediaID     graphql.Int
				Progress    graphql.Int
				Score       graphql.Int `graphql:"score(format: POINT_100)"`
				StartedAt   date
				CompletedAt date
				UpdatedAt   graphql.Int
				Media       struct {
					ID       graphql.Int
					IDMal    graphql.Int
					Title    animeTitle
					Format   graphql.String
					Episodes graphql.Int
					Status   graphql.String
				}
				Repeat graphql.Int
			}
		}
	} `graphql:"MediaListCollection(userId:$userId,forceSingleCompletedList:true,type:ANIME)"`
}

type findAnimeByTitleQuery struct {
	Page struct {
		Media []struct {
			ID       graphql.Int
			IDMal    graphql.Int
			Title    animeTitle
			Format   graphql.String
			Episodes graphql.Int
			Status   graphql.String
		} `graphql:"media(search:$title,type:ANIME)"`
	} `graphql:"Page(page:$pageId)"`
}

// -------------------------- mutations --------------------------
type updateEntryMutation struct {
	SaveMediaListEntry struct {
		ID graphql.Int
	} `graphql:"SaveMediaListEntry(mediaId:$id,progress:$progress,status:$status,scoreRaw:$scoreRaw)"`
}

func (dt date) ToDateTime() null.Time {
	if dt.Year == 0 && dt.Month == 0 && dt.Day == 0 ||
		dt.Year < 0 { // idk what this is, but it happened to me once with yamada-kun
		return null.TimeFromPtr(nil)
	}

	if dt.Year == 0 {
		dt.Year = 1
	}
	if dt.Month == 0 {
		dt.Month = 1
	}
	if dt.Day == 0 {
		dt.Day = 1
	}

	result := time.Date(
		int(dt.Year),
		time.Month(dt.Month),
		int(dt.Day),
		0, 0, 0, 0, time.UTC)

	return null.TimeFrom(result)
}
