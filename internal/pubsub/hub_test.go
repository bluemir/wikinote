package pubsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/bluemir/wikinote/internal/pubsub"
)

type key struct{}

var tkey key = struct{}{}

func testContext(t *testing.T, timeout time.Duration) (context.Context, func()) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)

	ctx := context.Background()
	ctx = context.WithValue(ctx, tkey, t)
	return context.WithTimeout(ctx, timeout)
}
func From(ctx context.Context) *testing.T {
	return ctx.Value(tkey).(*testing.T)
}

func TestSendEvent(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	h := &CounterHandler{}

	hub.AddHandler("*", h)

	hub.Publish("test", nil)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 1, h.GetCount())
}
func TestSendMultiple(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	recoder := NewRecoder(ctx, hub)

	hub.Publish("test-1", nil)
	hub.Publish("test-2", nil)
	hub.Publish("test-3", nil)
	hub.Publish("test-4", nil)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 4, len(recoder.recodes))
	assert.Equal(t, []string{
		"test-1",
		"test-2",
		"test-3",
		"test-4",
	}, recoder.History())
}

func TestAddEventHandlerWithNull(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	hub.AddHandler("test", pubsub.Handler(nil))

	hub.Publish("test", nil)
}

func TestEventKind(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	counter := &CounterHandler{}
	hub.AddHandler("test-1", counter)

	hub.Publish("test-2", nil)

	assert.Equal(t, 0, counter.GetCount())
}

func TestFireEventInsideEventHandler(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	recoder := NewRecoder(ctx, hub)

	hub.AddHandler("button.down", FowardHandler{"click"})

	hub.Publish("button.down", nil)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, []string{
		"button.down",
		"click",
	}, recoder.History())
}

func TestAddHandlerInsideEventHandler(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	recoder := NewRecoder(ctx, hub)
	hub.AddHandler("do", ReplaceSelfHandler{})
	hub.Publish("do", nil)
	hub.Publish("do", nil)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, []string{
		"do",
		"do",
		"done",
	}, recoder.History())
}

func TestListenWithStar(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	counter := &CounterHandler{}
	hub.AddHandler("*", counter)

	hub.Publish("test", nil)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 1, counter.GetCount())
}
func ignoreTestListenWithStarInWord(t *testing.T) {
	ctx, cancel := testContext(t, 1*time.Second)
	defer cancel()

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	counter := &CounterHandler{}
	hub.AddHandler("test.*", counter)

	hub.Publish("test.test", nil)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 1, counter.GetCount())
}
