package oauth

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

const defaultBaseURL = "https://id.twitch.tv"

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	random     io.Reader
	clock      helix.Clock
	valid      bool
	sleeper    func(context.Context, time.Duration) error
}

type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }

func New(options ...Option) (*Client, error) {
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{
		httpClient: cloneHTTPClient(http.DefaultClient),
		baseURL:    baseURL,
		random:     rand.Reader,
		clock:      wallClock{},
		sleeper:    sleepWithTimer,
	}
	for _, option := range options {
		if option == nil {
			return nil, ErrInvalidOption
		}
		if err := option(client); err != nil {
			return nil, err
		}
	}
	if client.httpClient == nil || client.baseURL == nil || client.random == nil || client.clock == nil {
		return nil, ErrInvalidClient
	}
	client.valid = true
	return client, nil
}

func cloneHTTPClient(caller *http.Client) *http.Client {
	clone := *caller
	clone.CheckRedirect = func(*http.Request, []*http.Request) error {
		return errors.New("oauth: redirect rejected")
	}
	return &clone
}

func sleepWithTimer(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (client *Client) validClient() error {
	if client == nil || !client.valid {
		return ErrInvalidClient
	}
	return nil
}

func (client *Client) endpoint(name string) string {
	base := *client.baseURL
	base.RawQuery = ""
	base.Fragment = ""
	base.Path = strings.TrimSuffix(base.Path, "/")
	if !strings.HasSuffix(base.Path, "/oauth2") {
		base.Path += "/oauth2"
	}
	base.Path += "/" + strings.TrimPrefix(name, "/")
	return base.String()
}
