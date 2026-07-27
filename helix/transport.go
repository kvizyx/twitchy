package helix

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

type transportExecutor struct {
	httpClient *http.Client
	source     TokenSource
	policy     RateLimitPolicy
	clock      Clock
	sleeper    Sleeper
	refresh    *refreshCoordinator
}

func newTransportExecutor(httpClient *http.Client, source TokenSource, policy RateLimitPolicy, clock Clock, sleeper Sleeper) *transportExecutor {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if clock == nil {
		clock = wallClock{}
	}
	if sleeper == nil {
		sleeper = timerSleeper{}
	}
	return &transportExecutor{
		httpClient: httpClient,
		source:     source,
		policy:     policy,
		clock:      clock,
		sleeper:    sleeper,
		refresh:    newRefreshCoordinator(),
	}
}

func (executor *transportExecutor) execute(ctx context.Context, request *http.Request, operation manifest.Operation, credential CredentialSnapshot) (*http.Response, ResponseMeta, error) {
	if err := ctx.Err(); err != nil {
		return nil, ResponseMeta{}, err
	}
	if request == nil {
		return nil, ResponseMeta{}, fmt.Errorf("helix: nil request")
	}
	replaySafe := requestReplaySafe(request, operation)
	state := retryState{seen: make(map[retryCause]bool)}
	var retryAfter string
	for attempt := 1; attempt <= maxHTTPAttempts; attempt++ {
		attemptRequest, err := replayRequest(ctx, request, attempt, credential)
		if err != nil {
			return nil, ResponseMeta{}, err
		}
		response, err := executor.httpClient.Do(attemptRequest)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ResponseMeta{}, ctxErr
			}
			return nil, ResponseMeta{}, newTransportError(operation.OperationID, err, credential.AccessToken())
		}
		meta := newResponseMetaAt(response.StatusCode, response.Header, attempt, executor.clock.Now())
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			if retryAfter != "" && meta.Header().Get("Retry-After") == "" {
				meta.headers.Set("Retry-After", retryAfter)
			}
			return response, meta, nil
		}

		switch response.StatusCode {
		case http.StatusUnauthorized:
			if replaySafe && state.canReplay(retryUnauthorized, attempt) && executor.canRefresh(credential) {
				if err := drainAndClose(response.Body); err != nil {
					return nil, meta, err
				}
				credential, err = executor.refreshCredential(ctx, credential)
				if err != nil {
					return nil, meta, err
				}
				continue
			}
		case http.StatusServiceUnavailable:
			retryAfter = headerValue(response.Header, "Retry-After")
			if replaySafe && state.canReplay(retryUnavailable, attempt) {
				if err := drainAndClose(response.Body); err != nil {
					return nil, meta, err
				}
				continue
			}
		case http.StatusTooManyRequests:
			if replaySafe && operation.Replay.BucketWaitable && executor.policy.Wait && state.canReplay(retryRateLimit, attempt) && meta.RateLimit().Valid() {
				if err := drainAndClose(response.Body); err != nil {
					return nil, meta, err
				}
				if err := executor.waitForRateLimit(ctx, meta.RateLimit()); err != nil {
					return nil, meta, err
				}
				continue
			}
		}
		return nil, meta, executor.responseError(response, operation.OperationID, meta, credential.AccessToken())
	}
	return nil, ResponseMeta{}, fmt.Errorf("helix: request attempt limit exceeded")
}

func (executor *transportExecutor) canRefresh(credential CredentialSnapshot) bool {
	_, refreshable := executor.source.(RefreshableTokenSource)
	return credential.Refreshable() && credential.TokenClass() == TokenClassUser && refreshable
}

func (executor *transportExecutor) refreshCredential(ctx context.Context, credential CredentialSnapshot) (CredentialSnapshot, error) {
	source, ok := executor.source.(RefreshableTokenSource)
	if !ok {
		return CredentialSnapshot{}, ErrNotRefreshable
	}
	current, err := source.Token(ctx)
	if err != nil {
		return CredentialSnapshot{}, fmt.Errorf("helix: read current credential: %w", err)
	}
	if current.Generation() > credential.Generation() {
		return current, nil
	}
	refreshed, err := executor.refresh.refresh(ctx, source, credential)
	if err != nil {
		return CredentialSnapshot{}, fmt.Errorf("helix: refresh credential: %w", err)
	}
	return refreshed, nil
}

func (executor *transportExecutor) waitForRateLimit(ctx context.Context, rate RateLimit) error {
	duration := rate.Reset().Sub(executor.clock.Now())
	if duration <= 0 {
		return nil
	}
	if maxWait := executor.policy.maxWait(); duration > maxWait {
		duration = maxWait
	}
	return executor.sleeper.Sleep(ctx, duration)
}

func (executor *transportExecutor) responseError(response *http.Response, operation string, meta ResponseMeta, secret string) error {
	body, _, err := readErrorExcerpt(response.Body, BodyLimits{})
	closeErr := response.Body.Close()
	if err != nil {
		return newProtocolError(errorInput{operation: operation, statusCode: response.StatusCode, meta: meta, cause: err, secrets: []string{secret}})
	}
	if closeErr != nil {
		return newProtocolError(errorInput{operation: operation, statusCode: response.StatusCode, meta: meta, body: body, cause: closeErr, secrets: []string{secret}})
	}
	return newAPIError(errorInput{operation: operation, statusCode: response.StatusCode, meta: meta, body: body, secrets: []string{secret}})
}

func drainAndClose(body io.ReadCloser) error {
	_, readErr := io.CopyN(io.Discard, body, 4096)
	closeErr := body.Close()
	if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
		return fmt.Errorf("helix: drain response body: %w", readErr)
	}
	if closeErr != nil {
		return fmt.Errorf("helix: close response body: %w", closeErr)
	}
	return nil
}
