package helix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const maxErrorExcerptBytes = 64 * 1024

var (
	ErrNotRefreshable   = errors.New("helix: credential is not refreshable")
	ErrCursorCycle      = errors.New("helix: pagination cursor cycle")
	ErrPageLimit        = errors.New("helix: pagination page limit reached")
	ErrCredentialCommit = errors.New("helix: credential commit failed")
	ErrInvalidSession   = errors.New("helix: invalid session")
	ErrSessionClosed    = errors.New("helix: session is closed")
)

var (
	bearerSecretPattern = regexp.MustCompile(`(?i)(\bbearer\s+)[^\s,]+`)
	credentialPattern   = regexp.MustCompile(`(?i)(["']?(?:access[_-]?token|refresh[_-]?token|client[_-]?secret|authorization|token)["']?\s*[:=]\s*["']?)[^"',\s}]+`)
	headerSecretPattern = regexp.MustCompile(`(?i)(?:authorization|client-secret|refresh-token):\s*[^\s,]+`)
	querySecretPattern  = regexp.MustCompile(`(?i)([?&](?:access_token|refresh_token|client_secret|client-secret|token)=)[^&\s]+`)
)

// TransportError reports a failure before an HTTP response was received.
type TransportError struct {
	operation string
	cause     error
	secrets   []string
}

func (e *TransportError) Error() string {
	if e == nil {
		return "helix transport error: <nil>"
	}
	return formatError("transport error", e.operation, 0, "", e.cause, e.secrets)
}

func (e *TransportError) Operation() string { return e.operation }
func (e *TransportError) Unwrap() error     { return e.cause }

type errorDetails struct {
	operation string
	status    int
	meta      ResponseMeta
	cause     error
	excerpt   string
	secrets   []string
}

func (e errorDetails) Error(kind string) string {
	return formatError(kind, e.operation, e.status, e.excerpt, e.cause, e.secrets)
}

func (e errorDetails) Unwrap() error { return e.cause }

// ProtocolError reports a response that violated the declared wire format.
type ProtocolError struct{ details errorDetails }

func (e *ProtocolError) Error() string {
	if e == nil {
		return "helix protocol error: <nil>"
	}
	return e.details.Error("protocol error")
}

func (e *ProtocolError) Operation() string  { return e.details.operation }
func (e *ProtocolError) StatusCode() int    { return e.details.status }
func (e *ProtocolError) Meta() ResponseMeta { return e.details.meta }
func (e *ProtocolError) Unwrap() error      { return e.details.Unwrap() }

// APIError reports a non-authentication, non-rate-limit Helix response error.
type APIError struct{ details errorDetails }

func (e *APIError) Error() string {
	if e == nil {
		return "helix API error: <nil>"
	}
	return e.details.Error("API error")
}

func (e *APIError) Operation() string  { return e.details.operation }
func (e *APIError) StatusCode() int    { return e.details.status }
func (e *APIError) Meta() ResponseMeta { return e.details.meta }
func (e *APIError) Unwrap() error      { return e.details.Unwrap() }

// AuthError reports an authentication or authorization HTTP status.
type AuthError struct{ details errorDetails }

func (e *AuthError) Error() string {
	if e == nil {
		return "helix auth error: <nil>"
	}
	return e.details.Error("auth error")
}

func (e *AuthError) Operation() string  { return e.details.operation }
func (e *AuthError) StatusCode() int    { return e.details.status }
func (e *AuthError) Meta() ResponseMeta { return e.details.meta }
func (e *AuthError) Unwrap() error      { return e.details.Unwrap() }

// RateLimitError reports a response rejected by a rate limit.
type RateLimitError struct{ details errorDetails }

func (e *RateLimitError) Error() string {
	if e == nil {
		return "helix rate-limit error: <nil>"
	}
	return e.details.Error("rate-limit error")
}

func (e *RateLimitError) Operation() string  { return e.details.operation }
func (e *RateLimitError) StatusCode() int    { return e.details.status }
func (e *RateLimitError) Meta() ResponseMeta { return e.details.meta }
func (e *RateLimitError) Unwrap() error      { return e.details.Unwrap() }

type errorInput struct {
	operation  string
	statusCode int
	meta       ResponseMeta
	body       []byte
	cause      error
	secrets    []string
}

func newTransportError(operation string, cause error, secrets ...string) *TransportError {
	return &TransportError{operation: operation, cause: cause, secrets: append([]string(nil), secrets...)}
}

func newProtocolError(input errorInput) *ProtocolError {
	return &ProtocolError{details: errorDetails{
		operation: input.operation,
		status:    input.statusCode,
		meta:      input.meta,
		cause:     input.cause,
		excerpt:   boundedRedactedExcerpt(input.body, input.secrets),
		secrets:   append([]string(nil), input.secrets...),
	}}
}

func newAPIError(input errorInput) error {
	payload, parseErr := parseAPIErrorPayload(input.body)
	if parseErr != nil {
		if input.cause == nil {
			input.cause = parseErr
		}
	} else {
		input.body = []byte(payload.Message)
	}
	details := errorDetails{
		operation: input.operation,
		status:    input.statusCode,
		meta:      input.meta,
		cause:     input.cause,
		excerpt:   boundedRedactedExcerpt(input.body, input.secrets),
		secrets:   append([]string(nil), input.secrets...),
	}
	switch input.statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthError{details: details}
	case http.StatusTooManyRequests:
		return &RateLimitError{details: details}
	default:
		return &APIError{details: details}
	}
}

type twitchAPIError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func parseAPIErrorPayload(body []byte) (twitchAPIError, error) {
	if len(body) == 0 {
		return twitchAPIError{}, errors.New("empty API error body")
	}
	if len(body) > maxErrorExcerptBytes {
		return twitchAPIError{}, errors.New("API error body exceeds limit")
	}
	var payload twitchAPIError
	if err := json.Unmarshal(body, &payload); err != nil {
		return twitchAPIError{}, fmt.Errorf("decode API error: %w", err)
	}
	if payload.Error == "" || payload.Status <= 0 || payload.Message == "" {
		return twitchAPIError{}, errors.New("malformed API error payload")
	}
	return payload, nil
}

func boundedRedactedExcerpt(body []byte, secrets []string) string {
	if len(body) == 0 {
		return ""
	}
	redacted := redactSecrets(string(body), secrets)
	if len(redacted) > maxErrorExcerptBytes {
		redacted = redacted[:maxErrorExcerptBytes]
	}
	return strings.ToValidUTF8(redacted, "�")
}

func redactSecrets(value string, secrets []string) string {
	for _, secret := range secrets {
		if secret != "" {
			value = strings.ReplaceAll(value, secret, "[redacted]")
		}
	}
	value = bearerSecretPattern.ReplaceAllString(value, `${1}[redacted]`)
	value = credentialPattern.ReplaceAllString(value, `${1}[redacted]`)
	value = headerSecretPattern.ReplaceAllString(value, "[redacted]")
	return querySecretPattern.ReplaceAllString(value, `${1}[redacted]`)
}

func formatError(kind, operation string, status int, excerpt string, cause error, secrets []string) string {
	parts := []string{"helix " + kind}
	if operation != "" {
		parts = append(parts, "operation="+redactSecrets(operation, secrets))
	}
	if status != 0 {
		parts = append(parts, fmt.Sprintf("status=%d", status))
	}
	if excerpt != "" {
		parts = append(parts, "excerpt="+redactSecrets(excerpt, secrets))
	}
	if cause != nil {
		parts = append(parts, "cause="+redactSecrets(cause.Error(), secrets))
	}
	return strings.Join(parts, ": ")
}
