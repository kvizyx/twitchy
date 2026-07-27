package helix

import (
	"context"
	"strings"
)

type SendWhisperRequest struct {
	FromUserID string `query:"from_user_id" json:"-"`
	ToUserID   string `query:"to_user_id" json:"-"`
	Message    string `query:"-" json:"message"`
}

type SendWhisperData struct{}

type whisperError struct {
	cause  error
	secret string
}

func (e *whisperError) Error() string {
	return strings.ReplaceAll(e.cause.Error(), e.secret, "[redacted]")
}

func (e *whisperError) Unwrap() error { return e.cause }

func (s *WhispersService) SendWhisper(ctx context.Context, req SendWhisperRequest) (*Response[SendWhisperData], error) {
	result, err := executeTask26EndpointWithSecrets[SendWhisperData](s.client, ctx, "send-whisper", req, req, task26Authorization{scopeSets: [][]AuthorizationScope{{ScopeUserManageWhispers}}, subjectID: req.FromUserID}, []string{req.Message})
	if err != nil {
		return nil, &whisperError{cause: err, secret: req.Message}
	}
	return result, nil
}
