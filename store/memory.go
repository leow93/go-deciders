package store

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type streams map[string][]Message
type categories map[string][]Message
type MemoryStore struct {
	globalPosition int64
	streams        streams
	categories     categories
	mutexes        sync.Map
}

func NewMemoryStore() *MemoryStore {
	streams := make(map[string][]Message)
	categories := make(map[string][]Message)
	store := MemoryStore{
		globalPosition: int64(0),
		streams:        streams,
		categories:     categories,
	}
	return &store
}

func filter(msgs []Message, fn func(Message) bool) []Message {
	var result []Message
	for _, x := range msgs {
		if fn(x) {
			result = append(result, x)
		}
	}
	return result
}

func (store *MemoryStore) ReadStream(stream string, fromPosition int64) (error, []Message) {
	events := store.streams[stream][fromPosition:]
	return nil, events
}

func currentPosition(events []Message) int64 {
	return int64(len(events))
}

func categoryName(streamName string) string {
	return strings.Split(streamName, "-")[0]
}

func (store *MemoryStore) lockCategory(key string) func() {
	mutex, _ := store.mutexes.LoadOrStore(key, &sync.Mutex{})
	mutex.(*sync.Mutex).Lock()
	return func() {
		mutex.(*sync.Mutex).Unlock()
	}
}

func (store *MemoryStore) AppendEvents(events []Message, stream string, expectedPosition int64) error {
	category := categoryName(stream)
	unlock := store.lockCategory(category)
	defer unlock()
	err, currentEvents := store.ReadStream(stream, 0)
	if err != nil {
		return err
	}
	position := currentPosition(currentEvents)

	if expectedPosition != position {
		return errors.New(fmt.Sprintf("Expected %d, got %d", expectedPosition, position))
	}
	var result = currentEvents
	for i, ev := range events {
		event := Message{
			Id:             uuid.NewString(),
			Stream:         stream,
			Data:           ev.Data,
			Meta:           ev.Meta,
			Position:       position + int64(1+i),
			GlobalPosition: store.globalPosition + int64(1+i),
		}
		result = append(result, event)
	}
	store.streams[stream] = result
	// todo: fix
	store.categories[category] = append(store.categories[category], events...)
	return nil
}

func (store *MemoryStore) SubscribeToStream(stream string, fromPosition int64) (error, <-chan Message) {
	err, events := store.ReadStream(stream, fromPosition)
	if err != nil {
		return err, nil
	}
	ch := make(chan Message)
	go func() {
		for _, event := range events {
			ch <- event
		}
	}()
	return nil, ch
}

func (store *MemoryStore) SubscribeToCategory(category string, fromPosition int64) (error, <-chan Message) {
	return nil, make(chan Message)
}
