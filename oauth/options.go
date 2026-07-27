package oauth

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/kvizyx/twitchy/helix"
)

type Option func(*Client) error

func WithHTTPClient(caller *http.Client) Option {
	return func(client *Client) error {
		if caller == nil {
			return ErrInvalidOption
		}
		clone := *caller
		clone.CheckRedirect = func(*http.Request, []*http.Request) error {
			return errors.New("oauth: redirect rejected")
		}
		client.httpClient = &clone
		return nil
	}
}

func WithBaseURL(raw string) Option {
	return func(client *Client) error {
		parsed, err := parseBaseURL(raw)
		if err != nil {
			return ErrInvalidOption
		}
		client.baseURL = parsed
		return nil
	}
}

func WithRandom(randomSource io.Reader) Option {
	return func(client *Client) error {
		if isNilValue(randomSource) {
			return ErrInvalidOption
		}
		client.random = randomSource
		return nil
	}
}

func WithClock(clock helix.Clock) Option {
	return func(client *Client) error {
		if isNilValue(clock) {
			return ErrInvalidOption
		}
		client.clock = clock
		return nil
	}
}

func isNilValue(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func parseBaseURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, ErrInvalidOption
	}
	parsed.Scheme = strings.ToLower(parsed.Scheme)
	if parsed.Scheme != "https" && !isLoopbackHTTP(parsed) {
		return nil, ErrInvalidOption
	}
	return parsed, nil
}

func isLoopbackHTTP(parsed *url.URL) bool {
	if parsed.Scheme != "http" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
