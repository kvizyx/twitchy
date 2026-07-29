package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/kvizyx/twitchy/oauth"
	"github.com/redis/go-redis/v9"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/oauth2/validate":
			writeJSON(writer, map[string]any{
				"client_id":  "local-client",
				"login":      "local-login",
				"scopes":     []string{"user:read:email"},
				"user_id":    "local-user",
				"expires_in": 3600,
			})
		case "/oauth2/token":
			writeJSON(writer, map[string]any{
				"access_token":  "rotated-access-placeholder",
				"refresh_token": "rotated-refresh-placeholder",
				"expires_in":    3600,
				"scope":         []string{"user:read:email"},
				"token_type":    "bearer",
			})
		default:
			writer.WriteHeader(http.StatusNotFound)
			writeJSON(writer, map[string]any{"error": "Not Found", "status": http.StatusNotFound, "message": "offline example route"})
		}
	}))
	defer server.Close()

	client, err := oauth.New(oauth.WithBaseURL(server.URL), oauth.WithHTTPClient(server.Client()))
	if err != nil {
		panic(err)
	}

	hook := func(_ context.Context, pair oauth.TokenPair) error {
		fmt.Printf("refresh hook: type=%s expires=%s scopes=%d\n", pair.TokenType, pair.ExpiresIn, len(pair.Scopes))
		return nil
	}
	source, err := oauth.NewRefreshingTokenSource(client, oauth.TokenPair{
		AccessToken:  "initial-access-placeholder",
		RefreshToken: "initial-refresh-placeholder",
		ExpiresIn:    time.Minute,
		TokenType:    "bearer",
	}, hook)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if _, err := source.Token(ctx); err != nil {
		panic(err)
	}
	session, err := oauth.NewManagedSession(ctx, source)
	if err != nil {
		panic(err)
	}
	if err := session.Close(); err != nil {
		panic(err)
	}
	fmt.Println("managed session validated hourly by default and closed explicitly")

	runCoordinatedExample(client)
}

func runCoordinatedExample(client *oauth.Client) {
	address := os.Getenv("TWITCHY_EXAMPLE_REDIS_ADDR")
	if address == "" {
		fmt.Println("coordinated registry example skipped: set TWITCHY_EXAMPLE_REDIS_ADDR to run it")
		return
	}
	redisClient := redis.NewClient(&redis.Options{Addr: address})
	defer redisClient.Close()

	coordinator, err := oauth.NewRedisRefreshCoordinator(redisClient,
		func(userID string) string { return "twitchy:example:refresh:" + userID })
	if err != nil {
		panic(err)
	}
	registry, err := oauth.NewCoordinatedRegistry(client, coordinator)
	if err != nil {
		panic(err)
	}
	defer registry.Close()

	var stored oauth.TokenPair
	loader := func(context.Context, string) (oauth.TokenPair, error) {
		if stored.AccessToken == "" {
			return oauth.TokenPair{
				AccessToken:  "initial-access-placeholder",
				RefreshToken: "initial-refresh-placeholder",
				ExpiresIn:    time.Minute,
				TokenType:    "bearer",
			}, nil
		}
		return stored, nil
	}
	hook := func(_ context.Context, pair oauth.TokenPair) error {
		stored = pair
		fmt.Printf("coordinated hook: type=%s expires=%s scopes=%d\n", pair.TokenType, pair.ExpiresIn, len(pair.Scopes))
		return nil
	}
	ctx := context.Background()
	if err := registry.AddCoordinatedUser(ctx, "example-user", loader, hook, "chat"); err != nil {
		panic(err)
	}
	source, err := registry.SourceForUser("example-user")
	if err != nil {
		panic(err)
	}
	if _, err := source.Token(ctx); err != nil {
		panic(err)
	}
	fmt.Println("coordinated registry rotated and committed one credential under the Redis lease")
}

func writeJSON(writer http.ResponseWriter, value any) {
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		panic(err)
	}
}
