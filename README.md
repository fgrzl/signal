[![ci](https://github.com/fgrzl/tickle/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/tickle/actions/workflows/ci.yml)
[![Dependabot Updates](https://github.com/fgrzl/tickle/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/tickle/actions/workflows/dependabot/dependabot-updates)

# Tickle

In-process, token-based pub/sub notifications with blocking waits.

## Features

- Subscribe to one or more string tokens
- `Tickle` subscribers from any goroutine
- `Wait`, `WaitTimeout`, and `WaitContext` on subscriptions
- Safe disposal and unsubscribe
- `Subscribe(nil, ...)` uses `context.Background()`
- Buffered notifications win over dispose for wait calls

## Installation

```bash
go get github.com/fgrzl/tickle
```

## Quick example

```go
sm := tickle.NewTickler()
sub := sm.Subscribe(context.Background(), "token1")

go func() {
    if sub.Wait() {
        fmt.Println("Received notification for token1")
    }
}()

time.Sleep(time.Second)
sm.Tickle("token1")
```

## Documentation

Full guides: **[docs/](docs/README.md)**

- [Overview](docs/overview.md)
- [Getting started](docs/getting-started.md)

## Running tests

```bash
go test ./...
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

See repository license file.
