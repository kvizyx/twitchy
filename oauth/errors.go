package oauth

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidOption    = errors.New("oauth: invalid option")
	ErrInvalidClient    = errors.New("oauth: invalid client")
	ErrStateMismatch    = errors.New("oauth: authorization state mismatch")
	ErrRedirectMismatch = errors.New("oauth: authorization redirect mismatch")
	ErrDeviceExpired    = errors.New("oauth: device authorization expired")
)

type TransportError struct {
	operation string
	cause     error
}

func (e *TransportError) Error() string {
	if e == nil {
		return "oauth transport error: <nil>"
	}
	return fmt.Sprintf("oauth transport error during %s", e.operation)
}

func (e *TransportError) Unwrap() error     { return e.cause }
func (e *TransportError) Operation() string { return e.operation }

type ProtocolError struct {
	operation string
	status    int
	cause     error
}

func (e *ProtocolError) Error() string {
	if e == nil {
		return "oauth protocol error: <nil>"
	}
	return fmt.Sprintf("oauth protocol error during %s (status %d)", e.operation, e.status)
}

func (e *ProtocolError) Unwrap() error     { return e.cause }
func (e *ProtocolError) Operation() string { return e.operation }
func (e *ProtocolError) StatusCode() int   { return e.status }

type OAuthError struct {
	operation   string
	status      int
	code        string
	description string
	retryable   bool
	cause       error
}

func (e *OAuthError) Error() string {
	if e == nil {
		return "oauth error: <nil>"
	}
	message := fmt.Sprintf("oauth error during %s (status %d)", e.operation, e.status)
	if e.code != "" {
		message += ": " + e.code
	}
	if e.description != "" {
		message += ": " + e.description
	}
	return message
}

func (e *OAuthError) Unwrap() error       { return e.cause }
func (e *OAuthError) Operation() string   { return e.operation }
func (e *OAuthError) StatusCode() int     { return e.status }
func (e *OAuthError) Code() string        { return e.code }
func (e *OAuthError) Description() string { return e.description }
func (e *OAuthError) Retryable() bool     { return e.retryable }

type DeviceAuthorizationError struct {
	oauthError *OAuthError
}

func (e *DeviceAuthorizationError) Error() string {
	if e == nil || e.oauthError == nil {
		return "oauth device authorization error: <nil>"
	}
	return e.oauthError.Error()
}

func (e *DeviceAuthorizationError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.oauthError
}

func (e *DeviceAuthorizationError) Operation() string   { return e.oauthError.Operation() }
func (e *DeviceAuthorizationError) StatusCode() int     { return e.oauthError.StatusCode() }
func (e *DeviceAuthorizationError) Code() string        { return e.oauthError.Code() }
func (e *DeviceAuthorizationError) Description() string { return e.oauthError.Description() }
func (e *DeviceAuthorizationError) Retryable() bool     { return e.oauthError.Retryable() }

func sanitizeText(value string, secrets ...string) string {
	for _, secret := range secrets {
		if secret != "" {
			value = strings.ReplaceAll(value, secret, "[redacted]")
		}
	}
	return value
}
