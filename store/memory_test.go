package store

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)
import "github.com/google/uuid"

func streamName() string {
	return "Test-" + uuid.NewString()
}

var store = NewMemoryStore()

func TestMemoryStore(t *testing.T) {
	t.Run("appending events to a stream", func(t *testing.T) {
		stream := streamName()
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

	t.Run("concurrently appending events to a stream", func(t *testing.T) {
		stream := streamName()
		event := Message{
			Data: []byte("data"),
			Meta: []byte("meta"),
		}
		errors := make(chan error, 2)
		done := make(chan bool, 2)

		for i := 0; i < 2; i++ {
			go func() {
				var events []Message
				events = append(events, event)
				r := rand.Float32()
				time.Sleep(time.Millisecond * time.Duration(r*1000))
				err := store.AppendEvents(events, stream, int64(0))
				if err != nil {
					errors <- err
				}
				done <- true
			}()
		}
		<-done
		<-done

		if len(errors) != 1 {
			t.Errorf("Expected one error, got %d", len(errors))
		}
	})

	t.Run("reading events from a stream", func(t *testing.T) {
		stream := streamName()
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
		err, events = store.ReadStream(stream, int64(0))
		if err != nil {
			panic(err)
		}
		if len(events) != 1 {
			t.Errorf("Expected one event, got %d", len(events))
		}
	})
	t.Run("reading events from a non-existent stream", func(t *testing.T) {
		err, events := store.ReadStream("empty", int64(0))
		if err != nil {
			panic(err)
		}
		if len(events) != 0 {
			t.Errorf("Expected no events, got %d", len(events))
		}
	})

	t.Run("reading events from a position", func(t *testing.T) {
		stream := streamName()
		// append another event
		event := Message{
			Data: []byte("more-data"),
			Meta: []byte("more-meta"),
		}
		var events []Message
		for i := 0; i < 2; i++ {
			events = append(events, event)
		}
		err := store.AppendEvents(events, stream, int64(0))
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

	t.Run("subscribing to a stream", func(t *testing.T) {
		streamName := streamName()
		var events []Message
		for i := 0; i < 5; i++ {
			event := Message{
				Data: []byte(fmt.Sprint(i)),
				Meta: []byte("meta"),
			}
			events = append(events, event)
		}
		err := store.AppendEvents(events, streamName, int64(0))
		if err != nil {
			panic(err)
		}
		err, ch := store.SubscribeToStream(streamName, int64(0))

		if err != nil {
			panic(err)
		}

		timeoutCh := make(chan bool)

		go func() {
			time.Sleep(3 * time.Second)
			timeoutCh <- true
		}()

		for i := 0; i < 5; i++ {
			select {
			case <-timeoutCh:
				t.Errorf("Timeout")
				return
			case ev := <-ch:
				if string(ev.Data) != fmt.Sprint(i) {
					t.Errorf("Expected %d, got %s", i, string(ev.Data))
				}
			}
		}
	})

	t.Run("subscribing to a stream from a position", func(t *testing.T) {
		streamName := streamName()
		var events []Message
		for i := 0; i < 5; i++ {
			event := Message{
				Data: []byte(fmt.Sprint(i)),
				Meta: []byte("meta"),
			}
			events = append(events, event)
		}
		err := store.AppendEvents(events, streamName, int64(0))
		if err != nil {
			panic(err)
		}
		err, ch := store.SubscribeToStream(streamName, int64(2))

		if err != nil {
			panic(err)
		}

		timeoutCh := make(chan bool)

		go func() {
			time.Sleep(1 * time.Second)
			timeoutCh <- true
		}()

		for i := 2; i < 5; i++ {
			select {
			case <-timeoutCh:
				t.Errorf("Timeout")
				return
			case ev := <-ch:
				if string(ev.Data) != fmt.Sprint(i) {
					t.Errorf("Expected %d, got %s", i, string(ev.Data))
				}
			}
		}
	})
}
