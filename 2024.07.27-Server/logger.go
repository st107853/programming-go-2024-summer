package main

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type EventType byte

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, title, artist, prise string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type Event struct {
	Sequence  uint64    // A unique record ID
	EventType EventType //The action token
	Key       string    // The key affected by this transaction
	Title     string    // The value of a PUT the transaction
	Artist    string    //
	Prise     string    //
}
