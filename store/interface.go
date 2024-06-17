package store

type Message struct {
	Id             string
	Stream         string
	Data           []byte
	Meta           []byte
	Position       int64
	GlobalPosition int64
}
