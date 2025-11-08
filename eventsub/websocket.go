package eventsub

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/coder/websocket"
	"github.com/kvizyx/twitchy/eventsub/eventtracker"
	"github.com/kvizyx/twitchy/internal/json"
)

var (
	ErrConnectionUnused = errors.New("connection is unused as subscription time limit exceeded")
	ErrReconnectTimeout = errors.New("reconnect grace time expired")
)

const websocketURL = "wss://eventsub.wss.twitch.tv/ws"

// Default timeout for keepalive messages in seconds.
const defaultKeepaliveTimeout = 600

// Websocket is an EventSub websocket client.
type Websocket struct {
	client       *http.Client
	eventTracker eventtracker.EventTracker

	serverURL          string
	serverReconnectURL string

	retryAttempts uint
	retryDelay    time.Duration

	conn            *websocket.Conn
	connDialOptions *websocket.DialOptions

	keepaliveSeconds uint
	lastKeepalive    atomic.Value

	isActive       atomic.Bool
	isReconnecting atomic.Bool

	restart     chan struct{}
	welcome     chan struct{}
	reconnected chan struct{}
	disconnect  chan struct{}

	onWelcome    func(WebsocketWelcomeMessage)
	onKeepalive  func(WebsocketKeepaliveMessage)
	onPing       func()
	onReconnect  func(WebsocketReconnectMessage)
	onError      func(error)
	onDisconnect func()

	callback[WebsocketNotificationMetadata]
}

func newWebsocket(eventTracker eventtracker.EventTracker, options ...WebsocketOption) *Websocket {
	ws := &Websocket{
		client:             http.DefaultClient,
		eventTracker:       eventTracker,
		serverURL:          websocketURL,
		serverReconnectURL: websocketURL,
		retryAttempts:      5,
		retryDelay:         1 * time.Second,
		keepaliveSeconds:   defaultKeepaliveTimeout,
		reconnected:        make(chan struct{}),
		welcome:            make(chan struct{}),
		restart:            make(chan struct{}),
	}

	for _, option := range options {
		option(ws)
	}

	ws.connDialOptions = &websocket.DialOptions{
		HTTPClient: ws.client,
		OnPingReceived: func(_ context.Context, _ []byte) bool {
			if ws.onPing != nil {
				go ws.onPing()
			}

			return true
		},
	}

	return ws
}

func (ws *Websocket) Connect(ctx context.Context) error {
	if !ws.setActive() {
		return nil
	}

	defer func() {
		ws.setInactivate()

		if ws.onDisconnect != nil {
			go ws.onDisconnect()
		}
	}()

	// We should recreate disconnect channel on each connect so user can do connect-disconnect cycle as many times as he
	// wants and don't get panic on double-close already closed channel.
	ws.disconnect = make(chan struct{})

	if err := ws.connect(ctx, ws.serverURL); err != nil {
		return err
	}

	// Context of the whole instance including side workers (e.g. keepalive worker).
	lifecycleCtx, stopLifecycle := context.WithCancel(context.Background())
	defer stopLifecycle()

	// Context of the connection itself that can be canceled to restart for example.
	connectionCtx, stopConnection := context.WithCancel(lifecycleCtx)

	stop := make(chan error)

	go ws.startKeepaliveWorker(lifecycleCtx)
	go func() {
		stop <- ws.startReadWorker(connectionCtx)
	}()

	for {
		select {
		case err := <-stop:
			if ws.isReconnecting.Load() {
				continue
			}
			stopConnection()
			return err
		case <-ws.restart:
			stopConnection()

			// Restore connection context to reuse it on next restart (reconnect).
			connectionCtx, stopConnection = context.WithCancel(lifecycleCtx)

			go func() {
				stop <- ws.startReadWorker(connectionCtx)
			}()
		case <-ws.disconnect:
			stopLifecycle()
		}
	}
}

func (ws *Websocket) Disconnect() error {
	if !ws.setInactivate() {
		return nil
	}

	ws.disconnect <- struct{}{}
	close(ws.disconnect)

	if err := ws.conn.Close(websocket.StatusNormalClosure, "client is shutting down"); err != nil {
		return fmt.Errorf("close connection: %w", err)
	}

	return nil
}

func (ws *Websocket) connect(ctx context.Context, serverURL string) error {
	if ws.retryAttempts > 0 {
		return retry.Do(
			func() error {
				return ws.connectWithoutRetry(ctx, serverURL)
			},
			retry.Attempts(ws.retryAttempts),
			retry.Delay(ws.retryDelay),
			retry.Context(ctx),
		)
	}

	return ws.connectWithoutRetry(ctx, serverURL)
}

func (ws *Websocket) connectWithoutRetry(ctx context.Context, serverURL string) error {
	rawURL, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("parse server url: %w", err)
	}

	keepaliveTimeoutSeconds := strconv.FormatUint(uint64(ws.keepaliveSeconds), 10)

	query := rawURL.Query()
	query.Add("keepalive_timeout_seconds", keepaliveTimeoutSeconds)

	rawURL.RawQuery = query.Encode()

	conn, _, err := websocket.Dial(ctx, rawURL.String(), ws.connDialOptions)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	ws.conn = conn
	return nil
}

