# Codebase Overview

This file provides instructions and overview for agentic coding agents working in this repository.

## Project Overview

`github.com/kvizyx/twitchy` is a pure Go library (requires at least Go 1.24) for the Twitch EventSub and Helix API.
It supports both Webhook and WebSocket transports. There are no non-Go files to build or compile.

```
twitchy/
├── eventsub/                   # Main public package
│   ├── messagetracker/         # Duplicate-message tracking (in-memory and Redis)
│   ├── eventsub.go             # Top-level EventSub struct and New()
│   ├── option.go               # Functional options
│   ├── callback.go             # Generic callback store
│   ├── event.go                # All event payload structs
│   ├── event_type.go           # EventType string constants
│   ├── condition.go            # All subscription condition types
│   ├── object.go               # Shared sub-objects and enums
│   ├── subscription.go         # Generic Subscription[C, T] type
│   ├── webhook.go / websocket.go
│   └── ...
├── helix/                      # Future Helix API package (stubs only for now)
├── internal/
│   ├── json/                   # Swappable JSON marshal/unmarshal vars
│   └── shardedset/             # Generic concurrent sharded set with TTL eviction
├── _examples/                  # Standalone main packages (not part of the API)
├── Taskfile.yaml
├── .golangci.yaml
└── .editorconfig
```

---

## Build, Lint, and Test Commands

### Task runner

