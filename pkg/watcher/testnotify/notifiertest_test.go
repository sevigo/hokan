package testnotify

import (
	"testing"

	"github.com/sevigo/notify/event"
	"github.com/stretchr/testify/assert"
)

func TestFakeWatcher_StartWatching(t *testing.T) {
	w := New()
	go w.StartWatching("/test", nil)
	e := <-w.Event()
	assert.Equal(t, "/test", e.Path)
	assert.Equal(t, event.FileAdded, e.Action)
}
