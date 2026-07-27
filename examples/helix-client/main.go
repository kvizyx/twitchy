package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

func main() {
	streamFailures := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Ratelimit-Limit", "800")
		writer.Header().Set("Ratelimit-Remaining", "799")
		writer.Header().Set("Ratelimit-Reset", fmt.Sprint(time.Now().Add(time.Hour).Unix()))

		switch request.URL.Path {
		case "/helix/users":
			writeJSON(writer, map[string]any{"data": []map[string]any{{"id": "local-user", "login": "local-user"}}})
		case "/helix/streams":
			if streamFailures == 0 {
				streamFailures++
				writer.WriteHeader(http.StatusBadRequest)
				writeJSON(writer, map[string]any{"error": "Bad Request", "status": http.StatusBadRequest, "message": "offline example failure"})
				return
			}
			writeJSON(writer, map[string]any{"data": []map[string]any{{"id": "local-stream", "user_id": "local-user", "type": "live"}}})
		default:
			writer.WriteHeader(http.StatusNotFound)
			writeJSON(writer, map[string]any{"error": "Not Found", "status": http.StatusNotFound, "message": "offline example route"})
		}
	}))
	defer server.Close()

	client, err := helix.New(
		helix.WithBaseURL(server.URL+"/helix"),
		helix.WithHTTPClient(server.Client()),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "local-placeholder",
			ClientID:    "local-client",
			TokenClass:  helix.TokenClassApp,
		}),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	users, err := client.Users.GetUsers(ctx, helix.GetUsersRequest{Logins: []string{"local-user"}})
	if err != nil {
		panic(err)
	}
	rate := users.Meta.RateLimit()
	if !rate.Valid() {
		panic("offline example did not receive rate-limit metadata")
	}
	fmt.Printf("users=%d rate-limit=%d/%d\n", len(users.Data), rate.Remaining(), rate.Limit())

	_, err = client.Streams.GetStreams(ctx, helix.GetStreamsRequest{UserLogin: []string{"local-user"}})
	if err == nil {
		panic("offline example expected a typed API error")
	}
	var apiErr *helix.APIError
	if !errors.As(err, &apiErr) {
		panic(fmt.Errorf("offline example error type = %T", err))
	}
	fmt.Printf("typed error=%s status=%d\n", apiErr.Operation(), apiErr.StatusCode())

	streams, err := client.Streams.GetStreams(ctx, helix.GetStreamsRequest{UserLogin: []string{"local-user"}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("streams=%d attempts=%d\n", len(streams.Data), streams.Meta.Attempts())
}

func writeJSON(writer http.ResponseWriter, value any) {
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		panic(err)
	}
}
