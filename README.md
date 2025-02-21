[![ci](https://github.com/fgrzl/tickle/actions/workflows/ci.yml/badge.svg)](https://github.com/fgrzl/tickle/actions/workflows/ci.yml)
[![Dependabot Updates](https://github.com/fgrzl/tickle/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/fgrzl/tickle/actions/workflows/dependabot/dependabot-updates)

# tickle

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
- Tickle subscribers based on tokens
- Wait for tickles with or without timeouts
- Dispose of subscriptions safely

## Installation

To install the library, use `go get`:

```sh
go get github.com/fgrzl/tickle
```

## Usage

Here is an example of how to use the library:

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
    ctx := context.Background()
    sub := sm.Add(ctx, "token1")

    go func() {
        if sub.Wait() {
            fmt.Println("Received notification for token1")
        }
    }()

    time.Sleep(1 * time.Second)
    sm.Tickle("token1")
}
```
