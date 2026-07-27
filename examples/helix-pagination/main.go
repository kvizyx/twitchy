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
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		writer.Header().Set("Content-Type", "application/json")
		if request.URL.Query().Get("after") == "" {
			writeJSON(writer, map[string]any{
				"data":       []map[string]any{{"id": "local-stream-1", "type": "live"}},
				"pagination": map[string]string{"cursor": "local-next"},
			})
			return
		}
		writeJSON(writer, map[string]any{"data": []map[string]any{{"id": "local-stream-2", "type": "live"}}})
	}))
	defer server.Close()

	client, err := helix.New(
		helix.WithBaseURL(server.URL+"/helix"),
		helix.WithHTTPClient(server.Client()),
		helix.WithStaticToken(helix.Credential{AccessToken: "local-placeholder", TokenClass: helix.TokenClassApp}),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	single, err := client.Streams.GetStreams(ctx, helix.GetStreamsRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("single page=%d requests=%d\n", len(single.Data), requests)

	pager, err := client.Streams.GetStreamsPager(helix.GetStreamsRequest{}, helix.WithPageLimit(2))
	if err != nil {
		panic(err)
	}
	if requests != 1 {
		panic("pager construction performed I/O")
	}

	pages := 0
	items := 0
	for pager.Next(ctx) {
		pages++
		items += len(pager.Page().Data)
	}
	if err := pager.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("lazy pages=%d items=%d requests=%d\n", pages, items, requests)
}

func writeJSON(writer http.ResponseWriter, value any) {
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		panic(err)
	}
}
