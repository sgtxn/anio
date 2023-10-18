//nolint:revive // can't do without nested structs here
package anilist

import "github.com/hasura/go-graphql-client"

type MediaListStatus struct{}

type updateEntryQuery struct {
	SaveMediaListEntry struct {
		ID graphql.Int
	} `graphql:"SaveMediaListEntry(mediaId:$id,progress:$progress,status:$status,scoreRaw:$scoreRaw)"`
}
