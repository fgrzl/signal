package tickle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/fgrzl/tickle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test basic subscription addition and notification.
func TestSubscriptionManager_Tickle(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub := sm.Subscribe(ctx, "token1")

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = sub.Wait()
	}()

	// Act
	sm.Tickle("token1")

	// Assert
	wg.Wait()
	assert.True(t, result, "Wait() should return true after notification")
}

// Test waiting with a timeout.
func TestSubscription_WaitTimeout(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub := sm.Subscribe(ctx, "token1")

	// Act & Assert (timeout case)
	require.False(t, sub.WaitTimeout(10*time.Millisecond), "WaitTimeout should return false when it times out")

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = sub.WaitTimeout(100 * time.Millisecond)
	}()

	// Act
	sm.Tickle("token1")

	// Assert
	wg.Wait()
	assert.True(t, result, "WaitTimeout should return true after notification")
}

// Test removing a subscription.
func TestSubscriptionManager_Remove(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub := sm.Subscribe(ctx, "token1")

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = sub.WaitTimeout(1000 * time.Millisecond)
	}()

	// Act
	sm.Unsubscribe(sub)
	sm.Tickle("token1")

	// Assert
	wg.Wait()
	assert.False(t, result, "Subscription should not receive a notification after being removed")
}

// Test disposing a subscription.
func TestSubscription_Dispose(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub := sm.Subscribe(ctx, "token1")

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = sub.WaitTimeout(10000 * time.Millisecond)
	}()

	// Act
	sub.Dispose()

	// Assert
	wg.Wait()
	assert.False(t, result, "WaitTimeout should return false after Dispose()")
	assert.NotPanics(t, func() { sub.Tickle() }, "Tickle() should not panic after Dispose()")
}

// Test that multiple subscribers can wait for different tokens.
func TestSubscriptionManager_MultipleSubscriptions(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub1 := sm.Subscribe(ctx, "token1")
	sub2 := sm.Subscribe(ctx, "token2")

	var result1, result2 bool
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		result1 = sub1.Wait()
	}()
	go func() {
		defer wg.Done()
		result2 = sub2.WaitTimeout(10 * time.Millisecond)
	}()

	// Act
	sm.Tickle("token1")

	// Assert
	wg.Wait()
	assert.True(t, result1, "sub1 should be notified")
	assert.False(t, result2, "sub2 should not be notified")
}

// Test that multiple subscriptions can receive the same notification.
func TestSubscriptionManager_MultipleSubscribersSameToken(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub1 := sm.Subscribe(ctx, "token1")
	sub2 := sm.Subscribe(ctx, "token1") // Both subscriptions listen to the same token

	var result1, result2 bool
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		result1 = sub1.Wait()
	}()
	go func() {
		defer wg.Done()
		result2 = sub2.Wait()
	}()

	// Act
	sm.Tickle("token1")

	// Assert
	wg.Wait()
	assert.True(t, result1, "sub1 should be notified")
	assert.True(t, result2, "sub2 should be notified")
}

// Test that Tickleing a non-existent token does nothing.
func TestSubscriptionManager_TickleNonExistentToken(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	ctx := context.Background()
	sub := sm.Subscribe(ctx, "token1")

	var result bool
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		result = sub.WaitTimeout(100 * time.Millisecond)
	}()

	// Act
	sm.Tickle("token2")

	// Assert
	wg.Wait()
	assert.False(t, result, "Subscription should not receive a notification for an unrelated token")
}

func TestSubscription_SubscribeNilContext(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	var nilCtx context.Context

	var sub *tickle.Subscription
	assert.NotPanics(t, func() {
		sub = sm.Subscribe(nilCtx, "token1")
	}, "Subscribe(nil, ...) should not panic")
	require.NotNil(t, sub)

	// Act
	sm.Tickle("token1")

	// Assert
	assert.True(t, sub.WaitTimeout(100*time.Millisecond), "Subscription created with nil context should receive notifications")
}

func TestSubscription_WaitContext_CancelledContext(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	sub := sm.Subscribe(context.Background(), "token1")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act / Assert
	assert.False(t, sub.WaitContext(ctx), "WaitContext should return false when caller context is canceled first")
}

func TestSubscription_WaitContext_DisposedSubscription(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	sub := sm.Subscribe(context.Background(), "token1")
	sub.Dispose()

	// Act / Assert
	assert.False(t, sub.WaitContext(context.Background()), "WaitContext should return false when the subscription is disposed first")
}

func TestSubscription_WaitContext_NilContext(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	sub := sm.Subscribe(context.Background(), "token1")
	var nilCtx context.Context

	// Act
	sm.Tickle("token1")

	// Assert
	assert.True(t, sub.WaitContext(nilCtx), "WaitContext(nil) should be treated as context.Background()")
}

func TestSubscription_NotificationWinsWaitWhenBufferedAndDisposed(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	sub := sm.Subscribe(context.Background(), "token1")
	sm.Tickle("token1")
	sub.Dispose()

	// Act / Assert
	assert.True(t, sub.Wait(), "Wait should return true when a notification is already buffered")
}

func TestSubscription_NotificationWinsWaitTimeoutWhenBufferedAndDisposed(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	sub := sm.Subscribe(context.Background(), "token1")
	sm.Tickle("token1")
	sub.Dispose()

	// Act / Assert
	assert.True(t, sub.WaitTimeout(1*time.Millisecond), "WaitTimeout should return true when a notification is already buffered")
}

func TestSubscription_NotificationWinsWaitContextWhenBufferedAndDisposed(t *testing.T) {
	// Arrange
	sm := tickle.NewTickler()
	sub := sm.Subscribe(context.Background(), "token1")
	ctx, cancel := context.WithCancel(context.Background())
	sm.Tickle("token1")
	sub.Dispose()
	cancel()

	// Act / Assert
	assert.True(t, sub.WaitContext(ctx), "WaitContext should return true when a notification is already buffered")
}
