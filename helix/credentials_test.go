package helix

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kvizyx/twitchy/helix/internal/manifest"
)

func TestCredentialSnapshot_copiesCredentialState(t *testing.T) {
	// Given a credential whose caller-owned scopes can be mutated
	scopes := []AuthorizationScope{ScopeUserReadChat}
	credential := Credential{
		AccessToken: "access-token-secret",
		ClientID:    "client-id",
		TokenClass:  TokenClassUser,
		UserID:      "user-id",
		Scopes:      scopes,
		ExpiresAt:   time.Unix(123, 0),
		Refreshable: true,
		Generation:  7,
	}

	// When the snapshot is created and returned slices are mutated
	snapshot := NewCredentialSnapshot(credential)
	scopes[0] = ScopeUserWriteChat
	returnedScopes := snapshot.Scopes()
	returnedScopes[0] = ScopeUserManageChatColor

	// Then all snapshot values remain bound to their original immutable state
	if snapshot.AccessToken() != credential.AccessToken {
		t.Fatalf("AccessToken() = %q, want %q", snapshot.AccessToken(), credential.AccessToken)
	}
	if snapshot.ClientID() != credential.ClientID || snapshot.TokenClass() != credential.TokenClass || snapshot.UserID() != credential.UserID {
		t.Fatal("snapshot identity fields changed")
	}
	if got := snapshot.Scopes(); len(got) != 1 || got[0] != ScopeUserReadChat {
		t.Fatalf("Scopes() = %#v, want original scope", got)
	}
	if !snapshot.ExpiresAt().Equal(credential.ExpiresAt) || !snapshot.Refreshable() || snapshot.Generation() != credential.Generation {
		t.Fatal("snapshot metadata changed")
	}
}

func TestCredentialFormatting_redactsAccessToken(t *testing.T) {
	// Given credentials containing a secret
	credential := Credential{AccessToken: "access-token-secret"}
	snapshot := NewCredentialSnapshot(credential)

	// When credentials are formatted for diagnostics
	formatted := fmt.Sprintf("%v %#v %v %#v", credential, credential, snapshot, snapshot)

	// Then no formatted representation contains the access token
	if strings.Contains(formatted, credential.AccessToken) {
		t.Fatalf("formatted credential contains access token: %s", formatted)
	}
}

func TestTokenSource_static_returnsSnapshotAndHonorsContext(t *testing.T) {
	// Given a static source and a canceled context
	source := NewStaticTokenSource(Credential{AccessToken: "access-token", Generation: 3})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// When requesting a token with cancellation
	_, err := source.Token(ctx)

	// Then cancellation wins before returning credentials
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Token() error = %v, want context.Canceled", err)
	}

	// When requesting the static token with an active context
	snapshot, err := source.Token(context.Background())

	// Then the source returns the immutable credential and is not refreshable
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if snapshot.AccessToken() != "access-token" || snapshot.Generation() != 3 {
		t.Fatalf("unexpected static snapshot: %#v", snapshot)
	}
	var tokenSource TokenSource = source
	if _, ok := tokenSource.(RefreshableTokenSource); ok {
		t.Fatal("static source unexpectedly implements RefreshableTokenSource")
	}
}

func TestCredentialValidation_rejectsKnownManifestConstraints(t *testing.T) {
	// Given a manifest operation with deterministic authorization requirements
	operation := manifest.Operation{
		OperationID:    "known-operation",
		TokenClasses:   []manifest.TokenClass{manifest.TokenClassUser},
		Scopes:         []string{"user:read:chat"},
		SubjectBinding: "user_id",
	}
	valid := NewCredentialSnapshot(Credential{
		ClientID:   "client-a",
		TokenClass: TokenClassUser,
		UserID:     "user-a",
		Scopes:     []AuthorizationScope{ScopeUserReadChat},
	})

	tests := []struct {
		name       string
		credential CredentialSnapshot
		clientID   string
		subjectID  string
	}{
		{
			name: "wrong token class",
			credential: NewCredentialSnapshot(Credential{
				ClientID:   "client-a",
				TokenClass: TokenClassApp,
				Scopes:     []AuthorizationScope{ScopeUserReadChat},
			}),
			clientID:  "client-a",
			subjectID: "user-a",
		},
		{name: "client ID mismatch", credential: valid, clientID: "client-b", subjectID: "user-a"},
		{
			name: "missing scope",
			credential: NewCredentialSnapshot(Credential{
				ClientID:   "client-a",
				TokenClass: TokenClassUser,
				UserID:     "user-a",
			}),
			clientID:  "client-a",
			subjectID: "user-a",
		},
		{name: "subject mismatch", credential: valid, clientID: "client-a", subjectID: "user-b"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// When local authorization is checked before request execution
			err := validateCredentialForOperation(test.credential, operation, test.clientID, test.subjectID)

			// Then a typed AuthError is returned without exposing credentials
			var authErr *AuthError
			if !errors.As(err, &authErr) {
				t.Fatalf("error = %v, want AuthError", err)
			}
			if test.credential.AccessToken() != "" && strings.Contains(err.Error(), test.credential.AccessToken()) {
				t.Fatal("authorization error exposed access token")
			}
		})
	}
}

func TestCredentialValidation_passesUnknownManifestRules(t *testing.T) {
	// Given a manifest operation whose server-enforced rules are unknown
	operation := manifest.Operation{
		OperationID:    "unknown-operation",
		TokenClasses:   []manifest.TokenClass{manifest.TokenClassUnknown},
		Scopes:         []string{"unknown"},
		SubjectBinding: "unknown",
	}
	credential := NewCredentialSnapshot(Credential{TokenClass: TokenClassApp})

	// When local authorization is checked
	err := validateCredentialForOperation(credential, operation, "", "different-user")

	// Then unknown server rules are not guessed locally
	if err != nil {
		t.Fatalf("validateCredentialForOperation() error = %v, want nil", err)
	}
}

func TestCredentialSnapshot_isSafeForConcurrentReads(t *testing.T) {
	// Given one immutable snapshot shared by concurrent callers
	snapshot := NewCredentialSnapshot(Credential{AccessToken: "access-token", Scopes: []AuthorizationScope{ScopeUserReadChat}})
	source := NewStaticTokenSource(Credential{AccessToken: "access-token", Scopes: []AuthorizationScope{ScopeUserReadChat}})

	// When many goroutines read the snapshot and source
	var waitGroup sync.WaitGroup
	for index := 0; index < 32; index++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for count := 0; count < 100; count++ {
				_ = snapshot.Scopes()
				_, _ = source.Token(context.Background())
			}
		}()
	}
	waitGroup.Wait()

	// Then all concurrent reads complete without mutation or error
	if snapshot.AccessToken() != "access-token" {
		t.Fatal("snapshot access token changed")
	}
}
