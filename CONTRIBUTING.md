# Contributing to Tickle

Thanks for your interest in making Tickle better!

## Philosophy

Tickle is a minimal, focused library. The goal is to do one thing well: token-based pub/sub notifications with blocking waits. Keep this in mind when proposing changes.

## Getting Started

1. Fork and clone the repository
2. Go 1.25.4+ is required
3. Install dependencies: `go mod download`
4. Run tests: `go test ./...`

## Before You Open a PR

- Ensure all tests pass: `go test ./...`
- Keep changes focused and minimal
- Add tests for new behavior
- Update documentation if APIs change
- Run `go fmt ./...`

## Testing

Tests should be minimal and clear. Use table-driven tests for multiple scenarios. Every test should validate a specific behavior—avoid testing implementation details.

## Code Style

- Follow standard Go conventions
- Use meaningful variable names
- Keep functions small and focused
- Write self-documenting code over verbose comments

## What Won't Be Merged

- Features that expand scope beyond pub/sub notifications
- Complex abstractions or alternative APIs
- Significant performance optimizations that sacrifice clarity
- Dependencies (except testing libraries)

## Questions?

Open an issue to discuss before investing time in a large change. We want to keep Tickle simple, and early discussion prevents wasted effort.
