# Overview

Tickle provides **cooperative notification** inside a single process: subscribers register interest in string **tokens**, publishers **tickle** those tokens, and subscribers **wait** until a notification arrives.

## Core types

| Type | Role |
|------|------|
| `Tickler` | Registry of subscriptions; call `Tickle` to notify |
| `Subscription` | Per-subscriber handle; `Wait`, `WaitTimeout`, `WaitContext` |

## Guarantees

- **Buffered delivery** — once a notification is buffered, wait calls return promptly even if the subscription is disposed afterward
- **Token fan-out** — one `Tickle("a", "b")` notifies every subscription that includes either token
- **Context** — `Subscribe(nil, ...)` uses `context.Background()`; canceling the subscribe context unblocks `Wait`, `WaitTimeout`, and `WaitContext` (via `Done()`). `WaitContext` also respects the caller's context.
- **Thread-safe** — safe concurrent subscribe, tickle, and wait from multiple goroutines

## What Tickle is not

- Not distributed messaging — no network, persistence, or cross-process delivery
- Not a replacement for channels when you only need one producer/consumer pair
- Not fair queuing across subscribers — each subscription has its own buffer

## Typical use cases

- **Processor scheduling** — wake workers when work arrives on a logical token
- **Test synchronization** — block until an async step signals completion
- **In-process event coalescing** — tickle once, many waiters resume
