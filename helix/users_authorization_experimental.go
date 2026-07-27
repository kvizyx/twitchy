package helix

import "context"

type GetAuthorizationByUserRequest struct {
	UserIDs []string `query:"user_id"`
}

type UserAuthorization struct {
	UserID        string               `json:"user_id"`
	UserName      string               `json:"user_name"`
	UserLogin     string               `json:"user_login"`
	Scopes        []AuthorizationScope `json:"scopes"`
	HasAuthorized bool                 `json:"has_authorized"`
}

type GetAuthorizationByUserData []UserAuthorization

func (s *ExperimentalUsersService) GetAuthorizationByUser(ctx context.Context, req GetAuthorizationByUserRequest) (*Response[GetAuthorizationByUserData], error) {
	return executeTask26Endpoint[GetAuthorizationByUserData](s.client, ctx, "get-authorization-by-user", req, nil, nil, "")
}
