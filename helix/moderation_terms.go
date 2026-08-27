package helix

import "context"

type GetBlockedTermsRequest struct {
	BroadcasterID string  `query:"broadcaster_id"`
	ModeratorID   string  `query:"moderator_id"`
	First         *int    `query:"first,omitempty"`
	After         *string `query:"after,omitempty"`
}

type BlockedTerm struct {
	BroadcasterID string     `json:"broadcaster_id"`
	ModeratorID   string     `json:"moderator_id"`
	ID            string     `json:"id"`
	Text          string     `json:"text"`
	CreatedAt     Timestamp  `json:"created_at"`
	UpdatedAt     Timestamp  `json:"updated_at"`
	ExpiresAt     *Timestamp `json:"expires_at"`
}

type GetBlockedTermsData []BlockedTerm
type AddBlockedTermData []BlockedTerm
type RemoveBlockedTermData struct{}

type AddBlockedTermRequest struct {
	BroadcasterID string `query:"broadcaster_id" json:"-"`
	ModeratorID   string `query:"moderator_id" json:"-"`
	Text          string `query:"-" json:"text"`
}

type RemoveBlockedTermRequest struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	ID            string `query:"id"`
}

func (s *ModerationService) GetBlockedTerms(ctx context.Context, req GetBlockedTermsRequest) (*Response[GetBlockedTermsData], error) {
	scopes := [][]AuthorizationScope{{ScopeModeratorReadBlockedTerms}, {ScopeModeratorManageBlockedTerms}}
	return executeModerationEndpoint[GetBlockedTermsData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "get-blocked-terms", query: req, auth: moderationAuthorization{userScopeSets: scopes, subjectIDs: []string{req.ModeratorID}}})
}

func (s *ModerationService) GetBlockedTermsPager(req GetBlockedTermsRequest, opts ...PagerOption) (*Pager[GetBlockedTermsData], error) {
	initialCursor := ""
	if req.After != nil {
		initialCursor = *req.After
	}
	request := req
	request.After = nil
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetBlockedTermsData], error) {
		pageRequest := request
		if cursor != "" {
			pageRequest.After = &cursor
		}
		return s.GetBlockedTerms(ctx, pageRequest)
	}, initialCursor, opts...)
}

func (s *ModerationService) AddBlockedTerm(ctx context.Context, req AddBlockedTermRequest) (*Response[AddBlockedTermData], error) {
	return executeModerationEndpoint[AddBlockedTermData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "add-blocked-term", query: req, body: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageBlockedTerms), subjectIDs: []string{req.ModeratorID}}})
}

func (s *ModerationService) RemoveBlockedTerm(ctx context.Context, req RemoveBlockedTermRequest) (*Response[RemoveBlockedTermData], error) {
	return executeModerationEndpoint[RemoveBlockedTermData](moderationEndpointSpec{client: s.client, ctx: ctx, anchor: "remove-blocked-term", query: req, auth: moderationAuthorization{userScopeSets: chatReadScopes(ScopeModeratorManageBlockedTerms), subjectIDs: []string{req.ModeratorID}}})
}
