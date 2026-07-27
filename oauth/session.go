package oauth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/kvizyx/twitchy/helix"
)

const validationInterval = time.Hour

type SessionOption func(*sessionOptions) error

type sessionOptions struct {
	interval time.Duration
	hook     func(error)
}

func WithValidationInterval(interval time.Duration) SessionOption {
	return func(options *sessionOptions) error {
		options.interval = interval
		return nil
	}
}

func WithValidationErrorHook(hook func(error)) SessionOption {
	return func(options *sessionOptions) error {
		if hook == nil {
			return ErrInvalidOption
		}
		options.hook = hook
		return nil
	}
}

type ManagedSession struct {
	source   *RefreshingTokenSource
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
	hook     func(error)
	started  chan struct{}
}

type timerClock interface {
	NewTimer(time.Duration) helix.Timer
}

func NewManagedSession(ctx context.Context, source *RefreshingTokenSource, options ...SessionOption) (*ManagedSession, error) {
	if ctx == nil || source == nil {
		return nil, ErrInvalidOption
	}
	configuration := sessionOptions{interval: validationInterval}
	for _, option := range options {
		if option == nil {
			return nil, ErrInvalidOption
		}
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}
	if configuration.interval <= 0 {
		return nil, ErrInvalidOption
	}
	if configuration.interval != validationInterval {
		if _, productionClock := source.clock.(wallClock); productionClock {
			return nil, ErrInvalidOption
		}
		if _, ok := source.clock.(timerClock); !ok {
			return nil, ErrInvalidOption
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := source.validate(ctx); err != nil {
		if source.invalidValidation(err) {
			source.terminalize(helix.ErrInvalidSession)
			if source.invalidSessionStatus(err) {
				return nil, helix.ErrInvalidSession
			}
		}
		return nil, err
	}
	sessionCtx, cancel := context.WithCancel(ctx)
	session := &ManagedSession{source: source, ctx: sessionCtx, cancel: cancel, interval: configuration.interval, hook: configuration.hook, started: make(chan struct{})}
	go session.run()
	<-session.started
	return session, nil
}

func (source *RefreshingTokenSource) validate(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	source.mu.Lock()
	if err := source.lifecycleErrorLocked(); err != nil {
		source.mu.Unlock()
		return err
	}
	token := source.current.AccessToken()
	source.mu.Unlock()
	validationContext, cancel := context.WithCancel(ctx)
	stopSourceCancel := context.AfterFunc(source.ctx, cancel)
	defer func() {
		stopSourceCancel()
		cancel()
	}()
	validation, err := source.client.Validate(validationContext, ValidateRequest{AccessToken: token})
	if err != nil {
		return err
	}
	source.applyValidation(*validation)
	return nil
}

func (source *RefreshingTokenSource) invalidValidation(err error) bool {
	if errors.Is(err, helix.ErrInvalidSession) {
		return true
	}
	var protocolErr *ProtocolError
	if errors.As(err, &protocolErr) {
		return protocolErr.StatusCode() == 200
	}
	var oauthErr *OAuthError
	if errors.As(err, &oauthErr) {
		status := oauthErr.StatusCode()
		return status == http.StatusUnauthorized || (status >= 400 && status < 500 && status != http.StatusTooManyRequests)
	}
	return false
}

func (source *RefreshingTokenSource) invalidSessionStatus(err error) bool {
	var oauthErr *OAuthError
	return errors.As(err, &oauthErr) && oauthErr.StatusCode() == http.StatusUnauthorized
}

func (source *RefreshingTokenSource) terminalize(err error) {
	source.mu.Lock()
	if !source.closed && source.terminal == nil {
		source.terminal = err
	}
	source.mu.Unlock()
}

func (session *ManagedSession) run() {
	provider, ok := session.source.clock.(timerClock)
	var timer helix.Timer
	if ok {
		timer = provider.NewTimer(session.interval)
	} else {
		timer = wallClock{}.NewTimer(session.interval)
	}
	close(session.started)
	defer timer.Stop()
	for {
		select {
		case <-timer.C():
			timer.Stop()
			if err := session.source.validate(session.ctx); err != nil {
				if session.source.invalidValidation(err) {
					session.source.terminalize(helix.ErrInvalidSession)
					return
				}
				if session.hook != nil {
					session.hook(err)
				}
			}
			if err := session.source.lifecycleError(); err != nil {
				return
			}
			if ok {
				timer = provider.NewTimer(session.interval)
			} else {
				timer = wallClock{}.NewTimer(session.interval)
			}
		case <-session.ctx.Done():
			_ = session.source.Close()
			return
		case <-session.source.done:
			return
		}
	}
}

func (source *RefreshingTokenSource) lifecycleError() error {
	source.mu.Lock()
	defer source.mu.Unlock()
	return source.lifecycleErrorLocked()
}

func (session *ManagedSession) Close() error {
	session.cancel()
	return session.source.Close()
}
