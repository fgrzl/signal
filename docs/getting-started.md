# Getting started

## Install

```bash
go get github.com/fgrzl/tickle
```

## Basic pattern

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/fgrzl/tickle"
)

func main() {
    sm := tickle.NewTickler()
    sub := sm.Subscribe(context.Background(), "job-done")

    go func() {
        if sub.Wait() {
            fmt.Println("notified")
        }
    }()

    time.Sleep(100 * time.Millisecond)
    sm.Tickle("job-done")
}
```

## Wait with timeout

```go
if sub.WaitTimeout(5 * time.Second) {
    // received
} else {
    // timed out
}
```

## Wait with context

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
if sub.WaitContext(ctx) {
    // notification received
} else {
    // context canceled, disposed, or timed out without a tickle
}
```

## Cleanup

```go
sub.Dispose()
sm.Unsubscribe(sub)
```

## Tests

```bash
go test ./...
```

See [CONTRIBUTING](../CONTRIBUTING.md) for PR expectations.
