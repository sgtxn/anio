package anilist

import (
	"context"
	"fmt"
	"net/http"

	"anio/config/outputs"

	"github.com/hasura/go-graphql-client"
)

type Client struct {
	httpClient    *http.Client
	graphqlClient *graphql.Client
	userID        int
}

func New(ctx context.Context, cfg *outputs.AnilistAuthConfig) (*Client, error) {
	client, userID, err := authenticate(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("anilist authentication error: %w", err)
	}

	return &Client{
		userID:        userID,
		httpClient:    client,
		graphqlClient: graphql.NewClient("https://graphql.anilist.co", client),
	}, nil
}
