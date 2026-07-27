package helix

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
)

type Option func(*Client) error

func WithHTTPClient(caller *http.Client) Option {
	return func(client *Client) error {
		if caller == nil {
			return ErrInvalidOption
		}
		client.httpClient = cloneHTTPClient(caller)
		return nil
	}
}

func WithBaseURL(raw string) Option {
	return func(client *Client) error {
		baseURL, err := parseBaseURL(raw)
		if err != nil {
			return ErrInvalidOption
		}
		client.baseURL = baseURL
		return nil
	}
}

func WithUserAgent(userAgent string) Option {
	return func(client *Client) error {
		client.userAgent = userAgent
		return nil
	}
}

func WithTokenSource(source TokenSource) Option {
	return func(client *Client) error {
		if isNilTokenSource(source) {
			return ErrInvalidOption
		}
		client.tokenSource = source
		return nil
	}
}

func WithStaticToken(credential Credential) Option {
	return func(client *Client) error {
		copy := credential
		copy.Scopes = append([]AuthorizationScope(nil), credential.Scopes...)
		client.staticToken = &copy
		return nil
	}
}

func WithRateLimitPolicy(policy RateLimitPolicy) Option {
	return func(client *Client) error {
		client.rateLimitPolicy = policy
		return nil
	}
}

func cloneHTTPClient(caller *http.Client) *http.Client {
	if caller == nil {
		return nil
	}
	clone := *caller
	clone.CheckRedirect = func(*http.Request, []*http.Request) error {
		return errors.New("helix: redirect rejected")
	}
	return &clone
}

func parseBaseURL(raw string) (*url.URL, error) {
	baseURL, err := url.Parse(raw)
	if err != nil || baseURL.Scheme == "" || baseURL.Host == "" || baseURL.User != nil || baseURL.Fragment != "" || baseURL.RawQuery != "" {
		return nil, ErrInvalidOption
	}
	if baseURL.Scheme != "http" && baseURL.Scheme != "https" {
		return nil, ErrInvalidOption
	}
	return baseURL, nil
}

func isNilTokenSource(source TokenSource) bool {
	if source == nil {
		return true
	}
	value := reflect.ValueOf(source)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