This project uses [Task](https://taskfile.dev) instead of `make`. There is no `Makefile`.

```sh
# Run all tests with coverage report
task test
```

This is equivalent to:
```sh
mkdir -p .coverage
go test -coverprofile=.coverage/coverage.out ./...
go tool cover -html=.coverage/coverage.out -o ./.coverage/coverage.html
```

### Running tests manually

```sh
# All tests
go test ./...

# All tests with race detector (matches CI)
go test -race -v ./...

# Single package
go test github.com/kvizyx/twitchy/eventsub

# Single test function (in any matching package)
go test ./eventsub/... -run TestFunctionName

# Single test function in a specific package
go test github.com/kvizyx/twitchy/eventsub -run TestFunctionName

# With verbose output and race detector
go test -race -v ./eventsub/... -run TestFunctionName
```

### Linting

```sh
# Run golangci-lint (uses .golangci.yaml automatically)
golangci-lint run ./...

# Run on a specific package
golangci-lint run ./eventsub/...
```

Config is in `.golangci.yaml` with `disable-all: true` (opt-in model). Key constraints:
- Max line length: **120 characters**
- Max function length: **100 lines / 50 statements**
- Cyclomatic complexity: **30**
- `gofumpt` formatting (stricter than `gofmt`)
- No `init()` functions (`gochecknoinits`) — the one exception is `internal/json/json.go`
- Comments must end with a period (`godot`)

### Formatting

```sh
gofumpt -w .
goimports -w .
```

---

## Code Style Guidelines

### Imports

Use two groups separated by a blank line: stdlib first, then all external and internal packages together.
Internal packages (`github.com/kvizyx/twitchy/...`) are **not** placed in a separate third group.

```go
import (
    "context"
    "errors"
    "fmt"
    "net/http"

    "github.com/avast/retry-go/v4"
    "github.com/coder/websocket"
    "github.com/kvizyx/twitchy/eventsub/messagetracker"
    "github.com/kvizyx/twitchy/internal/json"
)
```

### Naming Conventions

| Entity | Convention | Example |
|---|---|---|
| Files | `snake_case` with subject prefix | `websocket_message_handler.go` |
| Exported types | `PascalCase` with full domain context | `ChannelPointsCustomRewardRedemptionAddEvent` |
| Unexported types | `camelCase` | `callbackStore` |
| Interfaces | `PascalCase`, semantic names | `MessageTracker`, `Transport` |
| Constructors | `New()` or `New<Type>()` | `NewInMemoryMessageTracker()`, `newWebsocket()` |
| Exported functions/methods | `PascalCase` | `ServeHTTP`, `Connect` |
| Unexported functions/methods | `camelCase` | `handleMessage`, `isSafeMessage` |
| Constants/string enums | `TypeName` + `ValueName` | `EventTypeChannelFollow`, `SubscriptionTierOne` |
| Sentinel errors | `Err` prefix | `ErrConnectionUnused`, `ErrInvalidWebhookSecret` |
| Functional option types | `Option func(*T)` | `Option`, `WebsocketOption` |
| Functional option constructors | `With<Feature>()` | `WithUnmarshal()`, `WebsocketWithConnectTimeout()` |
| Generic type parameters | Short and descriptive | `[K comparable, V any]`, `[C Condition]` |

### Types and Generics

- Prefer **concrete types** over interfaces; use interfaces only at genuine polymorphism points.
- Use **marker interfaces** (unexported methods) to constrain generic type parameters to domain types:
  ```go
  type Condition interface{ condition() }
  type Transport interface{ transport() }
  ```
- Add compile-time interface assertions where a type must satisfy an external interface:
  ```go
  var _ http.Handler = (*Webhook)(nil)
  var _ MessageTracker = (*InMemoryMessageTracker)(nil)
  ```
- Use the **functional options** pattern for all optional configuration (`Option func(*T)`).
- When options are implementation-specific, prefix them with the implementation name:
  `InMemoryOption`, `InMemoryWithMessageTTL()`, `RedisOption`, `RedisWithKeyBuilder()`.

### Error Handling

- Wrap errors with context using `fmt.Errorf("description: %w", err)`:
  ```go
  return fmt.Errorf("parse server url: %w", err)
  ```
- Declare sentinel errors at package level with `errors.New(...)` and the `Err` prefix.
- Use `errors.Is(err, Target)` for checking sentinels — never string comparison.
- Dispatch goroutine errors to the user-registered `onError` callback; do not silently drop them.
- In `http.Handler` implementations, convert errors to `http.Error(w, msg, statusCode)`.
- Explicitly ignored errors use the blank identifier with a comment:
  ```go
  _, _ = hasher.Write(data) // hash.Hash.Write never returns an error
  _ = r.Body.Close()
  ```
- The unconventional `return error, bool` signature is used in a few internal helpers for
  combined error+gate semantics — do not copy this pattern in new code without strong justification.

### Comments and Documentation

- All exported symbols **must** have a doc comment.
- Doc comments **must end with a period** (enforced by `godot` linter).
- Begin doc comments with the identifier name:
  ```go
  // EventSub is a Twitch event-sub module.
  // Webhook returns a new event-sub Webhook handler...
  ```
- Include Twitch API reference links where relevant:
  ```go
  // Reference: https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/#channelfollow
  ```
- Use `// TODO: implement me` for stubs that are intentionally unfinished.
- Unexported types and functions should have comments where the logic is non-obvious,
  especially for concurrency decisions.

### Package and File Organization

- **One API domain = one package.** Keep `eventsub`, `helix`, etc. separate.
- Split large packages across multiple focused files by concept, not by type:
  `option.go`, `callback.go`, `websocket.go`, `websocket_option.go`, `websocket_message.go`, etc.
- When a concept has multiple implementations with their own options, split options into separate
  files per implementation (e.g., `message_tracker_in_memory_options.go`,
  `message_tracker_redis_options.go`).
- Place reusable utilities that must not be exported under `internal/`.
- Create a sub-package (like `eventsub/messagetracker`) only when there is a standalone concept
  with its own interface and multiple concrete implementations.
- `_examples/` are standalone `main` packages and are not part of the importable API.

### Concurrency

- Use `sync/atomic.Bool` and `atomic.Value` for hot-path state (e.g., `isActive`, `isReconnecting`).
- Use channels for goroutine coordination and signaling rather than shared mutable state with mutexes
  where practical.
- Prefer explicit `Stop()` methods over `context.Context` for lifecycle management of background
  goroutines in internal utilities (e.g., `ShardedSet.Stop()`).
- Every user-facing callback is dispatched in a new goroutine so user code cannot block the
  internal read loop:
  ```go
  if handler := cb.onChannelFollow; handler != nil {
      go handler(event, metadata)
  }
  ```

### JSON

- Do not import `encoding/json` directly. Use `github.com/kvizyx/twitchy/internal/json` instead.
  This provides swappable `Unmarshal`/`Marshal` function variables that users can override
  (e.g., to use `sonic` or `jsoniter`) via `WithUnmarshal()`.

### General Rules

- Do not introduce `init()` functions. The only allowed exception is `internal/json/json.go`.
- Do not use `github.com/pkg/errors` or `golang.org/x/net/context` (both are excluded by IDE config).
- Keep functions under **100 lines / 50 statements** and cyclomatic complexity under **30**.
- All lines must be at most **120 characters** wide.
- Use tabs for indentation (size 4), spaces for YAML/JSON (size 2), per `.editorconfig`.
- Spell error messages and comments in US English (`misspell` linter is active).