// reconnect reconnects to the new eventsub edge server according to the reconnect message flow.
//
// Reference: https://dev.twitch.tv/docs/eventsub/handling-websocket-events/#reconnect-message.
func (ws *Websocket) reconnect(ctx context.Context, reconnectURL string) error {
	if !ws.setReconnecting() {
		return nil
	}

	defer func() {
		ws.isReconnecting.Store(false)

		select {
		case ws.reconnected <- struct{}{}:
		default:
		}
	}()

	// Save old connection to close it later after connecting to the new edge server and set new reconnection URL.
	oldConnection := ws.conn
	ws.serverReconnectURL = reconnectURL

	if err := ws.connect(ctx, reconnectURL); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	// Signal read worker to stop reading from old connection, restart with new one and wait for welcome message.
	ws.restart <- struct{}{}
	<-ws.welcome

	// Close old connection to complete reconnect flow.
	if err := oldConnection.Close(websocket.StatusNormalClosure, "connected to the new edge server"); err != nil {
		return fmt.Errorf("close old connection: %w", err)
	}

	return nil
}

// startReadWorker starts and blocks on reading and processing messages from the connection.
func (ws *Websocket) startReadWorker(ctx context.Context) error {
	for {
		_, data, err := ws.conn.Read(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			switch websocket.CloseStatus(err) {
			case websocket.StatusNormalClosure:
				if ws.isReconnecting.Load() {
					<-ws.reconnected
					continue
				}

				return nil
			case websocket.StatusCode(4003):
				return ErrConnectionUnused
			case websocket.StatusCode(4004):
				return ErrReconnectTimeout
			}

			return fmt.Errorf("read from connection: %w", err)
		}

		var message websocketRawMessage

		if err = json.Unmarshal(data, &message); err != nil {
			return fmt.Errorf("unmarshal data from connection: %w", err)
		}

		metadata := WebsocketNotificationMetadata{
			MessageId:           message.Metadata.MessageId,
			MessageType:         message.Metadata.MessageType,
			MessageTimestamp:    message.Metadata.MessageTimestamp,
			SubscriptionType:    message.Metadata.SubscriptionType,
			SubscriptionVersion: message.Metadata.SubscriptionVersion,
		}

		err, safe := isSafeMessage(
			ctx,
			ws.onDuplicate,
			ws.eventTracker,
			metadata,
			metadata.MessageId,
			metadata.MessageTimestamp,
		)
		if err != nil {
			return fmt.Errorf("is safe message: %w", err)
		}

		if !safe {
			continue
		}

		if err = ws.handleMessage(ctx, message); err != nil {
			return fmt.Errorf("handle message: %w", err)
		}
	}
}

func (ws *Websocket) startKeepaliveWorker(ctx context.Context) {
	const delay = 1 * time.Second

	keepalive := time.NewTicker(time.Second*time.Duration(ws.keepaliveSeconds) + delay)
	defer keepalive.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-keepalive.C:
			var (
				now, _               = timestampUTCFromString(time.Now().String())
				lastKeepaliveSeconds = now.Sub(ws.getLastKeepalive()).Seconds()
			)

			if lastKeepaliveSeconds < float64(ws.keepaliveSeconds) {
				continue
			}

			// We must reconnect to the eventsub server if keepalive timeout expired too.
			if err := ws.reconnect(ctx, ws.serverReconnectURL); err != nil {
				if ws.onError == nil {
					continue
				}

				go ws.onError(err)
			}
		}
	}
}

// setInactivate atomically sets Websocket in reconnecting state (meaning that it's reconnecting to the server) if it was
// not in reconnecting state and returns does Websocket was not reconnecting.
func (ws *Websocket) setReconnecting() bool {
	return ws.isReconnecting.CompareAndSwap(false, true)
}

// setInactivate atomically sets Websocket in active state (meaning that it's running and eventually will be connected to
// the server or already connected) if it was in inactive state and returns does Websocket was in inactive state.
func (ws *Websocket) setActive() bool {
	return ws.isActive.CompareAndSwap(false, true)
}

// setInactivate atomically sets Websocket in inactive state (meaning that it's not connected to the server nor running) if it
// was in active state and returns does Websocket was in active state.
func (ws *Websocket) setInactivate() bool {
	return ws.isActive.CompareAndSwap(true, false)
}

// setLastKeepalive atomically sets timestamp for the last time keepalive message was received from server.
func (ws *Websocket) setLastKeepalive(lastKeepalive TimestampUTC) {
	ws.lastKeepalive.Store(lastKeepalive)
}

// getLastKeepalive atomically loads and returns the last time keepalive message was received from the server.
// It may return zeroed time if there was no keepalive message since the application started, so you should consider this case.
func (ws *Websocket) getLastKeepalive() time.Time {
	lastKeepalive := ws.lastKeepalive.Load()
	if lastKeepalive == nil {
		return time.Time{}
	}

	return lastKeepalive.(TimestampUTC).Time
}
