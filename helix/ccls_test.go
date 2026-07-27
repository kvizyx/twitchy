package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestCCLsGetContentClassificationLabels_encodesLocaleAndLabels(t *testing.T) {
	fixture := testkit.ContractFixture{
		Request: testkit.ContractRequest{
			Query: urlValues("locale", "fr-FR"),
			Headers: http.Header{
				"Authorization": {"Bearer app-token"},
				"Client-Id":     {"client-id"},
			},
		},
		Response: testkit.ContractResponse{
			Status:  http.StatusOK,
			Body:    `{"data":[{"id":"DebatedSocialIssuesAndPolitics","description":"Questions sociales et politiques sensibles","name":"Politique et questions sociales sensibles"},{"id":"DrugsIntoxication","description":"Drogues et intoxication","name":"Drogues, intoxication ou tabac excessif"}]}`,
			Success: true,
		},
		Want: testkit.ContractExpectation{Attempts: 1},
	}
	transport := testkit.NewRecordingRoundTripper(testkit.ResponseFromFixture(fixture.Response))
	client, err := newTask17Client(transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	if err != nil {
		t.Fatal(err)
	}

	result, callErr := client.CCLs.GetContentClassificationLabels(context.Background(), helix.GetContentClassificationLabelsRequest{Locale: stringPtr("fr-FR")})
	if callErr != nil {
		t.Fatal(callErr)
	}
	if err := testkit.RunManifestContract(context.Background(), manifestOperation(t, "get-content-classification-labels"), fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return result.Meta, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 2 || result.Data[0].ID != "DebatedSocialIssuesAndPolitics" || result.Data[0].Name == "" || result.Data[1].Description == "" {
		t.Fatalf("classification labels = %#v", result.Data)
	}
}

func TestCCLsGetContentClassificationLabels_rejectsWrongTokenClassBeforeNetwork(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper()
	client, err := newTask17Client(transport, helix.Credential{AccessToken: "extension-token", ClientID: "client-id", TokenClass: helix.TokenClassExtension})
	if err != nil {
		t.Fatal(err)
	}

	_, callErr := client.CCLs.GetContentClassificationLabels(context.Background(), helix.GetContentClassificationLabelsRequest{})
	var authErr *helix.AuthError
	if !errors.As(callErr, &authErr) {
		t.Fatalf("error = %T %v, want AuthError", callErr, callErr)
	}
	if len(transport.Requests()) != 0 {
		t.Fatal("wrong CCL token class reached the network")
	}
}

func TestCCLsGetContentClassificationLabels_returnsProtocolErrorForMalformedJSON(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{
		StatusCode: http.StatusOK,
		Body:       "not-json",
	})
	client, err := newTask17Client(transport, helix.Credential{AccessToken: "app-token", ClientID: "client-id", TokenClass: helix.TokenClassApp})
	if err != nil {
		t.Fatal(err)
	}

	_, callErr := client.CCLs.GetContentClassificationLabels(context.Background(), helix.GetContentClassificationLabelsRequest{})
	var protocolErr *helix.ProtocolError
	if !errors.As(callErr, &protocolErr) {
		t.Fatalf("error = %T %v, want ProtocolError", callErr, callErr)
	}
}
