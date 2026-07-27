package helix

import "context"

type EntitlementFulfillmentStatus string

const (
	EntitlementFulfillmentStatusClaimed   EntitlementFulfillmentStatus = "CLAIMED"
	EntitlementFulfillmentStatusFulfilled EntitlementFulfillmentStatus = "FULFILLED"
)

type EntitlementUpdateStatus string

const (
	EntitlementUpdateStatusInvalidID    EntitlementUpdateStatus = "INVALID_ID"
	EntitlementUpdateStatusNotFound     EntitlementUpdateStatus = "NOT_FOUND"
	EntitlementUpdateStatusSuccess      EntitlementUpdateStatus = "SUCCESS"
	EntitlementUpdateStatusUnauthorized EntitlementUpdateStatus = "UNAUTHORIZED"
	EntitlementUpdateStatusFailed       EntitlementUpdateStatus = "UPDATE_FAILED"
)

type GetDropsEntitlementsRequest struct {
	IDs               []string                     `query:"id,omitempty"`
	UserID            string                       `query:"user_id,omitempty"`
	GameID            string                       `query:"game_id,omitempty"`
	FulfillmentStatus EntitlementFulfillmentStatus `query:"fulfillment_status,omitempty"`
	After             *string                      `query:"after,omitempty"`
	First             *int                         `query:"first,omitempty"`
}

type DropsEntitlement struct {
	ID                string                       `json:"id"`
	BenefitID         string                       `json:"benefit_id"`
	Timestamp         Timestamp                    `json:"timestamp"`
	UserID            string                       `json:"user_id"`
	GameID            string                       `json:"game_id"`
	FulfillmentStatus EntitlementFulfillmentStatus `json:"fulfillment_status"`
	LastUpdated       Timestamp                    `json:"last_updated"`
}

type Entitlement = DropsEntitlement

type GetDropsEntitlementsData []DropsEntitlement

type UpdateDropsEntitlementsRequest struct {
	EntitlementIDs    []string                     `json:"entitlement_ids"`
	FulfillmentStatus EntitlementFulfillmentStatus `json:"fulfillment_status"`
}

type UpdatedDropsEntitlement struct {
	Status EntitlementUpdateStatus `json:"status"`
	IDs    []string                `json:"ids"`
}

type UpdateDropsEntitlementsData []UpdatedDropsEntitlement

func (s *EntitlementsService) GetDropsEntitlements(ctx context.Context, req GetDropsEntitlementsRequest) (*Response[GetDropsEntitlementsData], error) {
	return executeEndpoint[GetDropsEntitlementsData](s.client, ctx, "get-drops-entitlements", req)
}

func (s *EntitlementsService) GetDropsEntitlementsPager(req GetDropsEntitlementsRequest, opts ...PagerOption) (*Pager[GetDropsEntitlementsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetDropsEntitlementsData], error) {
		request := any(req)
		if cursor != "" {
			cursorRequest, cursorErr := plan.withCursor(req, cursor)
			if cursorErr != nil {
				return nil, cursorErr
			}
			request = cursorRequest
		}
		return s.GetDropsEntitlements(ctx, request.(GetDropsEntitlementsRequest))
	}, opts...)
}

func (s *EntitlementsService) UpdateDropsEntitlements(ctx context.Context, req UpdateDropsEntitlementsRequest) (*Response[UpdateDropsEntitlementsData], error) {
	return executeEndpointWithBody[UpdateDropsEntitlementsData](s.client, ctx, "update-drops-entitlements", nil, req)
}
