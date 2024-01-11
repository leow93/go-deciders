package store

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type streams map[string][]Message
type MemoryStore struct {
	globalPosition int64
	streams        streams
}

func NewMemoryStore() *MemoryStore {
	store := MemoryStore{
		globalPosition: int64(0),
		streams:        map[string][]Message{},
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

func (store *MemoryStore) readStream(stream string, fromPosition int64) (error, []Message) {
	events, exists := store.streams[stream]
	if !exists {
		var result []Message
		return nil, result
	}
	return nil, filter(events, func(message Message) bool {
		return message.Position > fromPosition
	})
}

func currentPosition(events []Message) int64 {
	return int64(len(events))
}

func (store *MemoryStore) appendEvents(events []StreamMessage, stream string, expectedPosition int64) error {
	err, currentEvents := store.readStream(stream, 0)
	if err != nil {
		return err
	}
	position := currentPosition(currentEvents)
	if expectedPosition != position {
		return errors.New(fmt.Sprintf("Expected %d, got %d", expectedPosition, position))
	}
	var result = currentEvents
	for _, ev := range events {
		event := Message{
			Id:             uuid.NewString(),
			Stream:         stream,
			Data:           ev.Data,
			Meta:           ev.Meta,
			Position:       position + int64(1),
			GlobalPosition: store.globalPosition + int64(1),
		}
		result = append(result, event)
	}
	store.streams[stream] = result
	return nil
}
