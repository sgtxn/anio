package anilist

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"

	"anio/config"

	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type Auth struct {
	Client *http.Client
}

func Authenticate(ctx context.Context, cfg *config.AnilistAuthConfig) (*Auth, error) {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  "https://anilist.co/api/v2/oauth/pin",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://anilist.co/api/v2/oauth/authorize",
			TokenURL: "https://anilist.co/api/v2/oauth/token",
		},
	}

	url := oauthCfg.AuthCodeURL("auth_code")
	log.Info().Msg("launching the browser to authenticate in anilist...")

	if err := browser.OpenURL(url); err != nil {
		return nil, fmt.Errorf("failed to open the auth url in browser: %w", err)
	}

	log.Info().Msg("Enter the authorization code:")

	var code string
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &code); err != nil {
		return nil, fmt.Errorf("failed to read the authorization code: %w", err)
	}

	token, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange the authorization code for an access token: %w", err)
	}

	client := oauthCfg.Client(ctx, token)

	return &Auth{Client: client}, nil
}
