package anilist

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"anio/config/outputs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func authenticate(ctx context.Context, cfg *outputs.AnilistAuthConfig) (*http.Client, int, error) {
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
		return nil, 0, fmt.Errorf("failed to open the auth url in browser: %w", err)
	}

	log.Info().Msg("Enter the authorization code:")

	var code string
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &code); err != nil {
		return nil, 0, fmt.Errorf("failed to read the authorization code: %w", err)
	}

	token, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to exchange the authorization code for an access token: %w", err)
	}

	userID, err := extractUserID(token.AccessToken)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to extract user ID from an access token: %w", err)
	}

	client := oauthCfg.Client(ctx, token)

	return client, userID, nil
}

func extractUserID(token string) (int, error) {
	tokenParsed, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return 0, fmt.Errorf("failed to parse the access token token: %w", err)
	}

	userIDString, ok := tokenParsed.Claims.(jwt.MapClaims)["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("user ID not found in anilist token claims")
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user ID as a number: %w", err)
	}

	return int(userID), nil
}
