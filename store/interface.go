package store

type Message struct {
	Id             string
	Stream         string
	Data           struct{}
	Meta           struct{}
	Position       int64
	GlobalPosition int64
}

type StreamMessage struct {
	Data struct{}
	Meta struct{}
}

type EventStore interface {
	readStream(stream string, fromPosition int64) (error, []Message)
	appendEvents(events []StreamMessage, stream string, expectedPosition int64) error
}
