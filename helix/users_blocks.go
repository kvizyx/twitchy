package helix

import "context"

type GetUserBlockListRequest struct {
	BroadcasterID string  `query:"broadcaster_id"`
	First         *int    `query:"first,omitempty"`
	After         *string `query:"after,omitempty"`
}

type BlockedUser struct {
	UserID      string `json:"user_id"`
	UserLogin   string `json:"user_login"`
	DisplayName string `json:"display_name"`
}

type GetUserBlockListData []BlockedUser

type BlockUserSourceContext string

const (
	BlockUserSourceContextChat    BlockUserSourceContext = "chat"
	BlockUserSourceContextWhisper BlockUserSourceContext = "whisper"
)

type BlockUserReason string

const (
	BlockUserReasonHarassment BlockUserReason = "harassment"
	BlockUserReasonSpam       BlockUserReason = "spam"
	BlockUserReasonOther      BlockUserReason = "other"
)

type BlockUserRequest struct {
	TargetUserID  string                  `query:"target_user_id" json:"-"`
	SourceContext *BlockUserSourceContext `query:"source_context,omitempty" json:"-"`
	Reason        *BlockUserReason        `query:"reason,omitempty" json:"-"`
}

type BlockUserData struct{}

type UnblockUserRequest struct {
	TargetUserID string `query:"target_user_id"`
}

type UnblockUserData struct{}

func (s *UsersService) GetUserBlockList(ctx context.Context, req GetUserBlockListRequest) (*Response[GetUserBlockListData], error) {
	return executeTask26Endpoint[GetUserBlockListData](s.client, ctx, "get-user-block-list", req, nil, nil, req.BroadcasterID)
}

func (s *UsersService) GetUserBlockListPager(req GetUserBlockListRequest, opts ...PagerOption) (*Pager[GetUserBlockListData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	initialCursor := ""
	baseRequest := req
	if req.After != nil {
		initialCursor = *req.After
		baseRequest.After = nil
	}
	return newPagerAt(func(ctx context.Context, cursor string) (*Response[GetUserBlockListData], error) {
		if cursor == "" {
			return s.GetUserBlockList(ctx, baseRequest)
		}
		next, err := plan.withCursor(baseRequest, cursor)
		if err != nil {
			return nil, err
		}
		request, ok := next.(GetUserBlockListRequest)
		if !ok {
			return nil, &paginationRequestError{reason: "block-list request has unexpected type"}
		}
		return s.GetUserBlockList(ctx, request)
	}, initialCursor, opts...)
}

func (s *UsersService) BlockUser(ctx context.Context, req BlockUserRequest) (*Response[BlockUserData], error) {
	return executeTask26Endpoint[BlockUserData](s.client, ctx, "block-user", req, nil, nil, "")
}

func (s *UsersService) UnblockUser(ctx context.Context, req UnblockUserRequest) (*Response[UnblockUserData], error) {
	return executeTask26Endpoint[UnblockUserData](s.client, ctx, "unblock-user", req, nil, nil, "")
}
