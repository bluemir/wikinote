package pubsub_test

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bluemir/wikinote/internal/pubsub"
)

// --- Test Setup & Helpers ---
func testContext(t *testing.T, timeout time.Duration) (context.Context, func()) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "test", t)
	// logrus.SetLevel(logrus.TraceLevel)
	// logrus.SetReportCaller(true)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	return ctx, cancel
}

// Define some event types for testing
type EventForTest struct {
	Message string
}

// Simple recording handler
type RecordingHandler struct {
	mu     sync.Mutex
	Events []pubsub.Event
}

func NewRecordingHandler() *RecordingHandler {
	return &RecordingHandler{
		Events: make([]pubsub.Event, 0),
	}
}

func (rh *RecordingHandler) Handle(ctx context.Context, evt pubsub.Event) {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	// evt.Context = nil // Avoid storing context if comparing events directly
	rh.Events = append(rh.Events, evt)
}

func (rh *RecordingHandler) Count() int {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	return len(rh.Events)
}

func (rh *RecordingHandler) Kinds() []string {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	kinds := make([]string, len(rh.Events))
	for i, e := range rh.Events {
		kinds[i] = e.Kind // Kind is already string in v2 Event
	}
	return kinds
}

// --- Test Cases ---

func TestHub_SingleEvent(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	require.NoError(t, err)

	recorder := NewRecordingHandler()
	// Note: AddHandler uses reflect.TypeOf(kind), so pass an instance
	hub.AddHandler(EventForTest{}, recorder)

	eventToSend := EventForTest{Message: "hello"}
	hub.Publish(ctx, eventToSend)

	// Wait briefly for event processing (better: use Watch or sync)
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, recorder.Count(), "Handler should have received one event")
	require.Len(t, recorder.Events, 1)
	assert.Equal(t, reflect.TypeOf(eventToSend).String(), recorder.Events[0].Kind)
	assert.Equal(t, eventToSend, recorder.Events[0].Detail)
}

func TestHub_MultipleEvents(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	require.NoError(t, err)

	recoder := NewRecordingHandler()

	hub.AddHandler(EventForTest{}, recoder)

	hub.Publish(context.Background(), EventForTest{Message: "hello"})
	hub.Publish(context.Background(), EventForTest{Message: "world"})

	assert.Len(t, recoder.Events, 2)
	assert.Equal(t, reflect.TypeOf(EventForTest{}).String(), recoder.Events[0].Kind)
	assert.Equal(t, reflect.TypeOf(EventForTest{}).String(), recoder.Events[1].Kind)
	assert.Equal(t, "hello", recoder.Events[0].Detail.(EventForTest).Message)
	assert.Equal(t, "world", recoder.Events[1].Detail.(EventForTest).Message)
}

func TestHub_ConcurrentPublish(t *testing.T) {
	ctx, cancel := testContext(t, 5*time.Second) // Longer timeout for concurrency
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	require.NoError(t, err)

	recoder := NewRecordingHandler()

	hub.AddHandler(EventForTest{}, recoder)

	numGoroutines := 10
	numEventsPerG := 5
	totalEvents := numGoroutines * numEventsPerG

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for j := 0; j < numEventsPerG; j++ {
				hub.Publish(ctx, EventForTest{Message: fmt.Sprintf("g%d-e%d", gID, j)})
				// time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond) // Optional: add jitter
			}
		}(i)
	}

	wg.Wait() // Wait for all publishers to finish

	assert.Equal(t, totalEvents, recoder.Count(), "Should receive all events published concurrently")
}

// Handler that adds/removes other handlers
type ModifyingHandler struct {
	recoder *RecordingHandler
}

func (h ModifyingHandler) Handle(ctx context.Context, evt pubsub.Event) {
	hub := pubsub.HubFrom(ctx)

	hub.AddHandler(EventForTest{}, h.recoder)
}

type EventAddHandler struct {
}

func TestHub_ModifyHandlersInHandler(t *testing.T) {
	ctx, cancel := testContext(t, 2*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	require.NoError(t, err)

	recoder := NewRecordingHandler()

	hub.AddHandler(EventAddHandler{}, ModifyingHandler{recoder}) //

	hub.Publish(ctx, EventForTest{})

	assert.Equal(t, 0, recoder.Count(), "AddedRecorder count mismatch")

	hub.Publish(ctx, EventAddHandler{}) // trigger add handler

	hub.Publish(ctx, EventForTest{})

	assert.Equal(t, 1, recoder.Count(), "AddedRecorder count mismatch")
}
