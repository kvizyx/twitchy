package helix

import "context"

type GetUserActiveExtensionsRequest struct {
	UserID *string `query:"user_id,omitempty"`
}

type UserActiveExtension struct {
	Active  bool   `json:"active"`
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

type GetUserActiveExtensionsData struct {
	Panel     map[string]UserActiveExtension `json:"panel"`
	Overlay   map[string]UserActiveExtension `json:"overlay"`
	Component map[string]UserActiveExtension `json:"component"`
}

type UserExtensionUpdate struct {
	Active  bool    `json:"active"`
	ID      *string `json:"id,omitempty"`
	Version *string `json:"version,omitempty"`
	X       *int    `json:"x,omitempty"`
	Y       *int    `json:"y,omitempty"`
}

type UpdateUserExtensionsRequest struct {
	Data map[string]map[string]UserExtensionUpdate `query:"-" json:"data"`
}

type UpdateUserExtensionsData struct {
	Panel     map[string]UserActiveExtension `json:"panel"`
	Overlay   map[string]UserActiveExtension `json:"overlay"`
	Component map[string]UserActiveExtension `json:"component"`
}

func (s *UsersService) GetUserActiveExtensions(ctx context.Context, req GetUserActiveExtensionsRequest) (*Response[GetUserActiveExtensionsData], error) {
	subjectID := ""
	if req.UserID != nil {
		subjectID = *req.UserID
	}
	return executeTask26EndpointWithSecrets[GetUserActiveExtensionsData](s.client, ctx, "get-user-active-extensions", req, nil, task26Authorization{subjectID: subjectID, appSubjectRequired: true}, nil)
}

func (s *UsersService) UpdateUserExtensions(ctx context.Context, req UpdateUserExtensionsRequest) (*Response[UpdateUserExtensionsData], error) {
	return executeTask26Endpoint[UpdateUserExtensionsData](s.client, ctx, "update-user-extensions", nil, req, nil, "")
}
