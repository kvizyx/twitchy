package oauth

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
)

type coordinatedOAuthServer struct {
	server         *httptest.Server
	refreshes      atomic.Int32
	validations    atomic.Int32
	refreshInputs  chan string
	refreshStarted chan struct{}

	mu        sync.Mutex
	refreshFn func(http.ResponseWriter, *http.Request, int32)
}

const coordinatedValidationResponse = `{"client_id":"client","login":"streamer",` +
	`"scopes":[],"user_id":"42","expires_in":3600}`

const coordinatedRotationResponse = `{"access_token":"rotated-access",` +
	`"refresh_token":"rotated-refresh","expires_in":3600,"token_type":"bearer"}`

func newCoordinatedOAuthServer(t *testing.T) *coordinatedOAuthServer {
	t.Helper()
	fixture := &coordinatedOAuthServer{
		refreshInputs:  make(chan string, 64),
		refreshStarted: make(chan struct{}, 64),
	}
	fixture.server = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/oauth2/validate":
			fixture.validations.Add(1)
			_, _ = io.WriteString(writer, coordinatedValidationResponse)
		case "/oauth2/token":
			if err := request.ParseForm(); err != nil {
				http.Error(writer, "invalid form", http.StatusBadRequest)
				return
			}
			count := fixture.refreshes.Add(1)
			fixture.refreshInputs <- request.Form.Get("refresh_token")
			fixture.refreshStarted <- struct{}{}
			fixture.mu.Lock()
			refreshFn := fixture.refreshFn
			fixture.mu.Unlock()
			if refreshFn != nil {
				refreshFn(writer, request, count)
				return
			}
			_, _ = io.WriteString(writer, coordinatedRotationResponse)
		default:
			http.Error(writer, "unexpected path", http.StatusNotFound)
		}
	}))
	t.Cleanup(fixture.server.Close)
	return fixture
}

func (fixture *coordinatedOAuthServer) setRefreshFn(refreshFn func(http.ResponseWriter, *http.Request, int32)) {
	fixture.mu.Lock()
	fixture.refreshFn = refreshFn
	fixture.mu.Unlock()
}

func newCoordinatedSource(
	t *testing.T,
	clock *testClock,
	server *coordinatedOAuthServer,
	coordinator RefreshCoordinator,
	store *coordinatedTestStore,
) *RefreshingTokenSource {
	t.Helper()
	client, err := New(WithBaseURL(server.server.URL), WithHTTPClient(server.server.Client()), WithClock(clock))
	if err != nil {
		t.Fatal(err)
	}
	registry, err := NewCoordinatedRegistry(client, coordinator)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = registry.Close() })
	if err = registry.AddCoordinatedUser(context.Background(), "42", store.load, store.hook, "chat"); err != nil {
		t.Fatal(err)
	}
	source, err := registry.SourceForUser("42")
	if err != nil {
		t.Fatal(err)
	}
	refreshing, ok := source.(*RefreshingTokenSource)
	if !ok {
		t.Fatalf("source type = %T", source)
	}
	return refreshing
}

func requireNoActivation(t *testing.T, source *RefreshingTokenSource, accessToken string) {
	t.Helper()
	source.mu.Lock()
	defer source.mu.Unlock()
	if got := source.current.AccessToken(); got != accessToken {
		t.Fatalf("active access token = %q, want unchanged", got)
	}
}
