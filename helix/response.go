package helix

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response[T any] struct {
	Data       T
	Pagination Pagination
	Meta       ResponseMeta
}

func (r *Response[T]) UnmarshalJSON(data []byte) error {
	type responseWire[T any] struct {
		Data       T          `json:"data"`
		Pagination Pagination `json:"pagination"`
	}
	var wire responseWire[T]
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	r.Data = wire.Data
	r.Pagination = wire.Pagination
	return nil
}

func (r Response[T]) MarshalJSON() ([]byte, error) {
	type responseWire[T any] struct {
		Data       T           `json:"data"`
		Pagination *Pagination `json:"pagination,omitempty"`
	}
	var pagination *Pagination
	if r.Pagination.Cursor() != "" {
		pagination = &r.Pagination
	}
	return json.Marshal(responseWire[T]{Data: r.Data, Pagination: pagination})
}

type Pagination struct {
	cursor string
}

type RateLimit struct {
	limit     int
	remaining int
	reset     time.Time
	valid     bool
}

func (r RateLimit) Limit() int       { return r.limit }
func (r RateLimit) Remaining() int   { return r.remaining }
func (r RateLimit) Reset() time.Time { return r.reset }
func (r RateLimit) Valid() bool      { return r.valid }

type ResponseMeta struct {
	statusCode int
	headers    http.Header
	requestID  string
	rateLimit  RateLimit
	attempts   int
}

func (m ResponseMeta) StatusCode() int      { return m.statusCode }
func (m ResponseMeta) Header() http.Header  { return cloneHeader(m.headers) }
func (m ResponseMeta) RequestID() string    { return m.requestID }
func (m ResponseMeta) RateLimit() RateLimit { return m.rateLimit }
func (m ResponseMeta) Attempts() int        { return m.attempts }

func cloneHeader(headers http.Header) http.Header {
	if headers == nil {
		return make(http.Header)
	}
	return headers.Clone()
}

func newResponseMetaAt(statusCode int, headers http.Header, attempts int, now time.Time) ResponseMeta {
	copyHeaders := cloneHeader(headers)
	return ResponseMeta{
		statusCode: statusCode,
		headers:    copyHeaders,
		requestID:  responseRequestID(copyHeaders),
		rateLimit:  parseRateLimitAt(copyHeaders, now),
		attempts:   attempts,
	}
}

func responseRequestID(headers http.Header) string {
	for _, name := range []string{"X-Request-ID", "Request-ID", "X-Correlation-ID"} {
		if value := headerValue(headers, name); value != "" {
			return value
		}
	}
	return ""
}

func headerValue(headers http.Header, name string) string {
	if value := headers.Get(name); value != "" {
		return value
	}
	for key, values := range headers {
		if strings.EqualFold(key, name) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func parseRateLimitAt(headers http.Header, now time.Time) RateLimit {
	limit, limitOK := parseNonNegativeInt(headerValue(headers, "Ratelimit-Limit"))
	remaining, remainingOK := parseNonNegativeInt(headerValue(headers, "Ratelimit-Remaining"))
	resetUnix, resetOK := parseNonNegativeInt64(headerValue(headers, "Ratelimit-Reset"))
	if !limitOK || !remainingOK || !resetOK {
		return RateLimit{}
	}
	reset := time.Unix(resetUnix, 0)
	if !reset.After(now) || remaining > limit {
		return RateLimit{}
	}
	return RateLimit{limit: limit, remaining: remaining, reset: reset, valid: true}
}

func parseNonNegativeInt(value string) (int, bool) {
	parsed, err := strconv.ParseInt(value, 10, 0)
	if err != nil || parsed < 0 {
		return 0, false
	}
	return int(parsed), true
}

func parseNonNegativeInt64(value string) (int64, bool) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 0 {
		return 0, false
	}
	return parsed, true
}
