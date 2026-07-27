package helix_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/kvizyx/twitchy/helix"
	"github.com/kvizyx/twitchy/helix/internal/testkit"
)

func TestPredictionsCreatePrediction_preservesExactWireAndOutcomes(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusOK, Header: interactiveRateHeaders(), Body: task23Fixture(t, "prediction.json")})
	client := interactiveClient(t, transport)
	result, err := client.Predictions.CreatePrediction(context.Background(), helix.CreatePredictionRequest{
		BroadcasterID:    "123456",
		Title:            "Will it rain?",
		Outcomes:         []helix.PredictionOutcome{{Title: "Yes"}, {Title: "No"}},
		PredictionWindow: 120,
	})
	if err != nil {
		t.Fatal(err)
	}
	request := transport.Requests()[0]
	if string(request.Body) != `{"broadcaster_id":"123456","title":"Will it rain?","outcomes":[{"title":"Yes"},{"title":"No"}],"prediction_window":120}` {
		t.Fatalf("prediction request = %#v", request)
	}
	if result.Data[0].Outcomes[0].Title != "Yes" || result.Data[0].Outcomes[0].Color != helix.PredictionOutcomeColorBlue || result.Data[0].Status != helix.PredictionStatusActive {
		t.Fatalf("prediction response = %#v", result.Data[0])
	}
}

func TestPredictionsGetPredictionsPager_sendsCursorAndDecodesNullableFields(t *testing.T) {
	first := `{"data":[{"id":"prediction-1","broadcaster_id":"123456","title":"Question","winning_outcome_id":null,"outcomes":[{"id":"outcome-1","title":"Yes","users":3,"channel_points":300,"top_predictors":null,"color":"BLUE"}],"prediction_window":120,"status":"ACTIVE","created_at":"2024-01-02T03:04:05Z","ended_at":null,"locked_at":null}],"pagination":{"cursor":"next"}}`
	second := `{"data":[{"id":"prediction-2","status":"RESOLVED","winning_outcome_id":"outcome-2","outcomes":[],"created_at":"2024-01-02T03:04:05Z","ended_at":"2024-01-02T03:05:05Z","locked_at":"2024-01-02T03:04:15Z"}]}`
	transport := testkit.NewRecordingRoundTripper(interactiveResponse(first), interactiveResponse(second))
	client := interactiveClient(t, transport)
	pager, err := client.Predictions.GetPredictionsPager(helix.GetPredictionsRequest{BroadcasterID: "123456", First: intPointer(1)})
	if err != nil {
		t.Fatal(err)
	}
	if !pager.Next(context.Background()) || pager.Page().Data[0].WinningOutcomeID != nil || !pager.Next(context.Background()) || *pager.Page().Data[0].WinningOutcomeID != "outcome-2" || pager.Next(context.Background()) || pager.Err() != nil {
		t.Fatalf("prediction pager state: page=%#v err=%v", pager.Page(), pager.Err())
	}
	requests := transport.Requests()
	if len(requests) != 2 || requests[1].Path != "/helix/predictions?after=next&broadcaster_id=123456&first=1" {
		t.Fatalf("prediction pager requests = %#v", requests)
	}
}

func TestPredictionsEndPrediction_isOneAttemptOnUnauthorized(t *testing.T) {
	transport := testkit.NewRecordingRoundTripper(testkit.RoundTripResponse{StatusCode: http.StatusUnauthorized, Header: interactiveRateHeaders(), Body: `{"error":"Unauthorized","status":401,"message":"invalid"}`})
	client := interactiveClient(t, transport)
	_, err := client.Predictions.EndPrediction(context.Background(), helix.EndPredictionRequest{BroadcasterID: "123456", ID: "prediction-1", Status: helix.PredictionStatusResolved, WinningOutcomeID: stringPointer("outcome-1")})
	var authErr *helix.AuthError
	if !errors.As(err, &authErr) || len(transport.Requests()) != 1 {
		t.Fatalf("prediction auth error=%T %v requests=%d", err, err, len(transport.Requests()))
	}
}
