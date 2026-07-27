package helix_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestExtensionsConfigurationUsesExtensionAuthorization(t *testing.T) {
	responses := []testkit.RoundTripResponse{
		extensionResponse(http.StatusOK, extensionBody(t, "extension_configuration.json")),
		extensionResponse(http.StatusNoContent, ""),
		extensionResponse(http.StatusNoContent, ""),
		extensionResponse(http.StatusNoContent, ""),
		extensionResponse(http.StatusNoContent, ""),
	}
	transport := testkit.NewRecordingRoundTripper(responses...)
	client := extensionClient(t, transport, helix.Credential{AccessToken: "jwt-token", ClientID: "client-id", TokenClass: helix.TokenClassExtension})

	configuration, err := client.Extensions.GetExtensionConfigurationSegment(context.Background(), helix.GetExtensionConfigurationSegmentRequest{BroadcasterID: "broadcaster-1", ExtensionID: "extension-1", Segment: []string{"broadcaster", "global"}})
	if err != nil {
		t.Fatal(err)
	}
	extensionMetaContract(t, "get-extension-configuration-segment", extensionFixture(
		map[string][]string{"broadcaster_id": {"broadcaster-1"}, "extension_id": {"extension-1"}, "segment": {"broadcaster", "global"}}, "", extensionHeaders("Extension", "jwt-token"), testkit.ContractResponse{Status: http.StatusOK, Headers: http.Header{"Ratelimit-Limit": {"8000"}, "Ratelimit-Remaining": {"7999"}, "Ratelimit-Reset": {extensionRateReset}}, Body: extensionBody(t, "extension_configuration.json"), Success: true}), transport, configuration.Meta, nil)

	if _, err := client.Extensions.SetExtensionConfigurationSegment(context.Background(), helix.SetExtensionConfigurationSegmentRequest{ExtensionID: "extension-1", Segment: "global", Content: stringPointer("hello config"), Version: stringPointer("1")}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Extensions.SetExtensionRequiredConfiguration(context.Background(), helix.SetExtensionRequiredConfigurationRequest{BroadcasterID: "broadcaster-1", ExtensionID: "extension-1", ExtensionVersion: "1", RequiredConfiguration: "RCS"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Extensions.SendExtensionPubSubMessage(context.Background(), helix.SendExtensionPubSubMessageRequest{Target: []string{"broadcast"}, BroadcasterID: "broadcaster-1", Message: "hello"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Extensions.SendExtensionChatMessage(context.Background(), helix.SendExtensionChatMessageRequest{BroadcasterID: "broadcaster-1", Text: "Hello", ExtensionID: "extension-1", ExtensionVersion: "1"}); err != nil {
		t.Fatal(err)
	}

	requests := transport.Requests()
	if got := requests[0].Header.Get("Authorization"); got != "Extension jwt-token" {
		t.Fatalf("extension authorization = %q", got)
	}
	wantPaths := []string{
		"/helix/extensions/configurations?broadcaster_id=broadcaster-1&extension_id=extension-1&segment=broadcaster&segment=global",
		"/helix/extensions/configurations",
		"/helix/extensions/required_configuration?broadcaster_id=broadcaster-1",
		"/helix/extensions/pubsub",
		"/helix/extensions/chat?broadcaster_id=broadcaster-1",
	}
	for index, want := range wantPaths {
		if requests[index].Path != want || requests[index].Header.Get("Authorization") != "Extension jwt-token" {
			t.Fatalf("request[%d] = %#v, want path %q and extension auth", index, requests[index], want)
		}
	}
	if got := string(requests[1].Body); got != `{"extension_id":"extension-1","segment":"global","content":"hello config","version":"1"}` {
		t.Fatalf("configuration body = %s", got)
	}
	if got := string(requests[2].Body); got != `{"extension_id":"extension-1","extension_version":"1","required_configuration":"RCS"}` {
		t.Fatalf("required configuration body = %s", got)
	}
	if got := string(requests[3].Body); got != `{"target":["broadcast"],"broadcaster_id":"broadcaster-1","message":"hello"}` {
		t.Fatalf("pubsub body = %s", got)
	}
}

func TestExtensionsSecretsAndListingsUseExactSchemes(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		extensionResponse(http.StatusOK, extensionBody(t, "extension_secrets.json")),
		extensionResponse(http.StatusOK, extensionBody(t, "extension_secrets.json")),
		extensionResponse(http.StatusOK, extensionBody(t, "extensions.json")),
		extensionResponse(http.StatusOK, extensionBody(t, "extensions.json")),
		extensionResponse(http.StatusOK, extensionBody(t, "extension_bits_products.json")),
		extensionResponse(http.StatusOK, extensionBody(t, "extension_bits_products.json")),
	)
	client := extensionClient(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})

	if _, err := client.Extensions.GetExtensionSecrets(context.Background(), helix.GetExtensionSecretsRequest{ExtensionID: "extension-1"}); err == nil {
		t.Fatal("app token was accepted for extension-secret auth")
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("unsupported extension auth reached network")
	}

	client = extensionClient(t, transport, helix.Credential{AccessToken: "jwt-token", ClientID: "client-id", TokenClass: helix.TokenClassExtension})
	secrets, err := client.Extensions.GetExtensionSecrets(context.Background(), helix.GetExtensionSecretsRequest{ExtensionID: "extension-1"})
	if err != nil || len(secrets.Data) != 1 || secrets.Data[0].Secrets[0].Content != "secret-value" {
		t.Fatalf("secrets = %#v, err = %v", secrets, err)
	}
	if _, err := client.Extensions.CreateExtensionSecret(context.Background(), helix.CreateExtensionSecretRequest{ExtensionID: "extension-1", Delay: intPointer(600)}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Extensions.GetExtensions(context.Background(), helix.GetExtensionsRequest{ExtensionID: "extension-1", ExtensionVersion: stringPointer("1")}); err != nil {
		t.Fatal(err)
	}
	appClient := extensionClient(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	if _, err := appClient.Extensions.GetReleasedExtensions(context.Background(), helix.GetReleasedExtensionsRequest{ExtensionID: "extension-1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := appClient.Extensions.GetExtensionBitsProducts(context.Background(), helix.GetExtensionBitsProductsRequest{ShouldIncludeAll: boolPointer(true)}); err != nil {
		t.Fatal(err)
	}
	if _, err := appClient.Extensions.UpdateExtensionBitsProduct(context.Background(), helix.UpdateExtensionBitsProductRequest{SKU: "sku-1", Cost: helix.ExtensionBitsProductCost{Amount: 100, Type: "bits"}, DisplayName: "Product"}); err != nil {
		t.Fatal(err)
	}

	requests := transport.Requests()
	for index, request := range requests[:3] {
		if request.Header.Get("Authorization") != "Extension jwt-token" {
			t.Fatalf("extension request[%d] authorization = %q", index, request.Header.Get("Authorization"))
		}
	}
	for index, request := range requests[3:] {
		if request.Header.Get("Authorization") != "Bearer app-token" {
			t.Fatalf("bearer request[%d] authorization = %q", index, request.Header.Get("Authorization"))
		}
	}
}

func TestExtensionsLiveChannelsStringPagination(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(
		extensionResponse(http.StatusOK, extensionBody(t, "extension_live_channels_page_1.json")),
		extensionResponse(http.StatusOK, extensionBody(t, "extension_live_channels_page_2.json")),
	)
	client := extensionClient(t, transport, helix.Credential{AccessToken: "user-token", ClientID: "client-id", TokenClass: helix.TokenClassUser})
	pager, err := client.Extensions.GetExtensionLiveChannelsPager(helix.GetExtensionLiveChannelsRequest{ExtensionID: "extension-1", First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || len(pager.Page().Data) != 1 || pager.Page().Data[0].BroadcasterID != "broadcaster-1" {
		t.Fatalf("first live page = %#v, err = %v", pager.Page(), pager.Err())
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].BroadcasterID != "broadcaster-2" || pager.Err() != nil {
		t.Fatalf("second live page = %#v, err = %v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/extensions/live?after=next-cursor&extension_id=extension-1&first=1" {
		t.Fatalf("live pagination requests = %#v", requests)
	}
}

func TestExtensionsRejectsUnsupportedAuthBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client := extensionClient(t, transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	_, err := client.Extensions.GetExtensionConfigurationSegment(context.Background(), helix.GetExtensionConfigurationSegmentRequest{ExtensionID: "extension-1", Segment: []string{"global"}})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 0 {
		t.Fatalf("error = %T %v, requests = %d", err, err, len(transport.Requests()))
	}
}

func TestExtensionsRedirectNeverSendsJWTToHostileServer(t *testing.T) {
	var hostileRequests atomic.Int32
	var hostileAuthorization atomic.Value
	hostile := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		hostileRequests.Add(1)
		hostileAuthorization.Store(request.Header.Get("Authorization"))
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer hostile.Close()

	redirector := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, hostile.URL+"/helix/extensions/jwt/secrets?extension_id=extension-1", http.StatusFound)
	}))
	defer redirector.Close()

	client, err := helix.New(
		helix.WithBaseURL(redirector.URL+"/helix"),
		helix.WithStaticToken(helix.Credential{AccessToken: "jwt-secret", ClientID: "client-id", TokenClass: helix.TokenClassExtension}),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Extensions.GetExtensionSecrets(context.Background(), helix.GetExtensionSecretsRequest{ExtensionID: "extension-1"})
	if err == nil || strings.Contains(err.Error(), "jwt-secret") {
		t.Fatalf("redirect error = %v, want rejected redirect without JWT", err)
	}
	if hostileRequests.Load() != 0 {
		t.Fatal("hostile redirect received a request")
	}
	if value := hostileAuthorization.Load(); value != nil && value.(string) != "" {
		t.Fatalf("hostile authorization = %q", value)
	}
}

func TestExtensionsSecretTokenIsRedactedFromErrors(t *testing.T) {
	secret := "extension-secret"
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusBadRequest,
		Body:       `{"error":"Bad Request","status":400,"message":"` + secret + ` rejected"}`,
	})
	client := extensionClient(t, transport, helix.Credential{AccessToken: secret, ClientID: "client-id", TokenClass: helix.TokenClassExtension})
	_, err := client.Extensions.GetExtensionSecrets(context.Background(), helix.GetExtensionSecretsRequest{ExtensionID: "extension-1"})
	if err == nil || strings.Contains(err.Error(), secret) || !strings.Contains(err.Error(), "[redacted]") {
		t.Fatalf("secret error = %v, want redacted token", err)
	}
}
