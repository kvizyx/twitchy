package helix

import "context"

type GetHypeTrainStatusRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type HypeTrainContribution struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Type      string `json:"type"`
	Total     int    `json:"total"`
}

type HypeTrainParticipant struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

type HypeTrainCurrent struct {
	ID                      string                  `json:"id"`
	BroadcasterUserID       string                  `json:"broadcaster_user_id"`
	BroadcasterUserLogin    string                  `json:"broadcaster_user_login"`
	BroadcasterUserName     string                  `json:"broadcaster_user_name"`
	Level                   int                     `json:"level"`
	Total                   int                     `json:"total"`
	Progress                int                     `json:"progress"`
	Goal                    int                     `json:"goal"`
	TopContributions        []HypeTrainContribution `json:"top_contributions"`
	SharedTrainParticipants []HypeTrainParticipant  `json:"shared_train_participants"`
	StartedAt               Timestamp               `json:"started_at"`
	ExpiresAt               Timestamp               `json:"expires_at"`
	Type                    string                  `json:"type"`
	IsSharedTrain           bool                    `json:"is_shared_train"`
}

type HypeTrainRecord struct {
	Level      int       `json:"level"`
	Total      int       `json:"total"`
	AchievedAt Timestamp `json:"achieved_at"`
}

type HypeTrainStatus struct {
	Current           *HypeTrainCurrent `json:"current"`
	AllTimeHigh       *HypeTrainRecord  `json:"all_time_high"`
	SharedAllTimeHigh *HypeTrainRecord  `json:"shared_all_time_high"`
}

type GetHypeTrainStatusData []HypeTrainStatus

func (s *HypeTrainService) GetHypeTrainStatus(ctx context.Context, req GetHypeTrainStatusRequest) (*Response[GetHypeTrainStatusData], error) {
	return executeBroadcasterEndpoint[GetHypeTrainStatusData](s.client, ctx, "get-hype-train-status", req, req.BroadcasterID)
}
