package anilist

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"

	"anio/providers/anilist/config"

	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type Auth struct {
	Client *http.Client
}

func Authenticate(ctx context.Context, cfg *config.AuthConfig) (*Auth, error) {
	conf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  "https://anilist.co/api/v2/oauth/pin",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://anilist.co/api/v2/oauth/authorize",
			TokenURL: "https://anilist.co/api/v2/oauth/token",
		},
	}

	url := conf.AuthCodeURL("auth_code")
	log.Info().Msg("Launching the browser to authenticate in anilist...")

	if err := browser.OpenURL(url); err != nil {
		return nil, fmt.Errorf("failed to open the auth url in browser: %w", err)
	}

	fmt.Println("Enter the authorization code:")

	var code string
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &code); err != nil {
		return nil, fmt.Errorf("failed to read the authorization code: %w", err)
	}

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange the authorization code for an access token: %w", err)
	}

	client := conf.Client(ctx, token)

	return &Auth{Client: client}, nil
}
