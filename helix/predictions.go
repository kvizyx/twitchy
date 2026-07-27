package helix

import "context"

type PredictionStatus string

const (
	PredictionStatusActive   PredictionStatus = "ACTIVE"
	PredictionStatusCanceled PredictionStatus = "CANCELED"
	PredictionStatusLocked   PredictionStatus = "LOCKED"
	PredictionStatusResolved PredictionStatus = "RESOLVED"
)

type PredictionOutcomeColor string

const (
	PredictionOutcomeColorBlue PredictionOutcomeColor = "BLUE"
	PredictionOutcomeColorPink PredictionOutcomeColor = "PINK"
)

type PredictionTopPredictor struct {
	UserID            string `json:"user_id"`
	UserName          string `json:"user_name"`
	UserLogin         string `json:"user_login"`
	ChannelPointsUsed int64  `json:"channel_points_used"`
	ChannelPointsWon  int64  `json:"channel_points_won"`
}

type PredictionOutcome struct {
	ID            string                    `json:"id,omitempty"`
	Title         string                    `json:"title"`
	Users         int64                     `json:"users,omitempty"`
	ChannelPoints int64                     `json:"channel_points,omitempty"`
	TopPredictors *[]PredictionTopPredictor `json:"top_predictors,omitempty"`
	Color         PredictionOutcomeColor    `json:"color,omitempty"`
}

type Prediction struct {
	ID               string              `json:"id"`
	BroadcasterID    string              `json:"broadcaster_id"`
	BroadcasterName  string              `json:"broadcaster_name"`
	BroadcasterLogin string              `json:"broadcaster_login"`
	Title            string              `json:"title"`
	WinningOutcomeID *string             `json:"winning_outcome_id"`
	Outcomes         []PredictionOutcome `json:"outcomes"`
	PredictionWindow int                 `json:"prediction_window"`
	Status           PredictionStatus    `json:"status"`
	CreatedAt        Timestamp           `json:"created_at"`
	EndedAt          *Timestamp          `json:"ended_at"`
	LockedAt         *Timestamp          `json:"locked_at"`
}

type GetPredictionsRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	IDs           []string `query:"id,omitempty"`
	First         *int     `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
}

type GetPredictionsData []Prediction

type CreatePredictionRequest struct {
	BroadcasterID    string              `query:"-" json:"broadcaster_id"`
	Title            string              `query:"-" json:"title"`
	Outcomes         []PredictionOutcome `query:"-" json:"outcomes"`
	PredictionWindow int                 `query:"-" json:"prediction_window"`
}

type CreatePredictionData []Prediction

type EndPredictionRequest struct {
	BroadcasterID    string           `query:"-" json:"broadcaster_id"`
	ID               string           `query:"-" json:"id"`
	Status           PredictionStatus `query:"-" json:"status"`
	WinningOutcomeID *string          `query:"-" json:"winning_outcome_id,omitempty"`
}

type EndPredictionData []Prediction

func (s *PredictionsService) GetPredictions(ctx context.Context, req GetPredictionsRequest) (*Response[GetPredictionsData], error) {
	return executeTask23Endpoint[GetPredictionsData](s.client, ctx, "get-predictions", req, nil, req.BroadcasterID, ScopeChannelReadPredictions, ScopeChannelManagePredictions)
}

func (s *PredictionsService) GetPredictionsPager(req GetPredictionsRequest, opts ...PagerOption) (*Pager[GetPredictionsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetPredictionsData], error) {
		next, err := plan.withCursor(req, cursor)
		if err != nil {
			return nil, err
		}
		return s.GetPredictions(ctx, next.(GetPredictionsRequest))
	}, opts...)
}

func (s *PredictionsService) CreatePrediction(ctx context.Context, req CreatePredictionRequest) (*Response[CreatePredictionData], error) {
	return executeTask23Endpoint[CreatePredictionData](s.client, ctx, "create-prediction", req, req, req.BroadcasterID)
}

func (s *PredictionsService) EndPrediction(ctx context.Context, req EndPredictionRequest) (*Response[EndPredictionData], error) {
	return executeTask23Endpoint[EndPredictionData](s.client, ctx, "end-prediction", req, req, req.BroadcasterID)
}
