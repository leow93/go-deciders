package store

import "testing"
import "github.com/google/uuid"

func streamName() string {
	return "Test-" + uuid.NewString()
}

var store = NewMemoryStore()

func TestMemoryStore(t *testing.T) {
	stream := streamName()
	t.Run("appending events to a stream", func(t *testing.T) {
		event := StreamMessage{
			Data: struct{}{},
			Meta: struct{}{},
		}
		var events []StreamMessage
		events = append(events, event)
		err := store.appendEvents(events, stream, int64(0))
		if err != nil {
			panic(err)
		}
	})

	t.Run("reading events from a stream", func(t *testing.T) {
		err, events := store.readStream(stream, int64(0))
		if err != nil {
			panic(err)
		}
		if len(events) != 1 {
			t.Errorf("Expected one event, got %d", len(events))
		}
	})
	t.Run("reading events from a non-existed stream", func(t *testing.T) {
		err, events := store.readStream("empty", int64(0))
		if err != nil {
			panic(err)
		}
		if len(events) != 0 {
			t.Errorf("Expected no events, got %d", len(events))
		}
	})

	t.Run("reading events from a position", func(t *testing.T) {
		// append another event
		event := StreamMessage{
			Data: struct{}{},
			Meta: struct{}{},
		}
		var events []StreamMessage
		events = append(events, event)
		err := store.appendEvents(events, stream, int64(1))
		if err != nil {
			panic(err)
		}
		err, result := store.readStream(stream, int64(1))
		if err != nil {
			panic(err)
		}
		if len(result) != 1 {
			t.Errorf("Expected one event, got %d", len(result))
		}
	})
}
