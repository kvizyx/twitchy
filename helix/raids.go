package helix

import "context"

type StartRaidRequest struct {
	FromBroadcasterID string `query:"from_broadcaster_id"`
	ToBroadcasterID   string `query:"to_broadcaster_id"`
}

type Raid struct {
	CreatedAt Timestamp `json:"created_at"`
	IsMature  bool      `json:"is_mature"`
}

type StartRaidData []Raid

type CancelRaidRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type CancelRaidData struct{}

func (s *RaidsService) StartRaid(ctx context.Context, req StartRaidRequest) (*Response[StartRaidData], error) {
	return executeTask23Endpoint[StartRaidData](s.client, ctx, "start-a-raid", req, nil, req.FromBroadcasterID)
}

func (s *RaidsService) CancelRaid(ctx context.Context, req CancelRaidRequest) (*Response[CancelRaidData], error) {
	return executeTask23Endpoint[CancelRaidData](s.client, ctx, "cancel-a-raid", req, nil, req.BroadcasterID)
}
