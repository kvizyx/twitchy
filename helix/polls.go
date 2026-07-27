package helix

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type PollStatus string

const (
	PollStatusActive     PollStatus = "ACTIVE"
	PollStatusCompleted  PollStatus = "COMPLETED"
	PollStatusTerminated PollStatus = "TERMINATED"
	PollStatusArchived   PollStatus = "ARCHIVED"
	PollStatusModerated  PollStatus = "MODERATED"
	PollStatusInvalid    PollStatus = "INVALID"
)

type PollChoice struct {
	ID                 string `json:"id,omitempty"`
	Title              string `json:"title"`
	Votes              int64  `json:"votes,omitempty"`
	ChannelPointsVotes int64  `json:"channel_points_votes,omitempty"`
	BitsVotes          int64  `json:"bits_votes,omitempty"`
}

type Poll struct {
	ID                         string       `json:"id"`
	BroadcasterID              string       `json:"broadcaster_id"`
	BroadcasterName            string       `json:"broadcaster_name"`
	BroadcasterLogin           string       `json:"broadcaster_login"`
	Title                      string       `json:"title"`
	Choices                    []PollChoice `json:"choices"`
	BitsVotingEnabled          bool         `json:"bits_voting_enabled"`
	BitsPerVote                int          `json:"bits_per_vote"`
	ChannelPointsVotingEnabled bool         `json:"channel_points_voting_enabled"`
	ChannelPointsPerVote       int          `json:"channel_points_per_vote"`
	Status                     PollStatus   `json:"status"`
	Duration                   int          `json:"duration"`
	StartedAt                  Timestamp    `json:"started_at"`
	EndedAt                    *Timestamp   `json:"ended_at"`
}

type GetPollsRequest struct {
	BroadcasterID string   `query:"broadcaster_id"`
	IDs           []string `query:"id,omitempty"`
	First         *int     `query:"first,omitempty"`
	After         *string  `query:"after,omitempty"`
}

type GetPollsData []Poll

type CreatePollRequest struct {
	BroadcasterID              string       `query:"-" json:"broadcaster_id"`
	Title                      string       `query:"-" json:"title"`
	Choices                    []PollChoice `query:"-" json:"choices"`
	Duration                   int          `query:"-" json:"duration"`
	ChannelPointsVotingEnabled *bool        `query:"-" json:"channel_points_voting_enabled,omitempty"`
	ChannelPointsPerVote       *int         `query:"-" json:"channel_points_per_vote,omitempty"`
}

type CreatePollData []Poll

type EndPollRequest struct {
	BroadcasterID string     `query:"-" json:"broadcaster_id"`
	ID            string     `query:"-" json:"id"`
	Status        PollStatus `query:"-" json:"status"`
}

type EndPollData []Poll

func (s *PollsService) GetPolls(ctx context.Context, req GetPollsRequest) (*Response[GetPollsData], error) {
	return executeTask23Endpoint[GetPollsData](s.client, ctx, "get-polls", req, nil, req.BroadcasterID, ScopeChannelReadPolls, ScopeChannelManagePolls)
}

func (s *PollsService) GetPollsPager(req GetPollsRequest, opts ...PagerOption) (*Pager[GetPollsData], error) {
	plan, err := newPaginationPlan(req, "after")
	if err != nil {
		return nil, err
	}
	return newPager(func(ctx context.Context, cursor string) (*Response[GetPollsData], error) {
		next, err := plan.withCursor(req, cursor)
		if err != nil {
			return nil, err
		}
		return s.GetPolls(ctx, next.(GetPollsRequest))
	}, opts...)
}

func (s *PollsService) CreatePoll(ctx context.Context, req CreatePollRequest) (*Response[CreatePollData], error) {
	return executeTask23Endpoint[CreatePollData](s.client, ctx, "create-poll", req, req, req.BroadcasterID)
}

func (s *PollsService) EndPoll(ctx context.Context, req EndPollRequest) (*Response[EndPollData], error) {
	return executeTask23Endpoint[EndPollData](s.client, ctx, "end-poll", req, req, req.BroadcasterID)
}

func executeTask23Endpoint[T any](client *Client, ctx context.Context, anchor string, query, body any, broadcasterID string, readScopes ...AuthorizationScope) (*Response[T], error) {
	if err := client.validClient(); err != nil {
		return nil, err
	}
	operation, err := manifest.OperationByAnchor(anchor)
	if err != nil {
		return nil, err
	}
	credential, err := client.credential(ctx)
	if err != nil {
		return nil, err
	}
	if !operationAllowsTokenClass(credential.TokenClass(), operation.TokenClasses) || credential.TokenClass() != TokenClassUser || credential.UserID() != broadcasterID {
		return nil, localCredentialAuthError(operation.OperationID)
	}
	if len(readScopes) > 0 {
		allowed := false
		for _, scope := range readScopes {
			if snapshotHasScope(credential, scope) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, localCredentialAuthError(operation.OperationID)
		}
	} else if err := validateCredentialForOperation(credential, operation, "", ""); err != nil {
		return nil, err
	}
	request, err := buildRequest(requestSpec{Context: ctx, Method: operation.Method, URL: client.endpointURL(operation.Path), Query: query, Body: body})
	if err != nil {
		return nil, fmt.Errorf("build %s request: %w", operation.OperationID, err)
	}
	request.Header.Set("User-Agent", client.userAgent)
	response, meta, err := client.executor.execute(ctx, request, operation, credential)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := decodeResponse[T](response.StatusCode, response.Body, DecodeOptions{})
	if err != nil {
		return nil, newProtocolError(errorInput{operation: operation.OperationID, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{credential.AccessToken()}})
	}
	result.Meta = meta
	return result, nil
}
