package store

type Message struct {
	Id             string
	Stream         string
	Data           []byte
	Meta           []byte
	Position       int64
	GlobalPosition int64
}

type EventStore interface {
	ReadStream(stream string, fromPosition int64) (error, []Message)
	AppendEvents(events []Message, stream string, expectedPosition int64) error
	SubscribeToStream(stream string, fromPosition int64) (error, <-chan Message)
	SubscribeToCategory(category string, fromPosition int64) (error, <-chan Message)
}
