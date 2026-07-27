package oauth

import (
	"time"

	"github.com/kvizyx/twitchy/helix"
)

type ResponseType string

const (
	ResponseTypeCode  ResponseType = "code"
	ResponseTypeToken ResponseType = "token"
)

type AuthorizationRequest struct {
	ClientID     string
	RedirectURI  string
	ResponseType ResponseType
	Scopes       []helix.AuthorizationScope
	State        string
	ForceVerify  bool
}

type Authorization struct {
	URL          string
	State        string
	RedirectURI  string
	ResponseType ResponseType
}

type AuthorizationResult struct {
	Code        string
	AccessToken string
	Scopes      []helix.AuthorizationScope
	State       string
	TokenType   string
}

type ClientCredentialsRequest struct {
	ClientID     string
	ClientSecret string
}

type ExchangeCodeRequest struct {
	ClientID     string
	ClientSecret string
	Code         string
	RedirectURI  string
}

type DeviceAuthorizationRequest struct {
	ClientID string
	Scopes   []helix.AuthorizationScope
}

type DeviceAuthorization struct {
	DeviceCode      string
	ExpiresIn       time.Duration
	Interval        time.Duration
	UserCode        string
	VerificationURI string
}

type PollDeviceTokenRequest struct {
	ClientID   string
	Scopes     []helix.AuthorizationScope
	DeviceCode string
	ExpiresAt  time.Time
}

type RefreshRequest struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
	Scopes       []helix.AuthorizationScope
	TokenType    string
}

type ValidateRequest struct {
	AccessToken string
}

type Validation struct {
	ClientID  string
	Login     string
	Scopes    []helix.AuthorizationScope
	UserID    string
	ExpiresIn time.Duration
}

type RevokeRequest struct {
	ClientID    string
	AccessToken string
}
