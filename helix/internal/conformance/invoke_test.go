package conformance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/manifest"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func executeHappyContract(t *testing.T, operation manifest.Operation) {
	t.Helper()
	status := successfulStatus(operation)
	body := responseBody(operation, status)
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: status, Body: body})
	client, err := newClient(operation, transport)
	if err != nil {
		t.Fatal(err)
	}
	method, err := resolveMethod(operation, client)
	if err != nil {
		t.Fatal(err)
	}
	result, callErr := invoke(method, operationRequest(operation, method))
	if callErr != nil {
		t.Fatalf("happy call: %v", callErr)
	}
	meta, err := responseMeta(result)
	if err != nil {
		t.Fatal(err)
	}
	fixture, err := requestFixture(operation, transport, testkit.ContractResponse{Status: status, Body: body, Success: true}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if operation.Response.Format != "text" {
		if err := testkit.RunManifestContract(context.Background(), operation, fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
			return meta, nil
		}); err != nil {
			t.Fatal(err)
		}
	}
}

func executeNegativeContract(t *testing.T, operation manifest.Operation) {
	t.Helper()
	const status = http.StatusUnauthorized
	body := `{"error":"Unauthorized","status":401,"message":"conformance negative fixture"}`
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: status, Body: body})
	client, err := newClient(operation, transport)
	if err != nil {
		t.Fatal(err)
	}
	method, err := resolveMethod(operation, client)
	if err != nil {
		t.Fatal(err)
	}
	_, callErr := invoke(method, operationRequest(operation, method))
	if callErr == nil {
		t.Fatal("negative call unexpectedly succeeded")
	}
	meta, err := errorMeta(callErr)
	if err != nil {
		t.Fatal(err)
	}
	fixture, err := requestFixture(operation, transport, testkit.ContractResponse{Status: status, Body: body, Success: false}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := testkit.RunManifestContract(context.Background(), operation, fixture, transport, func(context.Context) (helix.ResponseMeta, error) {
		return meta, callErr
	}); err != nil {
		t.Fatal(err)
	}
}

func invoke(method reflect.Value, request reflect.Value) (reflect.Value, error) {
	results := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), request})
	if errValue := results[1].Interface(); errValue != nil {
		return results[0], errValue.(error)
	}
	return results[0], nil
}

func errorMeta(err error) (helix.ResponseMeta, error) {
	var metaError interface{ Meta() helix.ResponseMeta }
	if !errors.As(err, &metaError) {
		return helix.ResponseMeta{}, fmt.Errorf("conformance: error %T has no response metadata", err)
	}
	return metaError.Meta(), nil
}

func newClient(operation manifest.Operation, transport *testkit.RecordingRoundTripper) (*helix.Client, error) {
	credential := helix.Credential{
		AccessToken: "conformance-token",
		ClientID:    "conformance-client",
		TokenClass:  conformanceTokenClass(operation),
		UserID:      "conformance-user",
		Scopes:      operationScopes(operation),
	}
	return helix.New(
		helix.WithBaseURL("http://127.0.0.1/helix"),
		helix.WithHTTPClient(&http.Client{Transport: testkit.NewOfflineRoundTripper(transport)}),
		helix.WithStaticToken(credential),
	)
}

func conformanceTokenClass(operation manifest.Operation) helix.TokenClass {
	if operation.Implementation.ServiceType == "EventSubService" {
		return helix.TokenClassApp
	}
	if operation.Anchor == "send-chat-announcement" || operation.Anchor == "send-chat-message" {
		return helix.TokenClassApp
	}
	if isExtensionOnly(operation.Anchor) {
		return helix.TokenClassExtension
	}
	for _, wanted := range []string{"user", "app", "extension"} {
		for _, tokenClass := range operation.TokenClasses {
			if string(tokenClass) == wanted {
				return helix.TokenClass(wanted)
			}
		}
	}
	for _, tokenClass := range operation.TokenClasses {
		switch string(tokenClass) {
		case "app":
			return helix.TokenClassApp
		case "user":
			return helix.TokenClassUser
		case "extension":
			return helix.TokenClassExtension
		}
	}
	return helix.TokenClassUser
}

func isExtensionOnly(anchor string) bool {
	return strings.HasPrefix(anchor, "get-extension-configuration") ||
		strings.HasPrefix(anchor, "set-extension-") ||
		strings.HasPrefix(anchor, "send-extension-") ||
		anchor == "get-extension-secrets" || anchor == "create-extension-secret" ||
		anchor == "get-extensions"
}

func operationScopes(operation manifest.Operation) []helix.AuthorizationScope {
	scopes := make([]helix.AuthorizationScope, 0, len(operation.Scopes))
	for _, scope := range operation.Scopes {
		if scope != "" && scope != "unknown" {
			scopes = append(scopes, helix.AuthorizationScope(scope))
		}
	}
	return scopes
}

func successfulStatus(operation manifest.Operation) int {
	for _, status := range operation.Response.Status {
		if status >= http.StatusOK && status < http.StatusMultipleChoices {
			return status
		}
	}
	return http.StatusNoContent
}

func responseBody(operation manifest.Operation, status int) string {
	if status == http.StatusNoContent {
		return ""
	}
	if operation.Response.Format == "text" {
		return "BEGIN:VCALENDAR\n"
	}
	return `{"data":null}`
}

func populatedValue(typ reflect.Type, operation manifest.Operation) reflect.Value {
	value := reflect.New(typ).Elem()
	populate(value, operation, "")
	return value
}

func populate(value reflect.Value, operation manifest.Operation, fieldName string) {
	switch value.Kind() {
	case reflect.Pointer:
		value.Set(reflect.New(value.Type().Elem()))
		populate(value.Elem(), operation, fieldName)
	case reflect.Struct:
		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			name := value.Type().Field(index).Name
			if field.CanSet() && !skipConformanceField(operation, name) {
				populate(field, operation, name)
			}
		}
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return
		}
		value.Set(reflect.MakeSlice(value.Type(), 1, 1))
		populate(value.Index(0), operation, fieldName)
	case reflect.String:
		value.SetString(conformanceString(value.Type(), fieldName))
	case reflect.Bool:
		value.SetBool(false)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(1)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(1)
	}
}

func skipConformanceField(operation manifest.Operation, fieldName string) bool {
	if operation.Anchor == "get-eventsub-subscriptions" {
		return fieldName == "Type" || fieldName == "UserID" || fieldName == "SubscriptionID" || fieldName == "ConduitID" || fieldName == "First"
	}
	if operation.Anchor == "get-teams" {
		return fieldName == "ID"
	}
	return fieldName == "Before" || fieldName == "VideoID" || fieldName == "TeamName"
}

func conformanceString(typ reflect.Type, fieldName string) string {
	if strings.Contains(typ.String(), "EventSubTransportMethod") || fieldName == "TransportMethod" || fieldName == "Method" {
		return "webhook"
	}
	if strings.HasSuffix(fieldName, "ID") {
		return "conformance-user"
	}
	if strings.HasSuffix(fieldName, "ClientID") {
		return "conformance-client"
	}
	return "conformance-value"
}
