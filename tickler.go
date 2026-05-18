package tickle

import (
	"context"
	"sync"
	"time"
)

// Tickler manages subscriptions and dispatches notifications by token.
type Tickler struct {
	mu   sync.Mutex
	subs map[string][]*Subscription
}

// NewTickler creates a new Tickler.
func NewTickler() *Tickler {
	return &Tickler{
		subs: make(map[string][]*Subscription),
	}
}

// Subscribe creates a subscription for the given tokens.
// The subscription is automatically removed when its context is canceled.
func (t *Tickler) Subscribe(ctx context.Context, tokens ...string) *Subscription {
	sub := newSubscription(ctx, tokens...)

	t.mu.Lock()
	for _, token := range tokens {
		t.subs[token] = append(t.subs[token], sub)
	}
	t.mu.Unlock()

	go func() {
		<-sub.Done()
		t.remove(sub)
	}()

	return sub
}

// Unsubscribe disposes a subscription and removes it from all associated tokens.
func (t *Tickler) Unsubscribe(sub *Subscription) {
	sub.Dispose()
}

// remove deletes the subscription from internal tracking.
func (t *Tickler) remove(sub *Subscription) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, token := range sub.tokens {
		subs := t.subs[token]
		for i, s := range subs {
			if s == sub {
				subs[i] = subs[len(subs)-1]
				subs[len(subs)-1] = nil
				subs = subs[:len(subs)-1]
				break
			}
		}
		if len(subs) == 0 {
			delete(t.subs, token)
		} else {
			t.subs[token] = subs
		}
	}
}

// Tickle notifies all subscribers of the given tokens.
func (t *Tickler) Tickle(tokens ...string) {
	t.mu.Lock()
	var targets []*Subscription
	for _, token := range tokens {
		targets = append(targets, t.subs[token]...)
	}
	t.mu.Unlock()

	for _, sub := range targets {
		sub.Tickle()
	}
}

// Subscription represents a notification listener bound to one or more tokens.
type Subscription struct {
	ctx    context.Context
	done   chan struct{}
	tokens []string
	ch     chan struct{}
	once   sync.Once
}

func newSubscription(ctx context.Context, tokens ...string) *Subscription {
	if ctx == nil {
		ctx = context.Background()
	}

	s := &Subscription{
		ctx:    ctx,
		done:   make(chan struct{}),
		tokens: tokens,
		ch:     make(chan struct{}, 1),
	}

	go func() {
		select {
		case <-ctx.Done():
			s.Dispose()
		case <-s.done:
		}
	}()

	return s
}

func (s *Subscription) disposed() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}

func (s *Subscription) tryConsumeTickle() bool {
	select {
	case <-s.ch:
		return true
	default:
		return false
	}
}

// Tickle signals the subscription. No-op if disposed.
func (s *Subscription) Tickle() {
	if s.disposed() {
		return
	}
	select {
	case s.ch <- struct{}{}:
	default:
	}
}

// Wait blocks until a notification is received or the subscription is disposed.
func (s *Subscription) Wait() bool {
	if s.tryConsumeTickle() {
		return true
	}

	select {
	case <-s.ch:
		return true
	case <-s.done:
		return s.tryConsumeTickle()
	}
}

// WaitContext blocks until a notification is received, the context is canceled, or the subscription is disposed.
func (s *Subscription) WaitContext(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	if s.tryConsumeTickle() {
		return true
	}

	select {
	case <-s.ch:
		return true
	case <-s.done:
		return s.tryConsumeTickle()
	case <-ctx.Done():
		return s.tryConsumeTickle()
	}
}

// WaitTimeout blocks until a notification, timeout, or disposal.
func (s *Subscription) WaitTimeout(timeout time.Duration) bool {
	if s.tryConsumeTickle() {
		return true
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-s.ch:
		return true
	case <-s.done:
		return s.tryConsumeTickle()
	case <-timer.C:
		return s.tryConsumeTickle()
	}
}

// Done returns a channel that is closed when the subscription is disposed.
func (s *Subscription) Done() <-chan struct{} {
	return s.done
}

// Dispose closes the subscription. Safe to call multiple times.
func (s *Subscription) Dispose() {
	s.once.Do(func() {
		close(s.done)
	})
}
