[![ci](https://github.com/fgrzl/signal/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/signal/actions/workflows/ci.yml)
[![Dependabot Updates](https://github.com/fgrzl/signal/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/signal/actions/workflows/dependabot/dependabot-updates)

# Signal

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Running Tests](#running-tests)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- Add and remove subscriptions
- Notify subscribers based on tokens
- Wait for notifications with or without timeouts
- Dispose of subscriptions safely

## Installation

To install the library, use `go get`:

```sh
go get github.com/fgrzl/signal
```

## Usage

Here is an example of how to use the library:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/fgrzl/signal"
)

func main() {
    sm := signal.NewSubscriptionManager()
    ctx := context.Background()
    sub := sm.Add(ctx, "token1")

    go func() {
        if sub.Wait() {
            fmt.Println("Received notification for token1")
        }
    }()

    time.Sleep(1 * time.Second)
    sm.Notify("token1")
}
```
