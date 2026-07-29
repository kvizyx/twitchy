package oauth

import "github.com/kvizyx/twitchy/helix"

// CredentialRotationError means a coordinated remote refresh may have rotated
// the durable refresh token without returning a pair safe to persist. The user
// must reauthorize because retrying the previous refresh token is unsafe.
type CredentialRotationError struct {
	cause error
}

func (err *CredentialRotationError) Error() string {
	return "oauth: credential rotation outcome unknown; reauthorization required"
}

func (err *CredentialRotationError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

func (err *CredentialRotationError) RequiresReauthorization() bool { return true }

type coordinatedCommitError struct {
	cause error
}

func (err *coordinatedCommitError) Error() string {
	return "oauth: rotated credential was not durably committed"
}

func (err *coordinatedCommitError) Unwrap() []error {
	if err == nil {
		return nil
	}
	return []error{helix.ErrCredentialCommit, err.cause}
}

type coordinatedLoadError struct {
	cause error
}

func (err *coordinatedLoadError) Error() string {
	return "oauth: durable credential reload failed"
}

func (err *coordinatedLoadError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

func wrapCoordinatedLoadError(err error) error {
	if err == nil {
		return nil
	}
	return &coordinatedLoadError{cause: err}
}
