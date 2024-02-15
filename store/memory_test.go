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
		event := Message{
			Data: []byte("data"),
			Meta: []byte("meta"),
		}
		var events []Message
		events = append(events, event)
		err := store.AppendEvents(events, stream, int64(0))
		if err != nil {
			panic(err)
		}
	})

	t.Run("reading events from a stream", func(t *testing.T) {
		err, events := store.ReadStream(stream, int64(0))
		if err != nil {
			panic(err)
		}
		if len(events) != 1 {
			t.Errorf("Expected one event, got %d", len(events))
		}
	})
	t.Run("reading events from a non-existed stream", func(t *testing.T) {
		err, events := store.ReadStream("empty", int64(0))
		if err != nil {
			panic(err)
		}
		if len(events) != 0 {
			t.Errorf("Expected no events, got %d", len(events))
		}
	})

	t.Run("reading events from a position", func(t *testing.T) {
		// append another event
		event := Message{
			Data: []byte("more-data"),
			Meta: []byte("more-meta"),
		}
		var events []Message
		events = append(events, event)
		err := store.AppendEvents(events, stream, int64(1))
		if err != nil {
			panic(err)
		}
		err, result := store.ReadStream(stream, int64(1))
		if err != nil {
			panic(err)
		}
		if len(result) != 1 {
			t.Errorf("Expected one event, got %d", len(result))
		}
	})
}
