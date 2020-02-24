package event

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/stretchr/testify/assert"
)

const WatchDirStartData = "dir:/foo/bar"

var eventCreator = New(Config{})
var totalEventCounter uint64
var wgDone sync.WaitGroup
var wgSubscriberReady sync.WaitGroup

func TestSubscriptionpubSub(t *testing.T) {
	wgDone.Add(2)
	wgSubscriberReady.Add(2)

	go subscribeTester(t)
	go subscribeTester(t)

	wgSubscriberReady.Wait()

	err := eventCreator.Publish(context.TODO(), &core.EventData{
		Type: core.WatchDirStart,
		Data: []byte(WatchDirStartData),
	})
	assert.NoError(t, err)

	wgDone.Wait()
	assert.Equal(t, uint64(2), totalEventCounter)
}

func subscribeTester(t *testing.T) {
	dataChan := eventCreator.Subscribe(context.TODO(), core.WatchDirStart)
	wgSubscriberReady.Done()
	data := <-dataChan
	assert.Equal(t, WatchDirStartData, string(data.Data))
	assert.Equal(t, core.WatchDirStart, data.Type)
	atomic.AddUint64(&totalEventCounter, 1)
	wgDone.Done()
}
