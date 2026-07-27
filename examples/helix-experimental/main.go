package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/kvizyx/twitchy/helix"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if request.URL.Path != "/helix/bits/custom_power_ups" {
			writer.WriteHeader(http.StatusNotFound)
			writeJSON(writer, map[string]any{"error": "Not Found", "status": http.StatusNotFound, "message": "offline example route"})
			return
		}
		writeJSON(writer, map[string]any{"data": []map[string]any{}})
	}))
	defer server.Close()

	client, err := helix.New(
		helix.WithBaseURL(server.URL+"/helix"),
		helix.WithHTTPClient(server.Client()),
		helix.WithStaticToken(helix.Credential{
			AccessToken: "local-placeholder",
			TokenClass:  helix.TokenClassUser,
			UserID:      "local-broadcaster",
			Scopes:      []helix.AuthorizationScope{"bits:read"},
		}),
	)
	if err != nil {
		panic(err)
	}

	response, err := client.Experimental.Bits.GetCustomPowerUp(context.Background(), helix.GetCustomPowerUpRequest{BroadcasterID: "local-broadcaster"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("experimental custom power-ups=%d\n", len(response.Data))
	fmt.Println("NEW/BETA access is opt-in and has no stable compatibility promise")
}

func writeJSON(writer http.ResponseWriter, value any) {
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		panic(err)
	}
}
