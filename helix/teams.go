package helix

import (
	"context"
)

type GetChannelTeamsRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type TeamUser struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
}

type Team struct {
	Users              []TeamUser `json:"users"`
	BroadcasterID      string     `json:"broadcaster_id"`
	BroadcasterLogin   string     `json:"broadcaster_login"`
	BroadcasterName    string     `json:"broadcaster_name"`
	BackgroundImageURL *string    `json:"background_image_url"`
	Banner             *string    `json:"banner"`
	CreatedAt          string     `json:"created_at"`
	UpdatedAt          string     `json:"updated_at"`
	Info               string     `json:"info"`
	ThumbnailURL       string     `json:"thumbnail_url"`
	TeamName           string     `json:"team_name"`
	TeamDisplayName    string     `json:"team_display_name"`
	ID                 string     `json:"id"`
}

type GetChannelTeamsData []Team

type GetTeamsRequest struct {
	Name string `query:"name,omitempty"`
	ID   string `query:"id,omitempty"`
}

type GetTeamsData []Team

func (s *TeamsService) GetChannelTeams(ctx context.Context, req GetChannelTeamsRequest) (*Response[GetChannelTeamsData], error) {
	return executeEndpoint[GetChannelTeamsData](s.client, ctx, "get-channel-teams", req)
}

func (s *TeamsService) GetTeams(ctx context.Context, req GetTeamsRequest) (*Response[GetTeamsData], error) {
	nameSet := req.Name != ""
	idSet := req.ID != ""
	if nameSet == idSet {
		return nil, &RequestEncodingError{Reason: "teams request requires exactly one of name or id"}
	}
	return executeEndpoint[GetTeamsData](s.client, ctx, "get-teams", req)
}
