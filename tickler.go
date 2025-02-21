package tickle

import (
	"context"
	"sync"
	"time"
)

// NewTickler creates a new subscription manager.
func NewTickler() *Tickler {
	return &Tickler{}
}

type Tickler struct {
	subscriptions sync.Map // Key: token, Value: []*Subscription
}

// Subscribe creates a new subscription and associates it with one or more tokens.
func (sm *Tickler) Subscribe(ctx context.Context, tokens ...string) *Subscription {
	sub := NewSubscription(ctx, tokens...)

	for _, token := range tokens {
		existing, _ := sm.subscriptions.Load(token)
		var subs []*Subscription
		if existing != nil {
			if list, ok := existing.([]*Subscription); ok {
				subs = list
			}
		}
		subs = append(subs, sub)
		sm.subscriptions.Store(token, subs)
	}

	// Automatically unsubscribe when context is done
	go func() {
		<-ctx.Done()
		sm.Unsubscribe(sub)
	}()

	return sub
}

// Unsubscribe disposes of a subscription and removes it from all associated tokens.
func (sm *Tickler) Unsubscribe(sub *Subscription) {
	sub.Dispose() // Ensure it's disposed before removing.

	for _, token := range sub.tokens {
		value, ok := sm.subscriptions.Load(token)
		if !ok {
			continue // Avoid processing if not found
		}

		subs, ok := value.([]*Subscription)
		if !ok {
			continue // Prevent panic in case of incorrect type
		}

		// Remove the specific subscription.
		var updatedSubs []*Subscription
		for _, s := range subs {
			if s != sub {
				updatedSubs = append(updatedSubs, s)
			}
		}

		if len(updatedSubs) == 0 {
			sm.subscriptions.Delete(token) // Remove token entry if no subscriptions left.
		} else {
			sm.subscriptions.Store(token, updatedSubs)
		}
	}
}

// Tickle sends a notifies to all subscribers of the given tokens.
func (sm *Tickler) Tickle(tokens ...string) {
	for _, token := range tokens {
		value, ok := sm.subscriptions.Load(token)
		if !ok {
			continue
		}

		subs, ok := value.([]*Subscription)
		if !ok {
			continue
		}

		for _, sub := range subs {
			sub.Tickle()
		}
	}
}

// Subscription represents a single event listener.
type Subscription struct {
	ctx      context.Context
	cancel   context.CancelFunc
	tokens   []string
	ch       chan struct{}
	mu       sync.Mutex
	disposed bool
}

// NewSubscription creates a new subscription.
func NewSubscription(ctx context.Context, tokens ...string) *Subscription {
	ctx, cancel := context.WithCancel(ctx)
	return &Subscription{
		ctx:    ctx,
		cancel: cancel,
		tokens: tokens,
		ch:     make(chan struct{}, 1),
	}
}

// Tickle signals the subscription if it is still active.
func (s *Subscription) Tickle() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return // Prevent sending to a closed channel
	}

	select {
	case s.ch <- struct{}{}:
	default:
	}
}

// Wait blocks until a notification is received or the context is canceled.
func (s *Subscription) Wait() bool {
	select {
	case <-s.ctx.Done():
		return false
	case <-s.ch:
		return true
	}
}

// WaitTimeout waits for a signal until the timeout expires.
// Returns true if a signal was received, false if the subscription is canceled or timed out.
func (s *Subscription) WaitTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-s.ctx.Done():
		return false
	case <-s.ch:
		return !s.disposed
	case <-timer.C:
		return false
	}
}

func (s *Subscription) Dispose() {
	s.mu.Lock()
	if s.disposed {
		s.mu.Unlock()
		return
	}
	s.disposed = true

	// Drain the channel to remove any lingering notifications
	for len(s.ch) > 0 {
		<-s.ch
	}

	close(s.ch) // Safe since only Dispose() closes it.
	s.cancel()
	s.mu.Unlock()
}
