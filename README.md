# go-logger

`go-logger` is a small structured logger for Go projects. It writes log records with a level, message, and optional fields, and lets you configure both the output writer and encoder.

> [!WARNING]
> Versions before `v1.0.0` do not guarantee backward compatibility between minor releases. For example, upgrading from `v0.2.x` to `v0.3.x` may include breaking API changes.

## Installation

```bash
go get github.com/Palladium-blockchain/go-logger
```

## Standard Usage

The `logger` package provides a default logger through package-level functions:

```go
package main

import (
	"errors"

	log "github.com/Palladium-blockchain/go-logger/pkg/logger"
)

func main() {
	log.Debug("loading node config")

	log.Info(
		"node started",
		log.WithField("height", 42),
		log.WithField("role", "validator"),
	)

	err := errors.New("connection refused")
	log.Error("failed to connect to peer", log.WithError(err))
}
```

By default, the logger uses the plain encoder and writes to `os.Stdout`.

Example output:

```text
debug | loading node config
info | node started height: 42  role: validator
error | failed to connect to peer error: connection refused
```

## Creating a Logger

Use `logger.New` when you need a separate logger instance with custom settings:

```go
package main

import (
	"os"

	"github.com/Palladium-blockchain/go-logger/pkg/encoder/plain"
	"github.com/Palladium-blockchain/go-logger/pkg/logger"
)

func main() {
	log := logger.New(
		logger.WithEncoder(plain.NewEncoder()),
		logger.WithLockingWriter(os.Stdout),
	)

	log.Info("service started", logger.WithField("port", 8080))
}
```

Available options:

| Option | Description |
| --- | --- |
| `logger.WithEncoder(encoder)` | Sets the encoder used to write log records. |
| `logger.WithWriter(writer)` | Sets the writer directly. Use this when the writer already handles synchronization. |
| `logger.WithLockingWriter(writer)` | Wraps an `io.Writer` with a mutex-protected writer for concurrent logging. |
| `logger.WithErrorHandler(handlerFn)` | Sets a callback for encoder or writer errors. |
| `logger.WithLogLevel(level)` | Sets the minimum level to write. Lower-priority records are skipped. |

## Fields

Fields add structured context to a log record:

```go
log.Info(
	"block imported",
	logger.WithField("height", 1201),
	logger.WithField("hash", "0xabc123"),
	logger.WithField("finalized", true),
)
```

Use `logger.WithError(err)` for errors. It stores the error under the `error` key:

```go
if err != nil {
	log.Error("failed to import block", logger.WithError(err))
}
```

## Encoders

The repository includes two encoder packages:

| Encoder | Package | Purpose |
| --- | --- | --- |
| Plain encoder | `pkg/encoder/plain` | The default encoder for `logger.New`. Writes readable text lines. |
| JSON encoder | `pkg/encoder/json` | Writes log records as newline-delimited JSON. |

### Plain Encoder

The plain encoder is useful for console output:

```go
log := logger.New(
	logger.WithEncoder(plain.NewEncoder()),
)

log.Warn("peer is slow", logger.WithField("peer_id", "12D3KooW..."))
```

Example output:

```text
warning | peer is slow peer_id: 12D3KooW...
```

### JSON Encoder

Use the JSON encoder when you want structured logs:

```go
package main

import (
	"os"

	logjson "github.com/Palladium-blockchain/go-logger/pkg/encoder/json"
	"github.com/Palladium-blockchain/go-logger/pkg/logger"
)

func main() {
	log := logger.New(
		logger.WithEncoder(&logjson.Encoder{}),
		logger.WithLockingWriter(os.Stdout),
	)

	log.Info(
		"node started",
		logger.WithField("height", 42),
		logger.WithField("role", "validator"),
	)
}
```

Example output:

```json
{"height":42,"level":"info","message":"node started","role":"validator"}
```

The JSON encoder writes fields in a flat object. User fields named `level` or `message` are prefixed as `field.level` and `field.message` to avoid overwriting the built-in fields.

## Custom Encoder

A custom encoder must implement `core.Encoder`:

```go
package main

import (
	"fmt"
	"io"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
	"github.com/Palladium-blockchain/go-logger/pkg/logger"
)

type CompactEncoder struct{}

func (CompactEncoder) Encode(w io.Writer, record core.Record) error {
	_, err := fmt.Fprintf(w, "[%s] %s\n", record.Level, record.Message)
	return err
}

func main() {
	log := logger.New(logger.WithEncoder(CompactEncoder{}))
	log.Info("ready")
}
```

## Log Levels

The logger supports four levels:

| Method | Level value |
| --- | --- |
| `Debug` | `debug` |
| `Info` | `info` |
| `Warn` | `warning` |
| `Error` | `error` |

Use `logger.WithLogLevel` to skip records below a minimum level:

```go
log := logger.New(logger.WithLogLevel(core.Warn))

log.Info("skipped")
log.Warn("written")
log.Error("written")
```
