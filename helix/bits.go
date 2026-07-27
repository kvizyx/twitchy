package helix

import "context"

type GetBitsLeaderboardRequest struct {
	Count     *int    `query:"count,omitempty"`
	Period    *string `query:"period,omitempty"`
	StartedAt *string `query:"started_at,omitempty"`
	UserID    string  `query:"user_id,omitempty"`
}

type BitsLeaderboardEntry struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	Rank      int    `json:"rank"`
	Score     int64  `json:"score"`
}

type BitsDateRange struct {
	StartedAt Timestamp `json:"started_at"`
	EndedAt   Timestamp `json:"ended_at"`
}

type GetBitsLeaderboardData []BitsLeaderboardEntry

type GetCheermotesRequest struct {
	BroadcasterID string `query:"broadcaster_id,omitempty"`
}

type CheermoteImages map[string]map[string]map[string]string

type CheermoteTier struct {
	MinBits        int64           `json:"min_bits"`
	ID             string          `json:"id"`
	Color          string          `json:"color"`
	Images         CheermoteImages `json:"images"`
	CanCheer       bool            `json:"can_cheer"`
	ShowInBitsCard bool            `json:"show_in_bits_card"`
}

type Cheermote struct {
	Prefix       string          `json:"prefix"`
	Tiers        []CheermoteTier `json:"tiers"`
	Type         string          `json:"type"`
	Order        int             `json:"order"`
	LastUpdated  Timestamp       `json:"last_updated"`
	IsCharitable bool            `json:"is_charitable"`
}

type GetCheermotesData []Cheermote

type GetExtensionTransactionsRequest struct {
	ExtensionID string   `query:"extension_id"`
	IDs         []string `query:"id,omitempty"`
	First       *int     `query:"first,omitempty"`
	After       *string  `query:"after,omitempty"`
}

type TransactionCost struct {
	Amount int64  `json:"amount"`
	Type   string `json:"type"`
}

type TransactionProductData struct {
	SKU           string          `json:"sku"`
	Domain        string          `json:"domain"`
	Cost          TransactionCost `json:"cost"`
	InDevelopment bool            `json:"inDevelopment"`
	DisplayName   string          `json:"displayName"`
	Expiration    string          `json:"expiration"`
	Broadcast     bool            `json:"broadcast"`
}

type ExtensionTransaction struct {
	ID               string                 `json:"id"`
	Timestamp        Timestamp              `json:"timestamp"`
	BroadcasterID    string                 `json:"broadcaster_id"`
	BroadcasterLogin string                 `json:"broadcaster_login"`
	BroadcasterName  string                 `json:"broadcaster_name"`
	UserID           string                 `json:"user_id"`
	UserLogin        string                 `json:"user_login"`
	UserName         string                 `json:"user_name"`
	ProductType      string                 `json:"product_type"`
	ProductData      TransactionProductData `json:"product_data"`
}

type GetExtensionTransactionsData []ExtensionTransaction

func (s *BitsService) GetBitsLeaderboard(ctx context.Context, req GetBitsLeaderboardRequest) (*Response[GetBitsLeaderboardData], error) {
	return executeEndpoint[GetBitsLeaderboardData](s.client, ctx, "get-bits-leaderboard", req)
}

func (s *BitsService) GetCheermotes(ctx context.Context, req GetCheermotesRequest) (*Response[GetCheermotesData], error) {
	return executeEndpoint[GetCheermotesData](s.client, ctx, "get-cheermotes", req)
}

func (s *BitsService) GetExtensionTransactions(ctx context.Context, req GetExtensionTransactionsRequest) (*Response[GetExtensionTransactionsData], error) {
	return executeEndpoint[GetExtensionTransactionsData](s.client, ctx, "get-extension-transactions", req)
}

func (s *BitsService) GetExtensionTransactionsPager(req GetExtensionTransactionsRequest, opts ...PagerOption) (*Pager[GetExtensionTransactionsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetExtensionTransactionsData], error) {
		request := any(req)
		if cursor != "" {
			cursorRequest, cursorErr := plan.withCursor(req, cursor)
			if cursorErr != nil {
				return nil, cursorErr
			}
			request = cursorRequest
		}
		return s.GetExtensionTransactions(ctx, request.(GetExtensionTransactionsRequest))
	}, opts...)
}
